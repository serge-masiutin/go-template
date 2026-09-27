# Configuration

## Summary

Configuration belongs in validated types with explicit sources. Preserve the source chapter's separation between schema, loading and use; replace framework config inheritance with ordinary Go structs.

## Layer Placement

Environment/file/secret-provider loading is infrastructure at startup. Each consumer receives the settings it owns. A domain rule may take a typed policy value; it must not read environment variables.

## Key Principles

Load once, validate before serving, document defaults at the boundary, protect secrets, and keep configuration immutable during normal request handling. Do not scatter `os.Getenv` through business code.

## Implementation with Typed Configuration

### Basic Configuration

```go
type MailConfig struct {
    Host string
    Port int
    Timeout time.Duration
}
func (c MailConfig) Validate() error {
    if c.Host == "" || c.Port < 1 || c.Port > 65535 || c.Timeout <= 0 {
        return ErrInvalidMailConfiguration
    }
    return nil
}
```

This illustrates the configuration boundary. The starter's actual `config.Mail` and `config.AI` live in `internal/config/config.go`; caarlos0/env parses their tags before semantic validation. Read those types before changing environment keys.

### Shared Configuration

Use composition for related settings. Do not create a base config class or expose every secret to every consumer. Pass a narrow config value to a constructor.

### Usage

Parse and validate in `main`, construct adapters explicitly, then start accepting requests. Startup errors must be nonzero exits. Do not print database URLs, passwords or provider keys.

### Nested Configuration

Group values by owner: database pool limits belong together; email transport settings are separate. Nested structs should clarify ownership rather than mirror arbitrary environment-variable naming.

### Feature Flags

A static flag is a typed boolean with a documented default. Dynamic flag providers require an explicit freshness/failure contract and should not be introduced merely to hold a boolean.

### Validation

Reject malformed integers, invalid enum values, negative timeouts and absent required production settings. A fallback is legitimate only when it is the documented default for an omitted value, not for an invalid value.

### Environment-Specific Behavior

Production requires a public HTTPS origin and secure cookies. Test configuration should point to isolated resources. Environment names are validated enums, not arbitrary strings that accidentally activate development behavior.

## Without a Configuration Framework

The standard library (`os`, `strconv`, `time`, `net/url`) is enough for this starter. `.env` is sourced by local scripts; the Go binary receives environment variables and does not search parent directories for secret files.

### Secret Provider Integration

If deploying with a secret manager, load secrets at the process boundary and inject them. Do not commit credentials or add a home-grown encrypted credential format. Rotation requires a deliberate reload/restart strategy.

## Testing Configuration

Use `t.Setenv` in nonparallel tests or test a loader accepting an explicit lookup function. Cover omitted defaults, invalid values, production HTTPS and missing secrets. Assert errors never contain secret values.

## Anti-Patterns

### Scattered Environment Access

Hidden reads make a function's behavior depend on process state and make test order matter. Move them to the loader.

### Unvalidated Configuration

A misspelled timeout should fail startup instead of becoming zero and disabling a limit.

### Mutable Configuration

Do not update package-level settings per request. If live reconfiguration is needed, define atomic replacement and reader ownership explicitly.

## File Organization

Use `internal/config` for this starter's loader and tests. Consumers accept only their needed settings. Avoid duplicating defaults in Docker, Go code, scripts and docs.

## Configuration Sources and Priority

The starter uses environment variables and explicit development defaults. Local scripts load the chosen `.env` file. No YAML merge, Rails credentials or implicit precedence chain exists. If adding another source, specify precedence, validate the final value and test conflicting values.
