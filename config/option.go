package config

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
)

// ClientConfig holds all configurable options for the SDK clients.
type ClientConfig struct {
	Network         Network
	Credentials     *auth.Credentials
	HTTPClient      *http.Client
	Logger          *slog.Logger
	MaxRetries      int
	RetryBaseWait   time.Duration
	WSDialTimeout   time.Duration
	WSPingInterval  time.Duration
	WSReconnect     bool
	WSMaxReconnects int
	WSReconnectWait time.Duration
	AccountNumber   string
	UserAgent       string
}

// DefaultClientConfig returns sensible defaults.
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Network:         Production,
		HTTPClient:      &http.Client{Timeout: 30 * time.Second},
		MaxRetries:      3,
		RetryBaseWait:   500 * time.Millisecond,
		WSDialTimeout:   10 * time.Second,
		WSPingInterval:  5 * time.Second,
		WSReconnect:     false,
		WSMaxReconnects: 10,
		WSReconnectWait: 1 * time.Second,
		UserAgent:       UserAgent,
	}
}

// ClientOption configures the SDK clients.
type ClientOption func(*ClientConfig)

// WithNetwork sets the Thalex network (Production or Testnet).
func WithNetwork(n Network) ClientOption {
	return func(c *ClientConfig) { c.Network = n }
}

// WithCredentials sets the API credentials for authentication.
func WithCredentials(creds *auth.Credentials) ClientOption {
	return func(c *ClientConfig) { c.Credentials = creds }
}

// WithHTTPClient sets a custom HTTP client for REST requests.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *ClientConfig) { c.HTTPClient = client }
}

// WithLogger sets a structured logger for the SDK.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *ClientConfig) { c.Logger = logger }
}

// WithMaxRetries sets the maximum number of retry attempts for failed requests.
func WithMaxRetries(n int) ClientOption {
	return func(c *ClientConfig) { c.MaxRetries = n }
}

// WithRetryBaseWait sets the base wait duration between retries (exponential backoff).
func WithRetryBaseWait(d time.Duration) ClientOption {
	return func(c *ClientConfig) { c.RetryBaseWait = d }
}

// WithWSDialTimeout sets the timeout for establishing WebSocket connections.
func WithWSDialTimeout(d time.Duration) ClientOption {
	return func(c *ClientConfig) { c.WSDialTimeout = d }
}

// WithWSPingInterval sets the interval between WebSocket ping frames.
func WithWSPingInterval(d time.Duration) ClientOption {
	return func(c *ClientConfig) { c.WSPingInterval = d }
}

// WithWSReconnect enables or disables automatic WebSocket reconnection.
func WithWSReconnect(enabled bool) ClientOption {
	return func(c *ClientConfig) { c.WSReconnect = enabled }
}

// WithWSMaxReconnects sets the maximum number of reconnection attempts.
func WithWSMaxReconnects(n int) ClientOption {
	return func(c *ClientConfig) { c.WSMaxReconnects = n }
}

// WithWSReconnectWait sets the base wait duration between reconnection attempts.
func WithWSReconnectWait(d time.Duration) ClientOption {
	return func(c *ClientConfig) { c.WSReconnectWait = d }
}

// WithAccountNumber sets the default account number for private endpoints.
func WithAccountNumber(accountNumber string) ClientOption {
	return func(c *ClientConfig) { c.AccountNumber = accountNumber }
}

// WithUserAgent overrides the default user agent string.
func WithUserAgent(ua string) ClientOption {
	return func(c *ClientConfig) { c.UserAgent = ua }
}
