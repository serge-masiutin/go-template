# Sources and coverage

## Primary source

Martin Kleppmann and Chris Riccomini, *Designing Data-Intensive Applications*, second edition, O'Reilly Media, 2026, ISBN 978-1-098-11906-5. The supplied Russian PDF identifies the second-edition release as February 18, 2026 and contains 673 PDF pages, including front matter, references, glossary, and index.

This skill was created from the user's supplied `Kleppmann_Riccomini_DDIA_2nd_ru_mono.pdf` using the [book-to-skill converter](https://github.com/virgiliojr94/book-to-skill), with manual structural and technical synthesis. Preparation date: September 27, 2026.

Source SHA-256: `10ce36e97534ef4f01bdfdee30e774bbb195dbc907fb3ab6d0125f8954a0e45b`.

The capability name is intentionally independent of the book title. The content is an original practical synthesis of its mechanisms and trade-offs, with compact reconstructed or explicitly original examples. It is not a copy of the book or an exhaustive paraphrase of every page, reference, and product example. Attribution does not transfer rights in the source book. The PDF, extracted book text, and source illustrations are not included.

## Standalone operation

Every operational rule, worked example, decision table, glossary entry, and evaluation case needed by the skill is bundled. Source page numbers are provenance, not instructions to open the book. There is no runtime PDF path, extraction script, model download, API key, or network requirement.

For implementation-specific behavior, the target project's actual versions/configuration and current official documentation remain the source of truth. The skill explains how to reason about a decision; it does not certify all releases of any named product.

## Source coverage map

Printed page numbers differ from PDF page numbers by 24 in the main text. Ranges below include the chapter's trailing references; instructional content is mapped into the linked capability-oriented chapter.

| Source chapter | Printed pages | Bundled guidance |
|---|---|---|
| 1. Trade-offs in Data Systems Architecture | 1–32 | [Ownership, OLTP/OLAP, cloud and distribution choices](chapters/ch01-architecture-tradeoffs.md) |
| 2. Defining Nonfunctional Requirements | 33–64 | [Workloads, latency, reliability, scalability, maintainability](chapters/ch02-requirements-and-reliability.md) |
| 3. Data Models and Query Languages | 65–114 | [Relational/document/graph models, CQRS, event sourcing, arrays](chapters/ch03-data-models.md) |
| 4. Storage and Retrieval | 115–160 | [LSM/B-trees, amplification, analytical and search indexes](chapters/ch04-storage-and-indexes.md) |
| 5. Encoding and Evolution | 161–196 | [Compatibility, formats, RPC, messaging, durable execution](chapters/ch05-encoding-and-evolution.md) |
| 6. Replication | 197–250 | [Leader topologies, lag, conflicts, local-first, version vectors](chapters/ch06-replication.md) |
| 7. Sharding | 251–276 | [Keys, skew, movement, routing, local/global indexes](chapters/ch07-sharding.md) |
| 8. Transactions | 277–344 | [ACID, anomalies, serializability, 2PC, deduplication](chapters/ch08-transactions.md) |
| 9. The Trouble with Distributed Systems | 345–400 | [Partial failure, time, fencing, system models and verification](chapters/ch09-distributed-failures.md) |
| 10. Consistency and Consensus | 401–450 | [Linearizability, clocks, ordered logs, quorum/leadership limits](chapters/ch10-consistency-and-consensus.md) |
| 11. Batch Processing | 451–486 | [Bounded dataflow, shuffle/joins, orchestration, publication](chapters/ch11-batch-processing.md) |
| 12. Stream Processing | 487–538 | [Brokers, CDC, windows, joins, checkpoints and effect scope](chapters/ch12-stream-processing.md) |
| 13. The Philosophy of Dataflow Systems | 539–584 | [Integration, evolution, end-to-end integrity and audit](chapters/ch13-dataflow-and-correctness.md) |
| 14. Doing the Right Thing | 585–602 | [Accountability, feedback loops, privacy and data lifecycle](chapters/ch14-responsible-data-systems.md) |

## Extraction and interpretation limits

The installed converter used its `pdftotext` fallback because Docling was unavailable. The full text layer was extracted; analysis used the actual table of contents and bounded chapter/topic excerpts rather than loading the entire book into one context. The converter's numeric heading detector incorrectly counted 155 chapters, including repeated running headers; the verified 14-chapter table above replaces that count.

The PDF producer metadata identifies an AI-generated translation. The text includes broken words, displaced code, and inconsistent translations of terms such as consistency, integrity, compaction, and shuffle. Definitions were normalized to precise English technical concepts. Damaged extracted code was not copied as runnable code.

Selected rendered pages were inspected to verify schemas and mechanisms, including HNSW, writer/reader schema resolution, write skew, fencing, shuffle, processing-time artifacts, and dataflow transfer. Not every illustration was visually inspected. No claim that every image was read follows from the extractor's image-count metadata.

## Explicit implementation clarifications

Two external primary sources were checked during preparation to avoid carrying over overbroad/version-sensitive statements:

- [Temporal Activity Definition](https://docs.temporal.io/activity-definition): activity retries can repeat effects; idempotency is needed at the activity's effect boundary. Durable workflow replay is not a blanket exactly-once guarantee for arbitrary external calls. This clarification is included in chapter 5 and reflected in cases 4–5.
- [Apache Avro 1.12.0 specification](https://avro.apache.org/docs/1.12.0/specification/): reader defaults and union matching are version-specific. Chapter 5 avoids presenting an older first-union-branch formulation as a timeless rule. Verify the actual producer/consumer runtimes when implementing.

Additional operational applications are labeled as such: expand–migrate–contract, outbox/inbox packaging, versioned publication, and concrete review templates. They organize the source's underlying compatibility, transaction, log, and end-to-end reasoning for engineering use rather than claiming to be verbatim named frameworks from the authors.

The skill deliberately retains critical limits: quorum intersection is not a linearizability proof; snapshot isolation permits write skew; a lease without resource-side enforcement cannot stop stale writes; deduplication retention depends on replay horizons; checksums and consensus do not prove all application behavior; compensation requires an accepted business contract.
