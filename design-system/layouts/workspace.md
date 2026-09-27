---
sourcesHash: 357a2ba55c574cc881718ecb0364b8f243c3a8dcab5aeeeb9f2aef1ae3b20220
id: workspace-layout
description: Arrange authenticated workspace navigation and one page's main content within a shared responsive container, with role-aware administration navigation and sign out.
status: discoverable
sources:
  - web/src/components/Layout.tsx
  - web/src/styles.css
tests:
  - tests/browser/workspace.spec.ts
examples:
  - web/src/components/Layout.stories.tsx
---
# Workspace layout

## When to use

Use when both conditions hold:

- The page is part of the signed-in workspace.
- Its navigation includes the workspace home, note tools, and sign out, with administration shown for an administrator.

## When not to use

- The page is a public authentication screen. Use a standalone main region because authenticated navigation would misrepresent access.
- The composition is a section inside an existing workspace main region. Place the section as children rather than nesting another layout and duplicating landmarks.

## Public API

Import `Layout` from `web/src/components/Layout.tsx`.

```tsx
<Layout user={user}><h1>Workspace</h1><section aria-label="Notes">…</section></Layout>
```

`user` is the authenticated server DTO; its `admin` flag selects the visible administration link. `children` is the page body. The layout must be inside the real Inertia application provider with the shared `csrfToken` prop. Navigation is presentation, not authorization; the server enforces all permissions.

The layout owns Inertia links to `/`, `/tools`, and conditional `/admin`, and a POST sign-out form with the CSRF header and a secondary [Button](../components/button.md). Consumers do not supply or replace those routes through styling.

## Composition

Required:

- Layout owns the centered maximum width, horizontal gutters, header border, header spacing and main vertical padding from [Layout.tsx](../../web/src/components/Layout.tsx) and [styles.css](../../web/src/styles.css).
- Header brand precedes the main navigation; navigation wraps as available width decreases. The main region follows the header in document order.
- The consumer supplies one descriptive page heading and semantic sections inside `children`; it owns content-specific grids, field spacing and record lists.
- Do not nest another `main` landmark or another workspace layout inside the body. Consumers cannot override header spacing or hide required navigation with local CSS.

Recommended: retain the page's source order at narrow widths and let content wrap. A two-column tools composition may collapse to one column without changing the task order.

## Accessibility

The layout supplies a named Main navigation landmark and a main content landmark. Native links and the sign-out button remain keyboard operable. The page supplies its heading hierarchy, content region names and any live regions. At 320 CSS px the header wraps within its gutters without horizontal overflow.
