package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// NewHTTPTransport
// ---------------------------------------------------------------------------

func TestNewHTTPTransport(t *testing.T) {
	t.Run("default client when nil", func(t *testing.T) {
		tr := NewHTTPTransport(HTTPTransportConfig{})
		if tr.client == nil {
			t.Fatal("expected default client to be set")
		}
	})

	t.Run("default maxRetries when 0", func(t *testing.T) {
		tr := NewHTTPTransport(HTTPTransportConfig{})
		if tr.maxRetries != 3 {
			t.Errorf("maxRetries = %d; want 3", tr.maxRetries)
		}
	})

	t.Run("default maxRetries when negative", func(t *testing.T) {
		tr := NewHTTPTransport(HTTPTransportConfig{MaxRetries: -1})
		if tr.maxRetries != 3 {
			t.Errorf("maxRetries = %d; want 3", tr.maxRetries)
		}
	})

	t.Run("default retryBaseWait when 0", func(t *testing.T) {
		tr := NewHTTPTransport(HTTPTransportConfig{})
		if tr.retryBaseWait != 500*time.Millisecond {
			t.Errorf("retryBaseWait = %v; want 500ms", tr.retryBaseWait)
		}
	})

	t.Run("default retryBaseWait when negative", func(t *testing.T) {
		tr := NewHTTPTransport(HTTPTransportConfig{RetryBaseWait: -1})
		if tr.retryBaseWait != 500*time.Millisecond {
			t.Errorf("retryBaseWait = %v; want 500ms", tr.retryBaseWait)
		}
	})

	t.Run("custom values preserved", func(t *testing.T) {
		customClient := &http.Client{Timeout: 60 * time.Second}
		tr := NewHTTPTransport(HTTPTransportConfig{
			Client:        customClient,
			BaseURL:       "https://example.com",
			UserAgent:     "test-agent",
			MaxRetries:    5,
			RetryBaseWait: 2 * time.Second,
			AccountNumber: "ACC-123",
		})
		if tr.client != customClient {
			t.Error("expected custom client to be preserved")
		}
		if tr.baseURL != "https://example.com" {
			t.Errorf("baseURL = %q; want %q", tr.baseURL, "https://example.com")
		}
		if tr.userAgent != "test-agent" {
			t.Errorf("userAgent = %q; want %q", tr.userAgent, "test-agent")
		}
		if tr.maxRetries != 5 {
			t.Errorf("maxRetries = %d; want 5", tr.maxRetries)
		}
		if tr.retryBaseWait != 2*time.Second {
			t.Errorf("retryBaseWait = %v; want 2s", tr.retryBaseWait)
		}
		if tr.accountNumber != "ACC-123" {
			t.Errorf("accountNumber = %q; want %q", tr.accountNumber, "ACC-123")
		}
	})
}

// ---------------------------------------------------------------------------
// apiError
// ---------------------------------------------------------------------------

