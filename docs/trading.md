# Trading

The SDK provides a typed, builder-pattern API for placing, amending, and canceling orders. All trading operations work identically on both REST and WebSocket clients.

## Insert Orders

### Builder Functions

Create order parameters using convenience constructors:

```go
import (
    "github.com/amiwrpremium/go-thalex/types"
    "github.com/amiwrpremium/go-thalex/enums"
)

// Buy order.
params := types.NewBuyOrderParams("BTC-PERPETUAL", 0.01)

// Sell order.
params := types.NewSellOrderParams("ETH-PERPETUAL", 0.1)

// Explicit direction.
params := types.NewInsertOrderParams(enums.DirectionBuy, "BTC-PERPETUAL", 0.01)
```

### Builder Methods

Chain `.With*()` methods to configure order parameters:

```go
params := types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
    WithPrice(95000).                                 // Limit price
    WithOrderType(enums.OrderTypeLimit).              // "limit" or "market"
    WithTimeInForce(enums.TimeInForceGoodTillCancelled). // GTC or IOC
    WithPostOnly(true).                               // Post-only mode
    WithRejectPostOnly(true).                         // Reject if would cross
    WithReduceOnly(true).                             // Reduce-only
    WithCollar(enums.CollarReject).                   // Collar handling
    WithLabel("my-bot-order").                        // User label
    WithClientOrderID(12345).                         // Client-supplied ID
    WithSTP(enums.STPLevelAccount, enums.STPActionCancelAggressor) // Self-trade prevention
```

### Builder Method Reference

| Method | Type | Description |
|--------|------|-------------|
| `WithPrice(v)` | `float64` | Limit price (omit for market order) |
| `WithOrderType(v)` | `enums.OrderType` | `"limit"` or `"market"` |
| `WithTimeInForce(v)` | `enums.TimeInForce` | `"good_till_cancelled"` or `"immediate_or_cancel"` |
| `WithPostOnly(v)` | `bool` | If true, order rejected if it would match immediately |
| `WithRejectPostOnly(v)` | `bool` | Book-or-cancel with PostOnly |
| `WithReduceOnly(v)` | `bool` | Only reduce existing position |
| `WithCollar(v)` | `enums.Collar` | `"ignore"`, `"reject"`, or `"clamp"` |
| `WithLabel(v)` | `string` | Custom label for tracking |
| `WithClientOrderID(v)` | `uint64` | Client-supplied order ID |
| `WithSTP(level, action)` | `STPLevel, STPAction` | Self-trade prevention config |

### InsertOrderParams Struct

```go
type InsertOrderParams struct {
    Direction      enums.Direction
    InstrumentName string
    Legs           []InsertLeg       // For combination orders
    Amount         float64
    Price          *float64
    OrderType      enums.OrderType
    TimeInForce    enums.TimeInForce
    PostOnly       *bool
    RejectPostOnly *bool
    ReduceOnly     *bool
    Collar         enums.Collar
    Label          string
    ClientOrderID  *uint64
    STPLevel       enums.STPLevel
    STPAction      enums.STPAction
}
```

### Placing an Order

```go
// Via REST.
order, err := restClient.Insert(ctx, params)

// Via WebSocket.
order, err := wsClient.Insert(ctx, params)
```

### Quick Market Orders

Shortcut methods for simple market orders:

```go
// Market buy.
order, err := client.Buy(ctx, "BTC-PERPETUAL", 0.01)

// Market sell.
order, err := client.Sell(ctx, "ETH-PERPETUAL", 0.1)
```

### Combination Orders

For multi-leg orders:

```go
params := types.NewComboInsertOrderParams(
    enums.DirectionBuy,
    []types.InsertLeg{
        {InstrumentName: "BTC-28MAR25-100000-C", Quantity: 1},
        {InstrumentName: "BTC-28MAR25-110000-C", Quantity: -1},
    },
    0.01,
).WithPrice(0.005)

order, err := client.Insert(ctx, params)
```

## Order Response: OrderStatus

Every order operation returns `types.OrderStatus`:

```go
type OrderStatus struct {
    OrderID         string
    OrderType       enums.OrderType
    TimeInForce     enums.TimeInForce
    InstrumentName  string
    Legs            []Leg
    Direction       enums.Direction
    Price           *float64
    Amount          float64
    FilledAmount    float64
    RemainingAmount float64
    Label           string
    ClientOrderID   *uint64
    Status          enums.OrderStatusValue
    Fills           []OrderFill
    ChangeReason    enums.ChangeReason
    DeleteReason    enums.DeleteReason
    InsertReason    enums.InsertReason
    CreateTime      float64
    CloseTime       *float64
    ReduceOnly      bool
    Persistent      bool
}
```

**JSON example:**

