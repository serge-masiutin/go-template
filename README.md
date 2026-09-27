# Go Template

A Go + Inertia + React starter for building web applications in one repository. Go owns routes, authentication, and data; React renders pages and handles local interactions. Start with a working application that includes private notes, sign-in, background jobs, email, an optional AI assistant, and an admin page.

- **Server-driven pages.** Standard `net/http` routing and Gonertia connect Go handlers to React without a separate public API or client-side router.
- **PostgreSQL throughout.** GORM's generics API, pgx, Goose migrations, SCS sessions, and River jobs share the database. Operations and their jobs commit together.
- **Working examples.** Create and delete your own notes, email them to yourself, or ask the assistant a question about them. The server checks ownership on every request.
- **React development.** TypeScript, Vite HMR, Tailwind, local Martian Mono fonts, and a separate Storybook component catalog.
- **Local tools.** Air reloads the Go server and worker. Mailpit captures development email; an optional River UI shows the queues.
- **Checks and packaging.** Go race tests, PostgreSQL integration tests, Vitest, Playwright, Storybook builds, dependency checks, and a Docker image for the server and worker.

[Quick start](#quick-start) · [How it works](#how-it-works) · [Project layout](#project-layout) · [Development commands](#development-commands) · [Deployment](#deployment) · [Documentation](#documentation)

## Quick start

Install Git and [mise](https://mise.jdx.dev/), and start Docker with Compose. Go and Node versions are pinned in [mise.toml](mise.toml).

Create a repository with **Use this template** on GitHub, then clone your new repository:

```sh
git clone git@github.com:YOUR_ACCOUNT/YOUR_APP.git
cd YOUR_APP
mise trust
mise install
mise exec -- bin/configure github.com/YOUR_ACCOUNT/YOUR_APP
mise exec -- bin/setup
```

`bin/configure` sets the Go module paths, imports, and npm package name for your project. `bin/setup` creates `.env` if needed, installs dependencies, starts PostgreSQL and Mailpit, applies migrations, and builds the frontend. Re-running setup preserves your configuration and database volume.

Create an administrator:

```sh
mise exec -- bin/manage create-user --email you@example.com --admin
```

Enter a password of 12–72 bytes and press Enter. The command reads standard input; the terminal displays what you type. See [creating a user](docs/development.md#create-a-user) for a command that hides input. Omit `--admin` to create a regular user.

Start development:

```sh
mise exec -- bin/dev
```

Open **http://localhost:3000**, sign in, and create a note. **Admin** shows the user count. **Note tools** lets you email your notes to the local Mailpit inbox; enable an AI provider separately to use the assistant.

| Service | Local address | How to start it |
| --- | --- | --- |
| Application | `http://localhost:3000` | `mise exec -- bin/dev` |
| Mailpit inbox | `http://localhost:8025` | Started by `bin/setup` |
| Storybook | `http://localhost:6006` | `mise exec -- npm run storybook` |
| River UI | `http://localhost:8087` | `docker compose --profile tools up -d riverui` after setup |

PostgreSQL listens on `127.0.0.1:5437`; Mailpit accepts SMTP on `127.0.0.1:1025`. Development email stays in Mailpit. AI is disabled by default and requires a model and API key; see [AI configuration](docs/ai.md).

## How it works

The first page request returns HTML. Subsequent Inertia navigation requests return a page name and JSON props, which React renders. The Go server and built frontend use one origin in production.

```text
React form → POST /notes → validate JSON, session, and CSRF
           → save the note for the signed-in user
           → 303 redirect → GET / → render updated page props
```

Go handlers validate input, check access, and call feature packages under `internal/`. Public DTOs list the fields sent to React; their TypeScript counterparts live in `web/src/types.ts`. Database records are not serialized directly into pages.

Email and AI requests create an operation record and a River job in one database transaction. A separate worker loads the operation and performs the external call. These jobs use a single attempt and an explicit state transition before sending: a timeout may leave an unknown outcome, so the template does not automatically repeat an uncertain external effect. See [background jobs and email](docs/background.md).

Sessions live in PostgreSQL. Mutations require CSRF tokens, note queries enforce ownership, and admin requests recheck the current database role. These contracts are covered by integration and browser tests.

## Project layout

| Path | Purpose |
| --- | --- |
| `cmd/server` | HTTP server and application startup |
| `cmd/worker` | River workers for email, AI, and reconciliation |
| `cmd/manage` | Database migrations and user creation |
| `internal/accounts`, `internal/notes` | Feature types, validation, and persistence |
| `internal/notemail`, `internal/assistant` | Email and AI operations and adapters |
| `internal/httpapp` | Routes, middleware, Inertia responses, and session flash |
| `internal/database` | Shared database pool and embedded migrations |
| `web/src` | React pages, components, TypeScript contracts, and styles |
| `.storybook` | Component catalog configuration |
| `tests/browser` | Playwright application flows |
| `.agents/skills` | Project-local guidance for coding agents |

Keep simple features in cohesive packages. Extract an operation when it coordinates multiple records or invariants; add an interface when a real boundary needs one. The [architecture guide](docs/architecture.md) explains the request, persistence, and background-work contracts.

## Development commands

Run commands from the repository root:

| Command | Purpose |
| --- | --- |
| `mise exec -- bin/dev` | Run Vite and reload the Go server and worker |
| `mise exec -- bin/manage migrate` | Apply application and River migrations |
| `mise exec -- bin/worker` | Run the worker independently |
| `mise exec -- npm run storybook` | Browse and develop components |
| `mise exec -- bin/ci` | Check Go, TypeScript, formatting, tests, builds, and skill integrity |

`bin/ci` does not require PostgreSQL. Integration and browser tests use a separate database with a name ending in `_test`; follow the [testing guide](docs/testing.md) to create it and run those suites. [GitHub Actions](.github/workflows/ci.yml) runs both, checks Go vulnerabilities, builds the Docker image, and smoke-tests the server and worker.

## Deployment

```sh
docker build -t my-app .
```

The production image contains the built frontend and three Go binaries: server, worker, and management CLI. Run migrations before starting a release, then run the server and worker as separate processes with PostgreSQL and HTTPS configured. See [deployment](docs/deployment.md) for environment variables, commands, health checks, and shutdown behavior.

The starter includes CLI-managed accounts and an admin page showing the user count. Public registration, password recovery, file uploads, realtime features, and hosting infrastructure are extension points rather than installed features.

## Documentation

- [Development](docs/development.md): configuration, users, migrations, and local tools.
- [Architecture](docs/architecture.md): package boundaries, routes, sessions, and data contracts.
- [Libraries](docs/stack.md): installed dependencies and the reasons for each choice.
- [Background jobs and email](docs/background.md): transactions, worker behavior, SMTP, and River UI.
- [AI assistant](docs/ai.md): provider setup, tool boundaries, and evaluation.
- [Testing](docs/testing.md): unit, integration, browser, and AI checks.
- [Deployment](docs/deployment.md): containers, reverse proxy, process supervision, and backups.
- [Agent skills](docs/skills.md): the 31 bundled skills, their sources, and update procedures.

## License

Application code is licensed under [MIT](LICENSE). The project builds on [rails-template](https://github.com/serge-masiutin/rails-template); source attribution and third-party licenses are listed in [THIRD_PARTY.md](THIRD_PARTY.md).
