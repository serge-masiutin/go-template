# Architecture decision record

Use this as a fill-in structure, not a demand for a long document. Remove irrelevant sections and replace prompts with evidence. Label estimates and unknowns.

## Decision

What will change, for which operation, and what observable behavior will result?

## Context and requirements

- Workload, data size/growth, skew, latency/throughput targets, and measured bottleneck.
- Invariants and required read/acknowledgment semantics.
- Failure domains, recovery targets, operating/cost constraints.
- Evidence: relevant code, schema, configuration, traces, tests, and documentation versions.
- Assumptions that could change the decision.

## Ownership and dataflow

Identify authoritative facts and writers, derived views, transaction boundaries, asynchronous edges, external effects, and progress/operation identities. Use a small diagram or table if it clarifies the flow.

## Options and trade-offs

Compare the current/simple baseline and plausible alternatives on correctness, normal/tail performance, failure recovery, cost, operational burden, and reversibility. Explain why the selected option satisfies the requirements and when another would be preferable.

## Failure and recovery behavior

For relevant histories, state the enforcing mechanism and caller-visible result: concurrent writes, unknown outcomes, duplicate/reordered work, stale readers, leader/owner change, overload, and restore. Name any intentionally weakened guarantee.

## Evolution and rollout

Compatibility matrix; schema/data migration; snapshot/live-change boundary; bounded backfill; comparison/acceptance criteria; reader cutover; rollback conditions; old-state retirement; deletion/retention implications.

## Validation and operations

Tests or measurements performed and results; unverified claims; required additional checks; observability signals; recovery owner/procedure; conditions that trigger revisiting the decision.
