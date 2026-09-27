//go:build integration

package httpapp

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/database"
)

func TestAccountAndOwnershipFlow(t *testing.T) {
	t.Chdir("../..")
	connectionString := os.Getenv("TEST_DATABASE_URL")
	if connectionString == "" {
		t.Fatal("TEST_DATABASE_URL is required for integration tests")
	}
	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(cfg.ConnConfig.Database, "_test") {
		t.Fatal("integration database name must end in _test")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	schema := fmt.Sprintf("test_%d", time.Now().UnixNano())
	if _, err := pool.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer pool.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	testPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer testPool.Close()
	if err := database.Migrate(ctx, testPool); err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(ctx, testPool); err != nil {
		t.Fatalf("repeat migration: %v", err)
	}
	users, err := accounts.New(testPool)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := users.Create(ctx, strings.Repeat("a", 255)+"@example.test", "integration-password", false); err == nil {
		t.Fatal("created an account whose address exceeds the login limit")
	}
	owner, err := users.Create(ctx, "owner@example.test", "integration-password", true)
	if err != nil {
		t.Fatal(err)
	}
	other, err := users.Create(ctx, "other@example.test", "integration-password", false)
	if err != nil {
		t.Fatal(err)
	}
	var foreignNote int64
	if err := testPool.QueryRow(ctx, "INSERT INTO notes(user_id,body) VALUES($1,'private to another user') RETURNING id", other.ID).Scan(&foreignNote); err != nil {
		t.Fatal(err)
	}
	sessionStore := pgxstore.NewWithCleanupInterval(testPool, 0)
	sessions := scs.New()
	sessions.Store = sessionStore
	handler, err := New(config.Config{}, testPool, sessions, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	manifest, err := os.ReadFile("web/build/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	version := fmt.Sprintf("%x", md5.Sum(manifest))
	request := func(method, path, body, csrf string, extra ...map[string]string) (int, http.Header, []byte) {
		t.Helper()
		r, err := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("X-Inertia", "true")
		r.Header.Set("X-Inertia-Version", version)
		r.Header.Set("Content-Type", "application/json")
		if csrf != "" {
			r.Header.Set("X-CSRF-Token", csrf)
		}
		for _, headers := range extra {
			for name, value := range headers {
				r.Header.Set(name, value)
			}
		}
		response, err := client.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		content, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		return response.StatusCode, response.Header, content
	}
	page := func(path string) map[string]json.RawMessage {
		t.Helper()
		status, _, body := request("GET", path, "", "")
		if status != 200 {
			t.Fatalf("GET %s: %d %s", path, status, body)
		}
		var result struct {
			Props map[string]json.RawMessage `json:"props"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatal(err)
		}
		return result.Props
	}
	csrfFrom := func(props map[string]json.RawMessage) string {
		t.Helper()
		var token string
		if err := json.Unmarshal(props["csrfToken"], &token); err != nil {
			t.Fatal(err)
		}
		if token == "" {
			t.Fatal("missing CSRF")
		}
		return token
	}
	status, _, _ := request("GET", "/", "", "")
	if status != 303 {
		t.Fatalf("anonymous home: %d", status)
	}
	csrf := csrfFrom(page("/login"))
	status, _, _ = request("POST", "/login", `{}`, csrf, map[string]string{"Sec-Fetch-Site": "cross-site"})
	if status != 403 {
		t.Fatalf("cross-site request accepted: %d", status)
	}
	status, staleHeaders, _ := request("GET", "/login", "", "", map[string]string{"X-Inertia-Version": "stale"})
	if status != 409 || staleHeaders.Get("X-Inertia-Location") == "" {
		t.Fatalf("stale asset version: %d", status)
	}
	status, _, partial := request("GET", "/login", "", "", map[string]string{"X-Inertia-Partial-Component": "Login", "X-Inertia-Partial-Data": "errors"})
	if status != 200 || !strings.Contains(string(partial), `"csrfToken"`) {
		t.Fatalf("CSRF omitted from partial response: %d %s", status, partial)
	}
	oldCSRF := csrf
	status, _, _ = request("POST", "/login", `{"email":"owner@example.test","password":"integration-password"}`, "")
	if status != 403 {
		t.Fatalf("CSRF bypass: %d", status)
	}
	status, _, _ = request("POST", "/login", `{"email":"owner@example.test","password":"wrong"}`, csrf)
	if status != 303 {
		t.Fatal(status)
	}
	if string(page("/login")["errors"]) == "{}" {
		t.Fatal("missing validation error after redirect")
	}
	if string(page("/login")["errors"]) != "{}" {
		t.Fatal("validation error was not consumed")
	}
	status, headers, _ := request("POST", "/login", `{"email":"owner@example.test","password":"integration-password"}`, csrf)
	if status != 303 || headers.Get("Location") != "/" {
		t.Fatalf("login: %d", status)
	}
	props := page("/")
	csrf = csrfFrom(props)
	if csrf == oldCSRF {
		t.Fatal("CSRF was not rotated at login")
	}
	status, _, _ = request("POST", "/notes", `{"body":"invalid old session token"}`, oldCSRF)
	if status != 403 {
		t.Fatalf("old CSRF still accepted: %d", status)
	}
	status, _, _ = request("POST", "/notes", `{"body":"must not be created"}`, csrf, map[string]string{"Precognition": "true"})
	if status != 400 {
		t.Fatalf("unsupported validation request accepted: %d", status)
	}
	if strings.Contains(string(page("/")["notes"]), "must not be created") {
		t.Fatal("validation-only request mutated state")
	}
	if strings.Contains(string(props["user"]), "password") || strings.Contains(string(props["notes"]), "private") {
		t.Fatal("private data leaked")
	}
	status, _, _ = request("POST", "/notes", `{"body":"","user_id":99}`, csrf)
	if status != 400 {
		t.Fatalf("unknown field accepted: %d", status)
	}
	status, _, _ = request("POST", "/notes", `{"body":""}`, csrf)
	if status != 303 {
		t.Fatal(status)
	}
	if string(page("/")["errors"]) == "{}" {
		t.Fatal("missing note error")
	}
	status, _, _ = request("POST", "/notes", `{"body":"my own note"}`, csrf)
	if status != 303 {
		t.Fatal(status)
	}
	if !strings.Contains(string(page("/")["notes"]), "my own note") {
		t.Fatal("note not persisted")
	}
	status, _, _ = request("DELETE", fmt.Sprintf("/notes/%d", foreignNote), "", csrf)
	if status != 404 {
		t.Fatalf("foreign delete: %d", status)
	}
	status, _, _ = request("GET", "/admin", "", "")
	if status != 200 {
		t.Fatal(status)
	}
	if _, err := testPool.Exec(ctx, "UPDATE users SET admin=false WHERE id=$1", owner.ID); err != nil {
		t.Fatal(err)
	}
	status, _, _ = request("GET", "/admin", "", "")
	if status != 403 {
		t.Fatalf("revoked admin: %d", status)
	}
	status, _, _ = request("POST", "/logout", "{}", csrf)
	if status != 303 {
		t.Fatal(status)
	}
	status, _, _ = request("GET", "/", "", "")
	if status != 303 {
		t.Fatalf("session survived logout: %d", status)
	}
}
