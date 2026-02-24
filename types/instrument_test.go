package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- Instrument type predicates ----------

func TestInstrument_IsOption(t *testing.T) {
	tests := []struct {
		name string
		typ  enums.InstrumentType
		want bool
	}{
		{"option_true", enums.InstrumentTypeOption, true},
		{"future_false", enums.InstrumentTypeFuture, false},
		{"perpetual_false", enums.InstrumentTypePerpetual, false},
		{"combination_false", enums.InstrumentTypeCombination, false},
		{"empty_false", enums.InstrumentType(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &types.Instrument{Type: tt.typ}
			if got := inst.IsOption(); got != tt.want {
				t.Errorf("Instrument{Type: %q}.IsOption() = %v, want %v", tt.typ, got, tt.want)
			}
		})
	}
}

func TestInstrument_IsFuture(t *testing.T) {
	tests := []struct {
		name string
		typ  enums.InstrumentType
		want bool
	}{
		{"future_true", enums.InstrumentTypeFuture, true},
		{"option_false", enums.InstrumentTypeOption, false},
		{"perpetual_false", enums.InstrumentTypePerpetual, false},
		{"combination_false", enums.InstrumentTypeCombination, false},
		{"empty_false", enums.InstrumentType(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &types.Instrument{Type: tt.typ}
			if got := inst.IsFuture(); got != tt.want {
				t.Errorf("Instrument{Type: %q}.IsFuture() = %v, want %v", tt.typ, got, tt.want)
			}
		})
	}
}

func TestInstrument_IsPerpetual(t *testing.T) {
	tests := []struct {
		name string
		typ  enums.InstrumentType
		want bool
	}{
		{"perpetual_true", enums.InstrumentTypePerpetual, true},
		{"option_false", enums.InstrumentTypeOption, false},
		{"future_false", enums.InstrumentTypeFuture, false},
		{"combination_false", enums.InstrumentTypeCombination, false},
		{"empty_false", enums.InstrumentType(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &types.Instrument{Type: tt.typ}
			if got := inst.IsPerpetual(); got != tt.want {
				t.Errorf("Instrument{Type: %q}.IsPerpetual() = %v, want %v", tt.typ, got, tt.want)
			}
		})
	}
}

func TestInstrument_IsCombination(t *testing.T) {
	tests := []struct {
		name string
		typ  enums.InstrumentType
		want bool
	}{
		{"combination_true", enums.InstrumentTypeCombination, true},
		{"option_false", enums.InstrumentTypeOption, false},
		{"future_false", enums.InstrumentTypeFuture, false},
		{"perpetual_false", enums.InstrumentTypePerpetual, false},
		{"empty_false", enums.InstrumentType(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &types.Instrument{Type: tt.typ}
			if got := inst.IsCombination(); got != tt.want {
				t.Errorf("Instrument{Type: %q}.IsCombination() = %v, want %v", tt.typ, got, tt.want)
			}
		})
	}
}

// ---------- Ticker.Spread ----------

func TestTicker_Spread(t *testing.T) {
	tests := []struct {
		name    string
		bid     *float64
		ask     *float64
		wantNil bool
		want    float64
	}{
		{
			name:    "both_set",
			bid:     types.Ptr(100.0),
			ask:     types.Ptr(105.0),
			wantNil: false,
			want:    5.0,
		},
		{
			name:    "both_set_negative_spread",
			bid:     types.Ptr(105.0),
			ask:     types.Ptr(100.0),
			wantNil: false,
			want:    -5.0,
		},
		{
			name:    "both_set_zero_spread",
			bid:     types.Ptr(100.0),
			ask:     types.Ptr(100.0),
			wantNil: false,
			want:    0.0,
		},
		{
			name:    "bid_nil",
			bid:     nil,
			ask:     types.Ptr(105.0),
			wantNil: true,
		},
		{
			name:    "ask_nil",
			bid:     types.Ptr(100.0),
			ask:     nil,
			wantNil: true,
		},
		{
			name:    "both_nil",
			bid:     nil,
			ask:     nil,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticker := &types.Ticker{
				BestBidPrice: tt.bid,
				BestAskPrice: tt.ask,
			}
			got := ticker.Spread()
			if tt.wantNil {
				if got != nil {
					t.Errorf("Ticker.Spread() = %v, want nil", *got)
				}
				return
			}
			if got == nil {
				t.Fatal("Ticker.Spread() = nil, want non-nil")
			}
			if *got != tt.want {
				t.Errorf("Ticker.Spread() = %v, want %v", *got, tt.want)
			}
		})
	}
}

