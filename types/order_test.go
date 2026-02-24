package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- NewInsertOrderParams ----------

func TestNewInsertOrderParams(t *testing.T) {
	p := types.NewInsertOrderParams(enums.DirectionBuy, "BTC-PERPETUAL", 1.5)
	if p == nil {
		t.Fatal("NewInsertOrderParams returned nil")
	}
	if p.Direction != enums.DirectionBuy {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionBuy)
	}
	if p.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "BTC-PERPETUAL")
	}
	if p.Amount != 1.5 {
		t.Errorf("Amount = %v, want %v", p.Amount, 1.5)
	}
	// Optional fields should be zero-valued.
	if p.Price != nil {
		t.Errorf("Price = %v, want nil", *p.Price)
	}
	if p.Legs != nil {
		t.Errorf("Legs = %v, want nil", p.Legs)
	}
}

// ---------- NewBuyOrderParams / NewSellOrderParams ----------

func TestNewBuyOrderParams(t *testing.T) {
	p := types.NewBuyOrderParams("ETH-PERPETUAL", 10.0)
	if p == nil {
		t.Fatal("NewBuyOrderParams returned nil")
	}
	if p.Direction != enums.DirectionBuy {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionBuy)
	}
	if p.InstrumentName != "ETH-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "ETH-PERPETUAL")
	}
	if p.Amount != 10.0 {
		t.Errorf("Amount = %v, want %v", p.Amount, 10.0)
	}
}

func TestNewSellOrderParams(t *testing.T) {
	p := types.NewSellOrderParams("ETH-PERPETUAL", 5.0)
	if p == nil {
		t.Fatal("NewSellOrderParams returned nil")
	}
	if p.Direction != enums.DirectionSell {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionSell)
	}
	if p.InstrumentName != "ETH-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "ETH-PERPETUAL")
	}
	if p.Amount != 5.0 {
		t.Errorf("Amount = %v, want %v", p.Amount, 5.0)
	}
}

// ---------- NewComboInsertOrderParams ----------

func TestNewComboInsertOrderParams(t *testing.T) {
	legs := []types.InsertLeg{
		{InstrumentName: "BTC-28MAR25-100000-C", Quantity: 1},
		{InstrumentName: "BTC-28MAR25-110000-C", Quantity: -1},
	}
	p := types.NewComboInsertOrderParams(enums.DirectionBuy, legs, 2.0)
	if p == nil {
		t.Fatal("NewComboInsertOrderParams returned nil")
	}
	if p.Direction != enums.DirectionBuy {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionBuy)
	}
	if p.InstrumentName != "" {
		t.Errorf("InstrumentName = %q, want empty for combo order", p.InstrumentName)
	}
	if len(p.Legs) != 2 {
		t.Fatalf("len(Legs) = %d, want 2", len(p.Legs))
	}
	if p.Legs[0].InstrumentName != "BTC-28MAR25-100000-C" {
		t.Errorf("Legs[0].InstrumentName = %q, want %q", p.Legs[0].InstrumentName, "BTC-28MAR25-100000-C")
	}
	if p.Legs[1].Quantity != -1 {
		t.Errorf("Legs[1].Quantity = %v, want %v", p.Legs[1].Quantity, -1.0)
	}
	if p.Amount != 2.0 {
		t.Errorf("Amount = %v, want %v", p.Amount, 2.0)
	}
}

// ---------- Builder methods: chaining returns same pointer ----------

func TestInsertOrderParams_WithPrice(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithPrice(50000.0)
	if ret != p {
		t.Error("WithPrice should return the same pointer for chaining")
	}
	if p.Price == nil {
		t.Fatal("Price is nil after WithPrice")
	}
	if *p.Price != 50000.0 {
		t.Errorf("Price = %v, want %v", *p.Price, 50000.0)
	}
}

