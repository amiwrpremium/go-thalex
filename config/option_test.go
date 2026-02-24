package config_test

import (
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
)

func TestDefaultClientConfig(t *testing.T) {
	cfg := config.DefaultClientConfig()

	t.Run("Network", func(t *testing.T) {
		if cfg.Network != config.Production {
			t.Errorf("Network = %v, want Production", cfg.Network)
		}
	})

	t.Run("HTTPClient", func(t *testing.T) {
		if cfg.HTTPClient == nil {
			t.Fatal("HTTPClient should not be nil")
		}
		if cfg.HTTPClient.Timeout != 30*time.Second {
			t.Errorf("HTTPClient.Timeout = %v, want %v", cfg.HTTPClient.Timeout, 30*time.Second)
		}
	})

	t.Run("Credentials", func(t *testing.T) {
		if cfg.Credentials != nil {
			t.Error("Credentials should be nil by default")
		}
	})

	t.Run("Logger", func(t *testing.T) {
		if cfg.Logger != nil {
			t.Error("Logger should be nil by default")
		}
	})

	t.Run("MaxRetries", func(t *testing.T) {
		if cfg.MaxRetries != 3 {
			t.Errorf("MaxRetries = %d, want 3", cfg.MaxRetries)
		}
	})

	t.Run("RetryBaseWait", func(t *testing.T) {
		if cfg.RetryBaseWait != 500*time.Millisecond {
			t.Errorf("RetryBaseWait = %v, want %v", cfg.RetryBaseWait, 500*time.Millisecond)
		}
	})

	t.Run("WSDialTimeout", func(t *testing.T) {
		if cfg.WSDialTimeout != 10*time.Second {
			t.Errorf("WSDialTimeout = %v, want %v", cfg.WSDialTimeout, 10*time.Second)
		}
	})

	t.Run("WSPingInterval", func(t *testing.T) {
		if cfg.WSPingInterval != 5*time.Second {
			t.Errorf("WSPingInterval = %v, want %v", cfg.WSPingInterval, 5*time.Second)
		}
	})

	t.Run("WSReconnect", func(t *testing.T) {
		if cfg.WSReconnect != false {
			t.Errorf("WSReconnect = %v, want false", cfg.WSReconnect)
		}
	})

	t.Run("WSMaxReconnects", func(t *testing.T) {
		if cfg.WSMaxReconnects != 10 {
			t.Errorf("WSMaxReconnects = %d, want 10", cfg.WSMaxReconnects)
		}
	})

	t.Run("WSReconnectWait", func(t *testing.T) {
		if cfg.WSReconnectWait != 1*time.Second {
			t.Errorf("WSReconnectWait = %v, want %v", cfg.WSReconnectWait, 1*time.Second)
		}
	})

	t.Run("UserAgent", func(t *testing.T) {
		if cfg.UserAgent != config.UserAgent {
			t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, config.UserAgent)
		}
	})

	t.Run("AccountNumber", func(t *testing.T) {
		if cfg.AccountNumber != "" {
			t.Errorf("AccountNumber = %q, want empty", cfg.AccountNumber)
		}
	})
}

func TestWithNetwork(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithNetwork(config.Testnet)(&cfg)

	if cfg.Network != config.Testnet {
		t.Errorf("Network = %v, want Testnet", cfg.Network)
	}
}

func TestWithCredentials(t *testing.T) {
	cfg := config.DefaultClientConfig()
	creds := auth.NewCredentials("kid", nil)
	config.WithCredentials(creds)(&cfg)

	if cfg.Credentials != creds {
		t.Error("Credentials should match the provided credentials")
	}
}

func TestWithHTTPClient(t *testing.T) {
	cfg := config.DefaultClientConfig()
	customClient := &http.Client{Timeout: 60 * time.Second}
	config.WithHTTPClient(customClient)(&cfg)

	if cfg.HTTPClient != customClient {
		t.Error("HTTPClient should match the provided client")
	}
	if cfg.HTTPClient.Timeout != 60*time.Second {
		t.Errorf("HTTPClient.Timeout = %v, want %v", cfg.HTTPClient.Timeout, 60*time.Second)
	}
}

func TestWithLogger(t *testing.T) {
	cfg := config.DefaultClientConfig()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	config.WithLogger(logger)(&cfg)

	if cfg.Logger != logger {
		t.Error("Logger should match the provided logger")
	}
}

func TestWithMaxRetries(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithMaxRetries(5)(&cfg)

	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", cfg.MaxRetries)
	}
}

func TestWithRetryBaseWait(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithRetryBaseWait(2 * time.Second)(&cfg)

	if cfg.RetryBaseWait != 2*time.Second {
		t.Errorf("RetryBaseWait = %v, want %v", cfg.RetryBaseWait, 2*time.Second)
	}
}

func TestWithWSDialTimeout(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithWSDialTimeout(30 * time.Second)(&cfg)

	if cfg.WSDialTimeout != 30*time.Second {
		t.Errorf("WSDialTimeout = %v, want %v", cfg.WSDialTimeout, 30*time.Second)
	}
}

func TestWithWSPingInterval(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithWSPingInterval(15 * time.Second)(&cfg)

	if cfg.WSPingInterval != 15*time.Second {
		t.Errorf("WSPingInterval = %v, want %v", cfg.WSPingInterval, 15*time.Second)
	}
}

