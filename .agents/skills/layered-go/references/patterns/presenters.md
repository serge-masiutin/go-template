# Presenters

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns
- Testing
- Go tools

## Summary

Presenters encapsulate representation logic for views. They bridge models and the view layer, extracting UI-specific formatting, CSS classes, and display logic from domain objects.

## When to Use

- View-specific formatting (dates, names, statuses)
- CSS class generation based on state
- Combining data from multiple models for display
- Logic that's only needed in views

## When NOT to Use

- Domain logic (belongs in models)
- Data transformation for APIs (use serializers)
- Simple attribute access

## Key Principles

- **Presentation layer** — never leak presenters to lower layers
- **Explicit interface (closed)** vs **delegation (open)** — choose based on isolation needs
- **One presenter per view context** — avoid god presenters
- **Test presenters independently** — no full request cycle needed

## Implementation

### Presentation stays at the consumer

Use a small function or value struct for a stable server presentation shape. React owns
visual choices, locale-sensitive formatting, and component composition. Do not expose
all methods of a domain object through reflection or an embedded persistence model.

```go
type UserSummary struct {
    ID string `json:"id"`
    Name string `json:"name"`
    JoinedAt time.Time `json:"joinedAt"`
}
func Summarize(user User) UserSummary {
    return UserSummary{ID: strconv.FormatInt(user.ID, 10), Name: user.Name, JoinedAt: user.CreatedAt}
}
```

### Delegation and collections

A presenter receives all required values; it does not lazily query associations.
Preload data at the query boundary. Keep page-only labels out of reusable domain values.
Use React components and Storybook for visual variants; see
[presenter extraction](../../examples/view-logic-to-presenter.md).

## Anti-Patterns

### Representation Logic in Models

A domain status method may return a domain enum; CSS classes and translated display strings belong to React or a presentation mapper. Test domain meaning separately from presentation.

### Leaking Decorators

Never pass a presenter back to an operation expecting a domain entity. Keep a closed DTO instead of embedding the entity and exposing every field.

### Global Helpers with Prefixes

Repeated `userStatusLabel`, `userColor`, and `userBadge` functions may belong to one feature component or focused presentation module. Do not create a global helper registry.

## Testing

Use table-driven mapper tests for nullability, identifiers and timestamps. React interaction tests cover visible behavior; Storybook covers representative loading, empty, error and permission states. Mapping must not issue SQL.

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
