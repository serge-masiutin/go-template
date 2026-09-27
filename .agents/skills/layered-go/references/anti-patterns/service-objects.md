# Service Object Anti-Patterns

Common mistakes when introducing service objects.

## Contents

- Anemic Models
- Bag of Random Objects
- Premature Abstraction

## Anemic Models

**Problem:** All logic moved to services, entities become passive data containers.

```text
BAD: OrderValidator, OrderCalculator and OrderStateUpdater each manipulate fields
of the same passive Order struct. Its invariants can be bypassed by any caller.
```

**Fix:** Keep domain logic in models. Services orchestrate, models know their business rules.

```go
func (o *Order) Cancel() error {
    if o.State != Pending { return ErrCannotCancel }
    o.State = Cancelled
    return nil
}
```

## Bag of Random Objects

**Problem:** No conventions, each service is unique.

```text
BAD: one operation returns bool, another panics for a normal validation failure,
a third mutates an error field and a fourth returns map[string]any. Callers must
learn each convention independently.
```

**Fix:** Establish conventions.

```text
Use meaningful typed parameters and (value, error) results. Wrap causes with %w,
classify expected failures with errors.Is/errors.As, and establish transaction
ownership explicitly. Do not require universal Execute/Call interfaces.
```

## Premature Abstraction

**Problem:** Creating abstractions before patterns emerge.

```text
BAD: introducing BaseService, reflection registration, generic repositories and
a result-monad package before any repeated contract exists.
```

**Fix:** Wait for patterns to emerge. Start simple.

```text
Start with one explicit operation. When multiple operations share stable behavior,
extract a narrow function or composed collaborator and verify the reduced coupling.
```
