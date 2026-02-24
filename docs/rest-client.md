# REST Client

The REST client provides synchronous access to all Thalex API endpoints over HTTPS. It handles authentication, retries, and JSON serialization automatically.

## Creating a Client

```go
import (
    "github.com/amiwrpremium/go-thalex/config"
    "github.com/amiwrpremium/go-thalex/rest"
)

// Public-only client (no auth).
pub := rest.NewClient(
    config.WithNetwork(config.Testnet),
)

// Authenticated client.
client := rest.NewClient(
    config.WithNetwork(config.Production),
    config.WithCredentials(creds),
    config.WithMaxRetries(5),
)
```

See [Configuration](configuration.md) for all available options.

## Configuration Options

Options most relevant to the REST client:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithNetwork(n)` | `config.Network` | `Production` | API environment |
| `WithCredentials(c)` | `*auth.Credentials` | `nil` | API credentials |
| `WithHTTPClient(c)` | `*http.Client` | 30s timeout | Custom HTTP client |
| `WithMaxRetries(n)` | `int` | `3` | Max retry attempts |
| `WithRetryBaseWait(d)` | `time.Duration` | `500ms` | Base wait between retries |
| `WithUserAgent(ua)` | `string` | `"go-thalex/0.2.0"` | Custom user agent |
| `WithAccountNumber(a)` | `string` | `""` | Sub-account number |

## HTTP Retry Behavior

The REST client automatically retries failed requests with exponential backoff:

- **Max attempts:** 3 (configurable via `WithMaxRetries`)
- **Base wait:** 500ms (configurable via `WithRetryBaseWait`)
- **Backoff:** Exponential (500ms, 1s, 2s, ...)
- **Retried errors:** Network errors, 5xx server errors
- **Not retried:** 4xx client errors (including API errors)

## Error Handling

```go
import (
    "errors"
    "github.com/amiwrpremium/go-thalex/apierr"
)

ticker, err := client.Ticker(ctx, "BTC-PERPETUAL")
if err != nil {
    // Check for specific API error.
    if apiErr, ok := apierr.IsAPIError(err); ok {
        fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
        return
    }

    // Check for other error types.
    var connErr *apierr.ConnectionError
    if errors.As(err, &connErr) {
        fmt.Printf("Connection error: %s\n", connErr.Message)
        return
    }

    log.Fatal(err)
}
```

See [Error Handling](error-handling.md) for complete details.

## Public Endpoints -- Market Data

These endpoints do not require authentication.

### Instruments

```go
// All active instruments.
instruments, err := client.Instruments(ctx)

// All instruments including expired/inactive.
all, err := client.AllInstruments(ctx)

// Single instrument by name.
inst, err := client.Instrument(ctx, "BTC-PERPETUAL")
```

**Response type:** `types.Instrument`

```json
{
    "instrument_name": "BTC-PERPETUAL",
    "product": "FBTCUSD",
    "tick_size": 0.5,
    "volume_tick_size": 0.01,
    "min_order_amount": 0.01,
    "underlying": "BTCUSD",
    "type": "perpetual"
}
```

Helper methods on `Instrument`:
- `inst.IsOption()` -- true if option
- `inst.IsFuture()` -- true if future
- `inst.IsPerpetual()` -- true if perpetual
- `inst.IsCombination()` -- true if combination

### Ticker

```go
ticker, err := client.Ticker(ctx, "BTC-PERPETUAL")
fmt.Printf("Mark: %.2f\n", ticker.MarkPrice)

if mid := ticker.MidPrice(); mid != nil {
    fmt.Printf("Mid: %.2f\n", *mid)
}
if spread := ticker.Spread(); spread != nil {
    fmt.Printf("Spread: %.2f\n", *spread)
}
```

**Response type:** `types.Ticker`

```json
{
    "best_bid_price": 95100.5,
    "best_bid_amount": 1.25,
    "best_ask_price": 95101.0,
    "best_ask_amount": 0.80,
    "mark_price": 95100.75,
    "mark_timestamp": 1700000000.123,
    "volume_24h": 1250.5,
    "open_interest": 5432.1,
    "funding_rate": 0.0001
}
```

### Index Price

```go
idx, err := client.Index(ctx, "BTCUSD")
fmt.Printf("%s: %.2f\n", idx.IndexName, idx.Price)
```

### Order Book

```go
book, err := client.Book(ctx, "BTC-PERPETUAL")

