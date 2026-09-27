# Extract God Object with Associated Objects

Decompose a sprawling `User` model into focused collaborators using the associated-objects pattern.

## Before

```go
// BAD: unrelated dependencies are concentrated in one object.
type User struct {
    ID int64
    PasswordHash []byte
    PaymentClient *PaymentClient
    Mailer *Mailer
    PushClient *PushClient
}
// Authentication, billing, and notification delivery now share one lifecycle.
```

## After

```go
type User struct { ID int64; Email string }
type NotificationSettings struct {
    EmailEnabled bool
    PushEnabled bool
    DigestFrequency DigestFrequency
}
type Billing struct {
    subscriptions SubscriptionStore
    payments PaymentGateway
}
// Billing owns the workflow; the gateway implements provider-specific HTTP.
func (b *Billing) Subscribe(ctx context.Context, user User, plan Plan, key string) error {
    subscription, err := b.payments.CreateSubscription(ctx, user.ID, plan.PriceID, key)
    if err != nil { return fmt.Errorf("create subscription: %w", err) }
    return b.subscriptions.Record(ctx, user.ID, subscription)
}
// The external success/local failure case needs reconciliation using the same
// idempotency key. A database transaction cannot roll back the payment provider.
```
