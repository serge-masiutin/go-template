# State Machines

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Identifying Implicit State Machines
- Implementation
- Anti-Patterns
- Go tools

## Summary

State machines describe possible states, transitions between them, and triggering events. They make implicit state logic explicit and centralized, preventing scattered conditional logic across the codebase.

## When to Use

- Multiple related boolean/timestamp attributes tracking state
- Complex conditional logic based on object state
- State-dependent behavior
- Audit trail requirements

## When NOT to Use

- Linear progressions (just use enum)
- Simple two-state flags
- Before patterns emerge (premature)

## Key Principles

- **Identify implicit state machines** — scattered booleans/timestamps often hide state
- **Prefer events over direct transitions** — `order.Submit()` rather than unrestricted `order.State = Submitted`
- **Extract when not central** — standalone workflow for secondary state
- **Guards and callbacks sparingly** — don't recreate callback problems

## Identifying Implicit State Machines

Signs of hidden state: overlapping status booleans, timestamps used as competing sources of truth, impossible combinations, repeated transition conditions, and state changes without a named event.

## Implementation

### Typed state and guarded transition

An enum-like defined string type plus explicit methods is enough for a small machine.
Reject unknown states at storage/input boundaries. Do not install a workflow framework
solely to replace a Rails DSL.

```go
type State string
const (Pending State = "pending"; Paid State = "paid"; Shipped State = "shipped")
func (o *Order) Ship(now time.Time) error {
    if o.State != Paid { return ErrInvalidTransition }
    o.State, o.ShippedAt = Shipped, now
    return nil
}
```

### Persistence and concurrency

Protect the transition with a locked transaction or conditional update. Check affected
rows and distinguish an invalid transition from a database failure. State timestamps
are evidence of transitions; contradictory timestamp combinations are not the state model.

### Side effects and complex flows

Persist delivery intent atomically with the transition if loss is unacceptable. A worker
handles delivery with an explicit retry/idempotency contract. For a wizard, define
allowed back/forward edges and final validation. See
[state-machine extraction](../../examples/implicit-to-explicit-state-machine.md).

## Anti-Patterns

### Implicit State Machines

Replace competing booleans with a single validated state and named events when combinations become hard to reason about. Preserve existing transition behavior before adding new restrictions.

### Phantom Transitions

Every reported success must correspond to a committed allowed transition. Check affected rows for conditional updates; a zero-row update is a conflict or missing authorized record, not success.

### Excessive Guards

Guards should express domain prerequisites, not perform network calls or unrelated workflow steps. Precompute external facts in an operation, then apply the transition atomically.

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
