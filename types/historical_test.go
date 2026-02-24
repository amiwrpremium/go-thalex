package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

func TestMarkPriceHistoricalResult_PerpetualData(t *testing.T) {
	t.Run("parses perpetual data", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			InstrumentType: enums.InstrumentTypePerpetual,
			Mark: [][]any{
				{1700000000.0, 50000.0, 50100.0, 49900.0, 50050.0, 0.001},
				{1700003600.0, 50050.0, 50200.0, 49950.0, 50100.0, 0.002, []any{49990.0, 5.0, 50010.0, 3.0}},
			},
		}
		data := r.PerpetualData()
		if len(data) != 2 {
			t.Fatalf("expected 2 data points, got %d", len(data))
		}
		if data[0].Time != 1700000000.0 {
			t.Errorf("expected Time=1700000000, got %f", data[0].Time)
		}
		if data[0].Open != 50000.0 {
			t.Errorf("expected Open=50000, got %f", data[0].Open)
		}
		if data[0].FundingPayment != 0.001 {
			t.Errorf("expected FundingPayment=0.001, got %f", data[0].FundingPayment)
		}
		if data[0].TopOfBook != nil {
			t.Error("expected nil TopOfBook for first row")
		}
		if data[1].TopOfBook == nil {
			t.Fatal("expected non-nil TopOfBook for second row")
		}
		if data[1].TopOfBook.BidPrice == nil || *data[1].TopOfBook.BidPrice != 49990.0 {
			t.Errorf("expected BidPrice=49990, got %v", data[1].TopOfBook.BidPrice)
		}
	})

	t.Run("skips short rows", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{1.0, 2.0, 3.0},                // too short
				{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}, // ok
			},
		}
		data := r.PerpetualData()
		if len(data) != 1 {
			t.Fatalf("expected 1 data point (short row skipped), got %d", len(data))
		}
	})

	t.Run("empty mark", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{Mark: nil}
		data := r.PerpetualData()
		if len(data) != 0 {
			t.Fatalf("expected 0 data points, got %d", len(data))
		}
	})
}

func TestMarkPriceHistoricalResult_FutureData(t *testing.T) {
	t.Run("parses future data", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			InstrumentType: enums.InstrumentTypeFuture,
			Mark: [][]any{
				{1700000000.0, 50000.0, 50100.0, 49900.0, 50050.0},
				{1700003600.0, 50050.0, 50200.0, 49950.0, 50100.0, []any{49990.0, 5.0, 50010.0, 3.0}},
			},
		}
		data := r.FutureData()
		if len(data) != 2 {
			t.Fatalf("expected 2 data points, got %d", len(data))
		}
		if data[0].Close != 50050.0 {
			t.Errorf("expected Close=50050, got %f", data[0].Close)
		}
		if data[0].TopOfBook != nil {
			t.Error("expected nil TopOfBook for first row")
		}
		if data[1].TopOfBook == nil {
			t.Fatal("expected non-nil TopOfBook for second row")
		}
	})

	t.Run("skips short rows", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{{1.0, 2.0, 3.0}},
		}
		data := r.FutureData()
		if len(data) != 0 {
			t.Fatalf("expected 0 data points, got %d", len(data))
		}
	})
}

func TestMarkPriceHistoricalResult_OptionData(t *testing.T) {
	t.Run("parses option data", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			InstrumentType: enums.InstrumentTypeOption,
			Mark: [][]any{
				{1700000000.0, 0.05, 0.06, 0.04, 0.055, 0.5, 0.6, 0.4, 0.55},
				{1700003600.0, 0.055, 0.065, 0.045, 0.06, 0.55, 0.65, 0.45, 0.6, []any{0.04, 10.0, 0.06, 8.0}},
			},
		}
		data := r.OptionData()
		if len(data) != 2 {
			t.Fatalf("expected 2 data points, got %d", len(data))
		}
		if data[0].IVOpen != 0.5 {
			t.Errorf("expected IVOpen=0.5, got %f", data[0].IVOpen)
		}
		if data[0].IVClose != 0.55 {
			t.Errorf("expected IVClose=0.55, got %f", data[0].IVClose)
		}
		if data[0].TopOfBook != nil {
			t.Error("expected nil TopOfBook for first row")
		}
		if data[1].TopOfBook == nil {
			t.Fatal("expected non-nil TopOfBook for second row")
		}
	})

	t.Run("skips short rows", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}},
		}
		data := r.OptionData()
		if len(data) != 0 {
			t.Fatalf("expected 0 data points, got %d", len(data))
		}
	})
}

