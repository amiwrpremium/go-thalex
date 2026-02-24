package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/types"
)

func TestPtr_Int(t *testing.T) {
	v := 42
	p := types.Ptr(v)
	if p == nil {
		t.Fatal("Ptr[int] should not return nil")
	}
	if *p != v {
		t.Errorf("*Ptr(42) = %d, want %d", *p, v)
	}
	// Verify it is a new allocation (modifying the pointer does not affect v)
	*p = 99
	if v != 42 {
		t.Error("Ptr should return a new pointer, not a reference to the original")
	}
}

func TestPtr_String(t *testing.T) {
	v := "hello"
	p := types.Ptr(v)
	if p == nil {
		t.Fatal("Ptr[string] should not return nil")
	}
	if *p != v {
		t.Errorf("*Ptr(%q) = %q, want %q", v, *p, v)
	}
}

func TestPtr_Float64(t *testing.T) {
	v := 3.14
	p := types.Ptr(v)
	if p == nil {
		t.Fatal("Ptr[float64] should not return nil")
	}
	if *p != v {
		t.Errorf("*Ptr(%f) = %f, want %f", v, *p, v)
	}
}

func TestPtr_Bool(t *testing.T) {
	v := true
	p := types.Ptr(v)
	if p == nil {
		t.Fatal("Ptr[bool] should not return nil")
	}
	if *p != v {
		t.Errorf("*Ptr(%v) = %v, want %v", v, *p, v)
	}
}

func TestPtr_ZeroValues(t *testing.T) {
	t.Run("int_zero", func(t *testing.T) {
		p := types.Ptr(0)
		if p == nil {
			t.Fatal("Ptr should not return nil for zero int")
		}
		if *p != 0 {
			t.Errorf("*Ptr(0) = %d, want 0", *p)
		}
	})

	t.Run("string_empty", func(t *testing.T) {
		p := types.Ptr("")
		if p == nil {
			t.Fatal("Ptr should not return nil for empty string")
		}
		if *p != "" {
			t.Errorf("*Ptr(\"\") = %q, want empty", *p)
		}
	})

	t.Run("float64_zero", func(t *testing.T) {
		p := types.Ptr(0.0)
		if p == nil {
			t.Fatal("Ptr should not return nil for zero float64")
		}
		if *p != 0.0 {
			t.Errorf("*Ptr(0.0) = %f, want 0.0", *p)
		}
	})
}

// ---------- Leg JSON round-trip ----------

func TestLeg_JSONRoundTrip(t *testing.T) {
	leg := types.Leg{
		InstrumentName:  "BTC-28MAR25-100000-C",
		Quantity:        1.0,
		FilledAmount:    0.5,
		RemainingAmount: 0.5,
	}

	data, err := json.Marshal(leg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Leg
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != leg.InstrumentName {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, leg.InstrumentName)
	}
	if got.Quantity != 1.0 {
		t.Errorf("Quantity = %v, want 1.0", got.Quantity)
	}
	if got.FilledAmount != 0.5 {
		t.Errorf("FilledAmount = %v, want 0.5", got.FilledAmount)
	}
}

func TestLeg_ZeroRemainingOmitted(t *testing.T) {
	leg := types.Leg{
		InstrumentName: "BTC-PERPETUAL",
		Quantity:       1.0,
	}

	data, err := json.Marshal(leg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	if _, ok := raw["remaining_amount"]; ok {
		t.Error("expected remaining_amount to be omitted when zero")
	}
}

// ---------- InsertLeg JSON round-trip ----------

func TestInsertLeg_JSONRoundTrip(t *testing.T) {
	leg := types.InsertLeg{
		InstrumentName: "BTC-28MAR25-100000-C",
		Quantity:       2.0,
	}

	data, err := json.Marshal(leg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.InsertLeg
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != leg.InstrumentName {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, leg.InstrumentName)
	}
	if got.Quantity != 2.0 {
		t.Errorf("Quantity = %v, want 2.0", got.Quantity)
	}
}

// ---------- Asset JSON round-trip ----------

func TestAsset_JSONRoundTrip(t *testing.T) {
	a := types.Asset{
		AssetName: "BTC",
		Amount:    1.5,
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Asset
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.AssetName != "BTC" {
		t.Errorf("AssetName = %q, want %q", got.AssetName, "BTC")
	}
	if got.Amount != 1.5 {
		t.Errorf("Amount = %v, want 1.5", got.Amount)
	}
}

// ---------- PositionTransfer JSON round-trip ----------

func TestPositionTransfer_JSONRoundTrip(t *testing.T) {
	pt := types.PositionTransfer{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         2.5,
	}

	data, err := json.Marshal(pt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.PositionTransfer
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, "BTC-PERPETUAL")
	}
	if got.Amount != 2.5 {
		t.Errorf("Amount = %v, want 2.5", got.Amount)
	}
}

// ---------- Ptr with struct types ----------

func TestPtr_Struct(t *testing.T) {
	leg := types.InsertLeg{InstrumentName: "BTC", Quantity: 1.0}
	p := types.Ptr(leg)
	if p == nil {
		t.Fatal("Ptr[InsertLeg] should not return nil")
	}
	if p.InstrumentName != "BTC" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "BTC")
	}
}
