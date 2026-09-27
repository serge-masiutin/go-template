# Go Adaptation Sources

Reviewed on **2026-09-27**. Runtime baseline: **Go 1.27.1**. These sources supplement the original Layered Rails framework; they do not replace its extraction criteria, specification test or refactoring scenarios. Publication dates and review dates are different: older architectural advice is retained when it remains applicable to the pinned toolchain.

## Current Toolchain and Concurrency

- [Go release history](https://go.dev/doc/devel/release): Go 1.27.1 was released on 2026-09-01. Check this source when updating the template; do not infer installed versions from a blog post.
- [Go 1.27 release notes](https://go.dev/doc/go1.27): the current language supports generic methods, but their availability does not justify generic service/repository frameworks. Use new language features only where they simplify a concrete contract.
- [Vlad Saioc, Goroutine Leak Profiles, 2026-09-02](https://go.dev/blog/goroutine-leak-profiles): Go 1.27's `goroutineleak` profile helps diagnose permanently blocked goroutines. It does not detect every leak, particularly network/file I/O stalls or primitives reachable from live roots. Apply it through controlled diagnostic access, not a public debug endpoint.
- [Damien Neil, Testing Time, 2025-08-26](https://go.dev/blog/testing-time): `testing/synctest` is stable since Go 1.25. Use it for suitable in-process time/concurrency tests; PostgreSQL contention still needs real database connections.

## Packages and Abstractions

- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments): place small interfaces at their consumers when actually needed; return concrete implementations; give goroutines an explicit lifetime. Do not create interfaces only to make every type mockable.
- [Organizing a Go module](https://go.dev/doc/modules/layout): `cmd` and `internal` suit a server with multiple executables. The four conceptual layers do not require four directory trees.
- [Sameer Ajmani, Package names, 2015-02-04](https://go.dev/blog/package-names): name packages by a coherent responsibility and read exported names in their package context. Avoid `utils`, `common`, and one universal `interfaces` package.
- [Alex Edwards, Eleven Tips for Structuring Your Go Projects, 2025-01-22](https://www.alexedwards.net/blog/11-tips-for-structuring-your-go-projects): let package boundaries follow actual needs. A large file or a small number of packages is not by itself a defect; mechanically recreating Rails directories adds friction.
- [Alex Edwards, Go Naming Conventions, 2026-03-24](https://www.alexedwards.net/blog/go-naming-conventions): use Go naming and visibility conventions, including behavior-oriented interface names. An `-er` suffix is not itself a service smell.
- [Dave Cheney, Practical Go](https://dave.cheney.net/practical-go/presentations/qcon-china.html): evaluate designs through maintainability and the caller's needs. This is established design guidance, not a 2026 release announcement; use current standard-library error APIs rather than historical third-party examples.
- [Ardan Labs, Ultimate Go: Interfaces](https://tour.ardanlabs.com/tour/eng/interfaces/1): interfaces describe required behavior and enable composition; they are not a substitute for class inheritance.

## Context, Errors and Persistence

- [context package](https://pkg.go.dev/context): pass context explicitly, normally first; cancel derived contexts and use values for request-scoped metadata, not optional function arguments. This template additionally chooses explicit actor/tenant operation parameters so business authority is visible.
- [Damien Neil and Jonathan Amsterdam, Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors): use `errors.Is`/`errors.As` and deliberate wrapping. Exposing a wrapped driver error is an API commitment; translate at a stable abstraction boundary when callers should depend on a feature error instead.
- [Executing transactions](https://go.dev/doc/database/execute-transactions): keep all transactional work on the transaction handle, check commit errors and roll back failed work. Do not mix pool calls into a transaction unintentionally.
- [pgx v5.11.0](https://pkg.go.dev/github.com/jackc/pgx/v5@v5.11.0): verify the actual driver's behavior. Unlike `database/sql` transaction context behavior, cancellation of the context passed to pgx `Begin` does not automatically roll back the transaction; explicitly close the transaction.
- [Alex Edwards, The Fat Service Pattern, 2022-08-08](https://www.alexedwards.net/blog/the-fat-service-pattern): a service can be a useful typed use-case boundary shared by HTTP and CLI. It is one design option, not a mandate to remove domain behavior or create a service for every endpoint.
- [Alex Edwards, When Is It OK to Panic in Go?, 2025-03-31](https://www.alexedwards.net/blog/when-is-it-ok-to-panic-in-go): expected operational/input failures use errors. Do not translate Rails exceptions into panics for ordinary validation or database failures.

## Adaptation Decisions

These are this template's choices informed by the sources, not claims of universal Go consensus:

| Original Rails mechanism | Go decision | Why |
| --- | --- | --- |
| Lifecycle callbacks | Explicit domain methods and application transactions | Makes execution order and side effects visible |
| `Current` | Explicit actor/tenant; context for I/O lifecycle and request metadata | Prevents hidden authority and cross-request state |
| Active Record relations/scopes | Parameterized SQL in a feature store or focused query | Keeps persistence explicit without recreating an ORM |
| Concerns and base classes | Named composition and narrow consumer interfaces | Matches Go's type system and dependency ownership |
| Controller-only authorization | Shared operation enforcement, scoped entry-boundary reads | HTTP, CLI and workers must obey the same use-case contract |
| Generated one-line jobs | Typed worker entry points where they own queue semantics | Thin transport adapters can still have a real responsibility |
| File/LOC thresholds | Inspection signals, never automatic extraction rules | Cohesion and change cost matter more than size |

The detailed chapters retain the original methodology while making these differences explicit. See the repository's skill provenance manifest for exact source paths and revisions.
