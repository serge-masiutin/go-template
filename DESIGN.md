# Go Template design

The starter provides a signed-in workspace, notes, background-work feedback and
administration through Go, Inertia and React. Preserve its existing light palette
and local Martian Mono. There is no external Figma specification or theme switcher.
[AGENTS.md](AGENTS.md) and [architecture](docs/architecture.md) govern implementation:
the server owns routes, access and data; React owns presentation and local interaction.

## Sources and scope

- [Components](design-system/COMPONENTS.md), [layouts](design-system/LAYOUTS.md),
  and [patterns](design-system/PATTERNS.md) link public usage contracts.
- [Token roles](design-system/tokens.md) describe the semantic decisions defined in
  [styles.css](web/src/styles.css), including the local font and light theme.
- [React components](web/src/components) provide the public bindings;
  [pages](web/src/pages) compose them and own request-specific content.
- [Storybook configuration](.storybook/main.ts) and colocated
  `web/src/components/*.stories.tsx` render real components with application CSS.

A contract documents the task a component serves, its public settings, states and
consumer obligations. `Layout` owns the page shell; patterns describe reusable
composition recipes and need no runtime wrapper. Private helpers and routed pages
are not separate shared components. shadcn is not installed; the skill for it does
not imply an available component library.

## Shared rules

- Select a public component by its contract before copying markup. Components own
  internal styling; pages own content, region order, outer spacing and request data.
- Use semantic color roles instead of repeating values. Preserve the existing
  Tailwind spacing scale; separate roles even where their current values match.
- Use a field's public API to connect its label, hint and validation error to the
  native control. The enclosing Inertia `Form` owns submission, processing, reset
  and error state. Preserve input names, autocomplete, constraints and CSRF headers.
- Use `Button` for actions and Inertia `Link` for navigation. Disabled controls must
  remain unavailable; styling must not simulate a disabled action that still fires.
- `Notice` announces a form-level error; field-specific feedback belongs beside
  the field. `WorkStatus` labels background work, not a replacement for the failure
  explanation or the server-owned work state.
- Preserve existing routing, effects and handlers. New form presentation must not
  add a second form-state library or alter request/response contracts.
- Stories stay colocated, use the actual component and cover materially different
  states. Controls and fixtures must not send requests to live services.

## Development and verification

Use [development](docs/development.md), [testing](docs/testing.md),
[Inertia forms](.agents/skills/inertia-go-forms/SKILL.md) and
[Storybook skills](.agents/skills/sb-stories/SKILL.md) for engineering.
[design-system](.agents/skills/design-system/SKILL.md) owns contracts and discovery.
All project documentation and technical content are in English.

Run from the root through `mise exec --`:

```sh
bin/design-system-check
node .agents/skills/design-system/scripts/generate-indexes.mjs .
npm test
npm run test:storybook
npm run build:storybook
bin/test-browser
bin/ci
```

CI validates metadata, links, source hashes and index freshness with
`bin/design-system-check`. Review a changed public promise before updating its
snapshot and indexes:

```sh
node .agents/skills/design-system/scripts/check-contract.mjs --update-sources-hash design-system/components/field.md
node .agents/skills/design-system/scripts/generate-indexes.mjs .
```
 Hashes do not prove runtime correctness. Browser story tests
exercise rendered states and accessibility; the application browser suite checks
real Inertia integration. Catalog compilation alone does not execute story checks.

