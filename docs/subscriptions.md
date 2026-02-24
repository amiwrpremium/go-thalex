# Real-time Subscriptions

The WebSocket client supports subscribing to real-time data channels. You register typed handlers for specific channels, then subscribe. The SDK automatically dispatches incoming notifications to the correct handler in a goroutine.

## Subscription Workflow

1. **Register a handler** using `On*` methods
2. **Subscribe** to the channel(s)
3. **Receive data** in your callback function
4. **Unsubscribe** when done

```go
// 1. Register handler.
ch := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
wsClient.OnTicker(ch, func(t types.Ticker) {
    fmt.Printf("BTC mark=%.2f\n", t.MarkPrice)
})

// 2. Subscribe.
err := wsClient.Subscribe(ctx, ch)

// ... receives data until ...

// 3. Unsubscribe.
err := wsClient.Unsubscribe(ctx, ch)
```

You can subscribe to multiple channels at once:

```go
err := wsClient.Subscribe(ctx, tickerCh, bookCh, indexCh, types.ChannelInstruments)
```

## Public Channel Helpers

The `types` package provides helper functions to construct channel names. These ensure correct formatting for parameterized channels.

### TickerChannel

Full ticker data for an instrument at a specified delay.

```go
ch := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
// => "ticker.BTC-PERPETUAL.100ms"

ch := types.TickerChannel("ETH-PERPETUAL", enums.Delay1000ms)
// => "ticker.ETH-PERPETUAL.1000ms"

ch := types.TickerChannel("BTC-PERPETUAL", enums.DelayNone)
// => "ticker.BTC-PERPETUAL.raw"
```

**Handler signature:** `func(types.Ticker)`

```go
wsClient.OnTicker(ch, func(t types.Ticker) {
    fmt.Printf("Mark: %.2f Bid: %v Ask: %v\n",
        t.MarkPrice, t.BestBidPrice, t.BestAskPrice)
})
```

**Delay options:**

| Constant | Value | Description |
|----------|-------|-------------|
| `enums.DelayNone` | `"raw"` | No throttling, every update |
| `enums.Delay100ms` | `"100ms"` | At most every 100ms |
| `enums.Delay1000ms` | `"1000ms"` | At most every 1 second |

### BookChannel

Order book snapshots with configurable grouping and depth.

```go
ch := types.BookChannel("BTC-PERPETUAL", 1, 10, enums.Delay100ms)
// => "book.BTC-PERPETUAL.1.10.100ms"

ch := types.BookChannel("ETH-PERPETUAL", 5, 20, enums.Delay1000ms)
// => "book.ETH-PERPETUAL.5.20.1000ms"
```

**Parameters:**
- `instrument` -- Instrument name
- `grouping` -- Price grouping level (1 = tick size, 5 = 5x tick size, etc.)
- `nlevels` -- Number of price levels per side
- `delay` -- Throttle delay

**Handler signature:** `func(types.BookUpdate)`

```go
wsClient.OnBook(ch, func(b types.BookUpdate) {
    if len(b.Bids) > 0 {
        fmt.Printf("Best bid: %.2f x %.4f\n", b.Bids[0].Price(), b.Bids[0].Amount())
    }
    if len(b.Asks) > 0 {
        fmt.Printf("Best ask: %.2f x %.4f\n", b.Asks[0].Price(), b.Asks[0].Amount())
    }
    for _, trade := range b.Trades {
        fmt.Printf("Trade: %s %.4f @ %.2f\n", trade.Direction, trade.Amount, trade.Price)
    }
})
```

**BookUpdate JSON example:**

```json
{
    "bids": [[95100.5, 1.25, 1.25], [95100.0, 0.50, 0.50]],
    "asks": [[95101.0, 0.80, 0.80], [95101.5, 1.10, 1.10]],
    "last": 95100.5,
    "time": 1700000000.123,
    "trades": [
        {"d": "buy", "p": 95100.5, "a": 0.05, "t": 1700000000.100}
    ]
}
```

### LWTChannel

Lightweight ticker -- a condensed version with just the essentials.

