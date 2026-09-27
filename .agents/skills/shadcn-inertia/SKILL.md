---
name: shadcn-inertia
description: >-
  shadcn/ui component integration for Inertia Go React (NOT Next.js): forms, dialogs, tables, toasts, dark mode,
  command palette, and more. Use when building UI with shadcn/ui components in an Inertia app or adapting shadcn
  examples from Next.js. NEVER react-hook-form/zod — wire shadcn inputs to Inertia Form via name attribute. Flash
  toasts require a session-backed Gonertia flash provider and matching TypeScript types.
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


# shadcn/ui for Inertia Go

shadcn/ui patterns adapted for Inertia.js + Go + React. NOT Next.js.

**Before using a shadcn example, ask:**
- **Does it use `react-hook-form` + `zod`?** → Replace with Inertia `<Form headers={{ "X-CSRF-Token": csrfToken }}>` + `name` attributes. Inertia handles errors, redirects and processing state; add the project CSRF header — react-hook-form would fight all of this.
- **Does it use `'use client'`?** → Remove it. Inertia has no RSC — all components are client components.
- **Does it use `next/link`, `next/head`, `useRouter()`?** → Replace with Inertia `<Link>`, `<Head>`, `router`.

## Key Differences from Next.js Defaults

| shadcn default (Next.js) | Inertia equivalent |
|---|---|
| `'use client'` directive | Remove — not needed (no RSC) |
| `react-hook-form` + `zod` | Inertia `<Form headers={{ "X-CSRF-Token": csrfToken }}>` component |
| `FormField`, `FormItem`, `FormMessage` | Plain `<Input name="...">` + `errors.field` |
| `next-themes` | CSS class strategy + `@custom-variant` |
| `useRouter()` (Next) | `router` from `@inertiajs/react` |
| `next/link` | `<Link>` from `@inertiajs/react` |
| `next/head` | `<Head>` from `@inertiajs/react` |

**NEVER use shadcn's `FormField`, `FormItem`, `FormLabel`, `FormMessage` components** —
they depend on react-hook-form's `useFormContext` internally and will crash without it.
Use plain shadcn `Input`/`Label`/`Select` with `name` attributes inside Inertia `<Form headers={{ "X-CSRF-Token": csrfToken }}>`,
and render errors from the render function's `errors` object (see examples below).

## Setup

The starter does not preinstall shadcn. When requested, inspect the current CLI and initialize it against the existing React/Vite/Tailwind stack. Keep `@/` aligned in both `vite.config.ts` and `tsconfig.json`. Pin generated dependencies.

## shadcn Inputs in Inertia `<Form headers={{ "X-CSRF-Token": csrfToken }}>`

Use plain shadcn `Input`/`Label`/`Button` with `name` attributes inside Inertia `<Form headers={{ "X-CSRF-Token": csrfToken }}>`.
See `inertia-go-forms` skill for full `<Form headers={{ "X-CSRF-Token": csrfToken }}>` API — this section covers shadcn-specific adaptation only.

**The key pattern:** Replace shadcn's `FormField`/`FormItem`/`FormMessage` with plain
components + manual error display:

```tsx
// shadcn error display pattern (replaces FormMessage):
<Label htmlFor="name">Name</Label>
<Input id="name" name="name" />
{errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
```

**`<Select>` requires `name` prop** for Inertia `<Form headers={{ "X-CSRF-Token": csrfToken }}>` integration — shadcn examples
omit it because react-hook-form manages values differently:

```tsx
<Select name="role" defaultValue="member">
  <SelectTrigger><SelectValue placeholder="Select role" /></SelectTrigger>
  <SelectContent>
    <SelectItem value="admin">Admin</SelectItem>
    <SelectItem value="member">Member</SelectItem>
  </SelectContent>
</Select>
```

## Dialog with Inertia Navigation

```tsx
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { router } from '@inertiajs/react'

function UserDialog({ open, user }: { open: boolean; user: User }) {
  return (
    <Dialog
      open={open}
      onOpenChange={(isOpen) => {
        if (!isOpen) {
          router.replaceProp('show_dialog', false)
        }
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{user.name}</DialogTitle>
        </DialogHeader>
        {/* content */}
      </DialogContent>
    </Dialog>
  )
}
```

## Table with Server-Side Sorting

shadcn `<Table>` renders normally. The Inertia-specific part is sorting via `router.get`:

```tsx
const handleSort = (column: string) => {
  router.get('/users', { sort: column }, { preserveState: true })
}

<TableHead onClick={() => handleSort('name')} className="cursor-pointer">
  Name {sort === 'name' && '↑'}
</TableHead>
```

Use `<Link>` (not `<a>`) for row links to preserve SPA navigation.

## Toast with Flash Messages

Flash provider configuration is in `inertia-go-controllers`. Flash access
(`usePage().flash`) is in `inertia-go-pages`. This section covers **toast UI wiring only**.

**MANDATORY — READ ENTIRE FILE** when implementing flash-based toasts with Sonner:
[`references/flash-toast.md`](references/flash-toast.md) (~80 lines) — full `useFlash`
hook and Sonner toast provider. **Do NOT load** if only reading flash values without toast UI.

Key contract: the Go flash payload MUST match your `FlashData`
TypeScript type — do NOT use `success`/`error` unless you also update both.

## Dark Mode (No next-themes)

`npx shadcn@latest init` generates CSS variables for light/dark and
`@custom-variant dark (&:is(.dark *));` in your CSS (Tailwind v4). No extra
setup needed for the variables themselves.

To avoid a flash of the wrong theme, run theme selection before loading styles and
React. The starter's production CSP permits same-origin external scripts; it rejects
unhashed inline scripts. A feature can add `web/public/theme.js` (copied by Vite to
`web/build/theme.js`) and reference it in `internal/httpapp/root.html` before assets:

```html
<script src="/build/theme.js"></script>
```

Implement that script and the React appearance hook together. Both must use the same
`light` / `dark` / `system` values, storage key, and system-media-query semantics.
Use `.dark` on `<html>`; persist deliberate user choices, remove the media listener on
cleanup, and treat blocked browser storage as an explicit optional-preference case.
Test initial load and navigation under the production CSP. No theme hook is installed
in the base starter. Do not copy an import for an unimplemented `use-appearance` hook.

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `FormField`/`FormMessage` crash | Using shadcn form components that depend on react-hook-form | Replace with plain `Input`/`Label` + `errors.field` display |
| `Select` value not submitted | Missing `name` prop | Add `name="field"` to `<Select>` — shadcn examples omit it |
| Dialog closes unexpectedly | Missing or wrong `onOpenChange` handler | Use `onOpenChange={(open) => { if (!open) closeHandler() }}` |
| Flash of wrong theme (FOUC) | Theme selection runs after styles/React | Add the same-origin external theme script before assets and verify production CSP (see Dark Mode section) |

## Related Skills
- **Form component** → `inertia-go-forms` (`<Form headers={{ "X-CSRF-Token": csrfToken }}>` render function, useForm)
- **Flash config** → `inertia-go-controllers` (session flash provider)
- **Flash access** → `inertia-go-pages` (usePage().flash)
- **URL-driven dialogs** → `inertia-go-pages` (router.get pattern)

## References

Load [`references/components.md`](references/components.md) (~300 lines) when building
shadcn components beyond those shown above (Accordion, Sheet, Tabs, DropdownMenu,
AlertDialog with Inertia patterns).

**Do NOT load** `components.md` for basic Form, Select, Dialog, or Table usage —
the examples above are sufficient.

Examples above use `const { csrfToken } = usePage().props` inside the component. Include that binding/import where a snippet is used; never hardcode a token.
