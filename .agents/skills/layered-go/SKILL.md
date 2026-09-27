---
name: layered-go
description: Write, refactor, and review Go code using the layered design framework adapted from Layered Rails. Use for HTTP handlers, domain behavior, application operations, queries, policies, workers, serialization, and architecture reviews; inspect explicit contracts, dependency direction, transactions, and specification-test failures.
metadata:
  version: "2"
  source: "serge-masiutin/rails-template@a315bfbcf29ec14b83c9f05378abd595bc5228e0"
---

## Project Context

- Read [architecture](../../../docs/architecture.md), `go.mod`, and callers first. Keep small features in domain-named packages under `internal`; do not create four empty directory layers or one service per endpoint.
- The starter uses net/http, Gonertia, React, SCS, GORM/pgx, Goose, River, go-mail, and Genkit. Read the [installed stack contracts](references/installed-stack.md). Domain values own invariants; stores own persistence; handlers parse, authorize, delegate, and render. Pass actor identity explicitly; context carries cancellation/deadlines, not hidden business dependencies.
- Query/store implementations are infrastructure. Put an interface at its consumer only when a real boundary requires it; the domain does not import pgx or HTTP. Runtime calls and import direction are different graphs.
- Examples describe possible features, not installed APIs. The original Rails gem guides are replaced by Go mechanism references. River jobs, SMTP delivery, and a bounded read-only Genkit assistant are implemented. Examples for other providers, payments, and generic outboxes remain illustrative; verify their concrete implementation before use.
- Use [testing](../../../docs/testing.md). Do not generate Ruby, Active Record hooks, base-service inheritance, or a universal repository layer.

# Layered Go

Design and review Go applications using layered architecture principles.

## Go Evidence

Read [Go sources and adaptation decisions](references/go-sources.md) when changing architecture guidance. Reviewed 2026-09-27 against Go 1.27.1, including the Go team’s September 2026 goroutine diagnostics guidance.

## Quick Start

Reason about Go applications using four architecture layers with **unidirectional data flow**:

```
┌─────────────────────────────────────────┐
│           PRESENTATION LAYER            │
│   Handlers, React views, Streams, Workers*   │
│ Forms, Filters, Presenters, Serializers │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│           APPLICATION LAYER             │
│ Operations, Policies, Notification flows │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│             DOMAIN LAYER                │
│  Entities, Value Objects, Domain behavior   │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│          INFRASTRUCTURE LAYER           │
│  pgx, HTTP clients, Storage adapters    │
└─────────────────────────────────────────┘
```

\* Jobs are **internal inbound** entry points: design-wise they follow the same rules as controllers, with explicit job validation and execution authority.

**Core Rule:** Domain behavior must not depend on transport or concrete infrastructure.
Runtime calls move from entry points through operations to domain behavior and persistence.
Go adapters implement consumer-owned contracts; do not mistake runtime arrows for imports.

See [Architecture Layers Reference](references/core/architecture-layers.md) for the full layer responsibilities and the Four Rules deep-dive.

## What this skill is for

Use this skill when:

1. **Analyzing a codebase** — apply the [architecture analysis](workflows/analyze.md) for a full audit, or zoom in with the [service-layer audit](workflows/analyze-services.md), [callback analysis](workflows/analyze-callbacks.md), or [god-object analysis](workflows/analyze-gods.md).
2. **Reviewing code changes** — run the [code review workflow](workflows/review.md) on a diff, file, or branch.
3. **Running the specification test** — use the [spec-test workflow](workflows/spec-test.md) on a single file or directory to evaluate whether code belongs in its current layer.
4. **Planning gradual adoption** — generate a phased roadmap with the [layerification plan workflow](workflows/plan.md), focused on a goal like "introduce authorization" or "decompose god objects."
5. **Planning a feature** — I'll apply the layered principles below to whichever code you're about to write.
6. **Implementing a specific pattern** — authorization, notifications, view components, AI integration, etc. — see the Pattern Catalog and Topic References below.

Invoke these workflows by their linked files or ask in plain language. This repository
does not publish a Claude slash-command plugin.

## Workflows

Reusable procedures bundled inside this skill. Read the file and apply it to the target code:

- [Architecture analysis](workflows/analyze.md) — full layered-architecture audit of a Go codebase
- [Code review](workflows/review.md) — review a diff or file set for layer violations
- [Specification test](workflows/spec-test.md) — evaluate whether code belongs in its current layer
- [Service-layer audit](workflows/analyze-services.md) — deep audit of application operations and service-like types (per-cluster proposals, contracts, layer hygiene)
- [Callback analysis](workflows/analyze-callbacks.md) — score existing hooks and implicit side effects and find extraction candidates
- [God-object analysis](workflows/analyze-gods.md) — identify oversized models and recommend decomposition
- [Gradual layerification plan](workflows/plan.md) — incremental roadmap for adopting layered patterns
- [Architecture-check setup](workflows/archspec.md) — generate and verify a tailored Go import-boundary check enforcing the layer boundaries in CI

