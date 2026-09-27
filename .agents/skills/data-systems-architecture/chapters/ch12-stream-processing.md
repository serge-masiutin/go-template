# Streams, CDC, event time, and recovery

## Core idea

A stream is an unbounded sequence of events. Processing it correctly requires contracts for delivery, ordering, time, retained state, and visible effects. A low processing delay does not imply complete or correct results.

## Choose delivery semantics intentionally

Traditional task-queue brokers commonly distribute messages among workers and remove acknowledged work. Durable log brokers retain ordered partitions and let consumer groups track positions. Real products offer variants; inspect the configured retention and acknowledgment semantics.

| Need | Useful mechanism | Important boundary |
|---|---|---|
| Independent tasks with variable duration | Per-message acknowledgment and redistribution | Redelivery after worker loss; execution order may change |
| Replay and several independent projections | Retained partitioned log and consumer offsets | Retention limit and ordering only within its declared scope |
| Parallelism with ordered per-entity changes | Stable entity-to-partition assignment | Hot keys, partition migration, and consumer handoff |

Acknowledging before an effect commits risks loss. Applying an effect and acknowledging afterward permits duplicates after a crash. Atomic coordination or idempotent effects close that gap within their documented scope; a delivery label alone does not.

A log offset identifies a position in one partition, not a universal business operation and not a global order. Producer retries can create the same logical event at different offsets unless separately prevented. Preserve business identity when deduplication must cross producer retries or replayed topics.

## Treat overload as a data-management problem

Track consumer lag in both positions and age, arrival/processing rates, retained bytes, and recovery throughput. If processing rate never exceeds arrival rate, a backlog cannot drain. A finite retention period makes a prolonged outage a possible loss of rebuildability.

Choose bounded buffering, backpressure, admission control, or an explicitly lossy policy according to the product contract. Skipping poison events to unblock a stream may violate state integrity. Quarantine with an observable gap and repair path when appropriate; do not mark a projection complete while silently omitting required events.

## Use CDC for committed changes

Change data capture observes changes at the authoritative database and publishes a representation for downstream consumers. It avoids the race in independent database-plus-index writes when the derivation preserves the source's relevant order and applies changes correctly.

Define the envelope: source/table identity, key, operation, schema/version, source position, transaction information if needed, and delete semantics. A row-change record explains what changed; it is not always a domain event explaining why. An outbox can carry application-defined event intent in the same transaction as business state.

To bootstrap a consumer:

1. Obtain a consistent snapshot and its associated log position using a supported connector/database procedure.
2. Retain the required change history while the snapshot is copied and applied.
3. Resume changes at the correct boundary without a gap; tolerate duplicates if the protocol permits them.
4. Preserve order and reconcile updates/deletes that race the snapshot.
5. Verify the target against a matching source point before declaring it caught up.

Log compaction retains useful latest-per-key state and tombstones according to its policy; it is not a full audit history. Deleting tombstones or expiring history can make an old consumer's recovery unsafe. Document when a consumer must rebuild instead of continuing from an unavailable position.

## Distinguish event sourcing, state, and change logs

A current-state table can be viewed as the accumulated effect of a change stream; a stream can expose changes to a table. This duality is useful for incremental views, but the representations preserve different information.

Event sourcing retains accepted domain events and intent. CDC commonly exposes database mutations. A compacted stream may preserve current state but lose intermediate transitions. Choose the representation according to required replay, audit, and correction semantics rather than calling every message “an event” with the same guarantees.

Immutability simplifies replay only if interpretation and inputs are controlled. Corrective events, schema evolution, privacy deletion, and external effects require explicit policies.

## Define time and window semantics

**Event time** describes when the source says something happened. **Processing time** describes when an operator handles it. **Ingestion time** describes entry at a specified system boundary. Record which clock assigns each and what uncertainty it has.

After an outage, processing-time counts can show an artificial spike while replaying a steady event-time workload. For historical meaning use event-time windows with an explicit late-data policy. Device timestamps may be wrong; correction based on send/receive times depends on assumptions about network delay and clock drift and should retain uncertainty.

