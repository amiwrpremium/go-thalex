package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/internal/transport"
)

// newTestClient creates a REST Client with its transport pointed at the given
// httptest.Server. Because this file lives in package rest (internal test), we
// can set the unexported transport field directly.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	tr := transport.NewHTTPTransport(transport.HTTPTransportConfig{
		Client:        server.Client(),
		BaseURL:       server.URL,
		UserAgent:     config.UserAgent,
		MaxRetries:    0,
		RetryBaseWait: 1 * time.Millisecond,
	})

	return &Client{
		transport: tr,
		cfg:       config.DefaultClientConfig(),
	}
}

// newTestClientWithAuth creates a REST Client that sends an Authorization header.
func newTestClientWithAuth(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	tokenFunc := func() (string, error) {
		return "test-token-123", nil
	}

	tr := transport.NewHTTPTransport(transport.HTTPTransportConfig{
		Client:        server.Client(),
		BaseURL:       server.URL,
		UserAgent:     config.UserAgent,
		MaxRetries:    0,
		RetryBaseWait: 1 * time.Millisecond,
		TokenFunc:     tokenFunc,
		AccountNumber: "ACC-001",
	})

	cfg := config.DefaultClientConfig()
	cfg.AccountNumber = "ACC-001"
	return &Client{
		transport: tr,
		cfg:       cfg,
	}
}

// wrapResult wraps v into the Thalex REST envelope: {"result": ...}
func wrapResult(t *testing.T, v any) []byte {
	t.Helper()
	resultJSON, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	return []byte(fmt.Sprintf(`{"result":%s}`, resultJSON))
}

// apiErrorJSON returns a Thalex API error response body.
func apiErrorJSON(code int, message string) []byte {
	return []byte(fmt.Sprintf(`{"error":{"code":%d,"message":%q}}`, code, message))
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient()
	if c.cfg.Network != config.Production {
		t.Errorf("expected network Production, got %v", c.cfg.Network)
	}
	if c.cfg.MaxRetries != 3 {
		t.Errorf("expected MaxRetries=3, got %d", c.cfg.MaxRetries)
	}
	if c.cfg.UserAgent != config.UserAgent {
		t.Errorf("expected UserAgent=%q, got %q", config.UserAgent, c.cfg.UserAgent)
	}
	if c.cfg.WSReconnect != false {
		t.Errorf("expected WSReconnect=false")
	}
	if c.transport == nil {
		t.Fatal("expected transport to be non-nil")
	}
}

func TestNewClient_WithAllOptions(t *testing.T) {
	customHTTP := &http.Client{Timeout: 60 * time.Second}
	c := NewClient(
		config.WithNetwork(config.Testnet),
		config.WithHTTPClient(customHTTP),
		config.WithMaxRetries(5),
		config.WithRetryBaseWait(1*time.Second),
		config.WithUserAgent("custom-agent/1.0"),
		config.WithAccountNumber("ACC-123"),
	)

	if c.cfg.Network != config.Testnet {
		t.Errorf("expected network Testnet, got %v", c.cfg.Network)
	}
	if c.cfg.HTTPClient != customHTTP {
		t.Error("expected custom HTTP client")
	}
	if c.cfg.MaxRetries != 5 {
		t.Errorf("expected MaxRetries=5, got %d", c.cfg.MaxRetries)
	}
	if c.cfg.RetryBaseWait != 1*time.Second {
		t.Errorf("expected RetryBaseWait=1s, got %v", c.cfg.RetryBaseWait)
	}
	if c.cfg.UserAgent != "custom-agent/1.0" {
		t.Errorf("expected UserAgent=%q, got %q", "custom-agent/1.0", c.cfg.UserAgent)
	}
	if c.cfg.AccountNumber != "ACC-123" {
		t.Errorf("expected AccountNumber=%q, got %q", "ACC-123", c.cfg.AccountNumber)
	}
	if c.transport == nil {
		t.Fatal("expected transport to be non-nil")
	}
}

func TestNewClient_WithCredentials(t *testing.T) {
	// We cannot easily test with real credentials, but we can verify the
	// config stores them.
	c := NewClient(
		config.WithNetwork(config.Testnet),
	)
	if c.cfg.Credentials != nil {
		t.Error("expected nil credentials when none provided")
	}
}

func TestClient_APIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10001, "invalid parameter"))
	})

	_, err := c.Instruments(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "API error 10001: invalid parameter" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestClient_HTTPError_NoBody(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
	})

	_, err := c.Instruments(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClient_InvalidJSON(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	})

	_, err := c.Instruments(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestClient_NullResult(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"result":null}`))
	})

	instruments, err := c.Instruments(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if instruments != nil {
		t.Errorf("expected nil instruments for null result, got %v", instruments)
	}
}

func TestClient_ContextCanceled(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := c.Instruments(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
