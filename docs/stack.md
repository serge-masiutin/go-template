# Libraries

This stack supports the Go + Inertia + React application in this repository. Application versions are pinned in `go.mod`, CLI tools in `tools/go.mod`, and containers in `compose.yml`. The choices below were checked against maintainer documentation and published Go modules on September 27, 2026; they are project decisions, not a universal ranking.

## Installed components

| Purpose | Library | Use in the template |
| --- | --- | --- |
| HTTP and Inertia pages | `net/http` + Gonertia 3.0.0 | Routes, middleware, pages, forms, and session flash |
| Persistence | GORM 1.31.2, PostgreSQL driver 1.6.3, pgx 5.11.0 | Accounts, notes, email operations, and AI requests; generics API and explicit projections |
| Migrations | Goose 3.28.0 | SQL files, a transaction per migration, and a shared deployment lock with River migrations |
| Background jobs | River 0.47.0 | Transactional mail/AI enqueue, workers, and periodic state reconciliation |
| Queue administration | River UI 0.19.0 OSS | Local queue interface through the `tools` Compose profile |
| Email delivery | go-mail 0.8.1 + `notemail.Service` | SMTP, TLS, a text template, and explicit delivery through River |
| Local email preview | Mailpit 1.31.2 | SMTP sink and web UI; captured messages are not relayed to the internet |
| AI assistant | Genkit Go 1.13.1 | Typed flow, JSON output, and bounded tool loop; Gemini and OpenAI adapters |
| Configuration | caarlos0/env 11.4.1 | Typed nested structs, defaults, and separate invariant validation |
| Sessions | SCS 2.9.0 + pgxstore | Server-side sessions in the same PostgreSQL pool |
| Metrics | prometheus/client_golang 1.24.1 + httpsnoop | HTTP latency/status, runtime, database connections, and queue states |
| Development reload | Air 1.67.4 + Vite | Go server/worker rebuilds, React HMR, and coordinated shutdown |
| Dependency auditing | govulncheck 1.8.0 + npm audit | Reachable vulnerable Go functions and npm dependency advisories |

The example workflow lives at `/tools`: email your latest notes to yourself or ask AI a question about them. Provider configuration stays on the server. Queue payloads contain only operation IDs; questions and answers belong to the user and are not exposed to other accounts.

## Why these choices

**HTTP.** Modern [ServeMux supports methods and route parameters](https://go.dev/blog/routing-enhancements), and Gonertia connects the server to React. [Chi is compatible with net/http](https://github.com/go-chi/chi) and can help with more complex route groups. This application does not currently need a second HTTP framework. Standard `net/http` middleware can be used without rewriting handlers.

**ORM.** [GORM recommends its generics API](https://gorm.io/docs/the_generics_way.html): operations return typed results and errors, and ambiguous `Save` and `FirstOrCreate` methods are excluded from that API. [User values use placeholders](https://gorm.io/docs/security.html); column names, table names, and sort expressions must remain server-controlled. The starter does not use AutoMigrate, persistence callbacks with business effects, automatic association saves, or whole-record ORM serialization.

[Ent](https://entgo.io/docs/getting-started/) is an alternative for projects that need a generated schema-first API and a large relationship graph. This template uses less generation and ordinary SQL migrations. [Bun](https://bun.uptrace.dev/guide/queries.html) supports SQL-first queries, but [its placeholders interpolate and escape values, including removing NUL](https://bun.uptrace.dev/guide/placeholders.html). This template uses GORM/pgx binding and explicit errors for invalid data.

**Queues.** [River supports sharing a transaction with GORM](https://riverqueue.com/docs/gorm). The operation record and job either commit together or roll back together. This closes the [commit-to-enqueue failure window described by Brandur Leach](https://brandur.org/job-drain) without a separate Redis service. SMTP and model APIs do not share a transaction with PostgreSQL, so the starter does not promise exactly-once external effects or automatically retry an uncertain send.

**AI.** [Genkit tools](https://genkit.dev/docs/go/tool-calling/) provide checked types and a tool-execution loop. The application uses the stable flows/generate API rather than `genkit/exp`; the [full-stack Agents API](https://genkit.dev/docs/go/agents/overview/) was beta when this stack was selected. [Google ADK](https://github.com/google/adk-go) and [Eino](https://www.cloudwego.io/docs/eino/overview/) provide broader orchestration, which this single assistant with a read-only tool does not need. Model capability checks, tool boundaries, and tests determine whether a provider fits the application.

**Observability.** Genkit can write full content to telemetry. The application enables neither its exporter nor reflection UI and rejects variables that would automatically open that channel. This follows [Genkit's distinction between traces and content logs](https://genkit.dev/docs/go/observability/telemetry-collection/). Application metrics and logs contain safe operational fields. Users can inspect email and answers through the application and local Mailpit.

## Extensions

Authorization currently uses owner-scoped queries and a current admin-role check. Introduce Casbin or OpenFGA when the product needs a policy or relationship-based access model. Realtime, S3/uploads, payments, vector search, registration, and password recovery are separate product tasks. Add alternative ORMs, agent runtimes, or queues only with a concrete need and a defined contract between them.
