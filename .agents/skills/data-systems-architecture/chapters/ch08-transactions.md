# Transactions and concurrency invariants

## Core idea

Choose a transaction boundary and concurrency mechanism that preserve named business invariants. A transaction API or an “ACID” label does not imply every anomaly is prevented. Isolation names and implementation details vary.

## Start with the invariant and all its writers

State the property as a checkable rule: a seat has at most one accepted owner; inventory never becomes negative; at least one responder remains on duty; a logical request changes a ledger once. List every path that can modify the relevant state, including workers, imports, admin tools, and repairs.

Then identify the facts each operation reads to decide whether to write. If another operation changes those facts before commit, what prevents the decision becoming invalid? A preflight read in a controller followed by an unrelated write is not a concurrency guarantee.

## Interpret ACID precisely

- **Atomicity:** a transaction's changes commit together or are aborted. It does not mean every machine instruction is atomic, and it does not include effects outside the transaction.
- **Consistency:** transactions should preserve declared application invariants. The database enforces supported constraints; it cannot infer the business rule from arbitrary application code.
- **Isolation:** concurrent transactions are restricted from observing/interfering with one another according to the selected model.
- **Durability:** committed results survive the specified failures under the configured persistence/replication contract. It is not an unlimited promise against every disaster or bug.

Treat a lost connection during commit as an unknown outcome until the operation identity or database protocol resolves it. Repeating a transaction with a new business identity can apply it twice even if each individual transaction is atomic.

## Recognize the anomaly before choosing the fix

| Anomaly | Characteristic history | Candidate protection |
|---|---|---|
| Dirty read | Read another transaction's uncommitted state | Appropriate isolation, commonly read committed or stronger |
| Dirty write | Overwrite another uncommitted change | Database write-conflict/locking rules |
| Read skew | Read related facts from incompatible moments | Consistent transaction snapshot where sufficient |
| Lost update | Two read–modify–write cycles overwrite one another | Atomic update, conditional version write, lock, or engine-supported conflict detection |
| Write skew | Transactions read a shared premise but update different records | Serializable execution, correctly scoped locking, or an enforceable constraint |
| Phantom-related anomaly | A predicate's matching set changes, including previously absent rows | Serializable/predicate protection, a suitable constraint, or a common locked guard |

Read committed prevents some anomalies but commonly allows successive statements to observe different committed states. Snapshot isolation supplies a coherent snapshot and often detects some write conflicts; it can still allow write skew. A stable snapshot is not proof of serializability. “Repeatable read” must be interpreted using the actual engine/version's semantics.

MVCC preserves multiple versions so readers can observe an appropriate snapshot while writes proceed. Long-lived snapshots can retain old versions and create cleanup pressure. Snapshot visibility and write-conflict checks are different mechanisms; do not infer conditional-write semantics merely from what a plain SELECT sees.

## Prefer the smallest correct concurrency mechanism

**Atomic update:** express a simple change and its predicate at the authoritative database boundary. An illustrative reservation operation is conceptually `decrement available only if available > 0`. Check the affected-row result; zero means no accepted reservation, not success. Verify the chosen engine's conditional update semantics.

**Optimistic version check:** read a version, compute a new value, then update only if the authoritative current version still matches. Increment the version on every relevant mutation and require exactly the expected affected-row count. A compare-and-set that compares against a stale replica is not the same guarantee.

**Explicit locking:** lock all resources that protect the invariant, recheck the premise under those locks, make the change, and commit. Keep transactions short and use consistent lock ordering where possible. A convention is only correct if every conflicting writer follows it.

**Database constraint:** use unique, foreign-key, exclusion, check, or another supported constraint when it truly expresses the invariant. Application validation improves feedback; the authoritative constraint resolves races. Check NULL, collation, tenant scope, deferred evaluation, and partition limitations.

**Serializable transaction:** rely on a documented serializable implementation when invariants span predicates or many records. The application must handle serialization failures by retrying the entire decision transaction with fresh reads, within bounded limits.

Do not bolt on distributed locks when one database's constraints/transactions already cover the resource. A separate lock service creates another failure boundary and still needs fencing at the resource.

## Serializability implementation trade-offs

Serializability means committed behavior is equivalent to some serial execution. It does not require literal one-at-a-time execution and does not alone impose real-time order.

