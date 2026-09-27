---
name: inertia-go-controllers
description: >-
  Go HTTP handlers for Gonertia and React: explicit rendering, regular/deferred/optional/merge/scroll
  props, shared request data, session flash, 303 redirects, validation errors and authorization.
  Use when loading records or serving Inertia pages. External navigation uses Inertia.Location.

---

## Go project contract

This is an adaptation of the original EM Inertia skill, preserving the React examples
and decision workflow. Read `docs/architecture.md` and the actual Go handlers first.
Use Gonertia v3 APIs from `go.mod`; server examples below replace the Rails adapter.
All mutations carry `X-CSRF-Token` from `usePage().props.csrfToken`. The starter accepts
strict JSON; uploads need a separate bounded multipart handler before using upload examples.
Validation errors cross a 303 redirect through the session flash provider. They are not
bare 422 JSON responses. Shared user data is request-scoped, never global `ShareProp`.
Examples describe features to implement, not routes or libraries already installed.


# Inertia Go Controllers

Server-side patterns for Go handlers serving Inertia responses.

**Before adding a prop, ask:**
- **Needed on every page?** → Request middleware and `inertia.SetProp`.
- **Expensive to compute?** → `inertia.Defer` after measuring the query.
- **Only needed on partial reload?** → `inertia.Optional`.
- **Stable reference data?** → `inertia.Once`, with the installed adapter's limitations.

**NEVER:**
- Share a user, session, or request-dependent closure through global `ShareProp`.
- Serialize persistence records with secret/private fields.
- Return flat validation messages without field keys.
- Invent Rails adapter methods or automatic type/resource discovery.
- Use a normal external redirect for an Inertia navigation.

## Render Syntax

Every handler names the page and passes its props explicitly. There is no controller-name
inference, automatic instance-variable serialization, or inherited base-controller state.

```go
if err := i.Render(w, r, "Users/Index", inertia.Props{"users": summaries}); err != nil {
    fail(w, r, err)
    return
}
```

Use an explicit empty props set for a static page. Load records through an authorized query,
map them to DTOs, and handle errors before rendering. The page name must resolve in the
React entrypoint; an unknown page is a contract error.

## Prop Types

Use `inertia "github.com/romsar/gonertia/v3"`; helpers are package functions, not methods
on a Rails constant. Verify the pinned adapter before using combinations.

| Type | Syntax | Behavior |
| --- | --- | --- |
| Regular | `inertia.Props{"key": value}` | Included when requested by the response |
| Lazy | `func() (any, error) { return query(ctx) }` | Evaluated when the prop is selected |
| Optional | `inertia.Optional(loader)` | Omitted initially; explicitly requested later |
| Deferred | `inertia.Defer(loader)` | Separate authorized request after initial render |
| Grouped defer | `inertia.Defer(loader, "analytics")` | Named deferred group |
| Once | `inertia.Once(value)` | Basic once-prop support; no invented expiry API |
| Merge | `inertia.Merge(items)` | Merge into client prop |
| Deep merge | `inertia.DeepMerge(value)` | Deep-merge supported shapes |
| Always | `inertia.Always(value)` | Included during partial reloads |
| Scroll | `inertia.Scroll(items, inertia.WithMetadata(metadata))` | Explicit pagination metadata |

Wrap expensive work in a loader. `inertia.Defer(expensiveQuery())` already executed the
query and cannot defer that work. Each follow-up request must authenticate and authorize.

### Deferred Props — Full Stack Example

Server defers slow data, client shows fallback then swaps in content:

```go
props := inertia.Props{
    "basic_stats": summary,
    "analytics": inertia.Defer(func() (any, error) { return analytics.Compute(r.Context()) }),
}
if err := i.Render(w, r, "Dashboard", props); err != nil { fail(w, r, err) }
```

```tsx
// Page component — child reads deferred prop from page props
import { Deferred, usePage } from '@inertiajs/react'

export default function Dashboard({ basic_stats }: Props) {
  return (
    <>
      <QuickStats data={basic_stats} />
      <Deferred data="analytics" fallback={<div>Loading analytics...</div>}>
        <AnalyticsPanel />
      </Deferred>
    </>
  )
}

function AnalyticsPanel() {
  const { analytics } = usePage<{ analytics: Analytics }>().props
  return <div>{analytics.revenue}</div>
}
```

## Shared Data

Use `inertia.SetProp` / `SetProps` on the request context. Shared user data is an explicit
small DTO, not a session object or the entire stored user. Immutable application constants
can use `ShareProp` during initialization; never mutate it for an individual request.

