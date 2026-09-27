package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/database"
	"github.com/serge-masiutin/go-template/internal/diagnostics"
	"github.com/serge-masiutin/go-template/internal/httpapp"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		var invalid *config.Error
		if errors.As(err, &invalid) {
			logger.Error("invalid configuration", "reason", invalid.Error())
		} else {
			logger.Error("server stopped", diagnostics.Attributes(err)...)
		}
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	startup, cancel := context.WithTimeout(ctx, 10*time.Second)
	pool, err := database.Open(startup, cfg.DatabaseURL, cfg.MaxConnections)
	cancel()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer pool.Close()
	sessionStore := pgxstore.New(pool)
	defer sessionStore.StopCleanup()
	sessions := scs.New()
	sessions.Store = sessionStore
	sessions.Lifetime = 24 * time.Hour
	sessions.IdleTimeout = 2 * time.Hour
	sessions.Cookie.Name = "go_template_session"
	sessions.Cookie.Secure = cfg.Production
	sessions.Cookie.SameSite = http.SameSiteLaxMode
	handler, err := httpapp.New(cfg, pool, sessions, logger)
	if err != nil {
		return err
	}
	server := &http.Server{Addr: cfg.Address, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 32 << 10}
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	stopped := make(chan error, 1)
	go func() { stopped <- server.Serve(listener) }()
	logger.Info("server started", "address", listener.Addr().String())
	select {
	case err := <-stopped:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		err := server.Shutdown(shutdown)
		if err != nil {
			// Cancel remaining request contexts before closing their database pool.
			server.Close()
		}
		return err
	}
}