```json
{
    "order_id": "ORD-abc123",
    "order_type": "limit",
    "time_in_force": "good_till_cancelled",
    "instrument_name": "BTC-PERPETUAL",
    "direction": "buy",
    "price": 95000.0,
    "amount": 0.01,
    "filled_amount": 0.0,
    "remaining_amount": 0.01,
    "label": "my-order",
    "status": "open",
    "fills": [],
    "change_reason": "insert",
    "insert_reason": "client_request",
    "create_time": 1700000000.123,
    "reduce_only": false,
    "persistent": false
}
```

### Order Status Values

| Value | Constant | Description |
|-------|----------|-------------|
| `"open"` | `OrderStatusOpen` | Active, no fills yet |
| `"partially_filled"` | `OrderStatusPartiallyFilled` | Active, partially filled |
| `"filled"` | `OrderStatusFilled` | Fully filled (terminal) |
| `"cancelled"` | `OrderStatusCancelled` | Cancelled, no fills (terminal) |
| `"cancelled_partially_filled"` | `OrderStatusCancelledPartiallyFilled` | Cancelled after partial fill (terminal) |

Helper methods:

```go
if order.Status.IsActive() {
    // Order is still in the book.
}
if order.Status.IsFinal() {
    // Order is done.
}
```

## Amend Orders

Modify the price and amount of an existing order.

### By Order ID

```go
amended, err := client.Amend(ctx,
    types.NewAmendByOrderID("ORD-abc123", 95100, 0.02),
)
fmt.Printf("New price: %.2f  New amount: %.4f\n", *amended.Price, amended.Amount)
```

### By Client Order ID

```go
amended, err := client.Amend(ctx,
    types.NewAmendByClientOrderID(12345, 95100, 0.02),
)
```

### With Collar

```go
amended, err := client.Amend(ctx,
    types.NewAmendByOrderID("ORD-abc123", 95100, 0.02).
        WithCollar(enums.CollarClamp),
)
```

### AmendOrderParams

```go
type AmendOrderParams struct {
    OrderID       string
    ClientOrderID *uint64
    Price         float64
    Amount        float64
    Collar        enums.Collar
}
```

## Cancel Orders

### By Order ID

```go
cancelled, err := client.Cancel(ctx, types.CancelByOrderID("ORD-abc123"))
fmt.Printf("Cancelled: %s status=%s\n", cancelled.OrderID, cancelled.Status)
```

### By Client Order ID

```go
cancelled, err := client.Cancel(ctx, types.CancelByClientOrderID(12345))
```

### Cancel All

Cancel all open orders across all instruments:

```go
n, err := client.CancelAll(ctx)
fmt.Printf("Cancelled %d orders\n", n)
```

## Open Orders

```go
// All open orders.
orders, err := client.OpenOrders(ctx, "")

// Filtered by instrument.
orders, err := client.OpenOrders(ctx, "BTC-PERPETUAL")

for _, o := range orders {
    fmt.Printf("%s %s %s %.4f @ %v status=%s\n",
        o.OrderID, o.Direction, o.InstrumentName,
        o.Amount, o.Price, o.Status)
}
```

## Complete Trading Example

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

    // Place a limit buy order.
    order, err := client.Insert(ctx,
        types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
            WithPrice(90000).
            WithPostOnly(true).
            WithLabel("example").
            WithClientOrderID(42),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Placed: %s status=%s\n", order.OrderID, order.Status)

    // Amend price and amount.
    amended, err := client.Amend(ctx,
        types.NewAmendByOrderID(order.OrderID, 90100, 0.02),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Amended: price=%.2f amount=%.4f\n", *amended.Price, amended.Amount)

    // List open orders.
    open, err := client.OpenOrders(ctx, "BTC-PERPETUAL")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Open orders: %d\n", len(open))

    // Cancel the specific order.
    cancelled, err := client.Cancel(ctx, types.CancelByOrderID(order.OrderID))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Cancelled: %s reason=%s\n", cancelled.OrderID, cancelled.DeleteReason)
}
```

## Error Scenarios

### Common API Error Codes

| Scenario | Description |
|----------|-------------|
| Insufficient margin | Order rejected due to insufficient margin |
| Price outside collar | Order price outside allowed range |
| Instrument not found | Invalid instrument name |
| Order not found | Cancel/amend of non-existent order |
| Rate limit | Too many requests |
| Post-only rejected | Post-only order would cross the spread |

### Error Handling

```go
order, err := client.Insert(ctx, params)
if err != nil {
    if apiErr, ok := apierr.IsAPIError(err); ok {
        switch {
        case apiErr.Code == -32600:
            fmt.Println("Invalid request")
        default:
            fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
        }
        return
    }
    log.Fatal(err)
}
```

See [Error Handling](error-handling.md) for complete error handling patterns.

---

[< Subscriptions](subscriptions.md) | [Home](README.md) | [Market Making >](market-making.md)
