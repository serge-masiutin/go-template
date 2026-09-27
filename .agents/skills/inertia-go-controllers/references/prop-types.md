# Prop Types — Detailed Reference

Prop loading and merge behavior checked against Gonertia v3.0.0. Domain/store names in snippets illustrate a feature to implement.

## Table of Contents

- [Regular Props](#regular-props)
- [Optional Props](#optional-props)
- [Deferred Props](#deferred-props)
- [Once Props](#once-props)
- [Merge Props](#merge-props)
- [Always Props](#always-props)
- [Scroll Props](#scroll-props)
- [Resetting Merge/Scroll Props](#resetting-mergescroll-props)
- [Combining Prop Types](#combining-prop-types)
---

## Regular Props

Plain values are included initially. A loader must have one of Gonertia's supported signatures, such as `func(context.Context) (any, error)`; a function returning `([]User, error)` is not interchangeable with that signature.

```go
props := inertia.Props{
    "filters": filterDTO,
    "users": func(ctx context.Context) (any, error) {
        return store.List(ctx, actor.ID)
    },
}
```

During a partial reload only selected loaders run. Keep loaders pure queries, scoped to the current actor, cancellable, and safe for Gonertia's concurrent resolution. Do not share a single pgx connection/transaction between parallel loaders.

## Optional Props

```go
"exportData": inertia.Optional(func(ctx context.Context) (any, error) {
    return reports.ExportPreview(ctx, actor.ID)
}),
```

Optional data is absent initially and evaluated only when requested. The client can use `router.reload({ only: ['exportData'] })`. Large file exports use a download route, not a giant prop.

## Deferred Props

```go
"reviews": inertia.Defer(func(ctx context.Context) (any, error) {
    return reviews.List(ctx, actor.ID, courseID)
}),
"chart": inertia.Defer(func(ctx context.Context) (any, error) {
    return analytics.Chart(ctx, actor.ID, courseID)
}, "analytics"),
```

Deferred values load in subsequent requests, not a streaming response. Every request repeats authentication and authorization. The initial payload contains deferred metadata rather than a null value for each prop.

```tsx
<Deferred data="reviews" fallback={<ReviewsSkeleton />}>
  <ReviewsList />
</Deferred>
```

## Once Props

`inertia.Once(value)` in Gonertia v3.0.0 includes the prop initially and skips it on ordinary partial reloads unless requested. It does **not** implement every cross-navigation cache/expiry option from other server adapters. Do not describe it as a server session cache or use it for revocable permissions. Verify protocol support before adding newer client once-prop features.

## Merge Props

```go
"messages": inertia.Merge(messageDTOs).MatchOn("id"),
"feed": inertia.Merge(feedDTO).Append("data").MatchOn("data.id"),
```

Merge applies during partial reloads; a full visit replaces the prop. IDs must be stable. Avoid appending every historical row on each poll; define the batch/cursor contract and test duplicate handling.

```tsx
router.reload({ only: ['messages'] })
```

## Always Props

```go
ctx := inertia.SetProp(r.Context(), "csrfToken", inertia.Always(token))
```

Use Always for small values that must accompany partial responses. It does not bypass authorization or turn a global shared value into request-local state.

## Scroll Props

Gonertia does not include a paginator. Query a bounded page and supply explicit metadata. Its metadata contract uses `any` because page cursors need not be integers:

```go
type PageMetadata struct { Current int; Previous, Next any }
func (p PageMetadata) GetPageName() string { return "page" }
func (p PageMetadata) GetPreviousPage() any { return p.Previous }
func (p PageMetadata) GetNextPage() any { return p.Next }
func (p PageMetadata) GetCurrentPage() any { return p.Current }

prop := inertia.Scroll(posts, inertia.WithMetadata(metadata)).ConfigureMergeIntent(r)
```

The default wrapper is `data`, so the React contract is `{ posts: { data: Post[] } }`. `Previous`/`Next` are nil at the boundary. Validate page/cursor inputs, stable ordering and tenant scope independently. Do not port Pagy methods.

```tsx
<InfiniteScroll data="posts" loading={() => <PostsSkeleton />}>
  {posts.data.map(post => <PostCard key={post.id} post={post} />)}
</InfiniteScroll>
```

Use the installed client API for manual load-more controls. Test bidirectional merge intent and filter resets before claiming infinite scrolling works.

## Resetting Merge/Scroll Props

```tsx
router.reload({ only: ['messages'], reset: ['messages'] })
```

Reset accumulated data when filters/sort change. A fresh full visit also replaces prior data.

## Combining Prop Types

Gonertia v3.0.0 supports `inertia.Defer(loader, "feed").Merge()` and methods on `MergeProps`. Do not assume Rails keyword combinations such as optional+merge or arbitrary expiry exist. Check `response.go` in the pinned module and add response tests for every combination used by a feature.
