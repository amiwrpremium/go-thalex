# Configuration

Both the REST and WebSocket clients are configured using functional options passed to their constructors. All options are defined in the `config` package.

## Option Pattern

```go
import "github.com/amiwrpremium/go-thalex/config"

client := rest.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
    config.WithMaxRetries(5),
)

wsClient := ws.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
    config.WithWSReconnect(true),
    config.WithWSPingInterval(3 * time.Second),
)
```

Options are of type `config.ClientOption` (which is `func(*ClientConfig)`). They modify the shared `ClientConfig` struct.

## Full Option Reference

### Network

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithNetwork(n)` | `config.Network` | `Production` | Thalex API environment |

**Network values:**

| Constant | REST Base URL | WebSocket URL |
|----------|---------------|---------------|
| `config.Production` | `https://thalex.com/api/v2` | `wss://thalex.com/ws/api/v2` |
| `config.Testnet` | `https://testnet.thalex.com/api/v2` | `wss://testnet.thalex.com/ws/api/v2` |

```go
// Production (default).
config.WithNetwork(config.Production)

// Testnet for development.
config.WithNetwork(config.Testnet)
```

The `Network` type also provides URL accessor methods:

```go
n := config.Testnet
fmt.Println(n.BaseURL())      // "https://testnet.thalex.com/api/v2"
fmt.Println(n.WebSocketURL()) // "wss://testnet.thalex.com/ws/api/v2"
fmt.Println(n.String())       // "testnet"
```

### Authentication

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithCredentials(c)` | `*auth.Credentials` | `nil` | API credentials for authentication |
| `WithAccountNumber(a)` | `string` | `""` | Sub-account number |

```go
// From PEM file.
creds, err := auth.NewCredentialsFromPEM("key-id", pemData)
config.WithCredentials(creds)

// Sub-account.
config.WithAccountNumber("ACC-12345")
```

See [Authentication](authentication.md) for details on creating credentials.

### Common Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithUserAgent(ua)` | `string` | `"go-thalex/0.2.0"` | Custom user agent string |
| `WithLogger(l)` | `*slog.Logger` | `nil` | Structured logger |

```go
import "log/slog"

config.WithUserAgent("my-bot/1.0.0")

config.WithLogger(slog.Default())
```

### REST-Specific Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithHTTPClient(c)` | `*http.Client` | 30s timeout | Custom HTTP client |
| `WithMaxRetries(n)` | `int` | `3` | Maximum retry attempts for failed requests |
| `WithRetryBaseWait(d)` | `time.Duration` | `500ms` | Base wait between retries (exponential backoff) |

```go
import "net/http"

// Custom HTTP client with proxy.
config.WithHTTPClient(&http.Client{
    Timeout:   60 * time.Second,
    Transport: &http.Transport{
        Proxy: http.ProxyFromEnvironment,
    },
})

// More retries with longer wait.
config.WithMaxRetries(10)
config.WithRetryBaseWait(1 * time.Second)
```

### WebSocket-Specific Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithWSDialTimeout(d)` | `time.Duration` | `10s` | Timeout for establishing WebSocket connection |
| `WithWSPingInterval(d)` | `time.Duration` | `5s` | Interval between WebSocket ping frames |
| `WithWSReconnect(b)` | `bool` | `false` | Enable automatic reconnection |
| `WithWSMaxReconnects(n)` | `int` | `10` | Maximum reconnection attempts |
| `WithWSReconnectWait(d)` | `time.Duration` | `1s` | Base wait between reconnection attempts |

```go
// Production-ready WebSocket configuration.
config.WithWSDialTimeout(15 * time.Second)
config.WithWSPingInterval(3 * time.Second)
config.WithWSReconnect(true)
config.WithWSMaxReconnects(50)
config.WithWSReconnectWait(2 * time.Second)
```

## ClientConfig Struct

All options modify the `ClientConfig` struct:

```go
type ClientConfig struct {
    Network         Network          // API environment
    Credentials     *Credentials     // Authentication credentials
    HTTPClient      *http.Client     // Custom HTTP client (REST)
    Logger          *slog.Logger     // Structured logger
    MaxRetries      int              // Max retry attempts (REST)
    RetryBaseWait   time.Duration    // Base backoff (REST)
    WSDialTimeout   time.Duration    // Connection timeout (WS)
    WSPingInterval  time.Duration    // Ping interval (WS)
    WSReconnect     bool             // Auto-reconnect (WS)
    WSMaxReconnects int              // Max reconnect attempts (WS)
    WSReconnectWait time.Duration    // Reconnect backoff (WS)
    AccountNumber   string           // Sub-account number
    UserAgent       string           // User agent string
}
```

## Default Configuration

The `DefaultClientConfig()` function returns sensible defaults:

```go
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
        UserAgent:       "go-thalex/0.2.0",
    }
}
```

## Configuration Examples

### Minimal Public Client

```go
client := rest.NewClient(
    config.WithNetwork(config.Testnet),
)
```

### Production Trading Client

```go
client := rest.NewClient(
    config.WithNetwork(config.Production),
    config.WithCredentials(creds),
    config.WithMaxRetries(5),
    config.WithRetryBaseWait(1 * time.Second),
    config.WithUserAgent("my-trading-bot/2.0"),
)
```

### Market Making WebSocket Client

```go
wsClient := ws.NewClient(
    config.WithNetwork(config.Production),
    config.WithCredentials(creds),
    config.WithWSReconnect(true),
    config.WithWSMaxReconnects(100),
    config.WithWSReconnectWait(500 * time.Millisecond),
    config.WithWSPingInterval(3 * time.Second),
    config.WithWSDialTimeout(5 * time.Second),
)
```

### Sub-Account Access

```go
client := rest.NewClient(
    config.WithNetwork(config.Production),
    config.WithCredentials(creds),
    config.WithAccountNumber("ACC-12345"),
)
```

### Custom HTTP Transport

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
}

client := rest.NewClient(
    config.WithHTTPClient(&http.Client{
        Timeout:   30 * time.Second,
        Transport: transport,
    }),
)
```

### With Structured Logging

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client := rest.NewClient(
    config.WithLogger(logger),
)
```

## Which Options Apply Where?

| Option | REST | WebSocket |
|--------|------|-----------|
| `WithNetwork` | Yes | Yes |
| `WithCredentials` | Yes | Yes |
| `WithAccountNumber` | Yes | Yes |
| `WithUserAgent` | Yes | Yes |
| `WithLogger` | Yes | Yes |
| `WithHTTPClient` | Yes | No |
| `WithMaxRetries` | Yes | No |
| `WithRetryBaseWait` | Yes | No |
| `WithWSDialTimeout` | No | Yes |
| `WithWSPingInterval` | No | Yes |
| `WithWSReconnect` | No | Yes |
| `WithWSMaxReconnects` | No | Yes |
| `WithWSReconnectWait` | No | Yes |

All options can be passed to either client constructor without error, but WebSocket-specific options have no effect on the REST client and vice versa.

---

[< Error Handling](error-handling.md) | [Home](README.md) | [Examples >](examples.md)
