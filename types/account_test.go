package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- AccountSummary.MarginUtilization ----------

func TestAccountSummary_MarginUtilization(t *testing.T) {
	tests := []struct {
		name           string
		margin         float64
		requiredMargin float64
		want           float64
	}{
		{
			name:           "zero_margin_returns_zero",
			margin:         0,
			requiredMargin: 100,
			want:           0,
		},
		{
			name:           "zero_margin_zero_required",
			margin:         0,
			requiredMargin: 0,
			want:           0,
		},
		{
			name:           "half_utilization",
			margin:         1000,
			requiredMargin: 500,
			want:           0.5,
		},
		{
			name:           "full_utilization",
			margin:         1000,
			requiredMargin: 1000,
			want:           1.0,
		},
		{
			name:           "over_utilization",
			margin:         1000,
			requiredMargin: 1500,
			want:           1.5,
		},
		{
			name:           "low_utilization",
			margin:         10000,
			requiredMargin: 100,
			want:           0.01,
		},
		{
			name:           "no_required_margin",
			margin:         5000,
			requiredMargin: 0,
			want:           0,
		},
		{
			name:           "fractional_values",
			margin:         3333.33,
			requiredMargin: 1111.11,
			want:           1111.11 / 3333.33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &types.AccountSummary{
				Margin:         tt.margin,
				RequiredMargin: tt.requiredMargin,
			}
			got := s.MarginUtilization()
			if got != tt.want {
				t.Errorf("MarginUtilization() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------- PortfolioEntry.IsLong ----------

func TestPortfolioEntry_IsLong(t *testing.T) {
	tests := []struct {
		name     string
		position float64
		want     bool
	}{
		{"positive_position", 10.0, true},
		{"zero_position", 0, false},
		{"negative_position", -5.0, false},
		{"small_positive", 0.0001, true},
		{"small_negative", -0.0001, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &types.PortfolioEntry{Position: tt.position}
			if got := p.IsLong(); got != tt.want {
				t.Errorf("PortfolioEntry{Position: %v}.IsLong() = %v, want %v", tt.position, got, tt.want)
			}
		})
	}
}

// ---------- PortfolioEntry.IsShort ----------

func TestPortfolioEntry_IsShort(t *testing.T) {
	tests := []struct {
		name     string
		position float64
		want     bool
	}{
		{"negative_position", -10.0, true},
		{"zero_position", 0, false},
		{"positive_position", 5.0, false},
		{"small_negative", -0.0001, true},
		{"small_positive", 0.0001, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &types.PortfolioEntry{Position: tt.position}
			if got := p.IsShort(); got != tt.want {
				t.Errorf("PortfolioEntry{Position: %v}.IsShort() = %v, want %v", tt.position, got, tt.want)
			}
		})
	}
}

// ---------- IsLong and IsShort are mutually exclusive for non-zero ----------

func TestPortfolioEntry_LongShortMutuallyExclusive(t *testing.T) {
	tests := []struct {
		name     string
		position float64
	}{
		{"positive", 1.0},
		{"negative", -1.0},
		{"zero", 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &types.PortfolioEntry{Position: tt.position}
			long := p.IsLong()
			short := p.IsShort()
			if long && short {
				t.Errorf("Position %v: IsLong and IsShort are both true", tt.position)
			}
			if tt.position == 0 && (long || short) {
				t.Errorf("Position 0: expected both IsLong and IsShort to be false")
			}
		})
	}
}

// ---------- AccountSummary JSON round-trip ----------

func TestAccountSummary_JSONRoundTrip(t *testing.T) {
	collIndexPrice := 50000.0
	s := types.AccountSummary{
		Cash: []types.CashHolding{
			{
				Currency:             "BTC",
				Balance:              1.5,
				CollateralFactor:     0.95,
				CollateralIndexPrice: &collIndexPrice,
				Transactable:         true,
			},
			{
				Currency:         "USD",
				Balance:          10000.0,
				CollateralFactor: 1.0,
				Transactable:     true,
			},
		},
		UnrealisedPnl:      500.0,
		CashCollateral:     60000.0,
		Margin:             50000.0,
		RequiredMargin:     25000.0,
		RemainingMargin:    25000.0,
		SessionRealisedPnl: 100.0,
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.AccountSummary
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Cash) != 2 {
		t.Fatalf("len(Cash) = %d, want 2", len(got.Cash))
	}
	if got.Cash[0].Currency != "BTC" {
		t.Errorf("Cash[0].Currency = %q, want %q", got.Cash[0].Currency, "BTC")
	}
	if got.Cash[0].CollateralIndexPrice == nil || *got.Cash[0].CollateralIndexPrice != collIndexPrice {
		t.Errorf("Cash[0].CollateralIndexPrice = %v, want %v", got.Cash[0].CollateralIndexPrice, collIndexPrice)
	}
	if got.Cash[0].Transactable != true {
		t.Error("Cash[0].Transactable = false, want true")
	}
	if got.Margin != 50000.0 {
		t.Errorf("Margin = %v, want 50000.0", got.Margin)
	}
	if got.RequiredMargin != 25000.0 {
		t.Errorf("RequiredMargin = %v, want 25000.0", got.RequiredMargin)
	}

	// Verify MarginUtilization on round-tripped value
	if got.MarginUtilization() != 0.5 {
		t.Errorf("MarginUtilization() = %v, want 0.5", got.MarginUtilization())
	}
}

// ---------- CashHolding JSON round-trip ----------

func TestCashHolding_JSONRoundTrip(t *testing.T) {
	collIdx := 3000.0
	ch := types.CashHolding{
		Currency:             "ETH",
		Balance:              10.0,
		CollateralFactor:     0.9,
		CollateralIndexPrice: &collIdx,
		Transactable:         false,
	}

	data, err := json.Marshal(ch)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.CashHolding
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Currency != "ETH" {
		t.Errorf("Currency = %q, want %q", got.Currency, "ETH")
	}
	if got.CollateralFactor != 0.9 {
		t.Errorf("CollateralFactor = %v, want 0.9", got.CollateralFactor)
	}
	if got.CollateralIndexPrice == nil || *got.CollateralIndexPrice != collIdx {
		t.Errorf("CollateralIndexPrice = %v, want %v", got.CollateralIndexPrice, collIdx)
	}
}

// ---------- PortfolioEntry JSON round-trip ----------

func TestPortfolioEntry_JSONRoundTrip(t *testing.T) {
	iv := 0.65
	idx := 50000.0
	perpFundEntry := 1.0
	unrealFunding := 0.5

	pe := types.PortfolioEntry{
		InstrumentName:             "BTC-PERPETUAL",
		Position:                   5.0,
		MarkPrice:                  50000.0,
		IV:                         &iv,
		Index:                      &idx,
		StartPrice:                 49000.0,
		AveragePrice:               49500.0,
		UnrealisedPnl:              2500.0,
		RealisedPnl:                100.0,
		EntryValue:                 247500.0,
		PerpetualFundingEntryValue: &perpFundEntry,
		UnrealisedPerpetualFunding: &unrealFunding,
	}

	data, err := json.Marshal(pe)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.PortfolioEntry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, "BTC-PERPETUAL")
	}
	if got.Position != 5.0 {
		t.Errorf("Position = %v, want 5.0", got.Position)
	}
	if got.IV == nil || *got.IV != iv {
		t.Errorf("IV = %v, want %v", got.IV, iv)
	}
	if got.PerpetualFundingEntryValue == nil || *got.PerpetualFundingEntryValue != perpFundEntry {
		t.Errorf("PerpetualFundingEntryValue = %v, want %v", got.PerpetualFundingEntryValue, perpFundEntry)
	}

	// Verify IsLong/IsShort on round-tripped value
	if !got.IsLong() {
		t.Error("IsLong() = false, want true for position 5.0")
	}
	if got.IsShort() {
		t.Error("IsShort() = true, want false for position 5.0")
	}
}

