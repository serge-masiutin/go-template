# Dataflow integration and end-to-end correctness

## Core idea

Compose specialized stores through explicit derivation paths while preserving end-to-end invariants. A correct database, broker, and workflow engine can still form an incorrect application if operation identity, ordering, publication, or recovery breaks between them.

## Draw the write path and read path

For each user-visible operation, show where input is accepted, which authoritative facts commit, how changes propagate, and what the user reads afterward. Label synchronous dependencies, asynchronous boundaries, ordering scope, durable positions, and permitted lag.

Treat caches, search engines, analytical tables, and learned representations as derived state when they are reconstructible. Pick the authoritative write path and derive other copies from its committed history. Letting callers independently update every copy creates partial-failure and reordering races.

The same pair of concurrent updates can arrive in opposite order at two stores. Retrying until both writes succeed does not repair that ordering disagreement. A single authoritative order per conflict domain plus correctly applied derivation does.

## Use an outbox/inbox when it fits

This is a concrete engineering application of the source's local transactions, durable logs, and duplicate-suppression reasoning:

1. In one local transaction, change the business state and persist an outbox event with stable identity and sufficient payload/version.
2. A relay reads committed outbox entries and publishes them. It marks progress only under its documented acknowledgment protocol.
3. Assume the relay can publish twice if it crashes after publication but before saving progress.
4. A consumer commits its deduplication/inbox claim and local business effect atomically.
5. Retain identities long enough for retries and replay, and monitor undelivered age and failed consumers.

The outbox closes the gap between committing state and durably recording publication intent. It does not guarantee one broker delivery, global event order, or atomic application across consumers. Relay parallelism and transaction commit order need explicit handling when per-entity ordering matters.

CDC can transport outbox records or raw changes. Use supported source positions and snapshots; do not reverse-engineer an internal log format casually. Keep one clear writer/derivation authority for each projection.

## Distinguish timeliness from integrity

**Timeliness** concerns how soon readers see accepted changes. **Integrity** concerns whether facts are preserved and derived results are correct. A delayed index may become correct by catching up; a lost event or repeated debit will not repair itself merely by waiting.

Asynchronous designs can preserve integrity with durable delivery, deterministic processing, deduplication, and repair while accepting bounded or unbounded freshness under stated conditions. They do not thereby provide an atomic cross-system snapshot.

Where a caller needs read-your-writes, return an operation/position token and provide a status or read path that waits for the required projection progress. Return an honest pending state when appropriate. A successful enqueue acknowledgment and successful business completion are different events.

Read paths also deserve optimization: move repeated work into materialized views when justified, keep appropriate data near readers, and update subscriptions carefully. Rebuilding a view should preserve its observable contract; incremental notification delivery alone is not proof of a consistent snapshot-plus-subscription transition.

## Operation identity must span retries

Create one identity for one logical user intent and carry it through the request, authoritative transaction, messages, and downstream effects. Two intentional identical purchases need different identities; one retried purchase needs the same identity. A payload hash alone may conflate repeated legitimate operations.

Store the identity with the accepted request semantics and result. Reuse with a different payload is a conflict, not a request to silently return an unrelated old result. Scope keys by the appropriate tenant/principal/operation to avoid collisions or unauthorized result access.

Atomically record the identity and effect at each local boundary. Retrying after an uncertain commit should query/reuse that identity. If a downstream system supports idempotency, use its specified retention and scope; if it does not, plan reconciliation or a visible unresolved state. A refund after a double charge is compensation, not proof the original charge occurred once.

## Coordination belongs around the invariant

A strict uniqueness or scarce-resource rule requires conflicting claims to meet at an authoritative decision point. Partition by the conflict key where possible; all claims to the same normalized name can be checked in one partition. Per-partition ordering still does not make an operation involving several independent keys atomic.

An ordered request stream can decide claims sequentially and emit accepted/rejected outcomes. The caller must distinguish “claim recorded” from “claim accepted.” The decision processor's state, output, and recovery need integrity guarantees; a log alone does not supply them.

Some domains permit temporary violations followed by compensation. Treat this as an explicit business contract: maximum exposure, detection deadline, customer effect, and repair ownership. Never convert a hard invariant into a soft one just to simplify the architecture. Irreversible disclosure or a physical shipment may not be meaningfully compensatable.

## Rebuilds and migrations are normal operations

To introduce a new projection or representation:

1. Version its schema, transformation, and relevant external inputs.
2. Build a separate generation from a known snapshot/history.
3. Apply concurrent changes through a coordinated boundary; avoid gaps or old backfill records overwriting newer live state.
4. Compare results at matching source positions, including deletions and rejected inputs.
5. Catch up, switch readers through an explicit publication mechanism, and retain a bounded rollback path.
6. Retire the old version only when compatibility, recovery, and retention conditions allow it.

Batch and stream processing can share logic and contracts. Maintaining two unrelated implementations of the same business calculation invites drift. Conversely, forcing a historical rebuild through a low-throughput live path can be impractical; separate execution strategy from semantic definition.

Historical enrichment must remain historical. A purchase replay using today's currency rate or today's permission state does not necessarily reproduce the accepted facts. Preserve the decision inputs or specify that the new projection intentionally answers a different question.

## Worked example: an accepted order and its notification

**Original application of the source's end-to-end argument:** the browser sends order intent `O-17`. The server commits the order and an outbox event together, then loses its response.

- The browser retries `O-17`; a unique accepted-request record returns the existing result without creating another order.
- The relay publishes the event, crashes before recording progress, then publishes again.
- The inventory consumer's inbox and inventory mutation share one transaction, so the repeated event does not reserve twice.
- An email destination may have no idempotency API. The system cannot honestly promise exactly one email merely because the order and inventory are deduplicated. Choose a delivery policy and record/reconcile uncertain sends.
- Search projection lag does not invalidate the accepted order. The confirmation page reads the authoritative result or waits for its required projection position.
- Rebuilding search consumes order history without re-executing inventory or email effects.

If the inventory invariant actually requires order acceptance and reservation to be atomic, this asynchronous flow is insufficient without an explicit pending/reservation protocol. Start by deciding that contract, not by assuming the outbox solved it.

## Audit correctness independently

Assertions at write time are necessary but can be undermined by application bugs, engine bugs, corrupt storage, or faulty migrations. Add checks that compare independent representations or invariants: ledger conservation, missing/duplicate identities, referential relationships, projection counts at a matching position, and end-to-end checksums where meaningful.

An audit should identify the scope and source position checked, not compare a current source against an intentionally lagging target and call all differences corruption. Repair from authoritative evidence, with a controlled plan and observable outcome. Avoid repairing every mismatch by blindly overwriting whichever side is easiest to reach.

Maintain a provenance chain from source identity and schema through transformation version to output generation. This supports diagnosing both operational failures and incorrect business results.

Related: [authority](ch01-architecture-tradeoffs.md), [transaction effects](ch08-transactions.md), [stream recovery](ch12-stream-processing.md), [privacy and deletion](ch14-responsible-data-systems.md).
