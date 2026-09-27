# Storage engines, indexes, and query cost

## Core idea

An index trades extra work and storage at write time for less work at read time. Select physical structures using the actual query shapes, durability needs, and steady-state maintenance costs, not a product's peak benchmark.

## Understand the write and read paths

A simple append-only key/value log can locate the latest value through an in-memory hash index. Appending is straightforward, but retaining every obsolete value makes the log grow indefinitely. Compaction rewrites useful records and removes obsolete ones. Recovery needs a way to reconstruct the index and distinguish complete writes from partial/corrupt records.

Sorted string tables (**SSTables**) store keys in sorted immutable segments. A sparse index can locate a region without holding every key in memory. A **log-structured merge tree (LSM tree)** buffers writes in a sorted memory structure, flushes immutable tables, and merges tables in the background. A durability log is commonly needed before acknowledging buffered writes; inspect the engine's actual acknowledgment contract.

Reads may inspect multiple segments. Bloom filters can reject many absent point lookups; a conventional Bloom filter can have false positives but, under its correct operating assumptions, not false negatives. It does not generally solve arbitrary range scans. Range reads merge ordered streams and reconcile versions and tombstones.

A **B-tree** maintains sorted keys in a hierarchy of pages. Search follows a short path; range scans use key order. Updates may modify pages and split them. Write-ahead logging or another crash-consistency mechanism protects multi-page changes. A mutable page structure does not itself provide transactions or replication guarantees.

| Dimension | B-tree tendency | LSM-tree tendency |
|---|---|---|
| Point/range reads | Predictable page traversal; efficient ordered access | May examine/merge several sorted runs |
| Writes | Page updates, sometimes scattered | Buffered writes and larger sequential flushes |
| Maintenance | Page management, splits, engine-specific cleanup | Compaction consumes sustained I/O and CPU |
| Space | Page utilization and fragmentation matter | Old versions, tombstones, and temporary merge space matter |
| Latency risk | Contention, cache misses, page work | Compaction debt, flush pressure, stalls |

These are explanatory tendencies. Modern engines mix techniques; SSD behavior, caching, compaction policy, concurrency, and workload can reverse a simple comparison.

## Account for amplification and headroom

- **Write amplification:** physical work or bytes written exceed the logical write; state which layer is measured. WAL, indexes, compaction, replication, and device garbage collection can each contribute.
- **Read amplification:** answering a logical read requires extra pages, runs, records, or remote requests.
- **Space amplification:** physical storage exceeds live logical data due to redundant structures, obsolete versions, free space, or temporary maintenance.

Measure after the engine reaches sustained maintenance equilibrium. A short test can accumulate compaction debt while reporting excellent ingestion. Include disk headroom during merges and recovery; a nearly full device can prevent the maintenance needed to reclaim space.

A tombstone represents a deletion until older copies cannot reappear. Removing it too early can resurrect data from an old segment or replica. Garbage collection interacts with replication, snapshots, and retention; it is not merely a local storage optimization.

In-memory systems still need explicit durability and restart behavior. They can be fast because of simpler layouts and reduced I/O, but volatile memory alone does not preserve accepted changes after power loss. Identify logging, replication, checkpoint, and restoration guarantees.

## Choose indexes for queries

1. List predicates, joins, required ordering, result size, and frequency for the important queries.
2. Inspect representative query plans and data distribution.
3. Choose key order, included values, and filtering support that reduce the measured work.
4. Estimate and measure write, storage, and maintenance costs of each added index.
5. Test with realistic selectivity, cold data, and skew; recheck as the distribution changes.

A secondary index maps another attribute to record identities or stored values. A clustered layout colocates rows with an index; a covering index stores enough values to answer a query without a further record lookup. Such layouts can improve locality while duplicating more data.

A composite ordered index is not symmetric in its columns. For an illustrative `(tenant_id, created_at)` index, one tenant's date range is naturally contiguous; a global date range may require a different strategy. Verify the optimizer's behavior rather than promising that every database uses only a simple left-prefix rule.

Spatial queries over multiple dimensions may benefit from specialized structures such as R-trees. Do not confuse “multiple columns in a sort order” with a multidimensional access method.

## Analytical storage and execution

Column-oriented storage puts values of one column together. A query reading a few columns across many rows can avoid unrelated data, compress similar values, and use vectorized operations. Ordering, dictionaries, run-length encodings, bitmaps, and other compression methods can reduce work. The benefit depends on the query and data, not merely the file extension.

**Vectorized execution** processes batches of values with efficient CPU operations; it is distinct from semantic vector search. Query compilation reduces interpretation overhead. These techniques matter when CPU execution becomes a bottleneck after unnecessary I/O is removed.

Materialized aggregates and data cubes move computation ahead of reads. Specify supported dimensions, update semantics, freshness, and rebuild strategy. Preaggregating every combination can cost more than computing infrequent questions on demand.

Separating compute from object storage permits independent capacity choices but introduces remote access and request overhead. Partition pruning, metadata, caching, and appropriately sized files affect whether a query reads a useful subset or an expensive collection of tiny objects.

## Text and semantic search

An inverted index maps terms to documents or postings. Tokenization, normalization, stemming, language, and ranking are part of the search contract. N-gram indexes support substring-style matching at additional storage cost; typo tolerance is a separate behavior from semantic similarity.

An embedding maps content into a numerical representation; similarity search compares a query representation with stored representations. Preserve the embedding model/version, dimensions, distance metric, preprocessing, and document identity. Changing these can require re-embedding and rebuilding an index.

| Vector approach | Mechanism | Decision trade-off |
|---|---|---|
| Flat exact scan | Compare the query against every candidate | A useful correctness baseline; work grows with candidate count |
| IVF | Search selected clusters of vectors | More probes generally trade more work for better recall |
| HNSW | Navigate a hierarchy of proximity graphs | Approximate retrieval with construction, memory, and search trade-offs |

Approximate neighbors can miss relevant candidates. Evaluate recall and application relevance against a reference set, with filters and updates included. Nearest vectors do not prove factual truth or authorization. A vector extension in an existing database may suffice; introduce a separate vector service only for a demonstrated requirement.

## Worked example: an apparently slow orders query

**Original application example:** a page retrieves the newest 50 orders for one tenant. The existing database contains millions of orders across tenants and scans/sorts far more rows than it returns.

Start with a query plan and a representative large tenant. A composite tenant/time access path may make a bounded ordered range read possible. Check deterministic pagination with a stable tie-breaker, such as order identity. If returned columns require expensive row lookups, consider covering frequently read small values, weighing larger indexes and write cost.

Measure the revised page query and write workload. Do not add a cache first: that would introduce invalidation and freshness questions without addressing the underlying query. Conversely, an occasional report scanning all tenants may be better served by a columnar analytical copy than by adding more OLTP indexes.

Validate the plan under skew, concurrent writes, and realistic history. Treat any stated latency improvement as unproven until measured.

Related: [workloads](ch02-requirements-and-reliability.md), [sharding indexes](ch07-sharding.md), [derived state](ch13-dataflow-and-correctness.md).
