package ws

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// ---------------------------------------------------------------------------
// Instruments (public)
// ---------------------------------------------------------------------------

func TestInstruments_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/instruments" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"instrument_name":"BTC-PERPETUAL","product":"btc","tick_size":0.5,"volume_tick_size":1,"min_order_amount":1,"underlying":"BTCUSD","type":"perpetual"}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instruments, err := c.Instruments(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instruments) != 1 {
		t.Fatalf("expected 1 instrument, got %d", len(instruments))
	}
	if instruments[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %q", instruments[0].InstrumentName)
	}
}

func TestInstruments_Empty(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`[]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instruments, err := c.Instruments(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instruments) != 0 {
		t.Errorf("expected 0 instruments, got %d", len(instruments))
	}
}

// ---------------------------------------------------------------------------
// AllInstruments
// ---------------------------------------------------------------------------

func TestAllInstruments_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/all_instruments" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"instrument_name":"BTC-PERPETUAL","product":"btc","tick_size":0.5,"volume_tick_size":1,"min_order_amount":1,"underlying":"BTCUSD","type":"perpetual"},{"instrument_name":"ETH-PERPETUAL","product":"eth","tick_size":0.05,"volume_tick_size":1,"min_order_amount":1,"underlying":"ETHUSD","type":"perpetual"}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	instruments, err := c.AllInstruments(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instruments) != 2 {
		t.Errorf("expected 2 instruments, got %d", len(instruments))
	}
}

// ---------------------------------------------------------------------------
// Instrument (single)
// ---------------------------------------------------------------------------

func TestInstrument_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/instrument" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"instrument_name":"BTC-PERPETUAL","product":"btc","tick_size":0.5,"volume_tick_size":1,"min_order_amount":1,"underlying":"BTCUSD","type":"perpetual"}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	inst, err := c.Instrument(ctx, "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inst.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %q", inst.InstrumentName)
	}
}

func TestInstrument_NotFound(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10002, Message: "instrument not found"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.Instrument(ctx, "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// Ticker (public)
// ---------------------------------------------------------------------------

func TestTicker_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/ticker" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"mark_price":50000.0,"mark_timestamp":1700000000.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker, err := c.Ticker(ctx, "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ticker.MarkPrice != 50000.0 {
		t.Errorf("expected mark_price=50000, got %f", ticker.MarkPrice)
	}
}

// ---------------------------------------------------------------------------
// Index
// ---------------------------------------------------------------------------

func TestIndex_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/index" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"index_name":"BTCUSD","price":49999.0,"timestamp":1700000000.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	index, err := c.Index(ctx, "BTCUSD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if index.IndexName != "BTCUSD" {
		t.Errorf("expected index_name=BTCUSD, got %q", index.IndexName)
	}
	if index.Price != 49999.0 {
		t.Errorf("expected price=49999, got %f", index.Price)
	}
}

// ---------------------------------------------------------------------------
// Book
// ---------------------------------------------------------------------------

func TestBook_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/book" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"bids":[[49999.0,1.0,1.0]],"asks":[[50001.0,2.0,2.0]],"time":1700000000.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	book, err := c.Book(ctx, "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(book.Bids) != 1 {
		t.Errorf("expected 1 bid, got %d", len(book.Bids))
	}
	if len(book.Asks) != 1 {
		t.Errorf("expected 1 ask, got %d", len(book.Asks))
	}
}

// ---------------------------------------------------------------------------
// SystemInfo
// ---------------------------------------------------------------------------

func TestSystemInfo_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/system_info" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"environment":"production","banners":[]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := c.SystemInfo(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Environment != "production" {
		t.Errorf("expected environment=production, got %q", info.Environment)
	}
}

// ---------------------------------------------------------------------------
// MarkPriceHistoricalData
// ---------------------------------------------------------------------------

func TestMarkPriceHistoricalData_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/mark_price_historical_data" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"instrument_type":"perpetual","mark":[[1.0,50000,50100,49900,50050,0.001]]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.MarkPriceHistoricalData(ctx, "BTC-PERPETUAL", 1.0, 2.0, enums.Resolution1h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.InstrumentType != "perpetual" {
		t.Errorf("expected instrument_type=perpetual, got %q", result.InstrumentType)
	}
}

// ---------------------------------------------------------------------------
// IndexPriceHistoricalData
// ---------------------------------------------------------------------------

func TestIndexPriceHistoricalData_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "public/index_price_historical_data" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"index":[[1.0,49900,50100,49800,50000]]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.IndexPriceHistoricalData(ctx, "BTCUSD", 1.0, 2.0, enums.Resolution1h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Index) != 1 {
		t.Errorf("expected 1 index data point, got %d", len(result.Index))
	}
}
