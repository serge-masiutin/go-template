# Typed Configuration

This replaces Anyway Config's inheritance and source chain with an explicit Go loader.

## Basic Usage

Read configuration once in the composition root, validate it and pass typed settings to consumers. `internal/config/config.go` is the starter's source of truth.

## Shared and Nested Configuration

Compose structs by owner. A database adapter receives connection/pool settings; it does not receive model-provider or mail secrets. Do not create a global mutable singleton.

## Type Coercion and Validation

Use `strconv`, `time.ParseDuration` and URL parsing. An omitted optional value may use its documented boundary default; a supplied invalid value must return an error. Validate enum values, positive limits and required production settings.

## Environment Sources and Local Overrides

The Go binary reads its environment. Local scripts source the selected `.env`. There is no automatic YAML merge, Rails credentials, parent-directory discovery or hidden source priority. Document precedence before adding a second source.

## Secrets

Production supplies secrets through its deployment environment/secret provider. Keep `.env` ignored, never print credentials, and define restart/reload behavior for rotation. Do not commit a replacement encrypted-credentials format.

## Testing

Test every required field and invalid value, production HTTPS requirements and secret-free errors. Use isolated environment changes or inject a lookup function; tests changing process environment cannot run concurrently. See [configuration](../topics/configuration.md).