```go
ch := types.LWTChannel("BTC-PERPETUAL", enums.Delay100ms)
// => "lwt.BTC-PERPETUAL.100ms"
```

**Handler signature:** `func(types.LightweightTicker)`

```go
wsClient.OnLWT(ch, func(t types.LightweightTicker) {
    fmt.Printf("Mark: %.2f\n", t.MarkPrice)
    if t.BestBidPrice != nil && t.BestAskPrice != nil {
        fmt.Printf("Bid: %.2f Ask: %.2f\n", *t.BestBidPrice, *t.BestAskPrice)
    }
})
```

### RecentTradesChannel

Trade stream for an underlying or specific category.

```go
ch := types.RecentTradesChannel("BTCUSD", enums.RecentTradesCategoryAll)
// => "recent_trades.BTCUSD.all"

ch := types.RecentTradesChannel("ETHUSD", enums.RecentTradesCategoryBlock)
// => "recent_trades.ETHUSD.block"
```

**Categories:**

| Constant | Value | Description |
|----------|-------|-------------|
| `enums.RecentTradesCategoryAll` | `"all"` | All trade types |
| `enums.RecentTradesCategoryNormal` | `"normal"` | Normal trades only |
| `enums.RecentTradesCategoryBlock` | `"block"` | Block trades only |
| `enums.RecentTradesCategoryCombo` | `"combo"` | Combination trades only |

**Handler signature:** `func([]types.RecentTrade)`

```go
wsClient.OnRecentTrades(ch, func(trades []types.RecentTrade) {
    for _, t := range trades {
        fmt.Printf("%s %s %.4f @ %.2f (%s)\n",
            t.InstrumentName, t.Direction, t.Amount, t.Price, t.TradeType)
    }
})
```

### PriceIndexChannel

Index price updates for an underlying.

```go
ch := types.PriceIndexChannel("BTCUSD")
// => "price_index.BTCUSD"
```

**Handler signature:** `func(types.IndexPrice)`

```go
wsClient.OnPriceIndex(ch, func(idx types.IndexPrice) {
    fmt.Printf("%s = %.2f\n", idx.IndexName, idx.Price)
})
```

### BasePriceChannel

Forward (base) price for a specific underlying and expiration.

```go
ch := types.BasePriceChannel("BTCUSD", "2025-03-28")
// => "base_price.BTCUSD.2025-03-28"
```

### IndexComponentsChannel

Breakdown of an index price into its exchange components.

```go
ch := types.IndexComponentsChannel("BTCUSD")
// => "index_components.BTCUSD"
```

### UnderlyingStatisticsChannel

Open interest and statistics for an underlying.

```go
ch := types.UnderlyingStatisticsChannel("BTCUSD")
// => "underlying_statistics.BTCUSD"
```

### Instruments (Constant)

Instrument activation/deactivation notifications.

```go
// No helper needed -- use the constant directly.
wsClient.OnInstruments(func(instruments []types.Instrument) {
    for _, inst := range instruments {
        fmt.Printf("Instrument: %s type=%s\n", inst.InstrumentName, inst.Type)
    }
})

err := wsClient.Subscribe(ctx, types.ChannelInstruments)
```

## Public Channel Constants

Channels that require no parameters:

| Constant | Value | Description |
|----------|-------|-------------|
| `types.ChannelInstruments` | `"instruments"` | Instrument changes |
| `types.ChannelRfqs` | `"rfqs"` | Public RFQ notifications |
| `types.ChannelSystem` | `"system"` | System events |
| `types.ChannelBanners` | `"banners"` | Banner updates |

## Private Channels

Private channels require authentication (`Login()`) and use `SubscribePrivate()`.

