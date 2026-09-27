# Layer Violations

Cross-layer dependencies that break unidirectional data flow.

## Contents

- Ambient Identity in Models
- Request Objects in Services
- Notifications in Models
- Business Logic in Controllers

## Ambient Identity in Models

**Problem:** Models depend on presentation-layer context.

```go
var currentActor Actor
func (p *Post) Delete() { p.DeletedBy = currentActor.UserID } // BAD
```

**Issues:**
- Background execution has no request identity unless explicitly established
- Concurrent requests overwrite shared identity
- Hidden dependency makes testing harder
- Violates no-reverse-dependencies rule

**Fix:** Use explicit parameters.

```go
func (p *Post) DeleteBy(actorID int64) { p.DeletedBy = actorID }
```

## Request Objects in Services

**Problem:** Application layer depends on presentation layer.

```go
func (s *Registration) Register(r *http.Request) error { // BAD
    return s.create(r.Context(), r.FormValue("email"))
}
```

**Fix:** Extract value object in controller, pass to service.

```go
// The transport maps its DTO to an application-owned command.
command := accounts.RegistrationCommand{Email: input.Email, Name: input.Name}
user, err := registration.Register(r.Context(), command)
```

## Notifications in Models

**Problem:** Model triggers notifications, crossing into application layer.

```go
func (p *Post) Publish(ctx context.Context) error { // BAD
    p.Published = true
    return mail.SendPostPublished(ctx, p.ID)
}
```

**Issues:**
- Domain layer depends on application layer (reverse dependency)
- Model has side effects beyond state management
- Harder to test model in isolation
- Notification may fire unexpectedly from different call sites

**Fix:** Trace the call chain, move notification to existing orchestrator.

```text
Move notification orchestration to the existing Publish operation. It calls the
pure transition, persists the post and records delivery intent in one transaction.
The worker performs external delivery after commit with a stable event ID.
```

**Resolution process:**
1. Find the caller (handler, operation, worker)
2. If orchestrator exists → move notification there
3. If no orchestrator → suggest an operation if there is a multi-step use case; keep trivial delegation local

## Business Logic in Controllers

**Problem:** Presentation layer doing domain work.

```go
// BAD: the HTTP layer owns pricing rules.
if order.Total > 10000 { order.Discount = order.Total / 10 }
```

**Fix:** Move domain logic to model, orchestration to service if needed.

```go
// GOOD: the entity owns its pricing invariant; the operation owns persistence.
if err := order.ApplyDiscount(policy); err != nil { return err }
return store.Save(ctx, order)
```
