# Extraction Signals

How to identify code that should be extracted to a different layer or abstraction.

## Contents

- Callback Scoring System
- God Object Identification
- Composition Health Check
- Service Object Signals
- Handler Fat Signals
- Quick Reference
- Static Analysis Tools

## Callback Scoring System

Rate existing hooks (do not add hooks to implement these scores) to identify extraction candidates:

The canonical scale is [Callbacks → Scoring](../topics/callbacks.md#callback-scoring-system):
5 = transformation; 4 = internal consistency; 3 = timestamp; 2 = background trigger;
1 = external operation. Preserve local invariants, make all calls explicit, and move
business delivery intent into the owning operation. This is an audit scale, not a
recommendation to install hooks.

### Transformer Callbacks (Keep)

Compute or default required attribute values:

```go
func (p Post) WordCount() int { return len(strings.Fields(p.Content)) }
// Compute locally and explicitly. Do not invent persistence hooks in Go.
```

### Normalizer Callbacks (Keep)

Normalize user input explicitly at the boundary:

```go
func ParseTitle(raw string) (Title, error) {
    title := strings.TrimSpace(raw)
    if title == "" { return "", ErrEmptyTitle }
    return Title(title), nil
}
```

### Utility Callbacks (Keep)

Framework-level utilities like counter caches:

```sql
-- Maintain a required counter in the same transaction as its source row.
UPDATE posts SET comment_count=comment_count+1 WHERE id=$1;
```

### Operation Callbacks (Extract)

Signs of misplacement:
- Caller-dependent conditions or skip flags
- Collaboration with non-model objects (mailers, API clients)
- Remote peer communication

```go
// BAD: Save also sends mail and analytics without telling its caller.
func (s *Store) Save(ctx context.Context, user User) error {
    if err := s.insert(ctx, user); err != nil { return err }
    return s.mailer.Welcome(ctx, user)
}
```

**Extraction options:**
1. Make the use-case call explicit at the handler
2. Move to service object
3. Use event-driven approach

```go
// GOOD: one operation persists user + outbox intent atomically.
user, err := registration.Register(ctx, input)
if err != nil { return fmt.Errorf("register: %w", err) }
// A worker claims committed events; retry and delivery idempotency are explicit.
```

## God Object Identification

### Churn × Complexity Metric

**Churn** = how often a file changes (indicates ongoing modifications)
**Complexity** = measured complexity using the project Go analyzer

Files high in both are prime refactoring candidates.

```sh
git log --format=oneline -- internal/accounts/accounts.go
# Compare churn with actual responsibilities and measured complexity.
# Use the project's configured Go analyzer if one exists.
```

### Automated Tool

Use existing package analysis and inspect high-churn files:

```sh
go list ./...
go vet ./...
```

### Common God Object Names

Watch for these accumulating responsibilities:
- `User` / `Account`
- `Order` / `Transaction`
- `Project` / `Workspace`
- `Post` / `Article`

### Decomposition Strategies

1. **Extract composed behaviors** for shared behaviors
2. **Extract delegate objects** for complex operations
3. **Extract value objects** for groups of related attributes
4. **Create new models** for distinct concepts

```go
type User struct { ID int64; Email string }
type Profile struct { UserID int64; DisplayName string }
type NotificationPreferences struct { EmailEnabled bool }
// Account, billing, and notification operations own separate dependencies.
```

## Composition Health Check

### Good Compositions (Behavioral)

Can be tested in isolation, shared across models:

```go
type Publication struct { PublishedAt time.Time }
func (p *Publication) Publish(now time.Time) error {
    if !p.PublishedAt.IsZero() { return ErrAlreadyPublished }
    p.PublishedAt = now
    return nil
}
```

**Test:** Can you write specs for this composed behavior without instantiating the host model?

### Bad Compositions (Code-Slicing)

Groups code by file or technical artifact type, not behavior:

```go
// BAD: embedding an entire service exposes unrelated methods and state.
type Contact struct { *ApplicationServices }
// Prefer a focused value such as ContactDetails with explicit behavior.
```

**Test:** If removing this composed behavior breaks unrelated tests, it's code-slicing.

### Overgrown Compositions

Signs a composed behavior should be extracted:
- Unrelated reasons to change (size alone is not a verdict)
- Multiple responsibilities
- Complex internal state

Extract to:
- **Delegate object** for cohesive pure behavior; application operation for I/O orchestration
- **Value object** for attribute groups
- **Separate model** for distinct entity

## Service Object Signals

### When to Extract to Service

Extract from handler when you see:
- Multiple model operations in sequence
- Transaction spanning multiple models
- Complex error handling
- Reusable business operation

### When NOT to Use Services

Don't extract:
- Single model operations (keep in model)
- Simple CRUD (handler parses/delegates; store owns SQL)
- Domain logic (belongs in model, not service)

### Anemic Model Warning Signs

Your models might be anemic if:
- Services contain calculations that use only model data
- Models are pure data containers (associations + validations only)
- You have `CalculateXService` for model attributes
- Domain rules live in services, not models

```go
type Item struct { UnitPriceCents int64; Quantity int64 }
func (i Item) SubtotalCents() int64 { return i.UnitPriceCents * i.Quantity }
// Validate monetary ranges before arithmetic; production money types need
// currency, rounding and overflow contracts. Do not use floating-point money.
```

## Handler Fat Signals

### Extract When You See

- Business calculations (pricing, discounts)
- Multiple model updates
- Complex conditionals based on business rules
- External API calls
- Mixed transport and business responsibilities

### Keep in Handler

- Parameter parsing
- Authentication/authorization
- Response formatting
- Simple model operations

## Quick Reference

| Signal | Threshold | Action |
|--------|-----------|--------|
| Hook responsibility | Background/external effect | Make intent explicit in the operation |
| Model complexity | Mixed responsibilities plus measured complexity | Decompose |
| Model churn | High relative churn | Review for extraction |
| Composition size | Multiple independent responsibilities | Extract to delegate |
| Handler action | HTTP mixed with domain orchestration | Extract to service |
| Service with domain logic | Any calculations | Move to model |

## Static Analysis Tools

| Tool | Purpose | Command |
| --- | --- | --- |
| Go compiler | Import cycles, types, package boundaries | `go test ./...` |
| go vet | Suspicious language/API usage | `go vet ./...` |
| Race detector | Executed concurrent behavior | `go test -race ./...` |
| Git | Churn evidence | `git log -- path/to/file.go` |

Line counts are investigation signals, not extraction mandates. A thin HTTP wrapper
can be longer in Go because it handles errors explicitly. Do not add abstraction to
meet a line-count target or import Ruby analysis tools into a Go project.
