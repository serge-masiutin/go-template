# Form Objects

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns
- Go tools

## Summary

Form objects handle specific user interactions involving data submission. They belong to the presentation layer and are useful when the form doesn't map 1:1 to a model — creating multiple records, updating virtual attributes, or handling complex validation contexts.

## When to Use

- Multi-model forms (creating/updating multiple records)
- Model-less forms (feedback, search, settings)
- Context-specific validations (publish vs draft)
- Complex form logic that doesn't belong in models

## When NOT to Use

- Simple single-record input (a request DTO and store/domain call are enough)
- Duplicating model validations

## Key Principles

- **Presentation layer abstraction** — handles UI concerns, not domain logic
- **Use when N ≠ 1** — forms that create/update zero, multiple, or virtual resources
- **Typed boundary** — decode and validate an explicit input struct
- **Explicit atomicity** — operations own transactions; external effects happen after commit
- **Don't duplicate validations** — delegate to models and merge errors

## Implementation

### Typed inputs, explicit operations

Use a request struct and a boundary parser. Form input belongs to presentation; a
multi-record registration transaction belongs to the application operation. Do not
translate Active Model callbacks or validation DSLs literally into Go.

```go
type InvitationInput struct {
    Email string `json:"email"`
    SendCopy bool `json:"sendCopy"`
}
func (input InvitationInput) Validate() error {
    address, err := mail.ParseAddress(input.Email)
    if err != nil || address.Address != input.Email { return ErrInvalidEmail }
    return nil
}
```

### Context-specific and multi-record forms

Pass the authenticated actor separately from the input. Clients cannot submit ownership
or role fields unless the specific operation allows them. Delegate domain invariants
to the owning domain type. Commit all related changes together and create any reliable
notification intent in that transaction.

### Model-less forms

Feedback/search need no artificial persistent model. Parse the request into a typed
input and call the operation or query. Durable asynchronous feedback still needs durable
storage; a goroutine started by the request is not a queue.

### Usage in handlers

For this starter, use strict JSON decoding, Inertia validation errors persisted across
303 redirects, and React `<Form>` with the CSRF header. Malformed structure is 400.
A separate JSON API may use 422; do not send a bare 422 JSON object to an Inertia page.
See [complex input extraction](../../examples/complex-input-to-form-object.md).

### Wizard Forms (Multi-Step)

Store server-owned progress explicitly when it must survive navigation. Use typed states
and allowed transitions; validate the current step and the complete final submission.
Never treat a client-supplied completion flag as proof that earlier checks passed.

## Anti-Patterns

### Duplicating Model Validations

The boundary checks JSON shape and UI-specific constraints. The domain owns durable invariants. Convert typed domain failures into field errors without maintaining a second independent rule set.

### UI Logic in Model Callbacks

Publishing a draft may require different fields from saving it. Use explicit commands and validation contexts; do not add a model flag that changes hooks depending on which form submitted it.

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