func TestInsertOrderParams_WithOrderType(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithOrderType(enums.OrderTypeLimit)
	if ret != p {
		t.Error("WithOrderType should return the same pointer for chaining")
	}
	if p.OrderType != enums.OrderTypeLimit {
		t.Errorf("OrderType = %q, want %q", p.OrderType, enums.OrderTypeLimit)
	}

	p.WithOrderType(enums.OrderTypeMarket)
	if p.OrderType != enums.OrderTypeMarket {
		t.Errorf("OrderType = %q, want %q after update", p.OrderType, enums.OrderTypeMarket)
	}
}

func TestInsertOrderParams_WithTimeInForce(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithTimeInForce(enums.TimeInForceGoodTillCancelled)
	if ret != p {
		t.Error("WithTimeInForce should return the same pointer for chaining")
	}
	if p.TimeInForce != enums.TimeInForceGoodTillCancelled {
		t.Errorf("TimeInForce = %q, want %q", p.TimeInForce, enums.TimeInForceGoodTillCancelled)
	}

	p.WithTimeInForce(enums.TimeInForceImmediateOrCancel)
	if p.TimeInForce != enums.TimeInForceImmediateOrCancel {
		t.Errorf("TimeInForce = %q, want %q after update", p.TimeInForce, enums.TimeInForceImmediateOrCancel)
	}
}

func TestInsertOrderParams_WithPostOnly(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithPostOnly(true)
	if ret != p {
		t.Error("WithPostOnly should return the same pointer for chaining")
	}
	if p.PostOnly == nil || *p.PostOnly != true {
		t.Errorf("PostOnly = %v, want true", p.PostOnly)
	}

	p.WithPostOnly(false)
	if p.PostOnly == nil || *p.PostOnly != false {
		t.Errorf("PostOnly = %v, want false after update", p.PostOnly)
	}
}

func TestInsertOrderParams_WithRejectPostOnly(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithRejectPostOnly(true)
	if ret != p {
		t.Error("WithRejectPostOnly should return the same pointer for chaining")
	}
	if p.RejectPostOnly == nil || *p.RejectPostOnly != true {
		t.Errorf("RejectPostOnly = %v, want true", p.RejectPostOnly)
	}

	p.WithRejectPostOnly(false)
	if p.RejectPostOnly == nil || *p.RejectPostOnly != false {
		t.Errorf("RejectPostOnly = %v, want false after update", p.RejectPostOnly)
	}
}

func TestInsertOrderParams_WithReduceOnly(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithReduceOnly(true)
	if ret != p {
		t.Error("WithReduceOnly should return the same pointer for chaining")
	}
	if p.ReduceOnly == nil || *p.ReduceOnly != true {
		t.Errorf("ReduceOnly = %v, want true", p.ReduceOnly)
	}
}

func TestInsertOrderParams_WithCollar(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithCollar(enums.CollarReject)
	if ret != p {
		t.Error("WithCollar should return the same pointer for chaining")
	}
	if p.Collar != enums.CollarReject {
		t.Errorf("Collar = %q, want %q", p.Collar, enums.CollarReject)
	}
}

func TestInsertOrderParams_WithLabel(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithLabel("my-label")
	if ret != p {
		t.Error("WithLabel should return the same pointer for chaining")
	}
	if p.Label != "my-label" {
		t.Errorf("Label = %q, want %q", p.Label, "my-label")
	}
}

func TestInsertOrderParams_WithClientOrderID(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithClientOrderID(12345)
	if ret != p {
		t.Error("WithClientOrderID should return the same pointer for chaining")
	}
	if p.ClientOrderID == nil {
		t.Fatal("ClientOrderID is nil after WithClientOrderID")
	}
	if *p.ClientOrderID != 12345 {
		t.Errorf("ClientOrderID = %d, want %d", *p.ClientOrderID, 12345)
	}
}

