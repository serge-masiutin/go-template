// Browserfixture isolates browser tests in a disposable schema of a dedicated test database.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/database"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (result error) {
	raw := os.Getenv("TEST_DATABASE_URL")
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || !strings.HasSuffix(parsed.Path, "_test") {
		return errors.New("TEST_DATABASE_URL must be a PostgreSQL URL with a database name ending in _test")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	setup, cancelSetup := context.WithTimeout(ctx, 30*time.Second)
	defer cancelSetup()
	pool, err := pgxpool.New(setup, raw)
	if err != nil {
		return errors.New("invalid test database configuration")
	}
	defer pool.Close()
	random := make([]byte, 12)
	if _, err := rand.Read(random); err != nil {
		return err
	}
	schema := "browser_" + hex.EncodeToString(random)
	if _, err := pool.Exec(setup, "CREATE SCHEMA "+pgx.Identifier{schema}.Sanitize()); err != nil {
		return errors.New("cannot create browser test schema")
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := pool.Exec(cleanup, "DROP SCHEMA "+pgx.Identifier{schema}.Sanitize()+" CASCADE"); err != nil {
			result = errors.Join(result, errors.New("cannot remove browser test schema "+schema))
		}
	}()
	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	isolated, err := database.Open(setup, parsed.String(), 10)
	if err != nil {
		return errors.New("invalid isolated test database configuration")
	}
	defer isolated.Close()
	if err := database.Migrate(setup, isolated); err != nil {
		return errors.New("browser fixture migration failed")
	}
	users, err := accounts.New(isolated)
	if err != nil {
		return err
	}
	const email = "browser@example.test"
	const password = "browser-testing-password"
	if _, err := users.Create(setup, email, password, true); err != nil {
		return errors.New("browser fixture account creation failed")
	}
	command := exec.CommandContext(ctx, "npm", "run", "test:browser")
	command.Env = append(os.Environ(), "DATABASE_URL="+parsed.String(), "BROWSER_TEST_EMAIL="+email, "BROWSER_TEST_PASSWORD="+password, "AI_ENABLED=false", "MAIL_ENABLED=false")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
