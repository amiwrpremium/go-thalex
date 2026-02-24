package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- Rfq JSON round-trip ----------

func TestRfq_JSONRoundTrip(t *testing.T) {
	validUntil := 1700001000.0
	volTickSize := 0.01
	tradePrice := 50000.0
	tradeAmount := 1.0
	closeTime := 1700002000.0

	rfq := types.Rfq{
		RfqID: "rfq-123",
		Legs: []types.RfqLeg{
			{InstrumentName: "BTC-PERPETUAL", Quantity: 1.0, FeeQuantity: 0.5},
		},
		Amount:         2.0,
		CreateTime:     1700000000.0,
		ValidUntil:     &validUntil,
		Label:          "my-rfq",
		InsertReason:   enums.RfqInsertReasonClientRequest,
		DeleteReason:   "",
		VolumeTickSize: &volTickSize,
		QuotedBid:      &types.RfqQuotedSide{Price: 49000.0, Amount: 1.0},
		QuotedAsk:      &types.RfqQuotedSide{Price: 51000.0, Amount: 1.0},
		TradePrice:     &tradePrice,
		TradeAmount:    &tradeAmount,
		CloseTime:      &closeTime,
		Event:          enums.RfqEventCreated,
	}

	data, err := json.Marshal(rfq)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Rfq
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.RfqID != rfq.RfqID {
		t.Errorf("RfqID = %q, want %q", got.RfqID, rfq.RfqID)
	}
	if len(got.Legs) != 1 {
		t.Fatalf("len(Legs) = %d, want 1", len(got.Legs))
	}
	if got.Legs[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("Legs[0].InstrumentName = %q, want %q", got.Legs[0].InstrumentName, "BTC-PERPETUAL")
	}
	if got.Legs[0].FeeQuantity != 0.5 {
		t.Errorf("Legs[0].FeeQuantity = %v, want 0.5", got.Legs[0].FeeQuantity)
	}
	if got.ValidUntil == nil || *got.ValidUntil != validUntil {
		t.Errorf("ValidUntil = %v, want %v", got.ValidUntil, validUntil)
	}
	if got.QuotedBid == nil || got.QuotedBid.Price != 49000.0 {
		t.Errorf("QuotedBid.Price = %v, want 49000.0", got.QuotedBid)
	}
	if got.QuotedAsk == nil || got.QuotedAsk.Amount != 1.0 {
		t.Errorf("QuotedAsk.Amount = %v, want 1.0", got.QuotedAsk)
	}
	if got.TradePrice == nil || *got.TradePrice != tradePrice {
		t.Errorf("TradePrice = %v, want %v", got.TradePrice, tradePrice)
	}
	if got.Event != enums.RfqEventCreated {
		t.Errorf("Event = %q, want %q", got.Event, enums.RfqEventCreated)
	}
}

func TestRfq_MinimalJSON(t *testing.T) {
	rfq := types.Rfq{
		RfqID:      "rfq-456",
		Amount:     1.0,
		CreateTime: 1700000000.0,
	}

	data, err := json.Marshal(rfq)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Omitempty fields should be absent
	for _, key := range []string{"valid_until", "volume_tick_size", "quoted_bid", "quoted_ask", "trade_price", "trade_amount", "close_time"} {
		if _, ok := raw[key]; ok {
			t.Errorf("expected key %q to be omitted", key)
		}
	}
}

// ---------- RfqOrder JSON round-trip ----------

func TestRfqOrder_JSONRoundTrip(t *testing.T) {
	clientOID := uint64(42)
	tradePrice := 50000.0
	tradeAmount := 1.0

	o := types.RfqOrder{
		RfqID:         "rfq-123",
		OrderID:       "order-456",
		ClientOrderID: &clientOID,
		Direction:     enums.DirectionBuy,
		Price:         49500.0,
		Amount:        2.0,
		Label:         "my-quote",
		TradePrice:    &tradePrice,
		TradeAmount:   &tradeAmount,
		DeleteReason:  enums.RfqDeleteReasonFilled,
		Event:         enums.RfqOrderEventFilled,
	}

	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.RfqOrder
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.RfqID != o.RfqID {
		t.Errorf("RfqID = %q, want %q", got.RfqID, o.RfqID)
	}
	if got.ClientOrderID == nil || *got.ClientOrderID != clientOID {
		t.Errorf("ClientOrderID = %v, want %v", got.ClientOrderID, clientOID)
	}
	if got.Direction != enums.DirectionBuy {
		t.Errorf("Direction = %q, want %q", got.Direction, enums.DirectionBuy)
	}
	if got.Event != enums.RfqOrderEventFilled {
		t.Errorf("Event = %q, want %q", got.Event, enums.RfqOrderEventFilled)
	}
}