## Core Principles

### The Four Rules

1. **Unidirectional Use-Case Flow** - Entry points delegate; domain behavior does not call transport
2. **No Reverse Dependencies** - Domain packages do not import HTTP, UI, or concrete database adapters
3. **Abstraction Boundaries** - Each abstraction belongs to exactly one layer
4. **Minimize Connections** - Fewer inter-layer connections = looser coupling

### Common Violations

| Violation | Example | Fix |
|-----------|---------|-----|
| Model uses Current | global actor in domain method | Pass user as explicit parameter |
| Service accepts request | `*http.Request` in operation | Extract value object from request |
| Controller has business logic | Pricing calculations in action | Extract to service or model |
| Anemic models | All logic in services | Keep domain logic in models |

| Category | Reference |
|----------|-----------|
| Layer violations (ambient actor in domain, request in services, notifications in models, business logic in controllers) | [layer-violations.md](references/anti-patterns/layer-violations.md) |
| Service objects (anemic models, bag of random objects, premature abstraction) | [service-objects.md](references/anti-patterns/service-objects.md) |
| Callbacks (operation callbacks, skip callbacks, control flags) | [callbacks.md](references/anti-patterns/callbacks.md) |
| Concerns (code-slicing, overgrown) | [concerns.md](references/anti-patterns/concerns.md) |
| Helpers (HTML construction in helpers) | [helpers.md](references/anti-patterns/helpers.md) |
| Jobs (anemic jobs) | [jobs.md](references/anti-patterns/jobs.md) |
| Testing (testing wrong layer) | [testing.md](references/anti-patterns/testing.md) |

### The Specification Test

> If the specification of an object describes features beyond the primary responsibility of its abstraction layer, such features should be extracted into lower layers.

**How to apply:**
1. List responsibilities the code handles
2. Evaluate each against the layer's primary concern
3. Extract misplaced responsibilities to appropriate layers

See [Specification Test Reference](references/core/specification-test.md) for detailed guide.

## Pattern Catalog

| Pattern | Layer | Use When | Reference |
|---------|-------|----------|-----------|
| Service Object | Application | Orchestrating domain operations | [service-objects.md](references/patterns/service-objects.md) |
| Query Object | Infrastructure | Complex, reusable queries | [query-objects.md](references/patterns/query-objects.md) |
| Form Object | Presentation | Multi-model forms, complex validation | [form-objects.md](references/patterns/form-objects.md) |
| Filter Object | Presentation | Request parameter transformation | [filter-objects.md](references/patterns/filter-objects.md) |
| Presenter | Presentation | View-specific logic, multiple models | [presenters.md](references/patterns/presenters.md) |
| Serializer | Presentation | API response formatting | [serializers.md](references/patterns/serializers.md) |
| Policy Object | Application | Authorization decisions | [policy-objects.md](references/patterns/policy-objects.md) |
| Value Object | Domain | Immutable, identity-less concepts | [value-objects.md](references/patterns/value-objects.md) |
| Collaborator Object | Domain | A slice of one model's behavior in a typed delegate | [collaborator-objects.md](references/patterns/collaborator-objects.md) |
| State Machine | Domain | States, events, transitions | [state-machines.md](references/patterns/state-machines.md) |
| Composition | Domain | Shared cohesive behavior | [concerns.md](references/patterns/concerns.md) |
| Repository | Infrastructure | Use a concrete store for persistence; add a consumer-owned interface only when a real boundary needs substitution | [repositories.md](references/patterns/repositories.md) |

### Pattern Selection Guide

**"Where should this code go?"**

| If you have... | Consider... |
|----------------|-------------|
| Complex multi-model form | Form Object |
| Request parameter filtering/transformation | Filter Object |
| View-specific formatting | Presenter |
| Complex database query used in multiple places | Query Object |
| Business operation spanning multiple models | Service Object (as waiting room) |
| Authorization rules | Policy Object |
| Application operations and delivery ports | Delivery Object (explicit notification orchestration) |

**Remember:** Services are a "waiting room" for code until proper abstractions emerge. Don't let a generic services package become a bag of random objects.

## Refactoring Scenarios

Canonical before/after transformations for the most common layerification moves. The [layerification plan workflow](workflows/plan.md) uses these as reference templates when proposing phases.

| Scenario | Goal area | Reference |
|----------|-----------|-----------|
| Extract callbacks to service | callbacks, hidden persistence-hook chains | [callbacks-to-service.md](examples/callbacks-to-service.md) |
| Extract authorization to policy | authorization, permissions | [authorization-to-policy.md](examples/authorization-to-policy.md) |
| Extract query logic to query object | complex scopes, reporting queries | [query-to-query-object.md](examples/query-to-query-object.md) |
| Extract Current from model | ambient request state in domain | [current-from-model.md](examples/current-from-model.md) |
| Decompose god object with associated objects | god model, large User/Account | [god-object-decomposition.md](examples/god-object-decomposition.md) |
| Replace implicit state machine | timestamp-based status | [implicit-to-explicit-state-machine.md](examples/implicit-to-explicit-state-machine.md) |
| Extract view logic to presenter | template logic, formatting | [view-logic-to-presenter.md](examples/view-logic-to-presenter.md) |
| Form object for complex input | fat controllers, multi-model forms | [complex-input-to-form-object.md](examples/complex-input-to-form-object.md) |

