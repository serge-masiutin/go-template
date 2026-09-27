# Service-Layer Audit Workflow

Deep audit of application operations and service-shaped code, adapted from the original three-phase Layered Rails workflow. Preserve its empirical mirror, promote/demote classification, per-cluster proposal and evidence-driven reporting. Go package boundaries and consumer contracts replace Rails directories and inherited base classes.

## Purpose

Answer one question: **is the service layer healthy, and if not, what concrete proposals would shape it?**

This focused audit complements [architecture analysis](analyze.md). The parent workflow produces a short service-layer brief; this workflow supplies contracts, current pain, test benefits and placement options for each supported recommendation.

## Core Idea

Services are a waiting room for responsibilities that have not yet formed a stable abstraction. Two failure modes are equally bad: a bag of unrelated operations and passive entities whose business rules have all moved elsewhere.

Specialization means clearer responsibilities, not a mandatory directory or interface for every pattern. A small Go feature can keep entity behavior, a concrete store and an operation in one domain-named package. Split packages when doing so enforces a useful boundary.

## How This Report Justifies Recommendations

Every recommendation must point to observed consequences: duplicate setup, unclear failure semantics, brittle tests, missing transaction ownership, repeated authorization or costly changes across unrelated packages. A design slogan alone is insufficient.

Use both the **specification test** and **current pain**. Explain how the proposed contract changes testing and use, not merely how files move. Cite actual paths and lines; never invent counts, latency or test coverage.

## Phase 1 — Discovery

### 1.1 Waiting-Room Gate

Map actual operations first; Go has no required `services` directory. Use `rg --files -g '*.go' cmd internal` and `go list -json ./...`, then classify representative files by responsibility. Count operations and their share of application code only after classification.

The upstream thresholds (fewer than 10 services or less than 10% of code) are triage heuristics. When the operation layer is small, skip elaborate clustering; still inspect domain packages for hidden orchestration. Never create a services package merely to satisfy a metric.

### 1.2 Empirical Mirror

Sample 5–10 representative operations, including recent, large and small ones. Inspect inputs, dependency ownership, public methods, results, errors, transaction boundaries and effects.

Write one plain sentence describing what an operation means **in this project**. Example: “An operation receives an explicit actor and typed command, coordinates stores in one transaction, and returns a value plus a classified error.” If shapes vary, state that concretely; do not substitute an invented convention.

Keep supporting counts in the convention report. This mirror anchors every later proposed contract.

### 1.3 Hidden Services in Domain Packages

A pure calculator spanning several entities can remain domain behavior. SQL queries are persistence implementations. HTTP clients are infrastructure. A workflow coordinating those dependencies is application behavior. Classify by purpose and dependency ownership, not suffix or folder alone.

Inspect exported functions/methods, imports, SQL, provider calls and queue writes. A domain package name does not make an external call pure. Conversely, a file named `service.go` may hold a small pure calculation that needs no new package.

### Decision Gate and Models-First Stance

If there is no substantial operation cluster and no misplaced responsibility, report the healthy boundary briefly and stop. Do not create work for an empty waiting room.

A deliberate models-first approach is compatible with operations for cross-record workflows. Identify its tradeoffs without calling it defective by default: hidden side effects, difficult tests, mixed persistence/domain responsibilities and unclear authority. A pure entity cannot own mail, request state or a provider client merely because the project prefers rich models.

## Phase 2 — Clusters: Promote Up or Demote Down

### 2.1 Cluster Detection

Group candidates by shared purpose, dependency shape, input/output contracts and failure semantics. The original thresholds (five similar names or 10% of services) prompt inspection; they do not mandate extraction.

Read 3–5 bodies before proposing a shared abstraction. Prefix/suffix similarity is a hypothesis. Split heterogeneous groups rather than forcing every member behind one interface.

### 2.1.1 Cluster Smells

