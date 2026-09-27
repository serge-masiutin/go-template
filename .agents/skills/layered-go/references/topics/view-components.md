# View Components → React Components

## Summary

Components encapsulate reusable UI structure, presentation logic and interaction. Preserve the source's component boundaries, composition and previews; use React and Storybook already present in the starter.

## Layer Placement

Components belong to presentation. They receive explicit typed props and emit user intentions. Authoritative business state and authorization stay on the server; Inertia owns page navigation and prop delivery.

## When to Use

Extract repeated markup with stable semantics, components with meaningful variants, or independently testable interactions. A component may have one consumer when it isolates a substantial responsibility.

## When NOT to Use

Do not create a framework around one line of markup or move domain rules into a component because the result is visible in the UI.

## Key Principles

Use fixed variants, semantic tokens and accessible HTML. Prefer a closed prop contract to arbitrary styling props. Supply all server data through props; make loading, empty, forbidden and error states explicit.

## Implementation

### Basic Component

```tsx
type StatusProps = { state: "draft" | "published" };
export function Status({ state }: StatusProps) {
  return <span>{state === "published" ? "Published" : "Draft"}</span>;
}
```

### Component with Template

JSX is the template. Keep rendering declarative; do not concatenate HTML in Go or use `dangerouslySetInnerHTML` for ordinary user text. React escapes text values.

### Usage in Pages

Import the real component into the Inertia page and into its story. Do not maintain a second approximation for Storybook. Use `<Link>` for internal navigation and `<Form>` for server mutations with this project's CSRF header.

### Slots and Polymorphism

Use `children` or named `ReactNode` slots where composition is meaningful. Prefer explicit semantic variants before adding an unconstrained `as` prop. Preserve keyboard, label and focus behavior when changing the underlying element.

### Component Composition

Pages assemble feature components; shared primitives own consistent controls. Avoid a component that fetches unrelated data, authorizes users and renders every page state through dozens of flags.

## Testing Components

Test observable behavior with Testing Library and browser tests. Server contract tests own permissions and data shape; component tests own controls, messages, keyboard interactions and disabled/submitting states.

### Preview Components

Use the restored `sb-stories`, `sb-inventory` and `sb-health` skills. Stories render real components with deterministic fixtures and representative states. Storybook runs separately from production pages.

## Common Component Patterns

### Collection Component

Use stable domain IDs as keys. Keep pagination and filtering in the URL/server contract; a loaded empty list differs from an unrequested deferred prop.

### Inline Component

Keep simple one-off markup in its page until a stable abstraction appears. Repeated styles alone can be a token or variant rather than a new component.

### Local Interaction

React state handles open/closed dialogs, local drafts and focus. Do not add Stimulus to reproduce the Rails chapter. Clean up effects and subscriptions; browser-only work belongs in effects or event handlers.

## When to Extract

### From Helpers

HTML-building functions with multiple nested branches become JSX components with typed props and stories.

### From Presenters

Pure display-value mapping can remain a function. Markup and interaction belong to components; server presenters must not emit arbitrary HTML or CSS classes for every view.

## Anti-Patterns

### Data Fetching in Components

Do not duplicate Inertia page data with mount-time fetches. Use server props and supported deferred/partial reloads. Independent widgets with a separate API need an explicit reason and contract.

### Business Logic in Components

A hidden Delete button does not authorize deletion. The server checks every action independently.

### God Components

Split by cohesive behavior, not arbitrary line count. Avoid hundreds of props and unrelated flags; prefer composition with a clear owner of state.

## File Organization

The starter uses `web/src/components`, `web/src/pages`, colocated stories/tests and shared CSS tokens. Follow the local layout before inventing a second component system.
