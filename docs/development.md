# Development

## Tools and startup

Versions are defined in `mise.toml`, `go.mod`, `package.json`, and `package-lock.json`. After `mise trust` and `mise install`, run commands through `mise exec --`.

`bin/configure github.com/owner/app` updates the application and tooling module paths, Go imports, npm package name, and lockfile. It preserves skill sources and links to the original template. Run `bin/ci` after configuring the project.

`bin/setup` creates `.env` if it is missing, installs dependencies, starts PostgreSQL and Mailpit from `compose.yml`, applies migrations, and builds the frontend. Re-running it preserves `.env` and database contents. Run Docker Compose from the project root.

`bin/dev` starts Vite and two Air watchers: server and worker. Go rebuilds when relevant files change; React and CSS use HMR. Ctrl+C stops all processes; an unexpected watcher or Vite exit stops the others. Compilation errors appear in the terminal while Air waits for a fix. The supervisor reads `.env` at startup, so restart `bin/dev` after changing it.

The application defaults to `http://localhost:3000` and Vite to `http://localhost:5173`. Change the Go port with `HTTP_ADDR` and `PUBLIC_URL` in `.env`; Vite reads its CORS origin from `PUBLIC_URL`. To change Vite's own port, update `vite.config.ts`, hot-origin validation in `internal/httpapp/app.go`, and the related tests together.

## Create a user

```sh
mise exec -- bin/manage create-user --email you@example.com --admin
```

The CLI reads one password line from standard input. Passwords must be 12–72 bytes; email addresses are normalized, unique, and limited to 254 bytes. Public registration and password recovery are not implemented.

To hide password input in the terminal, use Bash:

```sh
bash -c 'read -r -s -p "Password: " starter_password; printf "\n" >&2; printf "%s\n" "$starter_password" | mise exec -- bin/manage create-user --email you@example.com --admin; unset starter_password'
```

Keep production passwords out of shell history, process arguments, and repository files. Omit `--admin` to create a regular user.

## Configuration

| Variable | Contract |
| --- | --- |
| `APP_ENV` | `development`, `test`, or `production`; empty means development |
| `DATABASE_URL` | Required PostgreSQL connection string |
| `HTTP_ADDR` | Listener address; default `127.0.0.1:3000` |
| `PUBLIC_URL` | HTTP origin without a path; production requires an explicit HTTPS origin |
| `DB_MAX_CONNECTIONS` | At least 2 connections; default 10; server and worker have separate budgets |
| `WORKER_CONCURRENCY` | 1–100 mail workers; default 4; AI uses one per process |
| `SHUTDOWN_TIMEOUT` | 1s–1m; default 10s |
| `METRICS_TOKEN` | Empty disables `/metrics`; otherwise at least 32 characters |

`PUBLIC_URL` sets Vite's CORS origin in development and validates deployment configuration. It is not a Host allowlist or an absolute-URL generator; restrict public hostnames at the reverse proxy. Shell wrappers load the local `.env`; compiled binaries read environment variables and do not load dotenv automatically.

## SQL and migrations

Goose reads SQL files from `internal/database/migrations`, embedded in the binary. Add the next numbered file with `-- +goose Up` / `-- +goose Down`; do not edit applied migrations. `mise exec -- bin/manage migrate` holds one advisory lock across application and River migrations on a dedicated PostgreSQL session. That session closes after migration and is not returned to the application pool. Each application migration is atomic; the entire migration sequence is not one transaction. Fix the cause of a failure and rerun the command; already applied versions remain recorded.

The first Go migration accepts the known migration history of earlier starter versions or creates the initial tables; unknown legacy history fails explicitly. The original DDL is retained in `internal/database/bootstrap`. Neither the server nor the worker runs migrations. The management CLI does not expose destructive down migrations; express a reversal as a new migration.

GORM uses its generics API and the same pgx pool through `stdlib.OpenDBFromPool`. Keep feature-specific queries in the feature's store or operation. Do not enable AutoMigrate, SQL debug logging, or business hooks. Operations that change multiple records open an explicit transaction; `internal/notemail/notemail.go` demonstrates using `InsertTx` in that transaction.

Connect to the local database with `docker compose exec postgres psql -U starter -d starter_development`. Stop it without deleting data with `docker compose stop postgres`.

## Components

Start with [DESIGN.md](../DESIGN.md) for token meanings, component contracts and composition rules. Styles and semantic tokens live in `web/src/styles.css`; base components live in `web/src/components`. `Field` owns visible labels and hint/error associations for native input and textarea controls; the containing Inertia form owns validation, request headers and submission state. Tailwind scans only `web/src` so documentation and skills do not affect CSS. `npm run storybook` serves the catalog on port 6006; `npm run build:storybook` builds it into `storybook-static`. Stories remain colocated as `web/src/components/*.stories.tsx`. `npm run test:storybook` runs their interaction and accessibility checks in Chromium; see [testing](testing.md#component-and-catalog-checks). The catalog is excluded from the production image. Storybook does not create the application's Vite hot file.

shadcn is not installed. If you add it, use `shadcn-inertia` while preserving CSRF, actual JSON types, and CSP. For component work, use the original EM `sb-*` and `tailwind-best-practices` skills.

## Additional tools

- Email and River UI: [background jobs and email](background.md). Mailpit is at `http://localhost:8025`; River UI is at `http://localhost:8087` after enabling the Compose profile.
- AI configuration and tests: [AI assistant](ai.md).
- `bin/tool air -v` and `bin/tool govulncheck ./...` use the separate `tools/go.mod`. This lets development tooling evolve without changing runtime dependencies. After updating tools, run `go -C tools mod tidy`.
- `/metrics` requires `Authorization: Bearer <METRICS_TOKEN>`. The Prometheus exporter includes Go runtime/process metrics, the database pool, HTTP duration/status by route pattern, and River job states. It excludes request bodies, notes, questions, answers, and raw URLs. A queue-metrics query failure fails the scrape rather than reporting an empty queue.

Prometheus/Grafana and an OTLP collector are not installed. Connect the exporter to your monitoring system and monitor worker availability and queue age separately.
