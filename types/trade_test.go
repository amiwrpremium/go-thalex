package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- Trade JSON round-trip ----------

func TestTrade_JSONRoundTrip(t *testing.T) {
	positionPnl := 10.5
	perpFundingPnl := -0.5
	idx := 50000.0
	fundingMark := 0.0001
	liqFee := 1.0
	clientOID := uint64(12345)
	trade := types.Trade{
		TradeType:            enums.TradeTypeNormal,
		TradeID:              "t-123",
		OrderID:              "o-456",
		InstrumentName:       "BTC-PERPETUAL",
		Direction:            enums.DirectionBuy,
		Price:                50000.0,
		Amount:               1.0,
		Label:                "my-label",
		Time:                 1700000000.0,
		PositionAfter:        5.0,
		SessionRealisedAfter: 100.0,
		PositionPnl:          &positionPnl,
		PerpetualFundingPnl:  &perpFundingPnl,
		Fee:                  0.5,
		Index:                &idx,
		FeeRate:              0.0005,
		FeeBasis:             1000.0,
		FundingMark:          &fundingMark,
		LiquidationFee:       &liqFee,
		ClientOrderID:        &clientOID,
		MakerTaker:           enums.MakerTakerMaker,
		BotID:                "bot-1",
		LegIndex:             0,
	}

	data, err := json.Marshal(trade)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Trade
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.TradeID != trade.TradeID {
		t.Errorf("TradeID = %q, want %q", got.TradeID, trade.TradeID)
	}
	if got.Direction != trade.Direction {
		t.Errorf("Direction = %q, want %q", got.Direction, trade.Direction)
	}
	if got.Price != trade.Price {
		t.Errorf("Price = %v, want %v", got.Price, trade.Price)
	}
	if got.PositionPnl == nil || *got.PositionPnl != positionPnl {
		t.Errorf("PositionPnl = %v, want %v", got.PositionPnl, positionPnl)
	}
	if got.PerpetualFundingPnl == nil || *got.PerpetualFundingPnl != perpFundingPnl {
		t.Errorf("PerpetualFundingPnl = %v, want %v", got.PerpetualFundingPnl, perpFundingPnl)
	}
	if got.ClientOrderID == nil || *got.ClientOrderID != clientOID {
		t.Errorf("ClientOrderID = %v, want %v", got.ClientOrderID, clientOID)
	}
	if got.LiquidationFee == nil || *got.LiquidationFee != liqFee {
		t.Errorf("LiquidationFee = %v, want %v", got.LiquidationFee, liqFee)
	}
	if got.FundingMark == nil || *got.FundingMark != fundingMark {
		t.Errorf("FundingMark = %v, want %v", got.FundingMark, fundingMark)
	}
	if got.BotID != trade.BotID {
		t.Errorf("BotID = %q, want %q", got.BotID, trade.BotID)
	}
	if got.Label != trade.Label {
		t.Errorf("Label = %q, want %q", got.Label, trade.Label)
	}
}

func TestTrade_JSONOmitsNilPointers(t *testing.T) {
	trade := types.Trade{
		TradeID:        "t-1",
		InstrumentName: "ETH-PERPETUAL",
		Direction:      enums.DirectionSell,
		Price:          3000.0,
		Amount:         2.0,
		Time:           1700000000.0,
		Fee:            0.1,
	}

	data, err := json.Marshal(trade)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Fields with omitempty and nil pointer should not appear
	for _, key := range []string{"position_pnl", "perpetual_funding_pnl", "index", "funding_mark", "liquidation_fee", "client_order_id"} {
		if _, ok := raw[key]; ok {
			t.Errorf("expected key %q to be omitted from JSON when nil", key)
		}
	}
}

// ---------- DailyMark JSON round-trip ----------

func TestDailyMark_JSONRoundTrip(t *testing.T) {
	fundingPnl := 1.5
	dm := types.DailyMark{
		Time:                1700000000.0,
		InstrumentName:      "BTC-PERPETUAL",
		Position:            5.0,
		MarkPrice:           50000.0,
		RealizedPositionPnl: 200.0,
		RealizedFundingPnl:  &fundingPnl,
	}

	data, err := json.Marshal(dm)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.DailyMark
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != dm.InstrumentName {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, dm.InstrumentName)
	}
	if got.Position != dm.Position {
		t.Errorf("Position = %v, want %v", got.Position, dm.Position)
	}
	if got.RealizedFundingPnl == nil || *got.RealizedFundingPnl != fundingPnl {
		t.Errorf("RealizedFundingPnl = %v, want %v", got.RealizedFundingPnl, fundingPnl)
	}
}