```go
ctx := inertia.SetProps(r.Context(), inertia.Props{"auth": AuthProps{User: userDTO}})
next.ServeHTTP(w, r.WithContext(ctx))
```

Page-specific fields stay on the page. Keep key names and optionality aligned with
`InertiaConfig.sharedPageProps` in TypeScript. Read
[configuration](references/configuration.md) for session and middleware order.

## Flash Messages

Go does not provide Rails flash automatically. Gonertia's `FlashProvider` is implemented
by `internal/httpapp/flash.go` over SCS sessions, stored in PostgreSQL. Keep field errors,
flash messages, and history-clear flags independent and consumed once.

```go
ctx := inertia.SetFlash(r.Context(), inertia.Flash{"notice": "Saved"})
i.Redirect(w, r.WithContext(ctx), "/users", http.StatusSeeOther)
```

Client flash keys must match the declared TypeScript flash shape. Do not put private
messages into global props. Verify the provider's storage and error behavior before
adding complex flash payloads.

## Redirects & Validation Errors

After create/update/delete, redirect with 303 (Post-Redirect-Get). Invalid form values
also redirect with field errors persisted across the request. Malformed JSON or types
produce 400. An Inertia form does not consume an ordinary API's bare 422 response.

```go
ctx := inertia.SetValidationErrors(r.Context(), inertia.ValidationErrors{
    "email": "Enter a valid email address.",
})
i.Redirect(w, r.WithContext(ctx), "/users/new", http.StatusSeeOther)
```

Keys must match input names exactly. Use a known local destination instead of trusting
a user-supplied `Referer`. Do not persist password input or arbitrary request bodies.

## Authorization as Props

Pass per-resource permissions only to control visibility; the server enforces them.
Use policy functions and ownership-scoped queries. Never accept actor IDs or roles from
the form as authority. Read [authorization](references/authorization.md) when adding `can` props.

## External Redirects (`Location`)

Use `i.Location(w, r, destination)` for navigation outside the Inertia application.
For an Inertia request it returns 409 with `X-Inertia-Location`; ordinary requests receive
a redirect. Validate the destination against the operation's trusted provider/allowlist.
Do not redirect to arbitrary form input.

```go
i.Location(w, r, checkout.URL)
```

## History Encryption

Use `inertia.WithEncryptHistory()` and `inertia.ClearHistory(r.Context())` for logout or
identity changes. Persist the clear flag over the redirect with the flash provider.
Encryption does not replace authorization, session revocation, or `Cache-Control: no-store`.
Read [configuration](references/configuration.md) for the complete boundary.

## Configuration

Read [configuration](references/configuration.md) for initialization, Vite assets,
versioning, flash, request-scoped props, and session middleware.

## Troubleshooting

| Symptom | Cause to inspect | Fix |
| --- | --- | --- |
| External redirect opens an Inertia error | Ordinary redirect to provider | Use `i.Location` |
| Field error disappears after redirect | Missing/broken FlashProvider or session middleware | Verify persistence and consume-once behavior |
| Wrong user's shared data | Per-request mutation of global props | Use request context |
| 403 on submit | Missing CSRF header | Pass current `csrfToken` through `X-CSRF-Token` |
| Empty or stale page | Wrong component name or partial reload keys | Check Go render call and React page contract |
| Once options missing | Rails-only expiry options copied | Use the supported Go API; recheck adapter source |

## Related Skills
- **Form error display** → `inertia-go-forms`
- **Flash toast UI** → `inertia-go-pages` (access) + `shadcn-inertia` (Sonner)
- **Deferred on client** → `inertia-go-pages` (`<Deferred>` component)
- **Type-safe props** → `inertia-go-typescript` or `go-serialization` (serializers)
- **Testing** → `inertia-go-testing`

## References

**MANDATORY — READ ENTIRE FILE** when using advanced prop types (`merge`,
`scroll`, `deep_merge`) or combining multiple prop options:
[`references/prop-types.md`](references/prop-types.md) (~180 lines) — detailed behavior,
edge cases, and combination rules for all prop types.

**Do NOT load** `prop-types.md` for basic `defer`, `optional`, `once`, or `always`
usage — the table above is sufficient.

Load [`references/configuration.md`](references/configuration.md) (~180 lines) only when
setting up Gonertia initialization for the first time or debugging configuration
issues. **Do NOT load** for routine controller work.
