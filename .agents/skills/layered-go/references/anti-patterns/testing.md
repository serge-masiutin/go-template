# Testing Anti-Patterns

Tests that verify the wrong layer's responsibilities.

## Testing Wrong Layer

**Problem:** HTTP tests enumerate domain rules instead of transport behavior.

```text
BAD: every HTTP test enumerates all discount thresholds and sets up sessions,
PostgreSQL and rendering merely to exercise arithmetic. A failure cannot identify
whether the HTTP contract or the pricing rule broke.
```

**Fix:** Test business logic in domain tests.

```go
func TestOrderDiscount(t *testing.T) {
    for _, test := range []struct{ total, want int64 }{{9000, 0}, {11000, 1100}} {
        order := Order{Total: test.total}
        if got := order.Discount(); got != test.want {
            t.Fatalf("total %d: got %d, want %d", test.total, got, test.want)
        }
    }
}
// Keep separate HTTP tests for decoding, authorization and response contracts.
```
