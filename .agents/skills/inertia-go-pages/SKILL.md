---
name: inertia-go-pages
description: >-
  Page components, persistent layouts, Link/router navigation, Head, Deferred, WhenVisible, InfiniteScroll,
  and URL-driven state for Inertia Go. React examples inline; Vue and Svelte equivalents in references.
  Use when building pages, adding navigation, implementing persistent layouts, infinite scroll, lazy-loaded
  sections, or working with client-side Inertia APIs (router.reload, router.replaceProp, prefetching).
---

## Go project contract

This is an adaptation of the original EM Inertia skill, preserving the React examples
and decision workflow. Read `docs/architecture.md` and the actual Go handlers first.
Use Gonertia v3 APIs from `go.mod`; server examples below replace the Rails adapter.
All mutations carry `X-CSRF-Token` from `usePage().props.csrfToken`. The starter accepts
strict JSON; uploads need a separate bounded multipart handler before using upload examples.
Validation errors cross a 303 redirect through the session flash provider. They are not
bare 422 JSON responses. Shared user data is request-scoped, never global `ShareProp`.
Examples describe features to implement, not routes or libraries already installed.


# Inertia Go Pages

Page components, layouts, navigation, and client-side APIs.

**Before building a page, ask:**
- **Does this page need a layout?** → Use persistent layout (React: `Page.layout = ...`; Vue: `defineOptions({ layout })`; Svelte: module script export) — wrapping in JSX/template remounts on every navigation, losing scroll position, audio playback, and component state
- **Does UI state come from the URL?** → Change BOTH controller (read `params`, pass as prop) AND component (derive from prop, no `useState`/`useEffect`) — use `router.get` to update URL
- **Need to refresh data without navigation?** → `router.reload({ only: [...] })` — never `useEffect` + `fetch`
- **Need to update a prop without server round-trip?** → `router.replaceProp` — no fetch, no reload

**NEVER:**
- Parse `window.location.search` or use `useSearchParams` — derive URL state from controller props
- Use `useState`/`useEffect` to sync URL ↔ React state — the controller passes URL-derived data as props; the component just reads them
- Treat the `<Deferred>` render argument as the loaded prop. In React 3.7.1 it contains `{ reloading }`; read loaded data via `usePage()`
- Access `usePage().props.flash` — flash is top-level: `usePage().flash`
- Wrap layout in JSX return for persistence — use `Page.layout = ...` or global layout inside createInertiaApp's resolve callback

## Page Component Structure

Pages are default exports receiving controller props as function arguments.
Use `type Props = { ... }` (not `interface` — causes TS2344 in React). Vue uses `defineProps<T>()`, Svelte uses `let { ... } = $props()`.

```tsx
type Props = {
  posts: Post[]
}

export default function Index({ posts }: Props) {
  return <PostList posts={posts} />
}
```

## Persistent Layouts

Layouts persist across navigations — no remounting, preserving scroll, audio, etc.

```tsx
import { AppLayout } from '@/layouts/app-layout'

export default function Show({ course }: Props) {
  return <CourseContent course={course} />
}

// Single layout
Show.layout = (page: React.ReactNode) => <AppLayout>{page}</AppLayout>
```

Default layout in entrypoint:
```tsx
// web/src/app.tsx
// Keep the starter's explicit `pages = { Login, Home, Admin }` map.
type PageWithLayout = (typeof pages)[keyof typeof pages] & {
  layout?: null | ((page: React.ReactNode) => React.ReactNode)
}

resolve(name) {
  if (!(name in pages)) throw new Error(`Unknown page: ${name}`)
  const page: PageWithLayout = pages[name as keyof typeof pages]
  if (page.layout === undefined) {
    page.layout = content => <AppLayout>{content}</AppLayout>
  }
  return page
}
```

## Navigation

### `<Link>` and `router`

Use `<Link href="...">` for internal navigation (not `<a>`) and `router.get/post/patch/delete`
for programmatic navigation. Key non-obvious features:

```tsx
// Prefetching — preloads page data on hover
<Link href="/users" prefetch>Users</Link>
<Link href="/users" prefetch cacheFor="30s">Users</Link>

// Prefetch with cache tags — invalidate after mutations
<Link href="/users" prefetch cacheTags="users">Users</Link>

// Programmatic prefetch (e.g., likely next destination)
router.prefetch('/settings', {}, { cacheFor: '1m' })

// Partial reload — refresh specific props without navigation
router.reload({ only: ['users'] })
```

