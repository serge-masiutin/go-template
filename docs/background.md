# Background jobs and email

`cmd/server` enqueues work and `cmd/worker` executes it. Both use PostgreSQL. `bin/dev` starts the server, worker, and Vite; run the worker independently with `mise exec -- bin/worker`. In production, the server and worker are separate processes from the same image.

## Transactions and execution

`notemail.Service.Request` and `assistant.Service.Request` create an operation record and call River's `InsertTx` inside one GORM transaction. River receives that transaction's `*sql.Tx`. GORM, SCS, and direct pgx access share one pool; the worker also uses a dedicated LISTEN/NOTIFY connection.

Job payloads contain only operation IDs. The types `note_email_v1`, `notes_assistant_v1`, and `reconcile_work_v1` are versioned. Do not change the meaning of queued payloads without supporting their old version. The worker reloads the operation, establishes its owner, and reads only that owner's notes. Account deletion removes dependent operations through a foreign-key cascade, preventing pending work from loading them.

Email uses `WORKER_CONCURRENCY` parallel workers; AI uses one per process. Additional processes increase total concurrency. Mail/AI jobs default to one attempt: River must not automatically repeat an external effect with an unknown result. A conditional `queued → sending/running` update precedes the external call. Re-executing a completed operation sends nothing; replaying a started operation requires investigation rather than another charge or message.

SMTP can accept a message before the connection or result write fails. A `failed` state can therefore include an unknown outcome; inspect the inbox or Mailpit before retrying. A user retry creates a new operation. For AI, a model or prompt-version change makes already queued work fail; a new request uses the current configuration.

The River timeout is six minutes; the rescue threshold is seven minutes. A periodic job runs once a minute and marks a stuck operation `failed` when its River job is terminal or has been removed. It does not repeat the external call. After a process crash, state may remain intermediate until a rescue/maintenance cycle runs. On shutdown, the worker stops taking jobs, waits for `SHUTDOWN_TIMEOUT`, then cancels unfinished work with an additional five-second limit.

## Email

`.env.example` routes SMTP to local Mailpit at `127.0.0.1:1025`; its web UI is `http://localhost:8025`. Mailpit in `compose.yml` has no relay. `/tools` sends notes only to the signed-in user's email address. The template is `internal/notemail/notes.tmpl`; `internal/mailing` implements the protocol with go-mail.

| Variable | Contract |
| --- | --- |
| `MAIL_ENABLED` | Defaults to false in the binary; true in the development example |
| `MAIL_HOST`, `MAIL_PORT`, `MAIL_FROM` | Host/from required when enabled; port defaults to 587 |
| `MAIL_TLS` | `starttls` (required STARTTLS), `tls` (implicit TLS), or `none` |
| `MAIL_USERNAME`, `MAIL_PASSWORD` | Set together; SMTP AUTH PLAIN requires TLS |
| `MAIL_TIMEOUT` | 1s–1m; default 10s |

Plaintext is allowed only outside production and without SMTP credentials. Email is sent outside the database transaction. `MAIL_TIMEOUT` bounds the entire send; job cancellation closes the connection, including during the greeting or TLS handshake. Secrets, SMTP dialogue, and message bodies are excluded from logs.

## Queue interface

After `bin/setup`, run:

```sh
docker compose --profile tools up -d riverui
```

Open `http://localhost:8087`. If the port is occupied, use `RIVERUI_PORT=8088 docker compose --profile tools up -d riverui`. The template uses OSS River UI, bound to loopback. This is a local operations interface without application authentication; do not expose it through a public proxy without separate authentication. Job arguments show only IDs, but the interface can manage queues. Retrying mail/AI through River UI does not bypass the operation's guard against repeating an external effect.
