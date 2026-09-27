# Architecture Layers

## Contents

- The Four Layers (Presentation, Application, Domain, Infrastructure)
- The Four Layering Rules
- Mapping Go Components to Layers
- Applying Layered Architecture in Practice
- Common Layering Mistakes

## The Four Layers (Presentation, Application, Domain, Infrastructure)

Go applications are organized into four architecture layers with unidirectional data flow:

| Layer | Responsibility | Go Examples |
|-------|----------------|----------------|
| **Presentation** | Handle user interactions, present information; all inbound entry points | HTTP handlers, React views, stream handlers, workers (internal inbound), Form objects, Presenters |
| **Application** | Organize domain objects for use cases | Operations, policies, notification orchestration |
| **Domain** | Entities, rules, invariants, application state | Entities, Value objects, Domain events |
| **Infrastructure** | Supporting technologies | pgx/database/sql, API clients, File storage |

```
Runtime calls: HTTP → application operation → domain behavior / persistence ports
Imports: adapters depend on consumer-owned contracts; domain must not import HTTP or pgx.
```

## The Four Layering Rules

### Rule 1: Unidirectional Data Flow

Runtime work moves from entry points to operations and domain behavior. Go import direction is a separate concern: infrastructure adapters may import domain values or implement consumer-owned interfaces; domain packages do not import concrete adapters. Do not force the runtime arrow into an import chain.

### Rule 2: No Reverse Dependencies

Lower layers must not depend on higher layers. A domain object should never depend on a handler or request object.

**Violations:**

```go
// BAD: the application consumes transport and ambient identity.
func HandleEvent(r *http.Request) error { return process(r.FormValue("event")) }
var currentUser User
func (p *Post) Delete() { p.DeletedBy = currentUser.ID }
```

**Correct:**

```go
// GOOD: explicit domain input; cancellation is infrastructure context only.
func HandleEvent(ctx context.Context, event GitHubEvent) error {
    return events.Record(ctx, event)
}
func (p *Post) DeleteBy(actor User) { p.DeletedBy = actor.ID }
```

### Rule 3: Abstraction Boundaries

Every abstraction layer must belong to a single architecture layer. An abstraction cannot span multiple architecture layers.

**Evaluating abstractions:**
- Does this object depend on objects from a higher layer? → Extract or refactor
- Does this object's responsibility match its architecture layer? → Move if not

### Rule 4: Minimize Inter-Layer Connections

Fewer connections = looser coupling = better testability and reusability.

**Good layering:**
```go
func (s *Posts) Publish(ctx context.Context, actor User, id int64) error {
    post, err := s.store.Find(ctx, id)
    if err != nil { return err }
    if !CanPublish(actor, post) { return ErrForbidden }
    if err := post.Publish(); err != nil { return err }
    return s.store.Update(ctx, post)
}
```

**Caveat — Architecture Sinkhole:** If you reduce connections to minimum (each layer only talks to adjacent layer), you may create objects that just proxy data through layers with no modification.

```go
// BAD: a separate abstraction that adds no independent responsibility.
func (s *FindPostService) Find(ctx context.Context, id int64) (Post, error) {
    return s.store.Find(ctx, id)
}
```

## Mapping Go Components to Layers

### Presentation Layer

**Purpose:** Handle user interactions, present information

**Includes:**
- Handlers (HTTP request/response)
- React components (UI rendering)
- Channels (WebSocket connections)
- Mailboxes (inbound email)
- Jobs (**internal inbound** entry points — same design rules as handlers, with explicit job validation and execution authority)
- API serializers
- Form objects (user input handling)
- Filter objects (request parameter transformation)
- Presenters (view-specific logic)

**Primary concerns:**
- Request parsing and validation
- Authentication
- Response formatting
- User interface logic
- Building execution contexts for units of work (web requests and background jobs alike)

### Application Layer

**Purpose:** Organize domain objects for specific use cases

**Includes:**
- Service objects (business operations)
- Policy objects (authorization rules; shared operations enforce them for every caller)
- Interactors/Commands
- Mailers, deliveries, notifiers (notification orchestration; the delivery adapters they use — SMTP, push drivers — are infrastructure)

**Primary concerns:**
- Orchestrating domain objects
- Transaction boundaries
- Use-case specific logic

**Sub-layers:** The application layer is sometimes split into **Business** (framework-agnostic rules — e.g., policies) and **Services** (implementation-aware orchestration — e.g., mailers, deliveries). Either way, the invariant holds: mailers always belong to a layer above the domain layer — they are neither presentation (no user interaction) nor infrastructure (no transport details).

**Warning:** This layer is often overused. Don't strip all logic from models into services (anemic models anti-pattern).

### Domain Layer

**Purpose:** Entities, rules, invariants, application state

**Includes:**
- Models (business entities)
- Value objects (immutable concepts)
- Domain events
- Query contracts near their consumers; SQL query and repository implementations in persistence adapters
- Composition (shared behaviors)
- Configuration structs (schema objects; their data sources — ENV, YAML, credentials — are infrastructure)

**Primary concerns:**
- Business rules and invariants
- Entity relationships
- Data transformations
- Domain-specific calculations

### Infrastructure Layer

**Purpose:** Supporting technologies

**Includes:**
- pgx/database/sql (database access)
- API clients (external services)
- File storage adapters
- Message queue adapters
- Mail/notification delivery adapters (SMTP, push drivers)
- Configuration sources (ENV, YAML, credentials)
- Cache implementations

**Primary concerns:**
- Persistence
- External communication
- Technical implementations

## Applying Layered Architecture in Practice

When designing or refactoring code:

1. **Identify the architecture layer** the code belongs to
2. **Check dependencies** — does it depend on higher layers?
3. **Apply specification test** — do tests verify appropriate responsibilities?
4. **Extract if needed** — move code to the correct layer

## Common Layering Mistakes

| Mistake | Problem | Solution |
|---------|---------|----------|
| Ambient actor in domain types | Hidden dependency on presentation context | Pass as explicit parameter |
| Request in services | Service depends on HTTP layer | Extract value object from request |
| Mailer in callbacks | Domain calls upward into the application layer | Move the notification to the service/delivery layer |
| SQL in handlers | Presentation doing infrastructure work | Use the owning store or a focused query |
| Business logic in views | Presentation doing domain work | Use presenters or model methods |
