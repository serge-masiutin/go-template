package httpapp

import (
	"fmt"
	"testing"
	"time"
)

func TestLoginLimit(t *testing.T) {
	limiter := newLoginLimiter()
	now := time.Now()
	for n := 0; n < 10; n++ {
		if !limiter.Allow("127.0.0.1:1234", now) {
			t.Fatal("early limit")
		}
	}
	if limiter.Allow("127.0.0.1:5678", now) {
		t.Fatal("port bypassed limit")
	}
	if !limiter.Allow("127.0.0.2:1234", now) {
		t.Fatal("unrelated peer limited")
	}
	if !limiter.Allow("127.0.0.1:1234", now.Add(time.Minute)) {
		t.Fatal("window did not expire")
	}
	if limiter.Allow("malformed", now) {
		t.Fatal("invalid peer accepted")
	}
}
func TestLoginLimitBounded(t *testing.T) {
	limiter := newLoginLimiter()
	now := time.Now()
	for n := 0; n < 4096; n++ {
		if !limiter.Allow(fmt.Sprintf("host%d:1", n), now) {
			t.Fatal("early capacity limit")
		}
	}
	if limiter.Allow("extra:1", now) {
		t.Fatal("capacity unbounded")
	}
	if !limiter.Allow("extra:1", now.Add(time.Minute)) {
		t.Fatal("expired entries not reclaimed")
	}
}
