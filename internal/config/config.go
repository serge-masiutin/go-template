package config

import (
	"net/url"
	"os"
	"strconv"
	"time"
)

// Error describes invalid configuration without including supplied secret values.
type Error struct{ message string }

func (e *Error) Error() string { return e.message }

type Config struct {
	Address         string
	DatabaseURL     string
	PublicURL       string
	Production      bool
	MaxConnections  int32
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	c := Config{Address: "127.0.0.1:3000", MaxConnections: 10, ShutdownTimeout: 10 * time.Second}
	if value := os.Getenv("HTTP_ADDR"); value != "" {
		c.Address = value
	}
	c.DatabaseURL = os.Getenv("DATABASE_URL")
	if c.DatabaseURL == "" {
		return c, &Error{message: "DATABASE_URL is required"}
	}
	switch os.Getenv("APP_ENV") {
	case "", "development", "test":
	case "production":
		c.Production = true
	default:
		return c, &Error{message: "APP_ENV must be development, test, or production"}
	}
	c.PublicURL = os.Getenv("PUBLIC_URL")
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
	if value := os.Getenv("DB_MAX_CONNECTIONS"); value != "" {
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil || n < 1 {
			return c, &Error{message: "DB_MAX_CONNECTIONS must be a positive integer"}
		}
		c.MaxConnections = int32(n)
	}
	return c, nil
}
