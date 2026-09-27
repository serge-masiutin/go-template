package diagnostics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestJoinedErrorsKeepClassificationWithoutPrivateContent(t *testing.T) {
	err := errors.Join(fmt.Errorf("private-provider-response: %w", context.DeadlineExceeded), &pgconn.PgError{Code: "23505", Detail: "private-account-email"})
	encoded, marshalErr := json.Marshal(Attributes(err))
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	for _, want := range []string{"23505", "deadline_exceeded", "*pgconn.PgError", "*fmt.wrapError"} {
		if !strings.Contains(string(encoded), want) {
			t.Fatalf("missing classification %s: %s", want, encoded)
		}
	}
	if strings.Contains(string(encoded), "private-") {
		t.Fatal("private error text reached diagnostics")
	}
}
