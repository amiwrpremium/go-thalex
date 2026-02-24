# Enum Types

The `enums` package defines all enumerated types used throughout the SDK. Each type is a named `string` type with constants, standard methods, and in some cases domain-specific helper methods.

Import path:

```go
import "github.com/amiwrpremium/go-thalex/enums"
```

## Standard Methods

Every enum type implements these methods:

| Method | Description |
|--------|-------------|
| `String() string` | Returns the string representation |
| `IsValid() bool` | Returns true if the value is a recognized constant |
| `<Type>Values() []<Type>` | Package-level function returning all valid values |

```go
d := enums.DirectionBuy
fmt.Println(d.String())   // "buy"
fmt.Println(d.IsValid())  // true

for _, v := range enums.DirectionValues() {
    fmt.Println(v) // "buy", "sell"
}
```

## Complete Enum Reference

### Direction

Trade or order side.

| Constant | Value |
|----------|-------|
| `DirectionBuy` | `"buy"` |
| `DirectionSell` | `"sell"` |

**Extra methods:**
- `Opposite() Direction` -- returns the opposite direction

```go
d := enums.DirectionBuy
fmt.Println(d.Opposite()) // "sell"
```

### OrderType

| Constant | Value |
|----------|-------|
| `OrderTypeLimit` | `"limit"` |
| `OrderTypeMarket` | `"market"` |

### TimeInForce

| Constant | Value |
|----------|-------|
| `TimeInForceGoodTillCancelled` | `"good_till_cancelled"` |
| `TimeInForceImmediateOrCancel` | `"immediate_or_cancel"` |

### OrderStatusValue

| Constant | Value |
|----------|-------|
| `OrderStatusOpen` | `"open"` |
| `OrderStatusPartiallyFilled` | `"partially_filled"` |
| `OrderStatusCancelled` | `"cancelled"` |
| `OrderStatusCancelledPartiallyFilled` | `"cancelled_partially_filled"` |
| `OrderStatusFilled` | `"filled"` |

**Extra methods:**
- `IsActive() bool` -- true for `open` and `partially_filled`
- `IsFinal() bool` -- true for terminal states

### InstrumentType

| Constant | Value |
|----------|-------|
| `InstrumentTypePerpetual` | `"perpetual"` |
| `InstrumentTypeFuture` | `"future"` |
| `InstrumentTypeOption` | `"option"` |
| `InstrumentTypeCombination` | `"combination"` |

### OptionType

| Constant | Value |
|----------|-------|
| `OptionTypeCall` | `"call"` |
| `OptionTypePut` | `"put"` |

### Collar

Price collar handling.

| Constant | Value | Description |
|----------|-------|-------------|
| `CollarIgnore` | `"ignore"` | Ignore collar violations |
| `CollarReject` | `"reject"` | Reject order if outside collar |
| `CollarClamp` | `"clamp"` | Clamp price to collar boundary |

### Target

Trigger source for conditional orders and bot signals.

| Constant | Value |
|----------|-------|
| `TargetLast` | `"last"` |
| `TargetMark` | `"mark"` |
| `TargetIndex` | `"index"` |

### Delay

Subscription throttle interval.

| Constant | Value |
|----------|-------|
| `DelayNone` | `"raw"` |
| `Delay100ms` | `"100ms"` |
| `Delay1000ms` | `"1000ms"` |

### MakerTaker

| Constant | Value |
|----------|-------|
| `MakerTakerMaker` | `"maker"` |
| `MakerTakerTaker` | `"taker"` |

### TradeType

| Constant | Value |
|----------|-------|
| `TradeTypeNormal` | `"normal"` |
| `TradeTypeBlock` | `"block"` |
| `TradeTypeCombo` | `"combo"` |
| `TradeTypeAmend` | `"amend"` |
| `TradeTypeDelete` | `"delete"` |
| `TradeTypeInternalTransfer` | `"internal_transfer"` |
| `TradeTypeExpiration` | `"expiration"` |
| `TradeTypeDailyMark` | `"daily_mark"` |
| `TradeTypeRfq` | `"rfq"` |
| `TradeTypeLiquidation` | `"liquidation"` |

### ChangeReason

Why an order status changed.

| Constant | Value |
|----------|-------|
| `ChangeReasonExisting` | `"existing"` |
| `ChangeReasonInsert` | `"insert"` |
| `ChangeReasonAmend` | `"amend"` |
| `ChangeReasonCancel` | `"cancel"` |
| `ChangeReasonFill` | `"fill"` |

