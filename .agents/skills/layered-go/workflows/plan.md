# Gradual Layerification Workflow

Plan incremental adoption of layered design around a concrete goal. Preserve behavior and the project's useful conventions while changing the smallest necessary boundary.

## Inputs

A goal, target scope, current contracts, tests and constraints. If the goal is vague, inspect the relevant code before proposing a pattern. Do not produce a repository-wide redesign for one endpoint.

## Process

### 1. Understand the Goal

State the desired observable outcome and what must remain compatible. Distinguish adding a feature from restructuring existing behavior.

### 2. Assess Current State

Read commands, package imports, composition roots, stores, routes and tests. Use `go list -json ./...` and actual callers. Inspect existing architecture checks rather than assuming an ArchSpec-like tool exists.

### 3. Find Existing Patterns

Find how this project handles operations, input DTOs, errors, permissions and transactions. Reuse coherent conventions. No existing service layer is not by itself a problem.

### 4. Analyze Relevant Code

Apply the specification test. Identify misplaced responsibility and actual costs. Consider explicit composition, a focused query, a policy or a value object only when it solves the observed problem.

### 5. Trace Call Chains

Map all affected entry points, background consumers and external effects. Check queued payload/API compatibility and transaction ownership before moving code.

### 6. Prioritize

Correctness and authority first, then clarity and local reversible changes. Delay optional generalization. Do not mix unrelated package renames with behavior changes.

### 7. Generate Phases

Each phase has a concrete change, affected callers, preserved invariants, test/check commands, migration/rollback considerations and a stop criterion. Characterization tests precede risky extraction; remove the old path only after consumers move.

## Output Format

Goal; current design and evidence; smallest approach; ordered phases with verification; alternatives deliberately deferred and why; remaining unknowns. Do not fill the plan with unneeded layers or unverified command names.

## Guidelines

Prefer functions/struct composition to inheritance, consumer-owned interfaces to a universal repository, explicit commands to callbacks and actor parameters to ambient identity. Package boundaries follow actual responsibility. Use [the eight refactoring examples](../SKILL.md#refactoring-scenarios) as scenarios, not copy-paste production APIs.
