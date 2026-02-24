package ws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------------------------------------------------------------------------
// Typed handler registration (On* methods)
// ---------------------------------------------------------------------------

func TestOnBook(t *testing.T) {
	c := NewClient()
	ch := "book.BTC-PERPETUAL.1.10.100ms"
	c.OnBook(ch, func(v types.BookUpdate) {})

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", ch)
	}
}

func TestOnTicker(t *testing.T) {
	c := NewClient()
	ch := "ticker.BTC-PERPETUAL.100ms"
	c.OnTicker(ch, func(v types.Ticker) {})

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", ch)
	}
}

func TestOnLWT(t *testing.T) {
	c := NewClient()
	ch := "lwt.BTC-PERPETUAL.100ms"
	c.OnLWT(ch, func(v types.LightweightTicker) {})

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", ch)
	}
}

func TestOnRecentTrades(t *testing.T) {
	c := NewClient()
	ch := "recent_trades.BTCUSD.all"
	c.OnRecentTrades(ch, func(v []types.RecentTrade) {})

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", ch)
	}
}

func TestOnPriceIndex(t *testing.T) {
	c := NewClient()
	ch := "price_index.BTCUSD"
	c.OnPriceIndex(ch, func(v types.IndexPrice) {})

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", ch)
	}
}

func TestOnInstruments(t *testing.T) {
	c := NewClient()
	c.OnInstruments(func(v []types.Instrument) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelInstruments]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelInstruments)
	}
}

func TestOnOrders(t *testing.T) {
	c := NewClient()
	c.OnOrders(func(v []types.OrderStatus) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountOrders]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountOrders)
	}
}

func TestOnPersistentOrders(t *testing.T) {
	c := NewClient()
	c.OnPersistentOrders(func(v []types.OrderStatus) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountPersistent]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountPersistent)
	}
}

func TestOnSessionOrders(t *testing.T) {
	c := NewClient()
	c.OnSessionOrders(func(v []types.OrderStatus) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelSessionOrders]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelSessionOrders)
	}
}

func TestOnPortfolio(t *testing.T) {
	c := NewClient()
	c.OnPortfolio(func(v []types.PortfolioEntry) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountPortfolio]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountPortfolio)
	}
}

func TestOnAccountSummary(t *testing.T) {
	c := NewClient()
	c.OnAccountSummary(func(v types.AccountSummary) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountSummary]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountSummary)
	}
}

func TestOnTradeHistory(t *testing.T) {
	c := NewClient()
	c.OnTradeHistory(func(v []types.Trade) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountTradeHistory]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountTradeHistory)
	}
}

func TestOnOrderHistory(t *testing.T) {
	c := NewClient()
	c.OnOrderHistory(func(v []types.OrderHistory) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountOrderHistory]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountOrderHistory)
	}
}

func TestOnConditionalOrders(t *testing.T) {
	c := NewClient()
	c.OnConditionalOrders(func(v []types.ConditionalOrder) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountConditional]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountConditional)
	}
}

func TestOnBots(t *testing.T) {
	c := NewClient()
	c.OnBots(func(v []types.Bot) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountBots]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountBots)
	}
}

func TestOnRfqs(t *testing.T) {
	c := NewClient()
	c.OnRfqs(func(v []types.Rfq) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelAccountRfqs]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelAccountRfqs)
	}
}

func TestOnMMRfqs(t *testing.T) {
	c := NewClient()
	c.OnMMRfqs(func(v []types.Rfq) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelMMRfqs]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelMMRfqs)
	}
}

func TestOnMMRfqQuotes(t *testing.T) {
	c := NewClient()
	c.OnMMRfqQuotes(func(v []types.RfqOrder) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelMMRfqQuotes]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelMMRfqQuotes)
	}
}

func TestOnMMProtection(t *testing.T) {
	c := NewClient()
	c.OnMMProtection(func(v types.MMProtectionUpdate) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelSessionMMProtection]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelSessionMMProtection)
	}
}

func TestOnNotifications(t *testing.T) {
	c := NewClient()
	c.OnNotifications(func(v types.Notification) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelUserNotifications]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelUserNotifications)
	}
}

func TestOnSystemEvent(t *testing.T) {
	c := NewClient()
	c.OnSystemEvent(func(v types.SystemEvent) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelSystem]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelSystem)
	}
}

func TestOnBanners(t *testing.T) {
	c := NewClient()
	c.OnBanners(func(v []types.Banner) {})

	c.subMu.RLock()
	_, ok := c.handlers[types.ChannelBanners]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", types.ChannelBanners)
	}
}

func TestOnRaw(t *testing.T) {
	c := NewClient()
	ch := "custom.channel"
	c.OnRaw(ch, func(v json.RawMessage) {})

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if !ok {
		t.Fatalf("handler not registered for channel %q", ch)
	}
}

// ---------------------------------------------------------------------------
// Handler replacement
// ---------------------------------------------------------------------------