func TestApiError_Error(t *testing.T) {
	tests := []struct {
		code    int
		message string
		want    string
	}{
		{400, "Bad Request", "API error 400: Bad Request"},
		{-1, "Unknown error", "API error -1: Unknown error"},
		{0, "Zero", "API error 0: Zero"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			e := &apiError{Code: tc.code, Message: tc.message}
			got := e.Error()
			if got != tc.want {
				t.Errorf("Error() = %q; want %q", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DoPublic
// ---------------------------------------------------------------------------

func TestDoPublic(t *testing.T) {
	t.Run("successful request returns parsed result", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{
				Result: json.RawMessage(`{"ticker":"BTCUSD","price":50000}`),
			})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			UserAgent:     "test",
			RetryBaseWait: time.Millisecond,
		})

		var result struct {
			Ticker string  `json:"ticker"`
			Price  float64 `json:"price"`
		}
		err := tr.DoPublic(context.Background(), "/ticker", nil, &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Ticker != "BTCUSD" {
			t.Errorf("Ticker = %q; want %q", result.Ticker, "BTCUSD")
		}
		if result.Price != 50000 {
			t.Errorf("Price = %f; want %f", result.Price, 50000.0)
		}
	})

	t.Run("sends GET method", func(t *testing.T) {
		var receivedMethod string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedMethod = r.Method
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		tr.DoPublic(context.Background(), "/test", nil, nil)
		if receivedMethod != http.MethodGet {
			t.Errorf("method = %q; want %q", receivedMethod, http.MethodGet)
		}
	})

	t.Run("sets User-Agent header", func(t *testing.T) {
		var receivedUA string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedUA = r.Header.Get("User-Agent")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			UserAgent:     "go-thalex/1.0",
			RetryBaseWait: time.Millisecond,
		})
		tr.DoPublic(context.Background(), "/test", nil, nil)
		if receivedUA != "go-thalex/1.0" {
			t.Errorf("User-Agent = %q; want %q", receivedUA, "go-thalex/1.0")
		}
	})

	t.Run("appends query params to URL", func(t *testing.T) {
		var receivedURL string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedURL = r.URL.String()
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		params := url.Values{"instrument": {"BTCUSD"}, "limit": {"10"}}
		tr.DoPublic(context.Background(), "/trades", params, nil)
		if !strings.Contains(receivedURL, "instrument=BTCUSD") {
			t.Errorf("URL %q does not contain instrument=BTCUSD", receivedURL)
		}
		if !strings.Contains(receivedURL, "limit=10") {
			t.Errorf("URL %q does not contain limit=10", receivedURL)
		}
		if !strings.Contains(receivedURL, "?") {
			t.Errorf("URL %q should contain '?'", receivedURL)
		}
	})

	t.Run("no query params - no question mark in URL", func(t *testing.T) {
		var receivedURL string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedURL = r.URL.String()
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		tr.DoPublic(context.Background(), "/info", nil, nil)
		if strings.Contains(receivedURL, "?") {
			t.Errorf("URL %q should not contain '?' with no params", receivedURL)
		}
	})

	t.Run("server returns 400 with API error JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(apiResponse{
				Error: &apiError{Code: 1001, Message: "Invalid parameter"},
			})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/bad", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		apiErr, ok := err.(*apiError)
		if !ok {
			t.Fatalf("expected *apiError, got %T: %v", err, err)
		}
		if apiErr.Code != 1001 {
			t.Errorf("Code = %d; want 1001", apiErr.Code)
		}
	})

	t.Run("server returns 400 with non-JSON body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			w.Write([]byte("Bad Request: missing field"))
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/bad", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "HTTP 400") {
			t.Errorf("error = %q; want to contain 'HTTP 400'", err.Error())
		}
		if !strings.Contains(err.Error(), "Bad Request: missing field") {
			t.Errorf("error = %q; want to contain body text", err.Error())
		}
	})

	t.Run("server returns 500 triggers retries", func(t *testing.T) {
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := attempts.Add(1)
			if count < 3 {
				w.WriteHeader(500)
				w.Write([]byte("Internal Server Error"))
				return
			}
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`"success"`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			MaxRetries:    5,
			RetryBaseWait: time.Millisecond,
		})
		var result string
		err := tr.DoPublic(context.Background(), "/retry", nil, &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "success" {
			t.Errorf("result = %q; want %q", result, "success")
		}
		if attempts.Load() < 3 {
			t.Errorf("expected at least 3 attempts, got %d", attempts.Load())
		}
	})

	t.Run("server always returns 500 returns max retries exceeded", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			MaxRetries:    2,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/fail", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "max retries exceeded") {
			t.Errorf("error = %q; want to contain 'max retries exceeded'", err.Error())
		}
	})

	t.Run("context cancelled returns context error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			MaxRetries:    10,
			RetryBaseWait: 50 * time.Millisecond,
		})

		ctx, cancel := context.WithCancel(context.Background())
		// Cancel after a small delay so the first request goes through but retry wait is interrupted.
		go func() {
			time.Sleep(20 * time.Millisecond)
			cancel()
		}()

		err := tr.DoPublic(ctx, "/slow", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "context canceled") {
			t.Errorf("error = %q; want to contain 'context canceled'", err.Error())
		}
	})

	t.Run("network error triggers retries", func(t *testing.T) {
		// Create a server and immediately close it to produce a network error.
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		serverURL := server.URL
		server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       serverURL,
			MaxRetries:    1,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/dead", nil, nil)
		if err == nil {
			t.Fatal("expected an error for dead server")
		}
		if !strings.Contains(err.Error(), "max retries exceeded") {
			t.Errorf("error = %q; want to contain 'max retries exceeded'", err.Error())
		}
	})
}

// ---------------------------------------------------------------------------
// DoPrivateGET
// ---------------------------------------------------------------------------

