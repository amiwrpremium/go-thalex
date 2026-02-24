# Conditional Orders

Conditional orders are server-side orders that activate (convert to a regular order) when a trigger condition is met. The SDK supports stop orders, stop-limit orders, bracket orders, and trailing stop orders.

## Order Types

| Type | Description | Constructor |
|------|-------------|-------------|
| Stop | Market order when trigger price is hit | `NewStopOrder` |
| Stop Limit | Limit order when trigger price is hit | `NewStopLimitOrder` |
| Bracket | Combined stop-loss and take-profit | `NewBracketOrder` |
| Trailing Stop | Stop that follows price at a fixed distance | `NewTrailingStopOrder` |

## Stop Orders

A stop order places a market order when the trigger price is reached.

```go
import (
    "github.com/amiwrpremium/go-thalex/enums"
    "github.com/amiwrpremium/go-thalex/types"
)

// Stop-loss for a long position: sell when price drops to 92000.
params := types.NewStopOrder(
    enums.DirectionSell,    // direction of the activated order
    "BTC-PERPETUAL",        // instrument
    0.1,                    // amount
    92000,                  // trigger (stop) price
).
    WithTarget(enums.TargetMark).  // Trigger on mark price
    WithReduceOnly(true).          // Only reduce existing position
    WithLabel("stop-loss")

co, err := client.CreateConditionalOrder(ctx, params)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Stop order: %s status=%s\n", co.OrderID, co.Status)
```

## Stop Limit Orders

A stop-limit order places a limit order (not market) when the trigger price is reached.

```go
// Stop-limit: when price hits 92000, place a limit sell at 91900.
params := types.NewStopLimitOrder(
    enums.DirectionSell,
    "BTC-PERPETUAL",
    0.1,
    92000,  // trigger price
    91900,  // limit price for the activated order
).
    WithTarget(enums.TargetLast).
    WithLabel("stop-limit")

co, err := client.CreateConditionalOrder(ctx, params)
```

### Checking Order Type

```go
if co.IsStopLimit() {
    fmt.Printf("This is a stop-limit with limit price: %.2f\n", *co.LimitPrice)
}
```

## Bracket Orders

A bracket order combines a stop-loss and a take-profit. When either price is hit, the order activates and the bracket is complete.

```go
// Bracket for a long position:
//   Stop-loss at 92000, take-profit at 99000.
params := types.NewBracketOrder(
    enums.DirectionSell,
    "BTC-PERPETUAL",
    0.1,
    92000,  // stop price (stop-loss)
    99000,  // bracket price (take-profit)
).
    WithTarget(enums.TargetMark).
    WithReduceOnly(true).
    WithLabel("bracket")

co, err := client.CreateConditionalOrder(ctx, params)
```

### Checking Bracket Status

```go
if co.IsBracket() {
    fmt.Printf("Bracket: stop=%.2f take-profit=%.2f\n",
        co.StopPrice, *co.BracketPrice)
}
```

## Trailing Stop Orders

A trailing stop follows the price at a fixed callback rate, activating when the price reverses by the specified distance.

```go
// Trailing stop: follows price up, triggers when price drops by 500.
params := types.NewTrailingStopOrder(
    enums.DirectionSell,
    "BTC-PERPETUAL",
    0.1,
    95000,  // initial stop price
    500,    // callback rate (trail distance)
).
    WithTarget(enums.TargetLast).
    WithReduceOnly(true).
    WithLabel("trail-stop")

co, err := client.CreateConditionalOrder(ctx, params)
```

### Checking Trailing Stop

```go
if co.IsTrailingStop() {
    fmt.Printf("Trailing stop: callback rate=%.2f\n", *co.TrailingStopCallbackRate)
}
```

## Builder Methods

All conditional order types share the same builder methods:

| Method | Type | Description |
|--------|------|-------------|
| `WithTarget(v)` | `enums.Target` | Trigger source: `"last"`, `"mark"`, or `"index"` |
| `WithLabel(v)` | `string` | Custom label |
| `WithReduceOnly(v)` | `bool` | Activated order is reduce-only |

### Target Values

| Value | Constant | Description |
|-------|----------|-------------|
| `"last"` | `enums.TargetLast` | Trigger on last trade price |
| `"mark"` | `enums.TargetMark` | Trigger on mark price |
| `"index"` | `enums.TargetIndex` | Trigger on index price |

## CreateConditionalOrderParams

```go
type CreateConditionalOrderParams struct {
    Direction                enums.Direction
    InstrumentName           string
    Amount                   float64
    StopPrice                float64
    LimitPrice               *float64  // Set for stop-limit orders
    BracketPrice             *float64  // Set for bracket orders
    TrailingStopCallbackRate *float64  // Set for trailing stop orders
    Target                   enums.Target
    Label                    string
    ReduceOnly               *bool
}
```

