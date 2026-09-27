---
sourcesHash: 5b68fa2266f6eff04b2dd0b4dd25686049c3bbe77fd69d68c9ccaf3195bf29d7
id: notice
description: Announce an actionable validation or operation failure in the current task, using visible danger styling and alert semantics.
status: discoverable
sources:
  - web/src/components/Notice.tsx
  - web/src/styles.css
tests:
  - web/src/components/Field.test.tsx
  - tests/browser/workspace.spec.ts
examples:
  - web/src/components/Notice.stories.tsx
---
# Notice

## When to use

Use when both conditions hold:

- The current action failed at form/task level, or the message is being composed internally by Field for a submitted value that needs correction.
- The message needs the end user's attention now and belongs beside the affected task.

## When not to use

- The message is routine guidance, success feedback, or an empty-state explanation. Use ordinary text because an alert would interrupt the task without a failure to resolve.
- The content is the current state of background work without explanatory failure guidance. Use [WorkStatus](work-status.md) to annotate the work item.
- A page is assembling its own separate error alongside a text control. Supply `error` to [Field](field.md) instead; Field already composes Notice with the necessary description relationship.

## Public API

Import `Notice` from `web/src/components/Notice.tsx`.

```tsx
<Notice>Delivery failed. Check your inbox before requesting another copy.</Notice>
```

`children` is required plain text explaining the problem and, when known, the next action. Optional `id` gives the message a document-unique target for a referring control. There are no success variants, dismissal handlers or arbitrary presentation overrides.

The component owns border, padding, typography and danger tokens from [styles.css](../../web/src/styles.css). Its parent owns placement and when the message appears.

## Behaviour and states

Renders the supplied message as an alert. It neither detects failures nor clears them. The consumer removes or replaces the message when the task state changes.

## Accessibility

`role="alert"` announces new failure content. Avoid repeating the same message in multiple simultaneous alerts. When a control references the message, its owner must supply `aria-describedby`; Field handles that composition automatically.