### Channel Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `types.ChannelAccountOrders` | `"account.orders"` | All order updates |
| `types.ChannelAccountPersistent` | `"account.persistent_orders"` | Persistent order updates |
| `types.ChannelSessionOrders` | `"session.orders"` | Session order updates |
| `types.ChannelAccountPortfolio` | `"account.portfolio"` | Position changes |
| `types.ChannelAccountSummary` | `"account.summary"` | Account summary changes |
| `types.ChannelAccountTradeHistory` | `"account.trade_history"` | New trades |
| `types.ChannelAccountOrderHistory` | `"account.order_history"` | Order history updates |
| `types.ChannelAccountConditional` | `"account.conditional_orders"` | Conditional order changes |
| `types.ChannelAccountBots` | `"account.bots"` | Bot status changes |
| `types.ChannelAccountRfqs` | `"account.rfqs"` | Account RFQ updates |
| `types.ChannelAccountRfqHistory` | `"account.rfq_history"` | RFQ history updates |
| `types.ChannelSessionMMProtection` | `"session.mm_protection"` | MM protection events |
| `types.ChannelUserNotifications` | `"user.inbox_notifications"` | Inbox notifications |
| `types.ChannelMMRfqs` | `"mm.rfqs"` | MM RFQ opportunities |
| `types.ChannelMMRfqQuotes` | `"mm.rfq_quotes"` | MM RFQ quote updates |

### Private Subscription Example

```go
// Authenticate first.
if err := wsClient.Login(ctx); err != nil {
    log.Fatal(err)
}

// Register handlers.
wsClient.OnOrders(func(orders []types.OrderStatus) {
    for _, o := range orders {
        fmt.Printf("Order %s: %s %s %.4f @ %v\n",
            o.OrderID, o.Direction, o.Status, o.Amount, o.Price)
    }
})

wsClient.OnPortfolio(func(positions []types.PortfolioEntry) {
    for _, p := range positions {
        fmt.Printf("Position %s: %.4f pnl=%.4f\n",
            p.InstrumentName, p.Position, p.UnrealisedPnl)
    }
})

wsClient.OnAccountSummary(func(s types.AccountSummary) {
    fmt.Printf("Margin: %.2f Required: %.2f Remaining: %.2f\n",
        s.Margin, s.RequiredMargin, s.RemainingMargin)
})

wsClient.OnTradeHistory(func(trades []types.Trade) {
    for _, t := range trades {
        fmt.Printf("Trade %s: %s %.4f @ %.2f\n",
            t.TradeID, t.Direction, t.Amount, t.Price)
    }
})

wsClient.OnBots(func(bots []types.Bot) {
    for _, b := range bots {
        fmt.Printf("Bot %s: %s status=%s\n", b.BotID, b.Strategy, b.Status)
    }
})

wsClient.OnConditionalOrders(func(orders []types.ConditionalOrder) {
    for _, o := range orders {
        fmt.Printf("Conditional %s: %s status=%s\n", o.OrderID, o.Direction, o.Status)
    }
})

wsClient.OnMMProtection(func(update types.MMProtectionUpdate) {
    fmt.Printf("MM Protection: product=%s reason=%s\n", update.Product, update.Reason)
})

// Subscribe to all private channels.
err := wsClient.SubscribePrivate(ctx,
    types.ChannelAccountOrders,
    types.ChannelAccountPortfolio,
    types.ChannelAccountSummary,
    types.ChannelAccountTradeHistory,
    types.ChannelAccountBots,
    types.ChannelAccountConditional,
    types.ChannelSessionMMProtection,
)
```

## All Typed Handler Methods

### Public Channel Handlers

| Method | Channel Type | Callback Signature |
|--------|-------------|-------------------|
| `OnTicker(ch, fn)` | Ticker | `func(types.Ticker)` |
| `OnBook(ch, fn)` | Book | `func(types.BookUpdate)` |
| `OnLWT(ch, fn)` | Lightweight ticker | `func(types.LightweightTicker)` |
| `OnRecentTrades(ch, fn)` | Recent trades | `func([]types.RecentTrade)` |
| `OnPriceIndex(ch, fn)` | Price index | `func(types.IndexPrice)` |
| `OnInstruments(fn)` | Instruments | `func([]types.Instrument)` |
| `OnSystemEvent(fn)` | System events | `func(types.SystemEvent)` |
| `OnBanners(fn)` | Banners | `func([]types.Banner)` |

### Private Channel Handlers

