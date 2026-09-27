//go:build integration

package database_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/serge-masiutin/go-template/internal/database"
	"github.com/serge-masiutin/go-template/internal/testdb"
)

func TestAdoptLegacyWithoutLosingData(t *testing.T) {
	db := testdb.Open(t)
	ctx := context.Background()
	original, err := os.ReadFile("bootstrap/001_accounts.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, string(original)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, "CREATE TABLE schema_migrations(name text PRIMARY KEY); INSERT INTO schema_migrations VALUES ('migrations/001_accounts.sql'); INSERT INTO users(email,password_hash) VALUES ('legacy@example.test',decode('00','hex'))"); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := database.Migrate(ctx, db); err != nil {
			t.Fatal(err)
		}
	}
	var email string
	if err := db.QueryRow(ctx, "SELECT email FROM users").Scan(&email); err != nil || email != "legacy@example.test" {
		t.Fatalf("legacy account lost: %q %v", email, err)
	}
	var tables int
	if err := db.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema=current_schema() AND table_name IN ('river_job','assistant_runs','note_emails')").Scan(&tables); err != nil || tables != 3 {
		t.Fatalf("new tables: %d, %v", tables, err)
	}
}
func TestUnknownLegacyHistoryFails(t *testing.T) {
	db := testdb.Open(t)
	if _, err := db.Exec(context.Background(), "CREATE TABLE schema_migrations(name text PRIMARY KEY); INSERT INTO schema_migrations VALUES ('custom.sql')"); err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(context.Background(), db); err == nil {
		t.Fatal("unknown legacy history was silently adopted")
	}
}

func TestConcurrentMigratorsShareDeploymentLock(t *testing.T) {
	db := testdb.Open(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	results := make(chan error, 2)
	for range 2 {
		go func() { results <- database.Migrate(ctx, db) }()
	}
	for range 2 {
		if err := <-results; err != nil {
			t.Error(err)
		}
	}
	if t.Failed() {
		return
	}
	var applied int
	if err := db.QueryRow(ctx, "SELECT count(*) FROM goose_db_version WHERE version_id IN (1,2) AND is_applied").Scan(&applied); err != nil || applied != 2 {
		t.Fatalf("migration history inconsistent: count=%d err=%v", applied, err)
	}
}