func TestWithWSReconnect(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithWSReconnect(true)(&cfg)

	if cfg.WSReconnect != true {
		t.Errorf("WSReconnect = %v, want true", cfg.WSReconnect)
	}

	config.WithWSReconnect(false)(&cfg)
	if cfg.WSReconnect != false {
		t.Errorf("WSReconnect = %v, want false after setting back", cfg.WSReconnect)
	}
}

func TestWithWSMaxReconnects(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithWSMaxReconnects(20)(&cfg)

	if cfg.WSMaxReconnects != 20 {
		t.Errorf("WSMaxReconnects = %d, want 20", cfg.WSMaxReconnects)
	}
}

func TestWithWSReconnectWait(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithWSReconnectWait(5 * time.Second)(&cfg)

	if cfg.WSReconnectWait != 5*time.Second {
		t.Errorf("WSReconnectWait = %v, want %v", cfg.WSReconnectWait, 5*time.Second)
	}
}

func TestWithAccountNumber(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithAccountNumber("acc-123")(&cfg)

	if cfg.AccountNumber != "acc-123" {
		t.Errorf("AccountNumber = %q, want %q", cfg.AccountNumber, "acc-123")
	}
}

func TestWithUserAgent(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithUserAgent("custom-agent/1.0")(&cfg)

	if cfg.UserAgent != "custom-agent/1.0" {
		t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, "custom-agent/1.0")
	}
}

func TestApplyMultipleOptions(t *testing.T) {
	cfg := config.DefaultClientConfig()
	options := []config.ClientOption{
		config.WithNetwork(config.Testnet),
		config.WithMaxRetries(7),
		config.WithWSReconnect(true),
		config.WithWSMaxReconnects(50),
		config.WithUserAgent("multi-test/2.0"),
		config.WithAccountNumber("acc-multi"),
		config.WithRetryBaseWait(1 * time.Second),
		config.WithWSDialTimeout(20 * time.Second),
		config.WithWSPingInterval(10 * time.Second),
		config.WithWSReconnectWait(3 * time.Second),
	}

	for _, opt := range options {
		opt(&cfg)
	}

	if cfg.Network != config.Testnet {
		t.Errorf("Network = %v, want Testnet", cfg.Network)
	}
	if cfg.MaxRetries != 7 {
		t.Errorf("MaxRetries = %d, want 7", cfg.MaxRetries)
	}
	if cfg.WSReconnect != true {
		t.Errorf("WSReconnect = %v, want true", cfg.WSReconnect)
	}
	if cfg.WSMaxReconnects != 50 {
		t.Errorf("WSMaxReconnects = %d, want 50", cfg.WSMaxReconnects)
	}
	if cfg.UserAgent != "multi-test/2.0" {
		t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, "multi-test/2.0")
	}
	if cfg.AccountNumber != "acc-multi" {
		t.Errorf("AccountNumber = %q, want %q", cfg.AccountNumber, "acc-multi")
	}
	if cfg.RetryBaseWait != 1*time.Second {
		t.Errorf("RetryBaseWait = %v, want %v", cfg.RetryBaseWait, 1*time.Second)
	}
	if cfg.WSDialTimeout != 20*time.Second {
		t.Errorf("WSDialTimeout = %v, want %v", cfg.WSDialTimeout, 20*time.Second)
	}
	if cfg.WSPingInterval != 10*time.Second {
		t.Errorf("WSPingInterval = %v, want %v", cfg.WSPingInterval, 10*time.Second)
	}
	if cfg.WSReconnectWait != 3*time.Second {
		t.Errorf("WSReconnectWait = %v, want %v", cfg.WSReconnectWait, 3*time.Second)
	}
}

func TestOptionsOverrideEachOther(t *testing.T) {
	cfg := config.DefaultClientConfig()

	config.WithMaxRetries(5)(&cfg)
	if cfg.MaxRetries != 5 {
		t.Fatalf("MaxRetries after first set = %d, want 5", cfg.MaxRetries)
	}

	config.WithMaxRetries(10)(&cfg)
	if cfg.MaxRetries != 10 {
		t.Errorf("MaxRetries after override = %d, want 10", cfg.MaxRetries)
	}
}

func TestWithCredentials_Nil(t *testing.T) {
	cfg := config.DefaultClientConfig()
	creds := auth.NewCredentials("kid", nil)
	config.WithCredentials(creds)(&cfg)

	// Now set back to nil
	config.WithCredentials(nil)(&cfg)
	if cfg.Credentials != nil {
		t.Error("Credentials should be nil after setting nil")
	}
}

func TestWithHTTPClient_Nil(t *testing.T) {
	cfg := config.DefaultClientConfig()
	config.WithHTTPClient(nil)(&cfg)
	if cfg.HTTPClient != nil {
		t.Error("HTTPClient should be nil after setting nil")
	}
}

func TestWithLogger_Nil(t *testing.T) {
	cfg := config.DefaultClientConfig()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	config.WithLogger(logger)(&cfg)
	if cfg.Logger == nil {
		t.Fatal("Logger should not be nil after setting")
	}
	config.WithLogger(nil)(&cfg)
	if cfg.Logger != nil {
		t.Error("Logger should be nil after setting nil")
	}
}