func TestIndexPriceHistoricalResult_Data(t *testing.T) {
	t.Run("parses index data", func(t *testing.T) {
		r := &types.IndexPriceHistoricalResult{
			Index: [][]any{
				{1700000000.0, 50000.0, 50100.0, 49900.0, 50050.0},
				{1700003600.0, 50050.0, 50200.0, 49950.0, 50100.0},
			},
		}
		data := r.Data()
		if len(data) != 2 {
			t.Fatalf("expected 2 data points, got %d", len(data))
		}
		if data[0].Time != 1700000000.0 {
			t.Errorf("expected Time=1700000000, got %f", data[0].Time)
		}
		if data[0].Open != 50000.0 {
			t.Errorf("expected Open=50000, got %f", data[0].Open)
		}
		if data[1].Close != 50100.0 {
			t.Errorf("expected Close=50100, got %f", data[1].Close)
		}
	})

	t.Run("skips short rows", func(t *testing.T) {
		r := &types.IndexPriceHistoricalResult{
			Index: [][]any{{1.0, 2.0}},
		}
		data := r.Data()
		if len(data) != 0 {
			t.Fatalf("expected 0 data points, got %d", len(data))
		}
	})

	t.Run("empty index", func(t *testing.T) {
		r := &types.IndexPriceHistoricalResult{Index: nil}
		data := r.Data()
		if len(data) != 0 {
			t.Fatalf("expected 0 data points, got %d", len(data))
		}
	})
}

// ---------- toFloat edge cases via PerpetualData ----------

func TestMarkPriceHistoricalResult_ToFloatEdgeCases(t *testing.T) {
	t.Run("int_values", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{int(100), int(200), int(300), int(400), int(500), int(1)},
			},
		}
		data := r.PerpetualData()
		if len(data) != 1 {
			t.Fatalf("expected 1 data point, got %d", len(data))
		}
		if data[0].Time != 100.0 {
			t.Errorf("Time = %v, want 100.0", data[0].Time)
		}
	})

	t.Run("int64_values", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{int64(100), int64(200), int64(300), int64(400), int64(500), int64(1)},
			},
		}
		data := r.PerpetualData()
		if len(data) != 1 {
			t.Fatalf("expected 1 data point, got %d", len(data))
		}
		if data[0].Time != 100.0 {
			t.Errorf("Time = %v, want 100.0", data[0].Time)
		}
	})

	t.Run("unsupported_type_returns_zero", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{"not_a_number", "a", "b", "c", "d", "e"},
			},
		}
		data := r.PerpetualData()
		if len(data) != 1 {
			t.Fatalf("expected 1 data point, got %d", len(data))
		}
		if data[0].Time != 0 {
			t.Errorf("Time = %v, want 0 for unsupported type", data[0].Time)
		}
	})
}

// ---------- parseTopOfBook edge cases via PerpetualData ----------

func TestMarkPriceHistoricalResult_ParseTopOfBookEdgeCases(t *testing.T) {
	t.Run("tob_not_array", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, "not_an_array"},
			},
		}
		data := r.PerpetualData()
		if len(data) != 1 {
			t.Fatalf("expected 1 data point, got %d", len(data))
		}
		if data[0].TopOfBook != nil {
			t.Error("expected nil TopOfBook when value is not an array")
		}
	})

	t.Run("tob_short_array", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, []any{1.0, 2.0}},
			},
		}
		data := r.PerpetualData()
		if len(data) != 1 {
			t.Fatalf("expected 1 data point, got %d", len(data))
		}
		if data[0].TopOfBook != nil {
			t.Error("expected nil TopOfBook when array is too short")
		}
	})

	t.Run("tob_with_nil_values", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, []any{nil, nil, nil, nil}},
			},
		}
		data := r.PerpetualData()
		if len(data) != 1 {
			t.Fatalf("expected 1 data point, got %d", len(data))
		}
		if data[0].TopOfBook == nil {
			t.Fatal("expected non-nil TopOfBook")
		}
		if data[0].TopOfBook.BidPrice != nil {
			t.Error("expected nil BidPrice")
		}
		if data[0].TopOfBook.BidSize != nil {
			t.Error("expected nil BidSize")
		}
		if data[0].TopOfBook.AskPrice != nil {
			t.Error("expected nil AskPrice")
		}
		if data[0].TopOfBook.AskSize != nil {
			t.Error("expected nil AskSize")
		}
	})

	t.Run("tob_mixed_nil_and_values", func(t *testing.T) {
		r := &types.MarkPriceHistoricalResult{
			Mark: [][]any{
				{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, []any{49990.0, nil, 50010.0, nil}},
			},
		}
		data := r.PerpetualData()
		if data[0].TopOfBook == nil {
			t.Fatal("expected non-nil TopOfBook")
		}
		if data[0].TopOfBook.BidPrice == nil || *data[0].TopOfBook.BidPrice != 49990.0 {
			t.Errorf("BidPrice = %v, want 49990.0", data[0].TopOfBook.BidPrice)
		}
		if data[0].TopOfBook.BidSize != nil {
			t.Error("expected nil BidSize")
		}
		if data[0].TopOfBook.AskPrice == nil || *data[0].TopOfBook.AskPrice != 50010.0 {
			t.Errorf("AskPrice = %v, want 50010.0", data[0].TopOfBook.AskPrice)
		}
		if data[0].TopOfBook.AskSize != nil {
			t.Error("expected nil AskSize")
		}
	})
}

