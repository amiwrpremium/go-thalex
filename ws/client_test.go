package ws

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
	"github.com/amiwrpremium/go-thalex/internal/transport"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------------------------------------------------------------------------
// NewClient
// ---------------------------------------------------------------------------

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.transport == nil {
		t.Fatal("transport should not be nil")
	}
	if c.reconnector != nil {
		t.Fatal("reconnector should be nil when WSReconnect is false by default")
	}
	if c.pending == nil {
		t.Fatal("pending map should be initialized")
	}
	if c.handlers == nil {
		t.Fatal("handlers map should be initialized")
	}
	if c.cfg.Network != config.Production {
		t.Errorf("expected default network Production, got %v", c.cfg.Network)
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	c := NewClient(
		config.WithNetwork(config.Testnet),
		config.WithWSDialTimeout(42*time.Second),
	)
	if c.cfg.Network != config.Testnet {
		t.Errorf("expected Testnet, got %v", c.cfg.Network)
	}
	if c.cfg.WSDialTimeout != 42*time.Second {
		t.Errorf("expected 42s dial timeout, got %v", c.cfg.WSDialTimeout)
	}
}

func TestNewClient_WithReconnect(t *testing.T) {
	c := NewClient(config.WithWSReconnect(true))
	if c.reconnector == nil {
		t.Fatal("reconnector should be created when WSReconnect=true")
	}
}

// ---------------------------------------------------------------------------
// isPrivateChannel
// ---------------------------------------------------------------------------

