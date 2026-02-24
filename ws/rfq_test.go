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
// CreateRfq
// ---------------------------------------------------------------------------

func TestCreateRfq_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/create_rfq" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"rfq_id":"rfq-1","status":"open","create_time":1.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CreateRfq(ctx, &types.CreateRfqParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RfqID != "rfq-1" {
		t.Errorf("expected rfq_id=rfq-1, got %q", result.RfqID)
	}
}

// ---------------------------------------------------------------------------
// CancelRfq
// ---------------------------------------------------------------------------

func TestCancelRfq_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/cancel_rfq" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.CancelRfq(ctx, "rfq-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// TradeRfq
// ---------------------------------------------------------------------------

func TestTradeRfq_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/trade_rfq" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"trade_type":"rfq","trade_id":"t1","order_id":"o1","instrument_name":"BTC-PERPETUAL","direction":"buy","price":50000,"amount":1,"time":1.0,"position_after":1,"fee":0.1,"fee_rate":0.0001,"fee_basis":50000,"leg_index":0}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	trades, err := c.TradeRfq(ctx, &types.TradeRfqParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}
}

// ---------------------------------------------------------------------------
// OpenRfqs
// ---------------------------------------------------------------------------

func TestOpenRfqs_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/open_rfqs" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"rfq_id":"rfq-1","status":"open","create_time":1.0}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rfqs, err := c.OpenRfqs(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rfqs) != 1 {
		t.Fatalf("expected 1 rfq, got %d", len(rfqs))
	}
}

// ---------------------------------------------------------------------------
// RfqHistory
// ---------------------------------------------------------------------------

func TestRfqHistory_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/rfq_history" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"rfq_id":"rfq-1","status":"filled","create_time":1.0}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rfqs, err := c.RfqHistory(ctx, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rfqs) != 1 {
		t.Fatalf("expected 1 rfq, got %d", len(rfqs))
	}
}

func TestRfqHistory_WithAllParams(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/rfq_history" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	from := 1.0
	to := 2.0
	offset := 0
	limit := 10
	rfqs, err := c.RfqHistory(ctx, &from, &to, &offset, &limit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rfqs == nil {
		t.Error("expected non-nil result")
	}
}

// ---------------------------------------------------------------------------
// MMRfqs
// ---------------------------------------------------------------------------

func TestMMRfqs_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/mm_rfqs" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rfqs, err := c.MMRfqs(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rfqs) != 0 {
		t.Errorf("expected 0 rfqs, got %d", len(rfqs))
	}
}

// ---------------------------------------------------------------------------
// MMRfqInsertQuote
// ---------------------------------------------------------------------------

func TestMMRfqInsertQuote_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/mm_rfq_insert_quote" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"order_id":"q1","rfq_id":"rfq-1"}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.MMRfqInsertQuote(ctx, &types.RfqQuoteInsertParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "q1" {
		t.Errorf("expected order_id=q1, got %q", result.OrderID)
	}
}

// ---------------------------------------------------------------------------
// MMRfqAmendQuote
// ---------------------------------------------------------------------------

func TestMMRfqAmendQuote_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/mm_rfq_amend_quote" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"order_id":"q1","rfq_id":"rfq-1"}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.MMRfqAmendQuote(ctx, &types.RfqQuoteAmendParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "q1" {
		t.Errorf("expected order_id=q1, got %q", result.OrderID)
	}
}

// ---------------------------------------------------------------------------
// MMRfqDeleteQuote
// ---------------------------------------------------------------------------

func TestMMRfqDeleteQuote_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/mm_rfq_delete_quote" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.MMRfqDeleteQuote(ctx, &types.RfqQuoteDeleteParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// MMRfqQuotes
// ---------------------------------------------------------------------------

func TestMMRfqQuotes_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/mm_rfq_quotes" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"order_id":"q1","rfq_id":"rfq-1"}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	quotes, err := c.MMRfqQuotes(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}
}
