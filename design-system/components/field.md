---
sourcesHash: c03b200dfe201f75dcd60c960ac43e598c0bf343e54c6163d6327abc0029e565
id: field
description: Edit one named text value through a labelled input or textarea, with associated instructions and validation feedback while preserving native form semantics.
status: discoverable
sources:
  - web/src/components/Field.tsx
  - web/src/components/Notice.tsx
  - web/src/styles.css
tests:
  - web/src/components/Field.test.tsx
  - tests/browser/workspace.spec.ts
examples:
  - web/src/components/Field.stories.tsx
---
# Field

## When to use

Use when both conditions hold:

- The end user enters or reviews one named text value in a form.
- That value can use a native text, email or password input, or a textarea for multiline content.

A note body and an account email both fit; the consumer chooses the control from the value's meaning rather than the screen's appearance.

## When not to use

- The value is a selection, file, date or rich-text document. Those interactions need their own native control or component contract.
- The content is only an output without an editing or form task. Use normal text instead of a read-only field purely for appearance.

## Public API

Import `Field` from `web/src/components/Field.tsx`.

```tsx
<Field id="email" name="email" label="Email" type="email"
  autoComplete="username" required error={errors.email} />
<Field multiline id="body" name="body" label="New note"
  rows={4} maxLength={2000} required error={errors.body} />
```

- `id`, `name`, and `label` are required. Use a document-unique `id`, the server's form key as `name`, and a visible name describing the value. Reserve the derived IDs `<id>-hint` and `<id>-error` for this component.
- `multiline` defaults to false. Set it to true when line breaks are part of the value, such as notes or questions; otherwise use a single-line input.
- `type` defaults to `text`. For single-line controls choose `email` for an email address and `password` for a credential. It is not part of the multiline API.
- Optional `hint` explains how to supply the value. Optional `error` contains the current field's validation message from the form. Omit it when the field has no error; do not copy unrelated service failures here.
- Native attributes for the selected control are forwarded, including `required`, `autoComplete`, `maxLength`, `rows`, `disabled`, `readOnly`, `defaultValue`/`value` and native events. Derive validation limits and autocomplete from the actual form contract. Set rows according to the expected editing task; it controls initial height, not a domain limit. For controlled values supply `onChange` with `value` under React's native input contract.
- `aria-describedby` can name additional existing description elements. The component appends its hint and error IDs. It owns `aria-invalid`, `aria-label`/`aria-labelledby` are excluded so they cannot replace the visible label, and it provides no `className` or `style` escape hatch.

The component owns internal label/control/message spacing and token-based input/focus styles. Its parent owns field order, external spacing, grouping, submission, CSRF, validation and persistence. Error appearance comes from [Notice](notice.md).

## Behaviour and states

Controls retain native editing, password masking, autocomplete, validation attributes, read-only behaviour, and disabled form exclusion. The component performs no validation or request. An error adds `aria-invalid=true` and the associated alert; clearing it removes both the invalid state and the obsolete description reference. A hint can remain alongside an error.

## Accessibility

The label targets the control; hints and errors are programmatically associated. Errors are announced through Notice's alert semantics. The control keeps a visible keyboard focus outline. Consumers must keep IDs unique, avoid duplicate copies of the same alert outside the field, and describe what can be corrected in error text. Required state is native; task instructions should make required/optional expectations clear.