- **Verb-prefix families:** determine whether each file names a real use case or merely repeats one action over several domain concepts. Put stable behavior near the concept; do not rename verbs mechanically.
- **Action-axis grouping:** repeated Add/Remove packages may hide one coherent feature. Group related operations when they share invariants and dependencies; do not build an untyped action dispatcher.
- **Mixed delivery and record management:** separate notification decisions from cleanup/update responsibilities before choosing placement.
- **Pure presentation mappers:** page-shaped DTO construction belongs to presentation, not a generic service registry.
- **Repeated cross-cutting effect:** if every wrapper adds the same audit/delivery step, identify that missing responsibility instead of consolidating unrelated use cases.

### 2.2 Cluster Classification: Promote vs Demote

| Candidate purpose | Home | Contract |
| --- | --- | --- |
| Coordinate a use case | Application operation | Typed command, explicit authority, result/error |
| Derive a pure value/rule | Entity method or domain function | Required values, no I/O |
| Bundle identity-less meaning | Value object | Validated value semantics |
| Own one entity's cohesive behavior | Composed collaborator | Narrow named behavior |
| Execute SQL | Feature store or focused query | Authorized scope, typed result, context |
| Speak a provider protocol | Infrastructure adapter | Bounded I/O and mapped errors |
| Shape a page/API | Presentation mapper | Explicit public DTO |

#### Purpose First, Search Pattern Second

Read the caller and downstream consumers. A function that inserts notification records exists to initiate delivery even if it makes no immediate network call. A pure arithmetic function does not become an operation because its filename ends in Service.

Direct mail/queue/provider calls are evidence of I/O. Distinguish protocol implementation (infrastructure) from workflow decisions (application). LLM calls are external I/O, never pure domain calculations.

#### Downstream Consumers

Follow persisted events and callbacks to their consumers. If creating a record triggers a workflow, make that intent visible. Extract a heavy recipient query into persistence without misclassifying the notification operation as domain logic.

#### Ambient Infrastructure

Tracing metadata is not business behavior, but a tenant lookup, identity service or storage read is still an I/O dependency. Do not excuse hidden calls as “ambient infrastructure.” Load facts at the appropriate boundary and pass them explicitly.

#### Domain Concepts Misfiled as Services

Separate eligibility/calculation from querying and persistence loops. The domain concept may own the rule; the query supplies facts; the operation coordinates writes. Unlike Active Record-oriented advice, SQL does not become pure domain code when moved into a domain-named file.

#### Eligibility and Feasibility

Distinguish permission (actor may act), validation (input/invariant is valid) and system feasibility (resources/capacity allow it). A small emerging family can remain explicit functions. Extract a shared rule abstraction only when semantics actually repeat.

#### External-Vocabulary Mappers

Mapping provider statuses into internal values belongs near the adapter that owns that vocabulary. Prefer a function or fixed mapping over a runnable service object. Provider changes should not spread through domain packages.

#### Wrong-Abstraction Brake

A small heterogeneous cluster may be healthier with some duplication. Name the least similar members and explain why a common interface would lose meaning. Extract only a coherent subset when dependencies differ.

#### Domain Method Alternative

A calculation, invariant or transition over one entity may become its method. Inputs are values, not HTTP requests or provider clients. Keep persistence and delivery in the caller's operation. A stateless operation may also be a function; Go does not require a struct merely to signal orchestration.

#### Destination and Source Complexity Brakes

Before moving behavior, inspect the destination's cohesion and churn. Do not make a god object larger. Pure secondary behavior may become a composed collaborator; cross-record orchestration remains application code. A complex source that performs I/O is not a domain collaborator regardless of its size.

When several phases always execute together, investigate a named pipeline with explicit transaction/effect boundaries. Do not automatically create one abstraction per phase.

### 2.3 The Cluster Proposal Block

For every supported cluster recommendation, report these eight elements in order:

