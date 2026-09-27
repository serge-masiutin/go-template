//go:build integration

package notemail

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/notes"
	"github.com/serge-masiutin/go-template/internal/testdb"
)

type captureSender struct {
	calls           int
	recipient, body string
	err             error
}

func (s *captureSender) Send(_ context.Context, recipient, subject, body string) error {
	s.calls++
	s.recipient, s.body = recipient, body
	return s.err
}

func TestAtomicQueueAndOwnerScopedDelivery(t *testing.T) {
	db := testdb.New(t)
	ctx := context.Background()
	users, err := accounts.New(db)
	if err != nil {
		t.Fatal(err)
	}
	owner, err := users.Create(ctx, "owner@example.test", "integration-password", false)
	if err != nil {
		t.Fatal(err)
	}
	other, err := users.Create(ctx, "other@example.test", "integration-password", false)
	if err != nil {
		t.Fatal(err)
	}
	store := notes.New(db)
	if err := store.Create(ctx, owner.ID, "owner note"); err != nil {
		t.Fatal(err)
	}
	if err := store.Create(ctx, other.ID, "foreign secret"); err != nil {
		t.Fatal(err)
	}
	queue, err := river.NewClient(riverdatabasesql.New(db.SQL), &river.Config{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})
	if err != nil {
		t.Fatal(err)
	}
	service := New(db, queue, users)
	// Force the queue boundary to fail after the delivery row is inserted.
	if _, err := db.Exec(ctx, "ALTER TABLE river_job ADD CONSTRAINT reject_email CHECK (kind <> 'note_email_v1')"); err != nil {
		t.Fatal(err)
	}
	if err := service.Request(ctx, owner.ID); err == nil {
		t.Fatal("enqueue unexpectedly succeeded")
	}
	deliveries, err := service.Recent(ctx, owner.ID)
	if err != nil || len(deliveries) != 0 {
		t.Fatalf("delivery survived rollback: %v %v", deliveries, err)
	}
	if _, err := db.Exec(ctx, "ALTER TABLE river_job DROP CONSTRAINT reject_email"); err != nil {
		t.Fatal(err)
	}
	var accepted atomic.Int32
	var wg sync.WaitGroup
	for range 6 {
		wg.Go(func() {
			err := service.Request(ctx, owner.ID)
			if err == nil {
				accepted.Add(1)
			} else if !errors.Is(err, ErrPending) {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if accepted.Load() != 1 {
		t.Fatalf("accepted %d concurrent requests", accepted.Load())
	}
	deliveries, err = service.Recent(ctx, owner.ID)
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("deliveries: %v", err)
	}
	sender := &captureSender{}
	for range 2 {
		if err := service.Deliver(ctx, deliveries[0].ID, sender); err != nil {
			t.Fatal(err)
		}
	}
	if sender.calls != 1 || sender.recipient != owner.Email || !strings.Contains(sender.body, "owner note") || strings.Contains(sender.body, "foreign secret") {
		t.Fatalf("incorrect delivery scope or duplicate: %+v", sender)
	}
	if err := service.Request(ctx, owner.ID); err != nil {
		t.Fatal(err)
	}
	deliveries, err = service.Recent(ctx, owner.ID)
	if err != nil {
		t.Fatal(err)
	}
	sender.err = errors.New("SMTP acceptance uncertain")
	if err := service.Deliver(ctx, deliveries[0].ID, sender); err == nil {
		t.Fatal("SMTP failure hidden")
	}
	if err := service.Deliver(ctx, deliveries[0].ID, sender); !errors.Is(err, ErrUncertain) {
		t.Fatalf("ambiguous delivery repeated: %v", err)
	}
	if sender.calls != 2 {
		t.Fatal("SMTP was retried")
	}
}
