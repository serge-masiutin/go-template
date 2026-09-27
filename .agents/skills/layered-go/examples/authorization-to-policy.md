# Extract Authorization to Policy

Replace duplicated handler-side authorization checks with a policy object.

## Before

```go
// BAD: the same permission rule is duplicated across HTTP handlers.
if actor.Admin || post.AuthorID == actor.ID {
    return store.Update(ctx, post)
}
return ErrForbidden

// A second handler grows an independent copy of the rule.
if actor.Admin || post.AuthorID == actor.ID {
    return store.Delete(ctx, post.ID)
}
return ErrForbidden
```

## After

```go
func CanUpdatePost(actor User, post Post) bool {
    return actor.Admin || post.AuthorID == actor.ID
}
func CanDeletePost(actor User, post Post) bool { return actor.Admin || post.AuthorID == actor.ID }

// An application operation enforces the policy for HTTP and worker callers.
func (s *Posts) Update(ctx context.Context, actor User, post Post) error {
    if !CanUpdatePost(actor, post) { return ErrForbidden }
    return s.store.Update(ctx, post)
}
```