## Conditional Order Response

```go
type ConditionalOrder struct {
    OrderID                  string
    InstrumentName           string
    Direction                enums.Direction
    Amount                   float64
    Target                   enums.Target
    StopPrice                float64
    LimitPrice               *float64
    BracketPrice             *float64
    TrailingStopCallbackRate *float64
    Label                    string
    Status                   enums.ConditionalOrderStatus
    CreateTime               float64
    UpdateTime               float64
    ConvertTime              *float64
    ConvertedOrderID         string
    RejectReason             string
    ReduceOnly               bool
}
```

**JSON example:**

```json
{
    "order_id": "COND-abc123",
    "instrument_name": "BTC-PERPETUAL",
    "direction": "sell",
    "amount": 0.1,
    "target": "mark",
    "stop_price": 92000.0,
    "bracket_price": 99000.0,
    "status": "active",
    "create_time": 1700000000.123,
    "update_time": 1700000000.123,
    "reduce_only": true,
    "label": "bracket"
}
```

### Conditional Order Status Values

| Value | Constant | Description |
|-------|----------|-------------|
| `"created"` | `ConditionalOrderStatusCreated` | Just created |
| `"active"` | `ConditionalOrderStatusActive` | Active, waiting for trigger |
| `"converted"` | `ConditionalOrderStatusConverted` | Triggered, converted to order |
| `"rejected"` | `ConditionalOrderStatusRejected` | Rejected (see RejectReason) |
| `"cancel requested"` | `ConditionalOrderStatusCancelRequested` | Cancel in progress |
| `"cancelled"` | `ConditionalOrderStatusCancelled` | Cancelled |

```go
if co.Status.IsActive() {
    fmt.Println("Conditional order is active, waiting for trigger")
}
```

## Listing Conditional Orders

```go
orders, err := client.ConditionalOrders(ctx)
if err != nil {
    log.Fatal(err)
}

for _, co := range orders {
    fmt.Printf("%s: %s %s stop=%.2f status=%s\n",
        co.OrderID, co.Direction, co.InstrumentName,
        co.StopPrice, co.Status)

    if co.IsBracket() {
        fmt.Printf("  Bracket take-profit: %.2f\n", *co.BracketPrice)
    }
    if co.IsTrailingStop() {
        fmt.Printf("  Trail distance: %.2f\n", *co.TrailingStopCallbackRate)
    }
    if co.IsStopLimit() {
        fmt.Printf("  Limit price: %.2f\n", *co.LimitPrice)
    }
}
```

## Canceling Conditional Orders

### Cancel a Specific Order

```go
err := client.CancelConditionalOrder(ctx, "COND-abc123")
```

### Cancel All Conditional Orders

```go
n, err := client.CancelAllConditionalOrders(ctx)
fmt.Printf("Cancelled %d conditional orders\n", n)
```

## Subscribing to Updates

Monitor conditional order changes in real-time:

```go
wsClient.OnConditionalOrders(func(orders []types.ConditionalOrder) {
    for _, co := range orders {
        fmt.Printf("Conditional update: %s status=%s\n", co.OrderID, co.Status)
        if co.Status == enums.ConditionalOrderStatusConverted {
            fmt.Printf("  Converted to order: %s at %.0f\n",
                co.ConvertedOrderID, *co.ConvertTime)
        }
        if co.Status == enums.ConditionalOrderStatusRejected {
            fmt.Printf("  Rejected: %s\n", co.RejectReason)
        }
    }
})

err := wsClient.SubscribePrivate(ctx, types.ChannelAccountConditional)
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

    // Create a bracket order (stop-loss + take-profit).
    bracket, err := client.CreateConditionalOrder(ctx,
        types.NewBracketOrder(enums.DirectionSell, "BTC-PERPETUAL",
            0.1, 92000, 99000).
            WithTarget(enums.TargetMark).
            WithReduceOnly(true).
            WithLabel("bracket-example"),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bracket: %s stop=%.0f tp=%.0f\n",
        bracket.OrderID, bracket.StopPrice, *bracket.BracketPrice)

    // Create a trailing stop.
    trail, err := client.CreateConditionalOrder(ctx,
        types.NewTrailingStopOrder(enums.DirectionSell, "BTC-PERPETUAL",
            0.1, 94000, 500).
            WithTarget(enums.TargetLast).
            WithReduceOnly(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Trailing stop: %s callback=%.0f\n",
        trail.OrderID, *trail.TrailingStopCallbackRate)

    // List all conditional orders.
    orders, err := client.ConditionalOrders(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nActive conditional orders: %d\n", len(orders))

    // Cancel all.
    n, err := client.CancelAllConditionalOrders(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Cancelled %d conditional orders\n", n)
}
```

---

[< Bot Management](bots.md) | [Home](README.md) | [Request for Quote >](rfq.md)
