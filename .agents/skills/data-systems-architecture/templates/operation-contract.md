# Effectful operation contract

Use for one command or event-processing operation. Keep the scope small enough that its invariants and failure outcomes can be stated precisely.

| Field | Define |
|---|---|
| Intent and actor | What business action is requested, and who may request it? |
| Input boundary | Schema/version, required fields, units, normalization, authorization |
| Logical identity | Stable request/event ID, scope, payload-conflict policy, retention |
| Authority | Owning database/resource and all competing writer paths |
| Invariant | Condition that must hold before/after every accepted operation |
| Decision mechanism | Constraint, atomic operation, isolation, locked resource, or coordinated protocol |
| Local atomic boundary | State, deduplication claim, result, and publication intent committed together |
| External effects | Destinations, idempotency contracts, timeouts, reconciliation paths |
| Acknowledgment | What the response confirms: recorded, accepted, committed, or completed |
| Unknown outcome | How a lost response is resolved without duplicating the intent |
| Delivery and order | Redelivery, per-key order, concurrent consumers, fencing/progress rules |
| Read visibility | Authoritative read, progress token, pending state, permitted lag |
| Retry | Which errors, which layer, whole-operation scope, budget/backoff |
| Recovery | Restore/replay inputs, retained history, cleanup and deletion behavior |
| Evidence | Concurrent/crash tests, invariants checked, observed results and limits |

Exercise at least the applicable commit and acknowledgment boundaries. A response field named `success` is not a substitute for defining what has become durable.
