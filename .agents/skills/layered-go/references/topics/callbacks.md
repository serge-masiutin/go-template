# Callbacks and Explicit Operations

## Summary

Callbacks are hooks into object lifecycle events. The original scoring framework diagnoses hidden behavior; it is not a reason to introduce a callback system in Go. Preserve the distinction between maintaining an entity and orchestrating a business process.

## Callback Scoring System

| Score | Responsibility | Go placement |
| --- | --- | --- |
| 5/5 | Pure transformation | Explicit parser or constructor |
| 4/5 | Internal consistency | Entity method or atomic storage operation |
| 3/5 | Timestamp | Explicit event time argument or database default |
| 2/5 | Background trigger | Application operation plus durable delivery intent |
| 1/5 | Synchronous external operation | Explicit application operation and infrastructure adapter |

The score measures extraction urgency for existing hooks. A 5/5 hook can still hide execution order; in new Go code call the transformation directly.

## Good Callbacks (4–5/5): Preserve the Responsibility

### Transformers

Normalize at the input boundary, then validate the normalized value. Do not silently turn malformed data into an empty value. Keep normalization idempotent and domain-specific: trimming an email is different from changing a password.

```go
func ParseEmail(raw string) (Email, error) {
    normalized := strings.ToLower(strings.TrimSpace(raw))
    address, err := mail.ParseAddress(normalized)
    if err != nil || address.Address != normalized { return "", ErrInvalidEmail }
    return Email(normalized), nil
}
```

### Maintainers

A stored counter must change in the same transaction as the rows it summarizes. A read-count-then-write sequence without locking loses concurrent changes. Prefer a constraint or atomic SQL update when it expresses the invariant directly; do not recalculate unrelated associations on every save.

### Safe Timestamps

Supply `now` to a domain transition so tests control time. A database default is appropriate for creation time. Do not infer a complex workflow from a collection of nullable timestamps.

## Problematic Callbacks (1–2/5)

### Background Triggers (2/5)

An enqueue after commit avoids sending uncommitted data but leaves a crash window between commit and enqueue. If delivery must survive process failure, persist an outbox record in the same transaction, then have a worker deliver it. A goroutine launched by a request is neither durable nor an after-commit guarantee.

### Operations (1/5): Extract These

Welcome messages, workspace provisioning and external billing are use-case steps. Move them to a named operation with an explicit transaction and external-effect contract. Do not send mail inside the database transaction or rely on rolling back an HTTP request to undo an external charge.

## Callback Anti-Patterns

### Conditional Complexity

A hook that branches on caller, request context, `skipEmail`, or an environment variable hides distinct commands. Name the supported operations and pass validated inputs.

### Callback Chains

Draw the actual execution graph before extraction: user save → workspace save → owner update → email. Preserve write ordering and constraints, remove cycles, and put required writes in one transaction.

### Testing Difficulties

If a domain test needs SMTP, queue and HTTP mocks merely to construct an entity, its lifecycle crosses layers. Test domain behavior directly and use database integration tests for transaction atomicity.

### Skip Callback Temptation

Do not solve imports or fixtures with global hook toggles. An import is an explicit operation with a documented notification policy. It must still enforce durable invariants.

## When Callbacks Are Appropriate

Library lifecycle hooks may be necessary for protocol/resource integration. Keep them local to that adapter, document invocation order and error propagation, and test failure paths. Do not hide business workflows in an ORM hook.

## Callback to Service Extraction

Use [the original registration scenario adapted to Go](../../examples/callbacks-to-service.md). First characterize existing behavior, then move orchestration, then remove the hook. Do not add new permissions or notifications during the structural change.

## Testing Callback-Light Models

Test pure normalization, invalid values, explicit transitions, rollback at each write, concurrent updates, and delivery after commit. Verify that constructing/loading an entity does not enqueue work. Test the crash window and duplicate processing when delivery is required.
