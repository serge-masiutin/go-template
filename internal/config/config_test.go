package config

import "testing"

func TestProductionConfiguration(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/starter")
	t.Setenv("APP_ENV", "production")
	t.Setenv("PUBLIC_URL", "http://example.com")
	if _, err := Load(); err == nil {
		t.Fatal("HTTP production origin accepted")
	}
	t.Setenv("PUBLIC_URL", "https://example.com")
	t.Setenv("DB_MAX_CONNECTIONS", "0")
	if _, err := Load(); err == nil {
		t.Fatal("zero connection budget accepted")
	}
	t.Setenv("DB_MAX_CONNECTIONS", "12")
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !got.Production || got.MaxConnections != 12 {
		t.Fatalf("unexpected config: %#v", got)
	}
}

func TestOptionalServicesFailAtConfigurationBoundary(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/starter")
	t.Setenv("APP_ENV", "production")
	t.Setenv("PUBLIC_URL", "https://example.test")
	t.Setenv("MAIL_ENABLED", "true")
	t.Setenv("MAIL_HOST", "smtp.example.test")
	t.Setenv("MAIL_FROM", "starter@example.test")
	t.Setenv("MAIL_TLS", "none")
	if _, err := Load(); err == nil {
		t.Fatal("plaintext production SMTP accepted")
	}
	t.Setenv("MAIL_TLS", "starttls")
	t.Setenv("AI_ENABLED", "true")
	if _, err := Load(); err == nil {
		t.Fatal("enabled AI without provider credentials accepted")
	}
	t.Setenv("AI_MODEL", "googleai/gemini-3.8-flash")
	t.Setenv("AI_API_KEY", "synthetic-key")
	if _, err := Load(); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GENKIT_TELEMETRY_SERVER", "http://localhost:9999")
	if _, err := Load(); err == nil {
		t.Fatal("private content exporter accepted")
	}
}
