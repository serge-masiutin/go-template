package httpapp

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSONRejectsBrokenContracts(t *testing.T) {
	for _, body := range []string{`null`, `[]`, `{"name":"ok","admin":true}`, `{"name":1}`, `{"name":"ok"} {}`, `{"name":`} {
		t.Run(body, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			var value struct {
				Name string `json:"name"`
			}
			if err := decodeJSON(r, &value); err == nil {
				t.Fatal("invalid JSON contract accepted")
			}
		})
	}
}
