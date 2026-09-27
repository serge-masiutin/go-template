# Extract Callbacks to Service

Move multi-step user-creation side effects out of persistence hooks into an explicit service.

## Before

```go
// BAD: callers cannot persist a user without triggering network effects.
func (s *Store) Save(ctx context.Context, user User) error {
    if err := s.insert(ctx, user); err != nil { return err }
    if err := s.workspaces.CreateDefault(ctx, user.ID); err != nil { return err }
    if err := s.mailer.Welcome(ctx, user); err != nil { return err }
    if err := s.mailer.NotifyAdmins(ctx, user.ID); err != nil { return err }
    return s.analytics.Signup(ctx, user.ID)
}
```

## After

```go
// A use case owns the transaction. The outbox and user commit together.
func (s *Registration) Register(ctx context.Context, command RegistrationCommand) (User, error) {
    user, err := NewUser(command.Email)
    if err != nil { return User{}, err }
    tx, err := s.pool.Begin(ctx)
    if err != nil { return User{}, err }
    defer tx.Rollback(ctx)
    if err := s.users.Insert(ctx, tx, &user); err != nil { return User{}, err }
    if err := s.workspaces.CreateDefault(ctx, tx, user.ID); err != nil { return User{}, err }
    if err := s.outbox.UserRegistered(ctx, tx, user.ID); err != nil { return User{}, err }
    if err := tx.Commit(ctx); err != nil { return User{}, err }
    return user, nil
}
// UserRegistered has three durable consumers: welcome email, admin notice,
// and signup analytics. Each consumer records its own acknowledgement/retry state.
// All three effects from Before remain; the default workspace is still created.
// The outbox API illustrates a feature to implement; it is not in this starter.
```

The example has two stages: first extract the same four responsibilities without changing
behavior; then deliberately strengthen reliability with atomic user/workspace/outbox writes
and asynchronous delivery. The latter changes timing and failure semantics and must be
reviewed/tested as a behavior change, not disguised as a pure refactor. Welcome email,
admin notification and signup analytics each consume the same stable event identity;
one successful consumer must not mark another failed consumer delivered.

Verify user/workspace/outbox rollback, all three delivery intents, independent retry state,
no delivery before commit, and duplicate handling after provider success/local-ack failure.
`RegistrationCommand` belongs to the application package; HTTP maps its request DTO to it.
