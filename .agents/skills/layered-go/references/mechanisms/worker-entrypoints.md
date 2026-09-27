# Typed Worker Entry Points

The original guide removes redundant one-line Rails jobs through generated methods. Go should keep explicit queue boundaries and remove only wrappers with no responsibility.

## When to Use

A worker validates a versioned payload, establishes authority, creates a bounded context and delegates to an operation. Even a short worker is useful when it owns those contracts.

## Direct Operation vs Worker Wrapper

Call an operation directly for synchronous work. Use a durable queue when work must survive process failure or happen later. Do not add `_later` methods to domain entities or use reflection to generate job types.

## Detection Signals

Look for an extra forwarding function between a typed worker and its operation, duplicated scheduling logic, and queue retry settings scattered into domain code. Length alone is not evidence of redundancy.

## Configuration

Specify queue, timeout, maximum attempts, retryable errors, backoff, lease expiry and dead-letter handling. Delayed execution stores an explicit scheduled time. Private implementation steps are not independent jobs unless they have a durable contract.

## Application-Wide Patterns

Reuse payload decoding and queue instrumentation only after their semantics repeat. Keep each job's permissions and idempotency visible. Do not introduce a universal untyped task registry.

## When NOT to Use

Do not move a one-step local calculation to a queue merely for abstraction. Do not remove a thin worker that provides validation or execution authority.

## Migration Path

Characterize payload compatibility, retry behavior and scheduling first. Remove one redundant wrapper, preserve public payload versions and run worker failure tests. Existing queued messages remain part of the contract.

## Layer Placement

Workers are inbound presentation adapters; application operations orchestrate; domain types own invariants; queue clients are infrastructure. No queue runtime is installed in the starter.
