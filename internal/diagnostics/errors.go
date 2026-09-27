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
	var visit func(error)
	visit = func(cause error) {
		if cause == nil {
			return
		}
		chain = append(chain, fmt.Sprintf("%T", cause))
		if joined, ok := cause.(interface{ Unwrap() []error }); ok {
			for _, nested := range joined.Unwrap() {
				visit(nested)
			}
		} else {
			visit(errors.Unwrap(cause))
		}
	}
	visit(err)
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
