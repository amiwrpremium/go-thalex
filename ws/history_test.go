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
// TradeHistory
// ---------------------------------------------------------------------------

func TestTradeHistory_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/trade_history" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"trade_type":"trade","trade_id":"t1","order_id":"o1","instrument_name":"BTC-PERPETUAL","direction":"buy","price":50000,"amount":1,"time":1.0,"position_after":1,"fee":0.1,"fee_rate":0.0001,"fee_basis":50000,"leg_index":0}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trades, err := c.TradeHistory(ctx, &types.TradeHistoryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}
	if trades[0].TradeID != "t1" {
		t.Errorf("expected trade_id=t1, got %q", trades[0].TradeID)
	}
}

func TestTradeHistory_Empty(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`[]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trades, err := c.TradeHistory(ctx, &types.TradeHistoryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trades) != 0 {
		t.Errorf("expected 0 trades, got %d", len(trades))
	}
}

// ---------------------------------------------------------------------------
// OrderHistory (method)
// ---------------------------------------------------------------------------

func TestOrderHistory_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/order_history" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"order_id":"o1","order_type":"limit","direction":"buy","amount":1.0,"filled_amount":1.0,"status":"filled","fills":[],"insert_reason":"user","create_time":1.0,"close_time":2.0}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders, err := c.OrderHistory(ctx, &types.OrderHistoryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(orders))
	}
	if orders[0].OrderID != "o1" {
		t.Errorf("expected order_id=o1, got %q", orders[0].OrderID)
	}
}

// ---------------------------------------------------------------------------
// DailyMarkHistory
// ---------------------------------------------------------------------------

func TestDailyMarkHistory_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/daily_mark_history" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"time":1.0,"instrument_name":"BTC-PERPETUAL","position":1.0,"mark_price":50000,"realized_position_pnl":100}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	marks, err := c.DailyMarkHistory(ctx, &types.DailyMarkHistoryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(marks) != 1 {
		t.Fatalf("expected 1 mark, got %d", len(marks))
	}
	if marks[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %q", marks[0].InstrumentName)
	}
}

// ---------------------------------------------------------------------------
// TransactionHistory
// ---------------------------------------------------------------------------

func TestTransactionHistory_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/transaction_history" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"transaction_id":"tx1","time":1.0,"type":"trade","instrument_name":"BTC-PERPETUAL","amount":100,"currency":"USD"}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	txns, err := c.TransactionHistory(ctx, &types.TransactionHistoryParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txns) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txns))
	}
	if txns[0].TransactionID != "tx1" {
		t.Errorf("expected transaction_id=tx1, got %q", txns[0].TransactionID)
	}
}
