# Explicit Policy Functions

This replaces the Action Policy integration guide while preserving rule design, scoping, composition, denial reasons, context and testing.

## Setup and Basic Usage

Use ordinary functions or a small policy type. No policy library is required. The authenticated actor and the resource are explicit arguments; constructor or method names do not trigger lookup.

```go
type PostPolicy struct { Actor Actor; Post Post }
func (p PostPolicy) CanUpdate() bool {
    return p.Actor.TenantID == p.Post.TenantID &&
        (p.Actor.Admin || p.Actor.UserID == p.Post.AuthorID)
}
```

The application operation calls the policy before mutation. The HTTP handler maps a typed forbidden result; React receives only the capabilities needed to render controls.

## Scoping

Read and write through authorized tenant predicates. Do not first load every row and filter in Go or React. A policy checking one record does not establish a safe collection query.

## Rule Aliases and CRUD Rules

Reuse a predicate only when permission semantics are identical. List/show/create/update/delete are a useful starting vocabulary; transfer/refund/cancel deserve distinct rules when they have distinct authority. A collection of related endpoints may represent a separate resource.

## Denial Reasons

A boolean is sufficient when callers only need permit/deny. If the UI needs a reason, define a bounded reason enum and avoid revealing another tenant's existence or sensitive policy facts. Do not return raw SQL/provider errors as permission messages.

## Policy Composition and Pre-Checks

Compose smaller predicates explicitly. Tenant membership must not be bypassed accidentally by an early admin check. Default-deny unfamiliar roles/actions. Keep broad overrides named and tested.

## Authorization Context

Pass actor, tenant and any relevant loaded facts. Optional context must be explicit in its type and rule; an absent required actor is a boundary error, not an anonymous permissive policy.

## Error Handling

Use the operation's typed/sentinel denial error. Authentication failure, forbidden action, missing authorized record and infrastructure failure are separate outcomes. Every deferred request rechecks current rights.

## Testing

Use table tests for owner, outsider, admin, cross-tenant actor and revoked membership. Request tests verify 303 login redirects, 403/404 disclosure policy, direct mutation access and partial/deferred requests. See [authorization](../topics/authorization.md).
