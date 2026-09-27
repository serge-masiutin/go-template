# Policy Objects

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns
- Go tools

## Summary

Policy objects encapsulate authorization rules, determining whether a user can perform specific actions on resources. They belong to the application layer, used by authorized application operations; transport establishes identity and maps denials to responses.

## When to Use

- Authorization logic beyond simple ownership checks
- Role-based or attribute-based access control
- Reusable authorization rules across handlers
- Testing authorization independently from handlers

## When NOT to Use

- Simple ownership checks (`post.AuthorID == actor.ID`)
- Authentication (who is this user?)

## Key Principles

- **Separate from authentication** — authentication identifies, authorization permits
- **Keep in services layer** — policies are application-layer abstractions
- **Never leak to domain** — models live in already-authorized context
- **Predicate interface** — methods return true/false
- **Explicit wiring** — call the relevant policy; no reflection or naming-based lookup

## Implementation

### Explicit actor and record

Policies are ordinary functions or small types. Deny by default. Keep the policy free of
HTTP, global request state, and rendering. Enforce it in the shared application operation. For a simple query without an operation, enforce it in its entry boundary;
hiding a button does not authorize the request.

```go
func CanUpdatePost(actor User, post Post) bool {
    return actor.Admin || actor.ID == post.AuthorID
}
func (s *Posts) Update(ctx context.Context, actor User, post Post) error {
    if !CanUpdatePost(actor, post) { return ErrForbidden }
    return s.store.Update(ctx, post)
}
```

### Authorized collections and tenancy

Scope the SQL query by the authenticated account/tenant before selecting a record.
Check ownership again in conditional writes where a race could change authorization.
Never fetch every row and filter it in React. Role revocation must affect the next request.

### Reuse, context, and caching

Pass all policy inputs explicitly. Share rules across HTTP, workers, and CLI when they
perform the same operation; those entry points may establish authority differently.
Do not cache authorization across requests without an explicit invalidation contract.
Test owner, non-owner, admin, revoked role, and tenant isolation. See
[authorization extraction](../../examples/authorization-to-policy.md).

## Anti-Patterns

### Authorization in Models

Do not make pure entities read sessions, roles from HTTP context, or global actors. Domain invariants such as "paid orders cannot be edited" still belong to the entity even when an admin requests the change.

### Mixed Enforcement Layers

Choose one required enforcement boundary per use case. A shared operation authorizes every caller; handlers authenticate and translate its denial. UI capabilities are hints, never enforcement. Workers validate their explicit execution authority and current tenant scope.

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
