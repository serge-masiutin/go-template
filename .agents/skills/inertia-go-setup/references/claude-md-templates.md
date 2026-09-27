# Agent Instruction Stack Template

Adapt this content to the existing instruction file; retain local rules and document only installed, verified capabilities.

```markdown
## Go + Inertia Stack

- Read go.mod and package-lock.json before using version-specific APIs.
- Server: net/http, Gonertia v3, pgx, SCS; frontend: React, TypeScript, Vite, Tailwind.
- Use inertia-go-architecture first for page features, then the focused controllers/forms/pages/testing skills.
- Parse strict JSON input; authorize on the server; use explicit public DTOs.
- Mutations include X-CSRF-Token from shared props. Validation uses session errors and 303 redirects.
- Use request-context SetProp/SetProps for shared values; never mutate global user props.
- Keep frontend types synchronized with Go response contracts; no automatic generator is installed.
- Use restored sb-* skills for Storybook and tailwind-best-practices for tokens.
- Follow docs/development.md and docs/testing.md for verified commands.
```

Optional capabilities such as uploads, shadcn, realtime, queues, mail and generated contracts belong here only after implementation. Do not add references to absent tools or dependency names copied from the Rails source.
