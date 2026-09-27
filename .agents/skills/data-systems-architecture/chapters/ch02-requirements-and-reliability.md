# Requirements, performance, reliability, and evolution

## Core idea

Turn “fast,” “reliable,” and “scalable” into observable requirements for specific workloads. Reliability is the ability to continue delivering the promised behavior despite specified faults; it is not simply the absence of component failures.

## Describe load before proposing capacity

Record request types, arrival rates, concurrency, read/write ratio, payload and result sizes, data volume and growth, key distribution, fan-out, burst duration, and background work. Include which clients or tenants create the most expensive operations.

The average user is often a poor design input. A social network's ordinary post and a celebrity post have different fan-out costs; a multi-tenant service may fit comfortably overall while one tenant overloads a partition. Measure distributions and correlations, not only global totals.

State the expected operating envelope and what happens beyond it: queueing, admission control, bounded degradation, or rejection. A queue changes when work happens; it does not remove the work or guarantee eventual catch-up.

## Separate latency concepts

- **Service time:** time spent doing the work at a component.
- **Queueing delay:** time waiting to begin or continue service.
- **Response time:** the delay observed by a caller, including relevant queues and network time.
- **Throughput:** completed work per unit time, with the success criterion specified.

Describe latency with distributions. Median says what a typical request sees; high percentiles expose slow experiences. A percentile is meaningful only with the workload, observation window, sample count, and measurement location. An average of instance-level p99 values is not the global p99; aggregate compatible distributions or raw observations.

Tail latency matters especially for fan-out. If one request needs every downstream answer, it waits for the slowest required dependency. For an illustrative independent model, if each of 100 calls is below a threshold with probability 0.99, all are below it with probability `0.99^100 ≈ 0.366`. Real dependencies can be correlated, so this calculation is a warning about amplification, not a production forecast.

Distinguish an **SLI** (the measurement), **SLO** (the target), and **SLA** (a service commitment with stated consequences). An SLO needs a numerator, denominator, window, and exclusions. Count timed-out or rejected eligible requests explicitly rather than removing inconvenient failures from the sample.

## Benchmark the system users actually exercise

1. Define a representative mix of reads, writes, expensive tenants, and background tasks.
2. Use production-shaped datasets and key distributions, with cold and warm conditions where relevant.
3. Generate arrivals without inadvertently reducing offered load when the system slows. A closed-loop test can hide queueing by waiting for each response before sending more work; report the arrival model.
4. Measure caller-visible latency, successful throughput, errors, queue depth, resource saturation, and cost together.
5. Run long enough to observe compaction, checkpointing, cache churn, and recovery effects.
6. Report the load at which the required behavior fails, not only peak successful throughput.

Do not extrapolate linearly without identifying the bottleneck. Increasing load while holding resources fixed asks how performance deteriorates; increasing resources while holding a performance target fixed asks what scaling costs. Both are useful and answer different questions.

## Design reliability around faults

A **fault** is a component deviating from its specification. A **failure** is the system failing to deliver its required service. State the tolerated fault model: one machine loss, one availability-zone loss, a network partition, an invalid deployment, corrupted storage, or a dependency outage. Redundancy does not guarantee independence.

Hardware faults may be isolated or correlated. Software faults, bad configuration, and overload commonly affect many identical replicas at once. Human error is a property to design around: safe interfaces, staged changes, reversible operations, useful feedback, and learning from incidents are stronger controls than assuming perfect operators.

For critical operations define:

- What success acknowledges and what survives each fault.
- Whether an unacknowledged operation may nevertheless have completed.
- A recovery time target and acceptable data-loss window, if applicable.
- Who detects the failure and what evidence drives recovery.
- How to distinguish restored availability from restored correctness.

Replication is not a backup against an accidental delete replicated to every node. A backup is not a recovery plan until restoration and reconciliation have been exercised.

## Make maintenance and change part of the architecture

**Operability:** make system state understandable, give operators bounded interventions, support controlled deployments and restoration, and provide actionable signals.

**Simplicity:** choose abstractions that hide appropriate implementation details while exposing relevant guarantees. A small API can conceal a very complex failure model; evaluate the reasoning burden for users and operators, not just code size.

**Evolvability:** keep compatibility and rollback practical. A migration that can be reversed safely has a different risk profile from one that destroys the old representation. Model data formats, mixed-version operation, and backfill behavior before changing a database or service boundary.

Shared-memory, shared-disk, and shared-nothing systems move coordination and resource-sharing costs to different places. Shared-nothing designs can scale partitions independently, but skew, network work, and distributed operations remain. A topology name is not evidence of linear scalability.

## Worked example: materialized feeds

**Reconstructed mechanism from the source, with illustrative smaller numbers:** assume 100,000 feed reads per second, 100 followed accounts per reader, 1,000 posts per second, and 100 followers per posting account on average.

- Computing a feed at read time may require roughly 10 million per-account lookups per second before batching and caching benefits.
- Updating materialized feeds at write time produces roughly 100,000 feed insertions per second, with cheap reads afterward.
- The latter trades write amplification and freshness lag for read efficiency. It also introduces a derivation pipeline and a rebuild obligation.
- One account with 5 million followers breaks the average. A hybrid can merge high-fan-out accounts at read time while materializing ordinary accounts.
- Accepting lag during a burst is an explicit product decision. Measure the oldest undelivered item, not only the queue's item count.

These counts are a workload model, not a database benchmark. Validate storage, batching, concurrency, and fan-out distributions before choosing infrastructure. Do not silently drop ordinary user-visible items to make the model fit; any lossy behavior needs its own product contract.

## Review and validation

Reject a performance proposal that lacks the workload and percentile target it addresses. Reject a reliability claim that does not specify faults and acknowledgment semantics. Request measured evidence for scaling claims and restoration evidence for recovery claims.

Useful checks include a burst beyond steady-state capacity, a hot tenant, a slow dependency, a full queue, a failed rollout, and a restore from a known backup. Record whether data invariants hold during and after each test, not just whether HTTP responses resume.

Related: [storage amplification](ch04-storage-and-indexes.md), [sharding](ch07-sharding.md), [partial failures](ch09-distributed-failures.md), [evolution](ch05-encoding-and-evolution.md).
