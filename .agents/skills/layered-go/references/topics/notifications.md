# Notifications

## Summary

Separate **what happened**, **who should be notified**, and **how delivery occurs**. Preserve the source's delivery abstraction and channel/preferences decisions without porting Active Delivery's dynamic API.

## Layer Placement

Domain types express state and events. Application operations decide notification intent and recipients. SMTP, HTTP push and queue adapters implement transport. A notification is not an entity lifecycle hook.

## Key Principles

Keep external effects out of database transactions. If losing a notification is unacceptable, store its intent atomically with the business change. Make duplicate delivery, retry limits and failure visibility explicit.

## Implementation with Delivery Ports

### Basic Setup

```go
type WelcomeMessage struct { EventID string; RecipientID int64 }
type WelcomeSender interface {
    SendWelcome(context.Context, WelcomeMessage) error
}
```

Define an interface at its consumer when a transport boundary actually exists. The starter does not include an email provider, worker or outbox; these are design contracts for adding one.

### Triggering Notifications

An operation saves the user and a versioned outbox event in one transaction. The worker claims committed events with a bounded lease, loads recipient details and calls the delivery adapter. It records success only after acknowledged delivery.

A crash after provider success but before acknowledgement can duplicate delivery. Use a stable event ID as the provider idempotency key where supported; otherwise document at-least-once delivery and duplicate tolerance. Outbox alone does not guarantee exactly once.

### Conditional Delivery

Channel selection belongs to the notification operation. Return a typed result distinguishing delivered, intentionally suppressed and failed. Do not turn an unknown transport error into a successful suppression.

### User Preferences

Load current opt-outs before dispatch and define which mandatory notices are exempt. Do not copy stale preferences into a long-lived job unless the product explicitly requires submission-time preferences.

### Testing Deliveries

Use fake transport adapters at the boundary. Test recipient/channel selection, opt-outs, template variables and provider failures. Integration tests cover transaction rollback, worker claims, lease expiry and duplicate processing.

## Without a Delivery Framework

One explicit operation and one provider adapter are enough for a single channel. Add shared channel orchestration only when independent notifications repeat the same semantics. Do not create a notification base class or reflection registry.

## Triggering Notifications from State Machines

The transition records state and notification intent together. A guard must not send mail. The worker handles external delivery after commit; a failed delivery does not roll the state backward without an explicit compensating business operation.

## Anti-Patterns

### Notifications in Models

A domain method that calls SMTP or enqueues a job crosses layers and makes ordinary state changes produce surprising effects. Move the call to the existing use-case operation.

### Transport Knows Too Much

An SMTP adapter should not decide whether an invoice is overdue or which account owns it. It receives an explicit message and returns a transport result/error.

### Inline Notification Logic

Repeated recipient selection and template assembly in handlers hides a stable responsibility. Extract that responsibility without creating an abstraction for each trivial message.

## Delivery Patterns

### Batch Notifications

Use bounded batches, explicit cursor ordering and per-recipient outcomes. Never load every recipient into memory or spawn an unbounded goroutine per user. Define retry behavior for partial batches.

### Notification Objects

Use a typed, versioned event payload with stable identifiers. Avoid storing entire mutable domain objects, secrets or raw provider responses in the queue. Persist only the data required by the delivery contract.
