# Authorization

## Summary

Authorization determines whether an actor can act on a resource. Preserve the original distinction between authentication, permission rules and enforcement. Policies are application rules; a domain invariant remains true even for an administrator.

## Layer Placement

Entry points authenticate and establish actor identity. A shared application operation enforces its policy so HTTP, CLI and worker callers cannot bypass it. A simple read with no application operation enforces access in its entry boundary and uses a scoped query. React capability flags are presentation only.

## Key Principles

Deny unknown actions, keep inputs explicit, reuse one rule across callers, and test all permitted and forbidden branches. Do not introduce reflection-based policy lookup or a universal permission framework for a few functions.

## Implementation with Explicit Policies

### Basic Policy

```go
func CanUpdatePost(actor Actor, post Post) bool {
    return actor.TenantID == post.TenantID &&
        (actor.Admin || actor.UserID == post.AuthorID)
}
```

A policy is pure when it receives loaded facts. If facts require queries, load them in a bounded query/operation and preserve consistency for the subsequent write.

### Rule Design

Start with resource actions: view, list, create, update, delete. Give a distinct rule to an independent domain action such as transferring ownership or refunding a payment. Do not invent aliases that accidentally broaden permissions. A family of independently addressable actions may deserve its own resource.

### Handler Enforcement

Handlers translate `ErrForbidden` to 403 and a missing record in an authorized scope to 404. Choose the disclosure policy explicitly; do not reveal existence through a different response for another tenant's record. The handler must not pre-authorize once and assume later deferred requests inherit that decision.

### Scoping-Based Authorization

```sql
SELECT id, title, author_id FROM posts
WHERE tenant_id = $1 AND id = $2;
```

Keep the tenant predicate independent of user-controlled filters. Conditional writes should include the same scope and any concurrency guard. Check rows affected.

### View Authorization

Return a small `can` DTO to decide which controls to show. The mutation independently enforces permissions. Revoked roles take effect on the next full, partial or deferred request.

### Role-Based Access Control

Role values come from trusted storage, not form inputs. Avoid a broad admin bypass where tenant isolation is still required. Name administrative capabilities and test their scope.

### Attribute-Based Access Control

Ownership, membership, record state and time windows can determine access. Keep facts typed. If authorization depends on mutable state, evaluate and write under an appropriate lock or conditional update.

## Error Handling

Distinguish unauthenticated (login redirect for pages), forbidden (403), inaccessible/missing (404 according to the disclosure contract), invalid input (400 or form validation redirect), and unexpected infrastructure failures (500 with safe diagnostics).

## Anti-Patterns

### Authorization in Models

Do not make domain entities fetch sessions or roles. An order's invariant forbidding edits after settlement is domain behavior; whether this actor may request an edit is a policy.

### Multiple Enforcement Points

Duplicated handler and operation rules drift. Enforce once at the shared operation boundary; UI hints and scoped reads may reuse that rule. Document the entry boundary for operations that intentionally run under a system capability.

### Implicit Authorization

A hidden button, a route name, a previous page visit and an encrypted history entry grant no authority. Every request re-establishes access.

## Testing

### Policy Unit Tests

Use a table containing owner, outsider, admin, cross-tenant admin, revoked membership and invalid action. A refactor must preserve this truth table; see [policy extraction](../../examples/authorization-to-policy.md).

### HTTP Authorization Tests

Test anonymous, owner, non-owner, current admin and revoked admin against actual handlers and PostgreSQL. Include partial/deferred requests and direct mutation requests without visiting the page first.

## N+1 Authorization

Load relevant memberships in one query, use an authorized SQL scope, or use a cache restricted to one request. Do not cache a permission across requests without a documented invalidation/revocation contract. Measure query counts before adding a cache.
