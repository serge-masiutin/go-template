---
name: inertia-go-setup
description: Detect and align the Go, Gonertia, React, Vite, PostgreSQL and Storybook stack when bootstrapping or changing this starter. Preserve existing configuration and add dependencies only for a concrete missing capability.
---

## Go project contract

This is an adaptation of the original EM Inertia skill, preserving the React examples
and decision workflow. Read `docs/architecture.md` and the actual Go handlers first.
Use Gonertia v3 APIs from `go.mod`; server examples below replace the Rails adapter.
All mutations carry `X-CSRF-Token` from `usePage().props.csrfToken`. The starter accepts
strict JSON; uploads need a separate bounded multipart handler before using upload examples.
Validation errors cross a 303 redirect through the session flash provider. They are not
bare 422 JSON responses. Shared user data is request-scoped, never global `ShareProp`.
Examples describe features to implement, not routes or libraries already installed.


# Inertia Go Project Setup

Preserve the original detect → choose missing capability → configure → document → verify workflow. Read manifests and existing configuration before changing anything.

## Step 1: Detect Current Stack

Read `go.mod`, `go.sum`, `package.json`, `package-lock.json`, `mise.toml`, Vite/TypeScript configs, routes, session middleware and CI.

| Evidence | Meaning |
| --- | --- |
| `romsar/gonertia/v3` | Server adapter; inspect pinned API |
| `@inertiajs/react` | React client adapter |
| GORM generics + pgx + Goose | PostgreSQL persistence and explicit SQL migrations |
| River + `cmd/worker` | Transactional mail/AI jobs and reconciliation |
| go-mail + Mailpit | SMTP adapter and local delivery inspection |
| Genkit + `internal/assistant` | Typed flow/tools, Gemini/OpenAI adapters |
| `tools/go.mod` + Air | Isolated development tools and Go reload |
| SCS + pgxstore | Server sessions and flash storage |
| `web/src/types.ts` | Explicit frontend contract and InertiaConfig augmentation |
| Vite + Tailwind plugins | Frontend build and semantic CSS tokens |
| `.storybook` | Component preview configuration |
| `components.json` | shadcn setup only if present; absent in the base starter |

## Step 2: Identify Missing Capabilities

The starter already has the core runtime, GORM, Goose, River and Genkit. Read [stack decisions](../../../docs/stack.md) before adding alternatives. Do not install a Rails serializer, route generator, second ORM/queue or UI framework merely because an upstream skill mentioned one. Go DTO mapping and net/http routing are explicit.

For a requested capability, compare the standard library/current dependencies with a specific addition. Explain the tradeoff when it affects the user's architecture; do not add optional dependencies without a task-related reason.

## Step 3: Configure the Selected Capability

Pin dependencies and update lockfiles. Vite and TypeScript each need the same `@/` alias; there is no vite-plugin-ruby supplying it. shadcn, if requested, must use the existing React/Vite/Tailwind configuration and semantic tokens.

Wire middleware explicitly: sessions, CSRF, Gonertia and routes. Use context-local props for per-request values. Validate production config and assets before serving. Do not invent generated routes/types or claim uploads are supported before implementing the backend contract.

## Step 4: Update Agent Instructions

Read existing AGENTS.md and linked docs. Update the relevant block without deleting manual rules or adding conflicting duplicates. Use [the stack template](references/claude-md-templates.md) as a content guide; the file name preserves upstream provenance and does not require Claude-specific tooling.

## Troubleshooting

| Symptom | Check |
| --- | --- |
| `@/` import fails | Both Vite alias and TypeScript paths/include |
| Blank first page | Root template entrypoint, manifest, page resolver and browser console |
| Inertia 409 after deployment | Asset version mismatch; browser must perform full visit |
| Form 403 | Session cookie and explicit X-CSRF-Token header |
| Errors disappear on redirect | Flash provider/session middleware ordering |
| Shared data crosses users | Global ShareProp used for request-specific values |
| Production requests Vite localhost | Development hot file must be ignored in production |

## Step 5: Verify and Report

Run the documented setup/build/tests for the changed boundary. State installed changes, applicable skills and unverified deployment assumptions. Do not report a planned command as executed.
