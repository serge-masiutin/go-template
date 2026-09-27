# Responsible data use, privacy, and accountable decisions

## Core idea

Data architecture determines who can observe, infer, and act on information about people. Collection, retention, derivation, and automated decisions are design choices with consequences; technical correctness alone does not establish that the system is appropriate.

Use this chapter when the task actually involves personal data, tracking, predictive decisions, or consequential derived profiles. Keep the analysis proportional and concrete. It is an engineering decision framework, not a substitute for current jurisdiction-specific legal advice.

## Make purpose and data flows explicit

For each category of personal or potentially identifying data, document:

- The product purpose and the smallest data representation that serves it.
- Who supplies it, who is affected, and whether it also reveals facts about non-users.
- Who can access the raw data and derived results, including vendors and internal teams.
- Retention periods and the reason for retaining each representation.
- Correction, export, restriction, and deletion behavior required by the product and applicable rules.
- The owner accountable for misuse, inaccurate data, or a failed deletion process.

Do not infer that information collected for one purpose is automatically suitable for another. A helpful user feature, targeted advertising, fraud prevention, and staff monitoring have different purposes and power relationships even if they use the same event stream.

Data minimization can mean lower precision, shorter retention, aggregation, local processing, or no collection. Each reduces different exposure. Choose after identifying the necessary function rather than collecting everything because storage is inexpensive.

## Privacy is control over context, not universal secrecy

A person can willingly share information with one audience for one purpose and reasonably reject another use. Publicly posting one fact does not mean all inferred attributes or linked histories are welcome. Access controls between users do not alone constrain how the service itself uses data.

Explain collection and downstream use in terms users can understand. If a feature relies on consent, provide meaningful choices and implement withdrawal semantics where required. Do not assume consent is the only possible legal basis or treat a generic checkbox as proof every future processing purpose is authorized. Determine actual legal requirements from current authoritative sources when making legal claims.

Assess foreseeable secondary access: breach, insider misuse, vendor transfer, account compromise, acquisition, or compelled disclosure. Retention is an exposure decision as well as a storage-cost decision. Sensitive inference can arise from seemingly ordinary location, device, or behavior records.

## Propagate correction and deletion through derivation

An authoritative delete is incomplete if derived copies keep serving the fact. Trace the entire lifecycle: primary rows, replicas, CDC streams, event logs, search indexes, caches, exports, analytical snapshots, features, models, backups, and observability systems.

For each representation choose a supported mechanism and a completion condition. Examples include expiring raw events, emitting deletion records, rebuilding affected projections, deleting a subject-specific object, or separating direct identifiers from long-lived events. Keep only the minimal evidence needed to establish that the operation completed.

Replays and restores can resurrect deleted data. A recovery procedure must apply the current deletion/retention policy before the restored dataset becomes available, according to the system's stated obligations. Tombstone retention, old snapshots, offline clients, and delayed consumers all affect that procedure.

Crypto-shredding may be a useful mechanism when all relevant copies are encrypted under appropriately scoped keys and key destruction actually removes access. It is not a universal deletion certificate: plaintext copies, derived outputs, key backups, and other reconstruction paths can remain. Validate the real design.

Aggregate data or model outputs can still expose information; removing a name is not proof of anonymity. Treat model unlearning, irreversible aggregation, and legal deletion questions as specific requirements to investigate rather than promising an unsupported generic solution.

## Evaluate predictive systems as decisions in a larger system

Distinguish a probabilistic prediction from an established fact about a person. Good aggregate accuracy can coexist with harmful errors for particular people or groups. Historical labels can encode past unequal treatment; excluding an explicitly sensitive column does not remove correlated proxies.

Before using a prediction to influence a consequential action:

1. State the decision, who benefits, who bears errors, and the cost of false positives/negatives.
2. Examine how labels were created, who is missing from the data, and which outcomes were never observed.
3. Evaluate relevant groups and contexts with appropriate statistical limits; do not hide small-sample uncertainty.
4. Define meaningful review, contestability, correction, and an accountable decision owner.
5. Monitor outcomes after deployment, including behavior changes caused by the system itself.

These are architecture and process requirements. They do not specify a universal fairness metric; conflicting objectives and context may require explicit human judgment. A model or vendor score does not transfer accountability away from the organization using it.

## Look for feedback loops

The system's decisions can change the future data used to evaluate or retrain it. A recommendation increases exposure, increased exposure produces clicks, and those clicks can be mistaken for independent evidence of preference. A rejection can prevent the positive outcome that would have contradicted a pessimistic prediction.

Map the loop: observation → score → action → changed opportunity/behavior → new observation. Identify which outcomes are censored by the action and which proxy objective may diverge from the intended benefit. Where justified, use independent evaluation, carefully governed exploration, and human review; do not equate more engagement with better welfare by default.

## Worked example: account-risk flag and data correction

**Reconstructed application of the source's accountability reasoning:** an account receives a risk flag from a model using past activity and a derived profile. The flag restricts access. A user reports that the profile contains another person's activity due to an identity-linking error.

A complete design identifies the source records, linkage rule/version, model version, and decision inputs. A responsible operator can correct the link, rebuild the relevant profile, reassess the decision, and track which downstream restrictions used the incorrect result. The system explains a useful basis for review without exposing other users' data.

If the only retained artifact is a score with no provenance, the team cannot distinguish model error from bad input or incorrect linkage. If every new model run reuses the uncorrected warehouse copy, correcting the primary account row does not solve the problem. If denial prevents the user from ever demonstrating a good outcome, the training data may reinforce the original mistake.

The architecture therefore needs lineage, correction propagation, review ownership, and appropriate retention. It does not need unlimited raw-history retention: retain what is justified for those purposes, with explicit access and expiry.

## Practical review output

Add only the relevant findings to the architectural decision: unnecessary collection, uncontrolled reuse, missing deletion paths, unverifiable prediction provenance, proxy discrimination risk, or an unowned correction process. For each finding identify the affected flow, consequence, concrete change, and verification method.

Do not claim compliance solely because a checklist is complete or because the source book describes a law. Current legal requirements depend on jurisdiction and context; obtain authoritative review where the task requires it. The engineering responsibility remains to make the actual data behavior visible and changeable.

Related: [ownership and retention](ch01-architecture-tradeoffs.md), [event histories](ch03-data-models.md), [derived-state rebuilds](ch13-dataflow-and-correctness.md).
