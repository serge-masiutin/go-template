# Patterns and selection criteria

Use these as routes into the linked mechanisms. A pattern name is not proof that its implementation preserves the invariant.

## 1. Authoritative state with derived views

**When:** several representations serve different queries.
**How:** define one authority per fact; derive through committed changes; version transformations and progress; demonstrate rebuilds.
**Trade-offs:** read flexibility and failure isolation versus lag, retained history, and repair work. [Details](chapters/ch13-dataflow-and-correctness.md).

## 2. Materialization at write time

**When:** frequent reads repeatedly compute an expensive result.
**How:** maintain a view from source changes, including deletes and relationship updates; measure write fan-out and read benefit.
**Trade-offs:** cheap reads versus extra writes, hot keys, freshness, and rebuild cost. Hybridize only where skew justifies it. [Details](chapters/ch02-requirements-and-reliability.md).

## 3. Atomic condition and mutation

**When:** an invariant fits an authoritative single-resource operation.
**How:** use a supported conditional update/constraint; inspect affected rows or conflict result; include all writer paths.
**Trade-offs:** a small, clear boundary, but insufficient for predicates or invariants outside its scope. [Details](chapters/ch08-transactions.md).

## 4. Serializable decision transaction

**When:** a decision depends on multiple records or absence predicates.
**How:** read and change state in a documented serializable transaction; retry the whole transaction on classified serialization failures; keep external effects out.
**Trade-offs:** simpler invariant reasoning versus contention, aborts, and retry load. [Details](chapters/ch08-transactions.md).

## 5. Transactional outbox and inbox

**When:** a committed local change must reliably cause downstream work.
**How:** commit state and outbox intent together; relay committed entries; atomically claim an event and apply a consumer's local effect.
**Trade-offs:** closes loss gaps while allowing duplicate delivery; ordering, retention, and external effects still need contracts. [Details](chapters/ch13-dataflow-and-correctness.md).

## 6. End-to-end idempotency key

**When:** a caller may retry an effect after an uncertain response.
**How:** preserve one logical identity across attempts; atomically store identity, request meaning, effect, and result where possible; reject conflicting reuse.
**Trade-offs:** bounded duplicate suppression requires retained state and an explicit replay horizon. Separate intentional identical actions receive different identities. [Details](chapters/ch13-dataflow-and-correctness.md).

## 7. Read-after-write progress token

**When:** replicas/projections may lag but a caller must see its own write.
**How:** carry a committed source position; serve only from a representation at or beyond it, wait, or return the defined pending/failure result.
**Trade-offs:** freshness versus latency/availability; cross-device use needs shared context. [Details](chapters/ch06-replication.md).

## 8. Fenced ownership

**When:** a paused or delayed old worker could corrupt a protected resource.
**How:** issue ordered ownership epochs; atomically persist and enforce the current epoch at the resource; reject older writes.
**Trade-offs:** requires resource support and a sound epoch protocol. Lease expiry or worker-side checking alone is insufficient. [Details](chapters/ch09-distributed-failures.md).

## 9. Snapshot plus change-stream catch-up

**When:** bootstrap a replica, projection, or migration target.
**How:** obtain a consistent snapshot with a matching log position; retain history; apply concurrent changes without a gap; compare matching positions before cutover.
**Trade-offs:** online evolution versus catch-up capacity, retention pressure, and overlap handling. [Details](chapters/ch12-stream-processing.md).

## 10. Expand–migrate–contract

**When:** schemas or representations evolve while old clients remain.
**How:** add compatible support, migrate/backfill, verify mixed versions, switch readers, retire old state after rollback/retention conditions permit.
**Trade-offs:** temporary coexistence costs buy reversibility; incompatible semantic changes require explicit handling. [Details](chapters/ch05-encoding-and-evolution.md).

## 11. Versioned dataset publication

**When:** readers need complete batch output or a rebuilt index.
**How:** write an isolated generation, validate it, publish via a supported commit/manifest mechanism, and retain a rollback version.
**Trade-offs:** extra storage and publication coordination avoid partial results. A buffering stream does not supply atomic activation. [Details](chapters/ch11-batch-processing.md).

## 12. Conflict-key partitioning

**When:** decisions can be scoped to an entity, tenant, or unique value.
**How:** route conflicting operations to one authoritative partition; colocate related data where useful; specify cross-partition operations separately.
**Trade-offs:** local reasoning and throughput versus skew and global query cost. [Details](chapters/ch07-sharding.md).

## 13. Domain-aware merge

**When:** offline or multi-leader writes intentionally diverge.
**How:** detect causal/concurrent versions; choose a suitable CRDT or explicit domain/user resolution; model deletion semantics.
**Trade-offs:** availability and local work versus conflict complexity. Convergence does not enforce every business invariant. [Details](chapters/ch06-replication.md).

## 14. Event-time correction protocol

**When:** delayed events must contribute to historically meaningful results.
**How:** define timestamps, windows, watermarks, allowed lateness, and replacement/retraction identity; recover state and progress together.
**Trade-offs:** completeness versus delay and retained state. Dropped late data requires an explicit observable policy. [Details](chapters/ch12-stream-processing.md).

## 15. Independent integrity audit

**When:** incorrect derived state or missed effects would be costly.
**How:** compare invariants or independent representations at compatible source positions; keep provenance and a bounded repair procedure.
**Trade-offs:** computation and operational ownership detect errors that lower-level guarantees miss. [Details](chapters/ch13-dataflow-and-correctness.md).

## 16. Purpose-bound retention and deletion propagation

**When:** personal data spreads into logs, projections, exports, or models.
**How:** minimize collection; assign retention and access per purpose; propagate corrections/deletes; apply deletion policy during replay and restore.
**Trade-offs:** less exposure versus limits on historical reconstruction; decide deliberately. [Details](chapters/ch14-responsible-data-systems.md).
