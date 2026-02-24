package rest

import (
	"context"
	"net/http"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

func TestInstruments_Success(t *testing.T) {
	instruments := []types.Instrument{
		{
			InstrumentName: "BTC-PERPETUAL",
			Product:        "BTC",
			TickSize:       0.5,
			VolumeTickSize: 0.001,
			MinOrderAmount: 0.001,
			Underlying:     "BTCUSD",
			Type:           enums.InstrumentTypePerpetual,
		},
		{
			InstrumentName: "ETH-PERPETUAL",
			Product:        "ETH",
			TickSize:       0.01,
			VolumeTickSize: 0.01,
			MinOrderAmount: 0.01,
			Underlying:     "ETHUSD",
			Type:           enums.InstrumentTypePerpetual,
		},
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/public/instruments" {
			t.Errorf("expected path /public/instruments, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, instruments))
	})

	result, err := c.Instruments(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 instruments, got %d", len(result))
	}
	if result[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %s", result[0].InstrumentName)
	}
	if result[0].Type != enums.InstrumentTypePerpetual {
		t.Errorf("expected perpetual type, got %s", result[0].Type)
	}
	if result[1].Underlying != "ETHUSD" {
		t.Errorf("expected ETHUSD underlying, got %s", result[1].Underlying)
	}
}

func TestInstruments_Empty(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, []types.Instrument{}))
	})

	result, err := c.Instruments(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 instruments, got %d", len(result))
	}
}

func TestInstruments_APIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10050, "rate limit"))
	})

	_, err := c.Instruments(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAllInstruments_Success(t *testing.T) {
	instruments := []types.Instrument{
		{InstrumentName: "BTC-PERPETUAL", Type: enums.InstrumentTypePerpetual},
		{InstrumentName: "BTC-20250328-C-50000", Type: enums.InstrumentTypeOption},
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/all_instruments" {
			t.Errorf("expected path /public/all_instruments, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, instruments))
	})

	result, err := c.AllInstruments(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 instruments, got %d", len(result))
	}
}

func TestInstrument_Single(t *testing.T) {
	inst := types.Instrument{
		InstrumentName: "BTC-PERPETUAL",
		Product:        "BTC",
		TickSize:       0.5,
		Underlying:     "BTCUSD",
		Type:           enums.InstrumentTypePerpetual,
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/instrument" {
			t.Errorf("expected path /public/instrument, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("instrument_name") != "BTC-PERPETUAL" {
			t.Errorf("expected instrument_name=BTC-PERPETUAL")
		}
		w.Write(wrapResult(t, inst))
	})

	result, err := c.Instrument(context.Background(), "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %s", result.InstrumentName)
	}
}

func TestTicker_Success(t *testing.T) {
	bidPrice := 49999.5
	bidAmt := 10.0
	askPrice := 50000.5
	askAmt := 5.0
	lastPrice := 50000.0
	ticker := types.Ticker{
		BestBidPrice:  &bidPrice,
		BestBidAmount: &bidAmt,
		BestAskPrice:  &askPrice,
		BestAskAmount: &askAmt,
		LastPrice:     &lastPrice,
		MarkPrice:     50000.0,
		MarkTimestamp: 1700000000.0,
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/ticker" {
			t.Errorf("expected path /public/ticker, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("instrument_name") != "BTC-PERPETUAL" {
			t.Errorf("expected instrument_name=BTC-PERPETUAL")
		}
		w.Write(wrapResult(t, ticker))
	})

	result, err := c.Ticker(context.Background(), "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.MarkPrice != 50000.0 {
		t.Errorf("expected mark_price=50000, got %f", result.MarkPrice)
	}
	if result.BestBidPrice == nil || *result.BestBidPrice != 49999.5 {
		t.Errorf("expected best_bid_price=49999.5")
	}
	if result.BestAskPrice == nil || *result.BestAskPrice != 50000.5 {
		t.Errorf("expected best_ask_price=50000.5")
	}

	// Test spread helper
	spread := result.Spread()
	if spread == nil {
		t.Fatal("expected non-nil spread")
	}
	if *spread != 1.0 {
		t.Errorf("expected spread=1.0, got %f", *spread)
	}

	// Test mid price helper
	mid := result.MidPrice()
	if mid == nil {
		t.Fatal("expected non-nil mid price")
	}
	if *mid != 50000.0 {
		t.Errorf("expected mid=50000.0, got %f", *mid)
	}
}

