---
id: server-form
description: Submit a server-owned task through an Inertia form, keeping field validation beside the affected value and the principal action available according to processing state.
status: discoverable
sources: []
tests:
  - tests/browser/workspace.spec.ts
examples:
  - web/src/pages/Login.tsx
  - web/src/pages/Home.tsx
  - web/src/pages/Tools.tsx
---
# Server-owned form

## When to use

Use when both conditions hold:

- The user submits one task to a server-owned Inertia route.
- The task uses native named controls or a single action with no editable values.

## When not to use

- The controls only adjust local display state or navigate to another URL. Use local controls or links without a server mutation form.
- The task requires independent API response handling. This recipe depends on Inertia's form/error lifecycle, not an arbitrary JSON client.

## Structure

1. A nearby heading or section name establishes the task.
2. An Inertia `Form` owns the actual action, method and `X-CSRF-Token` header.
3. Optional [Field](../components/field.md) instances collect named values in reading order; each receives its own server error. An action-level failure uses [Notice](../components/notice.md) before the action.
4. A primary [Button](../components/button.md) submits the form and is disabled while processing; any domain-specific availability condition is additionally owned by the application.

```tsx
<Form action="/notes" method="post" headers={{ "X-CSRF-Token": csrfToken }} resetOnSuccess>
  {({ errors, processing }) => <>
    <Field multiline id="body" name="body" label="New note" required
      rows={4} maxLength={2000} error={errors.body} />
    <Button type="submit" disabled={processing}>Add note</Button>
  </>}
</Form>
```

## Composition

Required: use the real route and form keys; keep validation errors tied to their corresponding controls; retain visible labels and unique control IDs. Group each label, value and message together. The containing page owns vertical spacing between fields and the action, while each Field owns its internal spacing. Controls expand inside their available parent width; reading order remains label, control, hint/error, action.

The application owns processing, CSRF, server validation, permissions and post-submit state. Use `resetOnSuccess` only when a successful task consumes the input (the existing note and assistant forms); do not clear credentials or errors by adding a component-side handler. Use the real processing state for disabling, not an invented delay.

## Verification

Check labels resolve to their native controls, server errors set invalid state and are included in the accessible description, and successful submission clears the stale error relationship. Native names and request headers must remain unchanged by presentation edits. Run `npm run test` and the database-backed `bin/test-browser` per [testing](../../docs/testing.md).

Reject a composition that renders a failure under the form but leaves its affected field without the associated error reference, or that uses a default `type="button"` for its submit action. A correct visual arrangement alone does not establish a working server contract.
