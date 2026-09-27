# Extract Current from Model

Remove ambient request-user access from a domain model by moving authorization to a policy and ownership to the handler.

## Before

```go
// BAD: hidden request identity makes domain behavior depend on HTTP middleware.
var currentUser User
func (post Post) CanEdit() bool {
    return post.AuthorID == currentUser.ID || currentUser.Admin
}
```

## After

```go
type Post struct { ID, AuthorID int64 }
type User struct { ID int64; Admin bool }

func CanEditPost(actor User, post Post) bool {
    return actor.Admin || post.AuthorID == actor.ID
}

// The handler resolves the authenticated actor. The domain receives values.
func NewPost(actor User) Post {
    return Post{AuthorID: actor.ID}
}
```
