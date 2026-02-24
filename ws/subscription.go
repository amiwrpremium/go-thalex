package ws

import (
	"context"
	"encoding/json"

	"github.com/amiwrpremium/go-thalex/types"
)

// Subscribe subscribes to public channels.
func (ws *Client) Subscribe(ctx context.Context, channels ...string) error {
	return ws.callNoResult(ctx, "public/subscribe", map[string]any{
		"channels": channels,
	})
}

// SubscribePrivate subscribes to private channels (requires login).
func (ws *Client) SubscribePrivate(ctx context.Context, channels ...string) error {
	return ws.callNoResult(ctx, "private/subscribe", map[string]any{
		"channels": channels,
	})
}

// Unsubscribe unsubscribes from public channels.
func (ws *Client) Unsubscribe(ctx context.Context, channels ...string) error {
	ws.subMu.Lock()
	for _, ch := range channels {
		delete(ws.handlers, ch)
	}
	ws.subMu.Unlock()
	return ws.callNoResult(ctx, "public/unsubscribe", map[string]any{
		"channels": channels,
	})
}

// UnsubscribePrivate unsubscribes from private channels.
func (ws *Client) UnsubscribePrivate(ctx context.Context, channels ...string) error {
	ws.subMu.Lock()
	for _, ch := range channels {
		delete(ws.handlers, ch)
	}
	ws.subMu.Unlock()
	return ws.callNoResult(ctx, "private/unsubscribe", map[string]any{
		"channels": channels,
	})
}

// --- Typed subscription handler registration ---

// OnBook registers a handler for order book updates on the given channel.
func (ws *Client) OnBook(channel string, fn func(types.BookUpdate)) {
	ws.subMu.Lock()
	ws.handlers[channel] = fn
	ws.subMu.Unlock()
}

// OnTicker registers a handler for ticker updates on the given channel.
func (ws *Client) OnTicker(channel string, fn func(types.Ticker)) {
	ws.subMu.Lock()
	ws.handlers[channel] = fn
	ws.subMu.Unlock()
}

// OnLWT registers a handler for lightweight ticker updates.
func (ws *Client) OnLWT(channel string, fn func(types.LightweightTicker)) {
	ws.subMu.Lock()
	ws.handlers[channel] = fn
	ws.subMu.Unlock()
}

// OnRecentTrades registers a handler for recent trade notifications.
func (ws *Client) OnRecentTrades(channel string, fn func([]types.RecentTrade)) {
	ws.subMu.Lock()
	ws.handlers[channel] = fn
	ws.subMu.Unlock()
}

// OnPriceIndex registers a handler for index price updates.
func (ws *Client) OnPriceIndex(channel string, fn func(types.IndexPrice)) {
	ws.subMu.Lock()
	ws.handlers[channel] = fn
	ws.subMu.Unlock()
}

// OnInstruments registers a handler for instrument change notifications.
func (ws *Client) OnInstruments(fn func([]types.Instrument)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelInstruments] = fn
	ws.subMu.Unlock()
}

// OnOrders registers a handler for order status updates.
func (ws *Client) OnOrders(fn func([]types.OrderStatus)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountOrders] = fn
	ws.subMu.Unlock()
}

// OnPersistentOrders registers a handler for persistent order updates.
func (ws *Client) OnPersistentOrders(fn func([]types.OrderStatus)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountPersistent] = fn
	ws.subMu.Unlock()
}

// OnSessionOrders registers a handler for session order updates.
func (ws *Client) OnSessionOrders(fn func([]types.OrderStatus)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelSessionOrders] = fn
	ws.subMu.Unlock()
}

// OnPortfolio registers a handler for portfolio updates.
func (ws *Client) OnPortfolio(fn func([]types.PortfolioEntry)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountPortfolio] = fn
	ws.subMu.Unlock()
}

// OnAccountSummary registers a handler for account summary updates.
func (ws *Client) OnAccountSummary(fn func(types.AccountSummary)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountSummary] = fn
	ws.subMu.Unlock()
}

// OnTradeHistory registers a handler for trade history notifications.
func (ws *Client) OnTradeHistory(fn func([]types.Trade)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountTradeHistory] = fn
	ws.subMu.Unlock()
}

// OnOrderHistory registers a handler for order history notifications.
func (ws *Client) OnOrderHistory(fn func([]types.OrderHistory)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountOrderHistory] = fn
	ws.subMu.Unlock()
}

// OnConditionalOrders registers a handler for conditional order updates.
func (ws *Client) OnConditionalOrders(fn func([]types.ConditionalOrder)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountConditional] = fn
	ws.subMu.Unlock()
}

// OnBots registers a handler for bot status updates.
func (ws *Client) OnBots(fn func([]types.Bot)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountBots] = fn
	ws.subMu.Unlock()
}

// OnRfqs registers a handler for RFQ notifications.
func (ws *Client) OnRfqs(fn func([]types.Rfq)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelAccountRfqs] = fn
	ws.subMu.Unlock()
}

// OnMMRfqs registers a handler for market maker RFQ notifications.
func (ws *Client) OnMMRfqs(fn func([]types.Rfq)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelMMRfqs] = fn
	ws.subMu.Unlock()
}

// OnMMRfqQuotes registers a handler for market maker RFQ quote updates.
func (ws *Client) OnMMRfqQuotes(fn func([]types.RfqOrder)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelMMRfqQuotes] = fn
	ws.subMu.Unlock()
}

// OnMMProtection registers a handler for market maker protection updates.
func (ws *Client) OnMMProtection(fn func(types.MMProtectionUpdate)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelSessionMMProtection] = fn
	ws.subMu.Unlock()
}

// OnNotifications registers a handler for inbox notification updates.
func (ws *Client) OnNotifications(fn func(types.Notification)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelUserNotifications] = fn
	ws.subMu.Unlock()
}

// OnSystemEvent registers a handler for system events.
func (ws *Client) OnSystemEvent(fn func(types.SystemEvent)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelSystem] = fn
	ws.subMu.Unlock()
}

// OnBanners registers a handler for banner updates.
func (ws *Client) OnBanners(fn func([]types.Banner)) {
	ws.subMu.Lock()
	ws.handlers[types.ChannelBanners] = fn
	ws.subMu.Unlock()
}

// OnRaw registers a raw JSON handler for any channel.
func (ws *Client) OnRaw(channel string, fn func(json.RawMessage)) {
	ws.subMu.Lock()
	ws.handlers[channel] = fn
	ws.subMu.Unlock()
}
