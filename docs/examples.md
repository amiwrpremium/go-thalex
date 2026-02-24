# Examples

The `examples/` directory contains five runnable programs demonstrating different SDK features. Each example is a standalone `main` package.

## Environment Variables

All authenticated examples require these environment variables:

| Variable | Description |
|----------|-------------|
| `THALEX_KEY_ID` | Your API key identifier |
| `THALEX_PEM_PATH` | Absolute path to your RSA private key PEM file |

Set them before running:

```bash
export THALEX_KEY_ID="my-key-id"
export THALEX_PEM_PATH="/path/to/private_key.pem"
```

## Running Examples

```bash
# From the repository root:
go run ./examples/rest_basic/
go run ./examples/ws_subscriptions/
go run ./examples/ws_trading/
go run ./examples/market_making/
go run ./examples/bot_management/
```

All examples connect to **testnet** by default.

---

## 1. REST Basic (`examples/rest_basic/`)

**File:** `examples/rest_basic/main.go`

Demonstrates basic REST API usage including public market data and authenticated trading.

**What it does:**
- Creates an unauthenticated client for public endpoints
- Fetches all active instruments
- Fetches the BTC-PERPETUAL ticker (mark price, mid price)
- Fetches the BTC-PERPETUAL order book (best bid/ask)
- If credentials are provided: fetches account summary, places a limit buy order, then cancels it

**Authentication:** Optional (public endpoints work without credentials)

**Key concepts demonstrated:**
- `rest.NewClient()` with options
- `client.Instruments()`, `client.Ticker()`, `client.Book()`
- `types.NewBuyOrderParams()` with builder pattern
- `types.CancelByOrderID()`
- `Ticker.MidPrice()` helper method
- `BookLevel.Price()` and `BookLevel.Amount()` accessors

```bash
# Public endpoints only:
go run ./examples/rest_basic/

# With authentication:
export THALEX_KEY_ID="my-key-id"
export THALEX_PEM_PATH="/path/to/key.pem"
go run ./examples/rest_basic/
```

**Related docs:** [REST Client](rest-client.md), [Trading](trading.md)

---

## 2. WebSocket Subscriptions (`examples/ws_subscriptions/`)

**File:** `examples/ws_subscriptions/main.go`

Demonstrates real-time data streaming via WebSocket subscriptions.

**What it does:**
- Creates a WebSocket client with auto-reconnect enabled
- Registers an error handler
- Subscribes to public channels:
  - BTC-PERPETUAL ticker (100ms delay)
  - ETH-PERPETUAL lightweight ticker (100ms delay)
  - BTC-PERPETUAL order book (grouping=1, 5 levels, 100ms delay)
  - BTCUSD price index
  - Instrument changes
- Streams data until Ctrl+C

**Authentication:** Optional (shown in code but not required for public channels)

**Key concepts demonstrated:**
- `ws.NewClient()` with `WithWSReconnect(true)`
- Channel helper functions: `TickerChannel()`, `LWTChannel()`, `BookChannel()`, `PriceIndexChannel()`
- Typed handlers: `OnTicker()`, `OnLWT()`, `OnBook()`, `OnPriceIndex()`, `OnInstruments()`
- `Subscribe()` for multiple channels
- Signal-based graceful shutdown

```bash
go run ./examples/ws_subscriptions/
# Press Ctrl+C to stop
```

**Related docs:** [WebSocket Client](ws-client.md), [Subscriptions](subscriptions.md)

---

## 3. WebSocket Trading (`examples/ws_trading/`)

**File:** `examples/ws_trading/main.go`

Demonstrates low-latency trading over WebSocket.

**What it does:**
- Creates an authenticated WebSocket client
- Connects and logs in
- Enables cancel-on-disconnect for safety
- Places a limit buy order for BTC-PERPETUAL
- Amends the order (changes price)
- Cancels the order

**Authentication:** Required

**Key concepts demonstrated:**
- `wsClient.Connect()` and `wsClient.Login()`
- `wsClient.SetCancelOnDisconnect()`
- `wsClient.Insert()` with builder pattern
- `wsClient.Amend()` with `NewAmendByOrderID()`
- `wsClient.Cancel()` with `CancelByOrderID()`

```bash
export THALEX_KEY_ID="my-key-id"
export THALEX_PEM_PATH="/path/to/key.pem"
go run ./examples/ws_trading/
```

**Related docs:** [WebSocket Client](ws-client.md), [Trading](trading.md)

---

## 4. Market Making (`examples/market_making/`)

**File:** `examples/market_making/main.go`

Demonstrates market making operations (WebSocket-only).

**What it does:**
- Creates an authenticated WebSocket client
- Connects, logs in, enables cancel-on-disconnect
- Configures MM protection for the FBTCUSD product
- Sends single-level mass quotes for BTC-PERPETUAL
- Sends multi-level mass quotes with 2 price levels per side
- Cancels all mass quotes

**Authentication:** Required

**Key concepts demonstrated:**
- `wsClient.SetMMProtection()` configuration
- `types.NewSingleLevelQuote()` for simple quotes
- `types.NewDoubleSidedQuote()` with `QuoteLevel` for multi-level quotes
- `types.NewMassQuoteParams()` with `WithPostOnly()` and `WithLabel()`
- `wsClient.MassQuote()` and response handling
- `wsClient.CancelMassQuote()` cleanup
- `DoubleSidedQuoteResult` with success/fail counts and error details

```bash
export THALEX_KEY_ID="my-key-id"
export THALEX_PEM_PATH="/path/to/key.pem"
go run ./examples/market_making/
```

**Related docs:** [Market Making](market-making.md)

---

## 5. Bot Management (`examples/bot_management/`)

**File:** `examples/bot_management/main.go`

Demonstrates creating and managing server-side trading bots.

**What it does:**
- Creates an authenticated REST client
- Creates an SGSL (Signal Go Stop-Loss) bot on BTC-PERPETUAL
- Creates a Grid bot with 5 price levels on BTC-PERPETUAL
- Creates a Delta Hedger bot on BTC-PERPETUAL
- Lists all bots (including inactive)
- Cancels all bots

**Authentication:** Required

**Key concepts demonstrated:**
- `types.NewSGSLBotParams()` with `WithMaxSlippage()` and `WithLabel()`
- `types.NewGridBotParams()` with grid price levels
- `types.NewDHedgeBotParams()` with `WithTargetDelta()`, `WithThreshold()`, `WithEndTime()`
- `client.CreateSGSLBot()`, `client.CreateGridBot()`, `client.CreateDHedgeBot()`
- `client.Bots()` with `includeInactive` flag
- `client.CancelAllBots()` cleanup
- Bot response fields: `BotID`, `Strategy`, `Status`, `InstrumentName`

```bash
export THALEX_KEY_ID="my-key-id"
export THALEX_PEM_PATH="/path/to/key.pem"
go run ./examples/bot_management/
```

**Related docs:** [Bot Management](bots.md)

---

## Example Quick Reference

| Example | Client | Auth Required | Key Feature |
|---------|--------|---------------|-------------|
| `rest_basic` | REST | Optional | Market data + basic trading |
| `ws_subscriptions` | WebSocket | Optional | Real-time data streaming |
| `ws_trading` | WebSocket | Yes | Low-latency order management |
| `market_making` | WebSocket | Yes | Mass quotes + MM protection |
| `bot_management` | REST | Yes | Server-side bot creation |

---

[< Configuration](configuration.md) | [Home](README.md)
