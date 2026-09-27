//go:build integration

package assistant

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/notes"
	"github.com/serge-masiutin/go-template/internal/testdb"
)

func TestRunTransactionIsolationAndDuplicateExecution(t *testing.T) {
	db := testdb.New(t)
	ctx := context.Background()
	users, err := accounts.New(db)
	if err != nil {
		t.Fatal(err)
	}
	owner, err := users.Create(ctx, "ai-owner@example.test", "integration-password", false)
	if err != nil {
		t.Fatal(err)
	}
	other, err := users.Create(ctx, "ai-other@example.test", "integration-password", false)
	if err != nil {
		t.Fatal(err)
	}
	reader := notes.New(db)
	if err := reader.Create(ctx, owner.ID, "authorized-note-marker"); err != nil {
		t.Fatal(err)
	}
	if err := reader.Create(ctx, other.ID, "foreign-note-marker"); err != nil {
		t.Fatal(err)
	}
	queue, err := river.NewClient(riverdatabasesql.New(db.SQL), &river.Config{})
	if err != nil {
		t.Fatal(err)
	}
	service := New(db, queue)
	if _, err := db.Exec(ctx, "ALTER TABLE river_job ADD CONSTRAINT reject_ai CHECK (kind <> 'notes_assistant_v1')"); err != nil {
		t.Fatal(err)
	}
	if err := service.Request(ctx, owner.ID, "Question", "fixture/model"); err == nil {
		t.Fatal("queue failure hidden")
	}
	runs, err := service.Recent(ctx, owner.ID)
	if err != nil || len(runs) != 0 {
		t.Fatal("run survived queue rollback")
	}
	if _, err := db.Exec(ctx, "ALTER TABLE river_job DROP CONSTRAINT reject_ai"); err != nil {
		t.Fatal(err)
	}
	if err := service.Request(ctx, owner.ID, "Question", "fixture/model"); err != nil {
		t.Fatal(err)
	}
	if err := service.Request(ctx, owner.ID, "Duplicate", "fixture/model"); !errors.Is(err, ErrPending) {
		t.Fatalf("duplicate run accepted: %v", err)
	}
	runs, err = service.Recent(ctx, owner.ID)
	if err != nil {
		t.Fatal(err)
	}
	g := genkit.Init(ctx)
	var calls atomic.Int32
	genkit.DefineModel(g, "fixture/model", &ai.ModelOptions{Supports: &ai.ModelSupports{Tools: true, Multiturn: true, SystemRole: true, Output: []string{"json"}}}, func(_ context.Context, request *ai.ModelRequest, _ ai.ModelStreamCallback) (*ai.ModelResponse, error) {
		if calls.Add(1) == 1 {
			return &ai.ModelResponse{Message: &ai.Message{Role: ai.RoleModel, Content: []*ai.Part{ai.NewToolRequestPart(&ai.ToolRequest{Name: "list_notes", Ref: "notes", Input: map[string]any{}})}}}, nil
		}
		serialized, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		if !strings.Contains(string(serialized), "authorized-note-marker") || strings.Contains(string(serialized), "foreign-note-marker") {
			t.Error("model received unauthorized notes or missed its own notes")
		}
		return &ai.ModelResponse{Message: ai.NewModelTextMessage(`{"answer":"A validated answer."}`)}, nil
	})
	agent := newAgent(g, config.AI{Model: "fixture/model", Timeout: time.Second, MaxTurns: 3, MaxOutputTokens: 128}, reader)
	for range 2 {
		if err := service.Execute(ctx, runs[0].ID, agent); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 2 {
		t.Fatal("completed operation was billed again")
	}
	runs, err = service.Recent(ctx, owner.ID)
	if err != nil || runs[0].Answer != "A validated answer." || runs[0].State != "completed" {
		t.Fatal("result not persisted")
	}
	if err := service.Request(ctx, other.ID, "Deleted account", "fixture/model"); err != nil {
		t.Fatal(err)
	}
	removed, err := service.Recent(ctx, other.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, "DELETE FROM users WHERE id=$1", other.ID); err != nil {
		t.Fatal(err)
	}
	if err := service.Execute(ctx, removed[0].ID, agent); err != nil || calls.Load() != 2 {
		t.Fatal("deleted account triggered AI generation")
	}
}