func TestInsertOrderParams_WithSTP(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0)
	ret := p.WithSTP(enums.STPLevelAccount, enums.STPActionCancelAggressor)
	if ret != p {
		t.Error("WithSTP should return the same pointer for chaining")
	}
	if p.STPLevel != enums.STPLevelAccount {
		t.Errorf("STPLevel = %q, want %q", p.STPLevel, enums.STPLevelAccount)
	}
	if p.STPAction != enums.STPActionCancelAggressor {
		t.Errorf("STPAction = %q, want %q", p.STPAction, enums.STPActionCancelAggressor)
	}
}

// ---------- Builder chaining (multiple With* calls) ----------

func TestInsertOrderParams_Chaining(t *testing.T) {
	p := types.NewSellOrderParams("BTC-PERPETUAL", 2.0).
		WithPrice(48000.0).
		WithOrderType(enums.OrderTypeLimit).
		WithTimeInForce(enums.TimeInForceGoodTillCancelled).
		WithPostOnly(true).
		WithCollar(enums.CollarClamp).
		WithLabel("chain-test").
		WithClientOrderID(99).
		WithSTP(enums.STPLevelCustomer, enums.STPActionCancelBoth).
		WithReduceOnly(false).
		WithRejectPostOnly(false)

	if p.Direction != enums.DirectionSell {
		t.Errorf("Direction = %q, want %q", p.Direction, enums.DirectionSell)
	}
	if p.Amount != 2.0 {
		t.Errorf("Amount = %v, want %v", p.Amount, 2.0)
	}
	if p.Price == nil || *p.Price != 48000.0 {
		t.Errorf("Price = %v, want 48000.0", p.Price)
	}
	if p.OrderType != enums.OrderTypeLimit {
		t.Errorf("OrderType = %q, want %q", p.OrderType, enums.OrderTypeLimit)
	}
	if p.TimeInForce != enums.TimeInForceGoodTillCancelled {
		t.Errorf("TimeInForce = %q, want %q", p.TimeInForce, enums.TimeInForceGoodTillCancelled)
	}
	if p.PostOnly == nil || *p.PostOnly != true {
		t.Errorf("PostOnly = %v, want true", p.PostOnly)
	}
	if p.RejectPostOnly == nil || *p.RejectPostOnly != false {
		t.Errorf("RejectPostOnly = %v, want false", p.RejectPostOnly)
	}
	if p.ReduceOnly == nil || *p.ReduceOnly != false {
		t.Errorf("ReduceOnly = %v, want false", p.ReduceOnly)
	}
	if p.Collar != enums.CollarClamp {
		t.Errorf("Collar = %q, want %q", p.Collar, enums.CollarClamp)
	}
	if p.Label != "chain-test" {
		t.Errorf("Label = %q, want %q", p.Label, "chain-test")
	}
	if p.ClientOrderID == nil || *p.ClientOrderID != 99 {
		t.Errorf("ClientOrderID = %v, want 99", p.ClientOrderID)
	}
	if p.STPLevel != enums.STPLevelCustomer {
		t.Errorf("STPLevel = %q, want %q", p.STPLevel, enums.STPLevelCustomer)
	}
	if p.STPAction != enums.STPActionCancelBoth {
		t.Errorf("STPAction = %q, want %q", p.STPAction, enums.STPActionCancelBoth)
	}
}

// ---------- AmendOrderParams ----------

func TestNewAmendByOrderID(t *testing.T) {
	p := types.NewAmendByOrderID("order-abc-123", 51000.0, 3.0)
	if p == nil {
		t.Fatal("NewAmendByOrderID returned nil")
	}
	if p.OrderID != "order-abc-123" {
		t.Errorf("OrderID = %q, want %q", p.OrderID, "order-abc-123")
	}
	if p.ClientOrderID != nil {
		t.Errorf("ClientOrderID = %v, want nil", p.ClientOrderID)
	}
	if p.Price != 51000.0 {
		t.Errorf("Price = %v, want %v", p.Price, 51000.0)
	}
	if p.Amount != 3.0 {
		t.Errorf("Amount = %v, want %v", p.Amount, 3.0)
	}
}

