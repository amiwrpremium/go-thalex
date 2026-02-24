# go-thalex

[![Go Reference](https://pkg.go.dev/badge/github.com/amiwrpremium/go-thalex.svg)](https://pkg.go.dev/github.com/amiwrpremium/go-thalex)
[![Go Report Card](https://goreportcard.com/badge/github.com/amiwrpremium/go-thalex)](https://goreportcard.com/report/github.com/amiwrpremium/go-thalex)
[![CI](https://github.com/amiwrpremium/go-thalex/actions/workflows/ci.yml/badge.svg)](https://github.com/amiwrpremium/go-thalex/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/amiwrpremium/go-thalex/branch/master/graph/badge.svg)](https://codecov.io/gh/amiwrpremium/go-thalex)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/amiwrpremium/go-thalex)](https://github.com/amiwrpremium/go-thalex/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/amiwrpremium/go-thalex)](https://go.dev/)
[![GitHub Stars](https://img.shields.io/github/stars/amiwrpremium/go-thalex?style=social)](https://github.com/amiwrpremium/go-thalex/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/amiwrpremium/go-thalex)](https://github.com/amiwrpremium/go-thalex/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/amiwrpremium/go-thalex)](https://github.com/amiwrpremium/go-thalex/pulls)
[![Last Commit](https://img.shields.io/github/last-commit/amiwrpremium/go-thalex)](https://github.com/amiwrpremium/go-thalex/commits/master)

A comprehensive Go SDK for the [Thalex](https://thalex.com) cryptocurrency derivatives exchange.

## Features

- **REST Client** — Full HTTP API coverage for trading, market data, account management, and more
- **WebSocket Client** — Low-latency JSON-RPC API with automatic reconnection
- **Real-time Subscriptions** — Typed handlers for order books, tickers, trades, portfolio updates, and all other channels
- **Builder Pattern** — Fluent API for constructing orders, conditional orders, bots, and mass quotes
- **Minimal Dependencies** — Only [`github.com/gorilla/websocket`](https://github.com/gorilla/websocket) for WebSocket; JWT authentication uses stdlib crypto
- **Type Safety** — Fully typed request/response structs, enums with validators, and helper methods

## Installation

```bash
go get github.com/amiwrpremium/go-thalex
```

Requires Go 1.21+.

## Package Structure

```
github.com/amiwrpremium/go-thalex/apierr   — Error types (APIError, AuthError, etc.)
github.com/amiwrpremium/go-thalex/auth     — API key credentials and JWT token generation
github.com/amiwrpremium/go-thalex/config   — Network selection and client configuration options
github.com/amiwrpremium/go-thalex/enums    — Enum types (Direction, OrderType, etc.)
github.com/amiwrpremium/go-thalex/types    — Request/response types
github.com/amiwrpremium/go-thalex/rest     — REST API client
github.com/amiwrpremium/go-thalex/ws       — WebSocket JSON-RPC client with subscriptions
```

## Quick Start

### Authentication

Generate an RSA key pair and register it with Thalex to get your API key ID.

```go
import "github.com/amiwrpremium/go-thalex/auth"

creds, err := auth.NewCredentialsFromPEM("my-key-id", pemData)
if err != nil {
    log.Fatal(err)
}
```

### REST API

```go
import (
    "github.com/amiwrpremium/go-thalex/auth"
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/rest"
    "github.com/amiwrpremium/go-thalex/types"
)

client := rest.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
)

// Public market data
instruments, err := client.Instruments(ctx)
ticker, err := client.Ticker(ctx, "BTC-PERPETUAL")
book, err := client.Book(ctx, "BTC-PERPETUAL")

// Place an order
order, err := client.Insert(ctx,
    types.NewBuyOrderParams("BTC-PERPETUAL", 0.1).
        WithPrice(45000).
        WithPostOnly(true),
)

// Cancel all orders
n, err := client.CancelAll(ctx)
```

### WebSocket API

```go
import (
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/types"
    "github.com/amiwrpremium/go-thalex/ws"
)

wsClient := ws.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
    config.WithWSReconnect(true),
)

if err := wsClient.Connect(ctx); err != nil {
    log.Fatal(err)
}
defer wsClient.Close()

// Login
if err := wsClient.Login(ctx); err != nil {
    log.Fatal(err)
}

// Place orders via WebSocket
order, err := wsClient.Insert(ctx,
    types.NewSellOrderParams("ETH-PERPETUAL", 1.0).
        WithPrice(3500).
        WithLabel("ws-order"),
)
```

### Real-time Subscriptions

```go
import (
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
)

// Register typed handlers
ch := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
wsClient.OnTicker(ch, func(t types.Ticker) {
    fmt.Printf("BTC mark: %.2f\n", t.MarkPrice)
})

// Subscribe to the channel
wsClient.Subscribe(ctx, ch)

// Private subscriptions
wsClient.OnOrders(func(orders []types.OrderStatus) {
    for _, o := range orders {
        fmt.Printf("Order %s: %s\n", o.OrderID, o.Status)
    }
})
wsClient.SubscribePrivate(ctx, types.ChannelAccountOrders)
```

### Market Making

```go
// Mass quote (WebSocket only)
result, err := wsClient.MassQuote(ctx,
    types.NewMassQuoteParams([]types.DoubleSidedQuote{
        types.NewSingleLevelQuote("BTC-PERPETUAL", 44900, 1.0, 45100, 1.0),
        types.NewSingleLevelQuote("ETH-PERPETUAL", 3490, 10.0, 3510, 10.0),
    }).WithPostOnly(true),
)

// Set market maker protection
wsClient.SetMMProtection(ctx, &types.MMProtectionParams{
    Product:     "FBTCUSD",
    TradeAmount: 10.0,
    QuoteAmount: 50.0,
})
```

### Bots

```go
import (
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
)

// Create an SGSL bot
bot, err := client.CreateSGSLBot(ctx,
    types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetMark,
        45000, 1.0, 44000, 0.0, endTime).
        WithMaxSlippage(100).
        WithLabel("my-bot"),
)

// Create a Grid bot
bot, err := client.CreateGridBot(ctx,
    types.NewGridBotParams("BTC-PERPETUAL",
        []float64{44000, 44500, 45000, 45500, 46000}, 0.1, endTime).
        WithBasePosition(0.5),
)
```

### Conditional Orders

```go
import (
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
)

// Stop order
order, err := client.CreateConditionalOrder(ctx,
    types.NewStopOrder(enums.DirectionSell, "BTC-PERPETUAL", 0.1, 40000).
        WithTarget(enums.TargetMark).
        WithReduceOnly(true),
)

// Bracket order
order, err := client.CreateConditionalOrder(ctx,
    types.NewBracketOrder(enums.DirectionSell, "BTC-PERPETUAL", 0.1, 40000, 50000),
)
```

## Configuration

All client options are in the `config` package and work with both `rest.Client` and `ws.Client`:

| Option | Description | Default |
|--------|-------------|---------|
| `WithNetwork(n)` | Production or Testnet | Production |
| `WithCredentials(c)` | API credentials | nil |
| `WithHTTPClient(c)` | Custom HTTP client | 30s timeout |
| `WithLogger(l)` | Structured logger (`slog.Logger`) | nil |
| `WithMaxRetries(n)` | HTTP retry attempts | 3 |
| `WithRetryBaseWait(d)` | Base wait between retries | 500ms |
| `WithWSDialTimeout(d)` | WebSocket dial timeout | 10s |
| `WithWSPingInterval(d)` | WebSocket ping interval | 5s |
| `WithWSReconnect(b)` | Auto-reconnect on disconnect | false |
| `WithWSMaxReconnects(n)` | Max reconnection attempts | 10 |
| `WithWSReconnectWait(d)` | Base wait between reconnects | 1s |
| `WithAccountNumber(s)` | Default account number | "" |
| `WithUserAgent(ua)` | Custom user agent string | "go-thalex/0.2.0" |

## Subscription Channels

### Public Channels

| Helper | Example |
|--------|---------|
| `types.BookChannel(instrument, grouping, nlevels, delay)` | `book.BTC-PERPETUAL.1.10.100ms` |
| `types.TickerChannel(instrument, delay)` | `ticker.BTC-PERPETUAL.100ms` |
| `types.LWTChannel(instrument, delay)` | `lwt.BTC-PERPETUAL.100ms` |
| `types.RecentTradesChannel(target, category)` | `recent_trades.BTCUSD.all` |
| `types.PriceIndexChannel(underlying)` | `price_index.BTCUSD` |
| `types.BasePriceChannel(underlying, expiration)` | `base_price.BTCUSD.2025-03-28` |
| `types.IndexComponentsChannel(underlying)` | `index_components.BTCUSD` |

### Private Channels (Constants)

| Constant | Description |
|----------|-------------|
| `types.ChannelAccountOrders` | Active order changes |
| `types.ChannelAccountPortfolio` | Portfolio position updates |
| `types.ChannelAccountSummary` | Account summary updates |
| `types.ChannelAccountTradeHistory` | Trade history notifications |
| `types.ChannelAccountBots` | Bot status changes |
| `types.ChannelAccountConditional` | Conditional order updates |
| `types.ChannelSessionMMProtection` | MM protection status |

## API Coverage

| Category | REST | WebSocket |
|----------|------|-----------|
| Trading (insert, amend, cancel) | Yes | Yes |
| Conditional Orders | Yes | Yes |
| Bot Management | Yes | Yes |
| Request for Quote | Yes | Yes |
| Market Making (mass quote) | — | Yes |
| Market Data | Yes | Yes |
| Account/Portfolio | Yes | Yes |
| Trade/Order History | Yes | Yes |
| Historical OHLC Data | Yes | Yes |
| Wallet (deposits, withdrawals) | Yes | Yes |
| Notifications | Yes | Yes |
| Subscriptions | — | Yes |

## Documentation

For detailed documentation, see the [docs/](docs/) folder:

- [Getting Started](docs/getting-started.md)
- [Authentication](docs/authentication.md)
- [REST Client](docs/rest-client.md)
- [WebSocket Client](docs/ws-client.md)
- [Subscriptions](docs/subscriptions.md)
- [Trading](docs/trading.md)
- [Market Making](docs/market-making.md)
- [Bots](docs/bots.md)
- [Conditional Orders](docs/conditional-orders.md)
- [Error Handling](docs/error-handling.md)
- [Configuration](docs/configuration.md)
- [Enums](docs/enums.md)
- [Examples](docs/examples.md)

## Examples

See the [examples/](examples/) directory for complete, runnable programs:

| Example | Description |
|---------|-------------|
| [rest_basic](examples/rest_basic/) | REST API: fetch instruments, tickers, place and cancel orders |
| [ws_trading](examples/ws_trading/) | WebSocket trading: connect, login, insert/amend/cancel orders |
| [ws_subscriptions](examples/ws_subscriptions/) | Real-time subscriptions: tickers, order books, index prices |
| [market_making](examples/market_making/) | Market making: mass quotes and MM protection |
| [bot_management](examples/bot_management/) | Bot management: create SGSL, Grid, and DHedge bots |

## Development

```bash
# Install dependencies
go mod download

# Run all checks
make check

# Run tests with coverage
make test

# Run linter
make lint

# Format code
make fmt

# See all targets
make help
```

## Test Coverage

This project maintains high test coverage across all packages. Coverage is tracked via [Codecov](https://codecov.io/gh/amiwrpremium/go-thalex) and enforced in CI.

```bash
# Run tests with coverage report
make test

# Open HTML coverage report in browser
make coverage
```

| Package | Coverage |
|---------|----------|
| `apierr/` | 100% |
| `auth/` | 94%+ |
| `config/` | 100% |
| `enums/` | 100% |
| `internal/jsonrpc/` | 93%+ |
| `internal/transport/` | 93%+ |
| `rest/` | 95%+ |
| `types/` | 100% |
| `ws/` | 99%+ |

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) guide for details.

## Disclaimer

This project is provided **as-is**, without warranty of any kind. The author is **not affiliated with, endorsed by, or associated with [Thalex](https://thalex.com)** in any way. This is an independent, community-maintained SDK.

Use this software at your own risk. The author assumes no responsibility for any financial losses, trading errors, or other damages resulting from the use of this SDK. Always test thoroughly on the Thalex testnet before using in production.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
