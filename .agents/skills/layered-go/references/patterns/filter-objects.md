# Filter Objects

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns
- Security Considerations
- Go tools

## Summary

Filter objects transform datasets based on user-provided parameters. They belong to the presentation layer, consuming request parameters and applying transformations to domain collections.

## When to Use

- Request parameter filtering/transformation
- Search interfaces with multiple optional filters
- Sorting and pagination logic
- Parameter sanitization before domain layer

## When NOT to Use

- Domain-specific queries (use query objects)
- Form submissions (use form objects)

## Key Principles

- **Presentation layer** — consumes user input, not domain logic
- **One filter per interface** — avoid universal filter objects
- **Can use query objects internally** — filters orchestrate, queries implement
- **Explicit inputs** — use typed fields and allowlisted values

## Implementation

### Typed query parameters

Use `net/url`, `strconv`, and an explicit allowlist; no Rubanok-style DSL is needed.
Filter parsing belongs to the HTTP boundary, while SQL lives in the query/store.

```go
type PostFilter struct { Status string; Limit int }
func ParsePostFilter(values url.Values) (PostFilter, error) {
    filter := PostFilter{Status: values.Get("status"), Limit: 50}
    switch filter.Status {
    case "", "draft", "published":
    default: return PostFilter{}, ErrInvalidStatus
    }
    if raw := values.Get("limit"); raw != "" {
        n, err := strconv.Atoi(raw)
        if err != nil || n < 1 || n > 100 { return PostFilter{}, ErrInvalidLimit }
        filter.Limit = n
    }
    return filter, nil
}
```

### Usage

Parse once, authorize scope, then pass the typed filter into the query. Reject malformed
or unknown enum values rather than silently dropping the filter. Keep defaults at this
configuration/input boundary.

### Filter vs Query Object

| Aspect | Filter | Query |
| --- | --- | --- |
| Layer | Presentation | Infrastructure implementation, consumer-owned contract |
| Input | Untrusted URL parameters | Validated values and authorized scope |
| Purpose | Parse UI choices | Retrieve a bounded result |
| Location | HTTP feature package | Owning persistence package |

Shared filters require identical semantics, not merely similar field names.

## Anti-Patterns

### Universal Filter Object

A universal map of arbitrary columns exposes storage details and bypasses meaning. Define a typed filter for this screen. Reuse only filters with the same input and query semantics.

### Filtering in Handler

A handler can parse a small filter directly. Extract the parser when combination rules or reuse justify it; keep SQL and query execution in the store.

## Security Considerations

Allowlist sort identifiers and enum values; SQL placeholders bind values, not identifiers. Bound page size and query duration. Apply the authorized tenant predicate independently of user filters, and test that no filter can remove it.

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
