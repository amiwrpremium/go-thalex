# Bot Management

Thalex provides server-side trading bots that run on the exchange infrastructure. The SDK supports creating, listing, and canceling bots via both REST and WebSocket clients.

## Bot Strategies

| Strategy | Constant | Description |
|----------|----------|-------------|
| SGSL | `enums.BotStrategySGSL` | Signal Go Stop-Loss -- enters at a signal price, exits at stop-loss |
| OCQ | `enums.BotStrategyOCQ` | Option Combo Quote -- quotes option combos around a signal |
| Levels | `enums.BotStrategyLevels` | Levels -- places orders at defined price levels |
| Grid | `enums.BotStrategyGrid` | Grid -- grid trading with symmetric levels |
| DHedge | `enums.BotStrategyDHedge` | Delta Hedger -- hedges delta exposure periodically |
| DFollow | `enums.BotStrategyDFollow` | Delta Follower -- follows delta of a target instrument |

## SGSL Bot (Signal Go Stop-Loss)

A directional bot that enters a position when a signal price is reached and exits at a stop-loss level.

### Creating an SGSL Bot

```go
import (
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
)

endTime := float64(time.Now().Add(24 * time.Hour).Unix())

params := types.NewSGSLBotParams(
    "BTC-PERPETUAL",   // instrument
    enums.TargetMark,  // signal source: "last", "mark", or "index"
    96000,             // entry price -- enter when signal crosses this
    0.1,               // target position (positive = long)
    92000,             // exit price (stop-loss)
    0,                 // exit position (usually 0)
    endTime,           // expiry timestamp (unix)
).
    WithMaxSlippage(200).   // Max slippage per trade in price units
    WithLabel("sgsl-long")  // Custom label

bot, err := client.CreateSGSLBot(ctx, params)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("SGSL bot: %s status=%s\n", bot.BotID, bot.Status)
```

### SGSLBotParams

```go
type SGSLBotParams struct {
    Strategy       enums.BotStrategy  // Always BotStrategySGSL
    InstrumentName string
    Signal         enums.Target       // "last", "mark", or "index"
    EntryPrice     float64
    TargetPosition float64
    ExitPrice      float64
    ExitPosition   float64
    EndTime        float64            // Unix timestamp
    MaxSlippage    *float64           // Optional
    Label          string             // Optional
}
```

## Grid Bot

A grid trading bot that places symmetric buy and sell orders at predefined price levels.

### Creating a Grid Bot

```go
endTime := float64(time.Now().Add(12 * time.Hour).Unix())

params := types.NewGridBotParams(
    "BTC-PERPETUAL",
    []float64{93000, 94000, 95000, 96000, 97000}, // grid levels
    0.01,      // step size per level
    endTime,
).
    WithBasePosition(0).              // Starting position
    WithTargetMeanPrice(95000).       // Target mean entry price
    WithUpsideExitPrice(98000).       // Exit if price goes above
    WithDownsideExitPrice(91000).     // Exit if price goes below
    WithMaxSlippage(100).
    WithLabel("grid-btc")

bot, err := client.CreateGridBot(ctx, params)
```

### GridBotParams

```go
type GridBotParams struct {
    Strategy          enums.BotStrategy
    InstrumentName    string
    Grid              []float64  // Price levels
    StepSize          float64    // Amount per level
    EndTime           float64
    BasePosition      *float64   // Optional
    TargetMeanPrice   *float64   // Optional
    UpsideExitPrice   *float64   // Optional
    DownsideExitPrice *float64   // Optional
    MaxSlippage       *float64   // Optional
    Label             string     // Optional
}
```

## Levels Bot

Similar to Grid, but with independent bid and ask levels.

### Creating a Levels Bot

```go
params := types.NewLevelsBotParams(
    "ETH-PERPETUAL",
    []float64{3400, 3420, 3440},       // bid levels
    []float64{3560, 3580, 3600},       // ask levels
    0.1,                                // step size
    endTime,
).
    WithBasePosition(0).
    WithLabel("levels-eth")

bot, err := client.CreateLevelsBot(ctx, params)
```

### LevelsBotParams