// ---------- HistoricalDataParams JSON round-trip ----------

func TestHistoricalDataParams_JSONRoundTrip(t *testing.T) {
	p := types.HistoricalDataParams{
		InstrumentName: "BTC-PERPETUAL",
		From:           1700000000.0,
		To:             1700003600.0,
		Resolution:     enums.Resolution1h,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.HistoricalDataParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", got.InstrumentName, "BTC-PERPETUAL")
	}
	if got.Resolution != enums.Resolution1h {
		t.Errorf("Resolution = %q, want %q", got.Resolution, enums.Resolution1h)
	}
	if got.From != 1700000000.0 {
		t.Errorf("From = %v, want 1700000000.0", got.From)
	}
}

func TestHistoricalDataParams_IndexName(t *testing.T) {
	p := types.HistoricalDataParams{
		IndexName:  "BTCUSD",
		From:       1700000000.0,
		To:         1700003600.0,
		Resolution: enums.Resolution1d,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.HistoricalDataParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.IndexName != "BTCUSD" {
		t.Errorf("IndexName = %q, want %q", got.IndexName, "BTCUSD")
	}
}

// ---------- MarkPriceHistoricalResult JSON round-trip ----------

func TestMarkPriceHistoricalResult_JSONRoundTrip(t *testing.T) {
	r := types.MarkPriceHistoricalResult{
		InstrumentType: enums.InstrumentTypePerpetual,
		Mark:           [][]any{{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}},
		NoData:         false,
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.MarkPriceHistoricalResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.InstrumentType != enums.InstrumentTypePerpetual {
		t.Errorf("InstrumentType = %q, want %q", got.InstrumentType, enums.InstrumentTypePerpetual)
	}
	if len(got.Mark) != 1 {
		t.Fatalf("len(Mark) = %d, want 1", len(got.Mark))
	}
}

func TestMarkPriceHistoricalResult_NoData(t *testing.T) {
	r := types.MarkPriceHistoricalResult{
		InstrumentType: enums.InstrumentTypeFuture,
		NoData:         true,
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.MarkPriceHistoricalResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !got.NoData {
		t.Error("NoData = false, want true")
	}
	if got.Mark != nil {
		t.Errorf("Mark = %v, want nil", got.Mark)
	}

	// Calling data accessors on empty/nil should return empty slices
	if len(got.PerpetualData()) != 0 {
		t.Error("PerpetualData() should be empty for NoData result")
	}
	if len(got.FutureData()) != 0 {
		t.Error("FutureData() should be empty for NoData result")
	}
	if len(got.OptionData()) != 0 {
		t.Error("OptionData() should be empty for NoData result")
	}
}

// ---------- IndexPriceHistoricalResult JSON round-trip ----------

func TestIndexPriceHistoricalResult_JSONRoundTrip(t *testing.T) {
	r := types.IndexPriceHistoricalResult{
		Index:  [][]any{{1.0, 2.0, 3.0, 4.0, 5.0}},
		NoData: false,
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.IndexPriceHistoricalResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Index) != 1 {
		t.Fatalf("len(Index) = %d, want 1", len(got.Index))
	}
}

// ---------- TopOfBook JSON round-trip ----------

func TestTopOfBook_JSONRoundTrip(t *testing.T) {
	bid := 49990.0
	bidSize := 5.0
	ask := 50010.0
	askSize := 3.0

	tob := types.TopOfBook{
		BidPrice: &bid,
		BidSize:  &bidSize,
		AskPrice: &ask,
		AskSize:  &askSize,
	}

	data, err := json.Marshal(tob)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.TopOfBook
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.BidPrice == nil || *got.BidPrice != bid {
		t.Errorf("BidPrice = %v, want %v", got.BidPrice, bid)
	}
	if got.AskSize == nil || *got.AskSize != askSize {
		t.Errorf("AskSize = %v, want %v", got.AskSize, askSize)
	}
}

// ---------- OHLC JSON round-trip ----------

func TestOHLC_JSONRoundTrip(t *testing.T) {
	ohlc := types.OHLC{
		Time:  1700000000.0,
		Open:  50000.0,
		High:  50100.0,
		Low:   49900.0,
		Close: 50050.0,
	}

	data, err := json.Marshal(ohlc)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OHLC
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Open != 50000.0 {
		t.Errorf("Open = %v, want 50000.0", got.Open)
	}
	if got.High != 50100.0 {
		t.Errorf("High = %v, want 50100.0", got.High)
	}
}
