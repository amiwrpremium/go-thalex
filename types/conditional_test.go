package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- ConditionalOrder predicates ----------

func TestConditionalOrder_IsStopLimit(t *testing.T) {
	tests := []struct {
		name       string
		limitPrice *float64
		want       bool
	}{
		{"with_limit_price", types.Ptr(50000.0), true},
		{"nil_limit_price", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &types.ConditionalOrder{LimitPrice: tt.limitPrice}
			if got := o.IsStopLimit(); got != tt.want {
				t.Errorf("IsStopLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOrder_IsBracket(t *testing.T) {
	tests := []struct {
		name         string
		bracketPrice *float64
		want         bool
	}{
		{"with_bracket_price", types.Ptr(45000.0), true},
		{"nil_bracket_price", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &types.ConditionalOrder{BracketPrice: tt.bracketPrice}
			if got := o.IsBracket(); got != tt.want {
				t.Errorf("IsBracket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOrder_IsTrailingStop(t *testing.T) {
	tests := []struct {
		name         string
		callbackRate *float64
		want         bool
	}{
		{"with_callback_rate", types.Ptr(0.05), true},
		{"nil_callback_rate", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &types.ConditionalOrder{TrailingStopCallbackRate: tt.callbackRate}
			if got := o.IsTrailingStop(); got != tt.want {
				t.Errorf("IsTrailingStop() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------- ConditionalOrder JSON round-trip ----------

func TestConditionalOrder_JSONRoundTrip(t *testing.T) {
	limitPrice := 51000.0
	bracketPrice := 45000.0
	callbackRate := 0.05
	convertTime := 1700001000.0

	o := types.ConditionalOrder{
		OrderID:                  "co-123",
		InstrumentName:           "BTC-PERPETUAL",
		Direction:                enums.DirectionBuy,
		Amount:                   1.0,
		Target:                   enums.TargetMark,
		StopPrice:                50000.0,
		LimitPrice:               &limitPrice,
		BracketPrice:             &bracketPrice,
		TrailingStopCallbackRate: &callbackRate,
		Label:                    "my-stop",
		Status:                   enums.ConditionalOrderStatusActive,
		CreateTime:               1700000000.0,
		UpdateTime:               1700000500.0,
		ConvertTime:              &convertTime,
		ConvertedOrderID:         "order-456",
		RejectReason:             "",
		ReduceOnly:               true,
	}

	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.ConditionalOrder
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.OrderID != o.OrderID {
		t.Errorf("OrderID = %q, want %q", got.OrderID, o.OrderID)
	}
	if got.Direction != o.Direction {
		t.Errorf("Direction = %q, want %q", got.Direction, o.Direction)
	}
	if got.LimitPrice == nil || *got.LimitPrice != limitPrice {
		t.Errorf("LimitPrice = %v, want %v", got.LimitPrice, limitPrice)
	}
	if got.BracketPrice == nil || *got.BracketPrice != bracketPrice {
		t.Errorf("BracketPrice = %v, want %v", got.BracketPrice, bracketPrice)
	}
	if got.TrailingStopCallbackRate == nil || *got.TrailingStopCallbackRate != callbackRate {
		t.Errorf("TrailingStopCallbackRate = %v, want %v", got.TrailingStopCallbackRate, callbackRate)
	}
	if got.Status != o.Status {
		t.Errorf("Status = %q, want %q", got.Status, o.Status)
	}
	if got.ReduceOnly != true {
		t.Errorf("ReduceOnly = %v, want true", got.ReduceOnly)
	}
	if got.ConvertedOrderID != "order-456" {
		t.Errorf("ConvertedOrderID = %q, want %q", got.ConvertedOrderID, "order-456")
	}
}

// ---------- NewStopOrder ----------

func TestNewStopOrder(t *testing.T) {
	p := types.NewStopOrder(enums.DirectionBuy, "BTC-PERPETUAL", 1.0, 50000.0)
	if p == nil {
		t.Fatal("NewStopOrder returned nil")
	}
	if p.Direction != enums.DirectionBuy {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionBuy)
	}
	if p.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "BTC-PERPETUAL")
	}
	if p.Amount != 1.0 {
		t.Errorf("Amount = %v, want 1.0", p.Amount)
	}
	if p.StopPrice != 50000.0 {
		t.Errorf("StopPrice = %v, want 50000.0", p.StopPrice)
	}
	if p.LimitPrice != nil {
		t.Errorf("LimitPrice = %v, want nil", p.LimitPrice)
	}
	if p.BracketPrice != nil {
		t.Errorf("BracketPrice = %v, want nil", p.BracketPrice)
	}
	if p.TrailingStopCallbackRate != nil {
		t.Errorf("TrailingStopCallbackRate = %v, want nil", p.TrailingStopCallbackRate)
	}
}

// ---------- NewStopLimitOrder ----------

func TestNewStopLimitOrder(t *testing.T) {
	p := types.NewStopLimitOrder(enums.DirectionSell, "ETH-PERPETUAL", 2.0, 3000.0, 2950.0)
	if p == nil {
		t.Fatal("NewStopLimitOrder returned nil")
	}
	if p.Direction != enums.DirectionSell {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionSell)
	}
	if p.StopPrice != 3000.0 {
		t.Errorf("StopPrice = %v, want 3000.0", p.StopPrice)
	}
	if p.LimitPrice == nil || *p.LimitPrice != 2950.0 {
		t.Errorf("LimitPrice = %v, want 2950.0", p.LimitPrice)
	}
	if p.BracketPrice != nil {
		t.Errorf("BracketPrice = %v, want nil", p.BracketPrice)
	}
}

// ---------- NewBracketOrder ----------

func TestNewBracketOrder(t *testing.T) {
	p := types.NewBracketOrder(enums.DirectionBuy, "BTC-PERPETUAL", 1.0, 50000.0, 45000.0)
	if p == nil {
		t.Fatal("NewBracketOrder returned nil")
	}
	if p.StopPrice != 50000.0 {
		t.Errorf("StopPrice = %v, want 50000.0", p.StopPrice)
	}
	if p.BracketPrice == nil || *p.BracketPrice != 45000.0 {
		t.Errorf("BracketPrice = %v, want 45000.0", p.BracketPrice)
	}
	if p.LimitPrice != nil {
		t.Errorf("LimitPrice = %v, want nil", p.LimitPrice)
	}
}

// ---------- NewTrailingStopOrder ----------

func TestNewTrailingStopOrder(t *testing.T) {
	p := types.NewTrailingStopOrder(enums.DirectionSell, "BTC-PERPETUAL", 1.0, 50000.0, 0.05)
	if p == nil {
		t.Fatal("NewTrailingStopOrder returned nil")
	}
	if p.TrailingStopCallbackRate == nil || *p.TrailingStopCallbackRate != 0.05 {
		t.Errorf("TrailingStopCallbackRate = %v, want 0.05", p.TrailingStopCallbackRate)
	}
	if p.LimitPrice != nil {
		t.Errorf("LimitPrice = %v, want nil", p.LimitPrice)
	}
	if p.BracketPrice != nil {
		t.Errorf("BracketPrice = %v, want nil", p.BracketPrice)
	}
}

// ---------- CreateConditionalOrderParams builder methods ----------

func TestCreateConditionalOrderParams_WithTarget(t *testing.T) {
	p := types.NewStopOrder(enums.DirectionBuy, "BTC-PERPETUAL", 1.0, 50000.0)
	ret := p.WithTarget(enums.TargetIndex)
	if ret != p {
		t.Error("WithTarget should return the same pointer for chaining")
	}
	if p.Target != enums.TargetIndex {
		t.Errorf("Target = %q, want %q", p.Target, enums.TargetIndex)
	}
}

func TestCreateConditionalOrderParams_WithLabel(t *testing.T) {
	p := types.NewStopOrder(enums.DirectionBuy, "BTC-PERPETUAL", 1.0, 50000.0)
	ret := p.WithLabel("stop-label")
	if ret != p {
		t.Error("WithLabel should return the same pointer for chaining")
	}
	if p.Label != "stop-label" {
		t.Errorf("Label = %q, want %q", p.Label, "stop-label")
	}
}

func TestCreateConditionalOrderParams_WithReduceOnly(t *testing.T) {
	p := types.NewStopOrder(enums.DirectionBuy, "BTC-PERPETUAL", 1.0, 50000.0)
	ret := p.WithReduceOnly(true)
	if ret != p {
		t.Error("WithReduceOnly should return the same pointer for chaining")
	}
	if p.ReduceOnly == nil || *p.ReduceOnly != true {
		t.Errorf("ReduceOnly = %v, want true", p.ReduceOnly)
	}
}

func TestCreateConditionalOrderParams_Chaining(t *testing.T) {
	p := types.NewStopLimitOrder(enums.DirectionSell, "BTC-PERPETUAL", 1.0, 50000.0, 49500.0).
		WithTarget(enums.TargetMark).
		WithLabel("chain-test").
		WithReduceOnly(true)

	if p.Direction != enums.DirectionSell {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionSell)
	}
	if p.Target != enums.TargetMark {
		t.Errorf("Target = %q, want %q", p.Target, enums.TargetMark)
	}
	if p.Label != "chain-test" {
		t.Errorf("Label = %q, want %q", p.Label, "chain-test")
	}
	if p.ReduceOnly == nil || *p.ReduceOnly != true {
		t.Errorf("ReduceOnly = %v, want true", p.ReduceOnly)
	}
}

// ---------- CreateConditionalOrderParams JSON round-trip ----------

func TestCreateConditionalOrderParams_JSONRoundTrip(t *testing.T) {
	p := types.NewStopLimitOrder(enums.DirectionBuy, "BTC-PERPETUAL", 1.0, 50000.0, 49500.0).
		WithTarget(enums.TargetLast).
		WithLabel("json-test").
		WithReduceOnly(false)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.CreateConditionalOrderParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Direction != enums.DirectionBuy {
		t.Errorf("Direction = %q, want %q", got.Direction, enums.DirectionBuy)
	}
	if got.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, "BTC-PERPETUAL")
	}
	if got.StopPrice != 50000.0 {
		t.Errorf("StopPrice = %v, want 50000.0", got.StopPrice)
	}
	if got.LimitPrice == nil || *got.LimitPrice != 49500.0 {
		t.Errorf("LimitPrice = %v, want 49500.0", got.LimitPrice)
	}
	if got.Target != enums.TargetLast {
		t.Errorf("Target = %q, want %q", got.Target, enums.TargetLast)
	}
	if got.Label != "json-test" {
		t.Errorf("Label = %q, want %q", got.Label, "json-test")
	}
}
