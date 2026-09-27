package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB owns one pgx connection pool. SQL and ORM share its connection budget.
type DB struct {
	*pgxpool.Pool
	SQL *sql.DB
	ORM *gorm.DB
}

func Open(ctx context.Context, connectionString string, maxConnections int32) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}
	cfg.MaxConns = maxConnections
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect database: %w", err)
	}
	// OpenDBFromPool disables idle SQL connections so SCS and pgx can share the pool.
	sqlDB := stdlib.OpenDBFromPool(pool)
	orm, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,           // Each single-statement write is atomic; operations open explicit transactions.
		Logger:                 logger.Discard, // SQL and driver errors can contain passwords, addresses and note bodies.
	})
	if err != nil {
		sqlDB.Close()
		pool.Close()
		return nil, err
	}
	return &DB{Pool: pool, SQL: sqlDB, ORM: orm}, nil
}

func (db *DB) Close() {
	db.SQL.Close()
	db.Pool.Close()
}