// ---------- Ticker.MidPrice ----------

func TestTicker_MidPrice(t *testing.T) {
	tests := []struct {
		name    string
		bid     *float64
		ask     *float64
		wantNil bool
		want    float64
	}{
		{
			name:    "both_set",
			bid:     types.Ptr(100.0),
			ask:     types.Ptr(110.0),
			wantNil: false,
			want:    105.0,
		},
		{
			name:    "both_set_equal",
			bid:     types.Ptr(100.0),
			ask:     types.Ptr(100.0),
			wantNil: false,
			want:    100.0,
		},
		{
			name:    "both_set_small_spread",
			bid:     types.Ptr(99.5),
			ask:     types.Ptr(100.5),
			wantNil: false,
			want:    100.0,
		},
		{
			name:    "bid_nil",
			bid:     nil,
			ask:     types.Ptr(110.0),
			wantNil: true,
		},
		{
			name:    "ask_nil",
			bid:     types.Ptr(100.0),
			ask:     nil,
			wantNil: true,
		},
		{
			name:    "both_nil",
			bid:     nil,
			ask:     nil,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticker := &types.Ticker{
				BestBidPrice: tt.bid,
				BestAskPrice: tt.ask,
			}
			got := ticker.MidPrice()
			if tt.wantNil {
				if got != nil {
					t.Errorf("Ticker.MidPrice() = %v, want nil", *got)
				}
				return
			}
			if got == nil {
				t.Fatal("Ticker.MidPrice() = nil, want non-nil")
			}
			if *got != tt.want {
				t.Errorf("Ticker.MidPrice() = %v, want %v", *got, tt.want)
			}
		})
	}
}

// ---------- BookLevel accessors ----------

func TestBookLevel_Accessors(t *testing.T) {
	tests := []struct {
		name           string
		level          types.BookLevel
		wantPrice      float64
		wantAmount     float64
		wantOutrightAm float64
	}{
		{
			name:           "typical_level",
			level:          types.BookLevel{50000.0, 1.5, 1.0},
			wantPrice:      50000.0,
			wantAmount:     1.5,
			wantOutrightAm: 1.0,
		},
		{
			name:           "zero_values",
			level:          types.BookLevel{0, 0, 0},
			wantPrice:      0,
			wantAmount:     0,
			wantOutrightAm: 0,
		},
		{
			name:           "fractional_values",
			level:          types.BookLevel{0.001, 0.005, 0.003},
			wantPrice:      0.001,
			wantAmount:     0.005,
			wantOutrightAm: 0.003,
		},
		{
			name:           "large_values",
			level:          types.BookLevel{100000.0, 999.99, 500.0},
			wantPrice:      100000.0,
			wantAmount:     999.99,
			wantOutrightAm: 500.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Price(); got != tt.wantPrice {
				t.Errorf("BookLevel.Price() = %v, want %v", got, tt.wantPrice)
			}
			if got := tt.level.Amount(); got != tt.wantAmount {
				t.Errorf("BookLevel.Amount() = %v, want %v", got, tt.wantAmount)
			}
			if got := tt.level.OutrightAmount(); got != tt.wantOutrightAm {
				t.Errorf("BookLevel.OutrightAmount() = %v, want %v", got, tt.wantOutrightAm)
			}
		})
	}
}

// ---------- Instrument JSON round-trip ----------

