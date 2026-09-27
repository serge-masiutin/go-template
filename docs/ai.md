# AI assistant

Genkit Go runs one typed flow, `notes_assistant_v1`, with the read-only `list_notes` tool. HTTP creates an operation and a River job. The worker receives the ID, loads the owner and question, calls the flow, and saves the answer. `/tools` shows the five latest requests and polls their state until work finishes.

## Enable the assistant

The application starts with `AI_ENABLED=false` and no API key. To enable it, set these values in your local `.env` or deployment secrets:

```sh
AI_ENABLED=true
AI_MODEL=googleai/gemini-3.8-flash
AI_API_KEY=your-provider-key
AI_TIMEOUT=30s
AI_MAX_TURNS=3
AI_MAX_OUTPUT_TOKENS=1024
```

The model ID above is used in local protocol tests; check its availability with your provider and account. For OpenAI, use the `openai/` prefix and a Chat Completions model ID supporting tools and structured output. This Genkit plugin version uses Chat Completions; Responses-only models are incompatible. The key must belong to the selected provider. Restart the server and worker after configuration changes.

Enabling AI requires both a model and key. Invalid configuration stops the process before any provider HTTP request. Constructing the client and starting the application do not generate content. Before submission, the interface tells users that their question and notes will be sent to the AI provider.

## Contract and boundaries

- Questions contain 1–500 Unicode characters; HTTP middleware limits the JSON body size.
- The tool accepts no owner ID. Each flow invocation creates an `ai.NewTool` that closes over the trusted user ID loaded for the operation. Reads are limited to the owner's 50 latest notes, each up to 2000 characters.
- The tool cannot write, send SMTP, run shell commands, browse, or fetch arbitrary URLs. Note text is data, not instructions.
- The prompt is `internal/assistant/prompts/notes-v1.txt`; the operation stores `PromptVersion` alongside the model.
- Output requires `answer`: nonempty plain text of at most 4000 characters. A successful answer without a tool call is rejected. React renders text, not model-generated HTML.
- `AI_TIMEOUT` is 1s–5m, defaulting to 30s for the entire flow; `AI_MAX_TURNS` is 1–10; `AI_MAX_OUTPUT_TOKENS` is 128–8192 per model call. Total cost depends on the number of turns and input notes.
- OpenAI SDK retries are disabled. The Gemini generate path used here has no automatic HTTP retries, verified with a 503 test. River mail/AI jobs also have one attempt.
- Raw provider errors are not saved in River or logs. The application retains wrapped causes for diagnosis through safe categories. Questions and answers live only in the owner-scoped table, not queue payloads or metrics.

Do not enable Genkit telemetry/reflection in production by copying a development quick start: those tools can retain full content. `GENKIT_ENV`, `GENKIT_TELEMETRY_SERVER`, and `GENKIT_REFLECTION_V2_SERVER` are rejected when AI is enabled. OTLP export, conversation memory, embeddings, and multi-agent orchestration are not configured. Define question/answer retention for your product; the starter does not remove them automatically except when the account is deleted.

## Validation

`go test ./internal/assistant` checks the schema, missing evidence, loop limits, provider outages, length boundaries, and both SDK HTTP contracts, including Gemini thought signatures. PostgreSQL integration tests check atomic enqueue, note isolation, cancellation after account deletion, and prevention of duplicate generation.

These checks use synthetic data and consume no provider credits. They verify application contracts, not a live model's reasoning quality. Before releasing the AI feature, run the [evaluation cases](../evals/notes-assistant.md) against your chosen model. Record its ID, prompt version, results, latency, and cost from provider data. Do not treat missing usage as zero.
