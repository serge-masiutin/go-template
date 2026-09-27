---
name: go-boot-profiling
description: Profile Go process startup and readiness to separate compilation, package initialization, configuration, database connection, and asset setup. Use for slow startup or readiness regressions, not ordinary request latency.
---

# Go Boot Time Profiling

Preserve the baseline → narrow scope → inspect external work → deep profile → export
workflow. Go imports are compiled; there is no runtime `require` tree or Rails initializer
pipeline. Use Go's actual tools instead of emulating Ruby loading mechanics.

## Project Context and Prerequisites

- Use `mise exec --` and the versions in `mise.toml` / `go.mod`.
- Build assets and the server before measuring process startup. `go run` includes build
  overhead and `time server` measures its whole lifetime, not readiness.
- Keep binary, environment, DB/cache state, hardware, and readiness criteria consistent.
  Separate cold build, warm build, cold process, and warm external dependencies.
- Write local profiles under ignored `tmp/`. Do not record credentials, connection URLs,
  request bodies, or user data. Do not enable public profiling endpoints.
- Baseline tooling is the standard library. CPU/trace instrumentation below is temporary
  diagnostic code, not a preinstalled startup-profile flag or management command.

## Step 1: Run a Full Startup Baseline

```sh
mise exec -- npm run build
mise exec -- go build -o tmp/server ./cmd/server
```

Load the same environment as `bin/dev` (without running Vite), start the compiled binary,
and measure elapsed monotonic time until `/health/ready` succeeds. Bound the wait, fail
if the child exits, terminate it with SIGTERM, and wait for its exit. Repeat enough times
to report the distribution and measurement resolution, not a single lucky sample.
The ready probe measures DB availability too; keep the probe interval fixed and report
it. Never time an already running server or silently reuse its port.

Capture package initialization separately with the standard runtime diagnostic:

```sh
GODEBUG=inittrace=1 ./tmp/server 2>tmp/inittrace.log
```

This starts the server in the foreground; use another terminal for the readiness probe,
then stop it. The log reports package initialization start time, clock duration, and
allocations. It excludes ordinary work performed in `main`, including the starter's DB
connection and asset setup. Keep existing `GODEBUG` settings if the target depends on
them, adding `inittrace=1` instead of overwriting the other diagnostic controls.

## Step 2: Narrow the Scope

Rank `inittrace` entries by initialization duration and allocated bytes, retaining the
original log. Search a package name with `rg` when useful: unlike a Ruby require tree,
these entries are flat package records, not nested call stacks. There is no runtime
`INIT_PROFILE_THRESHOLD` or `REQUIRE_PROFILE_FOCUS` flag.

A tiny init total with slow readiness directs the investigation to `main` and external
work. A large package init directs it to that package's variable initialization and
`init` functions. Check transitive imports with `go list -deps ./cmd/server`; package
count alone is not a latency diagnosis.

## Step 3: Check Configuration and External Activity

Inspect configuration parsing, DNS, TLS, database authentication, pool startup, remote
configuration, and asset manifest reads. This starter explicitly opens/pings PostgreSQL
under a startup deadline and loads the Vite manifest before serving. It does not fetch
remote config or migrate implicitly on server startup.

Add temporary monotonic timing around these real boundaries. Log the operation name,
duration, and a safe error category; never log DSNs or unredacted errors from drivers.
Compare cold/warm external dependencies without disabling validation or authentication.
A CPU profile cannot attribute time spent waiting for a remote server by itself.

## Step 4: Inspect the Explicit Initialization Pipeline

Follow `cmd/server.run`: config → database → session store → HTTP/asset setup → listener.
If the process reports startup before the listener has successfully bound, use the
actual readiness probe as the measurement endpoint. When startup has parallel work,
record its critical path; summing overlapping spans overstates wall-clock startup.
Keep timeouts, cancellation, required dependency checks, and cleanup intact.

## Step 5: Deep-Dive with pprof or Runtime Trace

After locating a costly phase, temporarily wrap only that phase in
`runtime/pprof.StartCPUProfile(file)` / `StopCPUProfile()`. Check file and profiler errors,
stop before closing the file, and close on all exits. Start profiling before the phase
and stop at readiness, not at eventual process shutdown. Instrumentation in `main`
cannot capture package initialization that already ran; use `inittrace` for that portion.
For very short CPU phases, sampling may collect too little data: repeat an isolated
representative operation in a benchmark instead of inventing precision.

```sh
mise exec -- go tool pprof -top ./tmp/server tmp/startup.cpu.pprof
mise exec -- go tool pprof -top -cum ./tmp/server tmp/startup.cpu.pprof
mise exec -- go tool pprof -list 'database.Open' ./tmp/server tmp/startup.cpu.pprof
mise exec -- go tool pprof -http=127.0.0.1:6060 ./tmp/server tmp/startup.cpu.pprof
```

Use the interactive `focus`, `ignore`, and `list` commands to retain the call graph.
Separate flat/self cost from cumulative/descendant cost. The web graph can require
Graphviz; text reports do not. Profiles and their matching binaries are local artifacts.
For blocked/concurrent startup, temporarily use `runtime/trace.Start(file)` / `Stop()`
around the phase, then `go tool trace tmp/startup.trace`. Observe instrumented overhead;
use uninstrumented repeated readiness measurements to confirm the improvement.

## Step 6: Export for Handoff

Share the native `.pprof` or trace file with the matching build revision/toolchain,
measurement method, sanitized timing results, and environment conditions. `inittrace`
is a text log, not a CPU profile. Do not rename a native pprof file to JSON or claim it
is Speedscope JSON. Use the Go toolchain viewer, or explicitly verified conversion if
the recipient needs another format. Never upload profiles without authorization.

## Tool Reference

| Tool | What it measures | Important limit |
| --- | --- | --- |
| readiness harness | process start through a successful ready probe | probe interval and external state affect results |
| `GODEBUG=inittrace=1` | package initialization clock time and allocations | not work in `main`; not a call tree |
| `go tool pprof` | sampled CPU or other captured runtime profiles | blocked wall time is not CPU cost |
| `go tool trace` | captured scheduler/runtime events | instrumentation changes timing |
| `go test -bench ... -benchmem` | repeatable isolated hot operations | not whole-process startup |

## Recommended Agent Workflow

1. Establish the actual complaint and an uninstrumented, repeated baseline.
2. Separate build time, package initialization, explicit setup, and dependency readiness.
3. Rank measured costs; investigate the largest contributor before optimizing.
4. Check boot-time side effects, retries, deadlines, and failure behavior.
5. Profile the responsible phase with enough samples and matching binaries.
6. Export evidence and state limitations and the next testable hypothesis.
7. Make one local change and remeasure with identical conditions. Move heavy work out
   of `init` into explicit setup when it clarifies ownership; defer optional work only
   when its new first-request latency and failure semantics are acceptable. Do not hide
   required dependency failures or add caching/concurrency without measured benefit.

## Sources

- [Runtime GODEBUG / inittrace](https://pkg.go.dev/runtime#hdr-Environment_Variables)
- [Go diagnostics](https://go.dev/doc/diagnostics)
- [runtime/pprof](https://pkg.go.dev/runtime/pprof) and [runtime/trace](https://pkg.go.dev/runtime/trace)
- [Go profiling](https://go.dev/blog/pprof)

Check source/toolchain compatibility when applying this workflow. For goroutine leak
analysis and current Go guidance, see [layered-go sources](../layered-go/references/go-sources.md).