func TestNewAmendByClientOrderID(t *testing.T) {
	p := types.NewAmendByClientOrderID(42, 52000.0, 4.0)
	if p == nil {
		t.Fatal("NewAmendByClientOrderID returned nil")
	}
	if p.OrderID != "" {
		t.Errorf("OrderID = %q, want empty", p.OrderID)
	}
	if p.ClientOrderID == nil {
		t.Fatal("ClientOrderID is nil, want non-nil")
	}
	if *p.ClientOrderID != 42 {
		t.Errorf("ClientOrderID = %d, want %d", *p.ClientOrderID, 42)
	}
	if p.Price != 52000.0 {
		t.Errorf("Price = %v, want %v", p.Price, 52000.0)
	}
	if p.Amount != 4.0 {
		t.Errorf("Amount = %v, want %v", p.Amount, 4.0)
	}
}

func TestAmendOrderParams_WithCollar(t *testing.T) {
	p := types.NewAmendByOrderID("order-1", 50000.0, 1.0)
	ret := p.WithCollar(enums.CollarIgnore)
	if ret != p {
		t.Error("AmendOrderParams.WithCollar should return the same pointer for chaining")
	}
	if p.Collar != enums.CollarIgnore {
		t.Errorf("Collar = %q, want %q", p.Collar, enums.CollarIgnore)
	}
}

// ---------- CancelOrderParams ----------

func TestCancelByOrderID(t *testing.T) {
	p := types.CancelByOrderID("order-xyz")
	if p == nil {
		t.Fatal("CancelByOrderID returned nil")
	}
	if p.OrderID != "order-xyz" {
		t.Errorf("OrderID = %q, want %q", p.OrderID, "order-xyz")
	}
	if p.ClientOrderID != nil {
		t.Errorf("ClientOrderID = %v, want nil", p.ClientOrderID)
	}
}

func TestCancelByClientOrderID(t *testing.T) {
	p := types.CancelByClientOrderID(777)
	if p == nil {
		t.Fatal("CancelByClientOrderID returned nil")
	}
	if p.OrderID != "" {
		t.Errorf("OrderID = %q, want empty", p.OrderID)
	}
	if p.ClientOrderID == nil {
		t.Fatal("ClientOrderID is nil, want non-nil")
	}
	if *p.ClientOrderID != 777 {
		t.Errorf("ClientOrderID = %d, want %d", *p.ClientOrderID, 777)
	}
}

// ---------- OrderStatus JSON round-trip ----------

func TestOrderStatus_JSONRoundTrip(t *testing.T) {
	price := 50000.0
	clientOID := uint64(42)
	closeTime := 1700001000.0

	os := types.OrderStatus{
		OrderID:         "order-123",
		OrderType:       enums.OrderTypeLimit,
		TimeInForce:     enums.TimeInForceGoodTillCancelled,
		InstrumentName:  "BTC-PERPETUAL",
		Direction:       enums.DirectionBuy,
		Price:           &price,
		Amount:          1.0,
		FilledAmount:    0.5,
		RemainingAmount: 0.5,
		Label:           "test-order",
		ClientOrderID:   &clientOID,
		Status:          enums.OrderStatusPartiallyFilled,
		Fills: []types.OrderFill{
			{TradeID: "t-1", Price: 50000.0, Amount: 0.5, Time: 1700000000.0, MakerTaker: enums.MakerTakerMaker, LegIndex: 0},
		},
		ChangeReason:       enums.ChangeReasonFill,
		DeleteReason:       "",
		InsertReason:       enums.InsertReasonClientRequest,
		ConditionalOrderID: "",
		BotID:              "",
		CreateTime:         1700000000.0,
		CloseTime:          &closeTime,
		ReduceOnly:         false,
		Persistent:         true,
	}

	data, err := json.Marshal(os)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OrderStatus
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.OrderID != "order-123" {
		t.Errorf("OrderID = %q, want %q", got.OrderID, "order-123")
	}
	if got.OrderType != enums.OrderTypeLimit {
		t.Errorf("OrderType = %q, want %q", got.OrderType, enums.OrderTypeLimit)
	}
	if got.Status != enums.OrderStatusPartiallyFilled {
		t.Errorf("Status = %q, want %q", got.Status, enums.OrderStatusPartiallyFilled)
	}
	if got.Price == nil || *got.Price != price {
		t.Errorf("Price = %v, want %v", got.Price, price)
	}
	if got.ClientOrderID == nil || *got.ClientOrderID != clientOID {
		t.Errorf("ClientOrderID = %v, want %v", got.ClientOrderID, clientOID)
	}
	if len(got.Fills) != 1 {
		t.Fatalf("len(Fills) = %d, want 1", len(got.Fills))
	}
	if got.Fills[0].TradeID != "t-1" {
		t.Errorf("Fills[0].TradeID = %q, want %q", got.Fills[0].TradeID, "t-1")
	}
	if got.Fills[0].MakerTaker != enums.MakerTakerMaker {
		t.Errorf("Fills[0].MakerTaker = %q, want %q", got.Fills[0].MakerTaker, enums.MakerTakerMaker)
	}
	if got.CloseTime == nil || *got.CloseTime != closeTime {
		t.Errorf("CloseTime = %v, want %v", got.CloseTime, closeTime)
	}
	if !got.Persistent {
		t.Error("Persistent = false, want true")
	}
}

