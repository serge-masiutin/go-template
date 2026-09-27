//go:build integration

package accounts

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/serge-masiutin/go-template/internal/testdb"
	"gorm.io/gorm"
)

func TestLoginSessionRevocationExpiryAndRotation(t *testing.T) {
	ctx := context.Background()
	db := testdb.New(t)
	store, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.Create(ctx, "sessions@example.test", "session-test-password", false)
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range []string{"first-device", "second-device"} {
		if err := store.CreateLoginSession(ctx, user.ID, token, "", time.Now().Add(time.Hour)); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.CreateLoginSession(ctx, user.ID, "rotated-first", "first-device", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	assertRejected := func(token string) {
		t.Helper()
		if _, err := store.FindBySession(ctx, token); !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("revoked or expired token accepted: %v", err)
		}
	}
	assertRejected("first-device")
	for range 2 {
		if err := store.RevokeSession(ctx, "rotated-first"); err != nil {
			t.Fatal(err)
		}
	}
	assertRejected("rotated-first")
	if found, err := store.FindBySession(ctx, "second-device"); err != nil || found != user {
		t.Fatalf("other device was revoked: %v", err)
	}
	if err := store.CreateLoginSession(ctx, user.ID, "expired", "", time.Now().Add(-time.Second)); err != nil {
		t.Fatal(err)
	}
	assertRejected("expired")
	if err := store.CreateLoginSession(ctx, user.ID, "new", "", time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	var expired int
	if err := db.QueryRow(ctx, "SELECT count(*) FROM login_sessions WHERE expires_at <= now()").Scan(&expired); err != nil {
		t.Fatal(err)
	}
	if expired != 0 {
		t.Fatal("sign-in did not reclaim expired grants")
	}
	if _, err := db.Exec(ctx, "DELETE FROM users WHERE id = $1", user.ID); err != nil {
		t.Fatal(err)
	}
	assertRejected("second-device")
}
