# React Components and Storybook

This replaces ViewComponent classes, ERB templates and Lookbook previews with the original React/Storybook workflow.

## Setup and Basic Usage

The starter already includes React, Vite, Tailwind and Storybook. Use `web/src/components` and colocated `*.stories.tsx`. Do not add ViewComponent, ERB or Stimulus.

## Props, Slots and Composition

Define a closed TypeScript prop type. Use `children` or named `ReactNode` slots for composition. Prefer fixed semantic variants to arbitrary class strings and unconstrained element polymorphism.

## Collections and Inline Components

Render lists with stable IDs and explicit empty/loading states. Keep one-off markup local until a coherent reusable responsibility emerges. Avoid recreating the entire domain object as a component API.

## Local Interaction

React owns UI state and effect cleanup. Inertia owns server navigation, page props and form visits. The server owns permissions and durable state.

## Previews

Use the restored `sb-*` skills to inventory actual components, build representative stories, audit tokens/accessibility and graduate experiments. Stories import production components; preview-only experiments stay separate until accepted.

## Testing

Test visible interaction, keyboard behavior and accessible labels. Browser tests cover server integration; Go tests cover authorization and prop contracts. See [view components](../topics/view-components.md) for extraction and anti-patterns.
