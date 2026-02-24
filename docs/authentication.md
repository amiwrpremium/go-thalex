# Authentication

Thalex API authentication uses RSA key pairs to sign JWT tokens. You generate a key pair, register the public key with your Thalex account, and use the private key in the SDK to authenticate requests.

## RSA Key Pair Generation

Generate a 4096-bit RSA key pair using OpenSSL:

```bash
# Generate private key (PKCS1 format).
openssl genrsa -out private_key.pem 4096

# Extract the public key.
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

Alternatively, generate in PKCS8 format (both are supported):

```bash
# Generate PKCS8 private key.
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:4096

# Extract the public key.
openssl pkey -in private_key.pem -pubout -out public_key.pem
```

Upload `public_key.pem` to your Thalex account settings. You will receive a **Key ID** (the `kid` claim). Keep `private_key.pem` secret.

## Creating Credentials from PEM

The most common approach -- load a PEM file and create credentials:

```go
import (
    "os"
    "github.com/amiwrpremium/go-thalex/auth"
)

pemData, err := os.ReadFile("private_key.pem")
if err != nil {
    log.Fatal(err)
}

creds, err := auth.NewCredentialsFromPEM("my-key-id", pemData)
if err != nil {
    log.Fatal(err)
}
```

`NewCredentialsFromPEM` supports both PEM block types:
- `RSA PRIVATE KEY` (PKCS1)
- `PRIVATE KEY` (PKCS8)

Any other block type returns an `*apierr.AuthError`.

### Credentials Struct

```go
type Credentials struct {
    KeyID      string           // API key identifier (kid)
    PrivateKey *rsa.PrivateKey  // RSA private key for signing
}
```

## Creating Credentials from a Pre-Parsed Key

If you already have an `*rsa.PrivateKey` (for example, from a key vault or HSM):

```go
import (
    "crypto/rsa"
    "github.com/amiwrpremium/go-thalex/auth"
)

var privateKey *rsa.PrivateKey // from your key management system

creds := auth.NewCredentials("my-key-id", privateKey)
```

Note that `NewCredentials` returns `*Credentials` directly (no error), since the key is already parsed.

## JWT Token Generation

The SDK generates JWT tokens automatically when making authenticated requests. The token structure is:

**Header:**
```json
{
    "alg": "RS512",
    "typ": "JWT",
    "kid": "my-key-id"
}
```

**Payload:**
```json
{
    "iat": 1700000000
}
```

The token is signed using **RS512** (RSA with SHA-512). You generally never need to call `GenerateToken()` directly, but it is available:

```go
token, err := creds.GenerateToken()
if err != nil {
    log.Fatal(err)
}
fmt.Println(token) // eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCIsImtpZCI6Im15LWtleS1pZCJ9...
```

## Using Credentials with Clients

### REST Client

```go
import (
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/rest"
)

client := rest.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
)

// All private endpoints automatically include the JWT token.
summary, err := client.AccountSummary(ctx)
```

### WebSocket Client

WebSocket requires an explicit `Login()` call after connecting:

```go
import (
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/ws"
)

wsClient := ws.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
)

if err := wsClient.Connect(ctx); err != nil {
    log.Fatal(err)
}

// Must call Login() before using private endpoints.
if err := wsClient.Login(ctx); err != nil {
    log.Fatal(err)
}
```

When auto-reconnect is enabled (`WithWSReconnect(true)`), the client automatically re-authenticates and re-subscribes after a reconnection.

## Sub-Account Authentication

To authenticate as a sub-account, provide the account number:

```go
client := rest.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
    config.WithAccountNumber("ACC-12345"),
)
```

For WebSocket, the account number is included in the `Login` call automatically when configured.

## Error Type: AuthError

Authentication errors are returned as `*apierr.AuthError`:

```go
type AuthError struct {
    Message string  // Human-readable description
    Err     error   // Underlying error, if any
}
```

Example error messages:
- `"failed to decode PEM block"` -- invalid PEM data
- `"unsupported PEM block type: EC PRIVATE KEY"` -- not an RSA key
- `"failed to parse PKCS8 private key"` -- corrupted key data
- `"PKCS8 key is not an RSA key"` -- PKCS8 container holds a non-RSA key
- `"private key is nil"` -- calling `GenerateToken()` on nil key
- `"no credentials configured"` -- calling `Login()` without credentials

### Checking for Auth Errors

```go
import (
    "errors"
    "github.com/amiwrpremium/go-thalex/apierr"
)

var authErr *apierr.AuthError
if errors.As(err, &authErr) {
    fmt.Printf("Auth failed: %s\n", authErr.Message)
    if authErr.Err != nil {
        fmt.Printf("Caused by: %v\n", authErr.Err)
    }
}
```

## Security Best Practices

1. **Never commit private keys** to version control. Use environment variables or secrets management.
2. **Use testnet** for development: `config.WithNetwork(config.Testnet)`.
3. **Rotate keys** periodically via the Thalex dashboard.
4. **Use 4096-bit keys** for production use.
5. **Enable cancel-on-disconnect** for WebSocket trading sessions to prevent orphaned orders.

---

[< Getting Started](getting-started.md) | [Home](README.md) | [REST Client >](rest-client.md)
