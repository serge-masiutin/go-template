//go:build integration

package httpapp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/testdb"
)

type sessionCommitPauseKey struct{}

type pausedSessionStore struct {
	*pgxstore.PostgresStore
	loaded  chan struct{}
	release chan struct{}
}

func (s *pausedSessionStore) CommitCtx(ctx context.Context, token string, b []byte, expiry time.Time) error {
	if ctx.Value(sessionCommitPauseKey{}) == true {
		close(s.loaded)
		select {
		case <-s.release:
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
			return fmt.Errorf("session commit barrier timed out")
		}
	}
	return s.PostgresStore.CommitCtx(ctx, token, b, expiry)
}

func TestLogoutCannotBeUndoneByConcurrentSessionCommit(t *testing.T) {
	t.Chdir("../..")
	for _, scenario := range []struct {
		name         string
		idle         time.Duration
		method, path string
		status       int
	}{
		{"idle refresh", 2 * time.Hour, "GET", "/", http.StatusOK},
		{"flash write without idle timeout", 0, "POST", "/notes", http.StatusSeeOther},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			testLogoutWithConcurrentCommit(t, scenario.idle, scenario.method, scenario.path, scenario.status)
		})
	}
}

func testLogoutWithConcurrentCommit(t *testing.T, idle time.Duration, method, path string, slowStatus int) {
	t.Helper()
	db := testdb.New(t)
	users, err := accounts.New(db)
	if err != nil {
		t.Fatal(err)
	}
	user, err := users.Create(context.Background(), "review@example.test", "review-password", false)
	if err != nil {
		t.Fatal(err)
	}
	store := &pausedSessionStore{PostgresStore: pgxstore.NewWithCleanupInterval(db.Pool, 0), loaded: make(chan struct{}), release: make(chan struct{})}
	sessions := scs.New()
	sessions.Store = store
	sessions.Lifetime = 24 * time.Hour
	sessions.IdleTimeout = idle
	sessions.Cookie.Name = "go_template_session"
	ctx, err := sessions.Load(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := sessions.RenewToken(ctx); err != nil {
		t.Fatal(err)
	}
	if err := users.CreateLoginSession(ctx, user.ID, sessions.Token(ctx), "", sessions.Deadline(ctx)); err != nil {
		t.Fatal(err)
	}
	sessions.Put(ctx, "csrf", "review-csrf")
	token, _, err := sessions.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var release sync.Once
	defer release.Do(func() { close(store.release) })
	handler, err := New(config.Config{Environment: "test"}, db, sessions, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	// A second process-equivalent handler shares only PostgreSQL with the slow request.
	otherSessions := scs.New()
	otherSessions.Store = pgxstore.NewWithCleanupInterval(db.Pool, 0)
	otherSessions.Cookie.Name = sessions.Cookie.Name
	otherHandler, err := New(config.Config{Environment: "test"}, db, otherSessions, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	request := func(method, path string, pause bool) *httptest.ResponseRecorder {
		body := `{"body":"unauthorized"}`
		if pause {
			body = `{"body":""}`
		}
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		r.AddCookie(&http.Cookie{Name: sessions.Cookie.Name, Value: token})
		r.Header.Set("X-CSRF-Token", "review-csrf")
		if pause {
			r = r.WithContext(context.WithValue(r.Context(), sessionCommitPauseKey{}, true))
		}
		w := httptest.NewRecorder()
		if pause {
			handler.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
		return w
	}
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() { done <- request(method, path, true) }()
	select {
	case <-store.loaded:
	case <-time.After(10 * time.Second):
		release.Do(func() { close(store.release) })
		t.Fatal("authenticated request did not reach session commit")
	}
	logout := request("POST", "/logout", false)
	before := request("GET", "/", false)
	release.Do(func() { close(store.release) })
	var slow *httptest.ResponseRecorder
	select {
	case slow = <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("slow request did not finish")
	}
	after := request("GET", "/", false)
	if logout.Code != 303 || before.Code != 303 || slow.Code != slowStatus || after.Code != 303 {
		t.Fatalf("unexpected history: logout=%d old cookie before release=%d slow=%d old cookie after release=%d", logout.Code, before.Code, slow.Code, after.Code)
	}
	mutation := request("POST", "/notes", false)
	if mutation.Code != 303 && mutation.Code != 403 {
		t.Fatalf("revoked cookie mutation: %d", mutation.Code)
	}
	var count int
	if err := db.QueryRow(context.Background(), "SELECT count(*) FROM notes").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("revoked session created a note")
	}
}
