# Error Handling

The SDK defines four structured error types, each representing a different category of failure. All error types are in the `apierr` package.

## Error Types

| Type | Description | Example Scenario |
|------|-------------|------------------|
| `*apierr.APIError` | Error returned by the Thalex API | Invalid order, insufficient margin |
| `*apierr.ConnectionError` | Connection-level failure | Network timeout, WebSocket closed |
| `*apierr.AuthError` | Authentication failure | Invalid PEM, nil key, bad credentials |
| `*apierr.TimeoutError` | Request timed out | No response within deadline |

## APIError

Represents an error returned by the Thalex API in a JSON-RPC response.

```go
type APIError struct {
    Code    int    `json:"code"`    // Numeric error code
    Message string `json:"message"` // Human-readable description
}
```

**Error string format:** `"thalex: API error <code>: <message>"`

### Checking for API Errors

Use the `IsAPIError` helper:

```go
import "github.com/amiwrpremium/go-thalex/apierr"

order, err := client.Insert(ctx, params)
if err != nil {
    if apiErr, ok := apierr.IsAPIError(err); ok {
        fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
        // Handle specific error codes.
        return
    }
    // Not an API error -- could be connection, auth, or timeout.
    log.Fatal(err)
}
```

### Using errors.As

You can also use the standard library `errors.As` pattern:

```go
import (
    "errors"

    "github.com/amiwrpremium/go-thalex/apierr"
)

var apiErr *apierr.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("Code: %d  Message: %s\n", apiErr.Code, apiErr.Message)
}
```

### Common API Error Codes

| Code | Meaning |
|------|---------|
| -32600 | Invalid request |
| -32601 | Method not found |
| -32602 | Invalid params |
| -32000 | Generic server error |

The exact error codes depend on the Thalex API. The `Message` field always contains a human-readable description.

### Example API Error Response

```json
{
    "error": {
        "code": -32602,
        "message": "instrument_name: unknown instrument 'INVALID-PERP'"
    }
}
```

## ConnectionError

Represents a connection-level failure.

```go
type ConnectionError struct {
    Message string  // Description of the connection error
    Err     error   // Underlying error (may be nil)
}
```

**Error string format:** `"thalex: connection error: <message>: <underlying>"`

### Checking for Connection Errors

```go
import (
    "errors"

    "github.com/amiwrpremium/go-thalex/apierr"
)

var connErr *apierr.ConnectionError
if errors.As(err, &connErr) {
    fmt.Printf("Connection error: %s\n", connErr.Message)
    if connErr.Err != nil {
        fmt.Printf("Caused by: %v\n", connErr.Err)
    }
}
```

### Unwrap

`ConnectionError` implements `Unwrap()` for error chain inspection:

```go
// Check if the underlying error is a specific type.
if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("Connection timed out")
}
```

### Common Connection Error Scenarios

- WebSocket connection refused
- Connection closed while waiting for response
- DNS resolution failure
- TLS handshake failure

## AuthError

Represents an authentication or credential error.

```go
type AuthError struct {
    Message string  // Description of the auth error
    Err     error   // Underlying error (may be nil)
}
```

**Error string format:** `"thalex: auth error: <message>: <underlying>"`

### Checking for Auth Errors

```go
import (
    "errors"

    "github.com/amiwrpremium/go-thalex/apierr"
)

var authErr *apierr.AuthError
if errors.As(err, &authErr) {
    fmt.Printf("Auth error: %s\n", authErr.Message)
}
```

### Common Auth Error Messages

| Message | Cause |
|---------|-------|
| `"failed to decode PEM block"` | Invalid PEM data |
| `"unsupported PEM block type: ..."` | Not an RSA key |
| `"failed to parse PKCS8 private key"` | Corrupted key |
| `"PKCS8 key is not an RSA key"` | EC or other key in PKCS8 |
| `"failed to parse RSA private key"` | Invalid PKCS1 key |
| `"private key is nil"` | Nil key in GenerateToken() |
| `"failed to sign JWT"` | Signing failure |
| `"no credentials configured"` | Login() without credentials |