### DeleteReason

Why an order was deleted.

| Constant | Value |
|----------|-------|
| `DeleteReasonClientCancel` | `"client_cancel"` |
| `DeleteReasonClientBulkCancel` | `"client_bulk_cancel"` |
| `DeleteReasonSessionEnd` | `"session_end"` |
| `DeleteReasonInstrumentDeactivated` | `"instrument_deactivated"` |
| `DeleteReasonMMProtection` | `"mm_protection"` |
| `DeleteReasonFailover` | `"failover"` |
| `DeleteReasonMarginBreach` | `"margin_breach"` |
| `DeleteReasonFilled` | `"filled"` |
| `DeleteReasonImmediateCancel` | `"immediate_cancel"` |
| `DeleteReasonAdminCancel` | `"admin_cancel"` |

### InsertReason

Why an order was inserted.

| Constant | Value |
|----------|-------|
| `InsertReasonClientRequest` | `"client_request"` |
| `InsertReasonConditionalOrder` | `"conditional_order"` |
| `InsertReasonLiquidation` | `"liquidation"` |

### ConditionalOrderStatus

| Constant | Value |
|----------|-------|
| `ConditionalOrderStatusCreated` | `"created"` |
| `ConditionalOrderStatusActive` | `"active"` |
| `ConditionalOrderStatusConverted` | `"converted"` |
| `ConditionalOrderStatusRejected` | `"rejected"` |
| `ConditionalOrderStatusCancelRequested` | `"cancel requested"` |
| `ConditionalOrderStatusCancelled` | `"cancelled"` |

**Extra methods:**
- `IsActive() bool` -- true for `created` and `active`

### BotStrategy

| Constant | Value |
|----------|-------|
| `BotStrategySGSL` | `"sgsl"` |
| `BotStrategyOCQ` | `"ocq"` |
| `BotStrategyLevels` | `"levels"` |
| `BotStrategyGrid` | `"grid"` |
| `BotStrategyDHedge` | `"dhedge"` |
| `BotStrategyDFollow` | `"dfollow"` |

### BotStatus

| Constant | Value |
|----------|-------|
| `BotStatusActive` | `"active"` |
| `BotStatusStopped` | `"stopped"` |

**Extra methods:**
- `IsActive() bool` -- true for `active`
- `IsFinal() bool` -- true for `stopped`

### BotStopReason

| Constant | Value |
|----------|-------|
| `BotStopReasonClientCancel` | `"client_cancel"` |
| `BotStopReasonClientBulkCancel` | `"client_bulk_cancel"` |
| `BotStopReasonEndTime` | `"end_time"` |
| `BotStopReasonInstrumentDeactivated` | `"instrument_deactivated"` |
| `BotStopReasonMarginBreach` | `"margin_breach"` |
| `BotStopReasonAdminCancel` | `"admin_cancel"` |
| `BotStopReasonConflict` | `"conflict"` |
| `BotStopReasonStrategy` | `"strategy"` |

### DepositStatus

| Constant | Value |
|----------|-------|
| `DepositStatusUnconfirmed` | `"unconfirmed"` |
| `DepositStatusConfirmed` | `"confirmed"` |

**Extra methods:**
- `IsPending() bool` -- true for `unconfirmed`
- `IsFinal() bool` -- true for `confirmed`

### WithdrawalStatus

| Constant | Value |
|----------|-------|
| `WithdrawalStatusPending` | `"pending"` |
| `WithdrawalStatusAwaitingConfirmation` | `"awaiting_confirmation"` |
| `WithdrawalStatusExecuting` | `"executing"` |
| `WithdrawalStatusExecuted` | `"executed"` |
| `WithdrawalStatusRejected` | `"rejected"` |

**Extra methods:**
- `IsPending() bool` -- true for `pending`, `awaiting_confirmation`, `executing`
- `IsFinal() bool` -- true for `executed`, `rejected`

### Resolution

Historical data time resolution.

| Constant | Value |
|----------|-------|
| `Resolution1m` | `"1m"` |
| `Resolution5m` | `"5m"` |
| `Resolution15m` | `"15m"` |
| `Resolution30m` | `"30m"` |
| `Resolution1h` | `"1h"` |
| `Resolution1d` | `"1d"` |
| `Resolution1w` | `"1w"` |

### Sort

Pagination sort order.

| Constant | Value |
|----------|-------|
| `SortAsc` | `"asc"` |
| `SortDesc` | `"desc"` |

