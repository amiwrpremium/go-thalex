# go-thalex SDK Documentation

A comprehensive Go SDK for the [Thalex](https://thalex.com) cryptocurrency derivatives exchange, providing typed access to all trading, market data, and account management endpoints via both REST and WebSocket APIs.

## SDK Structure

| Package | Import Path | Description |
|---------|-------------|-------------|
| apierr | `github.com/amiwrpremium/go-thalex/apierr` | Error types (APIError, ConnectionError, AuthError, TimeoutError) |
| auth | `github.com/amiwrpremium/go-thalex/auth` | Authentication, credentials, JWT token generation |
| config | `github.com/amiwrpremium/go-thalex/config` | Configuration, network selection, client options |
| enums | `github.com/amiwrpremium/go-thalex/enums` | Enum types (Direction, OrderType, TimeInForce, etc.) |
| types | `github.com/amiwrpremium/go-thalex/types` | Request/response types, builder functions, channel helpers |
| rest | `github.com/amiwrpremium/go-thalex/rest` | REST API client |
| ws | `github.com/amiwrpremium/go-thalex/ws` | WebSocket JSON-RPC client with real-time subscriptions |

## Table of Contents

### Getting Started

- [Getting Started](getting-started.md) -- Installation, prerequisites, and minimal working examples

### Core Concepts

- [Authentication](authentication.md) -- RSA key pairs, credentials, JWT token generation
- [Configuration](configuration.md) -- All `With*` options, networks, defaults
- [Error Handling](error-handling.md) -- Error types, helpers, retry behavior

### API Clients

- [REST Client](rest-client.md) -- REST API client creation, endpoints, retry behavior
- [WebSocket Client](ws-client.md) -- WebSocket client, connection lifecycle, reconnection
- [Real-time Subscriptions](subscriptions.md) -- Public/private channels, typed handlers

### Trading

- [Trading](trading.md) -- Insert, amend, cancel orders with builder pattern
- [Market Making](market-making.md) -- Mass quotes, multi-level quotes, MM protection
- [Conditional Orders](conditional-orders.md) -- Stop, bracket, trailing stop orders
- [Bot Management](bots.md) -- SGSL, Grid, DHedge, DFollow, OCQ, Levels bots

### Reference

- [Request for Quote (RFQ)](rfq.md) -- Creating and managing RFQs
- [Enum Types](enums.md) -- All 28 enum types with values and methods
- [Examples](examples.md) -- Runnable example programs

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/rest"
)

func main() {
    client := rest.NewClient(
        config.WithNetwork(config.Testnet),
    )

    ticker, err := client.Ticker(context.Background(), "BTC-PERPETUAL")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC mark price: %.2f\n", ticker.MarkPrice)
}
```

## Version

Current SDK version: `0.2.0`

---

[Getting Started >](getting-started.md)