func TestInstrument_JSONRoundTrip(t *testing.T) {
	strikePrice := 100000.0
	expTimestamp := int64(1711612800)
	settlementPrice := 50500.0
	settlementIndex := 50400.0

	inst := types.Instrument{
		InstrumentName:       "BTC-28MAR25-100000-C",
		Product:              "OBTCUSD",
		TickSize:             0.0005,
		VolumeTickSize:       0.1,
		MinOrderAmount:       0.1,
		Underlying:           "BTCUSD",
		Type:                 enums.InstrumentTypeOption,
		OptionType:           enums.OptionTypeCall,
		ExpiryDate:           "2025-03-28",
		ExpirationTimestamp:  &expTimestamp,
		StrikePrice:          &strikePrice,
		BaseCurrency:         "BTC",
		CreateTime:           1700000000.0,
		SettlementPrice:      &settlementPrice,
		SettlementIndexPrice: &settlementIndex,
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Instrument
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != inst.InstrumentName {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, inst.InstrumentName)
	}
	if got.Type != enums.InstrumentTypeOption {
		t.Errorf("Type = %q, want %q", got.Type, enums.InstrumentTypeOption)
	}
	if got.OptionType != enums.OptionTypeCall {
		t.Errorf("OptionType = %q, want %q", got.OptionType, enums.OptionTypeCall)
	}
	if got.StrikePrice == nil || *got.StrikePrice != strikePrice {
		t.Errorf("StrikePrice = %v, want %v", got.StrikePrice, strikePrice)
	}
	if got.ExpirationTimestamp == nil || *got.ExpirationTimestamp != expTimestamp {
		t.Errorf("ExpirationTimestamp = %v, want %v", got.ExpirationTimestamp, expTimestamp)
	}
	if got.SettlementPrice == nil || *got.SettlementPrice != settlementPrice {
		t.Errorf("SettlementPrice = %v, want %v", got.SettlementPrice, settlementPrice)
	}
}

