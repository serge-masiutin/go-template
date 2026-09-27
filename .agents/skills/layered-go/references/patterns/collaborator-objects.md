# Collaborator Objects

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns

## Summary

A collaborator object is a small, focused class that owns a slice of behavior tightly coupled to a primary model. It lives alongside the model, takes the model (or a piece of its data) as input, and exposes a narrow interface. The classic shape is a **delegate object**: a separate class — sometimes a separate database table — that the primary model delegates a coherent group of attributes and methods to.

A collaborator is **not** a service object. Services orchestrate operations across the application. Collaborators stay inside the Domain layer; they exist to keep the primary model from accumulating unrelated responsibilities.

## When to Use

- A model has a coherent slice of behavior that doesn't justify a full sibling model but bloats the primary class (contact information, profile metadata, billing details).
- Multiple models share the slice (polymorphic delegates work well here).
- The slice has its own validations, normalizations, or formatting rules that don't belong on the parent.
- A behavioral concern has grown beyond a single mixin's reasonable size and would be clearer as a typed object.

## When NOT to Use

- The slice is just data with no behavior — use a value object instead.
- The slice is a unit of work (creating, sending, syncing) — use a service object.
- The slice is reusable across heterogeneous models with no per-instance state — a concern is fine.
- The slice has its own identity and lifecycle — promote it to a sibling model.

## Key Principles

- **One primary association.** A collaborator is created for and owned by a single model (or polymorphic owner). It does not stand alone.
- **Narrow interface.** The collaborator exposes a focused set of methods — what its slice of behavior provides — not the full persistence API.
- **Domain-layer only.** A collaborator stays inside the Domain layer. It does not call mailers, jobs, services, or external APIs. If it needs to, it's a service in disguise.
- **Delegate explicitly.** Keep a named collaborator field or forwarding method only where it preserves a useful public API.

## Implementation

### Explicit composition

Keep the contact-information scenario, but use a named field and narrow methods instead
of Active Record polymorphic associations and runtime delegation.

```go
type ContactInformation struct {
    Phone string
    Visible bool
}
func (c ContactInformation) PublicPhone() (string, bool) {
    if !c.Visible { return "", false }
    return c.Phone, true
}
type User struct { ID int64; Contact ContactInformation }
```

### Plain delegate and mixed behavior

A pure formatter can accept the values it needs. Do not pass the whole parent when a
small value is sufficient. A collaborator that sends mail, charges a card, or enqueues
work is an application operation, not domain behavior. Persist related rows explicitly;
a Go field does not establish a database relationship or cascade.

## Anti-Patterns

### Cross-Layer Delegate


Calling mailers, jobs, or external APIs from a delegate makes it an application service in domain clothing. Move it to the owning application package or to the dedicated notification layer.

### Anemic Delegate


If the delegate has no behavior of its own, it is paperwork. Either inline the methods on the parent, or promote to a value object if the attributes deserve their own type.

### Delegate That Knows About Callers

A collaborator should not branch on who called it (a branch on an ambient handler actor). It owns a slice of one model's behavior; cross-cutting context belongs in the calling layer.
