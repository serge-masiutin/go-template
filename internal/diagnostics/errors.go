package diagnostics

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

// Attributes retains error classification without exposing SQL detail or URLs.
// Errors remain wrapped internally; only their safe diagnostic view is logged.
func Attributes(err error) []any {
	chain := make([]string, 0)
	for cause := err; cause != nil; cause = errors.Unwrap(cause) {
		chain = append(chain, fmt.Sprintf("%T", cause))
	}
	attributes := []any{"error_chain", chain}
	var postgres *pgconn.PgError
	if errors.As(err, &postgres) {
		attributes = append(attributes, "sqlstate", postgres.Code)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		attributes = append(attributes, "cause", "deadline_exceeded")
	}
	if errors.Is(err, context.Canceled) {
		attributes = append(attributes, "cause", "canceled")
	}
	return attributes
}