func TestInstrument_CombinationWithLegs(t *testing.T) {
	inst := types.Instrument{
		InstrumentName: "BTC-CS",
		Type:           enums.InstrumentTypeCombination,
		Legs: []types.Leg{
			{InstrumentName: "BTC-28MAR25-100000-C", Quantity: 1.0},
			{InstrumentName: "BTC-28MAR25-110000-C", Quantity: -1.0},
		},
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Instrument
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !got.IsCombination() {
		t.Error("expected IsCombination() to be true")
	}
	if len(got.Legs) != 2 {
		t.Fatalf("len(Legs) = %d, want 2", len(got.Legs))
	}
	if got.Legs[0].InstrumentName != "BTC-28MAR25-100000-C" {
		t.Errorf("Legs[0].InstrumentName = %q, want %q", got.Legs[0].InstrumentName, "BTC-28MAR25-100000-C")
	}
}

func TestInstrument_PerpetualMinimal(t *testing.T) {
	inst := types.Instrument{
		InstrumentName: "BTC-PERPETUAL",
		Type:           enums.InstrumentTypePerpetual,
		Underlying:     "BTCUSD",
		TickSize:       0.5,
		VolumeTickSize: 1.0,
		MinOrderAmount: 1.0,
	}

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Optional fields should be absent
	for _, key := range []string{"strike_price", "expiration_timestamp", "settlement_price", "settlement_index_price"} {
		if _, ok := raw[key]; ok {
			t.Errorf("expected key %q to be omitted for perpetual", key)
		}
	}
}

// ---------- Ticker JSON round-trip ----------

func TestTicker_JSONRoundTrip(t *testing.T) {
	bid := 49000.0
	bidAmt := 5.0
	ask := 51000.0
	askAmt := 3.0
	last := 50000.0
	iv := 0.65
	delta := 0.55
	idx := 50000.0
	fwd := 50100.0
	vol24h := 1000.0
	val24h := 50000000.0
	low24h := 48000.0
	high24h := 52000.0
	change24h := 0.05
	collarLow := 45000.0
	collarHigh := 55000.0
	oi := 500.0
	fr := 0.0001
	fm := 0.0002

	ticker := types.Ticker{
		BestBidPrice:  &bid,
		BestBidAmount: &bidAmt,
		BestAskPrice:  &ask,
		BestAskAmount: &askAmt,
		LastPrice:     &last,
		MarkPrice:     50050.0,
		MarkTimestamp: 1700000000.0,
		IV:            &iv,
		Delta:         &delta,
		Index:         &idx,
		Forward:       &fwd,
		Volume24h:     &vol24h,
		Value24h:      &val24h,
		LowPrice24h:   &low24h,
		HighPrice24h:  &high24h,
		Change24h:     &change24h,
		CollarLow:     &collarLow,
		CollarHigh:    &collarHigh,
		OpenInterest:  &oi,
		FundingRate:   &fr,
		FundingMark:   &fm,
	}

	data, err := json.Marshal(ticker)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Ticker
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.MarkPrice != ticker.MarkPrice {
		t.Errorf("MarkPrice = %v, want %v", got.MarkPrice, ticker.MarkPrice)
	}
	if got.BestBidPrice == nil || *got.BestBidPrice != bid {
		t.Errorf("BestBidPrice = %v, want %v", got.BestBidPrice, bid)
	}
	if got.IV == nil || *got.IV != iv {
		t.Errorf("IV = %v, want %v", got.IV, iv)
	}
	if got.FundingRate == nil || *got.FundingRate != fr {
		t.Errorf("FundingRate = %v, want %v", got.FundingRate, fr)
	}

	// Verify Spread and MidPrice on the round-tripped value
	spread := got.Spread()
	if spread == nil {
		t.Fatal("Spread() should not be nil")
	}
	if *spread != 2000.0 {
		t.Errorf("Spread() = %v, want 2000.0", *spread)
	}

	mid := got.MidPrice()
	if mid == nil {
		t.Fatal("MidPrice() should not be nil")
	}
	if *mid != 50000.0 {
		t.Errorf("MidPrice() = %v, want 50000.0", *mid)
	}
}

// ---------- IndexPrice JSON round-trip ----------

func TestIndexPrice_JSONRoundTrip(t *testing.T) {
	expAvg := 50000.0
	expProgress := 0.5
	expectedExpPrice := 50100.0
	prevSettlement := 49800.0

	ip := types.IndexPrice{
		IndexName:               "BTCUSD",
		Price:                   50000.0,
		Timestamp:               1700000000.0,
		ExpirationPrintAverage:  &expAvg,
		ExpirationProgress:      &expProgress,
		ExpectedExpirationPrice: &expectedExpPrice,
		PreviousSettlementPrice: &prevSettlement,
	}

	data, err := json.Marshal(ip)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.IndexPrice
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.IndexName != "BTCUSD" {
		t.Errorf("IndexName = %q, want %q", got.IndexName, "BTCUSD")
	}
	if got.Price != 50000.0 {
		t.Errorf("Price = %v, want 50000.0", got.Price)
	}
	if got.ExpirationPrintAverage == nil || *got.ExpirationPrintAverage != expAvg {
		t.Errorf("ExpirationPrintAverage = %v, want %v", got.ExpirationPrintAverage, expAvg)
	}
	if got.PreviousSettlementPrice == nil || *got.PreviousSettlementPrice != prevSettlement {
		t.Errorf("PreviousSettlementPrice = %v, want %v", got.PreviousSettlementPrice, prevSettlement)
	}
}

// ---------- Book JSON round-trip ----------

func TestBook_JSONRoundTrip(t *testing.T) {
	lastPrice := 50000.0
	book := types.Book{
		Bids: []types.BookLevel{
			{49000.0, 2.0, 1.5},
			{48900.0, 3.0, 2.0},
		},
		Asks: []types.BookLevel{
			{51000.0, 1.0, 0.8},
		},
		Last: &lastPrice,
		Time: 1700000000.0,
	}

	data, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Book
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Bids) != 2 {
		t.Fatalf("len(Bids) = %d, want 2", len(got.Bids))
	}
	if len(got.Asks) != 1 {
		t.Fatalf("len(Asks) = %d, want 1", len(got.Asks))
	}
	if got.Bids[0].Price() != 49000.0 {
		t.Errorf("Bids[0].Price() = %v, want 49000.0", got.Bids[0].Price())
	}
	if got.Asks[0].Amount() != 1.0 {
		t.Errorf("Asks[0].Amount() = %v, want 1.0", got.Asks[0].Amount())
	}
	if got.Last == nil || *got.Last != lastPrice {
		t.Errorf("Last = %v, want %v", got.Last, lastPrice)
	}
}

func TestBook_EmptyBook(t *testing.T) {
	book := types.Book{
		Time: 1700000000.0,
	}

	data, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Book
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Last != nil {
		t.Errorf("Last = %v, want nil", got.Last)
	}
}

// ---------- BookLevel JSON round-trip ----------

func TestBookLevel_JSONRoundTrip(t *testing.T) {
	level := types.BookLevel{50000.0, 1.5, 1.0}

	data, err := json.Marshal(level)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.BookLevel
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Price() != 50000.0 {
		t.Errorf("Price() = %v, want 50000.0", got.Price())
	}
	if got.Amount() != 1.5 {
		t.Errorf("Amount() = %v, want 1.5", got.Amount())
	}
	if got.OutrightAmount() != 1.0 {
		t.Errorf("OutrightAmount() = %v, want 1.0", got.OutrightAmount())
	}
}
