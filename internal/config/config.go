package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v11"
	"net/mail"
	"net/url"
	"os"
	"strings"
	"time"
)

// Error describes invalid configuration without including supplied secret values.
type Error struct{ message string }

func (e *Error) Error() string { return e.message }

type Config struct {
	Address           string        `env:"HTTP_ADDR" envDefault:"127.0.0.1:3000"`
	DatabaseURL       string        `env:"DATABASE_URL"`
	PublicURL         string        `env:"PUBLIC_URL"`
	Environment       string        `env:"APP_ENV" envDefault:"development"`
	Production        bool          `env:"-"`
	MaxConnections    int32         `env:"DB_MAX_CONNECTIONS" envDefault:"10"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	WorkerConcurrency int           `env:"WORKER_CONCURRENCY" envDefault:"4"`
	MetricsToken      string        `env:"METRICS_TOKEN"`
	Mail              Mail          `envPrefix:"MAIL_"`
	AI                AI            `envPrefix:"AI_"`
}

type Mail struct {
	Enabled  bool          `env:"ENABLED" envDefault:"false"`
	Host     string        `env:"HOST"`
	Port     int           `env:"PORT" envDefault:"587"`
	From     string        `env:"FROM"`
	Username string        `env:"USERNAME"`
	Password string        `env:"PASSWORD"`
	TLS      string        `env:"TLS" envDefault:"starttls"`
	Timeout  time.Duration `env:"TIMEOUT" envDefault:"10s"`
}

type AI struct {
	Enabled         bool          `env:"ENABLED" envDefault:"false"`
	Model           string        `env:"MODEL"`
	APIKey          string        `env:"API_KEY"`
	Timeout         time.Duration `env:"TIMEOUT" envDefault:"30s"`
	MaxTurns        int           `env:"MAX_TURNS" envDefault:"3"`
	MaxOutputTokens int           `env:"MAX_OUTPUT_TOKENS" envDefault:"1024"`
}

func Load() (Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		var invalid env.ParseError
		if errors.As(err, &invalid) {
			return c, &Error{message: "invalid configuration field: " + invalid.Name}
		}
		return c, &Error{message: "invalid environment configuration"}
	}
	if c.DatabaseURL == "" {
		return c, &Error{message: "DATABASE_URL is required"}
	}
	switch c.Environment {
	case "", "development", "test":
	case "production":
		c.Production = true
	default:
		return c, &Error{message: "APP_ENV must be development, test, or production"}
	}
	if c.PublicURL == "" && !c.Production {
		c.PublicURL = "http://localhost:3000"
	}
	parsed, err := url.Parse(c.PublicURL)
	if err != nil || parsed.Host == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return c, &Error{message: "PUBLIC_URL must be an HTTP origin without a path"}
	}
	if c.Production && parsed.Scheme != "https" {
		return c, &Error{message: "production PUBLIC_URL requires HTTPS"}
	}
	if c.MaxConnections < 2 {
		return c, &Error{message: "DB_MAX_CONNECTIONS must be at least 2"}
	}
	if c.WorkerConcurrency < 1 || c.WorkerConcurrency > 100 {
		return c, &Error{message: "WORKER_CONCURRENCY must be between 1 and 100"}
	}
	if c.ShutdownTimeout < time.Second || c.ShutdownTimeout > time.Minute {
		return c, &Error{message: "SHUTDOWN_TIMEOUT must be between 1s and 1m"}
	}
	if c.MetricsToken != "" && len(c.MetricsToken) < 32 {
		return c, &Error{message: "METRICS_TOKEN must contain at least 32 characters"}
	}
	if c.Mail.Enabled {
		parsed, err := mail.ParseAddress(c.Mail.From)
		if err != nil || parsed.Address != c.Mail.From || c.Mail.Host == "" || c.Mail.Port < 1 || c.Mail.Port > 65535 {
			return c, &Error{message: "mail requires a valid MAIL_FROM, MAIL_HOST and MAIL_PORT"}
		}
		if (c.Mail.Username == "") != (c.Mail.Password == "") {
			return c, &Error{message: "MAIL_USERNAME and MAIL_PASSWORD must be configured together"}
		}
		if c.Mail.TLS != "starttls" && c.Mail.TLS != "tls" && c.Mail.TLS != "none" {
			return c, &Error{message: "MAIL_TLS must be starttls, tls or none"}
		}
		if c.Mail.TLS == "none" && (c.Production || c.Mail.Username != "") {
			return c, &Error{message: "plaintext SMTP is allowed only for development without credentials"}
		}
		if c.Mail.Timeout < time.Second || c.Mail.Timeout > time.Minute {
			return c, &Error{message: "MAIL_TIMEOUT must be between 1s and 1m"}
		}
	}
	if c.AI.Enabled {
		provider, name, found := strings.Cut(c.AI.Model, "/")
		if !found || name == "" || (provider != "googleai" && provider != "openai") || c.AI.APIKey == "" {
			return c, &Error{message: "AI requires AI_API_KEY and AI_MODEL with googleai/ or openai/ prefix"}
		}
		if c.AI.Timeout < time.Second || c.AI.Timeout > 5*time.Minute || c.AI.MaxTurns < 1 || c.AI.MaxTurns > 10 || c.AI.MaxOutputTokens < 128 || c.AI.MaxOutputTokens > 8192 {
			return c, &Error{message: "invalid AI timeout, turn limit or output token budget"}
		}
		for _, key := range []string{"GENKIT_TELEMETRY_SERVER", "GENKIT_REFLECTION_V2_SERVER", "GENKIT_ENV"} {
			if os.Getenv(key) != "" {
				return c, &Error{message: fmt.Sprintf("%s is unsupported in the application worker; it can expose private AI content", key)}
			}
		}
	}

	return c, nil
}
