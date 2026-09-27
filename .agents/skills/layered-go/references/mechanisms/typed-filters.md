# Typed Request Filters

This replaces Rubanok's map/match DSL with explicit parsing and fixed SQL alternatives.

## Basic Usage

The HTTP boundary parses `url.Values` into an application/query-owned filter type. The store receives validated values and the authorized scope. Do not send raw request maps into SQL builders.

## Mapping Parameters

Parse each supported field once. Define defaults for omitted fields, reject invalid values and specify how duplicate query keys are handled. Composite filters validate relationships between fields, such as start time preceding end time.

## Matching and Failure Strategy

Use a switch over an enum or fixed allowlist. Unknown status/sort values produce a validation error; silently dropping a filter can expose a broader result set than requested.

## Nested Parameters

If nested filters are needed, define a strict JSON or URL-key contract. Go does not automatically decode Rails bracket syntax into validated nested objects. Do not invent an implicit parser.

## Sorting

SQL placeholders parameterize values, not column identifiers. Map public sort names to fixed SQL clauses. Use a stable tie-breaker such as ID and test ascending/descending behavior.

## Pagination

Bound page size. Offset pagination is simple; cursor pagination needs an explicit typed cursor, stable ordering and tamper-resistant/validated values. Authorization scope applies to every page independently.

## Form Integration

URL filters are navigation state, so a GET visit is appropriate. A submitted mutation uses the form contract and CSRF protection. Do not put mutation side effects in a query parser.

## Testing

Table-test omitted/invalid/duplicate inputs and combinations. PostgreSQL tests prove filtering, ordering, bounds and tenant predicates. See [filter objects](../patterns/filter-objects.md).
