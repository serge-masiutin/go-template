package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serge-masiutin/go-template/internal/background"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/database"
	"github.com/serge-masiutin/go-template/internal/diagnostics"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(diagnostics.LibraryLogger(logger))
	if err := run(logger); err != nil {
		var invalid *config.Error
		if errors.As(err, &invalid) {
			logger.Error("invalid configuration", "reason", invalid.Error())
		} else {
			logger.Error("worker stopped", diagnostics.Attributes(err)...)
		}
		os.Exit(1)
	}
}
func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	signals, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Keep work alive during graceful drain; StopAndCancel enforces the shutdown deadline.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startup, startupCancel := context.WithTimeout(ctx, 10*time.Second)
	db, err := database.Open(startup, cfg.DatabaseURL, cfg.MaxConnections)
	startupCancel()
	if err != nil {
		return err
	}
	defer db.Close()
	worker, err := background.New(ctx, cfg, db, logger)
	if err != nil {
		return err
	}
	if err := worker.Client.Start(ctx); err != nil {
		shutdown, done := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer done()
		return errors.Join(err, worker.Stop(shutdown))
	}
	logger.Info("worker started", "mail_concurrency", cfg.WorkerConcurrency, "ai_concurrency", 1)
	<-signals.Done()
	shutdown, done := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer done()
	return worker.Stop(shutdown)
}
