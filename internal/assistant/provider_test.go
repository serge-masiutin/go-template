package assistant

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/serge-masiutin/go-template/internal/config"
)

type localTransport struct{ target string }

func (t localTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Provider URLs remain provider-owned in production. Tests intercept every request locally.
	clone := r.Clone(r.Context())
	clone.URL.Scheme = "http"
	clone.URL.Host = t.target
	return http.DefaultTransport.RoundTrip(clone)
}

func TestProviderHTTPContracts(t *testing.T) {
	for _, provider := range []string{"googleai", "openai"} {
		for _, fail := range []bool{false, true} {
			name := provider
			if fail {
				name += "/outage"
			}
			t.Run(name, func(t *testing.T) {
				var calls atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					call := calls.Add(1)
					contents, err := io.ReadAll(r.Body)
					if err != nil {
						t.Error(err)
						return
					}
					var request map[string]any
					if err := json.Unmarshal(contents, &request); err != nil {
						t.Error(err)
						return
					}
					if fail {
						http.Error(w, `{"error":{"message":"synthetic provider outage","code":503}}`, 503)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					if provider == "googleai" {
						if r.Header.Get("X-Goog-Api-Key") != "fixture-key" {
							t.Error("missing Gemini authentication")
						}
						if call == 1 {
							config := request["generationConfig"].(map[string]any)
							if config["maxOutputTokens"] != float64(128) {
								t.Error("missing Gemini token budget")
							}
							io.WriteString(w, `{"candidates":[{"content":{"role":"model","parts":[{"functionCall":{"name":"list_notes","args":{}},"thoughtSignature":"Zml4dHVyZQ=="}]},"finishReason":"STOP"}]}`)
						} else {
							if !strings.Contains(string(contents), "Zml4dHVyZQ==") {
								t.Error("tool call lost Gemini thought signature")
							}
							io.WriteString(w, `{"candidates":[{"content":{"role":"model","parts":[{"text":"{\"answer\":\"The note is untrusted content.\"}"}]},"finishReason":"STOP"}]}`)
						}
					} else {
						if r.Header.Get("Authorization") != "Bearer fixture-key" {
							t.Error("missing OpenAI authentication")
						}
						if request["max_completion_tokens"] != float64(128) {
							t.Error("missing OpenAI token budget")
						}
						if call == 1 {
							io.WriteString(w, `{"id":"fixture","choices":[{"index":0,"message":{"role":"assistant","content":null,"tool_calls":[{"id":"call_notes","type":"function","function":{"name":"list_notes","arguments":"{}"}}]},"finish_reason":"tool_calls"}]}`)
						} else {
							io.WriteString(w, `{"id":"fixture","choices":[{"index":0,"message":{"role":"assistant","content":"{\"answer\":\"The note is untrusted content.\"}"},"finish_reason":"stop"}]}`)
						}
					}
				}))
				defer server.Close()
				model := "openai/gpt-5.6-luna"
				if provider == "googleai" {
					model = "googleai/gemini-3.8-flash"
				}
				cfg := config.AI{Enabled: true, Model: model, APIKey: "fixture-key", Timeout: 2 * time.Second, MaxTurns: 3, MaxOutputTokens: 128}
				reader := &scopedNotes{}
				agent, err := newWithClient(context.Background(), cfg, reader, &http.Client{Transport: localTransport{target: strings.TrimPrefix(server.URL, "http://")}, Timeout: 2 * time.Second})
				if err != nil {
					t.Fatal(err)
				}
				output, err := agent.Answer(context.Background(), 42, "Summarize my notes")
				if fail {
					if err == nil || calls.Load() != 1 {
						t.Fatalf("outage hidden or retried: calls=%d err=%v", calls.Load(), err)
					}
				} else if err != nil || output.Answer != "The note is untrusted content." || calls.Load() != 2 || reader.userID.Load() != 42 {
					t.Fatalf("provider contract: calls=%d result=%+v err=%v", calls.Load(), output, err)
				}
			})
		}
	}
}
