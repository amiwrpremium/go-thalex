package rest

import (
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/internal/transport"
)

// Client provides access to the Thalex REST API.
type Client struct {
	transport *transport.HTTPTransport
	cfg       config.ClientConfig
}

// NewClient creates a new REST API client.
func NewClient(opts ...config.ClientOption) *Client {
	cfg := config.DefaultClientConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	var tokenFunc func() (string, error)
	if cfg.Credentials != nil {
		tokenFunc = cfg.Credentials.GenerateToken
	}
	t := transport.NewHTTPTransport(transport.HTTPTransportConfig{
		Client:        cfg.HTTPClient,
		BaseURL:       cfg.Network.BaseURL(),
		UserAgent:     cfg.UserAgent,
		MaxRetries:    cfg.MaxRetries,
		RetryBaseWait: cfg.RetryBaseWait,
		TokenFunc:     tokenFunc,
		AccountNumber: cfg.AccountNumber,
	})
	return &Client{transport: t, cfg: cfg}
}
