# Instrumentation

## Summary

Instrumentation observes behavior without changing domain decisions. Preserve the source's separation of events, structured logs, metrics and tracing; use explicit Go boundaries instead of an implicit Rails event bus.

## Layer Placement

Request/job boundaries measure outcomes. Application operations expose meaningful outcomes and wrap errors. Infrastructure emits transport diagnostics. Pure domain types do not depend on a logging or metrics backend.

## Key Principles

Use bounded labels, explicit time units, causal errors and stable operation names. Never log passwords, tokens, full request bodies, raw prompts, provider responses or database URLs. Observability must not become a new required dependency for pure domain behavior.

## Event Instrumentation

### Subscribe to Events

Do not introduce a global event bus merely to count an operation. Call a narrow observer or use middleware when the same boundary is already present. If events serve business delivery, use a durable mechanism separate from telemetry.

### Custom Events

Name an event for a stable operation/outcome, not a user-supplied string. Include correlation and safe identifiers only under the project's data policy. Record state transition and external delivery separately.

### Observer Pattern

An observer must have explicit failure semantics. Metrics failure should not silently alter a successful business transaction; log/measure collector failures at the observability boundary. Required audit records are business data and may have stricter atomicity requirements.

## Structured Logging

### Contextual Logging

Use `log/slog`, a stable operation name, duration, result and request correlation. Context cancellation is expected at request boundaries and should be distinguished from an infrastructure outage.

### JSON Logging

The starter emits JSON logs. Log causes safely: PostgreSQL details can include private values, so emit SQLSTATE and operation context instead of raw query parameters or complete provider errors.

### Logger Choice

Use the installed logger first. A second logging library needs a concrete missing capability and a migration plan, not stylistic preference.

## Metrics

### Exporter Integration

Choose the deployed collector before adding a library. Expose metrics on a deliberately controlled endpoint. Never use email, resource ID, raw URL, error text or request ID as a metric label.

### Custom Metrics

Measure request/operation count, duration, failure class, pool saturation and queue age when relevant. Use route templates and bounded result enums. A histogram's buckets and unit must match the operating question.

## Service Instrumentation

### Shared Instrumentation

Composition or a narrow function wrapper can measure repeated operation boundaries. Do not add a base service type that every use case must inherit or a wrapper that hides retries and transactions.

### Manual Instrumentation

```go
started := time.Now()
result, err := operation.Run(ctx, command)
logger.InfoContext(ctx, "operation completed",
    "operation", "report.generate", "duration_ms", time.Since(started).Milliseconds(),
    "failed", err != nil)
return result, err
```

The example reports outcome, not private command content. The error's cause remains available to the owning boundary.

## Anti-Patterns

### Logging in Domain Models

Domain methods return values/errors. Callers decide how to report them once; avoid duplicated stackless logs at every layer.

### Metrics Scattered in Code

Keep metric names, units and label vocabulary consistent. Do not increment "success" before commit or provider acknowledgement.

### Verbose Debug Logging

Debug logging is not permission to expose payloads. Expensive serialization for a disabled log level also needs measurement.

## Testing Instrumentation

Capture logs using a test handler; assert event name, outcome and absence of secrets. Test unexpected-error and cancellation paths. Do not assert exact duration or global event ordering across goroutines.

## Performance Considerations

Profile serialization, contention and exporter overhead before sampling or buffering. If buffering is used, define bounds, drop behavior, flush on shutdown and monitoring for lost records.