func TestDoPrivateGET(t *testing.T) {
	t.Run("sets Authorization header with Bearer token", func(t *testing.T) {
		var receivedAuth string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
			TokenFunc: func() (string, error) {
				return "my-secret-token", nil
			},
		})
		tr.DoPrivateGET(context.Background(), "/private", nil, nil)
		if receivedAuth != "Bearer my-secret-token" {
			t.Errorf("Authorization = %q; want %q", receivedAuth, "Bearer my-secret-token")
		}
	})

	t.Run("sets X-Thalex-Account header when accountNumber is set", func(t *testing.T) {
		var receivedAccount string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAccount = r.Header.Get("X-Thalex-Account")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
			AccountNumber: "ACCT-42",
		})
		tr.DoPrivateGET(context.Background(), "/private", nil, nil)
		if receivedAccount != "ACCT-42" {
			t.Errorf("X-Thalex-Account = %q; want %q", receivedAccount, "ACCT-42")
		}
	})

	t.Run("sets User-Agent header", func(t *testing.T) {
		var receivedUA string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedUA = r.Header.Get("User-Agent")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			UserAgent:     "private-agent/2.0",
			RetryBaseWait: time.Millisecond,
		})
		tr.DoPrivateGET(context.Background(), "/private", nil, nil)
		if receivedUA != "private-agent/2.0" {
			t.Errorf("User-Agent = %q; want %q", receivedUA, "private-agent/2.0")
		}
	})

	t.Run("token function returning error propagates", func(t *testing.T) {
		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       "http://localhost",
			RetryBaseWait: time.Millisecond,
			TokenFunc: func() (string, error) {
				return "", fmt.Errorf("token expired")
			},
		})
		err := tr.DoPrivateGET(context.Background(), "/private", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "token expired") {
			t.Errorf("error = %q; want to contain 'token expired'", err.Error())
		}
		if !strings.Contains(err.Error(), "generating auth token") {
			t.Errorf("error = %q; want to contain 'generating auth token'", err.Error())
		}
	})

	t.Run("nil token function sets no Authorization header", func(t *testing.T) {
		var receivedAuth string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
			// TokenFunc is nil.
		})
		tr.DoPrivateGET(context.Background(), "/private", nil, nil)
		if receivedAuth != "" {
			t.Errorf("Authorization = %q; want empty", receivedAuth)
		}
	})

	t.Run("appends query params", func(t *testing.T) {
		var receivedURL string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedURL = r.URL.String()
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		params := url.Values{"currency": {"BTC"}}
		tr.DoPrivateGET(context.Background(), "/balances", params, nil)
		if !strings.Contains(receivedURL, "currency=BTC") {
			t.Errorf("URL %q does not contain currency=BTC", receivedURL)
		}
	})
}

// ---------------------------------------------------------------------------
// DoPrivatePOST
// ---------------------------------------------------------------------------

func TestDoPrivatePOST(t *testing.T) {
	t.Run("sends POST method with JSON body", func(t *testing.T) {
		var receivedMethod string
		var receivedBody []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedMethod = r.Method
			receivedBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		body := map[string]interface{}{"instrument": "BTCUSD", "amount": 1.5}
		tr.DoPrivatePOST(context.Background(), "/order", body, nil)

		if receivedMethod != http.MethodPost {
			t.Errorf("method = %q; want %q", receivedMethod, http.MethodPost)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(receivedBody, &parsed); err != nil {
			t.Fatalf("failed to parse body: %v", err)
		}
		if parsed["instrument"] != "BTCUSD" {
			t.Errorf("body.instrument = %v; want BTCUSD", parsed["instrument"])
		}
	})

	t.Run("Content-Type is application/json", func(t *testing.T) {
		var receivedCT string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedCT = r.Header.Get("Content-Type")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		tr.DoPrivatePOST(context.Background(), "/order", map[string]string{"a": "b"}, nil)
		if receivedCT != "application/json" {
			t.Errorf("Content-Type = %q; want %q", receivedCT, "application/json")
		}
	})

	t.Run("nil body sends empty JSON object", func(t *testing.T) {
		var receivedBody []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		tr.DoPrivatePOST(context.Background(), "/cancel", nil, nil)
		if string(receivedBody) != "{}" {
			t.Errorf("body = %q; want %q", string(receivedBody), "{}")
		}
	})

	t.Run("auth headers are set", func(t *testing.T) {
		var receivedAuth string
		var receivedUA string
		var receivedAccount string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedAuth = r.Header.Get("Authorization")
			receivedUA = r.Header.Get("User-Agent")
			receivedAccount = r.Header.Get("X-Thalex-Account")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`null`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			UserAgent:     "post-agent",
			RetryBaseWait: time.Millisecond,
			TokenFunc: func() (string, error) {
				return "post-token", nil
			},
			AccountNumber: "ACCT-POST",
		})
		tr.DoPrivatePOST(context.Background(), "/order", nil, nil)

		if receivedAuth != "Bearer post-token" {
			t.Errorf("Authorization = %q; want %q", receivedAuth, "Bearer post-token")
		}
		if receivedUA != "post-agent" {
			t.Errorf("User-Agent = %q; want %q", receivedUA, "post-agent")
		}
		if receivedAccount != "ACCT-POST" {
			t.Errorf("X-Thalex-Account = %q; want %q", receivedAccount, "ACCT-POST")
		}
	})

	t.Run("token error propagates", func(t *testing.T) {
		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       "http://localhost",
			RetryBaseWait: time.Millisecond,
			TokenFunc: func() (string, error) {
				return "", fmt.Errorf("refresh failed")
			},
		})
		err := tr.DoPrivatePOST(context.Background(), "/order", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "refresh failed") {
			t.Errorf("error = %q; want to contain 'refresh failed'", err.Error())
		}
	})
}

