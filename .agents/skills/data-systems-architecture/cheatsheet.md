# Architecture decision cheatsheet

| If you observe… | Do this, because… |
|---|---|
| A small workload fitting one database | Start with its queries, constraints, and transactions; extra services add failure boundaries. |
| An unspecified “scale problem” | Measure load, skew, query cost, and tail latency before choosing distribution. |
| An expensive selective read | Inspect plan/index/locality before adding cache invalidation. |
| Large scans disrupting transactional work | Compare isolated analytical processing against its pipeline/freshness cost. |
| Many shared references across documents | Evaluate normalization/joins; copying facts creates update obligations. |
| A check followed by a write | Identify all concurrent writers and protect the premise at an atomic boundary. |
| A predicate returning no rows | Row locks may protect nothing; consider constraints, serializable isolation, or a shared guard. |
| A timeout after sending an effect | Preserve the logical ID and resolve unknown outcome; do not assume rejection. |
| A database write plus a separate publish | Use a durable atomic intent boundary such as an outbox; expect redelivery. |
| A user cannot see an accepted update | Provide authoritative/progress-aware reads; arbitrary sleeps do not prove freshness. |
| A worker can outlive its lease | Enforce fencing/conditional ownership at the resource. |
| One hot key despite balanced shards | Split only safely separable work or change placement; more hashes do not split one key. |
| Search without the partition key | Price scatter/gather versus a global index and its update guarantees. |
| A new projection or online migration | Coordinate snapshot and live-change positions; validate a new generation before switching. |
| Historical replay changes enrichment | Preserve historical inputs or declare the changed question explicitly. |
| Events arrive after window completion | Define correction/retraction or observable rejection, plus retained state. |
| Strict global claims during a partition | State which operations stop; independent acceptance cannot preserve one current winner. |
| Sensitive facts in immutable history | Design deletion, retained evidence, and replay behavior before collecting them. |

## Choose the boundary

1. Can one supported atomic operation/constraint enforce the invariant? Use it.
2. Does one database transaction cover the decision? Choose the required isolation and retry contract.
3. Are other stores only derived copies? Use recoverable propagation with explicit freshness.
4. Is cross-system atomic visibility genuinely required? Evaluate a supported coordination protocol and its cost.
5. Is delayed completion/compensation acceptable? Define the business state machine and exposure; do not assume consent to weakening the invariant.

## Invalid shortcuts

| Claim | Missing proof |
|---|---|
| `R + W > N`, therefore linearizable | Concurrency, membership, version resolution, failed writes, and protocol details |
| Snapshot isolation prevents every race | Write skew and predicate-based decisions |
| A timestamp proves global order | Clock uncertainty and unobserved/causal events |
| A committed transaction means one user action | End-to-end retry identity |
| Exactly-once stream means exactly-once email | Destination effect and acknowledgment boundary |
| CRDT means no business conflicts | Domain invariant and merge semantics |
| Replicated means backed up | Recovery from replicated corruption/deletion |
| Backup exists, therefore recovery works | Tested restoration and reconciliation |
| Vector similarity means relevant and allowed | Retrieval evaluation, freshness, and authorization |

Evidence order: **requirements → source of truth → invariant → failure history → mechanism → test/measurement → rollout/recovery**. There is no universal QPS threshold, database winner, or mandatory microservice count.
