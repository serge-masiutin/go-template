# Value Objects

## Contents

- Summary
- When to Use
- When NOT to Use
- Key Principles
- Implementation
- Anti-Patterns
- Performance Note
- Go tools

## Summary

Value objects encapsulate domain concepts that are defined by their attributes rather than identity. They're immutable, comparable by value, and often created from multiple database columns or JSON stores.

## When to Use

- Domain concepts like Money, Address, Coordinates
- Groups of related attributes that travel together
- Attributes with behavior (formatting, validation, comparison)
- JSON/JSONB column wrappers with validation

## When NOT to Use

- Entities with identity (use models)
- Simple scalar values without behavior

## Key Principles

- **Immutable** — value objects don't change after creation
- **Equality by value** — two objects with same attributes are equal
- **No identity** — no `id` field, not persisted directly
- **Self-contained** — all behavior related to the concept

## Implementation

### Defined types and value semantics

Use a struct or defined scalar type with a constructor that validates external values.
Keep fields unexported when mutation would break invariants. Value receivers do not make
maps or slices immutable; copy mutable members when ownership crosses a boundary.

```go
type MediaType string
func ParseMediaType(raw string) (MediaType, error) {
    parsed, _, err := mime.ParseMediaType(raw)
    if err != nil { return "", err }
    return MediaType(parsed), nil
}
func (m MediaType) IsImage() bool { return strings.HasPrefix(string(m), "image/") }
```

### Multiple columns and JSON stores

Map columns or decoded JSON into one validated type at the persistence boundary.
Use explicit scanning/mapping functions. Do not recreate `composed_of`, runtime accessors,
or a JSON store-model DSL. A missing required address is an error, not an empty address.

### Money Example

Represent minor units in an integer with an explicit currency. Define supported currency
precision, rounding, and overflow rules before adding arithmetic. Reject currency mismatch;
never sum money in `float64` or assume every currency has two decimal places.

```go
type Money struct { minor int64; currency Currency }
func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency { return Money{}, ErrCurrencyMismatch }
    if other.minor > 0 && m.minor > math.MaxInt64-other.minor ||
       other.minor < 0 && m.minor < math.MinInt64-other.minor { return Money{}, ErrOverflow }
    return Money{minor: m.minor + other.minor, currency: m.currency}, nil
}
```

## Anti-Patterns

### Mutable Value Objects

Unexported fields protect invariants only if methods do not leak slices or maps. Copy mutable members or document exclusive ownership. A zero value must be either valid or rejected before use.

### Value Objects With Side Effects

Formatting money or comparing coordinates must not persist, fetch exchange rates or read global configuration. Supply external facts explicitly; currency conversion with a provider is an application operation.

## Performance Note

Measure Go allocations and execution time with `go test -bench . -benchmem` and profiles before changing value semantics. Ruby benchmarks do not transfer to Go. Copying small values may simplify ownership; copying large slices or retaining backing arrays requires measurement.

## Go tools

Use the standard library and the dependencies already pinned in `go.mod`. A framework
is not required for this pattern; add one only for an established requirement.
