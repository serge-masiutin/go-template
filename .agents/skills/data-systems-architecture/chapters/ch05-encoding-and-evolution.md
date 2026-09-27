# Encoding, compatibility, RPC, and durable workflows

## Core idea

Persisted bytes and network messages outlive the code that created them. A format is a contract among writers and readers of different versions, including rollback versions, delayed consumers, archives, and recovery tools.

## Define compatibility direction explicitly

- **Backward compatibility:** a newer reader understands data from an older writer.
- **Forward compatibility:** an older reader understands data from a newer writer.

Do not say only “compatible.” Name the writer, reader, format/version, and behavior being preserved. Parsing successfully is weaker than interpreting the same business meaning. A field whose unit changes from cents to dollars can remain syntactically valid while corrupting calculations.

Inventory data flows through databases, synchronous APIs, queues, files, and workflows. Database rows may persist for years; mobile clients and consumers may upgrade late; a queued event may outlive several deployments. A rolling upgrade and a rollback can expose combinations never seen in a same-version unit test.

## Compare representations

| Representation | Benefits | Contract risks |
|---|---|---|
| Language-specific object encoding | Convenient local objects | Runtime coupling, unsafe deserialization, hidden version assumptions |
| JSON/XML and textual formats | Readability and broad interoperability | Numeric precision, optionality, bytes, dates, and schema enforcement differ |
| Protocol Buffers | Tagged fields, compact representation, generated interfaces | Field-number reuse, presence/default semantics, enum handling, format conversion |
| Avro | Writer/reader schema resolution, compact records, useful generated schemas | Writer-schema availability, resolution rules, defaults and incompatible type changes |

Treat external serialized data as data. Avoid instantiating arbitrary runtime objects from untrusted bytes. Validate sizes, shapes, required fields, and supported versions at the boundary. Forward-compatible unknown fields may be allowed by a documented contract; malformed required data should still fail clearly.

JSON numeric handling varies across runtimes. Identifiers, exact money amounts, decimal scale, timestamps, time zones, and binary fields need explicit representations. Test the full producer-to-consumer path, including proxies and reserialization, rather than assuming a format name guarantees fidelity.

## Evolve schemas deliberately

With tagged formats, field identifiers become durable meaning. Do not reuse a retired tag for another fact. A rename, type change, or move into an envelope may be wire-compatible in one encoding and incompatible in another. Unknown-field preservation can also change when converting to a textual representation; verify the selected toolchain.

Avro decoding uses the writer's schema and resolves into a reader's schema. Missing reader fields can be supplied by **schema-declared** defaults; this is intentional compatibility logic, not permission to invent values for malformed inputs. Reader-only required fields without a valid default cannot simply be assumed present. Store or reliably reference the writer schema for each encoded record/file.

Implementation clarification: union-default rules depend on the Avro specification/runtime version. Avro 1.12.0 describes matching the default to the first matching branch; older specifications used a first-branch restriction. Keep `null` first with a matching null default when that is the intended cross-version contract, and validate with the actual runtimes. Do not generalize one edition's example to every release. See the [Avro specification](https://avro.apache.org/docs/1.12.0/specification/).

## Use an expand–migrate–contract rollout

This is an operational application of the source's mixed-version compatibility reasoning:

1. **Expand:** deploy readers and schema support that understand old and new representations without requiring every writer to change at once.
2. **Migrate:** enable new writes and backfill historical data through a bounded, restartable process. Define which concurrent update wins.
3. **Verify:** compare representative records, invariants, consumer progress, and compatibility fixtures. Observe lagging clients and old queued data.
4. **Switch:** move reads to the new representation with a specific rollback plan.
5. **Contract:** remove old fields or interpretations only after the supported old readers/writers, retained events, and rollback period no longer require them.

An additive database migration can still take locks or rewrite data depending on the operation and engine version. Verify execution behavior before calling it online. Backfills must not overwrite a fresh user edit with an old snapshot value.

An old service reading a record, modifying one known field, and writing the entire record back can erase fields it does not understand. Use a preservation strategy consistent with the format: patch only owned fields, preserve unknown data, or remove that old writer from the allowed deployment matrix.

## RPC is a partial-failure boundary

A remote call differs from a local function: requests and responses may be lost, the remote side may execute after the caller times out, and a retry may repeat an effect. Define deadlines, cancellation behavior, idempotency, error classification, and retry ownership.

REST or RPC syntax does not remove these issues. Generated stubs can make a remote call look local while preserving all its failure uncertainty. Bound payloads and concurrency; avoid retries at several layers multiplying load. A timeout says the caller lacks a timely result, not that the operation did not happen.

Use asynchronous messaging when buffering, independent processing, or delayed delivery is part of the requirement. Specify ordering scope, retention, acknowledgment, poison-message handling, and duplicate behavior. A queue changes coupling and failure handling; it does not create a transaction spanning arbitrary services.

## Durable execution has an effect boundary

A durable workflow records progress so execution can resume or replay after failures. Workflow state and recorded results can prevent completed logical steps from being scheduled again during replay. This is useful for long-running coordination, timeouts, and recovery.

An external activity can still complete before its result is durably recorded. Retrying that activity may repeat a payment or email. Temporal's documentation explicitly calls for idempotent Activities because retries can repeat execution; a workflow engine alone does not make arbitrary external effects exactly once. See [Activity Definition](https://docs.temporal.io/activity-definition).

Separate deterministic workflow decisions from effects. Give each effect a stable logical operation identity, use the destination's idempotency mechanism where available, and provide reconciliation for unknown outcomes. Version long-running workflow behavior and preserve the supported replay contract. Choose durable orchestration for a demonstrated recovery need, not merely to wrap a few local calls.

## Worked example: adding a shipping preference

**Original application example:** order schema v1 has no delivery preference; v2 adds an optional preference where absence intentionally means “standard handling.”

Test the four writer/reader pairs: v1→v1, v1→v2, v2→v1, and v2→v2. Verify that the v2 reader applies only the documented default to v1 data; that v1 readers can tolerate the additive field; and that a v1 writer cannot erase a preference by reading and replacing a v2 record. Add a fixture with an unrecognized preference value to distinguish extensibility from invalid data.

If v2 instead introduces a mandatory customs declaration, defaulting absence to an empty declaration is dishonest. Support an explicit legacy state, collect the missing information, or restrict the new operation until requirements are met. Deploying a new decoder cannot create facts that were never stored.

For a workflow calling a shipping provider, persist an operation ID before issuing the call. A lost response enters an unknown-outcome branch resolved by lookup or safe retry under the provider's contract; it must not create a new shipment ID on every attempt.

Related: [data models](ch03-data-models.md), [failure uncertainty](ch09-distributed-failures.md), [stream replay](ch12-stream-processing.md).
