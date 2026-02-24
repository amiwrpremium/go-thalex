# WebSocket Client

The WebSocket client provides low-latency access to the Thalex API via JSON-RPC over WebSocket. It supports real-time subscriptions, automatic reconnection, and all the same endpoints as the REST client.

## Creating a Client

```go
import (
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/ws"
)

wsClient := ws.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
    config.WithWSReconnect(true),
    config.WithWSPingInterval(5 * time.Second),
)
```

See [Configuration](configuration.md) for all available options.

## Configuration Options

Options most relevant to the WebSocket client:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithNetwork(n)` | `config.Network` | `Production` | API environment |
| `WithCredentials(c)` | `*auth.Credentials` | `nil` | API credentials |
| `WithWSDialTimeout(d)` | `time.Duration` | `10s` | Connection timeout |
| `WithWSPingInterval(d)` | `time.Duration` | `5s` | Ping keepalive interval |
| `WithWSReconnect(b)` | `bool` | `false` | Enable auto-reconnect |
| `WithWSMaxReconnects(n)` | `int` | `10` | Max reconnect attempts |
| `WithWSReconnectWait(d)` | `time.Duration` | `1s` | Base wait between reconnects |
| `WithAccountNumber(a)` | `string` | `""` | Sub-account number |

## Connection Lifecycle

### Connect

```go
ctx := context.Background()

if err := wsClient.Connect(ctx); err != nil {
    log.Fatal(err)
}
```

`Connect` establishes the WebSocket connection and starts the ping keepalive loop. If reconnection is enabled, the reconnector also starts monitoring the connection.

### Login (Authentication)

After connecting, you must call `Login()` before using any private endpoints:

```go
if err := wsClient.Login(ctx); err != nil {
    log.Fatal(err)
}
```

Login generates a JWT token from the configured credentials and sends it to the server via the `public/login` JSON-RPC method. If an account number is configured, it is included in the login request.

### Check Connection State

```go
if wsClient.IsConnected() {
    fmt.Println("WebSocket is connected")
}
```

### Close

```go
if err := wsClient.Close(); err != nil {
    log.Printf("Close error: %v", err)
}
```

`Close()` stops the reconnector (if active), cancels all pending RPC calls, and gracefully closes the underlying WebSocket connection.

### Full Lifecycle Example

```go
wsClient := ws.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
    config.WithWSReconnect(true),
)

wsClient.OnErrorHandler(func(err error) {
    log.Printf("WS error: %v", err)
})

if err := wsClient.Connect(ctx); err != nil {
    log.Fatal(err)
}
defer wsClient.Close()

if err := wsClient.Login(ctx); err != nil {
    log.Fatal(err)
}

// Enable cancel-on-disconnect.
if err := wsClient.SetCancelOnDisconnect(ctx, true); err != nil {
    log.Fatal(err)
}

// Now use the client for trading, subscriptions, etc.
```

## Cancel-On-Disconnect

A critical safety feature for trading applications. When enabled, the server cancels all non-persistent orders if the WebSocket connection drops:

```go
// Enable.
err := wsClient.SetCancelOnDisconnect(ctx, true)