for _, level := range book.Bids {
    fmt.Printf("Bid: %.2f x %.4f (outright: %.4f)\n",
        level.Price(), level.Amount(), level.OutrightAmount())
}
for _, level := range book.Asks {
    fmt.Printf("Ask: %.2f x %.4f\n", level.Price(), level.Amount())
}
```

**Response type:** `types.Book`

```json
{
    "bids": [[95100.5, 1.25, 1.25], [95100.0, 0.50, 0.50]],
    "asks": [[95101.0, 0.80, 0.80], [95101.5, 1.10, 1.10]],
    "last": 95100.5,
    "time": 1700000000.123
}
```

### System Info

```go
info, err := client.SystemInfo(ctx)
fmt.Printf("Environment: %s\n", info.Environment)
```

### Historical Data

```go
import "github.com/amiwrpremium/go-thalex/enums"

// Mark price OHLC data.
markData, err := client.MarkPriceHistoricalData(ctx,
    "BTC-PERPETUAL",
    1700000000, // from (unix timestamp)
    1700086400, // to
    enums.Resolution1h,
)

// Parse based on instrument type.
switch markData.InstrumentType {
case enums.InstrumentTypePerpetual:
    for _, d := range markData.PerpetualData() {
        fmt.Printf("Time=%.0f O=%.2f H=%.2f L=%.2f C=%.2f Funding=%.6f\n",
            d.Time, d.Open, d.High, d.Low, d.Close, d.FundingPayment)
    }
case enums.InstrumentTypeOption:
    for _, d := range markData.OptionData() {
        fmt.Printf("Price OHLC: %.2f/%.2f/%.2f/%.2f  IV: %.4f-%.4f\n",
            d.Open, d.High, d.Low, d.Close, d.IVLow, d.IVHigh)
    }
default:
    for _, d := range markData.FutureData() {
        fmt.Printf("%.2f/%.2f/%.2f/%.2f\n", d.Open, d.High, d.Low, d.Close)
    }
}

// Index price OHLC data.
indexData, err := client.IndexPriceHistoricalData(ctx,
    "BTCUSD", 1700000000, 1700086400, enums.Resolution1h,
)
for _, d := range indexData.Data() {
    fmt.Printf("Time=%.0f Close=%.2f\n", d.Time, d.Close)
}
```

## Private Endpoints -- Trading

All private endpoints require credentials. See [Trading](trading.md) for detailed usage.

### Insert Order

```go
order, err := client.Insert(ctx,
    types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
        WithPrice(95000).
        WithPostOnly(true).
        WithLabel("my-order"),
)
fmt.Printf("ID: %s Status: %s\n", order.OrderID, order.Status)
```

### Quick Buy / Sell (Market Orders)

```go
// Market buy.
order, err := client.Buy(ctx, "BTC-PERPETUAL", 0.01)

// Market sell.
order, err := client.Sell(ctx, "ETH-PERPETUAL", 0.1)
```

### Amend Order

```go
amended, err := client.Amend(ctx,
    types.NewAmendByOrderID("ord-123", 95100, 0.02),
)
```

### Cancel Order

```go
cancelled, err := client.Cancel(ctx, types.CancelByOrderID("ord-123"))

// Cancel by client order ID.
cancelled, err := client.Cancel(ctx, types.CancelByClientOrderID(42))
```

### Cancel All Orders

```go
n, err := client.CancelAll(ctx)
fmt.Printf("Cancelled %d orders\n", n)
```

### Open Orders

```go
// All open orders.
orders, err := client.OpenOrders(ctx, "")

// Filtered by instrument.
orders, err := client.OpenOrders(ctx, "BTC-PERPETUAL")
```

## Private Endpoints -- Account

### Portfolio

```go
positions, err := client.Portfolio(ctx)
for _, p := range positions {
    fmt.Printf("%s: position=%.4f pnl=%.4f\n",
        p.InstrumentName, p.Position, p.UnrealisedPnl)
}
```

### Account Summary

```go
summary, err := client.AccountSummary(ctx)
fmt.Printf("Margin: %.2f Required: %.2f Utilization: %.2f%%\n",
    summary.Margin, summary.RequiredMargin,
    summary.MarginUtilization()*100)
```

### Account Breakdown

```go
breakdown, err := client.AccountBreakdown(ctx)
```

### Margin

```go
// Required margin breakdown.
margin, err := client.RequiredMarginBreakdown(ctx)

// Margin impact of a hypothetical order.
impact, err := client.RequiredMarginForOrder(ctx, "BTC-PERPETUAL", 95000, 0.1)
fmt.Printf("Current margin: %.2f  With buy: %.2f  With sell: %.2f\n",
    impact.Current.RequiredMargin,
    impact.WithBuy.RequiredMargin,
    impact.WithSell.RequiredMargin)
```

## Private Endpoints -- History

```go
// Trade history.
trades, err := client.TradeHistory(ctx, &types.TradeHistoryParams{
    Limit: types.Ptr(50),
    Sort:  enums.SortDesc,
    InstrumentNames: []string{"BTC-PERPETUAL"},
})

