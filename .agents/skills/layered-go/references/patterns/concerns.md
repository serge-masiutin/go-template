# Concerns

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Concern Health Check
- Implementation
- Anti-Patterns
- File Organization
- Concern Checklist

## Summary

Concerns are focused composed types/functions that extend domain capabilities, but require discipline to avoid becoming a dumping ground for unrelated code.

## When to Use

- Shared behavior across multiple models
- Cohesive feature sets (soft delete, auditing, slugging)
- Interface extraction (common protocol for different models)
- Reducing model file size while maintaining cohesion

## When NOT to Use

- Single-model "organization" (hiding complexity)
- Unrelated methods grouped together
- Callbacks that should be explicit operations
- As a substitute for proper abstractions

## Key Principles

- **Cohesion over organization** — concerns group related behavior, not random methods
- **Explicit dependencies** — declare required values and dependencies
- **Test concerns independently** — don't rely on specific model setup
- **Prefer composition** — use an operation only for use-case orchestration

## Concern Health Check

Signs of healthy composition: one named responsibility, small explicit inputs, independent behavior tests, and reuse with identical semantics.

Signs of unhealthy composition: promoted methods expose unintended APIs; embedded types depend on outer fields; changing one behavior breaks unrelated tests.

## Implementation

### Basic composition

Go has no ActiveSupport concerns or mixin callbacks. Use named fields, focused functions,
and consumer-owned interfaces. Embedding promotes methods; use it only when that exposed
API is intentional, not as a substitute for inheritance.

```go
type Publication struct { PublishedAt time.Time }
func (p *Publication) Publish(now time.Time) error {
    if !p.PublishedAt.IsZero() { return ErrAlreadyPublished }
    p.PublishedAt = now
    return nil
}
type Post struct { ID int64; Publication Publication }
```

### Configuration and required interface

Pass dependencies and settings explicitly through constructors. Define a small interface
where a consumer needs interchangeable behavior. Do not manufacture a universal domain
interface or infer dependencies from naming conventions.

### Testing composition

Test the composed behavior directly with `testing`; keep storage filtering tests in the
persistence package. A soft-delete field does not automatically add a query scope: every
relevant query must express its visibility contract explicitly.

## Anti-Patterns

### Single-Model Concerns

Splitting one large type into files named validations/callbacks/methods moves text but preserves coupling. Keep the type together until a coherent behavior can be composed. A single consumer can still justify a collaborator when it owns a real concept.

### Dependency Hiding

Do not require undocumented fields on an outer type. A function accepts what it uses; a constructor receives dependencies. Named fields expose ownership better than embedding used as inheritance.

### Concern Chains

Avoid nested embedding that promotes a large accidental method set. Trace dependencies and replace the chain with one explicit composition boundary.

## File Organization

```
internal/posts/
├── post.go
├── publication.go
├── publication_test.go
└── store.go
```

## Concern Checklist

Before extracting composed behavior, verify:

- [ ] Will multiple models use this?
- [ ] Are all methods cohesively related?
- [ ] Is the interface documented?
- [ ] Does it orchestrate a use case, making an operation more appropriate?
- [ ] Are callbacks necessary or hiding complexity?