// Disable.
err := wsClient.SetCancelOnDisconnect(ctx, false)
```

This is especially important for market making to avoid leaving stale quotes after disconnection.

## Cancel Session

Cancel all non-persistent orders placed in the current WebSocket session:

```go
n, err := wsClient.CancelSession(ctx)
fmt.Printf("Cancelled %d session orders\n", n)
```

## Auto-Reconnect

When `WithWSReconnect(true)` is enabled, the client automatically handles disconnections:

1. Detects WebSocket disconnection
2. Waits `WSReconnectWait` duration (with exponential backoff)
3. Re-establishes the WebSocket connection
4. Re-authenticates (if credentials are configured)
5. Re-subscribes to all previously subscribed channels (both public and private)

```go
wsClient := ws.NewClient(
    config.WithNetwork(config.Testnet),
    config.WithCredentials(creds),
    config.WithWSReconnect(true),          // Enable reconnection
    config.WithWSMaxReconnects(20),        // Max 20 attempts
    config.WithWSReconnectWait(2*time.Second), // 2s base wait
)
```

### Reconnection Behavior

- **Max attempts:** Configurable via `WithWSMaxReconnects` (default: 10)
- **Wait strategy:** Exponential backoff from `WSReconnectWait` base
- **Re-authentication:** Automatic if credentials are configured
- **Re-subscription:** All registered handlers are automatically re-subscribed
- **Channel classification:** Channels prefixed with `account.`, `session.`, `user.`, or `mm.` are re-subscribed as private; all others as public

## Ping Keepalive

The client automatically sends WebSocket ping frames at a configurable interval to keep the connection alive:

```go
config.WithWSPingInterval(5 * time.Second) // default
```

If the server does not respond to pings, the client detects the disconnection and triggers reconnection (if enabled).

## Error Handler

Register a callback to receive connection-level errors:

```go
wsClient.OnErrorHandler(func(err error) {
    log.Printf("WebSocket error: %v", err)
})
```

This callback fires for:
- Read errors from the WebSocket
- JSON parsing errors on incoming messages
- Unexpected message types

It does **not** fire for API-level errors on individual RPC calls (those are returned from the method call itself).

## All Available Methods

The WebSocket client mirrors the REST client's endpoint coverage. All methods accept a `context.Context` as the first argument and block until the JSON-RPC response arrives.

### Public -- Market Data

```go
instruments, err := wsClient.Instruments(ctx)
all, err := wsClient.AllInstruments(ctx)
inst, err := wsClient.Instrument(ctx, "BTC-PERPETUAL")
ticker, err := wsClient.Ticker(ctx, "BTC-PERPETUAL")
idx, err := wsClient.Index(ctx, "BTCUSD")
book, err := wsClient.Book(ctx, "BTC-PERPETUAL")
info, err := wsClient.SystemInfo(ctx)
markData, err := wsClient.MarkPriceHistoricalData(ctx, "BTC-PERPETUAL", from, to, enums.Resolution1h)
indexData, err := wsClient.IndexPriceHistoricalData(ctx, "BTCUSD", from, to, enums.Resolution1h)
```

### Private -- Trading

```go
order, err := wsClient.Insert(ctx, params)
order, err := wsClient.Buy(ctx, "BTC-PERPETUAL", 0.01)
order, err := wsClient.Sell(ctx, "ETH-PERPETUAL", 0.1)
amended, err := wsClient.Amend(ctx, params)
cancelled, err := wsClient.Cancel(ctx, params)
n, err := wsClient.CancelAll(ctx)
orders, err := wsClient.OpenOrders(ctx, "BTC-PERPETUAL")
```

### Private -- Account

```go
positions, err := wsClient.Portfolio(ctx)
summary, err := wsClient.AccountSummary(ctx)
breakdown, err := wsClient.AccountBreakdown(ctx)
margin, err := wsClient.RequiredMarginBreakdown(ctx)
impact, err := wsClient.RequiredMarginForOrder(ctx, "BTC-PERPETUAL", 95000, 0.1)
```

### Private -- Market Making (WebSocket-Only)

These methods are only available via WebSocket, not REST:

```go
result, err := wsClient.MassQuote(ctx, params)
err := wsClient.CancelMassQuote(ctx)
err := wsClient.SetMMProtection(ctx, params)
```

See [Market Making](market-making.md) for detailed usage.

### Private -- Session

```go
err := wsClient.Login(ctx)
err := wsClient.SetCancelOnDisconnect(ctx, true)
n, err := wsClient.CancelSession(ctx)
```

### Private -- Bots

```go
bots, err := wsClient.Bots(ctx, true)
bot, err := wsClient.CreateSGSLBot(ctx, params)
bot, err := wsClient.CreateGridBot(ctx, params)
bot, err := wsClient.CreateDHedgeBot(ctx, params)
bot, err := wsClient.CreateDFollowBot(ctx, params)
bot, err := wsClient.CreateOCQBot(ctx, params)
bot, err := wsClient.CreateLevelsBot(ctx, params)
err := wsClient.CancelBot(ctx, "bot-id")
n, err := wsClient.CancelAllBots(ctx)
```

### Private -- Conditional Orders

```go
orders, err := wsClient.ConditionalOrders(ctx)
co, err := wsClient.CreateConditionalOrder(ctx, params)
err := wsClient.CancelConditionalOrder(ctx, "order-id")
n, err := wsClient.CancelAllConditionalOrders(ctx)
```

### Private -- RFQ

```go
rfq, err := wsClient.CreateRfq(ctx, params)
err := wsClient.CancelRfq(ctx, "rfq-id")
trades, err := wsClient.TradeRfq(ctx, params)
rfqs, err := wsClient.OpenRfqs(ctx)
history, err := wsClient.RfqHistory(ctx, from, to, offset, limit)
mmRfqs, err := wsClient.MMRfqs(ctx)
quote, err := wsClient.MMRfqInsertQuote(ctx, params)
quote, err := wsClient.MMRfqAmendQuote(ctx, params)
err := wsClient.MMRfqDeleteQuote(ctx, params)
quotes, err := wsClient.MMRfqQuotes(ctx)
```

### Private -- Wallet

```go
deposits, err := wsClient.CryptoDeposits(ctx)
withdrawals, err := wsClient.CryptoWithdrawals(ctx)
btcAddr, err := wsClient.BTCDepositAddress(ctx)
ethAddr, err := wsClient.ETHDepositAddress(ctx)
verify, err := wsClient.VerifyWithdrawal(ctx, params)
w, err := wsClient.Withdraw(ctx, params)
verify, err := wsClient.VerifyInternalTransfer(ctx, params)
err := wsClient.InternalTransfer(ctx, params)
```

### Private -- History

```go
trades, err := wsClient.TradeHistory(ctx, params)
orders, err := wsClient.OrderHistory(ctx, params)
marks, err := wsClient.DailyMarkHistory(ctx, params)
txns, err := wsClient.TransactionHistory(ctx, params)
```

### Private -- Notifications

```go
result, err := wsClient.NotificationsInbox(ctx, &limit)
err := wsClient.MarkNotificationAsRead(ctx, "notif-id", true)
```

### Subscriptions

```go
err := wsClient.Subscribe(ctx, channels...)
err := wsClient.SubscribePrivate(ctx, channels...)
err := wsClient.Unsubscribe(ctx, channels...)
err := wsClient.UnsubscribePrivate(ctx, channels...)
```

See [Subscriptions](subscriptions.md) for complete subscription documentation.

## REST vs WebSocket: When to Use Which

| Feature | REST | WebSocket |
|---------|------|-----------|
| Real-time subscriptions | No | Yes |
| Mass quotes | No | Yes |
| MM protection | No | Yes |
| Cancel-on-disconnect | No | Yes |
| Session management | No | Yes |
| Lower latency | No | Yes |
| Simpler setup | Yes | No |
| No persistent connection | Yes | No |
| Built-in HTTP retries | Yes | No (use reconnect) |

Use REST for simple scripts, one-off queries, and applications that don't need real-time data. Use WebSocket for trading bots, market making, and applications that need low latency or streaming data.

---

[< REST Client](rest-client.md) | [Home](README.md) | [Subscriptions >](subscriptions.md)
