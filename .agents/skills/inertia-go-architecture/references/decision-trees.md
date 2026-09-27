# Decision Trees

Quick flowcharts for common Inertia.js + Go decisions.

## "I need data in my component"

```
Is this data specific to this page?
├── YES → Controller prop (engine.Render with explicit Props)
│   └── Is it expensive to compute?
│       ├── Measured slow, noncritical data → inertia.Defer(loader)
│       ├── Needed immediately → Regular prop
│       └── Other data → Measure the tradeoff, including the extra request
├── NO, it's needed on every page → request-context SetProps
└── NO, it's from an external API → Dedicated API endpoint
    (Stripe status, third-party webhook, etc.)
```

## "I need to update data"

```
Is this a form submission (create/update/delete)?
├── YES → <Form> component + controller action + redirect
│   └── Need programmatic control? → useForm hook instead
└── NO
    ├── Need to refresh page data? → router.reload({ only: [...] })
    ├── Need real-time updates?
    │   ├── Core feature (chat, live feed)? → authorized SSE/WebSocket + router.reload
    │   └── Periodic freshness sufficient? → usePoll(interval, { only: [...] })
    ├── Need optimistic UI? → useState for optimistic + router.post with onError rollback
    └── Need search/filter? → router.visit with query params (preserveState)
```

## "I need state in my component"

```
Where does the data come from?
├── Server → It's a prop, not state
├── User interaction (modal open, dropdown) → useState
├── Form data → <Form> component (or useForm for complex cases)
├── Shared across pages (auth) → usePage().props (from request-context SetProps)
└── Multiple components need it → Lift to closest common parent as prop
    └── Still unwieldy? → Consider React Context (rare in Inertia apps)
```

## "Should I prefetch / poll / defer / use authorized SSE/WebSocket?"

```
PREFETCH — preload page data before navigation:
├── Frequently visited page (dashboard, main nav)? → YES, prefetch="mount"
├── Likely next click (nav links)? → YES, prefetch (hover, default)
├── Data changes constantly per user? → NO — cache will be stale immediately
├── Page requires POST data to load? → NO — prefetch only works with GET
└── Multiple pages share data? → Use cacheTags for coordinated invalidation

POLL — auto-refresh data on an interval:
├── Dashboard counters, queue status, leaderboard? → YES, usePoll with only: [...]
├── Query is expensive? → Reduce query cost and refresh frequency; push alone does not remove it
├── Updates are rare (<1/hour)? → NO — manual refresh or authorized SSE/WebSocket
├── Need real-time (<1s latency)? → NO — use authorized SSE/WebSocket/WebSockets
└── Need user control? → { autoStart: false } + start/stop

DEFER — load expensive data after initial render:
├── Measured slow and noncritical? → Consider defer
├── Cheap and needed immediately? → Regular prop
├── Data critical for initial render (form defaults, auth)? → NO — regular prop
└── Other data? → Measure the extra request versus time to useful content

ACTIONCABLE — server pushes updates to client:
├── Core real-time feature (chat, live feed, collaboration)? → YES
├── Updates must arrive <1s after change? → YES
├── Multiple users see the same resource? → YES — broadcast on change
├── Only current user's data, low frequency? → usePoll is simpler
└── Pattern: authorized SSE/WebSocket receives event → router.reload({ only: [...] })
```

## "I need to navigate"

```
Is this a link the user clicks?
├── YES → <Link href={...}> (with prefetch for common destinations)
├── NO, programmatic after action → router.visit / router.get
├── External URL from server? → engine.Location (CRITICAL — not redirect_to)
├── External URL from client? → window.location.href
└── Need to update URL params? → router.visit with preserveState
```

## "I need to show a notification"

```
Is it a one-time message (success, error)?
├── YES → SetFlash + session provider + usePage().flash
│   └── Need custom keys beyond notice/alert? → explicit Flash DTO and provider
├── Need it to persist across navigations? → request-context SetProps (shared prop)
└── Client-side only (no server)? → router.flash('key', 'value')
```

## "I need to redirect after a mutation"

```
Is the destination inside the Inertia app?
├── YES → engine.Redirect(w, r, path, 303)
│   └── With flash? → SetFlash context + Redirect 303
└── NO, external URL (Stripe, OAuth, etc.)
    └── engine.Location url (returns 409 + X-Inertia-Location header)
        NEVER: ordinary redirect to external_url (breaks Inertia)
```
