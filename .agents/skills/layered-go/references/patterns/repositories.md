# Repositories and Stores

## Summary

Repositories encapsulate persistence behind a small contract. Preserve the original distinction between a focused query and a data-access abstraction, but account for Go having no built-in Active Record model API: the starter uses GORM generics in feature stores, with pgx available for focused SQL and sessions. ORM records remain separate from public DTOs and pure domain values.

## When to Use

- A domain-named store owns SQL for one cohesive feature.
- An operation needs a substitutable persistence boundary with a small stable contract.
- Multiple persistence sources genuinely implement the same consumer semantics.
- A read model combines a bounded set of sources with explicit consistency rules.

## When NOT to Use

Do not create a generic CRUD repository, one interface per struct, or a forwarding wrapper around every pgx method. A concrete store can be used directly by a simple handler. SQL belongs in the store; pure entities do not import pgx.

## Key Principles

Define interfaces at the consumer only where substitution or a domain boundary helps. Use explicit methods such as `FindInTenant`, not `Find(any)`. Return typed values and errors. Pass context into cancellable I/O. State who owns transactions and row locks.

## Implementation

### Basic Repository

The following pgx example illustrates the consumer contract. For installed features, follow `internal/notes`: `gorm.G[record](db).Where("user_id = ?", userID)` with bound values, explicit projections and affected-row checks. Keep table/column/order expressions server-owned. Do not add AutoMigrate, business callbacks or generic CRUD interfaces.

```go
type PostReader interface {
    FindInTenant(context.Context, int64, int64) (Post, error)
}
type PostStore struct { pool *pgxpool.Pool }
func (s *PostStore) FindInTenant(ctx context.Context, tenantID, id int64) (Post, error) {
    var post Post
    err := s.pool.QueryRow(ctx,
        "SELECT id, tenant_id, title FROM posts WHERE tenant_id=$1 AND id=$2", tenantID, id,
    ).Scan(&post.ID, &post.TenantID, &post.Title)
    if err != nil { return Post{}, fmt.Errorf("find post in tenant: %w", err) }
    return post, nil
}
```

The interface above belongs to the consuming operation; it need not be declared alongside every concrete store. `pgx.ErrNoRows` can be translated to a feature-owned missing-record error at the persistence boundary when consumers must remain independent of pgx.

### Repository with Caching

Add a cache only after measurement. Define key scope, TTL, invalidation after commit and failure behavior. Cached authorization-sensitive data must not survive revocation without a specific freshness contract. Do not silently return stale data for every database error.

### Multi-Source Repository

Combining SQL and a remote provider needs explicit authority: which source owns each field, what partial failure means, and whether the result is consistent. A network call is not a database transaction. Return partial results only if partiality is part of the public type.

### Read Model Repository

A reporting DTO can differ from a write entity. Use explicit SQL projections, bounded pagination and stable ordering. Do not hydrate an entire aggregate to show a count or a list item.

### External Data Integration

Map provider vocabulary in its adapter. Validate the remote schema and preserve cause on errors. A domain type must not expose an SDK response object as its public contract.

## Repository vs Query Object

| Responsibility | Store/repository | Focused query |
| --- | --- | --- |
| Scope | Cohesive persistence operations | One report/search/result shape |
| Writes | Explicit commands with affected-row checks | None |
| Composition | Operation coordinates stores | Typed filters and fixed SQL |
| Transaction | Concrete transaction-bound implementation | Uses supplied transaction only when needed |

A complex query can be a method until it gains independent responsibility. Do not create a query layer solely because a SELECT exists.

## Transactions

A multi-record operation begins one transaction and passes transaction-bound persistence to each required write. Do not accidentally call the pool from inside an operation that expects the same transaction. Ensure rollback on every error and propagate commit failure. Persist external delivery intent in the transaction; deliver afterward. In this starter use `ORM.WithContext(ctx).Transaction` and pass its `*sql.Tx` to River `InsertTx`, as in `notemail.Service.Request`. A second pool operation or an enqueue after commit breaks atomicity. Goose SQL migrations own schema changes, including River migrations under the same deployment lock.

For concurrent state changes use conditional UPDATE/RETURNING or a row lock. An in-memory state check followed by an unconditional update is insufficient across processes.

## Testing

Use PostgreSQL integration tests for constraints, SQL mappings, empty results, pagination, tenant isolation, rollback and concurrent writes. A fake store verifies operation decisions but cannot prove SQL correctness. Use distinct connections to exercise database races.

## Anti-Patterns

### Thin Repository Wrapper

A wrapper exposing `Exec`, `Query` and `QueryRow` under new names adds no domain contract. Use pgx directly inside the concrete adapter, or expose the narrow operation a consumer needs.

### Repository with Business Logic

A store persists state and expresses atomic database conditions. Pricing, eligibility and allowed business transitions belong to domain/application code. Do not send notifications from an Insert method.

## File Organization

Start with `internal/<feature>/<feature>.go` and `store.go`. If a pure domain boundary becomes useful, move concrete persistence to an adapter package and keep the port at the consumer. The four conceptual layers do not require four empty directory trees.
