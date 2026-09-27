package assistant

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/riverqueue/river"
	"github.com/serge-masiutin/go-template/internal/database"
	"gorm.io/gorm"
)

var ErrPending = errors.New("an assistant request is already pending")
var ErrUncertain = errors.New("AI invocation already started; it will not be billed again automatically")

type Run struct {
	ID        int64     `json:"id,string"`
	Question  string    `json:"question"`
	State     string    `json:"state"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"createdAt"`
}

type record struct {
	ID            int64      `gorm:"primaryKey"`
	UserID        int64      `json:"-"`
	JobID         *int64     `json:"-"`
	Question      string     `json:"question"`
	State         string     `json:"state"`
	Answer        string     `json:"answer"`
	Model         string     `json:"-"`
	PromptVersion string     `json:"-"`
	CreatedAt     time.Time  `gorm:"autoCreateTime:false"`
	FinishedAt    *time.Time `json:"-"`
}

func (record) TableName() string { return "assistant_runs" }

type Args struct {
	RunID int64 `json:"run_id"`
}

func (Args) Kind() string { return "notes_assistant_v1" }

type Service struct {
	db    *database.DB
	queue *river.Client[*sql.Tx]
}

func New(db *database.DB, queue *river.Client[*sql.Tx]) *Service {
	return &Service{db: db, queue: queue}
}

func (s *Service) Request(ctx context.Context, userID int64, question, model string) error {
	question = strings.TrimSpace(question)
	if err := ValidateQuestion(question); err != nil {
		return err
	}
	err := s.db.ORM.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		run := record{UserID: userID, Question: question, Model: model, PromptVersion: PromptVersion, State: "queued"}
		if err := gorm.G[record](tx).Omit("CreatedAt", "FinishedAt").Create(ctx, &run); err != nil {
			return err
		}
		job, err := s.queue.InsertTx(ctx, tx.Statement.ConnPool.(*sql.Tx), Args{RunID: run.ID}, &river.InsertOpts{Queue: "ai", MaxAttempts: 1})
		if err != nil {
			return err
		}
		_, err = gorm.G[record](tx).Where("id = ?", run.ID).Update(ctx, "job_id", job.Job.ID)
		return err
	})
	var conflict *pgconn.PgError
	if errors.As(err, &conflict) && conflict.ConstraintName == "assistant_runs_one_active_idx" {
		return ErrPending
	}
	return err
}

func (s *Service) Recent(ctx context.Context, userID int64) ([]Run, error) {
	items, err := gorm.G[record](s.db.ORM).Select("id, question, state, answer, created_at").Where("user_id = ?", userID).Order("id DESC").Limit(5).Find(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]Run, len(items))
	for i, item := range items {
		result[i] = Run{ID: item.ID, Question: item.Question, State: item.State, Answer: item.Answer, CreatedAt: item.CreatedAt}
	}
	return result, nil
}

func (s *Service) Execute(ctx context.Context, id int64, agent *Agent) error {
	run, err := gorm.G[record](s.db.ORM).Where("id = ?", id).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} // Deleted accounts cannot leave an executable request.
	if err != nil {
		return err
	}
	if run.State == "completed" {
		return nil
	}
	if run.State != "queued" {
		return ErrUncertain
	}
	if run.Model != agent.cfg.Model || run.PromptVersion != PromptVersion {
		return errors.New("queued AI configuration differs from the worker; request a new run")
	}
	changed, err := gorm.G[record](s.db.ORM).Where("id = ? AND state = ?", id, "queued").Update(ctx, "state", "running")
	if err != nil {
		return err
	}
	if changed != 1 {
		return ErrUncertain
	}
	output, err := agent.Answer(ctx, run.UserID, run.Question)
	if err != nil {
		return err
	}
	result := s.db.ORM.WithContext(ctx).Model(&record{}).Where("id = ? AND state = ?", id, "running").Updates(map[string]any{"state": "completed", "answer": output.Answer, "finished_at": time.Now().UTC()})
	return result.Error
}

func (s *Service) Fail(ctx context.Context, id int64) error {
	return s.db.ORM.WithContext(ctx).Model(&record{}).Where("id = ? AND state IN ?", id, []string{"queued", "running"}).Updates(map[string]any{"state": "failed", "finished_at": time.Now().UTC()}).Error
}