// Order history.
orders, err := client.OrderHistory(ctx, &types.OrderHistoryParams{
    Limit: types.Ptr(20),
})

// Daily mark history.
marks, err := client.DailyMarkHistory(ctx, &types.DailyMarkHistoryParams{
    Limit: types.Ptr(10),
})

// Transaction history.
txns, err := client.TransactionHistory(ctx, &types.TransactionHistoryParams{
    Sort: enums.SortDesc,
})
```

## Private Endpoints -- Wallet

```go
// Get deposit addresses.
btcAddr, err := client.BTCDepositAddress(ctx)
ethAddr, err := client.ETHDepositAddress(ctx)

// Get deposits and withdrawals.
deposits, err := client.CryptoDeposits(ctx)
withdrawals, err := client.CryptoWithdrawals(ctx)

// Verify before withdrawing.
verify, err := client.VerifyWithdrawal(ctx, &types.WithdrawParams{
    AssetName:     "BTC",
    Amount:        0.01,
    TargetAddress: "bc1q...",
})

// Execute withdrawal.
w, err := client.Withdraw(ctx, &types.WithdrawParams{
    AssetName:     "BTC",
    Amount:        0.01,
    TargetAddress: "bc1q...",
})

// Internal transfer.
err = client.InternalTransfer(ctx, &types.InternalTransferParams{
    DestinationAccountNumber: "ACC-67890",
    Assets: []types.Asset{{AssetName: "BTC", Amount: 0.1}},
})
```

## Private Endpoints -- Bots

See [Bot Management](bots.md) for detailed examples.

```go
bots, err := client.Bots(ctx, true) // true = include inactive
bot, err := client.CreateSGSLBot(ctx, params)
err = client.CancelBot(ctx, "bot-id")
n, err := client.CancelAllBots(ctx)
```

## Private Endpoints -- Conditional Orders

See [Conditional Orders](conditional-orders.md) for detailed examples.

```go
orders, err := client.ConditionalOrders(ctx)
co, err := client.CreateConditionalOrder(ctx, params)
err = client.CancelConditionalOrder(ctx, "order-id")
n, err := client.CancelAllConditionalOrders(ctx)
```

## Private Endpoints -- RFQ

See [Request for Quote](rfq.md) for detailed examples.

```go
rfq, err := client.CreateRfq(ctx, params)
trades, err := client.TradeRfq(ctx, params)
rfqs, err := client.OpenRfqs(ctx)
err = client.CancelRfq(ctx, "rfq-id")
```

## Private Endpoints -- Notifications

```go
// Get inbox notifications.
limit := 20
result, err := client.NotificationsInbox(ctx, &limit)
for _, n := range result.Notifications {
    fmt.Printf("[%s] %s: %s\n", n.Category, n.Title, n.Message)
}