```go
type LevelsBotParams struct {
    Strategy          enums.BotStrategy
    InstrumentName    string
    Bids              []float64
    Asks              []float64
    StepSize          float64
    EndTime           float64
    BasePosition      *float64
    TargetMeanPrice   *float64
    UpsideExitPrice   *float64
    DownsideExitPrice *float64
    MaxSlippage       *float64
    Label             string
}
```

## DHedge Bot (Delta Hedger)

Periodically hedges delta exposure by trading the configured instrument.

### Creating a DHedge Bot

```go
params := types.NewDHedgeBotParams(
    "BTC-PERPETUAL",  // hedging instrument
    60,               // period in seconds
).
    WithTargetDelta(0).       // Target portfolio delta
    WithThreshold(0.5).       // Hedge when delta deviation exceeds this
    WithTolerance(0.1).       // Acceptable delta range
    WithMaxSlippage(200).
    WithEndTime(endTime).
    WithPosition("OBTCUSD").  // Position source (product to hedge)
    WithLabel("dhedge-btc")

bot, err := client.CreateDHedgeBot(ctx, params)
```

### DHedgeBotParams

```go
type DHedgeBotParams struct {
    Strategy       enums.BotStrategy
    InstrumentName string
    Period         float64     // Hedge check interval (seconds)
    Position       string      // Optional: position/product to hedge
    TargetDelta    *float64    // Optional: target delta (default: 0)
    Threshold      *float64    // Optional: trigger threshold
    Tolerance      *float64    // Optional: acceptable range
    MaxSlippage    *float64    // Optional
    EndTime        *float64    // Optional
    Label          string      // Optional
}
```

## DFollow Bot (Delta Follower)

Follows the delta of a target instrument by trading the configured instrument.

### Creating a DFollow Bot

```go
params := types.NewDFollowBotParams(
    "BTC-PERPETUAL",                  // trading instrument
    "BTC-28MAR25-100000-C",          // target instrument to follow
    1.0,                              // target amount
    30,                               // period (seconds)
    endTime,
).
    WithThreshold(0.1).
    WithTolerance(0.05).
    WithMaxSlippage(200).
    WithLabel("dfollow-btc")

bot, err := client.CreateDFollowBot(ctx, params)
```

### DFollowBotParams

```go
type DFollowBotParams struct {
    Strategy         enums.BotStrategy
    InstrumentName   string
    TargetInstrument string
    TargetAmount     float64
    Period           float64
    EndTime          float64
    Threshold        *float64
    Tolerance        *float64
    MaxSlippage      *float64
    Label            string
}
```

## OCQ Bot (Option Combo Quote)

Quotes option combinations around a signal price.

### Creating an OCQ Bot

```go
params := types.NewOCQBotParams(
    "BTC-28MAR25-100000-C",   // instrument
    enums.TargetMark,          // signal source
    0.005,                     // bid offset
    0.005,                     // ask offset
    0.1,                       // quote size
    -1.0,                      // min position
    1.0,                       // max position
    endTime,
).
    WithExitOffset(0.01).
    WithTargetPosition(0).
    WithLabel("ocq-btc")

bot, err := client.CreateOCQBot(ctx, params)
```

## Bot Response Type

All creation methods return `types.Bot`:

```go
type Bot struct {
    BotID           string
    Strategy        enums.BotStrategy
    Status          enums.BotStatus
    StopReason      enums.BotStopReason
    InstrumentName  string
    EndTime         *float64
    StartTime       float64
    StopTime        *float64
    Label           string
    RealisedPnl     float64
    Fee             float64
    AveragePrice    *float64
    PositionSize    *float64
    MarkPriceAtStop *float64
    // ... strategy-specific fields
}
```

**JSON example:**

```json
{
    "bot_id": "BOT-abc123",
    "strategy": "sgsl",
    "status": "active",
    "instrument_name": "BTC-PERPETUAL",
    "start_time": 1700000000.0,
    "end_time": 1700086400.0,
    "label": "sgsl-long",
    "realized_pnl": 0.0,
    "fee": 0.0
}
```

### Bot Status

| Value | Constant | Description |
|-------|----------|-------------|
| `"active"` | `BotStatusActive` | Bot is running |
| `"stopped"` | `BotStatusStopped` | Bot has stopped |

