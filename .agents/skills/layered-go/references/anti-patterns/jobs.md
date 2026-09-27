# Job Anti-Patterns

Background entry points and wrappers with unclear responsibility.

## Anemic Jobs

**Problem:** Wrappers that add no contract beyond a second layer of delegation.

```text
A typed job handler that validates a payload, establishes authority and calls an
operation is not anemic: it owns the queue protocol. A second wrapper that adds
none of these responsibilities is redundant.
```

**Issues:**
- Redundant wrappers hiding meaningful worker code
- Two places to maintain (job + model method)
- Wrapper owns neither payload validation nor queue semantics
- Makes the jobs folder noisy, hiding complex jobs that need attention

**Signal:** A forwarding wrapper adds no validation, authority, error mapping or queue integration.

**Fix:** Keep a typed queue entry point when it owns the queue boundary; remove only genuinely redundant wrappers. Do not port generated job methods into Go.

```go
type RebuildIndexJob struct { TenantID int64; DocumentID int64 }
func (w *Worker) Rebuild(ctx context.Context, job RebuildIndexJob) error {
    if job.TenantID <= 0 || job.DocumentID <= 0 { return ErrInvalidJob }
    return w.index.Rebuild(ctx, job.TenantID, job.DocumentID)
}
```

**Worker configuration:**
```text
Configure queue name, attempt limit, retryable error classes, backoff, timeout,
lease expiry and dead-letter behavior at the worker boundary. Persist stable IDs;
reload authorized state at execution. No queue is installed by this starter.
```

**When to keep separate job class:**
- Job has complex logic beyond single method call
- Job processes multiple records with custom batching
- Job needs extensive retry/error handling configuration
- Job is triggered from multiple unrelated models
