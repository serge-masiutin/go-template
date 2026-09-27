---
name: inertia-go-testing
description: Test Go Inertia responses with testing, httptest, Gonertia assertions and PostgreSQL integration tests. Verify components, safe props, CSRF, redirect/flash lifecycle, partial/deferred loading and authorization.
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


# Inertia Go Testing

Preserve the upstream request-level checks: component, props, absence of private data, flash after redirects, deferred metadata and partial reload behavior. Use Go's test stack instead of RSpec/Minitest.

## Setup

Read `internal/httpapp/integration_test.go` and `docs/testing.md`. Integration tests exercise the real middleware/session/adapter stack with isolated PostgreSQL schemas. Browser tests exercise the built React client.

## Assertions

Gonertia v3 exports `AssertFromReader`, `AssertFromBytes`, `AssertComponent`, `AssertProps`, `AssertFlash`, `AssertDeferredProps`, `AssertMergeProps` and history assertions. Check the installed source before using another helper; do not invent Rails-style matchers.

```go
page := inertia.AssertFromBytes(t, responseBody)
page.AssertComponent("Home")
```

For selective checks, decode a typed response or inspect known keys. `AssertProps` uses exact equality; JSON numbers decoded into `any` are float64. String IDs avoid precision loss. Exact public DTO key checks are useful for preventing leaks; avoid brittle equality on unrelated whole-page metadata.

## Redirects and Flash

After POST/PATCH/DELETE, first assert 303 and a safe Location. Follow the GET with the same cookie jar, then assert validation/flash. Assert that one-time data is consumed. The redirect response itself has no rendered page props.

## Shared Props

Verify CSRF is present on relevant page responses and private account fields are absent. Test two independent clients to detect cross-user shared state. Per-request values must not be globally shared on the Inertia engine.

## Deferred Props

On the initial response assert registration in `deferredProps` and absence of the deferred value. Send a second partial request and assert the value and query execution. Also revoke access between requests and verify the second request is denied.

## Partial Reloads

Set `X-Inertia: true`, the current `X-Inertia-Version`, `X-Inertia-Partial-Component` and `X-Inertia-Partial-Data` (or `...-Except`) as appropriate. The component must match. Verify nonrequested expensive loaders do not run and required Always props remain.

A stale asset version intentionally returns 409 with `X-Inertia-Location`; a test that omits the version can accidentally test version negotiation rather than its handler.

## External Redirects

For an Inertia request, `engine.Location` returns 409 plus `X-Inertia-Location`. Verify the destination is server-owned or allowlisted. Test a normal browser visit separately when its behavior matters.

## What to Test and Avoid

Test this project's CSRF wiring: Gonertia does not supply it automatically. Test anonymous and unauthorized direct requests, current role revocation, owner-scoped writes, malformed JSON, errors across redirects and asset-version behavior.

Do not duplicate the entire Inertia library suite, assert private implementation details or mock away the adapter when testing its integration. Keep domain-rule matrices below HTTP and SQL/transaction correctness against PostgreSQL.

## Related Skills

`inertia-go-controllers`, `inertia-go-forms`, `inertia-go-pages`, `go-serialization`.

See [request test patterns](references/http-tests.md) for protocol setup.

## Deferred Failure and Retry

For any deferred feature, force its query deadline and a query failure. Verify the
actual response and UI contract: Gonertia v3.0.0 returns Render errors, not rescuedProps.
An unexpected failure must reach the error boundary and visible HTTP/network error UI,
not a perpetual fallback or empty success. A deliberate recoverable typed result must
be logged and tested as its own branch. Test retry with the same filter and after
access is revoked; retries must not recover another actor's data. Check cleanup of
error listeners when navigating away. Field-validation onError is not a 500 handler.