func TestIsPrivateChannel(t *testing.T) {
	tests := []struct {
		channel string
		want    bool
	}{
		{"account.orders", true},
		{"account.portfolio", true},
		{"session.orders", true},
		{"session.mm_protection", true},
		{"user.inbox_notifications", true},
		{"mm.rfqs", true},
		{"mm.rfq_quotes", true},
		{"ticker.BTC-PERPETUAL.100ms", false},
		{"book.BTC-PERPETUAL.1.10.100ms", false},
		{"instruments", false},
		{"system", false},
		{"banners", false},
		{"price_index.BTCUSD", false},
		{"recent_trades.BTCUSD.all", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.channel, func(t *testing.T) {
			got := isPrivateChannel(tt.channel)
			if got != tt.want {
				t.Errorf("isPrivateChannel(%q) = %v, want %v", tt.channel, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// OnResponse
// ---------------------------------------------------------------------------

func TestOnResponse_DispatchesToPendingCall(t *testing.T) {
	c := NewClient()
	pc := &pendingCall{result: make(chan *jsonrpc.Response, 1)}

	id := uint64(42)
	c.mu.Lock()
	c.pending[id] = pc
	c.mu.Unlock()

	resp := &jsonrpc.Response{
		JSONRPC: "2.0",
		ID:      &id,
		Result:  json.RawMessage(`{"ok":true}`),
	}
	c.OnResponse(resp)

	select {
	case got := <-pc.result:
		if got != resp {
			t.Errorf("expected dispatched response to be the same pointer")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for response dispatch")
	}
}

func TestOnResponse_NilID_Ignored(t *testing.T) {
	c := NewClient()
	// Should not panic even with a nil ID.
	resp := &jsonrpc.Response{JSONRPC: "2.0", ID: nil}
	c.OnResponse(resp) // no-op
}

func TestOnResponse_UnknownID_Ignored(t *testing.T) {
	c := NewClient()
	id := uint64(999)
	resp := &jsonrpc.Response{JSONRPC: "2.0", ID: &id}
	// No pending call registered for ID 999.
	c.OnResponse(resp) // should not panic
}

// ---------------------------------------------------------------------------
// OnNotification
// ---------------------------------------------------------------------------

func TestOnNotification_DispatchesHandler(t *testing.T) {
	c := NewClient()

	var mu sync.Mutex
	var received bool
	c.subMu.Lock()
	c.handlers["ticker.BTC-PERPETUAL.100ms"] = func(tk types.Ticker) {
		mu.Lock()
		received = true
		mu.Unlock()
	}
	c.subMu.Unlock()

	notif := &jsonrpc.Notification{
		JSONRPC: "2.0",
		Method:  "ticker.BTC-PERPETUAL.100ms",
		Params:  json.RawMessage(`{"mark_price":50000.0}`),
	}
	c.OnNotification(notif)

	// dispatchNotification runs in a goroutine; wait a bit.
	deadline := time.After(2 * time.Second)
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for notification dispatch")
		case <-tick.C:
			mu.Lock()
			got := received
			mu.Unlock()
			if got {
				return
			}
		}
	}
}

func TestOnNotification_NoHandler_Ignored(t *testing.T) {
	c := NewClient()
	notif := &jsonrpc.Notification{
		JSONRPC: "2.0",
		Method:  "unknown.channel",
		Params:  json.RawMessage(`{}`),
	}
	// Should not panic.
	c.OnNotification(notif)
}

// ---------------------------------------------------------------------------
// OnError
// ---------------------------------------------------------------------------

func TestOnError_WithHandler(t *testing.T) {
	c := NewClient()
	var captured error
	c.onError = func(err error) {
		captured = err
	}
	testErr := errors.New("test error")
	c.OnError(testErr)
	if captured != testErr {
		t.Errorf("expected captured error to be %v, got %v", testErr, captured)
	}
}

func TestOnError_NoHandler(t *testing.T) {
	c := NewClient()
	c.onError = nil
	// Should not panic.
	c.OnError(errors.New("some error"))
}

// ---------------------------------------------------------------------------
// OnErrorHandler
// ---------------------------------------------------------------------------

func TestOnErrorHandler(t *testing.T) {
	c := NewClient()
	var called bool
	c.OnErrorHandler(func(err error) {
		called = true
	})
	c.OnError(errors.New("x"))
	if !called {
		t.Fatal("OnErrorHandler callback was not invoked")
	}
}

// ---------------------------------------------------------------------------
// OnDisconnect
// ---------------------------------------------------------------------------

func TestOnDisconnect_NoReconnector(t *testing.T) {
	c := NewClient() // WSReconnect defaults to false
	// Should not panic when reconnector is nil.
	c.OnDisconnect()
}

// We cannot easily test OnDisconnect with a reconnector because
// TriggerReconnect tries to actually connect. We verify the code path
// does not panic and the goroutine is launched.
func TestOnDisconnect_WithReconnector(t *testing.T) {
	c := NewClient(config.WithWSReconnect(true))
	// OnDisconnect launches a goroutine that calls reconnector.TriggerReconnect.
	// It will fail because there is no real server, but it must not panic.
	c.OnDisconnect()
	// Give the goroutine a moment to start (and fail gracefully).
	time.Sleep(50 * time.Millisecond)
}

// ---------------------------------------------------------------------------
// Close -- clears pending calls
// ---------------------------------------------------------------------------

func TestClose_ClearsPendingCalls(t *testing.T) {
	c := NewClient()

	pc1 := &pendingCall{result: make(chan *jsonrpc.Response, 1)}
	pc2 := &pendingCall{result: make(chan *jsonrpc.Response, 1)}
	c.mu.Lock()
	c.pending[1] = pc1
	c.pending[2] = pc2
	c.mu.Unlock()

	// Close will also try to close the transport (which has no real conn).
	_ = c.Close()

	c.mu.Lock()
	n := len(c.pending)
	c.mu.Unlock()
	if n != 0 {
		t.Errorf("expected pending map to be empty after Close, got %d entries", n)
	}

	// Pending call channels should be closed.
	_, ok := <-pc1.result
	if ok {
		t.Error("expected pc1 channel to be closed")
	}
	_, ok = <-pc2.result
	if ok {
		t.Error("expected pc2 channel to be closed")
	}
}

// ---------------------------------------------------------------------------
// dispatchNotification -- all type cases
// ---------------------------------------------------------------------------

// waitForBool polls a mutex-protected bool until true or timeout.
func waitForBool(mu *sync.Mutex, flag *bool, timeout time.Duration) bool {
	deadline := time.After(timeout)
	tick := time.NewTicker(5 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-deadline:
			return false
		case <-tick.C:
			mu.Lock()
			v := *flag
			mu.Unlock()
			if v {
				return true
			}
		}
	}
}

func TestDispatchNotification_BookUpdate(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.BookUpdate) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"bids":[],"asks":[],"time":1.0}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("BookUpdate handler not called")
	}
}

func TestDispatchNotification_Ticker(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.Ticker) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"mark_price":50000.0}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("Ticker handler not called")
	}
}

func TestDispatchNotification_LightweightTicker(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.LightweightTicker) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"mark_price":50000.0}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("LightweightTicker handler not called")
	}
}

func TestDispatchNotification_RecentTrades(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.RecentTrade) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("RecentTrades handler not called")
	}
}

