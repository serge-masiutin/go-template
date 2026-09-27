---
name: inertia-go-architecture
description: Server-driven Go + Gonertia + React architecture. Load first for pages, forms, CRUD, navigation and state ownership; routes to detailed skills and distinguishes Inertia pages from independent APIs.
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


# Inertia Go Architecture

Server-driven architecture for Go + Inertia.js + React when building pages,
forms, navigation, or data refresh. Inertia is NOT a traditional SPA — the
server owns routing, data, and auth. React owns rendering and local interaction state.

## The Core Mental Model

The server is the source of truth. React receives data as props and renders UI.
Use server routes and page props by default. Add an independent API only for a consumer or interaction that needs one.

**Before building any feature, ask:**
- **Where does the data come from?** → If server: controller prop. If user interaction: `useState`.
- **Who owns this state?** → If it's in the URL or DB: server owns it (use props). If it's ephemeral UI: React owns it.
- **Am I reaching for a React/SPA pattern?** → Check the decision matrix below first — Inertia likely has a server-driven equivalent.

## Decision Matrix

| Need | Solution | NOT This |
|------|----------|----------|
| Page data from server | Controller props | useEffect + fetch |
| Global data (auth, config) | request-scoped `inertia.SetProps` + `usePage()` | React Context / Redux |
| Flash messages / toasts | `inertia.SetFlash` + session provider + `usePage().flash` | inertia_share / React state |
| Form submission | `<Form>` component | fetch/axios + useState |
| Navigate between pages | `<Link>` / `router.visit` | react-router / window.location |
| Refresh specific data | `router.reload({ only: [...] })` | React Query / SWR |
| Expensive server data | `inertia.Defer` | useEffect + loading state |
| Infinite scroll | `inertia.Scroll` with explicit metadata + `<InfiniteScroll>` | Client-side pagination |
| Stable reference data | `inertia.Once` (verify adapter semantics) | Cache in React state |
| Real-time updates (core) | authorized SSE/WebSocket + `router.reload` when needed | Polling with setInterval |
| Simple polling (MVP/prototyping) | `usePoll` (auto-throttles in background tabs) | setInterval + router.reload |
| URL-driven UI state (dialogs, tabs) | Handler parses URL values → prop, `router.get` to update | useEffect + window.location |
| Ephemeral UI state | `useState` / `useReducer` | Server props |
| External API calls | Dedicated API endpoint | Mixing with Inertia props |

## Rules (by impact)

| Impact | Rule | Reason |
| --- | --- | --- |
| Critical | Server enforces every permission | UI capabilities are hints, including during partial/deferred requests |
| Critical | Submit forms through Inertia with the project CSRF header | Preserves redirect/error lifecycle; Gonertia does not add CSRF protection for you |
| High | Use page props for page data | Avoid a second mount-time fetch/cache lifecycle for the same resource |
| High | Use Link/router for internal visits | Preserves Inertia navigation and page/layout state |
| High | Share request data through context | Global mutable ShareProp can leak between users |
| High | Use the flash provider for one-visit messages | Shared props have different persistence semantics |
| Medium | Defer noncritical expensive data after measuring | Keep authority and form defaults available immediately |
| Medium | Choose polling or push from freshness/load needs | Polling can be appropriate in production; no realtime server is installed |
| Medium | Keep layout/state ownership explicit | Persistent layouts and local React state solve different needs |

External links/downloads can use ordinary anchors. Independent autocomplete widgets can use a bounded API. Do not turn the defaults into bans on legitimate React context, local state or separate consumers.

## Skill Map

Common workflows span multiple skills — load all listed for complete coverage:

| Workflow | Load these skills |
|----------|-------------------|
| New page with props | `inertia-go-controllers` + `inertia-go-pages` + `inertia-go-typescript` |
| Form with validation | `inertia-go-forms` + `inertia-go-controllers` |
| shadcn form inputs | `inertia-go-forms` + `shadcn-inertia` |
| Flash toasts | `inertia-go-controllers` + `inertia-go-pages` + `shadcn-inertia` |
| Deferred/lazy data | `inertia-go-controllers` + `inertia-go-pages` |
| URL-driven dialog/tabs | `inertia-go-controllers` + `inertia-go-pages` |
| Explicit DTO serialization | `go-serialization` + `inertia-go-typescript` |
| Testing controllers | `inertia-go-testing` + `inertia-go-controllers` |

## References

**MANDATORY — READ ENTIRE FILE** before building a new Inertia page or feature:
[`references/AGENTS.md`](references/AGENTS.md)  — full-stack examples for
each pattern in the decision matrix above.

**MANDATORY — READ ENTIRE FILE** when unsure which Inertia pattern to use:
[`references/decision-trees.md`](references/decision-trees.md)  — flowcharts
for choosing between prop types, navigation methods, and data strategies.

**Do NOT load** references for quick questions about a single pattern already
covered in the decision matrix above.

## When You DO Need a Separate API

Not everything belongs in Inertia's request cycle. Use a traditional API endpoint when:

| Signal | Why | Example |
|--------|-----|---------|
| Non-browser consumer | Inertia's JSON envelope (component, props, url, version) is designed for the frontend adapter — other consumers can't use it | Mobile API, CLI tools, payment webhooks |
| Large-dataset search | Dataset is too big to load as a prop; each input needs per-keystroke server filtering. Use raw fetch for the search, let Inertia handle post-selection side effects via props. | City/address autocomplete, postal code lookup |
| Binary/streaming response | Inertia can only deliver JSON props. Use a separate route with a standard download response. | PDF/CSV export, file downloads |
