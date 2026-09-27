package diagnostics

import (
	"context"
	"log/slog"
)

// LibraryLogger keeps service diagnostics but excludes arbitrary SQL/provider/job payloads.
// Application boundaries log the wrapped error classification separately.
func LibraryLogger(logger *slog.Logger) *slog.Logger {
	return slog.New(libraryHandler{base: logger.Handler()})
}

type libraryHandler struct {
	base  slog.Handler
	attrs []slog.Attr
}

func (h libraryHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}
func (h libraryHandler) Handle(ctx context.Context, record slog.Record) error {
	clean := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	appendSafe := func(attr slog.Attr) bool {
		switch attr.Key {
		case "service", "service_name", "job_id", "job_kind", "kind", "queue", "count", "num_jobs", "duration", "attempt":
			clean.AddAttrs(attr)
		}
		return true
	}
	for _, attr := range h.attrs {
		appendSafe(attr)
	}
	record.Attrs(appendSafe)
	return h.base.Handle(ctx, clean)
}
func (h libraryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.attrs = append(append([]slog.Attr(nil), h.attrs...), attrs...)
	return h
}
func (h libraryHandler) WithGroup(name string) slog.Handler {
	h.base = h.base.WithGroup(name)
	return h
}
