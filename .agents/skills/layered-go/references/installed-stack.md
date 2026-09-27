# Installed Stack Contracts

Read the actual `go.mod`, `internal/config`, and [stack decisions](../../../../docs/stack.md) before changing adapters. Preserve the layered extraction criteria; libraries implement mechanisms, not business ownership.

## Persistence and Migrations

- Use GORM generics in feature stores, explicit projections, bound values and owner predicates. Keep private persistence records separate from exported DTOs. SQL identifiers and sort expressions come from trusted code.
- Keep invariants in domain/application behavior and constraints in PostgreSQL. Do not move business effects into GORM hooks, association autosave or serializers. Do not use AutoMigrate.
- Use Goose SQL migrations and `bin/manage migrate` as a deployment step. Its lock covers application and River migrations. Server/worker startup does not migrate. Check each migration's transaction/lock requirements.
- `database.DB` shares pgx connections with GORM/database/sql and SCS. Do not open parallel independent ORM pools. The worker owns one extra LISTEN connection; transaction-mode PgBouncer is unsuitable for it.

## Jobs and External Effects

- Follow `notemail.Service.Request` / `assistant.Service.Request`: persist the operation and River `InsertTx` with the same `*sql.Tx` inside one GORM transaction. Propagate commit errors.
- Queue payloads carry versioned operation IDs, not credentials, questions, notes or recipient text. Workers reload state and derive authority from the persisted owner. Recheck data ownership when reading.
- Keep entry points in `internal/background` thin. Domain/application operations own transitions; adapters own SMTP/model protocols. No untracked goroutine substitutes for a durable job.
- Mail/AI jobs have one attempt. Conditional transitions prevent repeat execution after an uncertain provider outcome. Completed work is idempotent; reconciliation marks terminal/abandoned work failed without repeating external effects.
- New retry policies need a provider idempotency contract and tests for success followed by lost acknowledgement. Database-only reconciliation can be retried. Drain workers on shutdown, then cancel within a bound.

## Mail and AI

- Use go-mail through `internal/mailing`; require TLS in production, keep template rendering in the notification operation, and send after commit. Mailpit is the explicitly configured local sink.
- Use stable Genkit flow/generate/tool APIs through `internal/assistant`. Per-invocation tools close over the trusted actor; the model cannot supply another user's ID. Tools do not inherit authority from prompt text.
- Version prompt/model/schema, validate structured output, and set whole-flow deadlines, turns and output tokens. Preserve disabled-by-default AI and explicit provider selection. Never silently fall back to another model.
- Do not log raw SDK errors or enable Genkit content telemetry/reflection incidentally. River receives sanitized errors; diagnostics retain safe categories. Keep content out of job args, metrics and routine logs.
- Synthetic provider tests establish contracts, not model quality. Run the [quality evaluations](../../../../evals/notes-assistant.md) against the selected model before enabling a product AI feature; do not claim unrun evaluations passed.

## Verification

Use real PostgreSQL tests for rollback, ownership, concurrent enqueue and state reconciliation; fake only external SMTP/model boundaries. Verify provider protocol errors and limits through local HTTP fixtures. Keep raw text and high-cardinality actor IDs out of Prometheus labels. Read [background operations](../../../../docs/background.md) and [AI contracts](../../../../docs/ai.md) when extending these paths.