// ---------- OrderStatus with Legs (combo order) ----------

func TestOrderStatus_ComboLegsJSON(t *testing.T) {
	os := types.OrderStatus{
		OrderID:   "order-combo",
		Direction: enums.DirectionBuy,
		Legs: []types.Leg{
			{InstrumentName: "BTC-28MAR25-100000-C", Quantity: 1.0, FilledAmount: 0.5, RemainingAmount: 0.5},
			{InstrumentName: "BTC-28MAR25-110000-C", Quantity: -1.0, FilledAmount: -0.5, RemainingAmount: -0.5},
		},
		Amount:       1.0,
		Status:       enums.OrderStatusOpen,
		ChangeReason: enums.ChangeReasonInsert,
		InsertReason: enums.InsertReasonClientRequest,
	}

	data, err := json.Marshal(os)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OrderStatus
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Legs) != 2 {
		t.Fatalf("len(Legs) = %d, want 2", len(got.Legs))
	}
	if got.Legs[0].InstrumentName != "BTC-28MAR25-100000-C" {
		t.Errorf("Legs[0].InstrumentName = %q, want %q", got.Legs[0].InstrumentName, "BTC-28MAR25-100000-C")
	}
}

// ---------- OrderHistory JSON round-trip ----------

func TestOrderHistory_JSONRoundTrip(t *testing.T) {
	price := 50000.0
	oh := types.OrderHistory{
		OrderID:        "order-hist-1",
		OrderType:      enums.OrderTypeLimit,
		InstrumentName: "BTC-PERPETUAL",
		Direction:      enums.DirectionSell,
		Price:          &price,
		Amount:         2.0,
		FilledAmount:   2.0,
		Status:         enums.OrderStatusFilled,
		Fills: []types.OrderFill{
			{TradeID: "t-1", Price: 50000.0, Amount: 2.0, MakerTaker: enums.MakerTakerTaker},
		},
		InsertReason: enums.InsertReasonClientRequest,
		DeleteReason: enums.DeleteReasonFilled,
		CreateTime:   1700000000.0,
		CloseTime:    1700000100.0,
	}

	data, err := json.Marshal(oh)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OrderHistory
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.OrderID != "order-hist-1" {
		t.Errorf("OrderID = %q, want %q", got.OrderID, "order-hist-1")
	}
	if got.Status != enums.OrderStatusFilled {
		t.Errorf("Status = %q, want %q", got.Status, enums.OrderStatusFilled)
	}
	if got.DeleteReason != enums.DeleteReasonFilled {
		t.Errorf("DeleteReason = %q, want %q", got.DeleteReason, enums.DeleteReasonFilled)
	}
}