1. **Identity:** actual files, grouping evidence, destination responsibility or “wait.”
2. **Current pain:** concrete duplication, setup cost, failure ambiguity or coupling.
3. **Specification-test verdict:** quote actual test names/assertions showing mixed responsibilities.
4. **Suggested contract:** inputs, result, error classes, authority, transaction ownership and effects.
5. **Optional library:** only when existing tools lack a measured required capability; inspect its current API.
6. **Shared machinery:** exactly what repeats and what stays specific. Prefer composition to a base type.
7. **Benefits:** show how the next real consumer/test becomes simpler; do not invent savings.
8. **Placement options:** favor existing domain-named packages; separate adapters only when useful.

For borderline candidates, report evidence and why waiting is preferable; do not design a framework speculatively.

Example contract for an importer: a parser validates a versioned row; an operation commits a bounded batch; a result counts accepted/rejected rows with typed reasons. The adapter owns CSV/provider syntax. Partial success and retry behavior are part of the contract, not hidden behind `Import(any)`.

### 2.4 Special Cases

#### Operation Effects vs Event Broadcast

A required effect belongs in the operation's contract. A notification that must survive a crash requires durable intent. An in-process event callback is not a substitute for a transaction or a queue.

#### API Wrappers

Do not unify unrelated providers before there is polymorphic usage with equivalent semantics. A narrow adapter can still be valuable for timeout, authorization, schema validation or error translation. Remove a wrapper only after checking those responsibilities.

#### Small Heterogeneous Clusters

Do not create a minimal base class: Go has no class inheritance. Use explicit functions/types and a small consumer interface only when consumers actually need substitution.

#### God-Object Triage

Distinguish size, cohesion, coupling and placement. Splitting a file without changing responsibility is cosmetic. Preserve behavior and prioritize the change that reduces a real failure/change cost.

## Phase 3 — Cross-Cutting Findings

### 3.1 Convention Strength

Sample representative code or all operations when small. Report deviations from the empirical mirror in these axes: naming, typed inputs, dependency construction, context lifecycle, results/errors, transaction ownership and testing style.

Go has no required base-service class or universal Call method. A uniform `Execute(any) any` is weaker than several meaningful typed methods. Compare conventions within a coherent family; do not demand one verb across queries, deliveries and commands.

#### Naming Consistency

Read names in package context (`accounts.Register`, not `accounts.AccountRegistrationService`). A minority style can be a review finding when it creates ambiguity, but do not rename the repository merely to equalize suffix percentages.

#### Test Impact

Explain how current inconsistencies affect fakes, fixtures, error assertions and integration coverage. Healthy uniformity exposes contracts; false uniformity erases them.

### 3.2 Organization Shape

Inspect package size, imports, fan-in/fan-out and concept boundaries. Generic mixed packages such as `utils` deserve scrutiny. A flat but cohesive package may be healthy; many tiny packages may increase coupling and import cycles.

#### Vendor Naming

Make an adapter's role clear. A provider-named package under a feature/infrastructure boundary can be appropriate; a vendor brand must not hide business orchestration. Do not impose Rails autoload directory rules on Go.

#### Ambiguous Classification

When several names fit the same behavior, identify the missing contract before proposing a rename. A rename needs concrete benefit for callers, dependency ownership or tests. A new base type is not the default fix.

### 3.3 Naming Smells and Alternative Forms

#### Behavior Suffixes

`Reader`, `Writer`, `Sender` and similar names are idiomatic for behavior interfaces. Flag a vague `Manager` only when its responsibilities are actually unclear. An `-er` suffix alone is not evidence of a bad abstraction.

#### Tautological Names

Evaluate the full call expression. A type/method pair that repeats a verb without owning state or dependencies may become a function. A type that owns a provider client can still be useful even with one public method.

#### Ubiquitous Language

Compare names with schema, routes and product terminology. A name absent from the database is not automatically invented: operations and value concepts may have no table.

#### Scheduling-Frequency Names

Scheduling belongs to the worker/scheduler boundary; the underlying operation should name its business purpose. Remove redundant scheduling wrappers while preserving payload validation, retries and execution authority.

### 3.4 Implicit Workflows