// ---------- MarginState JSON round-trip ----------

func TestMarginState_JSONRoundTrip(t *testing.T) {
	ms := types.MarginState{
		Underlying:                   "BTCUSD",
		RequiredMargin:               5000.0,
		LossMargin:                   4000.0,
		RollContingencyMargin:        100.0,
		D1RollContingencyMargin:      50.0,
		OptionsRollContingencyMargin: 25.0,
		OptionsContingencyMargin:     75.0,
	}

	data, err := json.Marshal(ms)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.MarginState
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Underlying != "BTCUSD" {
		t.Errorf("Underlying = %q, want %q", got.Underlying, "BTCUSD")
	}
	if got.RequiredMargin != 5000.0 {
		t.Errorf("RequiredMargin = %v, want 5000.0", got.RequiredMargin)
	}
}

// ---------- MarginForOrderResult JSON round-trip ----------

func TestMarginForOrderResult_JSONRoundTrip(t *testing.T) {
	r := types.MarginForOrderResult{
		Current: types.MarginBreakdownSide{
			RequiredMargin: 5000.0,
			Underlying:     types.MarginState{Underlying: "BTCUSD", RequiredMargin: 5000.0},
		},
		WithBuy: types.MarginBreakdownSide{
			RequiredMargin: 6000.0,
			Underlying:     types.MarginState{Underlying: "BTCUSD", RequiredMargin: 6000.0},
		},
		WithSell: types.MarginBreakdownSide{
			RequiredMargin: 4000.0,
			Underlying:     types.MarginState{Underlying: "BTCUSD", RequiredMargin: 4000.0},
		},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.MarginForOrderResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Current.RequiredMargin != 5000.0 {
		t.Errorf("Current.RequiredMargin = %v, want 5000.0", got.Current.RequiredMargin)
	}
	if got.WithBuy.RequiredMargin != 6000.0 {
		t.Errorf("WithBuy.RequiredMargin = %v, want 6000.0", got.WithBuy.RequiredMargin)
	}
	if got.WithSell.RequiredMargin != 4000.0 {
		t.Errorf("WithSell.RequiredMargin = %v, want 4000.0", got.WithSell.RequiredMargin)
	}
}

// ---------- ScenarioPosition / ScenarioAsset JSON ----------

func TestScenarioPosition_JSONRoundTrip(t *testing.T) {
	assumedFilled := true
	sp := types.ScenarioPosition{
		InstrumentName: "BTC-PERPETUAL",
		Position:       5.0,
		InstrumentPnl:  1000.0,
		Pnl:            1000.0,
		CurrentPrice:   50000.0,
		ScenarioPrice:  55000.0,
		OpenOrder:      true,
		AssumedFilled:  &assumedFilled,
	}

	data, err := json.Marshal(sp)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.ScenarioPosition
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, "BTC-PERPETUAL")
	}
	if got.AssumedFilled == nil || *got.AssumedFilled != true {
		t.Errorf("AssumedFilled = %v, want true", got.AssumedFilled)
	}
}