func TestTicker_APIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10003, "instrument not found"))
	})

	_, err := c.Ticker(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestIndex_Success(t *testing.T) {
	index := types.IndexPrice{
		IndexName: "BTCUSD",
		Price:     50000.0,
		Timestamp: 1700000000.0,
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/index" {
			t.Errorf("expected path /public/index, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("underlying") != "BTCUSD" {
			t.Errorf("expected underlying=BTCUSD")
		}
		w.Write(wrapResult(t, index))
	})

	result, err := c.Index(context.Background(), "BTCUSD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IndexName != "BTCUSD" {
		t.Errorf("expected IndexName=BTCUSD, got %s", result.IndexName)
	}
	if result.Price != 50000.0 {
		t.Errorf("expected price=50000, got %f", result.Price)
	}
}

func TestBook_Success(t *testing.T) {
	last := 50000.0
	book := types.Book{
		Bids: []types.BookLevel{
			{49999.0, 5.0, 3.0},
			{49998.0, 10.0, 8.0},
		},
		Asks: []types.BookLevel{
			{50001.0, 3.0, 2.0},
			{50002.0, 7.0, 5.0},
		},
		Last: &last,
		Time: 1700000000.0,
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/book" {
			t.Errorf("expected path /public/book, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("instrument_name") != "BTC-PERPETUAL" {
			t.Errorf("expected instrument_name=BTC-PERPETUAL")
		}
		w.Write(wrapResult(t, book))
	})

	result, err := c.Book(context.Background(), "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Bids) != 2 {
		t.Fatalf("expected 2 bid levels, got %d", len(result.Bids))
	}
	if len(result.Asks) != 2 {
		t.Fatalf("expected 2 ask levels, got %d", len(result.Asks))
	}
	if result.Bids[0].Price() != 49999.0 {
		t.Errorf("expected first bid price=49999, got %f", result.Bids[0].Price())
	}
	if result.Bids[0].Amount() != 5.0 {
		t.Errorf("expected first bid amount=5.0, got %f", result.Bids[0].Amount())
	}
	if result.Bids[0].OutrightAmount() != 3.0 {
		t.Errorf("expected first bid outright_amount=3.0, got %f", result.Bids[0].OutrightAmount())
	}
	if result.Asks[0].Price() != 50001.0 {
		t.Errorf("expected first ask price=50001, got %f", result.Asks[0].Price())
	}
	if result.Last == nil || *result.Last != 50000.0 {
		t.Errorf("expected last=50000")
	}
}

func TestBook_Empty(t *testing.T) {
	book := types.Book{
		Bids: []types.BookLevel{},
		Asks: []types.BookLevel{},
		Time: 1700000000.0,
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, book))
	})

	result, err := c.Book(context.Background(), "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Bids) != 0 {
		t.Errorf("expected 0 bids, got %d", len(result.Bids))
	}
	if result.Last != nil {
		t.Errorf("expected nil last")
	}
}

func TestSystemInfo_Success(t *testing.T) {
	info := types.SystemInfo{
		Environment: "testnet",
		APIVersion:  "v2",
		Banners:     []types.Banner{},
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/system_info" {
			t.Errorf("expected path /public/system_info, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, info))
	})

	result, err := c.SystemInfo(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Environment != "testnet" {
		t.Errorf("expected environment=testnet, got %s", result.Environment)
	}
}
