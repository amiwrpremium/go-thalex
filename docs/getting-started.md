# Getting Started

This guide walks you through installing the go-thalex SDK, setting up credentials, and making your first API calls.

## Prerequisites

- **Go 1.21+** (the module uses Go 1.25, but is compatible with 1.21+)
- **RSA key pair** registered with your Thalex account (for authenticated endpoints)
- A Thalex account at [thalex.com](https://thalex.com) or the [testnet](https://testnet.thalex.com)

## Installation

```bash
go get github.com/amiwrpremium/go-thalex
```

This installs the root package and all sub-packages (`apierr`, `auth`, `config`, `enums`, `types`, `rest`, `ws`).

## Project Setup

A typical import block for a go-thalex project:

```go
import (
    "github.com/amiwrpremium/go-thalex/apierr"
    "github.com/amiwrpremium/go-thalex/auth"
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
    "github.com/amiwrpremium/go-thalex/rest"
    "github.com/amiwrpremium/go-thalex/ws"
)
```

The `config` package provides client options and network constants, `auth` provides credential management, and `apierr` provides typed API errors.

## Minimal REST Example

Fetch public market data without authentication:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/rest"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client := rest.NewClient(
        config.WithNetwork(config.Testnet),
    )

    // Fetch all active instruments.
    instruments, err := client.Instruments(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d active instruments\n", len(instruments))

    // Fetch a ticker.
    ticker, err := client.Ticker(ctx, "BTC-PERPETUAL")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC-PERPETUAL mark price: %.2f\n", ticker.MarkPrice)
}
```

## Minimal REST Example (Authenticated)

Place and cancel an order:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/amiwrpremium/go-thalex/auth"
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/rest"
    "github.com/amiwrpremium/go-thalex/types"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Load credentials from PEM file.
    pemData, err := os.ReadFile("private_key.pem")
    if err != nil {
        log.Fatal(err)
    }
    creds, err := auth.NewCredentialsFromPEM("my-key-id", pemData)
    if err != nil {
        log.Fatal(err)
    }

    client := rest.NewClient(
        config.WithNetwork(config.Testnet),
        config.WithCredentials(creds),
    )

    // Place a limit buy order far below market.
    order, err := client.Insert(ctx,
        types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
            WithPrice(30000).
            WithPostOnly(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Order placed: %s status=%s\n", order.OrderID, order.Status)

    // Cancel it.
    cancelled, err := client.Cancel(ctx, types.CancelByOrderID(order.OrderID))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Cancelled: %s\n", cancelled.OrderID)
}
```

## Minimal WebSocket Example

Stream real-time ticker data:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"

    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
    "github.com/amiwrpremium/go-thalex/ws"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    wsClient := ws.NewClient(
        config.WithNetwork(config.Testnet),
        config.WithWSReconnect(true),
    )

    wsClient.OnErrorHandler(func(err error) {
        log.Printf("WS error: %v", err)
    })

    if err := wsClient.Connect(ctx); err != nil {
        log.Fatal(err)
    }
    defer wsClient.Close()

    // Register handler then subscribe.
    ch := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
    wsClient.OnTicker(ch, func(t types.Ticker) {
        fmt.Printf("BTC mark=%.2f\n", t.MarkPrice)
    })

    if err := wsClient.Subscribe(ctx, ch); err != nil {
        log.Fatal(err)
    }

    <-ctx.Done()
}
```

## WebSocket Trading Example

Connect, authenticate, and place orders over WebSocket:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/amiwrpremium/go-thalex/auth"
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/types"
    "github.com/amiwrpremium/go-thalex/ws"
)

func main() {
    pemData, _ := os.ReadFile("private_key.pem")
    creds, _ := auth.NewCredentialsFromPEM("my-key-id", pemData)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    wsClient := ws.NewClient(
        config.WithNetwork(config.Testnet),
        config.WithCredentials(creds),
        config.WithWSReconnect(true),
    )

    if err := wsClient.Connect(ctx); err != nil {
        log.Fatal(err)
    }
    defer wsClient.Close()

    // Authenticate.
    if err := wsClient.Login(ctx); err != nil {
        log.Fatal(err)
    }

    // Enable cancel-on-disconnect for safety.
    wsClient.SetCancelOnDisconnect(ctx, true)

    // Place and cancel an order.
    order, err := wsClient.Insert(ctx,
        types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
            WithPrice(30000).
            WithPostOnly(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Order: %s status=%s\n", order.OrderID, order.Status)

    wsClient.Cancel(ctx, types.CancelByOrderID(order.OrderID))
}
```

## Environment Variables

The example programs use these environment variables:

| Variable | Description |
|----------|-------------|
| `THALEX_KEY_ID` | Your API key identifier |
| `THALEX_PEM_PATH` | Path to your RSA private key PEM file |

## Next Steps

- [Authentication](authentication.md) -- Generate and configure RSA key pairs
- [REST Client](rest-client.md) -- Full REST API reference
- [WebSocket Client](ws-client.md) -- Real-time WebSocket API
- [Configuration](configuration.md) -- All configuration options

---

[< Home](README.md) | [Authentication >](authentication.md)
