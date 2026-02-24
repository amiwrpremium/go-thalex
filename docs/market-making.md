# Market Making

Market making operations are **WebSocket-only** -- they are not available via the REST client. This includes mass quotes, mass quote cancellation, and market maker protection.

## Mass Quotes

Mass quotes allow you to insert or update bid and ask quotes for multiple instruments in a single request. They replace all existing mass quotes for the quoted instruments.

### Single-Level Quotes

The simplest form -- one bid and one ask per instrument:

```go
import "github.com/amiwrpremium/go-thalex/types"

result, err := wsClient.MassQuote(ctx,
    types.NewMassQuoteParams([]types.DoubleSidedQuote{
        types.NewSingleLevelQuote("BTC-PERPETUAL",
            95000, 0.1,  // bid price, bid amount
            95100, 0.1,  // ask price, ask amount
        ),
        types.NewSingleLevelQuote("ETH-PERPETUAL",
            3480, 1.0,  // bid price, bid amount
            3520, 1.0,  // ask price, ask amount
        ),
    }),
)
fmt.Printf("Success: %d  Fail: %d\n", result.NSuccess, result.NFail)
```

### Multi-Level Quotes

Provide multiple price levels on each side:

```go
result, err := wsClient.MassQuote(ctx,
    types.NewMassQuoteParams([]types.DoubleSidedQuote{
        types.NewDoubleSidedQuote("BTC-PERPETUAL",
            []types.QuoteLevel{
                {Price: 95000, Amount: 0.05},  // Best bid
                {Price: 94900, Amount: 0.10},  // Second level
                {Price: 94800, Amount: 0.20},  // Third level
            },
            []types.QuoteLevel{
                {Price: 95100, Amount: 0.05},  // Best ask
                {Price: 95200, Amount: 0.10},  // Second level
                {Price: 95300, Amount: 0.20},  // Third level
            },
        ),
    }),
)
```

### One-Sided Quotes

Quote only one side by passing `nil` for the other:

```go
// Bid only.
q := types.NewDoubleSidedQuote("BTC-PERPETUAL",
    []types.QuoteLevel{{Price: 95000, Amount: 0.1}},
    nil,
)

// Ask only.
q := types.NewDoubleSidedQuote("BTC-PERPETUAL",
    nil,
    []types.QuoteLevel{{Price: 95100, Amount: 0.1}},
)
```

### MassQuoteParams Builder

```go
params := types.NewMassQuoteParams(quotes).
    WithLabel("mm-bot").              // Label all quotes
    WithPostOnly(true).               // Post-only mode
    WithRejectPostOnly(true).         // Reject if would cross
    WithSTP(                          // Self-trade prevention
        enums.STPLevelAccount,
        enums.STPActionCancelAggressor,
    )
```

| Method | Type | Description |
|--------|------|-------------|
| `WithLabel(v)` | `string` | Label for all quotes |
| `WithPostOnly(v)` | `bool` | Post-only mode |
| `WithRejectPostOnly(v)` | `bool` | Reject if would cross the spread |
| `WithSTP(level, action)` | `STPLevel, STPAction` | Self-trade prevention |

### MassQuoteParams Struct

```go
type MassQuoteParams struct {
    Quotes         []DoubleSidedQuote
    Label          string
    PostOnly       *bool
    RejectPostOnly *bool
    STPLevel       enums.STPLevel
    STPAction      enums.STPAction
}
```

### Mass Quote Response

```go
type DoubleSidedQuoteResult struct {
    NSuccess int          // Number of successful quote sides
    NFail    int          // Number of failed quote sides
    Errors   []QuoteError // Details of failures
}

type QuoteError struct {
    Code    int
    Message string
    Side    string   // "bid" or "ask"
    Price   *float64
}
```

**JSON example:**

```json
{
    "n_success": 3,
    "n_fail": 1,
    "errors": [
        {
            "code": -32000,
            "message": "price outside collar",
            "side": "bid",
            "price": 95000.0
        }
    ]
}
```

### Handling Mass Quote Errors

```go
result, err := wsClient.MassQuote(ctx, params)
if err != nil {
    // RPC-level error (connection, auth, etc.)
    log.Fatal(err)
}

// Check per-quote errors.
if result.NFail > 0 {
    for _, e := range result.Errors {
        log.Printf("Quote error: side=%s code=%d msg=%s", e.Side, e.Code, e.Message)
    }
}
```

