# Extract View Logic to Presenter

Move complex view computations out of templates into a presenter wrapping the model.

## Before

```tsx
<span>{user.status === "active" ? "Active member" : user.status === "pending" ? "Awaiting approval" : "Inactive"}</span>
// Repeating display decisions across pages makes them drift.
```

## After

```go
type UserSummary struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Status string `json:"status"`
    JoinedAt time.Time `json:"joinedAt"`
}
// Keep machine values in the Go DTO. React owns locale and visual presentation.
```

```tsx
const labels = { active: "Active member", pending: "Awaiting approval", inactive: "Inactive" } as const;
function MemberStatus({ status }: { status: keyof typeof labels }) {
  return <span>{labels[status]}</span>;
}
// Use one shared component; add its distinct states to Storybook.
```