func TestOnHandler_Replaces_Previous(t *testing.T) {
	c := NewClient()
	ch := "ticker.ETH-PERPETUAL.100ms"

	first := func(v types.Ticker) {}
	second := func(v types.Ticker) {}

	c.OnTicker(ch, first)
	c.OnTicker(ch, second)

	c.subMu.RLock()
	count := 0
	for k := range c.handlers {
		if k == ch {
			count++
		}
	}
	c.subMu.RUnlock()

	if count != 1 {
		t.Errorf("expected exactly 1 handler entry for channel, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// Unsubscribe removes handlers
// ---------------------------------------------------------------------------

func TestUnsubscribe_RemovesHandlers(t *testing.T) {
	c := NewClient()
	ch1 := "ticker.BTC-PERPETUAL.100ms"
	ch2 := "ticker.ETH-PERPETUAL.100ms"

	c.OnTicker(ch1, func(v types.Ticker) {})
	c.OnTicker(ch2, func(v types.Ticker) {})

	// Unsubscribe calls callNoResult, which will fail because there is no
	// transport connection. But the handler removal happens before the RPC call.
	_ = c.Unsubscribe(context.TODO(), ch1)

	c.subMu.RLock()
	_, ok1 := c.handlers[ch1]
	_, ok2 := c.handlers[ch2]
	c.subMu.RUnlock()

	if ok1 {
		t.Errorf("handler for %q should be removed after Unsubscribe", ch1)
	}
	if !ok2 {
		t.Errorf("handler for %q should still be present", ch2)
	}
}

func TestUnsubscribePrivate_RemovesHandlers(t *testing.T) {
	c := NewClient()
	ch := types.ChannelAccountOrders

	c.OnOrders(func(v []types.OrderStatus) {})

	_ = c.UnsubscribePrivate(context.TODO(), ch)

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()

	if ok {
		t.Errorf("handler for %q should be removed after UnsubscribePrivate", ch)
	}
}

// ---------------------------------------------------------------------------
// Multiple handlers for different channels
// ---------------------------------------------------------------------------

func TestMultipleHandlers_DifferentChannels(t *testing.T) {
	c := NewClient()

	c.OnTicker("ticker.BTC-PERPETUAL.100ms", func(v types.Ticker) {})
	c.OnBook("book.BTC-PERPETUAL.1.10.100ms", func(v types.BookUpdate) {})
	c.OnOrders(func(v []types.OrderStatus) {})

	c.subMu.RLock()
	n := len(c.handlers)
	c.subMu.RUnlock()

	if n != 3 {
		t.Errorf("expected 3 handlers registered, got %d", n)
	}
}

// ---------------------------------------------------------------------------
// Subscribe / SubscribePrivate with connected client
// ---------------------------------------------------------------------------

func TestSubscribe_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/subscribe" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.Subscribe(ctx, "ticker.BTC-PERPETUAL.100ms", "book.BTC-PERPETUAL.1.10.100ms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSubscribePrivate_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/subscribe" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.SubscribePrivate(ctx, types.ChannelAccountOrders, types.ChannelAccountPortfolio)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnsubscribe_Connected_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/unsubscribe" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ch := "ticker.BTC-PERPETUAL.100ms"
	c.OnTicker(ch, func(v types.Ticker) {})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.Unsubscribe(ctx, ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if ok {
		t.Error("handler should be removed after Unsubscribe")
	}
}

func TestUnsubscribePrivate_Connected_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/unsubscribe" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ch := types.ChannelAccountOrders
	c.OnOrders(func(v []types.OrderStatus) {})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.UnsubscribePrivate(ctx, ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c.subMu.RLock()
	_, ok := c.handlers[ch]
	c.subMu.RUnlock()
	if ok {
		t.Error("handler should be removed after UnsubscribePrivate")
	}
}

// ---------------------------------------------------------------------------
// Subscribe / Unsubscribe multiple channels
// ---------------------------------------------------------------------------

func TestUnsubscribe_MultipleChannels(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ch1 := "ticker.BTC-PERPETUAL.100ms"
	ch2 := "ticker.ETH-PERPETUAL.100ms"
	ch3 := "book.BTC-PERPETUAL.1.10.100ms"
	c.OnTicker(ch1, func(v types.Ticker) {})
	c.OnTicker(ch2, func(v types.Ticker) {})
	c.OnBook(ch3, func(v types.BookUpdate) {})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Unsubscribe from two of three.
	err := c.Unsubscribe(ctx, ch1, ch2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c.subMu.RLock()
	_, ok1 := c.handlers[ch1]
	_, ok2 := c.handlers[ch2]
	_, ok3 := c.handlers[ch3]
	c.subMu.RUnlock()

	if ok1 {
		t.Errorf("handler for %q should be removed", ch1)
	}
	if ok2 {
		t.Errorf("handler for %q should be removed", ch2)
	}
	if !ok3 {
		t.Errorf("handler for %q should still be present", ch3)
	}
}
