package types

import (
	"fmt"

	"github.com/amiwrpremium/go-thalex/enums"
)

// BookUpdate represents an order book snapshot from a book subscription.
type BookUpdate struct {
	Bids   []BookLevel `json:"bids"`
	Asks   []BookLevel `json:"asks"`
	Last   *float64    `json:"last,omitempty"`
	Time   float64     `json:"time"`
	Trades []BookTrade `json:"trades,omitempty"`
}

// BookTrade represents a trade tick included in a book update.
type BookTrade struct {
	Direction enums.Direction `json:"d"`
	Price     float64         `json:"p"`
	Amount    float64         `json:"a"`
	Time      float64         `json:"t"`
}

// LightweightTicker represents a condensed ticker from an lwt subscription.
type LightweightTicker struct {
	BestBidPrice  *float64 `json:"best_bid_price"`
	BestAskPrice  *float64 `json:"best_ask_price"`
	MarkPrice     float64  `json:"mark_price"`
	IV            *float64 `json:"iv,omitempty"`
	LastPrice     *float64 `json:"last_price,omitempty"`
	BestBidAmount *float64 `json:"best_bid_amount,omitempty"`
	BestAskAmount *float64 `json:"best_ask_amount,omitempty"`
}

// RecentTrade represents a trade from a recent_trades subscription.
type RecentTrade struct {
	TradeID        string          `json:"trade_id"`
	InstrumentName string          `json:"instrument_name"`
	Direction      enums.Direction `json:"direction"`
	Price          float64         `json:"price"`
	Amount         float64         `json:"amount"`
	Time           float64         `json:"time"`
	TradeType      enums.TradeType `json:"trade_type"`
	Index          *float64        `json:"index,omitempty"`
}

// UnderlyingStatistics represents statistics for a single underlying.
type UnderlyingStatistics struct {
	Underlying string  `json:"underlying"`
	Time       float64 `json:"time"`
	OpenInterest
}

// OpenInterestSummary represents aggregated open interest for an underlying.
type OpenInterestSummary struct {
	TotalOI     float64        `json:"total_oi"`
	PutOI       float64        `json:"put_oi"`
	CallOI      float64        `json:"call_oi"`
	FuturesOI   float64        `json:"futures_oi"`
	Expirations []ExpirationOI `json:"expirations,omitempty"`
}

// ExpirationOI represents open interest for a single expiration.
type ExpirationOI struct {
	ExpiryDate string  `json:"expiry_date"`
	TotalOI    float64 `json:"total_oi"`
	PutOI      float64 `json:"put_oi"`
	CallOI     float64 `json:"call_oi"`
	FuturesOI  float64 `json:"futures_oi"`
}

// OpenInterest is embedded in types that carry open interest data.
type OpenInterest struct {
	Futures float64 `json:"futures,omitempty"`
	Calls   float64 `json:"calls,omitempty"`
	Puts    float64 `json:"puts,omitempty"`
}

// BasePrice represents a forward price for a specific expiration.
type BasePrice struct {
	Underlying string  `json:"underlying"`
	Expiration string  `json:"expiration"`
	Price      float64 `json:"price"`
	Time       float64 `json:"time"`
}

// InstrumentChange represents an instrument activation or deactivation event.
type InstrumentChange struct {
	Instruments []Instrument `json:"instruments"`
}

// IndexComponent represents a single component of an index price.
type IndexComponent struct {
	Exchange string  `json:"exchange"`
	Price    float64 `json:"price"`
	Weight   float64 `json:"weight"`
	Stale    bool    `json:"stale"`
}

// IndexComponents represents the full index composition.
type IndexComponents struct {
	IndexName  string           `json:"index_name"`
	Price      float64          `json:"price"`
	Time       float64          `json:"time"`
	Components []IndexComponent `json:"components"`
}

// MMProtectionUpdate represents a market maker protection status update.
type MMProtectionUpdate struct {
	Product enums.Product            `json:"product"`
	Reason  enums.MMProtectionReason `json:"reason"`
	Time    float64                  `json:"time"`
}

