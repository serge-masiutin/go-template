package notemail

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"strings"
	"text/template"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/riverqueue/river"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/database"
	"github.com/serge-masiutin/go-template/internal/notes"
	"gorm.io/gorm"
)

var ErrPending = errors.New("a note email is already pending")
var ErrUncertain = errors.New("delivery was started previously; inspect SMTP delivery before requesting another email")

//go:embed notes.tmpl
var emailTemplate string
var bodyTemplate = template.Must(template.New("notes").Parse(emailTemplate))

type Sender interface {
	Send(context.Context, string, string, string) error
}

type Delivery struct {
	ID        int64     `json:"id,string"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"createdAt"`
}

type record struct {
	ID         int64      `gorm:"primaryKey"`
	UserID     int64      `json:"-"`
	JobID      *int64     `json:"-"`
	State      string     `json:"state"`
	CreatedAt  time.Time  `gorm:"autoCreateTime:false"`
	FinishedAt *time.Time `json:"-"`
}

func (record) TableName() string { return "note_emails" }

type Args struct {
	DeliveryID int64 `json:"delivery_id"`
}

func (Args) Kind() string { return "note_email_v1" }

type Service struct {
	db    *database.DB
	queue *river.Client[*sql.Tx]
	users *accounts.Store
	notes *notes.Store
}

func New(db *database.DB, queue *river.Client[*sql.Tx], users *accounts.Store) *Service {
	return &Service{db: db, queue: queue, users: users, notes: notes.New(db)}
}

// Request records delivery intent and queues it in the same transaction.
func (s *Service) Request(ctx context.Context, userID int64) error {
	err := s.db.ORM.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		delivery := record{UserID: userID, State: "queued"}
		if err := gorm.G[record](tx).Omit("CreatedAt", "FinishedAt").Create(ctx, &delivery); err != nil {
			return err
		}
		job, err := s.queue.InsertTx(ctx, tx.Statement.ConnPool.(*sql.Tx), Args{DeliveryID: delivery.ID}, &river.InsertOpts{Queue: "mail", MaxAttempts: 1})
		if err != nil {
			return err
		}
		_, err = gorm.G[record](tx).Where("id = ?", delivery.ID).Update(ctx, "job_id", job.Job.ID)
		return err
	})
	var conflict *pgconn.PgError
	if errors.As(err, &conflict) && conflict.ConstraintName == "note_emails_one_active_idx" {
		return ErrPending
	}
	return err
}

func (s *Service) Recent(ctx context.Context, userID int64) ([]Delivery, error) {
	items, err := gorm.G[record](s.db.ORM).Select("id, state, created_at").Where("user_id = ?", userID).Order("id DESC").Limit(5).Find(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]Delivery, len(items))
	for i, item := range items {
		result[i] = Delivery{ID: item.ID, State: item.State, CreatedAt: item.CreatedAt}
	}
	return result, nil
}

func (s *Service) Deliver(ctx context.Context, id int64, sender Sender) error {
	delivery, err := gorm.G[record](s.db.ORM).Where("id = ?", id).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} // Account deletion cancels its deliveries by FK cascade.
	if err != nil {
		return err
	}
	if delivery.State == "sent" {
		return nil
	}
	if delivery.State != "queued" {
		return ErrUncertain
	}
	user, err := s.users.Find(ctx, delivery.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	items, err := s.notes.List(ctx, user.ID)
	if err != nil {
		return err
	}
	var body strings.Builder
	if err := bodyTemplate.Execute(&body, items); err != nil {
		return err
	}
	changed, err := gorm.G[record](s.db.ORM).Where("id = ? AND state = ?", id, "queued").Update(ctx, "state", "sending")
	if err != nil {
		return err
	}
	if changed != 1 {
		return ErrUncertain
	}
	// Claim before SMTP. A crash after acceptance must not cause an automatic duplicate.
	sendErr := sender.Send(ctx, user.Email, "Your notes", body.String())
	state := "sent"
	if sendErr != nil {
		state = "failed"
	}
	cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	result := s.db.ORM.WithContext(cleanup).Model(&record{}).Where("id = ? AND state = ?", id, "sending").Updates(map[string]any{"state": state, "finished_at": time.Now().UTC()})
	return errors.Join(sendErr, result.Error)
}

func (s *Service) Fail(ctx context.Context, id int64) error {
	return s.db.ORM.WithContext(ctx).Model(&record{}).Where("id = ? AND state IN ?", id, []string{"queued", "sending"}).Updates(map[string]any{"state": "failed", "finished_at": time.Now().UTC()}).Error
}