| Window | Meaning | Review question |
|---|---|---|
| Tumbling | Fixed, non-overlapping intervals | Are boundaries and time zones explicit? |
| Hopping | Fixed-length intervals emitted at a smaller step | Can events contribute to several outputs? |
| Sliding | Window relative to events or a moving boundary | What exact implementation semantics apply? |
| Session | Events grouped by inactivity gaps | Can a late event merge previously separate sessions? |

A watermark represents progress under assumptions about future timestamps; it is not a magical proof that no earlier event exists. Idle or delayed partitions complicate combined progress. Define allowed lateness, state retention, and whether late records trigger correction/retraction, a side output, or observable rejection. Some sources provide stronger progress guarantees than heuristic watermarks; document which is used.

Consumers must understand whether an output is provisional, updated, or final. Repeatedly publishing a revised count as an additive increment double-counts. Include output identity/version or changelog semantics.

## Choose join semantics

- **Stream–stream:** correlate events within a defined interval, retaining candidates from both sides. Handle either arrival order, unmatched records, and late matches.
- **Stream–table:** enrich events with state. Decide whether the lookup means current-at-processing-time or valid-at-event-time; the latter requires suitable history/versioning.
- **Table–table:** maintain a materialized join as either input changes. Inserts, updates, deletes, and relationship changes can all affect output.

Ordering within each input partition does not define order across input streams. If replay order changes which dimension value joins an event, the output is not reproducible without additional temporal/dependency semantics.

Budget state by key cardinality, retained event rate, windows, allowed lateness, and data sizes. A logically small operator can hold unbounded state if no expiry condition is defined. State TTL is a semantic choice when it can discard a future match.

## Recovery and exactly-once scope

Recover operator state and input positions consistently. Checkpoints, changelogs, replay, and microbatches are mechanisms; each has retention, coordination, and restoration costs. A local state snapshot without the matching input position can lose or duplicate updates.

Framework-managed exactly-once state can require transactions with supported sinks or a destination idempotency protocol. It does not automatically cover HTTP calls, email, or arbitrary databases. Identify the boundary encompassing inputs, state, output, and progress before making an exactly-once claim.

For an ordered per-partition database projection, a conditional update recording both the projected change and consumed position can reject replayed work if ownership and order assumptions hold. A single greatest-seen offset is unsafe if events finish out of order: accepting offset 12 before 11 can cause 11 to be incorrectly skipped. Use serial partition processing, a contiguous commit frontier, or per-event deduplication as appropriate. Fence old consumers during handoff where the sink requires it.

## Worked example: search impressions and clicks

**Reconstructed source scenario:** compute click-through rate by joining impression events and click events. A click can arrive before its impression, much later, or never.

Use a stable impression identity and a defined event-time matching interval. Retain both sides and join in either arrival order. Counting only clicks carrying embedded impression details loses the denominator of impressions with no click.

Suppose a five-minute window emits 100 impressions and 10 clicks. A delayed eligible click arrives after emission. Under a correction policy, publish a replacement/versioned aggregate of 11 clicks or an explicit delta with stable identity. Under an approved drop policy, keep 10 but record the dropped late event and the resulting completeness limit. Never silently choose between those meanings.

Crash after updating state but before progress is saved. Recovery should reproduce the intended aggregate without a second visible increment. Rebuild the report with side effects disabled; a metrics replay must not resend a user notification.

## Validation

Test duplicate logical events, reversed arrival order, late data, idle partitions, changing dimension values, deletion, source-history expiry, checkpoint recovery, and consumer overlap. Assert both aggregate correctness and the visible correction protocol. Measure time to catch up after a realistic outage.

Related: [encoding](ch05-encoding-and-evolution.md), [transactions/inbox](ch08-transactions.md), [batch](ch11-batch-processing.md), [dataflow integration](ch13-dataflow-and-correctness.md).
