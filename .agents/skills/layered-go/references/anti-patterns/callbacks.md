# Callback Anti-Patterns

Hidden persistence hooks that signal a layering problem.

## Contents

- Operation Callbacks
- Skip Callback Anti-Pattern
- Callback Control Flags

## Operation Callbacks

**Problem:** Business process steps disguised as callbacks.

```go
// BAD: saving a user also calls a remote provider.
func (s *UserStore) Save(ctx context.Context, user User) error {
    if err := s.insert(ctx, user); err != nil { return err }
    return s.mail.SendWelcome(ctx, user.ID)
}
```

**Fix:** Extract to the owning application operation and explicit delivery boundary.

```go
// GOOD: the operation commits state and durable notification intent atomically.
func (s *Registration) Register(ctx context.Context, input RegistrationCommand) error {
    return s.withTx(ctx, func(tx RegistrationTx) error {
        user, err := tx.CreateUser(ctx, input)
        if err != nil { return err }
        return tx.AddWelcomeEvent(ctx, user.ID)
    })
}
```

## Skip Callback Anti-Pattern

**Problem:** A global skip flag creates hidden dependencies.

```go
// BAD: callers toggle shared authentication behavior.
handlers.SkipAuthentication = true
handlers.ServeHTTP(w, r)
```

**Fix:** Use explicit middleware composition.

```go
// GOOD: route middleware composition declares which endpoints require identity.
mux.Handle("GET /login", loginPage)
mux.Handle("GET /admin", authenticate(requireAdmin(adminPage)))
```

## Callback Control Flags

**Problem:** Virtual attributes to skip callbacks.

```go
type User struct {
    Email string
    SkipWelcome bool // BAD: transport/workflow choice hidden in entity state.
}
```

**Fix:** Extract callbacks, call explicitly when needed.

```text
Use an explicit Register operation that records welcome intent, and a separately
named Import operation with its documented notification policy. Both enforce the
same account invariants; neither disables validation or mutates global hooks.
```
