package diagnostics

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

func TestLibraryLoggerKeepsContextWithoutPayloads(t *testing.T) {
	var output bytes.Buffer
	logger := LibraryLogger(slog.New(slog.NewJSONHandler(&output, nil))).With("queue", "ai", "args", "private note")
	logger.Error("job failed", "job_id", 42, "error", errors.New("private provider response"), "sql", "secret query")
	text := output.String()
	if strings.Contains(text, "private") || strings.Contains(text, "secret") || !strings.Contains(text, `"job_id":42`) || !strings.Contains(text, `"queue":"ai"`) {
		t.Fatal(text)
	}
}
