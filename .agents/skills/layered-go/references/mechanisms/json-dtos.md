# Explicit JSON DTOs

This replaces Alba serializers and automatic lookup with Go response structs and explicit mapping.

## Setup and Basic Usage

Use `encoding/json` and the installed Gonertia renderer. Define only public fields in a response DTO. Never marshal a persistence record containing password hashes or provider metadata.

```go
type UserDTO struct {
    ID int64 `json:"id,string"`
    Name string `json:"name"`
}
func PublicUser(user User) UserDTO { return UserDTO{ID: user.ID, Name: user.Name} }
```

## Associations and Collections

Load authorized associations before mapping; mappers issue no SQL. Initialize loaded empty slices so JSON contains `[]`, not `null`. Represent unrequested deferred data separately from a loaded empty result.

## Conditional Attributes

Use separate DTOs for distinct disclosure contracts, or explicit optional fields when optionality is part of one contract. An admin-only field must not be populated merely because a client asks for it.

## Computed Attributes and Root Keys

Pure calculations can happen during mapping. Choose stable root/page keys deliberately; avoid reflection-derived names. IDs crossing into JavaScript are decimal strings, timestamps are RFC3339 and monetary values have explicit units/currency.

## Shared Serialization and Key Naming

Reuse stable entity DTOs, then compose page-specific props. Use explicit JSON tags for camelCase. No base serializer or naming-convention lookup is required.

## TypeScript

The starter keeps explicit TypeScript types and response tests. Do not claim generated types exist. If introducing a generator, commit its inputs/configuration, run it in CI and fail on drift. See `go-serialization` and `inertia-go-typescript`.

## Performance and Testing

Select needed columns, avoid N+1 loads and bound collections. Measure before caching. Test exact keys, nullability, empty lists and absence of secret fields at the HTTP response boundary.
