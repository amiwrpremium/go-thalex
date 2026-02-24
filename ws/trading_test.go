package ws

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------------------------------------------------------------------------
// Insert
// ---------------------------------------------------------------------------

func TestInsert_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/insert" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"order_id":"ord-1","direction":"buy","amount":1.0,"filled_amount":0,"remaining_amount":1.0,"status":"open","fills":[],"change_reason":"new","insert_reason":"user","create_time":1.0,"persistent":false}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	result, err := c.Insert(ctx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-1" {
		t.Errorf("expected order_id=ord-1, got %q", result.OrderID)
	}
}

func TestInsert_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10001, Message: "insufficient margin"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := types.NewBuyOrderParams("BTC-PERPETUAL", 100.0)
	_, err := c.Insert(ctx, params)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// Buy
// ---------------------------------------------------------------------------

func TestBuy_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/buy" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"order_id":"buy-1","direction":"buy","amount":2.0,"filled_amount":2.0,"remaining_amount":0,"status":"filled","fills":[],"change_reason":"filled","insert_reason":"user","create_time":1.0,"persistent":false}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.Buy(ctx, "BTC-PERPETUAL", 2.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "buy-1" {
		t.Errorf("expected order_id=buy-1, got %q", result.OrderID)
	}
	if result.Direction != enums.DirectionBuy {
		t.Errorf("expected direction=buy, got %q", result.Direction)
	}
}

// ---------------------------------------------------------------------------
// Sell
// ---------------------------------------------------------------------------

func TestSell_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/sell" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"order_id":"sell-1","direction":"sell","amount":3.0,"filled_amount":3.0,"remaining_amount":0,"status":"filled","fills":[],"change_reason":"filled","insert_reason":"user","create_time":1.0,"persistent":false}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.Sell(ctx, "ETH-PERPETUAL", 3.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "sell-1" {
		t.Errorf("expected order_id=sell-1, got %q", result.OrderID)
	}
}

// ---------------------------------------------------------------------------
// Amend
// ---------------------------------------------------------------------------

func TestAmend_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/amend" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"order_id":"ord-1","direction":"buy","amount":2.0,"filled_amount":0,"remaining_amount":2.0,"status":"open","fills":[],"change_reason":"amended","insert_reason":"user","create_time":1.0,"persistent":false}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := types.NewAmendByOrderID("ord-1", 50000.0, 2.0)
	result, err := c.Amend(ctx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-1" {
		t.Errorf("expected order_id=ord-1, got %q", result.OrderID)
	}
}

// ---------------------------------------------------------------------------
// Cancel
// ---------------------------------------------------------------------------

func TestCancel_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/cancel" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"order_id":"ord-1","direction":"buy","amount":1.0,"filled_amount":0,"remaining_amount":0,"status":"cancelled","fills":[],"change_reason":"cancelled","delete_reason":"user","insert_reason":"user","create_time":1.0,"persistent":false}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := types.CancelByOrderID("ord-1")
	result, err := c.Cancel(ctx, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-1" {
		t.Errorf("expected order_id=ord-1, got %q", result.OrderID)
	}
}

// ---------------------------------------------------------------------------
// CancelAll
// ---------------------------------------------------------------------------

func TestCancelAll_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/cancel_all" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"n_cancelled":3}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := c.CancelAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 cancelled, got %d", n)
	}
}

func TestCancelAll_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "not authenticated"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := c.CancelAll(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if n != 0 {
		t.Errorf("expected 0 on error, got %d", n)
	}
}

// ---------------------------------------------------------------------------
// OpenOrders
// ---------------------------------------------------------------------------

func TestOpenOrders_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/open_orders" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"orders":[{"order_id":"ord-1","direction":"buy","amount":1.0,"filled_amount":0,"remaining_amount":1.0,"status":"open","fills":[],"change_reason":"new","insert_reason":"user","create_time":1.0,"persistent":false}]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders, err := c.OpenOrders(ctx, "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(orders))
	}
	if orders[0].OrderID != "ord-1" {
		t.Errorf("expected order_id=ord-1, got %q", orders[0].OrderID)
	}
}

func TestOpenOrders_EmptyInstrument(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"orders":[]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders, err := c.OpenOrders(ctx, "") // empty instrument name
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("expected 0 orders, got %d", len(orders))
	}
}
