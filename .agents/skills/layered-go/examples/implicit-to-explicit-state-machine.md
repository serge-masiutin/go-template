# Replace Implicit State Machine

Replace ad-hoc timestamp checks with an explicit state machine using typed states and guarded methods.

## Before

```go
func (o Order) CanShip() bool {
    return !o.PaidAt.IsZero() && o.ShippedAt.IsZero() && o.CancelledAt.IsZero()
}
// Multiple independent timestamps can represent contradictory states.
```

## After

```go
type State string
const (
    Pending State = "pending"
    Paid State = "paid"
    Shipped State = "shipped"
    Delivered State = "delivered"
    Cancelled State = "cancelled"
)
func (o *Order) Ship(now time.Time) error {
    if o.State != Paid { return ErrInvalidTransition }
    o.State, o.ShippedAt = Shipped, now
    return nil
}
// Persist with a conditional update or locked transaction, and check RowsAffected.
// WHERE id=$1 AND status='paid' prevents two racing shipments.
// Notification intent commits with the transition; delivery happens afterwards.
```