## TimeoutError

Represents a timeout on a specific operation.

```go
type TimeoutError struct {
    Message string  // What timed out
    Err     error   // Underlying error (may be nil)
}
```

**Error string format:** `"thalex: timeout: <message>: <underlying>"`

### Checking for Timeout Errors

```go
import (
    "errors"

    "github.com/amiwrpremium/go-thalex/apierr"
)

var timeoutErr *apierr.TimeoutError
if errors.As(err, &timeoutErr) {
    fmt.Printf("Timeout: %s\n", timeoutErr.Message)
}
```

## Comprehensive Error Handling Pattern

Here is a complete pattern for handling all error types:

```go
import (
    "errors"
    "fmt"
    "log"

    "github.com/amiwrpremium/go-thalex/apierr"
)

func handleError(err error) {
    if err == nil {
        return
    }

    // Check API error first (most common).
    if apiErr, ok := apierr.IsAPIError(err); ok {
        fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
        return
    }

    // Check auth error.
    var authErr *apierr.AuthError
    if errors.As(err, &authErr) {
        fmt.Printf("Authentication failed: %s\n", authErr.Message)
        return
    }

    // Check connection error.
    var connErr *apierr.ConnectionError
    if errors.As(err, &connErr) {
        fmt.Printf("Connection error: %s\n", connErr.Message)
        // May want to retry or reconnect.
        return
    }

    // Check timeout.
    var timeoutErr *apierr.TimeoutError
    if errors.As(err, &timeoutErr) {
        fmt.Printf("Timed out: %s\n", timeoutErr.Message)
        return
    }

    // Check context cancellation.
    if errors.Is(err, context.Canceled) {
        fmt.Println("Request was cancelled")
        return
    }
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Context deadline exceeded")
        return
    }

    // Unknown error.
    log.Printf("Unexpected error: %v", err)
}
```

## Usage Example

```go
order, err := client.Insert(ctx,
    types.NewBuyOrderParams("BTC-PERPETUAL", 100).
        WithPrice(95000),
)
if err != nil {
    handleError(err)
    return
}
fmt.Printf("Order placed: %s\n", order.OrderID)
```

## REST Retry Behavior

The REST client automatically retries on certain errors:

| Error Type | Retried? |
|-----------|----------|
| Network errors (DNS, TCP) | Yes |
| 5xx server errors | Yes |
| 4xx client errors | No |
| API errors (invalid params, etc.) | No |
| Context cancellation | No |

Configuration:

```go
client := rest.NewClient(
    config.WithMaxRetries(5),                        // Max retry attempts
    config.WithRetryBaseWait(500 * time.Millisecond), // Base backoff
)
```

The backoff is exponential: 500ms, 1s, 2s, 4s, 8s for 5 retries.

## WebSocket Error Handling

For the WebSocket client, register an error handler for connection-level errors:

```go
wsClient.OnErrorHandler(func(err error) {
    // This fires for read errors, parse errors, etc.
    // NOT for per-request API errors.
    log.Printf("WS connection error: %v", err)
})
```

Per-request errors are returned from the method call:

```go
order, err := wsClient.Insert(ctx, params)
if err != nil {
    // This could be APIError, ConnectionError, or context error.
    handleError(err)
}
```

If the WebSocket connection drops while a request is pending:

```go
// Returns ConnectionError with message "connection closed while waiting for response"
```

## Tips

1. **Always check errors** -- never ignore the error return value.
2. **Use `IsAPIError` first** -- it is the most common error type during normal operation.
3. **Use `errors.As` for structured handling** -- lets you access error fields.
4. **Use `errors.Is` for sentinel checks** -- for `context.Canceled`, `context.DeadlineExceeded`.
5. **Enable cancel-on-disconnect** for WebSocket trading to prevent orphaned orders on connection errors.

---

[< Enum Types](enums.md) | [Home](README.md) | [Configuration >](configuration.md)
