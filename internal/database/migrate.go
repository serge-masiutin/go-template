package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"errors"
	"io/fs"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
)

//go:embed migrations/*.sql bootstrap/*.sql
var migrations embed.FS

// Migrate uses Goose's deploy lock and atomic per-migration transactions, then River's versioned schema.
func Migrate(ctx context.Context, db *DB) (result error) {
	locker, err := lock.NewPostgresSessionLocker(lock.WithLockID(724319))
	if err != nil {
		return err
	}
	conn, err := db.SQL.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := locker.SessionLock(ctx, conn); err != nil {
		return err
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if err := locker.SessionUnlock(cleanup, conn); err != nil {
			// Never return a possibly locked PostgreSQL session to the application pool.
			conn.Raw(func(any) error { return driver.ErrBadConn })
			result = errors.Join(result, err)
		}
	}()
	files, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return err
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, db.SQL, files,
		goose.WithDisableGlobalRegistry(true),
		goose.WithGoMigrations(goose.NewGoMigration(1, &goose.GoFunc{RunTx: bootstrap}, nil)))
	if err != nil {
		return err
	}
	if _, err := provider.Up(ctx); err != nil {
		return err
	}
	river, err := rivermigrate.New(riverdatabasesql.New(db.SQL), nil)
	if err != nil {
		return err
	}
	_, err = river.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err
}

// Adopt the only pre-Goose migration without recreating tables or losing existing accounts.
func bootstrap(ctx context.Context, tx *sql.Tx) error {
	var legacy bool
	if err := tx.QueryRowContext(ctx, "SELECT to_regclass('schema_migrations') IS NOT NULL").Scan(&legacy); err != nil {
		return err
	}
	if legacy {
		var total, known int
		if err := tx.QueryRowContext(ctx, "SELECT count(*), count(*) FILTER (WHERE name = 'migrations/001_accounts.sql') FROM schema_migrations").Scan(&total, &known); err != nil {
			return err
		}
		if total != 1 || known != 1 {
			return errors.New("unrecognized legacy migration history; migrate explicitly before adopting Goose")
		}
		return nil
	}
	script, err := migrations.ReadFile("bootstrap/001_accounts.sql")
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, string(script))
	return err
}