// ---------- CreateRfqParams JSON round-trip ----------

func TestCreateRfqParams_JSONRoundTrip(t *testing.T) {
	p := types.CreateRfqParams{
		Legs: []types.InsertLeg{
			{InstrumentName: "BTC-28MAR25-100000-C", Quantity: 1.0},
			{InstrumentName: "BTC-28MAR25-110000-C", Quantity: -1.0},
		},
		Amount: 5.0,
		Label:  "spread-rfq",
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.CreateRfqParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Legs) != 2 {
		t.Fatalf("len(Legs) = %d, want 2", len(got.Legs))
	}
	if got.Amount != 5.0 {
		t.Errorf("Amount = %v, want 5.0", got.Amount)
	}
	if got.Label != "spread-rfq" {
		t.Errorf("Label = %q, want %q", got.Label, "spread-rfq")
	}
}

// ---------- TradeRfqParams JSON round-trip ----------

func TestTradeRfqParams_JSONRoundTrip(t *testing.T) {
	p := types.TradeRfqParams{
		RfqID:     "rfq-123",
		Direction: enums.DirectionSell,
		Price:     50000.0,
		Amount:    1.0,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.TradeRfqParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.RfqID != "rfq-123" {
		t.Errorf("RfqID = %q, want %q", got.RfqID, "rfq-123")
	}
	if got.Direction != enums.DirectionSell {
		t.Errorf("Direction = %q, want %q", got.Direction, enums.DirectionSell)
	}
}

// ---------- RfqQuoteInsertParams JSON round-trip ----------

func TestRfqQuoteInsertParams_JSONRoundTrip(t *testing.T) {
	clientOID := uint64(99)
	p := types.RfqQuoteInsertParams{
		RfqID:         "rfq-123",
		Direction:     enums.DirectionBuy,
		Amount:        2.0,
		Price:         49500.0,
		ClientOrderID: &clientOID,
		Label:         "quote-label",
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.RfqQuoteInsertParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.RfqID != "rfq-123" {
		t.Errorf("RfqID = %q, want %q", got.RfqID, "rfq-123")
	}
	if got.ClientOrderID == nil || *got.ClientOrderID != clientOID {
		t.Errorf("ClientOrderID = %v, want %v", got.ClientOrderID, clientOID)
	}
}

// ---------- RfqQuoteAmendParams JSON round-trip ----------

func TestRfqQuoteAmendParams_JSONRoundTrip(t *testing.T) {
	clientOID := uint64(88)
	p := types.RfqQuoteAmendParams{
		OrderID:       "order-1",
		ClientOrderID: &clientOID,
		Amount:        3.0,
		Price:         50000.0,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.RfqQuoteAmendParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.OrderID != "order-1" {
		t.Errorf("OrderID = %q, want %q", got.OrderID, "order-1")
	}
	if got.ClientOrderID == nil || *got.ClientOrderID != clientOID {
		t.Errorf("ClientOrderID = %v, want %v", got.ClientOrderID, clientOID)
	}
}

// ---------- RfqQuoteDeleteParams JSON round-trip ----------

func TestRfqQuoteDeleteParams_JSONRoundTrip(t *testing.T) {
	t.Run("by_order_id", func(t *testing.T) {
		p := types.RfqQuoteDeleteParams{OrderID: "order-123"}

		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var got types.RfqQuoteDeleteParams
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if got.OrderID != "order-123" {
			t.Errorf("OrderID = %q, want %q", got.OrderID, "order-123")
		}
	})

	t.Run("by_client_order_id", func(t *testing.T) {
		clientOID := uint64(77)
		p := types.RfqQuoteDeleteParams{ClientOrderID: &clientOID}

		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var got types.RfqQuoteDeleteParams
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if got.ClientOrderID == nil || *got.ClientOrderID != 77 {
			t.Errorf("ClientOrderID = %v, want 77", got.ClientOrderID)
		}
	})
}

// ---------- RfqQuotedSide JSON round-trip ----------

func TestRfqQuotedSide_JSONRoundTrip(t *testing.T) {
	s := types.RfqQuotedSide{Price: 50000.0, Amount: 5.0}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.RfqQuotedSide
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Price != 50000.0 {
		t.Errorf("Price = %v, want 50000", got.Price)
	}
	if got.Amount != 5.0 {
		t.Errorf("Amount = %v, want 5.0", got.Amount)
	}
}
