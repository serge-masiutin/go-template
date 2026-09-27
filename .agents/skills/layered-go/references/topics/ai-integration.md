# AI Integration

## Summary

AI features need the same explicit contracts as other external integrations. Preserve the source's provider boundary, orchestration, tools, structured outputs, embeddings and testing topics; use the installed Genkit flow/tool APIs without translating Active Agent inheritance into Go.

## Layer Placement

Provider HTTP/SDK code is infrastructure. A feature operation owns prompt selection, tool policy and domain workflow. Domain entities do not call a model. Prompts and schemas are versioned production artifacts beside the feature.

## Key Principles

Start with one bounded model call. Add orchestration, retrieval or agents only after representative evaluations demonstrate a need. Treat retrieved documents, model output and tool arguments as untrusted inputs, not authority.

## Implementation with Explicit Boundaries

### Basic Provider

```go
type SummaryRequest struct { Document string; PromptVersion string }
type Summary struct { Text string }
type Summarizer interface {
    Summarize(context.Context, SummaryRequest) (Summary, error)
}
```

This is an illustrative feature boundary, not an installed starter API. The application sets deadline, input size and output schema. The provider adapter validates output and preserves errors without logging document content.

### Tools

Each tool has a typed input/output, narrow responsibility, permission contract and explicit side effects. Validate arguments and authorize the actor independently of the model's request. A model-generated instruction cannot grant filesystem, network or account access.

### Service Integration

The operation loads authorized data, renders a known prompt version, calls the provider and validates its result before use. Separate proposing an external action from executing it when review is required. Never execute generated SQL or shell commands as trusted output.

## Without an Agent Framework

A wrapped provider client is sufficient for summarization or classification. Supply the model identifier, timeout, budget and schema explicitly. Retry only documented transient failures and only when the operation's side effects are safe to repeat.

### AI-Powered Operation

Return typed output plus error; distinguish schema failure, refusal, timeout, rate limit and provider outage. Do not substitute fabricated text on error. If degraded behavior is useful, make it a named product outcome with metrics.

## Embeddings and Vector Search

Begin with a retrieval baseline. Version corpus snapshot, preprocessing, chunking, embedding model and index. Enforce tenant scope before retrieved text enters a prompt. Avoid train/evaluation leakage and test missing, stale and adversarial documents. Do not add a vector database merely because embeddings exist.

## Testing AI Features

Use deterministic tests for parsing, authorization, tools and provider failure handling. Maintain representative quality evaluations for normal, malformed, ambiguous, adversarial and long-context inputs. Record prompt/model/schema versions and compare quality, latency and cost before rollout. HTTP mocks alone do not establish model quality.

## Anti-Patterns

### LLM Calls in Models

A domain method should not hide variable-latency network access. Put it behind the feature operation and provider port.

### Unhandled AI Failures

Timeouts, invalid output and rate limits are explicit outcomes. Cancellation propagates to the adapter. A broad catch returning an empty summary hides a production failure.

### Prompt Strings Everywhere

Keep one versioned prompt source per feature and explicit template inputs. Do not refer to hidden conversation memory or ask the model to infer missing contracts.

## Background Processing

Long-running work uses a durable job with bounded attempts, idempotency and progress/failure states. Do not launch an untracked goroutine after an HTTP response. The starter implements this through River and `internal/assistant`: the job holds a run ID, the operation reloads its actor, and a bounded Genkit flow reads only that actor's notes. Gemini/OpenAI adapters use explicit models and no automatic generation retries. Vector search and multi-agent orchestration are not installed. See [installed contracts](../installed-stack.md) before extending the flow.
