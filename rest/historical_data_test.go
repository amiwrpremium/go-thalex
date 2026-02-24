package rest

import (
	"context"
	"net/http"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

func TestMarkPriceHistoricalData_Success(t *testing.T) {
	expected := types.MarkPriceHistoricalResult{
		InstrumentType: enums.InstrumentTypePerpetual,
		Mark: [][]any{
			{1700000000.0, 50000.0, 50100.0, 49900.0, 50050.0, 0.001},
		},
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/public/mark_price_historical_data" {
			t.Errorf("expected path /public/mark_price_historical_data, got %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("instrument_name") != "BTC-PERPETUAL" {
			t.Errorf("expected instrument_name=BTC-PERPETUAL, got %s", q.Get("instrument_name"))
		}
		if q.Get("resolution") != string(enums.Resolution1h) {
			t.Errorf("expected resolution=%s, got %s", enums.Resolution1h, q.Get("resolution"))
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.MarkPriceHistoricalData(
		context.Background(),
		"BTC-PERPETUAL",
		1700000000.0,
		1700100000.0,
		enums.Resolution1h,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.InstrumentType != enums.InstrumentTypePerpetual {
		t.Errorf("expected instrument type perpetual, got %s", result.InstrumentType)
	}
	if len(result.Mark) != 1 {
		t.Fatalf("expected 1 data point, got %d", len(result.Mark))
	}
}

func TestIndexPriceHistoricalData_Success(t *testing.T) {
	expected := types.IndexPriceHistoricalResult{
		Index: [][]any{
			{1700000000.0, 50000.0, 50100.0, 49900.0, 50050.0},
		},
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/public/index_price_historical_data" {
			t.Errorf("expected path /public/index_price_historical_data, got %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("index_name") != "BTCUSD" {
			t.Errorf("expected index_name=BTCUSD, got %s", q.Get("index_name"))
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.IndexPriceHistoricalData(
		context.Background(),
		"BTCUSD",
		1700000000.0,
		1700100000.0,
		enums.Resolution1h,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Index) != 1 {
		t.Fatalf("expected 1 data point, got %d", len(result.Index))
	}
}