// SubscriptionNotification wraps a notification from a subscription channel.
type SubscriptionNotification struct {
	Channel string `json:"channel"`
	Data    any    `json:"data"`
}

// --- Channel name helpers ---

// BookChannel returns the subscription channel name for an order book feed.
//
//	BookChannel("BTC-PERPETUAL", 1, 10, enums.Delay100ms) => "book.BTC-PERPETUAL.1.10.100ms"
func BookChannel(instrument string, grouping int, nlevels int, delay enums.Delay) string {
	return fmt.Sprintf("book.%s.%d.%d.%s", instrument, grouping, nlevels, delay)
}

// TickerChannel returns the subscription channel name for a ticker feed.
//
//	TickerChannel("BTC-PERPETUAL", enums.Delay100ms) => "ticker.BTC-PERPETUAL.100ms"
func TickerChannel(instrument string, delay enums.Delay) string {
	return fmt.Sprintf("ticker.%s.%s", instrument, delay)
}

// LWTChannel returns the subscription channel name for a lightweight ticker feed.
//
//	LWTChannel("BTC-PERPETUAL", enums.Delay100ms) => "lwt.BTC-PERPETUAL.100ms"
func LWTChannel(instrument string, delay enums.Delay) string {
	return fmt.Sprintf("lwt.%s.%s", instrument, delay)
}

// RecentTradesChannel returns the subscription channel name for recent trades.
//
//	RecentTradesChannel("BTCUSD", enums.RecentTradesCategoryAll) => "recent_trades.BTCUSD.all"
func RecentTradesChannel(target string, category enums.RecentTradesCategory) string {
	return fmt.Sprintf("recent_trades.%s.%s", target, category)
}

// PriceIndexChannel returns the subscription channel name for an index price feed.
//
//	PriceIndexChannel("BTCUSD") => "price_index.BTCUSD"
func PriceIndexChannel(underlying string) string {
	return "price_index." + underlying
}

// UnderlyingStatisticsChannel returns the subscription channel for underlying statistics.
//
//	UnderlyingStatisticsChannel("BTCUSD") => "underlying_statistics.BTCUSD"
func UnderlyingStatisticsChannel(underlying string) string {
	return "underlying_statistics." + underlying
}

// BasePriceChannel returns the subscription channel for a forward price.
//
//	BasePriceChannel("BTCUSD", "2025-03-28") => "base_price.BTCUSD.2025-03-28"
func BasePriceChannel(underlying, expiration string) string {
	return fmt.Sprintf("base_price.%s.%s", underlying, expiration)
}

// IndexComponentsChannel returns the subscription channel for index components.
//
//	IndexComponentsChannel("BTCUSD") => "index_components.BTCUSD"
func IndexComponentsChannel(underlying string) string {
	return "index_components." + underlying
}

// Constant channel names for parameterless subscriptions.
const (
	ChannelInstruments         = "instruments"
	ChannelRfqs                = "rfqs"
	ChannelSystem              = "system"
	ChannelBanners             = "banners"
	ChannelAccountOrders       = "account.orders"
	ChannelAccountPersistent   = "account.persistent_orders"
	ChannelSessionOrders       = "session.orders"
	ChannelAccountTradeHistory = "account.trade_history"
	ChannelAccountOrderHistory = "account.order_history"
	ChannelAccountPortfolio    = "account.portfolio"
	ChannelAccountSummary      = "account.summary"
	ChannelAccountRfqs         = "account.rfqs"
	ChannelAccountRfqHistory   = "account.rfq_history"
	ChannelAccountConditional  = "account.conditional_orders"
	ChannelAccountBots         = "account.bots"
	ChannelUserNotifications   = "user.inbox_notifications"
	ChannelSessionMMProtection = "session.mm_protection"
	ChannelMMRfqs              = "mm.rfqs"
	ChannelMMRfqQuotes         = "mm.rfq_quotes"
)
