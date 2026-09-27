# Go Request Test Patterns

Use `httptest.NewServer` when cookie/redirect behavior matters and `httptest.NewRecorder` for a focused handler contract. Keep redirects disabled initially so a mutation's status is observable.

```go
client := &http.Client{
    Jar: jar,
    CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}
```

Fetch the page first, read its CSRF token and Inertia asset version, then submit strict JSON with the same cookie jar and `X-CSRF-Token`. Follow an expected 303 manually. Decode response bodies only after checking status/content type.

For partial requests set the component and requested prop headers explicitly. Prove deferred loaders did not run initially, then check their values and authorization on the second request. For private DTOs assert absent sensitive keys as well as expected public values.

Use an isolated test schema/database; migrations and constraints are real. Never point integration tests at production. Cleanup failures must fail the test, and asynchronous work must terminate before teardown.
