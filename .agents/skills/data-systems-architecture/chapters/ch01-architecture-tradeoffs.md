# Architecture trade-offs and data ownership

## Core idea

Choose an architecture by the workload, guarantees, and operational constraints it must satisfy. A database, cloud service, or distributed topology is a means to those ends. This chapter concerns system boundaries; it does not prescribe an application class hierarchy.

## Start with the workload

Distinguish operational work from analytics before selecting a storage engine.

| Question | Operational / OLTP tendency | Analytical / OLAP tendency |
|---|---|---|
| Who consumes results? | A user or another service completing a task | An analyst, report, model, or downstream application |
| Read shape | A few records selected by an identifier or selective predicate | Large scans, joins, grouping, aggregations |
| Write shape | Small, concurrent changes with business invariants | Bulk ingestion and transformations, sometimes continuous updates |
| Important performance dimension | Per-operation latency and concurrent throughput | Scan throughput, resource efficiency, time to result |
| Typical representation | Current business state | Historical facts and derived views |

These are workload tendencies, not product categories. An operational application can run an analytical query; an analytical result can return to an operational product. Measure the interference before putting both workloads on the same resources. A dedicated analytical copy adds delay and a pipeline to operate; it is justified when those costs buy necessary isolation or query capability.

Use a warehouse when curated, queryable analytical structures are central. Use object-based lake storage when heterogeneous source formats, inexpensive retention, and multiple processing consumers matter. A lake does not remove the need for ownership, catalogs, schemas, access control, or reproducible transformation. Retaining everything forever is not a prerequisite for analytics.

## Map authority and derivation

For every dataset, write down:

1. The fact it represents and the component authorized to change that fact.
2. Whether it is a **system of record** or **derived state**.
3. Its upstream inputs, transformation version, and update mechanism.
4. What lag or inconsistency a consumer can observe.
5. How to recover it after corruption or loss.

A system of record holds the authoritative version for its stated scope. A derived dataset is reproducible from retained authoritative inputs and the relevant transformation. Products do not determine these roles: the same relational database can hold both authoritative orders and derived order summaries.

Treat the ability to rebuild as an obligation to demonstrate. A search index is not fully reproducible if historical source data has expired, the tokenizer version is unknown, or an external enrichment changes every time it is called. Either retain/version those inputs, declare an explicitly lossy rebuild contract, or stop calling the dataset disposable.

When two stores disagree, an ownership map should explain which fact wins. If independent writers legitimately create competing facts, define a merge or conflict-resolution contract instead of pretending there is one obvious authoritative value.

## Compare deployment options by responsibilities

Managed services can transfer hardware maintenance, upgrades, and some recovery work to a provider. Application correctness, configuration, capacity limits, access policy, dependency failure, and restoration remain design responsibilities. Self-hosting gives more control and diagnostic access while assigning more operational work to the team.

Compare alternatives on actual constraints:

- Required capabilities and guarantees, including failures during upgrades.
- Team experience, on-call coverage, recovery procedures, and incident visibility.
- Storage, compute, request, network-transfer, and support costs under the expected load.
- Data location, access, retention, export, and deletion obligations.
- Portability of data and semantics, not just apparent API compatibility.
- Exit procedure: export bandwidth, downtime, migration verification, and retained rollback state.

Cloud-native architectures often separate durable storage from compute. This can allow independent scaling and replacement of workers, but a remote storage dependency adds network behavior, access costs, and new failure boundaries. Do not infer those properties merely from the label “cloud.” Identify the actual data path.

## Decide whether distribution is needed

First evaluate a single-node or single-database design against measured needs. Distribution may be necessary for geographic proximity, fault tolerance, datasets or workloads exceeding one machine, or independent ownership. Each reason demands a different topology; “we may grow” does not specify one.

Replication copies data and introduces questions about freshness and failover. Sharding splits data and introduces routing, skew, and cross-shard operations. Splitting a process into services introduces partial failures even if the database stays centralized. Serverless execution changes provisioning and scaling, but persistent state, concurrency, retries, and resource limits still need contracts.

Use the smallest topology that meets the requirements with a credible growth path. Record what observation would invalidate it, such as storage headroom, a measured throughput ceiling, or an explicit regional recovery requirement. Do not assert a universal traffic threshold for adding services or shards.

## Worked example: operational orders and reporting

**Scenario, reconstructed as a practical application:** a small marketplace needs transactional order placement, catalog search, and daily revenue reports. Current traffic fits one relational database.

1. Make orders, order lines, and accepted payments authoritative records in that database. Define transaction boundaries around the relevant invariants.
2. Start with database-supported search if its semantics and measured latency suffice. A dedicated index becomes an option when ranking or scale requires it.
3. Run reports on controlled snapshots or a separate analytical copy if scans interfere with order placement. Declare a report's “as of” time.
4. If adding a search engine, send committed changes through a recoverable derivation mechanism. Store progress and verify rebuilds; do not treat two independent writes as atomic.
5. On an index outage, preserve order correctness and expose the chosen search failure behavior. A degraded search experience may be acceptable; a silently wrong payment balance is not.
6. Include time to rebuild indexes, restore authoritative data, and reconcile projections in the operating cost comparison.

The initial recommendation is deliberately modest because no requirement justifies a distributed write path. A future need can change the decision without changing the ownership of existing facts.

## Failure modes and review questions

- **Technology-first diagrams:** boxes show products but no authoritative facts, update direction, or failure behavior. Add those contracts before discussing vendors.
- **Managed means automatic correctness:** a provider's durable database cannot make an application-side dual write atomic. Identify which guarantee ends at each boundary.
- **All copies are sources of truth:** independent edits to an index, cache, and database create unresolved ownership. Assign authority or explicitly model conflicts.
- **Unbounded retention:** storage cost omits breach impact, unnecessary exposure, and deletion work. Retain data for a stated purpose and duration.
- **Premature service splits:** organizational or deployment convenience is mistaken for evidence that transactions can be separated safely. Inspect invariants crossing the proposed boundary.

## Decision artifact

Produce a workload table, an authority/derivation map, two plausible options, the chosen trade-off, operational ownership, and the condition for revisiting the choice. List unmeasured assumptions rather than inventing capacity numbers.

Related: [requirements](ch02-requirements-and-reliability.md), [models](ch03-data-models.md), [dataflow integration](ch13-dataflow-and-correctness.md), [responsible data use](ch14-responsible-data-systems.md).
