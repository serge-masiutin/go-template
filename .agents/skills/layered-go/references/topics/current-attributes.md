# Current Attributes → Explicit Execution Context

## Summary

The original chapter separates execution context from domain state. Go has no request-local `Current` singleton: goroutines are not thread-local requests, and a package variable races across users. Carry business inputs explicitly.

## Layer Placement

HTTP middleware establishes identity. Application operations receive the actor/tenant they require. `context.Context` carries cancellation, deadlines and request correlation into I/O; it must not become a hidden bag of business dependencies.

## Key Principles

- No global current user, tenant, transaction, or request.
- Typed middleware context keys stay at the HTTP boundary; extract identity before calling an operation.
- Persisted jobs carry stable IDs, not request contexts, live pointers, or session tokens.
- Request logging metadata is different from a business decision about who may act.

## Implementation

### Replacing the Current Class

```go
type Actor struct { UserID int64; TenantID int64 }
func (s *Posts) Publish(ctx context.Context, actor Actor, id int64) error {
    post, err := s.store.FindInTenant(ctx, actor.TenantID, id)
    if err != nil { return fmt.Errorf("load post: %w", err) }
    if !CanPublish(actor, post) { return ErrForbidden }
    return s.publish(ctx, post)
}
```

`Actor` is an application input established by a trusted entry point. A client-supplied tenant ID alone is not authority. Store methods include the authorized tenant predicate.

### Setting in Handlers

The handler validates the session, reloads current account permissions and maps them to the operation's input. It passes `r.Context()` for cancellation. The operation does not retrieve a user from a context value or import `net/http`.

### Multi-Tenancy

Apply tenant scope in every query and conditional write. Test two tenants with overlapping resource IDs or similar data. Separate request pools and immutable inputs prevent cross-request mutation; connection reuse must not leak session-level database tenant settings.

### Background Jobs

Deserialize and validate a versioned payload; establish worker authority and current tenant state before executing. For user-initiated work, decide whether permission is evaluated at submission, execution, or both, and test revocation. A system job needs a named system capability, not a fabricated admin user.

## Acceptable Execution Metadata

### Audit Logging

The application records actor, action and resource identity from explicit inputs. Audit events requiring atomicity are saved in the business transaction. Security-sensitive audit data has an explicit retention/access policy.

### Request Logging

Correlation IDs can travel through context or a logger passed by the boundary. Do not put raw headers, cookies or request bodies into log attributes.

### Time Zone

Pass a validated timezone to presentation formatting. Keep persisted event timestamps in UTC. Locale and timezone must not be process globals that change per request.

## Anti-Patterns

### Ambient Identity in Models

A method reading `currentUser` or a context key depends on a hidden caller. Replace it with the smallest explicit argument: actor ID, tenant ID, or domain value.

### Context for Business Logic

`ctx.Value("discount")` makes required data optional at runtime. Use a typed command field; missing required input must fail at the boundary.

### Context in Operations

Accept context for cancellable work, never as service locator. Constructors receive collaborators. Unit tests should reveal dependencies without reconstructing a fake HTTP request.

### Testing Difficulties

Tests that mutate global Current require cleanup and cannot run safely in parallel. Remove the global rather than serializing the suite around it.

## Where Request Context Is Appropriate

Keep request identity lookup in middleware/handlers. Correlation, deadline, cancellation and tracing follow I/O; business identity remains explicit even when both travel together.

## Testing

Use parallel tests with distinct actors, cancellation tests for I/O, tenant isolation tests and worker permission-revocation tests. `go test -race` detects memory races, not missing SQL authorization predicates.

## Current vs Explicit Parameters

| Concern | Mechanism |
| --- | --- |
| Who may act | Explicit actor/capability |
| Which data is visible | Authorized tenant scope |
| When I/O stops | Context deadline/cancellation |
| Correlation | Boundary-owned request metadata |
| Background execution | Validated durable payload and worker context |
