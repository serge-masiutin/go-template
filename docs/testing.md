# Testing

## Core checks

```sh
mise exec -- npx playwright install chromium
mise exec -- bin/ci
```

This runs Go formatting, `go mod tidy -diff` and `go mod verify` for both modules, `go vet`, `go test -race`, TypeScript checks, a Vite build, Vitest, a Storybook build, Chromium story interaction/accessibility tests, design-system contract/index validation, and skill validation. PostgreSQL is not required. Install Chromium once after installing npm dependencies; Linux CI uses `npx playwright install --with-deps chromium`. Frontend source uses Prettier, and its formatting check is included.

## Component and catalog checks

`mise exec -- npm run test` exercises native field semantics and error relationships in jsdom. `mise exec -- npm run test:storybook` runs every story in Chromium through the Storybook Vitest addon, including keyboard interactions, the 320 px workspace layout, and accessibility checks. `parameters.a11y.test: "error"` makes axe violations fail the command. The separate browser config loads the same Tailwind processing, CSS and fonts as the catalog; the CssCheck story detects missing CSS.

`mise exec -- bin/design-system-check` validates contracts, source snapshots, links and generated indexes against [DESIGN.md](../DESIGN.md). It does not replace runtime tests or visual inspection. Browser-based axe checks complement manual keyboard and layout review; they do not prove every accessibility requirement.

## PostgreSQL and browser tests

Start the local database with `bin/setup`, then create a separate test database once:

```sh
docker compose exec postgres createdb -U starter starter_test
```

Use the existing database if it has already been created. Test commands refuse database names without the `_test` suffix. Each run creates its own schema and removes it on completion; it does not touch the development schema.

```sh
export TEST_DATABASE_URL='postgres://starter:local-development-only@127.0.0.1:5437/starter_test?sslmode=disable'
mise exec -- go test -race -tags=integration ./...
mise exec -- npx playwright install chromium
mise exec -- bin/test-browser
```

Run `bin/ci` first to build the current asset manifest. `bin/test-browser` creates an isolated schema and synthetic administrator, runs Playwright, and drops the schema. Port 3100 must be available. `APP_ENV=test` uses built assets and ignores the development hot file. On Linux, the browser may need system dependencies: `npx playwright install --with-deps chromium`.

Integration tests cover a delayed SCS commit racing with logout across independent handlers, login grant rotation/expiry/revocation, migrations, sign-in, one-time form errors, CSRF and its rotation, cross-site requests, strict JSON, version mismatch, partial props, note ownership, admin-role revocation, logout, and rejection of mutations during unsupported Precognition validation. The browser suite covers failed and successful sign-in, accessible field-error descriptions and clearing, note validation, creation, reload, deletion, narrow viewport overflow, and logout. JavaScript errors fail the test.

## Changing contracts

Test public behavior and failure cases. Check SQL, sessions, and authorization against real PostgreSQL. Type checks and browser flows complement Go tests. For a new deferred feature, add loading failure, retry, and renewed authorization checks; for uploads, test actual multipart payloads and limits.

GitHub Actions runs both suites and builds the Docker image. It also runs `bin/tool govulncheck ./...` separately. The core suite does not include load tests, production smoke tests on your hosting platform, or paid generation against live AI providers.

## Queues, email, and AI

`go test ./internal/assistant ./internal/mailing` uses synthetic data, local HTTP/SMTP servers, and a fake model to verify contracts. It covers tool-loop limits, the output schema, rejection of answers without reading notes, Gemini thought signatures, authentication headers, token budgets, and no HTTP retry after 503. It sends no real email and makes no paid generation calls.

The PostgreSQL suite also checks accepted legacy migration history, rollback when enqueue fails, concurrent requests, owner scope, duplicate execution, account deletion, a real River worker, failure cleanup, and reconciliation after a crash. Metrics tests cover Bearer access and the exclusion of user content.

Live-model answer quality has a separate [evaluation procedure](../evals/notes-assistant.md). HTTP fixtures do not establish model quality. `govulncheck` may report unused subpackages within modules: distinguish reachable symbols, imported packages, and module-only findings rather than combining them into one count.
