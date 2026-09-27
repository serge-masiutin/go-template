//go:build integration

package background

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/assistant"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/notemail"
	"github.com/serge-masiutin/go-template/internal/testdb"
)

func TestWorkerConsumesFailuresAndReconcilesKilledWork(t *testing.T) {
	db := testdb.New(t)
	ctx := context.Background()
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	queue, err := Producer(db, logger)
	if err != nil {
		t.Fatal(err)
	}
	users, err := accounts.New(db)
	if err != nil {
		t.Fatal(err)
	}
	user, err := users.Create(ctx, "worker@example.test", "integration-password", false)
	if err != nil {
		t.Fatal(err)
	}
	emails := notemail.New(db, queue, users)
	runs := assistant.New(db, queue)
	if err := emails.Request(ctx, user.ID); err != nil {
		t.Fatal(err)
	}
	if err := runs.Request(ctx, user.ID, "private-question-marker", "fixture/model"); err != nil {
		t.Fatal(err)
	}
	worker, err := New(ctx, config.Config{WorkerConcurrency: 2}, db, logger)
	if err != nil {
		t.Fatal(err)
	}
	if err := worker.Client.Start(ctx); err != nil {
		t.Fatal(err)
	}
	stopped := false
	defer func() {
		if !stopped {
			shutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := worker.Stop(shutdown); err != nil {
				t.Error(err)
			}
		}
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		mailRows, err := emails.Recent(ctx, user.ID)
		if err != nil {
			t.Fatal(err)
		}
		aiRows, err := runs.Recent(ctx, user.ID)
		if err != nil {
			t.Fatal(err)
		}
		if mailRows[0].State == "failed" && aiRows[0].State == "failed" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("worker did not finish disabled-provider jobs")
		}
		time.Sleep(20 * time.Millisecond)
	}
	shutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := worker.Stop(shutdown); err != nil {
		t.Fatal(err)
	}
	stopped = true
	// Simulate a process killed after claiming an effect; River has exhausted its single attempt.
	if err := emails.Request(ctx, user.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, "UPDATE note_emails SET state='sending' WHERE state='queued'; UPDATE river_job SET state='discarded', finalized_at=now() WHERE kind='note_email_v1' AND state='available'"); err != nil {
		t.Fatal(err)
	}
	reconcile := &reconcileWorker{db: db, logger: logger}
	if err := reconcile.Work(ctx, &river.Job[reconcileArgs]{JobRow: &rivertype.JobRow{ID: 1, Kind: "reconcile_work_v1"}}); err != nil {
		t.Fatal(err)
	}
	rows, err := emails.Recent(ctx, user.ID)
	if err != nil || rows[0].State != "failed" {
		t.Fatalf("killed delivery remains pending: %v %v", rows, err)
	}

	if strings.Contains(logs.String(), "private-question-marker") {
		t.Fatal("question leaked into logs")
	}
	var persisted string
	if err := db.QueryRow(ctx, "SELECT errors::text FROM river_job WHERE kind='notes_assistant_v1'").Scan(&persisted); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(persisted, "private-question-marker") || !strings.Contains(persisted, "inspect the application log") {
		t.Fatalf("unsafe persisted job error: %s", persisted)
	}
}

func TestProviderErrorsArePreservedButNotPersistedVerbatim(t *testing.T) {
	var logs bytes.Buffer
	cause := errors.New("private-provider-payload")
	failed := false
	err := perform(context.Background(), slog.New(slog.NewJSONHandler(&logs, nil)), 42, "notes_assistant_v1", func() error { return cause }, func(context.Context) error { failed = true; return nil })
	if !failed || !errors.Is(err, cause) || strings.Contains(err.Error(), "private-provider-payload") || strings.Contains(logs.String(), "private-provider-payload") {
		t.Fatal("provider error contract violated")
	}
}
