# Gonertia Configuration

## Initialization

Construct one Inertia instance at startup. Configure Vite from the build manifest or
an explicit development hot file. Do not infer production mode from whichever file
happens to exist. Asset versions must change when the build changes.

```go
i, err := inertia.New(rootHTML,
    inertia.WithFlashProvider(provider),
    inertia.WithEncryptHistory(),
)
if err != nil { return err }
```

## Shared Data

Immutable app constants may use `ShareProp`. User, authorization, locale, CSRF, and
request-specific data use `inertia.SetProps` on the request context. Concurrent requests
must never mutate a shared instance's user-specific state.

## Middleware Order

Session loading wraps request authentication/CSRF and Gonertia middleware. The session
must be available before the flash provider executes. Database errors must fail the
request; validate your provider before relying on redirect persistence.

## Flash and Validation

The SCS-backed provider keeps validation errors, flash messages, and history flags
separate. Read actual `internal/httpapp/flash.go` before changing serialization. Consume
messages once. Never store passwords or entire form bodies as flashed input.

## History Encryption

`WithEncryptHistory()` enables browser-history encryption. `ClearHistory` on logout or
identity change is transported across the redirect. This does not revoke sessions or
remove the need for server-side authorization and no-store responses.

## Assets and Versions

Vite uses `web/src/app.tsx` as entrypoint and `web/build/manifest.json` as its manifest.
Use the manifest checksum as production asset version. The dev hot file belongs in
ignored `tmp/`; production must never honor a stale hot file. Build and check the app
without a running Vite process before publishing the image.

## Errors and Verification

Verify first HTML load, Inertia JSON navigation, partial reload, validation redirect,
external Location, role revocation, and logout. Check the pinned Gonertia code when an
API differs from the upstream Rails guide; no Rails initializer exists in this project.