func TestScenarioAsset_JSONRoundTrip(t *testing.T) {
	sa := types.ScenarioAsset{
		AssetName:     "BTC",
		Position:      1.5,
		UnderlyingPnl: 500.0,
		Pnl:           500.0,
		CurrentPrice:  50000.0,
		ScenarioPrice: 55000.0,
	}

	data, err := json.Marshal(sa)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.ScenarioAsset
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.AssetName != "BTC" {
		t.Errorf("AssetName = %q, want %q", got.AssetName, "BTC")
	}
}

// ---------- AccountBreakdown JSON round-trip ----------

func TestAccountBreakdown_JSONRoundTrip(t *testing.T) {
	ab := types.AccountBreakdown{
		Cash: []types.CashHolding{
			{Currency: "BTC", Balance: 1.0, CollateralFactor: 0.95, Transactable: true},
		},
		Portfolio: []types.PortfolioEntry{
			{InstrumentName: "BTC-PERPETUAL", Position: 2.0, MarkPrice: 50000.0},
		},
		UnrealisedPnl:      100.0,
		CashCollateral:     48000.0,
		Margin:             40000.0,
		RequiredMargin:     20000.0,
		RemainingMargin:    20000.0,
		SessionRealisedPnl: 50.0,
	}

	data, err := json.Marshal(ab)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.AccountBreakdown
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Cash) != 1 {
		t.Fatalf("len(Cash) = %d, want 1", len(got.Cash))
	}
	if len(got.Portfolio) != 1 {
		t.Fatalf("len(Portfolio) = %d, want 1", len(got.Portfolio))
	}
	if got.Margin != 40000.0 {
		t.Errorf("Margin = %v, want 40000.0", got.Margin)
	}
}

// ---------- PortfolioMarginBreakdown JSON round-trip ----------

func TestPortfolioMarginBreakdown_JSONRoundTrip(t *testing.T) {
	pmb := types.PortfolioMarginBreakdown{}
	pmb.Portfolio.RequiredMargin = 10000.0
	pmb.Portfolio.Underlyings = []types.UnderlyingMarginDetail{
		{
			Underlying:     "BTCUSD",
			RequiredMargin: 10000.0,
			LossMargin:     8000.0,
			ScenarioUsed:   3,
			Scenarios: []types.MarginScenario{
				{
					UnderlyingChangePct: -0.1,
					VolChangePctPoint:   0.05,
					Pnl:                 -5000.0,
					RequiredMargin:      10000.0,
					Highlight:           true,
				},
			},
		},
	}

	data, err := json.Marshal(pmb)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.PortfolioMarginBreakdown
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Portfolio.RequiredMargin != 10000.0 {
		t.Errorf("Portfolio.RequiredMargin = %v, want 10000.0", got.Portfolio.RequiredMargin)
	}
	if len(got.Portfolio.Underlyings) != 1 {
		t.Fatalf("len(Underlyings) = %d, want 1", len(got.Portfolio.Underlyings))
	}
	if got.Portfolio.Underlyings[0].ScenarioUsed != 3 {
		t.Errorf("ScenarioUsed = %d, want 3", got.Portfolio.Underlyings[0].ScenarioUsed)
	}
	if len(got.Portfolio.Underlyings[0].Scenarios) != 1 {
		t.Fatalf("len(Scenarios) = %d, want 1", len(got.Portfolio.Underlyings[0].Scenarios))
	}
	if !got.Portfolio.Underlyings[0].Scenarios[0].Highlight {
		t.Error("Scenarios[0].Highlight = false, want true")
	}
}
