# Agent skills

The repository includes 32 skills in `.agents/skills`. They are copied with the GitHub template and available to agents working on the project. The template does not provide a separate skill website, discovery endpoint, or installer that changes global agent settings.

## Design system

The bundled [design-system](../.agents/skills/design-system/SKILL.md) comes from
[Yuri Mandrikov's ai-design-system](https://github.com/ymandrikov/ai-design-system/tree/05e0a9ede21cd5e66f5293164e14f6cd3f441dd0/skills/design-system),
under the [MIT license](../third_party/licenses/ai-design-system.txt). All 31 files match
the author's pinned revision and the existing Evil Martians catalog copy byte for byte.

Use `$design-system` for component selection, contracts, tokens, and source drift checks.
It complements the `sb-*` skills, which supply Storybook stories and visual evidence.
Connecting a project's design sources and generating contracts is a separate `setup`
workflow; the template includes the skill without prefilled `DESIGN.md` or contracts.

## Data systems architecture

Use [data-systems-architecture](../.agents/skills/data-systems-architecture/SKILL.md) to review data invariants, transaction boundaries, queues, retries, schema evolution, and recovery. It complements `layered-go`: package organization and data correctness are related but distinct concerns. Start from a concrete operation and load only the relevant bundled chapters.

The complete standalone package includes 14 chapters, a glossary, patterns, a cheatsheet, decision/operation templates, and 20 evaluation cases. It requires no source book, PDF, extraction tools, or network access. [SOURCES.md](../.agents/skills/data-systems-architecture/SOURCES.md) records its basis in *Designing Data-Intensive Applications*, second edition, and the limits of that synthesis. The pinned archive preserves the original local skill; the patch changes only agent UI metadata to English.

## Original Evil Martians skills

Eighteen skills come from the official EM catalog and are preserved byte for byte: `design-system`, `good-readme`, `intent-log`, `llms-visibility`, `secure-npm-package`, `skills-visibility`, `tailwind-best-practices`, and eleven `sb-*` skills: audit, explore, figma, flows, health, hub, inventory, setup, ship, stories, and wrappers.

These use the original React/Storybook workflows. Source URLs, downloaded-content digests, and per-file SHA-256 hashes are recorded in [skills-lock.json](../config/skills-lock.json). CI checks prevent unnoticed project-specific edits to those originals. Project constraints belong in the root `AGENTS.md`.

`clear-writing` is preserved byte for byte from the `rails-template` commit recorded in the lockfile. One `SKILL.md` selects a complete Russian or English edition by the target text's language. Each edition includes a guide, six chapters, patterns, a cheatsheet, a glossary, and review cases. The language of a request alone does not authorize translating the text being edited. When updating it, copy the entire directory and update the source commit, archive, mapping, and hashes together.

## Book-to-skill

The [original skill](../.agents/skills/book-to-skill/SKILL.md) converts books and documents into agent skills. The complete upstream package is preserved without changes: instructions, Python extractor, tools, documentation, tests, and [MIT license](../.agents/skills/book-to-skill/LICENSE.md). Its source is [virgiliojr94/book-to-skill at commit 80ae087](https://github.com/virgiliojr94/book-to-skill/tree/80ae087784ddbc21dbbfde355fe5509631e0e322); the archive and every file's SHA-256 hash are pinned in the lockfile.

In Codex, invoke `$book-to-skill` and provide the document path. The extractor requires Python 3.9+. Check available format handlers from the project root:

```sh
mise exec -- python3 .agents/skills/book-to-skill/scripts/extract.py --check
```

Format-specific Python packages are installed as needed and are not application runtime dependencies. MOBI/AZW requires Calibre; enhanced technical PDF extraction uses Docling. See the [upstream installation guide](../.agents/skills/book-to-skill/docs/install.md).

## Go adaptations

| Source skill | Go skill | Retained method and adapted mechanisms |
| --- | --- | --- |
| `layered-rails` | `layered-go` | All 57 source documents map to Go versions, with additional sources and installed-library contracts. Preserves layers, specification tests, extraction criteria, anti-patterns, workflows, and refactoring examples. |
| `inertia-rails-architecture` | `inertia-go-architecture` | Server-driven navigation, state ownership, and decision trees; transport and persistence use Go. |
| `inertia-rails-controllers` | `inertia-go-controllers` | Rendering, shared props, authorization, and loading with actual Gonertia signatures. |
| `inertia-rails-forms` | `inertia-go-forms` | React Form/useForm, uploads, and multi-step forms; strict JSON, CSRF, 303 responses, and error-bag/Precognition limits. |
| `inertia-rails-pages` | `inertia-go-pages` | React layouts, navigation, deferred data, and scrolling; explicit resolver and Gonertia contracts. |
| `inertia-rails-setup` | `inertia-go-setup` | Existing-stack inspection; Go modules, Vite, TypeScript, and explicit configuration. |
| `inertia-rails-testing` | `inertia-go-testing` | HTTP, forms, access, partial/deferred requests; `httptest`, PostgreSQL, and Playwright. |
| `inertia-rails-typescript` | `inertia-go-typescript` | React types and Inertia augmentation; manual DTOs replace Ruby-specific generation. |
| `shadcn-inertia` | `shadcn-inertia` | React component integration with Go CSRF, types, and CSP compatibility. Installing shadcn remains a separate decision. |
| `rails-serialization` | `go-serialization` | Entity/page/shared/loading contracts; Go DTOs and explicit serialization. |
| `rails-boot-profiling` | `go-boot-profiling` | Baseline, narrowing scope, external operations, deep profiling, and export; inittrace, pprof, and readiness measurements replace require-profiler. |

The Inertia/React foundation comes from preserved EM originals in `rails-template` history, before its Hotwire adaptation. Every source commit and path is pinned in the lockfile. Only standalone Vue/Svelte examples were removed from forms/pages because this starter uses React; the originals remain in the provenance archives.

[third_party/skills](../third_party/skills) contains a source archive and complete patch for each adaptation. The lockfile's `mapping` relates source and destination files, including ten renamed Rails mechanisms in `layered-go`. Use these records to review changes or continue an adaptation from a specific original.

## Go sources and recommendation boundaries

The material was checked on September 27, 2026 with Go 1.27.1. The [primary-source list](../.agents/skills/layered-go/references/go-sources.md) includes the Go team, Sameer Ajmani, Damien Neil, Jonathan Amsterdam, Vlad Saioc, Alex Edwards, Dave Cheney, and Ardan Labs. Publication dates are distinct from review dates; older principles are not presented as new developments.

`layered-go` accounts for the installed GORM stack without requiring a universal repository, an interface for every type, or a directory for every conceptual layer. Explicit operations and transactions replace callbacks; actor/tenant arguments replace `Current`; parameterized queries replace relations/scopes; composition replaces concerns. Error and concurrency guidance uses the actual contracts of `context`, pgx, and the pinned Go version.

Skills include extension examples, not a list of installed application features. Email, AI, and workers are implemented with go-mail, Genkit, and River; see the [library guide](stack.md) for current boundaries. Other delivery channels, a general outbox, realtime, uploads, and shadcn require task-specific implementation. New frontend-library capabilities do not imply automatic support in the server adapter.

## Updating skills

First edit the relevant skill and check its behavior with a realistic scenario. For local adaptations, review the changes, then run `python3 scripts/update-skill-lock.py` (Python 3 standard library only), followed by `mise exec -- bin/ci`. The script preserves source archives and recalculates hashes and patches; update mappings for deleted or renamed files explicitly.

For an upstream skill, obtain the original from the official catalog, verify its digest, review its changes, and update the URL, digest, and file hashes together. Recalculating hashes is not a substitute for review. CI checks integrity, links, and fences; semantic correctness also requires checking APIs and working scenarios.
