# Query Objects

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns
- Go tools

## Summary

Query objects encapsulate complex, context-specific database queries. They separate persistence concerns from domain models and provide reusable, testable query logic.

## When to Use

- Complex queries with multiple JOINs or subqueries
- Context-specific reports with their own result shape
- Queries reused across multiple handlers/services
- Queries that need parameterization

## When NOT to Use

- Simple queries (keep the owning store method)
- Single-condition filters (use a parameterized predicate)

## Key Principles

- **Simple queries in stores** — do not add a query abstraction for every lookup
- **Complex queries in query objects** — multi-condition, context-specific queries
- **Parameterized SQL** — avoid raw user values or identifiers in SQL strings
- **Explicit wiring** — construct a query with its declared dependencies

## Implementation

### Explicit query type

Keep simple operations on the store that owns persistence. Extract a named query for a
reused report with independent parameters and result shape. Query implementations are
infrastructure; pure domain types do not import pgx. Do not port an Active Record
relation chain into a homemade Go ORM.

```go
type RecentPosts struct { pool *pgxpool.Pool }
func (q *RecentPosts) ByAuthor(ctx context.Context, authorID int64, limit int) ([]Post, error) {
    rows, err := q.pool.Query(ctx,
        "SELECT id,title FROM posts WHERE author_id=$1 ORDER BY id DESC LIMIT $2", authorID, limit)
    if err != nil { return nil, err }
    defer rows.Close()
    posts := make([]Post, 0)
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.ID, &post.Title); err != nil { return nil, err }
        posts = append(posts, post)
    }
    return posts, rows.Err()
}
```

### Atomic vs Complex Queries

A single authorized lookup can stay a store method. A reporting join, reusable search,
or aggregation can become a focused query, with integration tests against PostgreSQL.
The [sales-report extraction](../../examples/query-to-query-object.md) preserves the
original reporting scenario with parameterized SQL and an explicit result type.

### Composition and ordering

Use typed filter values and fixed SQL alternatives. Parameterize values; allowlist
identifiers such as sort columns. Define one explicit final ordering and bound collections.
Do not infer a table from a type name, concatenate raw request values into SQL, or
hide authorization in a query constructor.

## Anti-Patterns

### Conflicting Query Composition

Do not combine fragments that silently replace ordering or duplicate joins. Define a single final ordering and test filter combinations.


### Context-Specific Queries on Domain Types

A sales report depends on persistence and reporting context. Put its SQL in a focused query, not a pure `Order` method.


## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
