---
sourcesHash: 39cdc233fe7df6f42dcf67ba69ecebda2653e8af873335fd0ecbc82ef43d74d0
id: button
description: Invoke a local action or submit a form through a native button, with primary or secondary emphasis and a native disabled state.
status: discoverable
sources:
  - web/src/components/Button.tsx
  - web/src/styles.css
tests:
  - web/src/components/Button.test.tsx
examples:
  - web/src/components/Button.stories.tsx
---
# Button

## When to use

Use when both conditions hold:

- The end user invokes an action in the current context, such as adding a note, starting work, or signing out.
- A native button can express the interaction; the consumer owns the handler or form submission.

## When not to use

- The action is navigation to a URL. Use an Inertia `Link` or a native anchor so browser navigation semantics remain available.
- The interaction selects a persistent on/off value. Use a control with selection semantics rather than presenting an action as a toggle.

## Public API

Import `Button` from `web/src/components/Button.tsx`.

```tsx
<Button type="submit" disabled={processing}>Add note</Button>
<Button variant="secondary" onClick={removeNote}>Delete</Button>
```

- `children` supplies a visible action label. Describe the action; use text to distinguish working states such as “Signing in…”.
- `variant` defaults to `primary` for the principal action in a task. Choose `secondary` for supporting actions, navigation-adjacent actions such as sign out, or repeated per-record controls. Secondary does not imply destructive confirmation.
- `type` defaults to `button`; choose `submit` only when activating the containing form. Native `form` may name a form elsewhere in the document.
- `disabled` is controlled by the consumer's current availability or processing state. The component never infers availability or starts work.
- Native button attributes and event handlers are forwarded. Consumers own accessible names, event consequences, form association and any `aria-*` state needed by their action. Do not pass component styling through `style`; `className` is excluded from the typed API. Place layout classes on the containing composition.

The component owns padding, radius, typography and focus treatment from [styles.css](../../web/src/styles.css); its consumer owns external spacing and positioning.

## Behaviour and states

Primary, secondary, disabled and consumer-provided working labels retain native button semantics. Disabled buttons do not invoke click handlers and are skipped in the tab order. No spinner, asynchronous state, confirmation, or navigation is built in.

## Accessibility

The native button supports keyboard activation. The component supplies a visible focus outline and disabled presentation; the consumer supplies a meaningful name and decides which state prevents activation. Do not use colour as the only explanation of an unavailable action.