| Method | Channel | Callback Signature |
|--------|---------|-------------------|
| `OnOrders(fn)` | `account.orders` | `func([]types.OrderStatus)` |
| `OnPersistentOrders(fn)` | `account.persistent_orders` | `func([]types.OrderStatus)` |
| `OnSessionOrders(fn)` | `session.orders` | `func([]types.OrderStatus)` |
| `OnPortfolio(fn)` | `account.portfolio` | `func([]types.PortfolioEntry)` |
| `OnAccountSummary(fn)` | `account.summary` | `func(types.AccountSummary)` |
| `OnTradeHistory(fn)` | `account.trade_history` | `func([]types.Trade)` |
| `OnOrderHistory(fn)` | `account.order_history` | `func([]types.OrderHistory)` |
| `OnConditionalOrders(fn)` | `account.conditional_orders` | `func([]types.ConditionalOrder)` |
| `OnBots(fn)` | `account.bots` | `func([]types.Bot)` |
| `OnRfqs(fn)` | `account.rfqs` | `func([]types.Rfq)` |
| `OnMMRfqs(fn)` | `mm.rfqs` | `func([]types.Rfq)` |
| `OnMMRfqQuotes(fn)` | `mm.rfq_quotes` | `func([]types.RfqOrder)` |
| `OnMMProtection(fn)` | `session.mm_protection` | `func(types.MMProtectionUpdate)` |
| `OnNotifications(fn)` | `user.inbox_notifications` | `func(types.Notification)` |

### Raw Handler

For channels without a typed handler, or to receive raw JSON:

```go
wsClient.OnRaw("some.custom.channel", func(data json.RawMessage) {
    fmt.Printf("Raw data: %s\n", string(data))
})
```

## Unsubscribing

```go
// Unsubscribe from public channels.
err := wsClient.Unsubscribe(ctx, tickerCh, bookCh)

// Unsubscribe from private channels.
err := wsClient.UnsubscribePrivate(ctx,
    types.ChannelAccountOrders,
    types.ChannelAccountPortfolio,
)
```

Unsubscribing removes the handler from the internal map and sends the unsubscribe request to the server.

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "time"

    "github.com/amiwrpremium/go-thalex/auth"
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
    "github.com/amiwrpremium/go-thalex/ws"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    pemData, _ := os.ReadFile(os.Getenv("THALEX_PEM_PATH"))
    creds, _ := auth.NewCredentialsFromPEM(os.Getenv("THALEX_KEY_ID"), pemData)

    wsClient := ws.NewClient(
        config.WithNetwork(config.Testnet),
        config.WithCredentials(creds),
        config.WithWSReconnect(true),
    )

    wsClient.OnErrorHandler(func(err error) {
        log.Printf("Error: %v", err)
    })

    if err := wsClient.Connect(ctx); err != nil {
        log.Fatal(err)
    }
    defer wsClient.Close()

    // Public subscriptions.
    tickerCh := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
    bookCh := types.BookChannel("BTC-PERPETUAL", 1, 5, enums.Delay100ms)

    wsClient.OnTicker(tickerCh, func(t types.Ticker) {
        fmt.Printf("[ticker] mark=%.2f\n", t.MarkPrice)
    })
    wsClient.OnBook(bookCh, func(b types.BookUpdate) {
        fmt.Printf("[book] %d bids, %d asks\n", len(b.Bids), len(b.Asks))
    })

    wsClient.Subscribe(ctx, tickerCh, bookCh)

    // Private subscriptions.
    loginCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    if err := wsClient.Login(loginCtx); err != nil {
        log.Fatal(err)
    }

    wsClient.OnOrders(func(orders []types.OrderStatus) {
        for _, o := range orders {
            fmt.Printf("[order] %s %s %s\n", o.OrderID, o.Direction, o.Status)
        }
    })
    wsClient.OnPortfolio(func(positions []types.PortfolioEntry) {
        for _, p := range positions {
            fmt.Printf("[position] %s = %.4f\n", p.InstrumentName, p.Position)
        }
    })

    wsClient.SubscribePrivate(ctx,
        types.ChannelAccountOrders,
        types.ChannelAccountPortfolio,
    )

    <-ctx.Done()
}
```

---

[< WebSocket Client](ws-client.md) | [Home](README.md) | [Trading >](trading.md)
