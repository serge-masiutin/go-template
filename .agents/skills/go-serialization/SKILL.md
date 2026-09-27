---
name: go-serialization
description: Serialize explicit Go JSON and Inertia boundaries without leaking persistence fields; keep entity, page, shared, loading, and TypeScript contracts separate. Use when changing response shapes or prop loading.
---

# Go Serialization

Keep reusable entity shapes, page-specific responses, shared data, loading policy,
and consumer types separate. Go renders Inertia pages through Gonertia; React consumes
explicit JSON contracts. Read actual `go.mod` and frontend types before changing them.

## Setup: Establish Whether Serialization Is Needed

Inspect the consumer, authentication, access scope, content limits, and existing DTO
before adding a new shape. The starter's `accounts.User` is already an explicit public
shape; its database password hash is never a response field. Do not create a parallel
serializer layer for every small struct. Use a distinct DTO when persistence and
presentation differ, or when reuse gives the shape a stable meaning.

## Entity Shape

List fields in a small typed struct. Conversion must not query, mutate, send notifications,
or consult globals. Pass presentation context explicitly when needed.

```go
type ProjectSummary struct {
    ID        int64     `json:"id,string"`
    Name      string    `json:"name"`
    UpdatedAt time.Time `json:"updatedAt"`
}

func summarize(project Project) ProjectSummary {
    return ProjectSummary{
        ID: project.ID, Name: project.Name, UpdatedAt: project.UpdatedAt,
    }
}
```

This is an illustrative feature, not an installed Project model. IDs cross JavaScript
boundaries as strings. `time.Time` marshals as RFC3339 with variable fractional precision;
normalize UTC at the boundary when the contract requires it. Fetch related records in
the query, not in per-item converters. Test collection growth for N+1.

## Page-Specific Response

Authorize and scope before serialization. Keep records, pagination, filters, and totals
distinct. For example, an authenticated handler can call
`projects.List(r.Context(), actor.ID, filter)` and then render
`engine.Render(w, r, "Projects/Index", inertia.Props{"projects": summaries})`.
The handler must handle both query and render errors through its HTTP error boundary.
Define cursor/navigation and maximum page size for a real collection; a bounded list
alone is not a complete pagination API. Return Inertia responses for Inertia visits;
use `encoding/json` for an independently designed API consumer.

## Shared Data

Use request-local `inertia.SetProp` for scoped data. Only process-invariant public data
may use global sharing; never mutate global props with a request's user or permissions.
The starter shares `csrfToken` through `inertia.Always` so partial reloads retain the
current token. That token is intentionally exposed to the same-origin frontend;
session cookies, password hashes, provider tokens, and arbitrary config are not.
Expose an explicit user DTO only on pages that require it, unless a real shared layout
contract requires otherwise. No translation generator is installed in this starter.

## Explicit Rendering and Naming

Use explicit component names and props; there is no inferred serializer lookup.
Name reusable DTOs by domain or consumer. Compose page-only metadata around entity
shapes instead of adding optional fields to every entity or introducing a universal
resource abstraction. An empty collection promised as an array must serialize as `[]`,
not a nil slice's `null`; initialize it in the query/DTO boundary and test the contract.

## Types and Schema

Read `inertia-go-typescript`. Define nullability, omitted fields, IDs, enums, date/time
units, order, and version compatibility. Go JSON tags and TypeScript properties must
agree. Test security-sensitive allowlists, including fields that must be absent.
TypeScript annotations do not validate arbitrary network input. There is no automatic
Go-to-TypeScript generator installed; update both sides and their contract tests.
If introducing a generator, record its version, inputs, outputs, command, and consumers.

## Loading Options

Read [loading options](references/loading-options.md) and the linked adapter reference for regular,
optional, deferred/grouped, once, merge, always, and scroll options supported by the
pinned Gonertia version. Loading belongs to the request boundary, never a DTO method.
Use one of Gonertia's exact loader function signatures, such as
`func(context.Context) (any, error)`. Reauthorize every deferred/partial request and
keep owner scope inside its query. Loaders can run concurrently: share a pool, not a
single `pgx.Conn` or transaction across concurrent loaders.

## Troubleshooting

| Symptom | Cause to investigate | Resolution |
| --- | --- | --- |
| Private fields appear | Persistence struct or provider payload serialized directly | Restore an explicit allowlist and negative access tests |
| Collection is slow | Per-record database reads | Query related data in batches and measure query count |
| Consumer misses field | JSON tags and TypeScript drift | Update producer and consumer together |
| Null silently becomes empty | Undocumented optionality | Define the contract at the boundary; reject malformed required input |
| Component not found | Name differs from React resolver | Fix the explicit render name and test initial HTML plus navigation |
| Types never generated | No generator configured | Use actual manual contracts; do not invent a command |
| Stale/foreign props | Missing access scope or unsafe caching | Reauthorize and scope every loader request |

## Gate

Verify exact sensitive-field allowlists, guest/owner/foreign/admin access as applicable,
malformed external input, empty and bounded collections, query behavior, TypeScript
compilation, and initial/partial/deferred response shapes. Use synthetic fixtures.
