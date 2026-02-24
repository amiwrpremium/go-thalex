# Request for Quote (RFQ)

The RFQ system allows traders to request quotes for custom multi-leg orders from market makers. Market makers can respond with quotes, and the requestor can trade on those quotes.

## RFQ Workflow

1. **Requestor** creates an RFQ with specific legs and amount
2. **Market makers** see the RFQ and insert quotes (bid/ask prices)
3. **Requestor** executes a trade on the best quote
4. The RFQ is filled or cancelled

## Creating an RFQ

```go
import "github.com/amiwrpremium/go-thalex/types"

rfq, err := client.CreateRfq(ctx, &types.CreateRfqParams{
    Legs: []types.InsertLeg{
        {InstrumentName: "BTC-28MAR25-100000-C", Quantity: 1},
        {InstrumentName: "BTC-28MAR25-110000-C", Quantity: -1},
    },
    Amount: 0.5,
    Label:  "call-spread-rfq",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("RFQ created: %s\n", rfq.RfqID)
```

### CreateRfqParams

```go
type CreateRfqParams struct {
    Legs   []InsertLeg  // Leg definitions (instrument + quantity)
    Amount float64      // Notional amount
    Label  string       // Optional label
}
```

## Canceling an RFQ

```go
err := client.CancelRfq(ctx, "RFQ-abc123")
```

## Trading on an RFQ

Execute a trade on a quoted RFQ:

```go
trades, err := client.TradeRfq(ctx, &types.TradeRfqParams{
    RfqID:     "RFQ-abc123",
    Direction: enums.DirectionBuy,
    Price:     0.015,
    Amount:    0.5,
})
if err != nil {
    log.Fatal(err)
}
for _, t := range trades {
    fmt.Printf("Trade: %s %s %.4f @ %.4f\n",
        t.InstrumentName, t.Direction, t.Amount, t.Price)
}
```

### TradeRfqParams

```go
type TradeRfqParams struct {
    RfqID     string
    Direction enums.Direction
    Price     float64
    Amount    float64
}
```

## Listing Open RFQs

```go
rfqs, err := client.OpenRfqs(ctx)
for _, r := range rfqs {
    fmt.Printf("RFQ %s: amount=%.4f legs=%d\n", r.RfqID, r.Amount, len(r.Legs))
    if r.QuotedBid != nil {
        fmt.Printf("  Best bid: %.4f x %.4f\n", r.QuotedBid.Price, r.QuotedBid.Amount)
    }
    if r.QuotedAsk != nil {
        fmt.Printf("  Best ask: %.4f x %.4f\n", r.QuotedAsk.Price, r.QuotedAsk.Amount)
    }
}
```

## RFQ History

```go
from := float64(time.Now().Add(-24 * time.Hour).Unix())
to := float64(time.Now().Unix())
limit := 50

rfqs, err := client.RfqHistory(ctx, &from, &to, nil, &limit)
for _, r := range rfqs {
    fmt.Printf("RFQ %s: event=%s\n", r.RfqID, r.Event)
}
```

## RFQ Response Type

```go
type Rfq struct {
    RfqID          string
    Legs           []RfqLeg
    Amount         float64
    CreateTime     float64
    ValidUntil     *float64
    Label          string
    InsertReason   enums.RfqInsertReason
    DeleteReason   string
    VolumeTickSize *float64
    QuotedBid      *RfqQuotedSide
    QuotedAsk      *RfqQuotedSide
    TradePrice     *float64
    TradeAmount    *float64
    CloseTime      *float64
    Event          enums.RfqEvent
}
```

**JSON example:**

```json
{
    "rfq_id": "RFQ-abc123",
    "legs": [
        {"instrument_name": "BTC-28MAR25-100000-C", "quantity": 1, "fee_quantity": 0.001},
        {"instrument_name": "BTC-28MAR25-110000-C", "quantity": -1, "fee_quantity": 0.001}
    ],
    "amount": 0.5,
    "create_time": 1700000000.123,
    "valid_until": 1700000060.123,
    "event": "Created",
    "insert_reason": "client_request"
}
```

### RFQ Events

