---
name: inertia-go-typescript
description: >-
  TypeScript type safety for Inertia Go (React, Vue, Svelte): shared props, flash, and errors via InertiaConfig
  module augmentation in globals.d.ts. Use when setting up TypeScript types, configuring shared props typing, fixing
  TS2344 or TS2339 errors in Inertia components, or adding new shared data.
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


# Inertia Go TypeScript Setup

Type-safe shared props, flash, and errors using InertiaConfig module augmentation.
The starter uses React. Its included TypeScript module augments InertiaConfig once.

**Before adding TypeScript types, ask:**
- **Shared props (auth, flash)?** → Update `SharedProps` (and define `FlashData` if the feature needs flash) in `web/src/types.ts` — InertiaConfig in the included augmentation module propagates them globally via `usePage()`
- **Page-specific props?** → `type Props = { ... }` in the page file only — never include shared props here

## InertiaConfig Module Augmentation

Define shared props type ONCE globally — never in individual page components.

InertiaConfig property names are EXACT — do not rename them:
- `sharedPageProps` (NOT sharedProps)
- `flashDataType` (NOT flashProps, NOT flashData)
- `errorValueType` (NOT errorBag, NOT errorType)

```typescript
// web/src/types.ts — a module included by tsconfig.json
export type SharedProps = { csrfToken: string };
declare module '@inertiajs/core' {
  interface InertiaConfig {
    sharedPageProps: SharedProps;
  }
}
```

If adding flash, augment `flashDataType` to the exact payload. If changing error values, augment `errorValueType` to that actual server contract; the starter uses strings. Do not declare auth or flash fields that the server does not provide.

## BAD vs GOOD Patterns

```tsx
// BAD — passing shared props as generics:
// usePage<{ users: User[], auth: AuthData, flash: FlashData }>()

// BAD — extending a SharedProps interface into page props:
// interface Props extends SharedData { users: User[] }

// BAD — declaring PageProps interface:
// interface PageProps { auth: AuthData; flash: FlashData }

// BAD — destructuring auth directly from usePage() (TS2339: 'auth' does not exist on Page):
// const { auth } = usePage()
// usePage() returns a Page object with { props, flash, component, url, ... }
// auth lives inside props, not on the Page itself

// BAD — declaring conflicting InertiaConfig property types in several modules

// GOOD — props from usePage().props, flash from usePage().flash:
const { props, flash } = usePage()
// props.csrfToken is typed through InertiaConfig
// Add an exact flashDataType declaration before using application flash keys.
```

**Important:** the included augmentation module configures InertiaConfig ONCE. When adding a new shared
prop, update its canonical SharedProps declaration in `web/src/types.ts`:

```typescript
// BEFORE — web/src/types.ts
export type SharedProps = {
  csrfToken: string
}

// AFTER — add the new key here, NOT in globals.d.ts
export type SharedProps = {
  csrfToken: string
  notifications: { unread_count: number }
}
```

InertiaConfig in the included augmentation module references `SharedProps` by name — it picks up the
change automatically. Keep one canonical declaration. Compatible TypeScript declaration merging is legal, but conflicting property types are not.

## Page-Specific Props

Page components type ONLY their own props. Shared props (like csrfToken) and flash come from
InertiaConfig automatically.

### `type` vs `interface` for page props (React-specific)

This constraint applies to **React only**. Vue's `defineProps<T>()` and Svelte's
`$props()` do not use `usePage<T>()` generics, so `interface` works fine there.

`usePage<T>()` requires `T` to have an index signature. `type` aliases have one
implicitly; `interface` declarations do not. Using `interface` with `usePage`
causes TS2344 at compile time.

| Pattern | Works with `usePage<T>()`? | Notes |
|---------|---------------------------|-------|
| `type Props = { users: User[] }` | Yes | Preferred — just works |
| `interface Props { users: User[] }` | **No** — TS2344 | Missing index signature |
| `usePage<Required<Props>>()` | Yes | Wraps interface to add index signature |

```tsx
// React
type Props = {
  users: User[]         // page-specific only
  // Shared csrfToken comes from InertiaConfig globally
}

export default function Index({ users }: Props) {
  // Access shared props separately:
  const { props, flash } = usePage()
  // props.csrfToken is typed via InertiaConfig
  // Flash keys require their own exact flashDataType declaration.
  return <UserList users={users} />
}
```

### Accessing shared props in Vue and Svelte

The template is React-only. The original upstream skill supports other adapters, but those packages and build pipelines are not installed here. Do not mix framework examples into this starter.

## Common TypeScript Errors

| Error | Cause | Fix |
|-------|-------|-----|
| **TS2344** on `usePage<Props>()` | `interface` lacks index signature | Use `type Props = { ... }` instead of `interface`, or wrap: `usePage<Required<Props>>()` |
| **TS2339** `'auth' does not exist on type Page` | Destructuring `auth` from `usePage()` directly | `usePage()` returns `{ props, flash, ... }` — use `usePage().props.auth`, not `usePage().auth` |
| **TS2339** `'flash' does not exist on type` | Accessing `usePage().props.flash` | Flash is top-level: `usePage().flash`, NOT `usePage().props.flash` |
| Shared props untyped | Missing InertiaConfig | Add the included augmentation module with module augmentation (see above) |
| InertiaConfig not taking effect | Augmentation module not included | Include the module in tsconfig; augmentation works in an imported `.ts` module or a `.d.ts` file |
| Types correct but IDE shows errors | the included augmentation module not included | Verify `tsconfig.json` includes the types directory in `include` array |

## Go Contract Synchronization

The Go starter uses explicit response DTOs and TypeScript declarations. There is no Typelizer or implicit generator. Update both contracts atomically and test JSON keys, decimal-string IDs, nullability and arrays. If adopting generation later, commit the generator/configuration and fail CI on drift.

## Related Skills
- **Shared props setup** → `inertia-go-controllers` (request-context SetProps)
- **Flash config** → `inertia-go-controllers` (session flash provider)
- **JSON contracts** → `go-serialization` (explicit DTO mapping)
- **Page component props** → `inertia-go-pages` (type Props pattern)
