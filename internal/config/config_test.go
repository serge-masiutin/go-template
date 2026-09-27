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
