# Architecture Analysis Workflow

Comprehensive layered architecture analysis of a Go codebase or a focused directory. Preserve the original survey: layer violations, abstraction opportunities, god objects, hidden hooks and missing responsibilities.

## Analysis Process

### 1. Structural Assessment

Read `go.mod`, commands, package imports, callers, tests and deployment configuration. Map responsibilities before directories: HTTP/CLI/workers are entry points; operations coordinate; domain types own invariants; stores/provider clients implement I/O; React renders page contracts.

Use `go list -json ./...` and `rg --files -g '*.go' cmd internal` as an inventory. Do not label every feature package a pure domain layer when it intentionally includes persistence. Note existing CI import rules; propose one only for a real declared boundary.

### Service-Layer Brief

Keep this brief under roughly 25 lines. Include the empirical operation definition, scale/triage gate, hidden orchestration, contextual architecture verdict, convention strength, up to three cluster headlines, observed unused abstractions, cross-cutting issues and anemic-domain risk. Omit unsupported categories.

Use [the service audit](analyze-services.md) for detailed classification. Counts and suffixes prompt inspection; they do not justify new packages or base classes. A small coherent codebase may need no service layer changes.

### 2. Layer Violation Detection

Inspect domain behavior for ambient user/tenant state, HTTP dependencies, mail/queue/provider calls and hidden configuration. Inspect operations for transport DTO ownership and response rendering. Inspect handlers for SQL and business calculations.

A shared operation may enforce authorization; authentication remains at the entry point. A context parameter for cancellation is valid; retrieving required business inputs from context values hides a dependency.

Follow both direct imports and calls through interfaces. Distinguish actual boundary violations from conceptual layers colocated in one feature package. Report only observed violations.

### 3. God-Object Identification

Combine churn (`git log --name-only`), cohesive responsibilities, dependency fan-out and test setup. Size alone is insufficient. Prioritize objects whose unrelated reasons to change cause defects or costly edits; see [god-object analysis](analyze-gods.md).

### 4. Callback Analysis

Inventory hooks and indirect effects, score their responsibility and trace chains. In new Go code prefer explicit calls. See [callback analysis](analyze-callbacks.md); the score does not authorize adding a callback framework.

### 5. Coupling via the Specification Test

List what each representative test actually proves. HTTP tests should cover transport/permissions/contracts; domain tests should cover rules; PostgreSQL tests should cover persistence and concurrency. Extract responsibility only when the evidence supports it.

### 6. Composition Health

Inspect embedding and shared interfaces for cohesion, hidden dependencies and accidental promoted APIs. File slicing is not decomposition. A concrete type can be preferable to an unused interface.

### 7. Anti-Pattern Detection

Look for redundant wrappers, untyped generic services, global identity, skipped validation flags, HTML string builders, swallowed errors, request goroutines pretending to be durable jobs and mixed pool/transaction calls.

A thin worker can own a meaningful queue boundary. An `-er` interface name is normal Go. Do not port these Rails heuristics mechanically.

### 8. Pattern Gap Analysis

Recommend forms/queries/policies/presenters/operations only where repeated responsibilities have a stable contract. Reuse the local architecture and propose the smallest reversible step.

## Reporting Principles

Report facts with paths/lines, then their effect, then a concrete remedy and test. Separate confirmed defects from risks and optional design preferences. Preserve healthy patterns explicitly; omit empty finding categories and raw scan noise.

## Output Format

```text
Architecture analysis: <scope>
Summary: <current design and main consequence>
Service-layer brief: <compact evidence>
Findings: <severity, location, behavior, proposed change, validation>
Preserve: <healthy boundaries>
Immediate steps: <small sequence>
Later options: <conditional on observed growth>
Checks and limitations: <what was actually inspected/run>
```

## Severity Levels

Critical: data/authorization/transaction integrity or a concrete correctness defect. Major: repeated coupling or hidden effects with demonstrated change/test cost. Minor: localized clarity or naming improvement. Architectural taste alone is not a critical finding.
