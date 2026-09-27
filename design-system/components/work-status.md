---
sourcesHash: c474b41da5059505b3976bee86c68235c108d07f1ab54bed623b353b2e797b8a
id: work-status
description: Identify the lifecycle state of one background email or assistant run using a concise text label attached to that work item.
status: discoverable
sources:
  - web/src/components/WorkStatus.tsx
  - web/src/styles.css
examples:
  - web/src/components/WorkStatus.stories.tsx
---
# Work status

## When to use

Use when both conditions hold:

- The end user is inspecting one previously requested background job.
- The server's state is queued, sending, sent, running, completed, or failed.

The label annotates the work item; a surrounding list or region identifies what work it describes.

## When not to use

- The content describes form submission in progress before a background job exists. The containing form controls its submit label and availability.
- The requested content is failure details or recovery advice rather than a state annotation. Use [Notice](notice.md) for that content, optionally alongside the status label.
- The text is a standalone heading or instruction rather than an attribute of a work item. Use ordinary semantic content.

## Public API

Import `WorkStatus` from `web/src/components/WorkStatus.tsx`.

```tsx
<li><WorkStatus state={email.state} /></li>
```

The required `state` comes directly from the typed server DTO. Values map to Queued, Sending, Sent, Working, Completed, and Needs attention respectively. Do not infer completion from elapsed time or translate an unknown state into a successful one. There is no default state or styling override.

The component owns its text size, weight and muted token from [styles.css](../../web/src/styles.css). The consumer owns layout and association with the work item.

## Behaviour and states

This is a static label for the supplied state. Polling, transition timing, retries and error details stay with the application. The label does not change permissions or disable actions.

## Accessibility

The state is readable text rather than colour alone. For asynchronously updated lists the containing composition provides a polite live region, as the tools page does. WorkStatus does not create a competing alert for every item.
