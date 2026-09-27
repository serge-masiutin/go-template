# Replication, freshness, and concurrent writes

## Core idea

Replication creates copies to improve durability, availability, proximity, or read capacity. Those benefits require explicit rules for acknowledgment, ordering, conflict handling, and recovery; copies do not automatically behave as one up-to-date database.

## Select a topology by the required write behavior

| Topology | Basic mechanism | Main design burden |
|---|---|---|
| Single leader | One leader orders writes; followers apply its changes | Failover safety, replication lag, leader capacity and location |
| Multiple leaders | Several leaders accept writes and exchange changes | Concurrent writes, merge semantics, loops and replication topology |
| Leaderless | Clients/coordinators send writes and reads to several replicas | Version resolution, quorum assumptions, repair, partial writes |

Single-leader replication is often the simplest starting point for centrally owned mutable state. Multiple leaders can support local writes across regions or disconnected devices, but shift conflict handling into the design. Leaderless replication can make different availability/latency trade-offs; do not infer its guarantees from replica count alone.

## Specify what a successful write means

Synchronous replication waits for defined replicas and a defined persistence/apply condition. It can improve recovery guarantees while increasing latency and making writes depend on those replicas' availability. “Synchronous” may mean receipt, durable logging, or query-visible application; inspect the actual contract.

Asynchronous replication acknowledges before all replicas have the change. It reduces coupling but can lose acknowledged writes when promoting a replica missing them. State the allowed recovery point and the failover policy. Replication of corruption or deletion is another reason backups and recovery procedures remain necessary.

When bootstrapping a replica, copy a consistent snapshot and retain/replay the matching log position. An uncoordinated table copy plus “start tailing now” can miss changes or duplicate them. Monitor whether required history is retained until catch-up completes.

Statement-based, physical write-ahead, and logical row-level replication carry different semantics. Nondeterministic statements, storage-version coupling, and schema changes can matter. Select the mechanism with the intended consumers and evolution contract in mind.

## Make failover a protocol

1. Establish the conditions for suspecting the old leader; a timeout is suspicion, not proof of death.
2. Select an eligible replica based on the required data-loss and committed-history rules.
3. Prevent the old leader from continuing to accept authoritative writes, using the system's fencing/term protocol.
4. Redirect clients and handle in-flight or outcome-unknown operations.
5. Reconcile or rejoin the old node without allowing stale history to overwrite current state.

Automatic failover improves response time only when these semantics are sound. A stale promoted leader can invalidate previously issued identities, operations, or downstream effects. Test failover at acknowledgment boundaries rather than only after the cluster is idle.

## Choose a freshness guarantee per read path

| Need | Meaning | Possible mechanism and limitation |
|---|---|---|
| Read-your-writes | A session sees its accepted changes | Route to authoritative/fresh replica or carry a commit-position token; cross-device sessions need shared metadata |
| Monotonic reads | A session does not move backward in observed version | Sticky replica or lower-bound version token; replica replacement must preserve the bound |
| Consistent prefix | Effects are not observed before their ordered prerequisites | Preserve relevant order/dependencies across replicas and shards |
| Eventual convergence | Replicas converge when updates stop and propagation/repair succeeds | No bounded freshness promise unless separately specified |

Reading from the leader is not a substitute for every consistency proof: stale leaders, failover, cached reads, and multi-object operations need their own guarantees. Likewise, routing a user to the leader for a fixed minute after a write is a heuristic unless the system proves the required lag bound.

A commit/log-position token is often more precise than wall-clock time. The next read can wait until the serving replica has applied that position, choose an eligible replica, or return an explicit freshness failure. Define timeout behavior; never silently claim the write vanished because a stale replica cannot see it.

## Understand quorum limits

For `N` fixed replicas of a key, write acknowledgment count `W` and read count `R` with `R + W > N` imply an intersection between those sets. This is a set-overlap fact, **not sufficient proof of linearizability**.

Reason about concurrent writes, reads racing writes, failed writes left on some replicas, membership changes, restoration from stale data, version comparison, and any sloppy quorum that uses substitute nodes outside the normal set. These cases can defeat the simple “one replica must know the latest value” story.

Read repair updates stale replicas encountered during reads. Anti-entropy compares replicas in the background, including infrequently read keys. Hinted handoff can retain a write for a temporarily unavailable owner, but delivery and convergence still depend on recovery. Monitor repair progress and stale-read behavior, not just healthy process counts.

## Resolve concurrency using domain meaning

Two operations are concurrent when neither causally depends on the other, even if their wall-clock timestamps differ. A version vector tracks per-replica progress and can distinguish a descendant from a concurrent sibling. It detects conflict; it does not choose the business-correct result.

Last-write-wins provides a simple winner but may discard accepted work. Wall-clock skew can make the winner particularly surprising. Use it only when that loss policy is acceptable and explicit. Do not use timestamp arbitration for scarce inventory or balances that require preserving all accepted operations.

Options include avoiding conflict through ownership, retaining siblings for user/domain resolution, or choosing a suitable CRDT with precisely defined concurrent semantics. A merge must account for deletions as well as additions. A set union can resurrect an item a user intended to remove. CRDT convergence is not proof that an arbitrary cross-record business invariant holds.

Local-first editing treats temporary divergence as part of normal operation. Design identity, synchronization, conflict semantics, authorization changes, and deletion propagation for disconnected clients. Offline acceptance and globally enforced uniqueness can conflict; make the chosen behavior visible.

## Worked example: a disappearing profile update

**Reconstructed from the source's stale-read scenario:** a user updates a profile, receives success, and is redirected to a follower that has not applied the change. The old value looks like a failed save.

An acceptable fix for a small application is to read that user's profile from the authoritative path. At larger scale, return a commit-position token with the write and require the subsequent read to meet that position. A fresh replica can serve it; a lagging replica must wait, reroute, or report the explicitly chosen freshness failure.

Test a long replication delay, replica failover, and switching from web to mobile. A token stored only in one browser cannot guarantee cross-device read-your-writes. Also test whether the acknowledged write survives leader loss; read routing cannot recover data that the failover policy permits losing.

## Required review evidence

Document acknowledgment semantics, replica placement/failure domains, promotion rules, read guarantees, conflict rules, and repair/rebuild behavior. Demonstrate the intended behavior with lag, partition, and recovery tests. For a multi-region proposal, show which operations cross a region and why their latency is compatible with the requirement.

Related: [sharding](ch07-sharding.md), [failure models](ch09-distributed-failures.md), [consistency and consensus](ch10-consistency-and-consensus.md).