## Cancel Mass Quotes

Cancel all outstanding mass quotes:

```go
err := wsClient.CancelMassQuote(ctx)
if err != nil {
    log.Fatal(err)
}
```

This cancels all quotes placed via `MassQuote` but does not affect regular orders placed via `Insert`.

## Market Maker Protection

MM Protection is a circuit breaker that automatically cancels all your quotes when your trading volume exceeds a configured threshold within a rolling window. This protects against adverse selection during fast markets.

### Set MM Protection

```go
err := wsClient.SetMMProtection(ctx, &types.MMProtectionParams{
    Product:     "FBTCUSD",   // Product group
    TradeAmount: 5.0,          // Max fill amount in rolling window
    QuoteAmount: 25.0,         // Max quoted amount in rolling window
})
```

### MMProtectionParams

```go
type MMProtectionParams struct {
    Product     enums.Product  // Product group (e.g., "FBTCUSD", "OBTCUSD")
    TradeAmount float64        // Threshold for filled amount
    QuoteAmount float64        // Threshold for outstanding quote amount
}
```

### Product Groups

Product groups typically follow the pattern:
- `FBTCUSD` -- BTC futures and perpetuals
- `FETHUSD` -- ETH futures and perpetuals
- `OBTCUSD` -- BTC options
- `OETHUSD` -- ETH options

### Monitoring MM Protection

Subscribe to protection events to know when your quotes are being cancelled:

```go
wsClient.OnMMProtection(func(update types.MMProtectionUpdate) {
    switch update.Reason {
    case enums.MMProtectionReasonTriggered:
        log.Printf("MM PROTECTION TRIGGERED for %s at %f", update.Product, update.Time)
        // All quotes for this product have been cancelled.
        // Wait for the reset before re-quoting.
    case enums.MMProtectionReasonReset:
        log.Printf("MM protection reset for %s", update.Product)
        // Safe to resume quoting.
    }
})

err := wsClient.SubscribePrivate(ctx, types.ChannelSessionMMProtection)
```

**MMProtectionUpdate JSON:**

```json
{
    "product": "FBTCUSD",
    "reason": "triggered",
    "time": 1700000000.123
}
```

## Complete Market Making Example

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

    if err := wsClient.Login(ctx); err != nil {
        log.Fatal(err)
    }

    // Safety: cancel orders on disconnect.
    wsClient.SetCancelOnDisconnect(ctx, true)

    // Configure MM protection.
    wsClient.SetMMProtection(ctx, &types.MMProtectionParams{
        Product:     "FBTCUSD",
        TradeAmount: 5.0,
        QuoteAmount: 25.0,
    })

    // Monitor protection events.
    wsClient.OnMMProtection(func(update types.MMProtectionUpdate) {
        log.Printf("MM Protection: %s %s", update.Product, update.Reason)
    })
    wsClient.SubscribePrivate(ctx, types.ChannelSessionMMProtection)

    // Stream ticker data to derive quotes.
    tickerCh := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
    wsClient.OnTicker(tickerCh, func(t types.Ticker) {
        if t.BestBidPrice == nil || t.BestAskPrice == nil {
            return
        }
        mid := (*t.BestBidPrice + *t.BestAskPrice) / 2
        spread := 50.0

        result, err := wsClient.MassQuote(ctx,
            types.NewMassQuoteParams([]types.DoubleSidedQuote{
                types.NewSingleLevelQuote("BTC-PERPETUAL",
                    mid-spread/2, 0.05,
                    mid+spread/2, 0.05,
                ),
            }).WithPostOnly(true).WithLabel("mm-bot"),
        )
        if err != nil {
            log.Printf("MassQuote error: %v", err)
            return
        }
        if result.NFail > 0 {
            for _, e := range result.Errors {
                log.Printf("Quote fail: %s %s", e.Side, e.Message)
            }
        }
    })
    wsClient.Subscribe(ctx, tickerCh)

    fmt.Println("Market making started. Press Ctrl+C to exit.")
    <-ctx.Done()

    // Clean up.
    cleanCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    wsClient.CancelMassQuote(cleanCtx)
}
```

---

[< Trading](trading.md) | [Home](README.md) | [Bot Management >](bots.md)
