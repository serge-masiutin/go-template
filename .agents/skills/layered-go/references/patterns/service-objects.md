# Service Objects

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns
- From Services to Abstractions
- Go tools

## Summary

Service objects represent single business operations between handlers and models. They fill the gap in an HTTP application where neither handlers (inbound layer) nor models (domain layer) are appropriate homes for business logic orchestration.

## When to Use

- Orchestrating multiple domain objects for a use case
- Operations that span multiple models with transaction boundaries
- Complex business logic that doesn't belong in any single model
- Reusable operations called from multiple handlers

## When NOT to Use

- Simple CRUD (handler parses/delegates; store owns SQL)
- Domain logic (belongs in models)
- Single-entity behavior (keep on its owning type)

## Key Principles

- **Waiting room metaphor** — services are temporary homes until proper abstractions emerge
- **Establish conventions early** — naming, interface, dependency ownership, return values
- **Services orchestrate, models know rules** — don't strip domain logic into services
- **Avoid bag of random objects** — decompose into specialized abstractions as patterns emerge

## Implementation

### Explicit operation and constructor

Go does not need `ApplicationService`, inheritance, `.call`, or a result monad.
Use a function for stateless work; use a struct when the operation owns dependencies.
Keep interfaces small and define them at the consumer only when substitution is useful.

```go
type Registration struct { users UserStore }
func NewRegistration(users UserStore) *Registration { return &Registration{users: users} }
func (s *Registration) Register(ctx context.Context, command RegistrationCommand) (User, error) {
    user, err := NewUser(command.Email)
    if err != nil { return User{}, err }
    if err := s.users.Insert(ctx, &user); err != nil { return User{}, fmt.Errorf("insert user: %w", err) }
    return user, nil
}
```

### Conventions to Establish

| Decision | Go convention |
| --- | --- |
| Naming | Domain package and meaningful operation, such as `accounts.Register` |
| Interface | Explicit parameters, `context.Context` first for cancellable I/O |
| Shared machinery | Composition only after repetition proves a need |
| Return values | Typed value plus `error` |
| Error handling | Sentinel or typed errors; `%w`, `errors.Is` / `errors.As` |

Transaction orchestration belongs here when a use case spans records. A one-line wrapper
around an existing store method is not a reason to create a service. See the
[callback extraction](../../examples/callbacks-to-service.md) for atomic notification intent.

## Anti-Patterns

### Anemic Models

If `OrderCalculator`, `OrderValidator`, and `OrderStateChanger` only manipulate one order, move stable rules onto `Order` or its values. Keep cross-record transactions and external effects in the operation.

### Bag of Random Objects

Use consistent typed inputs, context for cancellable I/O, `(value, error)` results and `%w` wrapping. Do not force every operation to implement a universal `Execute(any) any` interface.

### Premature Abstraction

A base service, reflection registry or generic transaction framework needs evidence of repeated stable semantics. Start with explicit functions/structs; extract shared machinery only after its contract becomes clear.

## From Services to Abstractions

As services accumulate, decompose into specialized patterns:

- **Form objects** — user input handling
- **Query objects** — complex queries
- **Policy objects** — authorization
- **Presenter objects** — view logic

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
