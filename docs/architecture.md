# Architecture

Use the bundled [data-systems-architecture](../.agents/skills/data-systems-architecture/SKILL.md) with `layered-go` to review invariants, transaction boundaries, retries, and recovery. Trace the authoritative write and external effect separately; load only the relevant chapters. The complete skill works without its source PDF or network access. See [skill sources and maintenance](skills.md#data-systems-architecture).

Go owns URLs, authentication, authorization, and data. Gonertia sends a React page name and props: the first request returns HTML, and subsequent Inertia visits return JSON. In production, the frontend and backend share one origin. There is no separate public JSON API or React Router.

## Boundaries

- `cmd/server`: configuration, database pool, session store, HTTP server, and signal-driven shutdown.
- `cmd/worker`: River queues, a dedicated LISTEN/NOTIFY connection, mail/AI adapters, and graceful draining.
- `cmd/manage`: migrations and user creation through the shared `accounts` package.
- `internal/accounts`, `internal/notes`: feature-specific types, validation, and GORM queries. Simple CRUD does not need empty service/repository/domain layers.
- `internal/httpapp`: routes, strict decoding, CSRF, HTTP authorization, Inertia, and session flash.
- `internal/database`: one pgx pool with SQL/GORM adapters; Goose and River migrations under a deployment lock.
- `internal/notemail`, `internal/assistant`: explicit operations, state records, transactional enqueue, and external adapters.
- `internal/background`: typed River workers, failure cleanup, and periodic reconciliation.
- `internal/observability`: a Prometheus exporter with bounded label cardinality.
- `web/src`: pages, components, TypeScript DTOs, and styles. `web/src/types.ts` holds shared types and the single `InertiaConfig` augmentation.

When an operation coordinates multiple records and invariants, extract a named use case with an explicit transaction. Move authorization into the operation's shared contract if callers extend beyond HTTP. External effects need a delivery guarantee after commit; an ordinary goroutine does not provide one. River's `InsertTx` stores delivery intent alongside operation data. See [background jobs and email](background.md) and [AI](ai.md).

## HTTP contract

| Route | Access and behavior |
| --- | --- |
| `GET /login` | Sign-in form |
| `POST /login` | Verify password, rotate session ID and CSRF token, redirect to `/` with 303 |
| `POST /logout` | Revoke the login grant, destroy session data, clear Inertia history, redirect to `/login` with 303 |
| `GET /` | Sign-in required; up to 50 latest notes belonging to the current user |
| `POST /notes` | 1–2000 Unicode characters after trimming; owner comes from the session |
| `DELETE /notes/{id}` | Delete the user's own note; another user's or missing note returns 404 |
| `GET /admin` | Recheck `admin` in the database; show the user count |
| `GET /tools` | Show the user's email/AI operations and current feature availability |
| `POST /tools/email` | Sign-in required; email only the current user; one active operation |
| `POST /tools/assistant` | Sign-in required; 1–500 characters; one active operation |
| `GET /metrics` | Available only with a configured Bearer token; no session authentication |
| `GET /health/live` | Process responds |
| `GET /health/ready` | PostgreSQL responds within two seconds |

Mutations accept JSON limited to 16 KiB and require `X-CSRF-Token`. Unknown fields and trailing JSON are rejected. File uploads would need a separate multipart contract and limits. Precognition requests are rejected because validation-only handling is not implemented.

Field errors are stored in session flash, followed by a 303 redirect and a normal GET containing `errors`. Error fields are flat; named error bags are not supported automatically. Invalid JSON returns 400, denied access returns 403/404, and unexpected failures return 500. `net/http.CrossOriginProtection` supplements session CSRF validation.

Public IDs serialize as strings. Password hashes, cookies, and private configuration values never enter props. CSRF uses `Always`, including partial reloads. The user is a page-specific prop, not a global variable. Every subsequent request checks access again.

## Sessions and operations

SCS stores sessions in PostgreSQL with a 24-hour absolute lifetime and a two-hour idle timeout. Cookies use Secure in production; HttpOnly and SameSite=Lax apply in every environment. Protected requests also require a nonexpired `login_sessions` grant, created only after password verification and token rotation. The grant stores a SHA-256 token digest, user ID, and absolute expiry. Logout deletes it before destroying SCS data. A concurrent SCS commit may restore old flash/data, but cannot restore authorization. Requests authorized before logout may finish; new requests with that cookie cannot authenticate. Admin requests read the current role from the database.

Sign-in removes expired grants and the previous token's grant in the same transaction that creates the new grant. Database errors fail the request. A failed SCS commit can leave an unreachable grant; its absolute expiry limits retention and the next sign-in reclaims it. Deleting an account removes its grants through a foreign key.

Sign-in is limited to ten attempts per minute per direct connection address; the limiter holds at most 4096 addresses. Forwarded headers are not trusted as the source IP. A reverse proxy may make the address shared: configure limits at the public ingress and consider a shared limit across application instances.

Production ignores the Vite hot file and uses the build manifest and asset version. A version mismatch makes Gonertia return 409 for a full reload. CSP allows same-origin scripts and blocks inline JavaScript by default. Errors are logged as safe categories and SQLSTATE values, without raw driver messages containing user values. An HTTP-boundary panic returns 500 and logs a stack trace without the panic value.

## Extending the application

`layered-go` retains the design method of `layered-rails` while replacing framework mechanisms with Go contracts. Four conceptual layers do not require four directory trees. For Inertia features, start with `inertia-go-architecture`, then load the relevant skill. Add deferred pages, uploads, shadcn, new notification channels, or other features together with error handling and tests; skill examples do not mean those features are installed.

## Persistence and background effects

GORM records and public DTOs are separate: the password hash is not a field of `accounts.User`. Background-operation projections select only the public fields explicitly mapped to DTOs. Do not decode HTTP input directly into ORM records. GORM query fragments, table names, and sort expressions belong to the server; user values use bind parameters.

GORM's `Transaction` supplies the shared `*sql.Tx` for record changes and River enqueue. External SMTP/model calls happen after commit and an atomic state claim. The template adds no business callbacks, AutoMigrate, or lazy associations. Pure domain behavior does not depend on GORM or HTTP.