## Topic References

For deep dives on specific topics:

| Topic | Reference |
|-------|-----------|
| Authorization (RBAC, ABAC, policies) | [authorization.md](references/topics/authorization.md) |
| Notifications (multi-channel delivery) | [notifications.md](references/topics/notifications.md) |
| View Components | [view-components.md](references/topics/view-components.md) |
| AI Integration (LLM, agents, RAG, MCP) | [ai-integration.md](references/topics/ai-integration.md) |
| Configuration | [configuration.md](references/topics/configuration.md) |
| Callbacks (scoring, extraction) | [callbacks.md](references/topics/callbacks.md) |
| Current Attributes | [current-attributes.md](references/topics/current-attributes.md) |
| Instrumentation (logging, metrics) | [instrumentation.md](references/topics/instrumentation.md) |

## Go Mechanism References

These preserve the upstream topic mapping, but describe the actual Go mechanism.
They do not recommend installing the original Ruby libraries.

| Original topic | Go mechanism | Reference |
|-----|---------|-----------|
| action_policy | Explicit policy functions | [action-policy.md](references/mechanisms/policy-functions.md) |
| view_component | React components and Storybook | [view-component.md](references/mechanisms/react-components.md) |
| anyway_config | Validated configuration structs | [anyway-config.md](references/mechanisms/typed-config.md) |
| active_delivery | Application operations and delivery ports | [active-delivery.md](references/mechanisms/delivery-ports.md) |
| alba | Explicit DTOs and encoding/json | [alba.md](references/mechanisms/json-dtos.md) |
| workflow | Typed states and guarded methods | [workflow.md](references/mechanisms/explicit-state-machines.md) |
| rubanok | Typed request filters | [rubanok.md](references/mechanisms/typed-filters.md) |
| active_agent | Explicit provider boundary and versioned prompts | [active-agent.md](references/mechanisms/ai-provider-boundary.md) |
| active_job-performs | Typed worker entry points | [active-job-performs.md](references/mechanisms/worker-entrypoints.md) |
| archspec | Check Go package imports in CI | [archspec.md](references/mechanisms/import-boundaries.md) |

## Extraction Signals

Go does not need a callback framework. Use the score only when reviewing existing hooks;
pure transformation and normalization become explicit functions. Keep reliable external
effects outside transactions, with atomic delivery intent when loss is unacceptable.


**When to extract from models:**

| Signal | Metric | Action |
|--------|--------|--------|
| God object | High churn × complexity | Decompose into composition, delegates, or separate models |
| Operation callback | Score 1-2/5 | Extract to service or event handler |
| Code-slicing concern | Groups by artifact type | Convert to composed behavior or extract |
| Current dependency | Model reads Current.* | Pass as explicit parameter |

**Callback scoring:** use the single [canonical scale](references/topics/callbacks.md#callback-scoring-system). Pure transformations and consistency remain explicit local behavior; background/external effects belong to an operation with a deliberate delivery contract.

See [Extraction Signals Reference](references/core/extraction-signals.md) for detailed guide.

## Model Organization

Keep related behavior in a domain-named Go package. Declare types and invariants, then
constructors and public behavior, then private implementation details. Keep persistence
in the package's store file or a separate adapter once the domain boundary needs it.

```go
type Publication struct { publishedAt time.Time }
func (p Publication) Published() bool { return !p.publishedAt.IsZero() }
func (p *Publication) Publish(now time.Time) error {
    if p.Published() { return ErrAlreadyPublished }
    p.publishedAt = now
    return nil
}
```

Use named fields and explicit methods. Do not add Active Record-style associations,
scopes, callbacks, or metaprogramming to reproduce Rails file conventions.

## Success Checklist

Well-layered code:

- [ ] No reverse dependencies (lower layers don't depend on higher)
- [ ] Domain types do not read ambient actor or request state
- [ ] Services don't accept request objects
- [ ] Controllers are thin (HTTP concerns only)
- [ ] Domain logic lives in models, not services
- [ ] Existing hooks are scored; new Go code uses explicit calls
- [ ] Composition is behavioral, not code-slicing
- [ ] Abstractions don't span multiple layers
- [ ] Tests verify appropriate layer responsibilities

## Guidelines

- **Use domain language** - Name models after business concepts (Participant, not User; Cloud, not GeneratedImage)
- **Patterns before abstractions** - Let code age before extracting; premature abstraction is worse than duplication
- **Services as waiting room** - Don't let a generic services package become permanent residence for code
- **Explicit over implicit** - Prefer explicit parameters over ambient state
- **Extraction thresholds** - Investigate mixed responsibilities or external effects; line counts are signals, not mandatory extraction thresholds
