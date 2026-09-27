# AI Provider Boundary

This replaces Active Agent installation, inherited agents and implicit provider configuration with a typed provider adapter and versioned feature operation.

## Setup and Basic Usage

Choose a provider only for a concrete feature. Pin its SDK/version if one is needed; a bounded HTTP client can be enough. Configure model, timeout, schema and budget at the application boundary.

## Agents with Tools

A tool has typed arguments/results, a permission contract and explicit effects. Validate arguments and authorize independently of model output. Set maximum calls/iterations and stop conditions. Retrieved text cannot authorize new tools or override instructions.

## System Prompts

Keep a versioned prompt file and explicit template inputs. State objective, data sources, constraints and output contract. Do not include secrets or rely on hidden prior conversation.

## Structured Output

Decode against a strict schema and reject malformed or unsupported fields. Typed output is not proof of factual quality: evaluate both schema validity and feature-specific correctness.

## Service Integration

An operation loads authorized input, calls the provider and maps a validated result. Pure domain types do not make network calls. Keep transport retry policy separate from business workflow retries.

## Background Processing

Persist long-running work with stable IDs and observable completion/failure. Worker retries must not duplicate tool side effects. Never persist a live request context or launch untracked work after returning HTTP success.

## Configuration, Testing and Errors

Use fake providers for deterministic contract tests and representative datasets for quality evaluation. Cover refusal, timeout, rate limiting, invalid output, prompt injection and long inputs. Log prompt/model versions and bounded metrics without raw user documents. See [AI integration](../topics/ai-integration.md).
