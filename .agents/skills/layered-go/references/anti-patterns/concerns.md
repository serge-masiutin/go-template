# Concern Anti-Patterns

Misuses of embedding and composed behavior that hide structural problems.

## Code-Slicing Concerns

**Problem:** Splitting a type by artifact kind, not behavior.

```text
BAD: user_validations.go, user_callbacks.go and user_methods.go partition code by
artifact kind while all methods still depend on all User fields. File movement
alone is not decomposition.
```

**Test:** If removing this concern breaks unrelated tests, it's code-slicing.

**Fix:** Keep in model or extract to value object.

```go
// GOOD: a cohesive value with its own behavior.
type ContactInformation struct { Phone string; Public bool }
func (c ContactInformation) PublicPhone() (string, bool) {
    return c.Phone, c.Public
}
```

## Overgrown Concerns

**Problem:** Concern has too many responsibilities.

```text
BAD: embedding AccountFeatures promotes billing, search, notification and export
methods into every entity. The large method set hides dependencies and ownership.
```

**Fix:** Extract to delegate object or separate concerns.

```text
Keep pure contact/billing values near their owner. Extract an export operation
when it orchestrates I/O. Compose fields by name rather than promoting unrelated
methods through nested embedding.
```