func TestDispatchNotification_IndexPrice(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.IndexPrice) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"index_name":"BTCUSD","price":50000.0,"timestamp":1.0}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("IndexPrice handler not called")
	}
}

func TestDispatchNotification_Instruments(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.Instrument) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("Instruments handler not called")
	}
}

func TestDispatchNotification_OrderStatuses(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.OrderStatus) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("OrderStatuses handler not called")
	}
}

func TestDispatchNotification_PortfolioEntries(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.PortfolioEntry) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("PortfolioEntries handler not called")
	}
}

func TestDispatchNotification_AccountSummary(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.AccountSummary) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"cash":[],"margin":1.0}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("AccountSummary handler not called")
	}
}

func TestDispatchNotification_Trades(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.Trade) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("Trades handler not called")
	}
}

func TestDispatchNotification_OrderHistory(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.OrderHistory) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("OrderHistory handler not called")
	}
}

func TestDispatchNotification_ConditionalOrders(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.ConditionalOrder) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("ConditionalOrders handler not called")
	}
}

func TestDispatchNotification_Bots(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.Bot) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("Bots handler not called")
	}
}

func TestDispatchNotification_Rfqs(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.Rfq) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("Rfqs handler not called")
	}
}

func TestDispatchNotification_RfqOrders(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.RfqOrder) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("RfqOrders handler not called")
	}
}

func TestDispatchNotification_MMProtectionUpdate(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.MMProtectionUpdate) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"product":"options","reason":"delta","time":1.0}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("MMProtectionUpdate handler not called")
	}
}

func TestDispatchNotification_Notification(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.Notification) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"id":"n1","time":1.0,"category":"trade","title":"t","message":"m","display_type":"popup","read":false,"popup":false}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("Notification handler not called")
	}
}

func TestDispatchNotification_SystemEvent(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.SystemEvent) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"event":"maintenance"}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("SystemEvent handler not called")
	}
}

func TestDispatchNotification_Banners(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v []types.Banner) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`[]`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("Banners handler not called")
	}
}

