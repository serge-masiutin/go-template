# Serializers

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

Serializers are specialized presenters for API responses. They form a dedicated abstraction layer between models and the external API contract, providing consistent JSON formatting and the ability to generate TypeScript types.

## When to Use

- API responses beyond simple JSON
- Multiple JSON formats for different contexts
- TypeScript frontend integration
- Versioned APIs

## When NOT to Use

- An already explicit, safe response DTO: marshal it directly
- Internal data transformation

## Key Principles

- **Keep mapping explicit** — avoid domain `MarshalJSON` methods that encode one UI contract
- **Explicit mapping** — call a mapper with loaded, authorized values
- **One serializer per context** — different serializers for list vs detail
- **Check TypeScript contracts** — this starter uses explicit types and response tests; introduce generation only with a real generator

## Implementation

### Explicit response DTO

Go's exported struct fields and `encoding/json` tags form the allowlist. Do not marshal
persistence structs containing password hashes, provider payloads, or private metadata.
Serialization is pure: no query, mutation, authorization lookup, or network request.

```go
type UserResponse struct {
    ID int64 `json:"id,string"`
    Email string `json:"email"`
}
func UserDTO(user User) UserResponse { return UserResponse{ID: user.ID, Email: user.Email} }
```

### Entity, page, and shared shapes

Compose entity DTOs into explicit page props. Keep shared props small and request-scoped.
A bigint ID is a decimal string for JavaScript; timestamps use RFC3339 with a documented
precision; empty arrays are `[]`; optional fields distinguish omission from null.

### Associations and loading

Authorize and load data before mapping it. Deferred/optional loading belongs to Gonertia
props, not an implicit lazy serializer. Pair Go response tests with TypeScript consumers.
See `go-serialization` for the installed APIs and `inertia-go-typescript` for types.

## Anti-Patterns

### UI-specific MarshalJSON on Domain Types

A domain type used by admin, public, and internal consumers cannot have one safe implicit representation. Define separate DTOs; verify exact exposed keys, not only expected values.

### Multiple Formats Without Serializers

Do not scatter anonymous map construction across endpoints. Reuse an entity DTO where semantics match, then compose page-specific props. Distinct permissions can require distinct DTOs.

## Testing

Marshal the DTO and assert public keys, decimal-string IDs, RFC3339 timestamps, empty arrays, nullable fields and absence of passwords. Add request tests proving unauthorized associations cannot enter the DTO. TypeScript compilation alone cannot prove the server matches its declarations.

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
