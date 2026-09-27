# Go Import-Boundary Checks

This replaces ArchSpec's Ruby constant rules with Go package/import analysis. The compiler already rejects cycles, but it does not know which imports violate your architecture.

## Division of Labor

Compile/type checking validates language contracts. An import check enforces selected package boundaries. Human review and specification tests assess responsibility, hidden effects and inappropriate abstractions.

## Setup

Inspect `go list -json ./...` and actual package ownership before writing a rule. This starter deliberately keeps small feature stores and behavior together; do not falsely classify the whole feature package as a pure domain package.

## Reference Rule

If a project later introduces `internal/domain/...`, its CI check can reject imports of `net/http`, pgx and application adapters from that subtree. Derive the module prefix from `go.mod`; inspect both direct and transitive dependencies when the rule requires it. Fail on tool errors rather than treating missing output as an empty import graph.

## Domain Services

Pure calculations may span entities without becoming application operations. Query implementations that import pgx are persistence adapters even if their contract is consumed by domain/application code.

## Exceptions

An exception names an exact package, a reason and an owner. Do not allow an entire directory to import everything merely because one adapter needs HTTP. Narrow exceptions must retain the other domain rules.

## Adopting in Existing Code

Start with observed boundaries and fix one violation at a time. If using a baseline, record exact existing edges and fail on new violations; update the baseline only through review. Do not hide failures behind `|| true`.

## Per-Package Boundaries

Use domain-named packages and consumer-owned interfaces where needed. Separate adapters only when it clarifies a real boundary. An architecture checker is not a reason to create empty layers.

## Limits

Import analysis cannot prove tenant isolation, transaction atomicity, correct permissions, absence of a network call behind an interface, or useful abstractions. Run [the specification test](../core/specification-test.md) and behavioral tests too.
