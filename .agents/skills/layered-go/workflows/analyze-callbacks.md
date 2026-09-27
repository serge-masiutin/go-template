# Callback Analysis Workflow

Use the original scoring method to audit existing hooks and hidden side effects. Go has no required lifecycle callback system; new code uses explicit domain methods and operations.

## Callback Scoring System

| Score | Responsibility | Default action |
| --- | --- | --- |
| 5 | Pure transformation | Keep responsibility; call explicitly |
| 4 | Internal consistency | Keep invariant; verify atomicity |
| 3 | Timestamp | Make clock and event explicit |
| 2 | Background trigger | Move delivery intent to operation |
| 1 | External business operation | Extract orchestration immediately when unsafe |

## Analysis Process

### 1. Find Hooks

Read store methods, constructors, `init`, middleware, event subscriptions and queue triggers. Search for known hook APIs only after inspecting dependencies. `go func`, provider calls and save methods are leads, not proof of a callback defect.

### 2. Score Each Hook

Record source, trigger, inputs, writes, external effects, error propagation and transaction context. A score alone is insufficient: explain why the responsibility belongs where it does.

### 3. Identify Chains

Trace indirect writes and subscribers until the workflow ends. Record cycles, repeated effects and order dependence. Distinguish required atomic writes from post-commit delivery.

### 4. Check Skip Patterns

Find caller flags, global toggles and test setup that disables hooks. Determine whether they represent distinct legitimate commands or bypass invariants. Do not preserve a skip-validation escape hatch as an implementation convenience.

## Output Format

For each owner: trigger, score, responsibility, current failure mode and proposed explicit caller. Then list chains, skip flags and prioritized extractions. Use actual findings; do not fill a report with invented example metrics.

## Extraction Patterns

Characterize behavior first. Move orchestration to the existing use case, keep pure transformations near the data, commit related writes together and record durable delivery intent when required. Remove the hook only after every caller uses the new path.

Test rollback at each write, no effect before commit, duplicate processing and crash after provider success. See [registration extraction](../examples/callbacks-to-service.md) and [callback topic](../references/topics/callbacks.md).
