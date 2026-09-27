# Delivery Ports

This replaces Active Delivery's channel registry and dynamic methods with explicit application orchestration and transport contracts.

## Setup and Basic Usage

Define the message the use case needs and a narrow delivery interface at its consumer. Wire the actual SMTP/HTTP adapter in the composition root. The starter uses `internal/mailing` (go-mail SMTP), `notemail.Service`, and River. Local delivery goes to the explicitly configured Mailpit sink; production requires its own SMTP configuration.

## Deliver Now and Deliver Later

Synchronous delivery is an explicit external call with timeout and error handling. Durable asynchronous delivery requires a persisted job/outbox, a worker and an acknowledgement contract. A goroutine does not provide delivery durability.

## Multiple Channels

A notification operation chooses email, push or webhook using product rules. Each adapter owns protocol details. Return per-channel outcomes if partial success is allowed; otherwise define the failure policy before combining channels.

## Conditional Delivery and Preferences

Load recipient preferences deliberately, distinguish suppression from failure, and identify mandatory notices. Avoid a global callback that sends on every persistence operation.

## Custom Delivery Adapters

Webhook delivery must sign a known payload, bound timeouts, allow only intended destinations and classify retries. A delivery retry uses the same stable event identity. Do not blindly retry non-idempotent provider calls.

## Testing and Development

Unit-test channel selection and payload mapping with fake boundary adapters. PostgreSQL tests prove business state and outbox intent commit together. Worker tests cover duplicate delivery and provider-success/local-ack-failure. A local development sink must be explicitly configured and must never silently replace the production provider.

See [notifications](../topics/notifications.md) for batching, idempotency and state-machine integration.
