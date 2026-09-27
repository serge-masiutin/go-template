# Loading Options at the Go Boundary

| Loading requirement | Gonertia v3.0.0 implementation | Contract |
| --- | --- | --- |
| Regular value | Direct prop or supported loader | Required fields present and authorized |
| Lazy evaluation | `func(context.Context) (any, error)` | No work for an excluded partial prop; not a serializer side effect |
| Optional | `inertia.Optional(loader)` | Omitted initially; explicit partial requests reauthorize |
| Deferred | `inertia.Defer(loader)` | Later authorized request; placeholder, error and ready UI |
| Deferred group | `inertia.Defer(loader, "group")` | Group related props without sharing a concurrent transaction |
| Once | `inertia.Once(loader)` | Verify actual adapter metadata; no assumed cross-navigation cache/expiry |
| Once with expiry | No direct expiry API in this pinned adapter | Choose an explicit freshness policy; expiry does not replace authorization |
| Always | `inertia.Always(value)` | Retained on partial requests; minimal shared allowlist |
| Merge/append | `inertia.Merge(...).Append(...).MatchOn(...)` | Stable keys, ordering and duplicate behavior |
| Deep merge | `inertia.Merge(...).DeepMerge()` | Explicit trusted response structure, not input authorization |
| Scroll | `inertia.Scroll(value, inertia.WithMetadata(metadata))` | Bounded query, stable navigation, nested data and reset behavior |
| Reset | Client visit `reset` + scoped query | Filter changes must not merge old rows into a new result |

See [prop types](../../inertia-go-controllers/references/prop-types.md) for exact signatures,
examples and supported combinations. Do not combine optional/deferred/merge methods
by analogy with the Rails adapter.

## Combining Options

For a deferred collection that later grows, authorize every request, group only related
props, use stable IDs for merges, and reset the prop when the filter changes. Query
pagination and UI merging are different responsibilities. Concurrent prop loaders
must not share a `pgx.Tx`, `*sql.Tx`, transaction-bound GORM handle, or single `pgx.Conn`.

Start rarely changing reference data as a normal prop. Only add a server cache after
measurement; include access/version/locale dimensions and explicit invalidation.
Clearing Inertia history does not expire server caches or revoke data already exposed.

Long AI generation runs through River after a transaction commits the run and job together.
`assistant.Service.Recent` returns owner-scoped persisted results; `/tools` polls only while
work is active. Serialization describes this output; it must not start model calls as
a hidden attribute conversion. Model/prompt metadata and queue IDs stay outside public DTOs.