| Value | Constant | Description |
|-------|----------|-------------|
| `"Created"` | `RfqEventCreated` | Newly created |
| `"Cancelled"` | `RfqEventCancelled` | Cancelled by requestor |
| `"Traded"` | `RfqEventTraded` | Executed |
| `"Existing"` | `RfqEventExisting` | Existing on subscription start |

## Market Maker RFQ Operations

Market makers use a separate set of endpoints to provide liquidity to RFQs.

### View Available RFQs

```go
rfqs, err := client.MMRfqs(ctx)
for _, r := range rfqs {
    fmt.Printf("MM RFQ %s: %d legs, amount=%.4f\n",
        r.RfqID, len(r.Legs), r.Amount)
}
```

### Insert a Quote

```go
quote, err := client.MMRfqInsertQuote(ctx, &types.RfqQuoteInsertParams{
    RfqID:     "RFQ-abc123",
    Direction: enums.DirectionBuy,
    Amount:    0.5,
    Price:     0.012,
    Label:     "mm-quote",
})
fmt.Printf("Quote inserted: %s\n", quote.OrderID)
```

### Amend a Quote

```go
amended, err := client.MMRfqAmendQuote(ctx, &types.RfqQuoteAmendParams{
    OrderID: "ORD-quote-123",
    Amount:  0.5,
    Price:   0.013,
})
```

### Delete a Quote

```go
err := client.MMRfqDeleteQuote(ctx, &types.RfqQuoteDeleteParams{
    OrderID: "ORD-quote-123",
})
```

### List Active Quotes

```go
quotes, err := client.MMRfqQuotes(ctx)
for _, q := range quotes {
    fmt.Printf("Quote %s: %s %.4f @ %.4f on RFQ %s\n",
        q.OrderID, q.Direction, q.Amount, q.Price, q.RfqID)
}
```

### RfqOrder Response

```go
type RfqOrder struct {
    RfqID         string
    OrderID       string
    ClientOrderID *uint64
    Direction     enums.Direction
    Price         float64
    Amount        float64
    Label         string
    TradePrice    *float64
    TradeAmount   *float64
    DeleteReason  enums.RfqDeleteReason
    Event         enums.RfqOrderEvent
}
```

### RFQ Order Events

| Value | Constant | Description |
|-------|----------|-------------|
| `"Inserted"` | `RfqOrderEventInserted` | Quote inserted |
| `"Amended"` | `RfqOrderEventAmended` | Quote amended |
| `"Cancelled"` | `RfqOrderEventCancelled` | Quote cancelled |
| `"Filled"` | `RfqOrderEventFilled` | Quote filled |
| `"Existing"` | `RfqOrderEventExisting` | Existing on subscription start |

## Subscribing to RFQ Updates

### Requestor Side

```go
wsClient.OnRfqs(func(rfqs []types.Rfq) {
    for _, r := range rfqs {
        fmt.Printf("RFQ %s: event=%s\n", r.RfqID, r.Event)
        if r.QuotedBid != nil {
            fmt.Printf("  Bid: %.4f x %.4f\n", r.QuotedBid.Price, r.QuotedBid.Amount)
        }
    }
})
wsClient.SubscribePrivate(ctx, types.ChannelAccountRfqs)
```

### Market Maker Side

```go
// RFQ opportunities.
wsClient.OnMMRfqs(func(rfqs []types.Rfq) {
    for _, r := range rfqs {
        fmt.Printf("New RFQ: %s %d legs amount=%.4f\n",
            r.RfqID, len(r.Legs), r.Amount)
    }
})

// Quote status updates.
wsClient.OnMMRfqQuotes(func(quotes []types.RfqOrder) {
    for _, q := range quotes {
        fmt.Printf("Quote update: %s event=%s\n", q.OrderID, q.Event)
    }
})

wsClient.SubscribePrivate(ctx,
    types.ChannelMMRfqs,
    types.ChannelMMRfqQuotes,
)
```

### Public RFQ Channel

Monitor all RFQs (anonymized):

```go
err := wsClient.Subscribe(ctx, types.ChannelRfqs)
```

---

[< Conditional Orders](conditional-orders.md) | [Home](README.md) | [Enum Types >](enums.md)
