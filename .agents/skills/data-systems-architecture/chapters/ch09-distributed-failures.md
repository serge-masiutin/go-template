# Partial failures, clocks, leases, and verification

## Core idea

A distributed system must make decisions while lacking complete knowledge of other nodes. Network delays, process pauses, and lost responses create uncertainty that cannot be removed by a generous timeout or a successful happy-path test.

## Represent uncertainty in the operation contract

After sending a request, the caller may not know whether the destination never received it, is still processing, completed it but lost the response, or failed after a partial effect. A timeout distinguishes “no timely answer” from “answer”; it does not distinguish all those histories.

For an effectful operation, model at least known success, known rejection/failure where the protocol establishes it, and unknown outcome. Resolve uncertainty through a stable operation identity, result lookup, safe retry under an idempotency contract, or reconciliation. Do not turn uncertainty into an invented failed or successful result.

TCP provides a byte stream for a connection; it does not promise a business operation executes once, a response arrives by a deadline, or the peer application remains healthy. Connections, proxies, and clients can time out at different points.

## Treat failure detection as suspicion

A slow node and a failed node can look identical from another node. Queueing, scheduling delays, garbage collection, disk stalls, overloaded network buffers, virtualization, and paused processes can all delay an otherwise correct participant.

Short timeouts detect some failures quickly but increase false suspicions and unnecessary failovers. Long timeouts reduce false positives but delay recovery and consume resources. Choose deadlines using observed distributions and the caller's budget; separate connection, request, and overall workflow deadlines where those distinctions matter.

Retry only classified conditions, with bounded attempts or elapsed-time budgets, backoff/jitter where appropriate, and a stable operation identity. Avoid multiplying retries across layers. A timeout retry against an overloaded service can prolong the overload. Cancellation is often a request to stop, not proof the destination stopped before committing an effect.

## Use the right clock

| Clock/ordering mechanism | Suitable purpose | What it does not establish |
|---|---|---|
| Wall-clock / time-of-day | Human timestamps and calendar interpretation | Reliable elapsed duration or global causal order |
| Monotonic clock | Local durations and timeout measurement | Comparable timestamps across machines or restarts |
| Logical clock | Event ordering based on communication | Accurate physical time or knowledge of unseen events |

Wall clocks can drift or jump. NTP reduces error but does not make timestamps infallible. Clock precision is not clock accuracy. When a correctness protocol depends on a bounded clock uncertainty, name the mechanism and what happens when the bound cannot be maintained.

Do not compare monotonic readings from unrelated processes/machines as if they share an epoch. Even local elapsed-time APIs have documented behavior around suspend and reboot; use the runtime's supported contract.

## Leases need resource-side enforcement

A lease expires so another worker can take over. The old worker can be paused past expiry and resume believing it still owns the resource. A delayed packet sent before expiry can also arrive after a new owner starts. Checking a lease immediately before writing leaves another possible pause between the check and the write.

For correctness-critical ownership, use a resource-supported fencing or conditional-write protocol:

1. The coordination mechanism issues a monotonically increasing ownership epoch/token with appropriate persistence and ordering guarantees.
2. Each protected write includes the epoch.
3. The resource atomically validates the epoch against durable accepted ownership state and performs the write only if authorized by that protocol.
4. A new owner establishes its epoch at the resource before treating takeover as complete.
5. Older owners cannot overwrite state once the newer epoch is established.

The resource must actually enforce the guard. Carrying a token in a log, checking it only in a worker, or using an unrelated wall-clock timestamp does not fence writes. Equal-token duplicate requests and ordering within one owner need separate idempotency/version rules.

If one authoritative store already offers the needed conditional writes, a separate lock service may be unnecessary. Prefer a well-specified existing primitive over a custom lease implementation. A lock used only to avoid duplicate computation has a different failure consequence from one protecting financial correctness; assess that consequence explicitly.

## Worked example: a paused worker

**Reconstructed source mechanism:** worker A receives epoch 41, reads a document, and pauses. Its lease expires. Worker B receives epoch 42, establishes it at storage, and writes a newer document. A resumes and submits its old result with epoch 41.

Without resource-side fencing, A overwrites B despite the coordination service having correctly granted only one current lease. With durable atomic epoch checking, storage rejects A's stale write.

Now move A's pause: it sends its request first, but the network holds that request until after B's write. The same guard must reject the delayed request. A client-side “lease still valid” check cannot protect this history.

Finally test resource restart. If the accepted epoch existed only in volatile memory and resets, the stale write may succeed after restart. Persist or reconstruct the fencing state according to the storage protocol. This is a system invariant, not merely a field added to a request.

## Separate safety, liveness, and assumptions

**Safety:** something forbidden never happens, such as two accepted owners for an exclusive resource. Once violated, a later recovery cannot erase the fact that the violation occurred.

**Liveness:** useful progress eventually happens under stated conditions. A system that rejects all requests can preserve some safety properties while being useless.

State the model: crash-stop versus crash-recovery nodes, asynchronous or partially synchronous timing assumptions, durable storage behavior, and permitted message failures. Byzantine faults include arbitrary or malicious behavior; ordinary crash-fault consensus is not protection against them. Checksums, authentication, and access boundaries address different risks and should not be conflated with a Byzantine consensus protocol.

When the environment violates a critical assumption, fail visibly or enter a defined recovery state. Do not silently downgrade a correctness guarantee to regain apparent availability.

## Verify histories, not just outcomes

Use the simplest verification method capable of exposing the risk:

- **Concurrent integration tests:** control the interleaving around reads, writes, commits, and acknowledgments.
- **Property-based tests:** generate operations and assert invariants against a small reference model.
- **Deterministic simulation:** control scheduling, network delivery, clocks, and crashes with reproducible seeds.
- **Model checking:** explore a bounded state model of a protocol; report the model and bounds.
- **Fault injection:** exercise the deployed implementation under faults in an authorized test environment.

Record invocation, completion, response, and uncertainty in a history. An operation that timed out may have taken effect, which the checker must allow when analyzing consistency. Check behavior during recovery, not only after waiting for convergence.

A model proof assumes a model; a bounded search covers its bounds; randomized tests cover sampled executions. None alone proves the implementation correct in every real environment. Preserve counterexample traces and turn them into regression cases.

Use fault injection only within the user's authorized environment and operational constraints. Architecture advice is not permission to kill production nodes or modify network rules.

Related: [transactions](ch08-transactions.md), [consensus](ch10-consistency-and-consensus.md), [stream recovery](ch12-stream-processing.md).