```go
if bot.Status.IsActive() {
    fmt.Println("Bot is running")
}
if bot.Status.IsFinal() {
    fmt.Printf("Bot stopped: %s\n", bot.StopReason)
}
```

### Bot Stop Reasons

| Value | Constant | Description |
|-------|----------|-------------|
| `"client_cancel"` | `BotStopReasonClientCancel` | Cancelled by user |
| `"client_bulk_cancel"` | `BotStopReasonClientBulkCancel` | Cancelled by bulk cancel |
| `"end_time"` | `BotStopReasonEndTime` | Reached end time |
| `"instrument_deactivated"` | `BotStopReasonInstrumentDeactivated` | Instrument delisted |
| `"margin_breach"` | `BotStopReasonMarginBreach` | Insufficient margin |
| `"admin_cancel"` | `BotStopReasonAdminCancel` | Cancelled by admin |
| `"conflict"` | `BotStopReasonConflict` | Conflicting bot/order |
| `"strategy"` | `BotStopReasonStrategy` | Strategy logic stopped it |

## Listing Bots

```go
// Active bots only.
bots, err := client.Bots(ctx, false)

// All bots including stopped ones.
bots, err := client.Bots(ctx, true)

for _, b := range bots {
    fmt.Printf("Bot %s: strategy=%s status=%s instrument=%s pnl=%.4f\n",
        b.BotID, b.Strategy, b.Status, b.InstrumentName, b.RealisedPnl)
}
```

## Canceling Bots

### Cancel a Specific Bot

```go
err := client.CancelBot(ctx, "BOT-abc123")
```

### Cancel All Bots

```go
n, err := client.CancelAllBots(ctx)
fmt.Printf("Cancelled %d bots\n", n)
```

## Subscribing to Bot Updates

Monitor bot status changes in real-time via WebSocket:

```go
wsClient.OnBots(func(bots []types.Bot) {
    for _, b := range bots {
        fmt.Printf("Bot update: %s status=%s\n", b.BotID, b.Status)
        if b.Status.IsFinal() {
            fmt.Printf("  Stop reason: %s  PnL: %.4f\n", b.StopReason, b.RealisedPnl)
        }
    }
})

err := wsClient.SubscribePrivate(ctx, types.ChannelAccountBots)
```

## Complete Example

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
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/rest"
    "github.com/amiwrpremium/go-thalex/types"
)

func main() {
    pemData, _ := os.ReadFile(os.Getenv("THALEX_PEM_PATH"))
    creds, _ := auth.NewCredentialsFromPEM(os.Getenv("THALEX_KEY_ID"), pemData)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    client := rest.NewClient(
        config.WithNetwork(config.Testnet),
        config.WithCredentials(creds),
    )

    endTime := float64(time.Now().Add(1 * time.Hour).Unix())

    // Create an SGSL bot.
    sgsl, err := client.CreateSGSLBot(ctx,
        types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetMark,
            96000, 0.1, 92000, 0, endTime).
            WithMaxSlippage(200).
            WithLabel("sgsl-example"),
    )
    if err != nil {
        log.Printf("SGSL: %v", err)
    } else {
        fmt.Printf("SGSL created: %s\n", sgsl.BotID)
    }

    // Create a Grid bot.
    grid, err := client.CreateGridBot(ctx,
        types.NewGridBotParams("BTC-PERPETUAL",
            []float64{93000, 94000, 95000, 96000, 97000},
            0.01, endTime).
            WithLabel("grid-example"),
    )
    if err != nil {
        log.Printf("Grid: %v", err)
    } else {
        fmt.Printf("Grid created: %s\n", grid.BotID)
    }

    // List all bots.
    bots, err := client.Bots(ctx, true)
    if err != nil {
        log.Fatal(err)
    }
    for _, b := range bots {
        fmt.Printf("  %s: %s status=%s\n", b.BotID, b.Strategy, b.Status)
    }

    // Cancel all bots.
    n, err := client.CancelAllBots(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Cancelled %d bots\n", n)
}
```

---

[< Market Making](market-making.md) | [Home](README.md) | [Conditional Orders >](conditional-orders.md)