// ---------- Transaction JSON round-trip ----------

func TestTransaction_JSONRoundTrip(t *testing.T) {
	tx := types.Transaction{
		TransactionID:  "tx-abc",
		Time:           1700000000.0,
		Type:           "trade",
		InstrumentName: "BTC-PERPETUAL",
		Amount:         100.0,
		Currency:       "BTC",
		TradeID:        "t-1",
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Transaction
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.TransactionID != tx.TransactionID {
		t.Errorf("TransactionID = %q, want %q", got.TransactionID, tx.TransactionID)
	}
	if got.Type != tx.Type {
		t.Errorf("Type = %q, want %q", got.Type, tx.Type)
	}
	if got.Currency != tx.Currency {
		t.Errorf("Currency = %q, want %q", got.Currency, tx.Currency)
	}
}

// ---------- TradeHistoryParams JSON round-trip ----------

func TestTradeHistoryParams_JSONRoundTrip(t *testing.T) {
	from := 1700000000.0
	to := 1700003600.0
	offset := 10
	limit := 50

	p := types.TradeHistoryParams{
		From:            &from,
		To:              &to,
		Offset:          &offset,
		Limit:           &limit,
		Sort:            enums.SortDesc,
		InstrumentNames: []string{"BTC-PERPETUAL", "ETH-PERPETUAL"},
		BotIDs:          []string{"bot-1"},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.TradeHistoryParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.From == nil || *got.From != from {
		t.Errorf("From = %v, want %v", got.From, from)
	}
	if got.Sort != enums.SortDesc {
		t.Errorf("Sort = %q, want %q", got.Sort, enums.SortDesc)
	}
	if len(got.InstrumentNames) != 2 {
		t.Errorf("InstrumentNames len = %d, want 2", len(got.InstrumentNames))
	}
	if len(got.BotIDs) != 1 {
		t.Errorf("BotIDs len = %d, want 1", len(got.BotIDs))
	}
}

func TestTradeHistoryParams_JSONOmitsNilFields(t *testing.T) {
	p := types.TradeHistoryParams{}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	for _, key := range []string{"from", "to", "offset", "limit", "instrument_names", "bot_ids"} {
		if _, ok := raw[key]; ok {
			t.Errorf("expected key %q to be omitted when zero/nil", key)
		}
	}
}

// ---------- OrderHistoryParams JSON ----------

func TestOrderHistoryParams_JSONRoundTrip(t *testing.T) {
	from := 100.0
	p := types.OrderHistoryParams{
		From:            &from,
		Sort:            enums.SortAsc,
		InstrumentNames: []string{"BTC-PERPETUAL"},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OrderHistoryParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.From == nil || *got.From != from {
		t.Errorf("From = %v, want %v", got.From, from)
	}
	if got.Sort != enums.SortAsc {
		t.Errorf("Sort = %q, want %q", got.Sort, enums.SortAsc)
	}
}

// ---------- DailyMarkHistoryParams JSON ----------

func TestDailyMarkHistoryParams_JSONRoundTrip(t *testing.T) {
	from := 100.0
	to := 200.0
	offset := 5
	limit := 20
	p := types.DailyMarkHistoryParams{
		From:   &from,
		To:     &to,
		Offset: &offset,
		Limit:  &limit,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.DailyMarkHistoryParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.From == nil || *got.From != from {
		t.Errorf("From = %v, want %v", got.From, from)
	}
	if got.Limit == nil || *got.Limit != limit {
		t.Errorf("Limit = %v, want %v", got.Limit, limit)
	}
}

// ---------- TransactionHistoryParams JSON ----------

func TestTransactionHistoryParams_JSONRoundTrip(t *testing.T) {
	from := 100.0
	p := types.TransactionHistoryParams{
		From: &from,
		Sort: enums.SortDesc,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.TransactionHistoryParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.From == nil || *got.From != from {
		t.Errorf("From = %v, want %v", got.From, from)
	}
	if got.Sort != enums.SortDesc {
		t.Errorf("Sort = %q, want %q", got.Sort, enums.SortDesc)
	}
}
