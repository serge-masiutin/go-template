# Specification-Test Workflow

The specification test asks whether an object's described behavior belongs to its abstraction layer. It is a design analysis, not a claim that a test suite has been executed.

## Process

1. Read the subject, callers, dependency types and existing tests.
2. List observable responsibilities, including errors and effects.
3. Identify the primary responsibility and layer.
4. Compare test descriptions/assertions with that responsibility.
5. Recommend the smallest extraction only for a demonstrated mismatch.
6. State what was inspected, what ran and what remains unverified.

## Layer Responsibilities

| Layer | Appropriate specification |
| --- | --- |
| Presentation | Decode input, establish identity, delegate, map errors, render a public contract |
| Application | Coordinate a use case, enforce its policy, own transaction/effect order |
| Domain | Preserve invariants and perform pure calculations/transitions |
| Infrastructure | Execute SQL/protocols, map storage/provider data, propagate cancellation/errors |

## Output Format

Subject and current contract; responsibility list; verdict with evidence; proposed test skeleton; existing coverage to preserve; extraction recommendation if needed; verification limits.

## Example

An order handler test covering every discount threshold is testing a domain rule through unnecessary HTTP setup. Move threshold cases to an order/value test. Keep handler tests for malformed input, permissions, operation failure and redirect/JSON shape. Keep PostgreSQL tests for concurrent state transitions and rollback.

A short worker that validates a versioned payload and establishes authority can pass the specification test even if it delegates the actual work in one line.

## Automation Level

Searches and import analysis find candidates; they cannot decide business responsibility alone. Never fabricate executed tests from a proposed skeleton. See [core specification test](../references/core/specification-test.md).
