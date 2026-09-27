# Data models and query boundaries

## Core idea

Choose a representation that makes important relationships, queries, and invariants clear. Different models make different operations easy; neither JSON nor SQL decides the entire architecture.

## Select by relationships and access patterns

| Model | Useful shape | Inspect before choosing |
|---|---|---|
| Relational | Shared entities, many-to-one/many-to-many relationships, joins and constraints | Join patterns, indexes, transaction boundaries, schema evolution |
| Document | Mostly self-contained aggregates commonly loaded together | Aggregate growth, partial updates, references across documents, explicit validation |
| Graph | Variable-depth traversal and heterogeneous relationships | Traversal selectivity, degree skew, cycle handling, update semantics |
| DataFrame / array | Transformations, scientific or analytical work, matrix-oriented algorithms | Types, axes, missing values, ordering, memory layout, reproducibility |
| Event log plus projections | History and intent are central; multiple reproducible views matter | Ordering scope, event evolution, deterministic replay, external effects, deletion |

First test whether the existing database can express the needed model adequately. A relational database can store documents and graphs; a document system can offer joins. Specialized engines may improve ergonomics or performance, but add migration, operation, and synchronization costs. Verify capabilities in the chosen product/version before relying on them.

## Normalize facts; denormalize deliberately

Normalization avoids repeating independently changeable facts. Refer to a shared organization by stable identity when its name, logo, or other attributes must change consistently. Copying those attributes into every profile accelerates some reads but creates an update and inconsistency problem.

Denormalization is appropriate when its read benefits justify maintaining redundant representations. For each copy define the source, freshness contract, update mechanism, and repair method. Distinguish a **historical snapshot** from a redundant current fact: the shipping address on an accepted order may intentionally remain as agreed even after the customer's current address changes. Calling both copies “duplication” would erase domain meaning.

Joins express relationships; moving joins into application loops does not make their work disappear. It can add round trips, consistency windows, and N+1 access patterns. Compare a query plan, a batched application join, and a materialized view using the actual workload rather than banning joins categorically.

## Use document locality without hiding constraints

A document model is attractive when one aggregate is read and changed together and links to other aggregates are limited. Check whether nested elements need stable identities, independent access, concurrent updates, or references from elsewhere. A giant document can create hot records and expensive reads or rewrites; the details depend on the engine's storage behavior.

“Schemaless” generally means the schema is implicit in readers. **Schema-on-read** moves validation and compatibility responsibility to consumers; **schema-on-write** validates stored data earlier. Flexible data still needs a contract at the boundary. Versioned heterogeneous records can be valid; unexpected missing fields should not be converted into plausible default values.

For changes in document shape, distinguish absent, null, and empty values, specify supported versions, and use an explicit migration or boundary decoder. Do not spread ad hoc legacy-format repair throughout domain logic. See [encoding and evolution](ch05-encoding-and-evolution.md).

## Match query language to the question

Declarative queries state a result and let an optimizer choose an execution strategy. Relational algebra and SQL work well for joins and aggregation; graph languages such as Cypher and SPARQL express paths; Datalog expresses recursive rules. Recursive SQL can cover many graph queries, although a dedicated language may be more convenient for complex traversal.

A property graph models vertices and edges with properties. A triple store models subject–predicate–object facts. In either case specify identities, relationship direction, multiplicity, and traversal bounds. An unconstrained traversal over a dense graph can be operationally expensive even if the query is short.

GraphQL is an API query language, not a guarantee of graph storage, efficient joins, or transaction semantics. A convenient client query can cause substantial backend fan-out. Bound query cost and preserve authorization at the data boundary.

For analytics, a star schema centers a fact table at a precisely stated grain with dimensions describing those facts. A snowflake normalizes dimensions further. Never aggregate before defining the grain: joining an order-level amount to multiple line items can multiply revenue. Record units, time zones, and whether dimensions reflect current attributes or historical ones.

DataFrames bridge tabular transformations and array-oriented computation. A matrix conversion needs a stable row/column mapping and a missing-value policy. Missing feedback is not equivalent to a rating of zero. Feature order and encoding are part of a downstream model's contract.

## Event sourcing and CQRS

Event sourcing stores the accepted sequence of meaningful events as authoritative history. A projection derives current state. CQRS separates command handling from read representations; it does not inherently require event sourcing, separate services, or separate databases.

Apply event sourcing when retaining intent and replaying history delivers a concrete benefit that warrants the extra contracts:

1. Separate a proposed command, which may be rejected, from a recorded event describing an accepted fact.
2. Define the ordering scope needed for each invariant, often an entity or aggregate. Do not assume independent log partitions produce one global order.
3. Store enough information for deterministic interpretation, including relevant historical external inputs.
4. Version event schemas and projection code; retain a supported interpretation of old events.
5. Track projection progress and specify how a caller observes its own accepted change.
6. Separate rebuilding state from performing external effects such as sending email or charging a card.
7. Design deletion and retention across the log, snapshots, projections, and backups.

An immutable log does not make erroneous facts true. Corrections can be represented by compensating events, but compensation is a domain operation; it cannot undo every real-world effect. Auditability and replay do not automatically resolve privacy requirements.

## Worked example: reservation history

**Reconstructed from the source's reservation example:** a seat reservation changes from active to canceled, seat assignments disappear, and a refund record appears. Those row changes show what changed, but not necessarily why.

An accepted `ReservationCanceled` event can express the cause. A seat projection releases the relevant assignment, a reporting projection counts cancellations, and a refund workflow decides whether money must be returned. The event should identify the reservation and carry the facts needed for the agreed policy.

Replaying the event to rebuild a report must not issue a second refund. If a historical conversion rate affected the refund, reading today's rate during replay changes the outcome. Preserve the accepted rate or a stable historical reference. If two cancellations race, the command boundary must ensure the same reservation is not accepted as newly canceled twice; a log alone does not enforce that invariant.

Compare this with a conventional transaction updating current rows plus an audit record. If that simpler design satisfies the history and recovery requirements, choose it. Event sourcing is an option with costs, not a maturity level.

## Takeaways and review questions

- Which facts are independently owned, and which are snapshots or derived copies?
- Which relationships must remain consistent under concurrent changes?
- Can the selected model answer the main queries without uncontrolled fan-out?
- What schema assumptions do readers enforce, including historic data?
- Can projections be replayed deterministically without repeating external effects?
- Does the model preserve deletion, retention, and migration requirements?

Related: [storage](ch04-storage-and-indexes.md), [transactions](ch08-transactions.md), [streams](ch12-stream-processing.md), [responsible data](ch14-responsible-data-systems.md).