// Mark as read.
err = client.MarkNotificationAsRead(ctx, "notif-id", true)
```

## Full Endpoint Reference

### Public (No Auth)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `Instruments(ctx)` | `GET /public/instruments` | Active instruments |
| `AllInstruments(ctx)` | `GET /public/all_instruments` | All instruments |
| `Instrument(ctx, name)` | `GET /public/instrument` | Single instrument |
| `Ticker(ctx, name)` | `GET /public/ticker` | Instrument ticker |
| `Index(ctx, underlying)` | `GET /public/index` | Index price |
| `Book(ctx, name)` | `GET /public/book` | Order book |
| `SystemInfo(ctx)` | `GET /public/system_info` | System status |
| `MarkPriceHistoricalData(...)` | `GET /public/mark_price_historical_data` | Mark price OHLC |
| `IndexPriceHistoricalData(...)` | `GET /public/index_price_historical_data` | Index price OHLC |

### Private -- Trading

| Method | Endpoint | Description |
|--------|----------|-------------|
| `Insert(ctx, params)` | `POST /private/insert` | Place order |
| `Buy(ctx, name, amount)` | `POST /private/buy` | Market buy |
| `Sell(ctx, name, amount)` | `POST /private/sell` | Market sell |
| `Amend(ctx, params)` | `POST /private/amend` | Modify order |
| `Cancel(ctx, params)` | `POST /private/cancel` | Cancel order |
| `CancelAll(ctx)` | `POST /private/cancel_all` | Cancel all orders |
| `OpenOrders(ctx, name)` | `GET /private/open_orders` | List open orders |

### Private -- Account

| Method | Endpoint | Description |
|--------|----------|-------------|
| `Portfolio(ctx)` | `GET /private/portfolio` | Positions |
| `AccountSummary(ctx)` | `GET /private/account_summary` | Financial summary |
| `AccountBreakdown(ctx)` | `GET /private/account_breakdown` | Detailed breakdown |
| `RequiredMarginBreakdown(ctx)` | `GET /private/required_margin_breakdown` | Margin breakdown |
| `RequiredMarginForOrder(...)` | `GET /private/required_margin_for_order` | Margin impact |

### Private -- History

| Method | Endpoint | Description |
|--------|----------|-------------|
| `TradeHistory(ctx, params)` | `GET /private/trade_history` | Trade history |
| `OrderHistory(ctx, params)` | `GET /private/order_history` | Order history |
| `DailyMarkHistory(ctx, params)` | `GET /private/daily_mark_history` | Daily marks |
| `TransactionHistory(ctx, params)` | `GET /private/transaction_history` | Transactions |

### Private -- Wallet

| Method | Endpoint | Description |
|--------|----------|-------------|
| `CryptoDeposits(ctx)` | `GET /private/crypto_deposits` | Deposits |
| `CryptoWithdrawals(ctx)` | `GET /private/crypto_withdrawals` | Withdrawals |
| `BTCDepositAddress(ctx)` | `GET /private/btc_deposit_address` | BTC address |
| `ETHDepositAddress(ctx)` | `GET /private/eth_deposit_address` | ETH address |
| `VerifyWithdrawal(ctx, params)` | `POST /private/verify_withdrawal` | Verify withdrawal |
| `Withdraw(ctx, params)` | `POST /private/withdraw` | Execute withdrawal |
| `VerifyInternalTransfer(...)` | `POST /private/verify_internal_transfer` | Verify transfer |
| `InternalTransfer(ctx, params)` | `POST /private/internal_transfer` | Execute transfer |

### Private -- Bots

| Method | Endpoint | Description |
|--------|----------|-------------|
| `Bots(ctx, includeInactive)` | `GET /private/bots` | List bots |
| `CreateSGSLBot(ctx, params)` | `POST /private/create_bot` | Create SGSL bot |
| `CreateOCQBot(ctx, params)` | `POST /private/create_bot` | Create OCQ bot |
| `CreateLevelsBot(ctx, params)` | `POST /private/create_bot` | Create Levels bot |
| `CreateGridBot(ctx, params)` | `POST /private/create_bot` | Create Grid bot |
| `CreateDHedgeBot(ctx, params)` | `POST /private/create_bot` | Create DHedge bot |
| `CreateDFollowBot(ctx, params)` | `POST /private/create_bot` | Create DFollow bot |
| `CancelBot(ctx, botID)` | `POST /private/cancel_bot` | Cancel a bot |
| `CancelAllBots(ctx)` | `POST /private/cancel_all_bots` | Cancel all bots |

### Private -- Conditional Orders

| Method | Endpoint | Description |
|--------|----------|-------------|
| `ConditionalOrders(ctx)` | `GET /private/conditional_orders` | List conditional orders |
| `CreateConditionalOrder(ctx, params)` | `POST /private/create_conditional_order` | Create conditional |
| `CancelConditionalOrder(ctx, id)` | `POST /private/cancel_conditional_order` | Cancel conditional |
| `CancelAllConditionalOrders(ctx)` | `POST /private/cancel_all_conditional_orders` | Cancel all |

### Private -- RFQ

| Method | Endpoint | Description |
|--------|----------|-------------|
| `CreateRfq(ctx, params)` | `POST /private/create_rfq` | Create RFQ |
| `CancelRfq(ctx, rfqID)` | `POST /private/cancel_rfq` | Cancel RFQ |
| `TradeRfq(ctx, params)` | `POST /private/trade_rfq` | Execute RFQ trade |
| `OpenRfqs(ctx)` | `GET /private/open_rfqs` | List open RFQs |
| `RfqHistory(ctx, ...)` | `GET /private/rfq_history` | RFQ history |
| `MMRfqs(ctx)` | `GET /private/mm_rfqs` | MM RFQ opportunities |
| `MMRfqInsertQuote(ctx, params)` | `POST /private/mm_rfq_insert_quote` | Insert RFQ quote |
| `MMRfqAmendQuote(ctx, params)` | `POST /private/mm_rfq_amend_quote` | Amend RFQ quote |
| `MMRfqDeleteQuote(ctx, params)` | `POST /private/mm_rfq_delete_quote` | Delete RFQ quote |
| `MMRfqQuotes(ctx)` | `GET /private/mm_rfq_quotes` | Active RFQ quotes |

### Private -- Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| `NotificationsInbox(ctx, limit)` | `GET /private/notifications_inbox` | Inbox |
| `MarkNotificationAsRead(ctx, id, read)` | `POST /private/mark_inbox_notification_as_read` | Mark read/unread |

---

[< Authentication](authentication.md) | [Home](README.md) | [WebSocket Client >](ws-client.md)
