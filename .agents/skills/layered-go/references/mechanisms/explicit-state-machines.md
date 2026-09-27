# Explicit State Machines

This replaces the Workflow gem's event DSL and callbacks with typed states and guarded Go methods.

## Basic Usage

```go
type OrderState string
const (Pending OrderState = "pending"; Paid OrderState = "paid"; Shipped OrderState = "shipped")
func (o *Order) Ship(now time.Time) error {
    if o.State != Paid { return ErrInvalidTransition }
    o.State, o.ShippedAt = Shipped, now
    return nil
}
```

Validate stored/input enum values. Prefer named events over unrestricted status assignment. For a two-state flag, a full framework is unnecessary.

## Triggering and Available Events

Expose allowed actions through pure predicates for presentation, but re-evaluate the transition in the operation. UI availability is not permission or a concurrency guarantee.

## Per-Event Effects and Guards

Guards express domain prerequisites using supplied facts. They do not call providers. The operation persists the transition and required outbox intent atomically; workers deliver after commit.

## State-Dependent Behavior

Keep behavior near the state owner. If the workflow is secondary and independent of the entity's lifecycle, compose a separate state value/type instead of expanding a god object.

## Generic Transition Hooks

An explicit audit observer may record transition facts, but it must not become a hidden workflow engine. Required audit records belong to the same transaction; telemetry has separate failure semantics.

## Querying by State

Write parameterized SQL with an explicit state filter. Go methods do not generate scopes. A conditional update or row lock protects concurrent transitions; check affected rows.

## Testing

Enumerate allowed and rejected transitions, timestamps and unknown states. Use two independent database connections to test concurrent commands. Test provider timeout and success-before-local-ack failure for payment-like workflows; preserve the same durable attempt/idempotency key during reconciliation.