Trace operation calls, especially chains, hubs and cycles. Long chains are inspection signals, not automatic defects. Identify who owns the complete use case, transaction, cancellation and external effects. A named orchestrator should make the workflow observable and testable without hiding its steps.

### 3.5 Layer Hygiene

Inspect HTTP/session/render dependencies in operations and external I/O in pure domain types. Request-scoped context metadata may live at transport boundaries; required actor/tenant values remain explicit operation inputs.

#### Authorization in Shared Operations

Unlike the original controller-only recommendation, shared Go operations enforce their policy for all callers. Transport authenticates and maps errors. A worker establishes explicit authority. Flag duplicated or missing enforcement, not the mere presence of a policy call in an operation.

#### Sinkholes

A pure forwarding wrapper may be removable. First check whether it owns a boundary: decoding, authorization, transaction, queue protocol or error translation. A short handler or worker is not automatically an architecture sinkhole.

### 3.6 Misplaced Code

#### Self-Separation Gate

A cohesive adapter or mapper colocated with its consumers need not move for taxonomy alone. Recommend relocation only for a concrete import, testing, ownership or clarity benefit.

#### Mixed-Layer Files

Separate protocol/SQL mechanics from domain rules when they obscure independent responsibilities. Moving the entire mixed file to a different folder does not solve the problem.

#### CLI Utilities

Interactive parsing and stdout/stderr belong to a command under `cmd`, while reusable operations live under `internal`. Inject only the I/O needed by the command; domain methods do not print progress.

#### Application-Shaped Domain Code

A pure domain package must not import HTTP/provider/queue implementations. A small feature package may intentionally contain both a store and behavior; report that actual boundary rather than claiming its directory is a strict domain layer.

#### Thin SDK Wrappers

Check the installed SDK's documented timeout, retry and serialization behavior. Keep adapters with a real contract; inline meaningless forwarding at the infrastructure boundary. Never move provider calls into a domain model to reduce file count.

### 3.7 Peer Files

Check one-to-one worker/operation pairs for distinct responsibilities. A worker handles durable payload and queue semantics; an operation can be shared by HTTP/CLI. Remove only redundant intermediate forwarding. Domain constructors/methods orchestrating multiple stores may need an application operation.

### 3.8 Tests and Specification

Discover the actual test conventions. Compare test names and setup with the subject's responsibility. Mock external boundaries, not private methods. Keep database tests for SQL and transactions; fake stores cannot establish constraint or concurrency correctness.

For each refactor explain the specific testing benefit and which integration contract remains covered. If removing duplicate tests, show the replacement boundary coverage. Do not claim fewer tests automatically means better design.

### 3.9 Anemic-Model Risk

Inspect whether invariants live beside the data they govern. Structs used solely as transport/query DTOs need not have behavior. Do not add methods to DTOs to satisfy a richness metric.

### 3.10 Target Architecture Signals

A healthy design has visible ownership, explicit effects, small useful contracts and tests at the right boundaries. Empty folders and unused interfaces are not specialization. A package with two cohesive files is not a “ghost” merely because it lacks a base class.

Report a contextual verdict: healthy decomposition, mixed responsibilities, or an overcrowded waiting room. Support it with observations rather than a numerical architecture score.

## Reporting Principles

Lead with the concrete risks and top three actions. Keep internal scan logs and raw metrics out of the narrative unless they support a finding. Separate observed facts, design inference and optional preference. Omit empty sections.

## Output Format

```text
Service-layer audit: <project>
In this project, an operation is: <empirical mirror>
Top actions: <up to three evidence-backed changes>
Recommendations: <per-cluster proposal blocks>
Codebase insights: <conventions, organization, hygiene, workflows, tests, domain health>
Preserve: <healthy patterns>
Verification: <checks performed and limits>
Next steps: <small behavior-preserving sequence>
```

## Related

- [Architecture analysis](analyze.md)
- [Review](review.md)
- [Specification test](spec-test.md)
- [Go sources and adaptation decisions](../references/go-sources.md)