- **Actual serial execution:** avoids concurrency anomalies by scheduling transactions sequentially within its scope. Keep operations short; slow interactive or remote work undermines throughput. Multiple partitions reintroduce coordination for shared operations.
- **Two-phase locking (2PL):** acquires locks, including predicate/range protection where needed, and retains appropriate locks until completion. Waiting, deadlocks, and long-tail latency are costs. Two-phase locking is unrelated to the two phases of distributed commit.
- **Serializable snapshot isolation (SSI):** combines snapshot reads with dependency/conflict tracking and aborts that prevent nonserializable outcomes. Read tracking need not block writes. It is not simply “abort every transaction that read anything stale.” High contention and retries can still become expensive.

Benchmark the transaction mix and retry rate rather than assuming serializable isolation is always prohibitively slow or always free.

## Worked example: write skew in a duty roster

**Reconstructed source example:** two responders, A and B, are on duty. The invariant is that at least one remains.

| Step | Transaction A | Transaction B |
|---|---|---|
| 1 | Snapshot reads A=on, B=on | Snapshot reads A=on, B=on |
| 2 | Sees two responders; sets A=off | Sees two responders; sets B=off |
| 3 | Commits its own row | Commits its different row |
| Result | Both off; invariant broken | No same-row lost update was necessary |

Protecting only each person's row does not protect the shared premise. Options include a serializable transaction with whole-transaction retries, or locking a stable shift/guard row first and rechecking current duty state under a protocol followed by every roster writer.

For booking a previously empty time interval, locking only rows returned by “find conflicting bookings” may lock nothing. Use a supported exclusion/uniqueness constraint, serializable predicate protection, or a common resource row whose lock covers all conflicting reservations. Include concurrent inserts in the test.

## Distributed commit

Two-phase commit (**2PC**) coordinates an atomic commit decision across participants:

1. The coordinator asks participants to prepare; a successful prepare means a participant has durably promised it can commit and retains necessary state/locks.
2. After every participant agrees, the coordinator durably records the commit decision and communicates it; otherwise the protocol decides abort where permitted.
3. A prepared participant with an unknown decision cannot safely choose independently merely because a timeout elapsed. Recovery must learn the authoritative decision.

Coordinator loss can leave transactions in doubt and resources locked. Replication of the coordinator can improve availability but does not erase participant and protocol requirements. Distinguish a database's integrated cross-shard transactions from heterogeneous XA-style coordination among independent systems.

2PC and consensus solve related coordination needs but are not interchangeable labels. A basic 2PC exchange waits for every participant's prepare and can block; consensus protocols have their own quorum, state, and progress conditions. Do not claim “two voting rounds” proves equivalent availability.

## Exactly-once database effect from redelivery

For a message with stable identity, combine an inbox/deduplication record and the business mutation in **one database transaction**:

```text
receive message with stable event_id and consumer identity
begin transaction
  attempt unique claim for (consumer, event_id)
  if the claim already committed: apply no additional business change
  otherwise: validate and apply the business change in this transaction
commit
acknowledge the message after known successful commit
```

This is protocol pseudocode, not a ready-to-run database API. The unique claim must serialize concurrent duplicate deliveries; a separate “SELECT then INSERT” without a constraint is insufficient. A failure before commit leaves no accepted effect. A failure after commit but before acknowledgment causes redelivery, which observes the committed claim.

The guarantee covers that database transaction. Email, payment gateways, another database, and broker publication remain outside it unless separately coordinated. Use an outbox for durable publication intent and destination idempotency/reconciliation for external effects.

Retain deduplication identities for the full possible retry, redelivery, manual replay, and restore horizon. A broker acknowledgment alone does not prove the logical event can never appear again. Reject reuse of a request identity with conflicting payload semantics.

## Validation and evidence

Use deterministic barriers to force competing transactions to read the same premise, then attempt the writes. Check the invariant and error/retry behavior. Test empty predicate matches, duplicate messages, crash-before-commit, crash-after-commit, and an ambiguous commit response. Verify against the selected database, not an in-memory mock that cannot reproduce its isolation.

Related: [consistency](ch10-consistency-and-consensus.md), [partial failures](ch09-distributed-failures.md), [outbox and end-to-end correctness](ch13-dataflow-and-correctness.md).
