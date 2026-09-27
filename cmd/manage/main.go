package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/serge-masiutin/go-template/internal/accounts"
	"github.com/serge-masiutin/go-template/internal/config"
	"github.com/serge-masiutin/go-template/internal/database"
	"github.com/serge-masiutin/go-template/internal/diagnostics"
)

func main() {
	if err := run(); err != nil {
		var invalid *config.Error
		if errors.As(err, &invalid) {
			slog.Error("invalid configuration", "reason", invalid.Error())
		} else {
			slog.Error("management command failed; verify arguments, configuration and database constraints", diagnostics.Attributes(err)...)
		}
		os.Exit(1)
	}
}
func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: manage migrate | create-user --email ADDRESS [--admin]; password from stdin")
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, cfg.DatabaseURL, cfg.MaxConnections)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer pool.Close()
	switch os.Args[1] {
	case "migrate":
		if err := database.Migrate(ctx, pool); err != nil {
			return err
		}
		fmt.Println("Migrations applied.")
		return nil
	case "create-user":
		flags := flag.NewFlagSet("create-user", flag.ContinueOnError)
		email := flags.String("email", "", "Account email")
		admin := flags.Bool("admin", false, "Grant administrator access")
		if err := flags.Parse(os.Args[2:]); err != nil {
			return err
		}
		password, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return fmt.Errorf("read newline-terminated password from stdin: %w", err)
		}
		store, err := accounts.New(pool)
		if err != nil {
			return err
		}
		if _, err := store.Create(ctx, *email, strings.TrimSuffix(strings.TrimSuffix(password, "\n"), "\r"), *admin); err != nil {
			return fmt.Errorf("create user failed; check address, password length and uniqueness: %w", err)
		}
		fmt.Println("Account created.")
		return nil
	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}
}
