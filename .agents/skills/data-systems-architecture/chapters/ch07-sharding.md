# Sharding, skew, routing, and rebalancing

## Core idea

Sharding divides ownership of a dataset so multiple resources can share its storage and work. The partition key determines which operations stay local, which fan out, and which require coordination. Replication and sharding are independent dimensions.

## Establish the reason to shard

Identify the present or credibly forecast limit: dataset size, write throughput, per-node resource pressure, isolation of large tenants, or another explicit requirement. Check query/index improvements, resource sizing, and read replicas first when they address that limit with less complexity.

Do not equate table partitioning inside one database with independent distributed shards. Internal partitions can help pruning and maintenance without changing the failure and transaction boundary. State which meaning is intended.

List the dominant queries and invariants before selecting a key:

- Which operations know the partition key?
- Which entities must be changed together?
- Which reads scan a range or aggregate many entities?
- How skewed are data size, read rate, and write rate independently?
- Can the largest tenant/key outgrow one partition?
- Which key values or access patterns reveal sensitive information?

## Range and hash partitioning

**Range partitioning** gives a shard an interval of ordered keys. It preserves locality and makes suitable range queries efficient. Monotonically increasing timestamps or IDs can concentrate all new writes at one end. Composite keys can spread writes across a meaningful leading dimension while preserving a useful local order.

**Hash partitioning** distributes distinct keys more evenly but loses original key order. Queries without the hashed key may need many shards. Use a stable hash contract across runtimes and deployments; a language's process-randomized hash is unsuitable for durable routing.

Neither approach automatically distributes workload evenly. One very popular key still maps to one shard. A tenant key may colocate transactions conveniently while concentrating all work of a large tenant.

## Separate logical partitions from machines

Using `hash(key) % number_of_servers` remaps many keys when the server count changes. Prefer a managed assignment of logical partitions to nodes, ranges that can split/move, or a suitable consistent/rendezvous hashing scheme. Evaluate movement cost and the system's actual routing guarantees; “consistent” hashing has no relation to transaction consistency.

More logical partitions provide placement flexibility but add metadata, file, scheduling, and recovery overhead. A fixed partition count may become a ceiling; dynamic splitting adds data movement and coordination. Choose based on growth, workload, and the engine's operational support rather than copying a universal partition count.

## Handle hot keys explicitly

Different hot spots need different remedies:

| Symptom | Candidate response | New cost or constraint |
|---|---|---|
| Hot reads on immutable or staleness-tolerant data | Replicas or caching | Freshness/invalidation and stampede behavior |
| Hot mutable counter or append stream | Subdivide independent writes, aggregate later | Read fan-in and loss of a simple single-key invariant |
| One large tenant | Dedicated placement or further tenant subdivision | Routing and cross-subpartition transactions |
| Recent-time write concentration | Composite key or distribution across writers | More expensive global time-range reads |
| Temporary viral key | Targeted adaptive placement or subdivision | Transition protocol and eventual scale-down |

Salting a key spreads writes only when operations can be split safely. Reading an aggregate may now require all salted keys. A uniqueness or balance invariant cannot simply be scattered and reconstructed after accepting conflicting updates.

## Design routing and movement together

Routing can live in a client library, a routing tier, or server-side forwarding. Each option needs a source of partition ownership and a way to handle stale metadata. A client should not guess ownership from an obsolete local map and issue independent authoritative writes to both locations.

For a move, specify the system-supported protocol for snapshot/copy, concurrent change catch-up, ownership transfer, stale-writer rejection, and old-copy retirement. Preserve a monotonic ownership epoch or equivalent guard where the protocol requires it. Copying bytes alone is not a safe ownership transition.

Operationally:

1. Budget transfer bandwidth, disk space, and CPU alongside foreground load.
2. Limit concurrent moves and pause them when latency or replication health deteriorates.
3. Make progress observable and restartable.
4. Verify data counts/invariants and the catch-up position before cutover.
5. Keep a bounded rollback/recovery path consistent with the new write authority.

Automatic rebalancing can create a feedback loop: overload looks like node failure, triggering movement that increases overload. Prefer deliberate, bounded automation with visible health conditions.

## Secondary indexes change the cost model

**Local / document-partitioned indexes:** each shard indexes its own records. A write updates local indexes efficiently, but a lookup without the primary partition key may scatter to all shards and gather results. More shards may increase storage capacity while making every such query more expensive.

**Global / term-partitioned indexes:** index entries are partitioned by the indexed value independently of base records. A selective lookup can reach the relevant index partition directly, then fetch base records. Writes may touch several partitions, and multiple search predicates can require distributed intersections.

Define whether index updates are atomic with the base record or asynchronous. An asynchronous global index may temporarily omit a newly written record; it must not be used as the authoritative uniqueness check unless the system provides an appropriate stronger mechanism.

Cross-shard joins, transactions, and aggregations are possible, but require coordination and network work. Co-location can reduce that cost, at the expense of distribution flexibility. Describe the trade-off for actual query paths rather than declaring joins impossible.

## Worked example: a tenant-centric service

**Original application example:** almost every order operation belongs to one tenant, but an administrator searches all orders by customer email.

Partitioning by tenant keeps normal order transactions local. Hashing tenant identity spreads ordinary tenants. It does not split the one tenant producing half the traffic; that tenant needs separate placement or a finer partition boundary with an explicit transaction plan.

The global email search cannot be assumed cheap with local indexes. Options are scatter/gather with a bounded admin workload, a global index with stated freshness, or a derived search service. Choose using frequency, latency, and update guarantees. If the search is used to enforce a global rule, asynchronous indexing alone is insufficient.

During tenant migration, route writes through one current ownership decision, catch up the target, and cut over with the database's supported protocol. A retry carrying old routing information must be rejected or redirected safely, not accepted by an abandoned owner.

## Review evidence

Request per-shard load distributions, the worst key/tenant, cross-shard query traces, movement headroom, and failure tests during cutover. Check the largest partition's restore time: overall aggregate capacity does not make that bottleneck disappear.

Related: [requirements](ch02-requirements-and-reliability.md), [index structures](ch04-storage-and-indexes.md), [transactions](ch08-transactions.md), [coordination](ch10-consistency-and-consensus.md).
