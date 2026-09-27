# Go + Inertia Feature Patterns

This reference replaces the original Rails server examples while keeping their page/form/navigation/state-ownership decisions. It is a skill reference, not a second repository-wide instruction file.

## Page Data and Shared Data

Authenticate and authorize before rendering. Call `engine.Render(w, r, "PageName", inertia.Props{...})`; check its returned error. A page name must match the React resolver. Map loaded values to public DTOs.

Use `inertia.SetProp`/`SetProps` on the request context for shared user-specific data. Do not mutate `engine.ShareProp` per request. The starter shares CSRF and supplies user/notes explicitly on authenticated pages.

## Forms and Validation

Use the existing strict JSON decoder and explicit input struct. Reject unknown fields and malformed structure. Do not accept owner/admin fields merely because the database has them. The operation/store enforces business rules and authorization.

```go
ctx := inertia.SetValidationErrors(r.Context(), inertia.ValidationErrors{"name": "Enter a name."})
engine.Redirect(w, r.WithContext(ctx), "/profile", http.StatusSeeOther)
```

The configured flash provider persists errors through the redirect. Successful mutations also return 303. React Form/useForm supplies the `X-CSRF-Token` header; neither Inertia nor Gonertia automatically implements the project's CSRF mechanism.

## Navigation and External Destinations

Use Link/router for internal visits. Parse URL query state on the server and return explicit props. Use `engine.Location` for an external Inertia redirect; destinations must be server-owned or allowlisted. Normal external anchors and file downloads do not need client-side routing.

## Deferred and Partial Loading

A supported loader is `func(context.Context) (any, error)`. Use Optional for explicitly requested data and Defer for post-render requests. Each request repeats permission checks. Loaders may run concurrently; use the pool safely and avoid a shared transaction/connection or mutable request map.

Read [prop types](../../inertia-go-controllers/references/prop-types.md) for actual merge/once/scroll capabilities. Do not copy unsupported options from another server adapter.

## Local State and Persistent Layouts

React owns open dialogs, focus and local drafts. The server/URL owns persisted resources and shareable filters. A persistent layout preserves its own state across visits; it must not cache authoritative permissions independently.

## Flash and Notifications

Use `inertia.SetFlash` plus the provider, then consume `usePage().flash`. Declare the exact flash type if adding notices. One-time notifications are not ordinary shared props. The base starter currently uses field validation errors and has no toast library.

## Polling and Realtime

`usePoll` is valid when bounded periodic refresh meets the product need. Add SSE/WebSocket only for a concrete latency/load requirement and implement channel authorization, reconnect and cleanup. An event may trigger `router.reload({ only: [...] })`; it does not grant access. No realtime runtime is installed.

## Independent APIs

Use a separate endpoint for non-browser consumers, high-frequency autocomplete, downloads or streaming. Keep its auth, input/output schema and error/status contract explicit. Do not create an API copy of every Inertia page by default.

## Verification

Run the relevant Go request tests, TypeScript build and browser scenario. Check field errors across 303, CSRF failure, unauthorized direct requests, asset version negotiation and deferred permission revocation when those paths change.
