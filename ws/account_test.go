package ws

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// ---------------------------------------------------------------------------
// Portfolio
// ---------------------------------------------------------------------------

func TestPortfolio_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/portfolio" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"instrument_name":"BTC-PERPETUAL","position":1.0,"mark_price":50000,"start_price":49000,"average_price":49000,"unrealised_pnl":1000,"realised_pnl":0,"entry_value":49000}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	portfolio, err := c.Portfolio(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(portfolio) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(portfolio))
	}
	if portfolio[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %q", portfolio[0].InstrumentName)
	}
	if portfolio[0].Position != 1.0 {
		t.Errorf("expected position=1.0, got %f", portfolio[0].Position)
	}
}

func TestPortfolio_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "not authenticated"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Portfolio(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// AccountSummary (method)
// ---------------------------------------------------------------------------

func TestAccountSummary_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/account_summary" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"cash":[],"unrealised_pnl":100,"cash_collateral":10000,"margin":5000,"required_margin":2000,"remaining_margin":3000,"session_realised_pnl":50}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	summary, err := c.AccountSummary(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Margin != 5000 {
		t.Errorf("expected margin=5000, got %f", summary.Margin)
	}
	if summary.RequiredMargin != 2000 {
		t.Errorf("expected required_margin=2000, got %f", summary.RequiredMargin)
	}
}

// ---------------------------------------------------------------------------
// AccountBreakdown
// ---------------------------------------------------------------------------

func TestAccountBreakdown_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/account_breakdown" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"cash":[],"portfolio":[],"unrealised_pnl":0,"cash_collateral":0,"margin":0,"required_margin":0,"remaining_margin":0,"session_realised_pnl":0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	breakdown, err := c.AccountBreakdown(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if breakdown.Cash == nil {
		t.Error("expected Cash to be initialized (even if empty)")
	}
}

// ---------------------------------------------------------------------------
// RequiredMarginBreakdown
// ---------------------------------------------------------------------------

func TestRequiredMarginBreakdown_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/required_margin_breakdown" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"portfolio":{"required_margin":1000,"underlyings":[]}}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.RequiredMarginBreakdown(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Portfolio.RequiredMargin != 1000 {
		t.Errorf("expected required_margin=1000, got %f", result.Portfolio.RequiredMargin)
	}
}

// ---------------------------------------------------------------------------
// RequiredMarginForOrder
// ---------------------------------------------------------------------------

func TestRequiredMarginForOrder_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/required_margin_for_order" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"current":{"required_margin":100,"underlying":{"underlying":"BTCUSD","required_margin":100,"loss_margin":50,"roll_contingency_margin":0,"d1_roll_contingency_margin":0,"options_roll_contingency_margin":0,"options_contingency_margin":0}},"with_buy":{"required_margin":200,"underlying":{"underlying":"BTCUSD","required_margin":200,"loss_margin":100,"roll_contingency_margin":0,"d1_roll_contingency_margin":0,"options_roll_contingency_margin":0,"options_contingency_margin":0}},"with_sell":{"required_margin":150,"underlying":{"underlying":"BTCUSD","required_margin":150,"loss_margin":75,"roll_contingency_margin":0,"d1_roll_contingency_margin":0,"options_roll_contingency_margin":0,"options_contingency_margin":0}}}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.RequiredMarginForOrder(ctx, "BTC-PERPETUAL", 50000, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Current.RequiredMargin != 100 {
		t.Errorf("expected current required_margin=100, got %f", result.Current.RequiredMargin)
	}
	if result.WithBuy.RequiredMargin != 200 {
		t.Errorf("expected with_buy required_margin=200, got %f", result.WithBuy.RequiredMargin)
	}
}