### Product

Product group identifier (free-form string type, no predefined constants).

```go
type Product string
```

Used for MM protection configuration. Examples: `"FBTCUSD"`, `"OBTCUSD"`, `"FETHUSD"`.

### Severity

Banner/notification severity.

| Constant | Value |
|----------|-------|
| `SeverityInfo` | `"info"` |
| `SeverityWarning` | `"warning"` |
| `SeverityCritical` | `"critical"` |

### DisplayType

Notification display style.

| Constant | Value |
|----------|-------|
| `DisplayTypeSuccess` | `"success"` |
| `DisplayTypeFailure` | `"failure"` |
| `DisplayTypeInfo` | `"info"` |
| `DisplayTypeWarning` | `"warning"` |
| `DisplayTypeCritical` | `"critical"` |

### RecentTradesCategory

| Constant | Value |
|----------|-------|
| `RecentTradesCategoryAll` | `"all"` |
| `RecentTradesCategoryNormal` | `"normal"` |
| `RecentTradesCategoryBlock` | `"block"` |
| `RecentTradesCategoryCombo` | `"combo"` |

### STPLevel

Self-trade prevention scope.

| Constant | Value |
|----------|-------|
| `STPLevelAccount` | `"account"` |
| `STPLevelCustomer` | `"customer"` |
| `STPLevelSubaccount` | `"subaccount"` |

### STPAction

Self-trade prevention action.

| Constant | Value |
|----------|-------|
| `STPActionCancelAggressor` | `"cancel_aggressor"` |
| `STPActionCancelPassive` | `"cancel_passive"` |
| `STPActionCancelBoth` | `"cancel_both"` |

### MMProtectionReason

| Constant | Value |
|----------|-------|
| `MMProtectionReasonTriggered` | `"triggered"` |
| `MMProtectionReasonReset` | `"reset"` |

### SystemEventType

| Constant | Value |
|----------|-------|
| `SystemEventTypeReconnect` | `"reconnect"` |

### RfqEvent

RFQ lifecycle event.

| Constant | Value |
|----------|-------|
| `RfqEventCreated` | `"Created"` |
| `RfqEventCancelled` | `"Cancelled"` |
| `RfqEventTraded` | `"Traded"` |
| `RfqEventExisting` | `"Existing"` |

### RfqInsertReason

| Constant | Value |
|----------|-------|
| `RfqInsertReasonClientRequest` | `"client_request"` |
| `RfqInsertReasonLiquidation` | `"liquidation"` |

### RfqDeleteReason

| Constant | Value |
|----------|-------|
| `RfqDeleteReasonClientCancel` | `"client_cancel"` |
| `RfqDeleteReasonSessionEnd` | `"session_end"` |
| `RfqDeleteReasonInstrumentDeactivated` | `"instrument_deactivated"` |
| `RfqDeleteReasonMMProtection` | `"mm_protection"` |
| `RfqDeleteReasonFailover` | `"failover"` |
| `RfqDeleteReasonMarginBreach` | `"margin_breach"` |
| `RfqDeleteReasonFilled` | `"filled"` |

### RfqOrderEvent

| Constant | Value |
|----------|-------|
| `RfqOrderEventInserted` | `"Inserted"` |
| `RfqOrderEventAmended` | `"Amended"` |
| `RfqOrderEventCancelled` | `"Cancelled"` |
| `RfqOrderEventFilled` | `"Filled"` |
| `RfqOrderEventExisting` | `"Existing"` |

## Usage in Code

### In Order Parameters

```go
params := types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
    WithOrderType(enums.OrderTypeLimit).
    WithTimeInForce(enums.TimeInForceGoodTillCancelled).
    WithCollar(enums.CollarReject)
```

### In Response Processing

```go
for _, order := range orders {
    if order.Status.IsActive() {
        fmt.Printf("Active: %s\n", order.OrderID)
    }
    if order.DeleteReason == enums.DeleteReasonMarginBreach {
        fmt.Printf("Margin breach: %s\n", order.OrderID)
    }
}
```

### In Switch Statements

```go
switch order.Status {
case enums.OrderStatusOpen:
    // Handle open
case enums.OrderStatusFilled:
    // Handle filled
case enums.OrderStatusCancelled, enums.OrderStatusCancelledPartiallyFilled:
    // Handle cancelled
}
```

---

[< Request for Quote](rfq.md) | [Home](README.md) | [Error Handling >](error-handling.md)
