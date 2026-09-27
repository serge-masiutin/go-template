# Consistency, ordering, and consensus

## Core idea

Name the observable guarantee an operation requires. Freshness, transaction isolation, causal order, unique identifiers, and agreement on a decision are related but distinct. A protocol or product label does not prove that every application path has those guarantees.

## Distinguish the guarantees

| Guarantee | What it says | What it does not say |
|---|---|---|
| Linearizability | Each operation appears to take effect at one instant between invocation and completion, respecting real-time order | Multiple independent reads/writes automatically form a transaction |
| Serializability | Committed transactions are equivalent to some serial execution | That serial order necessarily matches real-time completion order |
| Strict serializability | Serializable transactions also respect real-time precedence | External effects outside the transaction are included |
| Causal consistency | Causally related effects are observed in a compatible order | All concurrent operations have one real-time order |
| Eventual consistency | Replicas converge under the stated propagation/no-new-write conditions | A bound on lag or any particular intermediate read result |

Linearizability is a history property. If a write completes before a read begins, the read cannot return an earlier value unless another intervening write permits it. Reads overlapping a write have more freedom, but the entire history must fit one legal order; returning a new value and later reverting to an older one can violate it.

Use linearizable conditional operations for decisions needing one current winner, such as taking an exclusive resource. Plain linearizable reads followed by independent writes are still a race: the condition and mutation need the appropriate atomic boundary.

For multi-record decisions, analyze transaction isolation as well as per-key freshness. A “strongly consistent” setting with unspecified scope is not enough.

## Understand availability during a partition

If separated nodes cannot communicate, they cannot all independently accept operations requiring one linearizable shared decision while preserving that guarantee. Choose which operations wait/fail and which can proceed under a weaker, explicitly acceptable contract.

CAP is not “choose any two features.” Network partitions are a failure condition, CAP's consistency means linearizability, and its availability definition differs from a production uptime percentage. The theorem does not pick a database, describe every fault, or justify weakening unrelated invariants.

Outside partitions, communication required for coordination still affects latency. A multi-region architecture needs to show the quorum/leader paths and expected network cost. “Strong consistency with low latency” is incomplete without the operation, topology, and failure assumptions.

Do not weaken a payment or authorization invariant merely because a catalog read can tolerate staleness. Choose guarantees per operation, while accounting for dependencies among them.

## IDs are not a complete ordering protocol

Random unique identifiers can distinguish events but need not order them. Time-encoded identifiers can help locality or rough chronology without proving global causal order, real-time order, or gap-free sequences. Counter ranges assigned to different workers can be unique while appearing out of order.

Lamport clocks advance with local events and observed remote clocks. If event A causally precedes B, its Lamport timestamp is lower. The converse does not hold: a lower timestamp does not prove causality. Adding a node tie-breaker gives a total order over timestamps, but it does not tell a receiver that no lower-timestamp event remains in flight.

Hybrid logical clocks combine a physical-time component with logical advancement. They help timestamps remain near physical time while respecting tracked causal communication, but do not by themselves establish linearizability or a globally complete snapshot.

Vector clocks/version vectors retain multiple progress components and can distinguish concurrent histories at the cost of larger metadata. A causality representation detects a relationship; application merge rules still decide its meaning.

Ordering all known events is different from knowing which events are committed and safe to deliver. Do not implement global uniqueness or total-order delivery by sorting timestamps and assuming unseen proposals cannot exist.

## Consensus and replicated state machines

Consensus chooses an agreed value from proposals under a stated fault model. Its core properties include agreement, validity, and non-reversal of a decided result; progress requires additional conditions. A machine that always makes the same hardcoded decision does not satisfy the intended proposal/validity contract.

In replicated state-machine designs, nodes apply the same committed operations in the same order with deterministic interpretation. A shared ordered log is a useful abstraction for this. Random values, current time, or external lookups used during application must be captured or otherwise handled consistently.

Mature protocols such as Raft and Multi-Paxos use terms/epochs, persistent state, quorum intersection, and recovery rules to preserve decisions across leadership changes. Do not reimplement one from a sketch. The details of log matching, uncommitted entries, membership changes, snapshots, and reads are correctness-critical.

For a typical majority-based crash-fault cluster with `2f + 1` voting members, a reachable majority can tolerate up to `f` unavailable voters while preserving the protocol's assumptions. Failure-domain placement matters: five voters in one failed location do not provide geographic resilience. Different protocols/configurations have different quorum rules; the arithmetic is not universal for Byzantine or arbitrary quorum systems.

Consensus safety can survive a partition while progress stops without the needed quorum. FLP rules out unconditional deterministic termination in a fully asynchronous crash-fault model; practical systems rely on additional timing/progress assumptions rather than treating consensus as impossible.

## Reads and reconfiguration need proof too

A leader's local state may be stale if it no longer knows it is leader. A linearizable read requires the protocol's valid read mechanism, such as quorum confirmation/read-index or a correctly bounded lease with all its assumptions, plus application of the relevant committed state. Merely sending a read to the node labeled leader is insufficient.

Adding or removing voters changes which quorums intersect. Use the implementation's supported membership transition; independently editing lists can create disjoint decision groups. Force-promoting stale replicas may recover service while losing acknowledged history. If such disaster recovery is authorized, state the weakened guarantee and reconciliation obligations explicitly.

Consensus within each shard does not automatically supply a globally consistent cross-shard transaction or snapshot. Those require another suitable protocol spanning the relevant participants.

Coordination services are useful for small critical metadata, membership, ownership, and configuration. Keep bulk application traffic out unless the service is designed for it. Watches and notifications often signal that state changed; verify delivery and re-read semantics instead of treating every notification as an exactly-once business event.

## Worked example: privacy before publication

**Reconstructed source scenario:** a person changes an account from public to private on a laptop, waits for success, then uploads a photo from a phone. Permissions and photos live in different stores.

Independent logical timestamps can assign the photo a position that is inconsistent with the user's real-time sequence, because the photo store has not observed the privacy change. A reader might combine the new photo with old public permissions. Merely making both stores use hybrid clocks does not prove this cannot occur.

Identify the actual security requirement: publication/reading must use authorization state at the required freshness point. Options include an authoritative transaction/read boundary or an explicit dependency token enforced by the publication and read paths. The solution must also cover cross-device use, caches, and replicas.

A stale content count is a usability issue; stale permission state can disclose data. Test those semantics separately. Do not represent the requirement as “all reads eventually catch up.”

## Review evidence

Request the guarantee and scope, a short valid/invalid history, configured read/write paths, quorum/failure-domain assumptions, and behavior during leader loss and membership change. For existing software, verify the documented implementation and relevant configuration; do not infer guarantees solely from use of a consensus library.

Related: [replication](ch06-replication.md), [transactions](ch08-transactions.md), [failure assumptions](ch09-distributed-failures.md), [end-to-end correctness](ch13-dataflow-and-correctness.md).
