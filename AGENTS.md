# Go Template: agent contract

- Write project documentation, code comments, technical messages, and skill instructions in English. Respond to the user in their requested language. `clear-writing` has complete Russian and English editions and selects the guide by the target text's language.
- Read related files, `go.mod`, `package-lock.json`, and the relevant guide before editing. Update a contract and all its consumers together.
- Stack: Go `net/http`, Gonertia v3, GORM generics/pgx, Goose, SCS, River, go-mail, Genkit; React, TypeScript, Vite, Tailwind, Storybook. The server owns routes, access, and data; React owns presentation and local interactions.
- Use feature packages under `internal`; do not create a service, interface, and repository for every endpoint. See `layered-go` and `docs/architecture.md`.
- Decode strict, size-limited JSON at the HTTP boundary. Derive ownership from the authenticated user, not the request body. Parameterize SQL; run transactional queries through the same `*sql.Tx`/GORM transaction and enqueue River jobs with `InsertTx` inside it.
- Do not enable AutoMigrate, SQL debug logs, or business effects in ORM hooks. Mail/AI run through River; payloads contain IDs, workers reestablish scope, and external calls run outside transactions.
- Every mutation sends `X-CSRF-Token`. Form errors use session errors plus 303; invalid JSON returns 400; denied access returns 403/404. Precognition is unsupported and rejected.
- Authorize every request, including partial/deferred requests. Keep the user in request context, not global props. Gonertia loaders may run concurrently: a pool can be shared, a Tx/Conn cannot.
- DTOs enumerate public fields. JSON IDs are strings. Update Go JSON and TypeScript together; types are not generated automatically.
- Do not replace errors with empty results. Log safe context and categories; exclude passwords, DSNs, session/CSRF tokens, and user content.
- Preserve original EM skills; their files are checked with SHA-256. Put project constraints here. Record adaptations and sources in `docs/skills.md`.

## Context

| Task | Guide | Skills |
| --- | --- | --- |
| Go architecture | `docs/architecture.md` | `layered-go` |
| Data invariants, transactions, retries, and failures | `docs/architecture.md` | `data-systems-architecture`, then the relevant `layered-go` guidance |
| Libraries, ORM, queues, AI | `docs/stack.md`, `docs/background.md`, `docs/ai.md` | `layered-go` and its `references/installed-stack.md` |
| Pages, routes, forms | `docs/architecture.md` | `inertia-go-architecture`, then the relevant `inertia-go-*` |
| JSON and TypeScript | `docs/architecture.md` | `go-serialization`, `inertia-go-typescript` |
| Components and Storybook | `docs/development.md` | `tailwind-best-practices`, `sb-hub`, then the relevant `sb-*` |
| Testing | `docs/testing.md` | `inertia-go-testing` |
| Setup and infrastructure | `docs/development.md`, `docs/deployment.md` | `inertia-go-setup` |
| Slow startup | `docs/development.md` | `go-boot-profiling` |
| README and writing | `README.md` | `good-readme`, `clear-writing` |
| Books and documents into skills | `docs/skills.md` | `book-to-skill` |

## Validation

Run commands through `mise exec --`. For code and configuration changes, run `bin/ci`; for HTTP/SQL/session contracts, also run integration tests and `bin/test-browser` as described in `docs/testing.md`. Check changed behavior and failure cases. Report completed and omitted checks. Do not leave stubs or concealed errors.
