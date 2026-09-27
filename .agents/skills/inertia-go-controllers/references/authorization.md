# Authorization Props

The server owns access. A `can` prop controls UI visibility only; replaying a hidden
request must still fail. Preserve the original per-resource permission pattern.

```go
type Permissions struct { Update bool `json:"update"`; Delete bool `json:"delete"` }
permissions := Permissions{Update: CanUpdate(actor, post), Delete: CanDelete(actor, post)}
```

```tsx
{post.can.update && <Link href={`/posts/${post.id}/edit`}>Edit</Link>}
```

Resolve the actor from the authenticated session. Scope collections and mutations by
ownership/tenant in SQL. Keep pure policy rules testable without HTTP. If an operation
has worker/CLI callers, enforce the same policy at the shared operation boundary where
appropriate. Test owner, non-owner, missing session, admin, revoked admin, and foreign
record IDs. Never serialize private records just to let React hide them.
