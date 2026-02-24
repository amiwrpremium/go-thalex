package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

func TestCreateRfq_Success(t *testing.T) {
	expected := types.Rfq{
		RfqID:      "rfq-001",
		Amount:     1.0,
		CreateTime: 1700000000.0,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/private/create_rfq" {
			t.Errorf("expected path /private/create_rfq, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params types.CreateRfqParams
		if err := json.Unmarshal(body, &params); err != nil {
			t.Errorf("unmarshal body: %v", err)
		}
		if params.Amount != 1.0 {
			t.Errorf("expected amount 1.0, got %f", params.Amount)
		}
		w.Write(wrapResult(t, expected))
	})

	params := &types.CreateRfqParams{
		Legs:   []types.InsertLeg{{InstrumentName: "BTC-PERPETUAL", Quantity: 1}},
		Amount: 1.0,
		Label:  "test-rfq",
	}
	result, err := c.CreateRfq(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RfqID != "rfq-001" {
		t.Errorf("expected RfqID=rfq-001, got %s", result.RfqID)
	}
}

func TestCancelRfq_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/private/cancel_rfq" {
			t.Errorf("expected path /private/cancel_rfq, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params struct {
			RfqID string `json:"rfq_id"`
		}
		json.Unmarshal(body, &params)
		if params.RfqID != "rfq-001" {
			t.Errorf("expected rfq_id=rfq-001, got %s", params.RfqID)
		}
		w.Write([]byte(`{"result":null}`))
	})

	err := c.CancelRfq(context.Background(), "rfq-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTradeRfq_Success(t *testing.T) {
	expected := []types.Trade{{TradeID: "trade-001"}}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/trade_rfq" {
			t.Errorf("expected path /private/trade_rfq, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	params := &types.TradeRfqParams{
		RfqID:     "rfq-001",
		Direction: enums.DirectionBuy,
		Price:     50000.0,
		Amount:    1.0,
	}
	result, err := c.TradeRfq(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(result))
	}
	if result[0].TradeID != "trade-001" {
		t.Errorf("expected TradeID=trade-001, got %s", result[0].TradeID)
	}
}

func TestOpenRfqs_Success(t *testing.T) {
	expected := []types.Rfq{{RfqID: "rfq-001"}, {RfqID: "rfq-002"}}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/private/open_rfqs" {
			t.Errorf("expected path /private/open_rfqs, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.OpenRfqs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 rfqs, got %d", len(result))
	}
}

func TestRfqHistory_NoParams(t *testing.T) {
	expected := []types.Rfq{{RfqID: "rfq-hist-001"}}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/rfq_history" {
			t.Errorf("expected path /private/rfq_history, got %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %s", r.URL.RawQuery)
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.RfqHistory(context.Background(), nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 rfq, got %d", len(result))
	}
}

func TestRfqHistory_WithParams(t *testing.T) {
	expected := []types.Rfq{{RfqID: "rfq-hist-002"}}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("from") == "" {
			t.Error("expected from param")
		}
		if q.Get("to") == "" {
			t.Error("expected to param")
		}
		if q.Get("offset") != "0" {
			t.Errorf("expected offset=0, got %s", q.Get("offset"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		w.Write(wrapResult(t, expected))
	})

	from := 1700000000.0
	to := 1700100000.0
	offset := 0
	limit := 10
	result, err := c.RfqHistory(context.Background(), &from, &to, &offset, &limit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 rfq, got %d", len(result))
	}
}

func TestMMRfqs_Success(t *testing.T) {
	expected := []types.Rfq{{RfqID: "mm-rfq-001"}}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/mm_rfqs" {
			t.Errorf("expected path /private/mm_rfqs, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.MMRfqs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 rfq, got %d", len(result))
	}
}

func TestMMRfqInsertQuote_Success(t *testing.T) {
	expected := types.RfqOrder{RfqID: "rfq-001", OrderID: "ord-001"}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/mm_rfq_insert_quote" {
			t.Errorf("expected path /private/mm_rfq_insert_quote, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	params := &types.RfqQuoteInsertParams{
		RfqID:     "rfq-001",
		Direction: enums.DirectionBuy,
		Amount:    1.0,
		Price:     50000.0,
	}
	result, err := c.MMRfqInsertQuote(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-001" {
		t.Errorf("expected OrderID=ord-001, got %s", result.OrderID)
	}
}

func TestMMRfqAmendQuote_Success(t *testing.T) {
	expected := types.RfqOrder{RfqID: "rfq-001", OrderID: "ord-001"}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/mm_rfq_amend_quote" {
			t.Errorf("expected path /private/mm_rfq_amend_quote, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	params := &types.RfqQuoteAmendParams{
		OrderID: "ord-001",
		Amount:  2.0,
		Price:   51000.0,
	}
	result, err := c.MMRfqAmendQuote(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-001" {
		t.Errorf("expected OrderID=ord-001, got %s", result.OrderID)
	}
}

func TestMMRfqDeleteQuote_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/mm_rfq_delete_quote" {
			t.Errorf("expected path /private/mm_rfq_delete_quote, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"result":null}`))
	})

	params := &types.RfqQuoteDeleteParams{OrderID: "ord-001"}
	err := c.MMRfqDeleteQuote(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMMRfqQuotes_Success(t *testing.T) {
	expected := []types.RfqOrder{{RfqID: "rfq-001", OrderID: "ord-001"}}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/mm_rfq_quotes" {
			t.Errorf("expected path /private/mm_rfq_quotes, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.MMRfqQuotes(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(result))
	}
}
