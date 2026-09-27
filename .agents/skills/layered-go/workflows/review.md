# Layered Go Review Workflow

Review a diff, file set or branch using the source skill's layer boundaries, specification test, extraction signals and gradual-refactoring discipline.

## Philosophy

Favor coherent extraction over complication. Let patterns emerge before abstracting. Operations coordinate; domain types own invariants. Keep effects and authority explicit. Evaluate actual Go packages rather than requiring a Rails-like directory tree.

## Process

1. Read the diff, surrounding callers, `go.mod`, contracts and relevant tests.
2. Identify responsibilities and runtime/import dependencies.
3. Apply [service classification](analyze-services.md#22-cluster-classification-promote-vs-demote) to operations, provider wrappers and new domain behavior.
4. Apply [the specification test](../references/core/specification-test.md) to the important changes.
5. Inspect hooks, composition, transactions, authorization, context lifecycle and concurrency.
6. Report concrete actionable findings with location, consequence and verification.

## Review Principles

### 1. Layer Boundaries

Flag HTTP/request state in operations, ambient actor/tenant in domain types, SQL in handlers and provider calls hidden in entities. Small feature packages may colocate responsibilities deliberately; judge dependency direction and ownership, not folder names alone.

### 2. Specification Test

Ask whether testing a rule requires HTTP/session/database setup unnecessarily. Keep integration tests for actual integration contracts; reducing setup must not remove coverage of permissions or atomic writes.

### 3. Extraction Signals

Score hidden hooks by responsibility. Assess composition for cohesion and explicit dependencies. Combine churn and mixed reasons to change when identifying god objects. Do not mandate extraction based on line count.

### 4. Execution Context

Context carries cancellation/deadline and boundary metadata. Required business actor/tenant inputs are explicit. Workers validate payloads and establish execution authority; they do not inherit a request session automatically.

### 5. Service Critique

Check operation input/output/error conventions, transaction ownership and effects. Move pure rules near their data; retain an operation when it coordinates a real use case. Do not require base classes, universal Call interfaces or a service for each endpoint.

### 6. Worker Boundaries

A thin worker may own queue validation, retries and authority. Remove only redundant forwarding. A goroutine is not a durable job. Every goroutine needs a stop/join contract; review blocked sends after early returns.

### 7. Abstraction Assessment

Ask what each interface/package makes simpler for a real consumer. An interface defined only for hypothetical mocking is not automatically useful. An `-er` suffix is idiomatic for behavior contracts.

## Review Checklist

- Correct permissions, tenant scope and current role checks for all entry points.
- Explicit transactional writes, affected-row checks and commit failures.
- External effects outside transactions; durable intent/idempotency where required.
- Typed boundary parsing, pure response DTO mapping and no secret exposure.
- Error causes preserved internally and safely reported once at boundaries.
- Cancellation, bounded concurrency and graceful shutdown remain correct.
- Domain and integration tests cover the changed behavior.

## Resolving Violations

Trace callers before proposing a new abstraction. Move an effect to the existing orchestrator when possible. If no orchestrator exists, introduce one only for a real multi-step use case. Keep the behavioral change separate from structural extraction where practical.

## Reviewing Tests

Test public behavior rather than unexported implementation steps. Table tests are useful for pure rules; combine assertions around one expensive integration setup when they form a coherent scenario. Do not replace meaningful assertions with a test that merely verifies a mock was called.

When removing duplication, state which boundary test still proves delegation and which lower-level test proves the rule. Use real PostgreSQL connections for transaction races. Use `testing/synctest` for suitable in-process timing tests, not to simulate external database scheduling.

## Output Format and Severity

Each finding includes severity, exact path/line, trigger, consequence, remedy and relevant check. Critical findings are correctness, authorization or data-integrity defects. Warnings need demonstrated coupling/failure risk. Suggestions are optional improvements. If no defect is found, say so and state the review limits.

## Related

[Architecture analysis](analyze.md), [service audit](analyze-services.md), [Go sources](../references/go-sources.md).