Full `router` API, visit options, and event callbacks are in
`references/navigation.md` — see loading trigger below.

### Client-Side Prop Helpers

Update props without a server round-trip:

```tsx
// Replace a single prop (dot notation supported)
router.replaceProp('show_modal', false)
router.replaceProp('user.name', 'Jane Smith')

// With callback (receives current value + all props)
router.replaceProp('count', (current) => current + 1)

// Append/prepend to array props
router.appendToProp('messages', { id: 4, text: 'New' })
router.prependToProp('notifications', (current, props) => ({
  id: Date.now(),
  message: `Hello ${props.auth.user.name}`,
}))
```

These are shortcuts to `router.replace()` with `preserveScroll` and
`preserveState` automatically set to `true`.

**`router.replaceProp` vs `router.reload`:** Use `router.replaceProp` for client-only state changes
(toggling a modal, incrementing a counter) — no server round-trip. Use `router.reload`
when you need fresh data from the server (updated records, recalculated stats).

## URL-Driven State (Dialogs, Tabs, Filters)

URL state = server state = props. **ALWAYS implement both sides:**

1. **Controller** — read `params` and pass as a prop
2. **Component** — derive UI state from that prop (no `useState`, no `useEffect`)
3. **Update** — `router.get` with query params to change URL (triggers server round-trip, new props arrive)

**NEVER** use `useState` + `useEffect` to sync URL ↔ dialog/tab/filter state.
The server is the single source of truth — the component just reads props.

```go
// Parse and validate the optional query value at the HTTP boundary.
selectedID := r.URL.Query().Get("user_id")
// Validate selectedID when present; enforce visibility before exposing a record.
err := engine.Render(w, r, "Users", inertia.Props{
    "users": userDTOs,
    "selected_user_id": selectedID,
})
```

```tsx
// Step 2+3: Derive state from props, router.get to update URL

type Props = {
  users: User[]
  selected_user_id: string  // from controller
}

export default function Index({ users, selected_user_id }: Props) {
  // Derive — no useState, no useEffect, no window.location parsing
  const selectedUser = selected_user_id
    ? users.find(u => u.id === selected_user_id)
    : null

  const openDialog = (id: number) =>
    router.get('/users', { user_id: id }, {
      preserveState: true,
      preserveScroll: true,
    })

  const closeDialog = () =>
    router.get('/users', {}, {
      preserveState: true,
      preserveScroll: true,
    })

  return (
    <Dialog open={!!selectedUser} onOpenChange={(open) => !open && closeDialog()}>
      <DialogContent>{/* ... */}</DialogContent>
    </Dialog>
  )
}
```

**Why not useEffect?** When `router.get('/users', { user_id: 5 })` fires, Inertia
makes a request to the server → Go handler reads `r.URL.Query().Get("user_id")` as `"5"` →
returns new props with `selected_user_id: 5` → component re-renders with the
dialog open. The cycle is: URL → server → props → render. Parsing
`window.location` client-side duplicates what the server already does.

## Shared Props

Shared props (auth, flash) are typed globally via InertiaConfig (see `inertia-go-typescript` skill) — page components only type their OWN props:

```tsx
type Props = {
  users: User[]         // page-specific only
  // auth is NOT here — typed globally via InertiaConfig
}

export default function Index({ users }: Props) {
  const { props, flash } = usePage()
  // props.auth typed via InertiaConfig, flash.notice typed via InertiaConfig
  return <UserList users={users} />
}
```

## Flash Access

**Flash is top-level on the page object, NOT inside props** — this is the #1
flash mistake. Flash config is in `inertia-go-controllers`; toast UI is in `shadcn-inertia`.

```tsx
// BAD:  usePage().props.flash   ← WRONG, flash is not in props
// GOOD: usePage().flash         ← flash.notice, flash.alert
```

## Deferred Failure Contract

Gonertia v3.0.0 returns a loader error from `Render`; the starter's error boundary
responds with 500. It does not produce `rescuedProps`, so React's `rescue` prop cannot
be assumed to handle these failures. A feature using deferred data must define and
test its failure path as well as loading/ready states.

