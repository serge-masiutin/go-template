package background

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/assistant"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/database"
	"github.com/serge-masiutin/go-template/internal/diagnostics"
	"github.com/serge-masiutin/go-template/internal/mailing"
	"github.com/serge-masiutin/go-template/internal/notemail"
	"github.com/serge-masiutin/go-template/internal/notes"
)

func Producer(db *database.DB, logger *slog.Logger) (*river.Client[*sql.Tx], error) {
	return river.NewClient(riverdatabasesql.New(db.SQL), &river.Config{Logger: diagnostics.LibraryLogger(logger)})
}

type Runtime struct {
	Client   *river.Client[*sql.Tx]
	listener *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config, db *database.DB, logger *slog.Logger) (*Runtime, error) {
	users, err := accounts.New(db)
	if err != nil {
		return nil, err
	}
	emails := notemail.New(db, nil, users)
	runs := assistant.New(db, nil)
	var sender notemail.Sender
	if cfg.Mail.Enabled {
		sender, err = mailing.New(cfg.Mail)
		if err != nil {
			return nil, err
		}
	}
	var agent *assistant.Agent
	if cfg.AI.Enabled {
		agent, err = assistant.NewAgent(ctx, cfg.AI, notes.New(db))
		if err != nil {
			return nil, err
		}
	}
	workers := river.NewWorkers()
	river.AddWorker(workers, &mailWorker{service: emails, sender: sender, logger: logger})
	river.AddWorker(workers, &assistantWorker{service: runs, agent: agent, logger: logger})
	river.AddWorker(workers, &reconcileWorker{db: db, logger: logger})
	listenerCfg := db.Pool.Config().Copy()
	listenerCfg.MaxConns, listenerCfg.MinConns = 1, 0
	listener, err := pgxpool.NewWithConfig(ctx, listenerCfg)
	if err != nil {
		return nil, err
	}
	client, err := river.NewClient(riverdatabasesql.NewWithPgxListener(db.SQL, listener), &river.Config{
		Logger: diagnostics.LibraryLogger(logger), Workers: workers,
		Queues:     map[string]river.QueueConfig{"mail": {MaxWorkers: cfg.WorkerConcurrency}, "ai": {MaxWorkers: 1}, "default": {MaxWorkers: 1}},
		JobTimeout: 6 * time.Minute, RescueStuckJobsAfter: 7 * time.Minute,
		PeriodicJobs: []*river.PeriodicJob{river.NewPeriodicJob(river.PeriodicInterval(time.Minute), func() (river.JobArgs, *river.InsertOpts) {
			return reconcileArgs{}, &river.InsertOpts{MaxAttempts: 3, UniqueOpts: river.UniqueOpts{ByPeriod: time.Minute}}
		}, &river.PeriodicJobOpts{RunOnStart: true})},
	})
	if err != nil {
		listener.Close()
		return nil, err
	}
	return &Runtime{Client: client, listener: listener}, nil
}

// Stop drains active jobs before closing the dedicated notification connection.
func (r *Runtime) Stop(ctx context.Context) error {
	err := r.Client.Stop(ctx)
	if err != nil {
		cancel, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		err = errors.Join(err, r.Client.StopAndCancel(cancel))
	}
	r.listener.Close()
	return err
}

type mailWorker struct {
	river.WorkerDefaults[notemail.Args]
	service *notemail.Service
	sender  notemail.Sender
	logger  *slog.Logger
}

func (w *mailWorker) Work(ctx context.Context, job *river.Job[notemail.Args]) error {
	return perform(ctx, w.logger, job.ID, job.Kind, func() error {
		if w.sender == nil {
			return errors.New("mail is disabled in the worker")
		}
		return w.service.Deliver(ctx, job.Args.DeliveryID, w.sender)
	}, func(ctx context.Context) error { return w.service.Fail(ctx, job.Args.DeliveryID) })
}

type assistantWorker struct {
	river.WorkerDefaults[assistant.Args]
	service *assistant.Service
	agent   *assistant.Agent
	logger  *slog.Logger
}

func (w *assistantWorker) Work(ctx context.Context, job *river.Job[assistant.Args]) error {
	return perform(ctx, w.logger, job.ID, job.Kind, func() error {
		if w.agent == nil {
			return errors.New("AI is disabled in the worker")
		}
		return w.service.Execute(ctx, job.Args.RunID, w.agent)
	}, func(ctx context.Context) error { return w.service.Fail(ctx, job.Args.RunID) })
}

type safeJobError struct{ cause error }

func (e safeJobError) Error() string { return "job failed; inspect the application log using job_id" }
func (e safeJobError) Unwrap() error { return e.cause }

func perform(ctx context.Context, logger *slog.Logger, id int64, kind string, work func() error, fail func(context.Context) error) (err error) {
	started := time.Now()
	defer func() {
		if recovered := recover(); recovered != nil {
			logger.ErrorContext(ctx, "worker panic", "job_id", id, "stack", string(debug.Stack()))
			err = errors.New("worker panic")
		}
		if err != nil {
			cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			defer cancel()
			err = errors.Join(err, fail(cleanup))
			logger.With("job_id", id, "kind", kind).ErrorContext(ctx, "job failed", diagnostics.Attributes(err)...)
			err = safeJobError{cause: err} // River persists returned errors; never put provider responses in its tables.
			return
		}
		logger.InfoContext(ctx, "job completed", "job_id", id, "kind", kind, "duration_ms", time.Since(started).Milliseconds())
	}()
	return work()
}

type reconcileArgs struct{}

func (reconcileArgs) Kind() string { return "reconcile_work_v1" }

type reconcileWorker struct {
	river.WorkerDefaults[reconcileArgs]
	db     *database.DB
	logger *slog.Logger
}

func (w *reconcileWorker) Work(ctx context.Context, job *river.Job[reconcileArgs]) error {
	// A process kill can bypass failure cleanup. Terminal or removed River jobs must not leave the UI pending forever.
	return perform(ctx, w.logger, job.ID, job.Kind, func() error {
		for _, table := range []string{"note_emails", "assistant_runs"} {
			_, err := w.db.Exec(ctx, "UPDATE "+table+" AS work SET state='failed', finished_at=now() WHERE work.state IN ('queued','running','sending') AND NOT EXISTS (SELECT 1 FROM river_job WHERE id=work.job_id AND state IN ('available','pending','retryable','running','scheduled'))")
			if err != nil {
				return err
			}
		}
		return nil
	}, func(context.Context) error { return nil })
}
