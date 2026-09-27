# Architecture-Check Setup Workflow

Adapt the original detect/configure/verify/report procedure to Go package imports. The starter does not use the Ruby ArchSpec gem.

## 1. Inspect Existing Tools

Read `go.mod`, CI and `go list -json ./...`. Reuse an existing lint/import rule where available. The compiler rejects import cycles, but architectural restrictions require an explicit project rule.

## 2. Detect Actual Boundaries

Identify pure domain packages, application consumers and infrastructure adapters. Do not forbid pgx in a feature package that deliberately owns its store. A package split is justified by a useful boundary, not the checker itself.

## 3. Define the Smallest Rule

For each protected package scope, specify forbidden imports and narrow exceptions with reasons. Derive the module prefix from `go.mod`. Check tool errors and malformed output; never treat them as an empty dependency graph. Include tests/generated files only according to an explicit policy.

## 4. Verify

Run the rule against the current project and a temporary intentional violation. Confirm allowed adapters pass and forbidden edges fail. Run compilation and relevant behavior tests too; import rules cannot prove transaction or authorization correctness.

## 5. Report

State the enforced boundary, files/config changed, verification performed and known limitations. If no pure boundary warrants a rule yet, explain that finding rather than adding an empty architecture framework.

See [import-boundary mechanism](../references/mechanisms/import-boundaries.md).
