package rest

import (
	"context"
	"net/http"
	"testing"

	"github.com/amiwrpremium/go-thalex/types"
)

func TestPortfolio_Success(t *testing.T) {
	portfolio := []types.PortfolioEntry{
		{
			InstrumentName: "BTC-PERPETUAL",
			Position:       1.5,
			MarkPrice:      50000.0,
			StartPrice:     48000.0,
			AveragePrice:   49000.0,
			UnrealisedPnl:  1500.0,
			RealisedPnl:    500.0,
			EntryValue:     73500.0,
		},
		{
			InstrumentName: "ETH-PERPETUAL",
			Position:       -10.0,
			MarkPrice:      3000.0,
			StartPrice:     3100.0,
			AveragePrice:   3050.0,
			UnrealisedPnl:  500.0,
			RealisedPnl:    0,
			EntryValue:     30500.0,
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/private/portfolio" {
			t.Errorf("expected path /private/portfolio, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, portfolio))
	})

	result, err := c.Portfolio(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %s", result[0].InstrumentName)
	}
	if result[0].Position != 1.5 {
		t.Errorf("expected position=1.5, got %f", result[0].Position)
	}
	if !result[0].IsLong() {
		t.Error("expected first entry to be long")
	}
	if !result[1].IsShort() {
		t.Error("expected second entry to be short")
	}
}

func TestPortfolio_Empty(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, []types.PortfolioEntry{}))
	})

	result, err := c.Portfolio(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestPortfolio_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(apiErrorJSON(10001, "unauthorized"))
	})

	_, err := c.Portfolio(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountSummary_Success(t *testing.T) {
	summary := types.AccountSummary{
		Cash: []types.CashHolding{
			{Currency: "BTC", Balance: 1.5, CollateralFactor: 1.0, Transactable: true},
			{Currency: "ETH", Balance: 10.0, CollateralFactor: 0.95, Transactable: true},
		},
		UnrealisedPnl:      2000.0,
		CashCollateral:     75000.0,
		Margin:             100000.0,
		RequiredMargin:     20000.0,
		RemainingMargin:    80000.0,
		SessionRealisedPnl: 500.0,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/account_summary" {
			t.Errorf("expected path /private/account_summary, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, summary))
	})

	result, err := c.AccountSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Cash) != 2 {
		t.Fatalf("expected 2 cash holdings, got %d", len(result.Cash))
	}
	if result.Cash[0].Currency != "BTC" {
		t.Errorf("expected BTC, got %s", result.Cash[0].Currency)
	}
	if result.Margin != 100000.0 {
		t.Errorf("expected margin=100000, got %f", result.Margin)
	}
	if result.RequiredMargin != 20000.0 {
		t.Errorf("expected required_margin=20000, got %f", result.RequiredMargin)
	}

	// Test MarginUtilization
	util := result.MarginUtilization()
	if util != 0.2 {
		t.Errorf("expected margin utilization=0.2, got %f", util)
	}
}

func TestAccountSummary_ZeroMargin(t *testing.T) {
	summary := types.AccountSummary{
		Cash:   []types.CashHolding{},
		Margin: 0,
	}
	util := summary.MarginUtilization()
	if util != 0 {
		t.Errorf("expected utilization=0 when margin is 0, got %f", util)
	}
}

func TestAccountBreakdown_Success(t *testing.T) {
	breakdown := types.AccountBreakdown{
		Cash: []types.CashHolding{
			{Currency: "BTC", Balance: 1.0},
		},
		Portfolio: []types.PortfolioEntry{
			{InstrumentName: "BTC-PERPETUAL", Position: 0.5},
		},
		Margin:         50000.0,
		RequiredMargin: 10000.0,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/account_breakdown" {
			t.Errorf("expected path /private/account_breakdown, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, breakdown))
	})

	result, err := c.AccountBreakdown(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Cash) != 1 {
		t.Errorf("expected 1 cash holding, got %d", len(result.Cash))
	}
	if len(result.Portfolio) != 1 {
		t.Errorf("expected 1 portfolio entry, got %d", len(result.Portfolio))
	}
}

func TestRequiredMarginBreakdown_Success(t *testing.T) {
	breakdown := types.PortfolioMarginBreakdown{}
	breakdown.Portfolio.RequiredMargin = 15000.0
	breakdown.Portfolio.Underlyings = []types.UnderlyingMarginDetail{
		{Underlying: "BTCUSD", RequiredMargin: 15000.0},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/required_margin_breakdown" {
			t.Errorf("expected path /private/required_margin_breakdown, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, breakdown))
	})

	result, err := c.RequiredMarginBreakdown(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Portfolio.RequiredMargin != 15000.0 {
		t.Errorf("expected required_margin=15000, got %f", result.Portfolio.RequiredMargin)
	}
}

func TestRequiredMarginForOrder_Success(t *testing.T) {
	marginResult := types.MarginForOrderResult{
		Current:  types.MarginBreakdownSide{RequiredMargin: 10000},
		WithBuy:  types.MarginBreakdownSide{RequiredMargin: 12000},
		WithSell: types.MarginBreakdownSide{RequiredMargin: 8000},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/required_margin_for_order" {
			t.Errorf("expected path /private/required_margin_for_order, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("instrument_name") != "BTC-PERPETUAL" {
			t.Errorf("expected instrument_name=BTC-PERPETUAL")
		}
		if r.URL.Query().Get("price") != "50000" {
			t.Errorf("expected price=50000, got %s", r.URL.Query().Get("price"))
		}
		if r.URL.Query().Get("amount") != "1" {
			t.Errorf("expected amount=1, got %s", r.URL.Query().Get("amount"))
		}
		w.Write(wrapResult(t, marginResult))
	})

	result, err := c.RequiredMarginForOrder(context.Background(), "BTC-PERPETUAL", 50000, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Current.RequiredMargin != 10000 {
		t.Errorf("expected current required_margin=10000, got %f", result.Current.RequiredMargin)
	}
	if result.WithBuy.RequiredMargin != 12000 {
		t.Errorf("expected with_buy required_margin=12000, got %f", result.WithBuy.RequiredMargin)
	}
}

func TestAccountSummary_WithAccountNumber(t *testing.T) {
	summary := types.AccountSummary{
		Cash:   []types.CashHolding{{Currency: "BTC", Balance: 1.0}},
		Margin: 50000.0,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Thalex-Account") != "ACC-001" {
			t.Errorf("expected X-Thalex-Account=ACC-001, got %q", r.Header.Get("X-Thalex-Account"))
		}
		w.Write(wrapResult(t, summary))
	})

	result, err := c.AccountSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Margin != 50000.0 {
		t.Errorf("expected margin=50000, got %f", result.Margin)
	}
}
