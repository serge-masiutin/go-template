# God-Object Analysis Workflow

Identify objects with too many independent responsibilities and plan gradual decomposition. Preserve the original combination of churn, complexity and responsibility clustering; Go file length is not a verdict.

## God-Object Indicators

### Quantitative Signals

Record file/type size, exported method set, dependencies and test setup as inspection signals. A 500-line cohesive parser can be healthier than five tightly coupled 100-line packages.

### Churn Analysis

Use actual Git history to find frequently changed files, then inspect the reasons for change. High churn from one active feature differs from unrelated edits colliding in one type.

### Complexity Indicators

Look for independent lifecycles, mixed HTTP/SQL/domain/provider responsibilities, many unrelated flags, global context, broad interfaces and tests requiring irrelevant collaborators.

## Analysis Process

### 1. Identify Candidates

Inventory Go files and types, then cross-reference history and dependencies. Do not rely on LOC sorting alone. Include workers and operation hubs, not only entities.

### 2. Analyze Each Candidate

Group methods/fields by responsibility and caller. Draw state ownership, transaction boundaries and effects. Identify invariant clusters that must stay together. Check whether a proposed destination would itself become a god object.

### 3. Recommend Decomposition

Choose the smallest coherent extraction and preserve public behavior. Name the new contract, dependencies and tests; avoid moving methods into an untyped service bag.

## Output Format

For each candidate report actual metrics, responsibility clusters, concrete current pain, proposed extraction order and behavior-preserving checks. Include non-issues where a large cohesive object was inspected and correctly retained.

## Extraction Priority Matrix

Prioritize high churn plus mixed responsibilities and observed defects. Lower priority: large stable cohesive code with no demonstrated cost. Do not present heuristic scores as objective architecture quality.

## Decomposition Strategies

1. **Composition:** extract cohesive reusable behavior into named fields/functions; avoid embedding as inheritance.
2. **Associated collaborator:** a secondary domain concept owns its own behavior and data; persistence relationships remain explicit.
3. **Value object:** group identity-less values and invariants; protect mutable slice/map ownership.
4. **Operation:** move cross-record workflows and external effects into an application use case.
5. **State machine:** replace contradictory flags with typed states/events when transitions justify it; enforce concurrency in storage.

Run domain tests and relevant database/HTTP tests after each step. See [decomposition example](../examples/god-object-decomposition.md).