// ---------- InsertOrderParams JSON round-trip ----------

func TestInsertOrderParams_JSONRoundTrip(t *testing.T) {
	p := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0).
		WithPrice(50000.0).
		WithOrderType(enums.OrderTypeLimit).
		WithTimeInForce(enums.TimeInForceGoodTillCancelled).
		WithLabel("json-test").
		WithClientOrderID(42)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.InsertOrderParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Direction != enums.DirectionBuy {
		t.Errorf("Direction = %q, want %q", got.Direction, enums.DirectionBuy)
	}
	if got.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, "BTC-PERPETUAL")
	}
	if got.Price == nil || *got.Price != 50000.0 {
		t.Errorf("Price = %v, want 50000.0", got.Price)
	}
	if got.OrderType != enums.OrderTypeLimit {
		t.Errorf("OrderType = %q, want %q", got.OrderType, enums.OrderTypeLimit)
	}
	if got.ClientOrderID == nil || *got.ClientOrderID != 42 {
		t.Errorf("ClientOrderID = %v, want 42", got.ClientOrderID)
	}
}

// ---------- AmendOrderParams JSON round-trip ----------

func TestAmendOrderParams_JSONRoundTrip(t *testing.T) {
	p := types.NewAmendByOrderID("order-1", 51000.0, 2.0).
		WithCollar(enums.CollarClamp)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.AmendOrderParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.OrderID != "order-1" {
		t.Errorf("OrderID = %q, want %q", got.OrderID, "order-1")
	}
	if got.Price != 51000.0 {
		t.Errorf("Price = %v, want 51000.0", got.Price)
	}
	if got.Amount != 2.0 {
		t.Errorf("Amount = %v, want 2.0", got.Amount)
	}
	if got.Collar != enums.CollarClamp {
		t.Errorf("Collar = %q, want %q", got.Collar, enums.CollarClamp)
	}
}

func TestAmendByClientOrderID_JSONRoundTrip(t *testing.T) {
	p := types.NewAmendByClientOrderID(99, 52000.0, 3.0)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.AmendOrderParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.ClientOrderID == nil || *got.ClientOrderID != 99 {
		t.Errorf("ClientOrderID = %v, want 99", got.ClientOrderID)
	}
}

// ---------- CancelOrderParams JSON round-trip ----------

func TestCancelOrderParams_JSONRoundTrip(t *testing.T) {
	t.Run("by_order_id", func(t *testing.T) {
		p := types.CancelByOrderID("order-xyz")
		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var got types.CancelOrderParams
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if got.OrderID != "order-xyz" {
			t.Errorf("OrderID = %q, want %q", got.OrderID, "order-xyz")
		}
	})

	t.Run("by_client_order_id", func(t *testing.T) {
		p := types.CancelByClientOrderID(555)
		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var got types.CancelOrderParams
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if got.ClientOrderID == nil || *got.ClientOrderID != 555 {
			t.Errorf("ClientOrderID = %v, want 555", got.ClientOrderID)
		}
	})
}

// ---------- OrderFill JSON round-trip ----------

func TestOrderFill_JSONRoundTrip(t *testing.T) {
	fill := types.OrderFill{
		TradeID:    "t-abc",
		Price:      50000.0,
		Amount:     1.5,
		Time:       1700000000.0,
		MakerTaker: enums.MakerTakerTaker,
		LegIndex:   1,
	}

	data, err := json.Marshal(fill)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OrderFill
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.TradeID != "t-abc" {
		t.Errorf("TradeID = %q, want %q", got.TradeID, "t-abc")
	}
	if got.LegIndex != 1 {
		t.Errorf("LegIndex = %d, want 1", got.LegIndex)
	}
	if got.MakerTaker != enums.MakerTakerTaker {
		t.Errorf("MakerTaker = %q, want %q", got.MakerTaker, enums.MakerTakerTaker)
	}
}