// ---------------------------------------------------------------------------
// doWithRetry (tested via public methods)
// ---------------------------------------------------------------------------

func TestDoWithRetry(t *testing.T) {
	t.Run("successful on first attempt", func(t *testing.T) {
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`"first"`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		var result string
		err := tr.DoPublic(context.Background(), "/ok", nil, &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "first" {
			t.Errorf("result = %q; want %q", result, "first")
		}
		if attempts.Load() != 1 {
			t.Errorf("attempts = %d; want 1", attempts.Load())
		}
	})

	t.Run("successful on retry after 500", func(t *testing.T) {
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := attempts.Add(1)
			if count == 1 {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`"retry_ok"`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			MaxRetries:    3,
			RetryBaseWait: time.Millisecond,
		})
		var result string
		err := tr.DoPublic(context.Background(), "/test", nil, &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != "retry_ok" {
			t.Errorf("result = %q; want %q", result, "retry_ok")
		}
		if attempts.Load() != 2 {
			t.Errorf("attempts = %d; want 2", attempts.Load())
		}
	})

	t.Run("API error in 200 response body returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{
				Error: &apiError{Code: 5000, Message: "rate limited"},
			})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/rate-limit", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		apiErr, ok := err.(*apiError)
		if !ok {
			t.Fatalf("expected *apiError, got %T: %v", err, err)
		}
		if apiErr.Code != 5000 {
			t.Errorf("Code = %d; want 5000", apiErr.Code)
		}
	})

	t.Run("result is nil skips unmarshaling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(apiResponse{Result: json.RawMessage(`{"ignored":"data"}`)})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/skip", nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid JSON in response body returns parse error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(`not json at all`))
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/bad-json", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "parsing response") {
			t.Errorf("error = %q; want to contain 'parsing response'", err.Error())
		}
	})

	t.Run("invalid JSON in result field returns parse error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			// Valid outer JSON but result is a string that can't unmarshal into a struct.
			w.Write([]byte(`{"result":"not_a_json_object"}`))
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		var result struct {
			Field string `json:"field"`
		}
		err := tr.DoPublic(context.Background(), "/bad-result", nil, &result)
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "parsing result") {
			t.Errorf("error = %q; want to contain 'parsing result'", err.Error())
		}
	})

	t.Run("4xx error with valid API error JSON in 400 range", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(apiResponse{
				Error: &apiError{Code: 422, Message: "Unprocessable"},
			})
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/unprocessable", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		apiErr, ok := err.(*apiError)
		if !ok {
			t.Fatalf("expected *apiError, got %T: %v", err, err)
		}
		if apiErr.Code != 422 {
			t.Errorf("Code = %d; want 422", apiErr.Code)
		}
	})

	t.Run("4xx errors are not retried", func(t *testing.T) {
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(403)
			w.Write([]byte("Forbidden"))
		}))
		defer server.Close()

		tr := NewHTTPTransport(HTTPTransportConfig{
			BaseURL:       server.URL,
			MaxRetries:    5,
			RetryBaseWait: time.Millisecond,
		})
		err := tr.DoPublic(context.Background(), "/forbidden", nil, nil)
		if err == nil {
			t.Fatal("expected an error")
		}
		if attempts.Load() != 1 {
			t.Errorf("attempts = %d; want 1 (4xx should not retry)", attempts.Load())
		}
	})
}
