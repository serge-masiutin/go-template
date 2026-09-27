# Batch processing and reproducible derived datasets

## Core idea

A batch job transforms a bounded input dataset into an explicit output. Immutable, identified inputs and controlled publication make reruns, debugging, and recovery easier than a collection of tasks that mutate production state as they go.

## Specify a job as a data contract

Record input dataset identities/snapshots, schemas, transformation version, parameters, output grain/schema, partitioning, completeness criteria, and publication behavior. Include dependencies such as lookup tables, model versions, locale, time zone, random seeds, and external enrichment snapshots when they affect results.

“Run the query again” is reproducible only if the inputs and interpretation are recoverable. A query against today's mutable customer table does not reproduce yesterday's customer attributes. An unversioned remote API is not a stable input.

Keep logical transformation separate from scheduling and storage. A small pipeline may fit one machine or a database query. Use distributed processing when data size, throughput, resource needs, or execution time justify the overhead.

## Learn from composable dataflow

Unix-style pipelines illustrate useful properties: each stage has an explicit input/output format, stages compose through a common interface, and intermediate results can be inspected. This does not mean every production pipeline should be a shell script. Preserve the explicit contracts while using suitable tooling.

A distributed job has at least three responsibilities:

- **Storage:** durable inputs, intermediate/checkpoint state as needed, and published output.
- **Compute/dataflow engine:** operators, partitioning, exchanges, and task recovery.
- **Orchestration:** resource allocation, dependencies, scheduling, retries, and execution status.

A workflow scheduler ordering jobs is different from an engine parallelizing one query. Avoid adding both layers when the existing runtime already meets the need.

## MapReduce, dataflow, and shuffle

MapReduce maps input records into keyed intermediate values, groups values by key through a shuffle, then reduces each group. The shuffle routes equal keys to a common partition and may sort within that partition. In modern engines an exchange/shuffle need not always sort; distinguish redistribution from a specific sort-based implementation.

General dataflow engines support directed graphs of operators, fuse compatible stages, pipeline work, and avoid unnecessary materialization. SQL and DataFrame APIs can expose enough structure for optimization. Inspect the physical plan: a simple expression can still trigger an expensive all-to-all shuffle.

Shuffles consume network bandwidth, serialization CPU, disk spill, and coordination. A skewed key can dominate one reducer even when the input bytes are distributed evenly. Examine per-task durations and partition sizes, not only total cluster utilization.

## Choose a join strategy

| Strategy | When it fits | Failure/cost to inspect |
|---|---|---|
| Partitioned sort-merge | Both large inputs can be partitioned/sorted compatibly | Shuffle and sort volume; hot groups |
| Partitioned hash join | Inputs can be distributed by the join key | Memory/spill pressure and skew |
| Broadcast/map-side join | One input is reliably small enough for every worker | Replicated memory/network use; an unexpectedly growing side |
| Co-partitioned join | Inputs already share a compatible key/partition contract | Hidden repartitioning if metadata or key interpretation differs |

Define multiplicity and missing-side behavior. A duplicate dimension key can multiply facts. A missing customer record can mean invalid input, an expected optional join, or a delayed dimension; choose explicitly. Do not silently drop unmatched rows when they affect the result's meaning.

For aggregation, use local partial aggregates where the operation permits combining them. An average should combine counts and sums, not take an unweighted average of partition averages. Floating-point accumulation order can change low-order bits; choose an appropriate numerical and reproducibility contract.

## Storage and resource behavior

Distributed filesystems and object stores expose different APIs. Do not assume directory rename, file locking, append, or atomic multi-object publication exists because a tool presents a filesystem-like path. Verify the exact backend. Data-local scheduling may help a filesystem cluster; disaggregated object storage trades locality for independent storage/compute scaling.

Plan object/file sizes, partition pruning, compression, and metadata overhead. Too many small files can make listing and scheduling expensive; oversized partitions can limit parallelism and recovery. Resource requests should include shuffle spill and intermediate space, not only source data size.

Long-running jobs must tolerate worker loss and retried or speculative tasks. The output protocol should ensure that only an accepted task/job attempt becomes visible. An engine's internal retry semantics do not automatically cover arbitrary external side effects.

## Publish completed output safely

Prefer building a new version in isolation, validating it, then publishing through a supported atomic manifest/version pointer or destination-specific commit protocol. Readers should either see the old complete version or the new complete version when that is the contract. Do not assume object-store rename supplies atomic publication.

Writing one row at a time into the live serving database from every worker can overload it, expose incomplete output, and duplicate effects on task retries. Alternatives include controlled bulk import, a versioned staging dataset, or a buffered stream consumed at a safe rate.

A stream buffers load but does not by itself hide partial batch output. If readers require an all-or-nothing dataset, attach a generation identity, stage results, verify completion, and activate that generation only after a reliable completion decision. Keep failed generations distinguishable and removable.

## Worked example: joining activity with account attributes

**Reconstructed source mechanism:** enrich a large activity log with account attributes for analysis.

1. Fix the activity interval and identify the account snapshot or temporal interpretation required by the question.
2. If both inputs are large, partition each by account identity and colocate equal keys through a shuffle.
3. Within a partition, join one account record to its activity records; reject or explicitly handle duplicate account identities and missing records.
4. Aggregate using a defined output grain, such as activity count per region per day.
5. Check source counts, unmatched counts, and conservation of total eligible activity.
6. Write a new output generation and publish after validation.

Using today's account region answers “activity by current region.” Joining the region valid at the activity's event time answers a different question. Store that choice in the job contract so replay and readers agree.

For a heavy account, inspect whether it creates a giant join group. Partial aggregation before the join may reduce volume if it preserves semantics; salting can distribute some work but adds a merge step. Do not use an optimization that changes the analytical question.

## ML and derived serving data

Feature engineering, training, evaluation, and batch inference are dataset transformations. Preserve feature definitions, input snapshots, model/code versions, seeds, and output lineage. Apply the same temporal reasoning to avoid using future information in a historical evaluation. A notebook can explore a transformation, but a retained production pipeline needs an explicit reproducible execution path.

Version recommendations, search indexes, and model outputs like other derived datasets. Rebuilding them should not resend notifications or mutate authoritative business facts. Serving freshness and correctness remain separate requirements.

## Validation

Test a worker crash, duplicate task attempt, skewed key, missing dimension, schema mismatch, full spill volume, and publication interruption. Compare reruns on the same controlled inputs; where exact byte identity is not expected, define the permitted numerical or ordering equivalence. Record partial-output visibility and restoration behavior.

Related: [models and grain](ch03-data-models.md), [streams](ch12-stream-processing.md), [integration and rebuilds](ch13-dataflow-and-correctness.md).