For an unexpected HTTP/network failure, surface an error UI through Inertia's
`httpException` / `networkError` events, scoped to the affected visit, and provide an
explicit retry (`router.reload({ only: ['stats'] })`). Remove listeners on unmount and
do not leave an endless loading placeholder after an error. Do not mistake validation
`onError` (field errors) for a network/500 handler. If a particular dependency outage
is deliberately recoverable, a typed `ready | unavailable` prop can model it, with
boundary logging and a retry policy; never convert arbitrary query errors into an
empty successful result. Retry requests reauthorize and retain the active filter.

## `<Deferred>` Component

Renders fallback until deferred props arrive. Children can be plain `ReactNode`
or `() => ReactNode` render function. Either way, the child reads the deferred
prop from page props via `usePage()` — the render function receives `{ reloading }`, not the loaded prop value.

```tsx
import { Deferred } from '@inertiajs/react'

export default function Dashboard({ basic_stats }: Props) {
  return (
    <>
      <QuickStats data={basic_stats} />
      <Deferred data="detailed_stats" fallback={<Spinner />}>
        <DetailedStats />
      </Deferred>
    </>
  )
}

// Also valid — render function (no args, child still reads from usePage):
// <Deferred data="stats" fallback={<Spinner />}>
//   {() => <Stats />}
// </Deferred>

// BAD — render function does NOT receive data as argument:
// <Deferred data="stats">{(data) => <Stats data={data} />}</Deferred>
```

## `<InfiniteScroll>` Component

Automatic infinite scroll — loads next pages as user scrolls down. Pairs with
`inertia.Scroll` on the server (see `inertia-go-controllers`):

```tsx
import { InfiniteScroll } from '@inertiajs/react'

export default function Index({ posts }: Props) {
  return (
    <InfiniteScroll data="posts" loading={() => <Spinner />}>
      {posts.data.map(post => <PostCard key={post.id} post={post} />)}
    </InfiniteScroll>
  )
}
```

Props: `data` (prop name), `loading` (fallback), `manual` (button instead of auto),
`manualAfter={3}` (auto for first 3 pages, then button), `preserveUrl` (don't update URL).

## `<WhenVisible>` Component

Loads data when element enters viewport. Use for **lazy sections** (comments,
related items), NOT for infinite scroll (use `<InfiniteScroll>` above):

```tsx
import { WhenVisible } from '@inertiajs/react'

<WhenVisible data="comments" fallback={<Spinner />}>
  <CommentsList />
</WhenVisible>
```

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Layout remounts on every navigation | Wrapping layout in JSX return instead of `Page.layout` | Use persistent layout |
| `Deferred` children never render | Render function expects args `{(data) => ...}` | The argument is `{ reloading }`, not fetched data. Read the prop via `usePage()`; plain `<Child />` also works |
| Flash is `undefined` | Accessing `usePage().props.flash` | Flash is top-level: `usePage().flash`, not inside `props` |
| URL state lost on navigation | Parsing `window.location` in useEffect | Derive from props — controller reads `params` and passes as prop |
| `WhenVisible` never triggers | Element not in viewport or prop name wrong | `data` must match a prop name the controller provides on partial reload |
| Component state resets on `router.get` | Missing `preserveState: true` | Add `preserveState: true` to visit options for filter/sort/tab changes |
| Scroll jumps to top after form submit | Missing `preserveScroll` | Add `preserveScroll: true` to the visit or form options |

## Related Skills
- **Flash config** → `inertia-go-controllers` (session flash provider)
- **Flash toast UI** → `shadcn-inertia` (Sonner + useFlash)
- **Shared props typing** → `inertia-go-typescript` (InertiaConfig)
- **Deferred server-side** → `inertia-go-controllers` (inertia.Defer)
- **URL-driven dialogs** → `shadcn-inertia` (Dialog component)

## References

**MANDATORY — READ ENTIRE FILE** when implementing event callbacks (`onBefore`,
`onStart`, `onProgress`, `onFinish`, `onCancel`), client-side flash, or scroll
management:
[`references/navigation.md`](references/navigation.md) (~200 lines) — full callback
API, `router.flash()`, scroll regions, and history encryption.

**MANDATORY — READ ENTIRE FILE** when implementing nested layouts, conditional
layouts, or layout-level data sharing:
[`references/layouts.md`](references/layouts.md) (~180 lines) — nested layout patterns,
layout props, and default layout configuration.

**Do NOT load** references for basic `<Link>`, `router.visit`, or single-level
layout usage — the examples above are sufficient.
