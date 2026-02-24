// Package transport provides HTTP and WebSocket transport implementations
// for the Thalex API.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"
)

// HTTPTransport handles HTTP communication with the Thalex REST API.
type HTTPTransport struct {
	client        *http.Client
	baseURL       string
	userAgent     string
	maxRetries    int
	retryBaseWait time.Duration
	tokenFunc     func() (string, error)
	accountNumber string
}

// HTTPTransportConfig contains configuration for the HTTP transport.
type HTTPTransportConfig struct {
	Client        *http.Client
	BaseURL       string
	UserAgent     string
	MaxRetries    int
	RetryBaseWait time.Duration
	TokenFunc     func() (string, error)
	AccountNumber string
}

// NewHTTPTransport creates a new HTTP transport.
func NewHTTPTransport(cfg HTTPTransportConfig) *HTTPTransport {
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: 30 * time.Second}
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryBaseWait <= 0 {
		cfg.RetryBaseWait = 500 * time.Millisecond
	}
	return &HTTPTransport{
		client:        cfg.Client,
		baseURL:       cfg.BaseURL,
		userAgent:     cfg.UserAgent,
		maxRetries:    cfg.MaxRetries,
		retryBaseWait: cfg.RetryBaseWait,
		tokenFunc:     cfg.TokenFunc,
		accountNumber: cfg.AccountNumber,
	}
}

// apiResponse wraps the Thalex REST API response format.
type apiResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *apiError       `json:"error"`
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// DoPublic performs a public (unauthenticated) GET request.
func (t *HTTPTransport) DoPublic(ctx context.Context, path string, queryParams url.Values, result interface{}) error {
	u := t.baseURL + path
	if len(queryParams) > 0 {
		u += "?" + queryParams.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", t.userAgent)

	return t.doWithRetry(req, result)
}

// DoPrivateGET performs an authenticated GET request.
func (t *HTTPTransport) DoPrivateGET(ctx context.Context, path string, queryParams url.Values, result interface{}) error {
	u := t.baseURL + path
	if len(queryParams) > 0 {
		u += "?" + queryParams.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := t.setAuthHeaders(req); err != nil {
		return err
	}

	return t.doWithRetry(req, result)
}

// DoPrivatePOST performs an authenticated POST request with a JSON body.
func (t *HTTPTransport) DoPrivatePOST(ctx context.Context, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	} else {
		bodyReader = bytes.NewReader([]byte("{}"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := t.setAuthHeaders(req); err != nil {
		return err
	}

	return t.doWithRetry(req, result)
}

func (t *HTTPTransport) setAuthHeaders(req *http.Request) error {
	req.Header.Set("User-Agent", t.userAgent)

	if t.tokenFunc != nil {
		token, err := t.tokenFunc()
		if err != nil {
			return fmt.Errorf("generating auth token: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	if t.accountNumber != "" {
		req.Header.Set("X-Thalex-Account", t.accountNumber)
	}

	return nil
}

func (t *HTTPTransport) doWithRetry(req *http.Request, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		if attempt > 0 {
			wait := t.retryBaseWait * time.Duration(math.Pow(2, float64(attempt-1)))
			select {
			case <-req.Context().Done():
				return req.Context().Err()
			case <-time.After(wait):
			}
		}

		resp, err := t.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			// Only retry on network-level errors, not on context cancellation.
			if req.Context().Err() != nil {
				return lastErr
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("reading response body: %w", err)
			continue
		}

		// Retry on 5xx server errors.
		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error: HTTP %d", resp.StatusCode)
			continue
		}

		// Don't retry on 4xx client errors.
		if resp.StatusCode >= 400 {
			var apiResp apiResponse
			if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Error != nil {
				return apiResp.Error
			}
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		// Parse the response.
		var apiResp apiResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return fmt.Errorf("parsing response: %w", err)
		}

		if apiResp.Error != nil {
			return apiResp.Error
		}

		if result != nil {
			if err := json.Unmarshal(apiResp.Result, result); err != nil {
				return fmt.Errorf("parsing result: %w", err)
			}
		}

		return nil
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}