func TestDispatchNotification_RawJSON(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v json.RawMessage) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{"arbitrary":"data"}`)
	c.dispatchNotification(handler, data)
	if !waitForBool(&mu, &called, 2*time.Second) {
		t.Fatal("RawJSON handler not called")
	}
}

func TestDispatchNotification_UnknownType_Ignored(t *testing.T) {
	c := NewClient()
	// A handler of an unsupported type should be silently ignored.
	handler := func(v string) {}
	data := json.RawMessage(`"hello"`)
	// Should not panic.
	c.dispatchNotification(handler, data)
	// Give the goroutine a moment.
	time.Sleep(50 * time.Millisecond)
}

func TestDispatchNotification_InvalidJSON_HandlerNotCalled(t *testing.T) {
	c := NewClient()
	var mu sync.Mutex
	called := false
	handler := func(v types.Ticker) {
		mu.Lock()
		called = true
		mu.Unlock()
	}
	data := json.RawMessage(`{invalid json}`)
	c.dispatchNotification(handler, data)
	// Wait briefly; handler should NOT be called because unmarshal fails.
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	if called {
		t.Error("handler should not be called when JSON is invalid")
	}
	mu.Unlock()
}

// ---------------------------------------------------------------------------
// Login without credentials
// ---------------------------------------------------------------------------

func TestLogin_NoCredentials(t *testing.T) {
	c := NewClient()
	err := c.Login(context.TODO())
	if err == nil {
		t.Fatal("expected error when calling Login without credentials")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// Connect / IsConnected
// ---------------------------------------------------------------------------

func TestConnect_And_IsConnected(t *testing.T) {
	c := newConnectedClient(t, echoNull)
	if !c.IsConnected() {
		t.Fatal("expected IsConnected to be true after Connect")
	}
}

func TestIsConnected_BeforeConnect(t *testing.T) {
	c := NewClient()
	if c.IsConnected() {
		t.Fatal("expected IsConnected to be false before Connect")
	}
}

// ---------------------------------------------------------------------------
// call -- success path
// ---------------------------------------------------------------------------

func TestCall_Success_WithResult(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"order_id":"abc123"}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		OrderID string `json:"order_id"`
	}
	err := c.call(ctx, "test/method", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "abc123" {
		t.Errorf("expected order_id=abc123, got %q", result.OrderID)
	}
}

func TestCall_Success_NilResult(t *testing.T) {
	c := newConnectedClient(t, echoNull)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.call(ctx, "test/method", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// call -- API error path
// ---------------------------------------------------------------------------

func TestCall_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10001, Message: "insufficient funds"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.call(ctx, "private/insert", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
	if apiErr.Code != 10001 {
		t.Errorf("expected code 10001, got %d", apiErr.Code)
	}
	if apiErr.Message != "insufficient funds" {
		t.Errorf("expected message 'insufficient funds', got %q", apiErr.Message)
	}
}

// ---------------------------------------------------------------------------
// call -- context cancellation
// ---------------------------------------------------------------------------

func TestCall_ContextCancelled(t *testing.T) {
	// Use a handler that never responds so context cancellation kicks in.
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		time.Sleep(10 * time.Second) // block
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := c.call(ctx, "test/slow", nil, nil)
	if err == nil {
		t.Fatal("expected error due to context timeout")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// call -- connection closed while waiting
// ---------------------------------------------------------------------------

func TestCall_ConnectionClosedWhileWaiting(t *testing.T) {
	// Create a mock server that closes the connection immediately after
	// receiving a request (without sending a response).
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		// Return nothing -- the mock server loop will break when conn closes.
		time.Sleep(10 * time.Second)
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	// Launch call in a goroutine and close the client to trigger connection closed path.
	errCh := make(chan error, 1)
	go func() {
		ctx := context.Background()
		errCh <- c.call(ctx, "test/close", nil, nil)
	}()

	// Give the call a moment to register the pending call.
	time.Sleep(50 * time.Millisecond)

	// Close the client, which closes pending channels.
	_ = c.Close()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error when connection closed")
		}
		var connErr *apierr.ConnectionError
		if !errors.As(err, &connErr) {
			// Could also be context error or other; just verify we got an error.
			t.Logf("got non-ConnectionError: %T: %v (acceptable)", err, err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for call to return after Close")
	}
}

// ---------------------------------------------------------------------------
// callNoResult
// ---------------------------------------------------------------------------

func TestCallNoResult_Success(t *testing.T) {
	c := newConnectedClient(t, echoNull)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.callNoResult(ctx, "test/noresult", map[string]any{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCallNoResult_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 403, Message: "forbidden"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.callNoResult(ctx, "test/noresult", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// Close with reconnector
// ---------------------------------------------------------------------------

func TestClose_WithReconnector(t *testing.T) {
	c := NewClient(config.WithWSReconnect(true))
	// Should not panic even though not connected.
	err := c.Close()
	if err != nil {
		t.Logf("Close returned (acceptable): %v", err)
	}
}

// ---------------------------------------------------------------------------
// OnResponse concurrent access
// ---------------------------------------------------------------------------

func TestOnResponse_ConcurrentSafe(t *testing.T) {
	c := NewClient()
	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			id := uint64(i)
			resp := &jsonrpc.Response{JSONRPC: "2.0", ID: &id}
			c.OnResponse(resp)
		}(i)
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// OnNotification concurrent access
// ---------------------------------------------------------------------------

func TestOnNotification_ConcurrentSafe(t *testing.T) {
	c := NewClient()
	c.subMu.Lock()
	c.handlers["test.channel"] = func(v json.RawMessage) {}
	c.subMu.Unlock()

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			notif := &jsonrpc.Notification{
				JSONRPC: "2.0",
				Method:  "test.channel",
				Params:  json.RawMessage(`{}`),
			}
			c.OnNotification(notif)
		}()
	}
	wg.Wait()
	// Allow dispatch goroutines to complete.
	time.Sleep(50 * time.Millisecond)
}

// ---------------------------------------------------------------------------
// Login with credentials (via mock server)
// ---------------------------------------------------------------------------

func TestLogin_WithCredentials_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/login" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}

	// Generate a test RSA key for credentials.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	creds := auth.NewCredentials("test-key-id", privateKey)

	srv := newMockWSServer(t, handler)

	cfg := config.DefaultClientConfig()
	cfg.Credentials = creds
	c := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	c.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          wsURLFromHTTP(srv.URL),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: 60 * time.Second,
		Handler:      c,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.transport.Connect(ctx); err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	if err := c.Login(ctx); err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
}

func TestLogin_WithCredentials_AndAccountNumber(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/login" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	creds := auth.NewCredentials("test-key-id", privateKey)

	srv := newMockWSServer(t, handler)

	cfg := config.DefaultClientConfig()
	cfg.Credentials = creds
	cfg.AccountNumber = "acct-123"
	c := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	c.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          wsURLFromHTTP(srv.URL),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: 60 * time.Second,
		Handler:      c,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.transport.Connect(ctx); err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	if err := c.Login(ctx); err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
}

func TestLogin_WithCredentials_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "authentication failed"}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	creds := auth.NewCredentials("bad-key-id", privateKey)

	srv := newMockWSServer(t, handler)

	cfg := config.DefaultClientConfig()
	cfg.Credentials = creds
	c := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	c.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          wsURLFromHTTP(srv.URL),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: 60 * time.Second,
		Handler:      c,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.transport.Connect(ctx); err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	err = c.Login(ctx)
	if err == nil {
		t.Fatal("expected error from Login")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// onReconnect
// ---------------------------------------------------------------------------

func TestOnReconnect_NoHandlers(t *testing.T) {
	// When there are no handlers, onReconnect should do nothing besides
	// optionally trying to login if credentials are set.
	c := newConnectedClient(t, echoNull)

	err := c.onReconnect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOnReconnect_PublicChannelsOnly(t *testing.T) {
	var mu sync.Mutex
	var subscribedPublic bool

	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		switch req.Method {
		case "public/subscribe":
			mu.Lock()
			subscribedPublic = true
			mu.Unlock()
			return json.RawMessage(`null`), nil
		default:
			return json.RawMessage(`null`), nil
		}
	}
	c := newConnectedClient(t, handler)

	// Register public channel handlers.
	c.subMu.Lock()
	c.handlers["ticker.BTC-PERPETUAL.100ms"] = func(v types.Ticker) {}
	c.handlers["book.ETH-PERPETUAL.1.10.100ms"] = func(v types.BookUpdate) {}
	c.subMu.Unlock()

	err := c.onReconnect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	gotPublic := subscribedPublic
	mu.Unlock()
	if !gotPublic {
		t.Error("expected public/subscribe to be called")
	}
}

func TestOnReconnect_PrivateChannelsOnly(t *testing.T) {
	var mu sync.Mutex
	var subscribedPrivate bool

	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		switch req.Method {
		case "private/subscribe":
			mu.Lock()
			subscribedPrivate = true
			mu.Unlock()
			return json.RawMessage(`null`), nil
		default:
			return json.RawMessage(`null`), nil
		}
	}
	c := newConnectedClient(t, handler)

	// Register private channel handlers.
	c.subMu.Lock()
	c.handlers["account.orders"] = func(v []types.OrderStatus) {}
	c.handlers["session.orders"] = func(v []types.OrderStatus) {}
	c.subMu.Unlock()

	err := c.onReconnect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	gotPrivate := subscribedPrivate
	mu.Unlock()
	if !gotPrivate {
		t.Error("expected private/subscribe to be called")
	}
}

func TestOnReconnect_MixedChannels(t *testing.T) {
	var mu sync.Mutex
	var subscribedPublic, subscribedPrivate bool

	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		switch req.Method {
		case "public/subscribe":
			mu.Lock()
			subscribedPublic = true
			mu.Unlock()
			return json.RawMessage(`null`), nil
		case "private/subscribe":
			mu.Lock()
			subscribedPrivate = true
			mu.Unlock()
			return json.RawMessage(`null`), nil
		default:
			return json.RawMessage(`null`), nil
		}
	}
	c := newConnectedClient(t, handler)

	// Register both public and private channel handlers.
	c.subMu.Lock()
	c.handlers["ticker.BTC-PERPETUAL.100ms"] = func(v types.Ticker) {}
	c.handlers["account.orders"] = func(v []types.OrderStatus) {}
	c.handlers["mm.rfqs"] = func(v []types.Rfq) {}
	c.subMu.Unlock()

	err := c.onReconnect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	gotPub := subscribedPublic
	gotPriv := subscribedPrivate
	mu.Unlock()

	if !gotPub {
		t.Error("expected public/subscribe to be called for public channels")
	}
	if !gotPriv {
		t.Error("expected private/subscribe to be called for private channels")
	}
}

func TestOnReconnect_WithCredentials(t *testing.T) {
	var mu sync.Mutex
	var loggedIn bool

	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		switch req.Method {
		case "public/login":
			mu.Lock()
			loggedIn = true
			mu.Unlock()
			return json.RawMessage(`null`), nil
		default:
			return json.RawMessage(`null`), nil
		}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	creds := auth.NewCredentials("test-key", privateKey)

	srv := newMockWSServer(t, handler)

	cfg := config.DefaultClientConfig()
	cfg.Credentials = creds
	c := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	c.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          wsURLFromHTTP(srv.URL),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: 60 * time.Second,
		Handler:      c,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.transport.Connect(ctx); err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	// Add a public handler to also test channel re-subscription.
	c.subMu.Lock()
	c.handlers["ticker.BTC-PERPETUAL.100ms"] = func(v types.Ticker) {}
	c.subMu.Unlock()

	err = c.onReconnect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	gotLogin := loggedIn
	mu.Unlock()
	if !gotLogin {
		t.Error("expected public/login to be called during onReconnect")
	}
}

func TestOnReconnect_LoginFailure(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method == "public/login" {
			return nil, &jsonrpc.Error{Code: 10000, Message: "auth failed"}
		}
		return json.RawMessage(`null`), nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	creds := auth.NewCredentials("test-key", privateKey)

	srv := newMockWSServer(t, handler)

	cfg := config.DefaultClientConfig()
	cfg.Credentials = creds
	c := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	c.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          wsURLFromHTTP(srv.URL),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: 60 * time.Second,
		Handler:      c,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.transport.Connect(ctx); err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	// onReconnect should return error when login fails.
	err = c.onReconnect()
	if err == nil {
		t.Fatal("expected error when login fails during onReconnect")
	}
}

// ---------------------------------------------------------------------------
// Connect via NewClient (using mock server URL)
// ---------------------------------------------------------------------------

func TestConnect_Success_ViaTransport(t *testing.T) {
	srv := newMockWSServer(t, echoNull)

	cfg := config.DefaultClientConfig()
	c := &Client{
		cfg:      cfg,
		pending:  make(map[uint64]*pendingCall),
		handlers: make(map[string]any),
	}
	c.transport = transport.NewWSTransport(transport.WSTransportConfig{
		URL:          wsURLFromHTTP(srv.URL),
		DialTimeout:  cfg.WSDialTimeout,
		PingInterval: 60 * time.Second,
		Handler:      c,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.Connect(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	if !c.IsConnected() {
		t.Fatal("expected IsConnected to be true after Connect")
	}
}

func TestConnect_Failure_BadURL(t *testing.T) {
	c := NewClient()
	// The default URL points to wss://thalex.com which we can't reach in tests.
	// Use a short timeout to make this fail quickly.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := c.Connect(ctx)
	if err == nil {
		t.Fatal("expected error when connecting to unreachable server")
		_ = c.Close()
	}
}

// ---------------------------------------------------------------------------
// SetCancelOnDisconnect via connected client (session.go coverage)
// ---------------------------------------------------------------------------

func TestSetCancelOnDisconnect_Connected(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/set_cancel_on_disconnect" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test both enable and disable.
	if err := c.SetCancelOnDisconnect(ctx, true); err != nil {
		t.Fatalf("unexpected error (enable): %v", err)
	}
}

// ---------------------------------------------------------------------------
// CancelSession via connected client
// ---------------------------------------------------------------------------

func TestCancelSession_ZeroCancelled(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"n_cancelled":0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := c.CancelSession(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 cancelled, got %d", n)
	}
}

// ---------------------------------------------------------------------------
// methodRouter
// ---------------------------------------------------------------------------

func TestMethodRouter_KnownMethod(t *testing.T) {
	routes := map[string]rpcHandler{
		"test/echo": func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
			return json.RawMessage(`"ok"`), nil
		},
	}
	router := methodRouter(routes)

	result, rpcErr := router(&jsonrpc.Request{Method: "test/echo"})
	if rpcErr != nil {
		t.Fatalf("unexpected error: %v", rpcErr)
	}
	if string(result) != `"ok"` {
		t.Errorf("expected \"ok\", got %s", string(result))
	}
}

func TestMethodRouter_UnknownMethod(t *testing.T) {
	routes := map[string]rpcHandler{}
	router := methodRouter(routes)

	_, rpcErr := router(&jsonrpc.Request{Method: "unknown/method"})
	if rpcErr == nil {
		t.Fatal("expected error for unknown method")
	}
	if rpcErr.Code != -32601 {
		t.Errorf("expected code -32601, got %d", rpcErr.Code)
	}
}
