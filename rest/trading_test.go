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

func TestInsert_Success(t *testing.T) {
	expected := types.OrderStatus{
		OrderID:        "ord-001",
		OrderType:      enums.OrderTypeLimit,
		TimeInForce:    enums.TimeInForceGoodTillCancelled,
		InstrumentName: "BTC-PERPETUAL",
		Direction:      enums.DirectionBuy,
		Amount:         1.0,
		FilledAmount:   0,
		Status:         enums.OrderStatusOpen,
		ChangeReason:   enums.ChangeReasonInsert,
		InsertReason:   enums.InsertReasonClientRequest,
		CreateTime:     1700000000.0,
		Persistent:     true,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/private/insert" {
			t.Errorf("expected path /private/insert, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token-123" {
			t.Errorf("expected auth header, got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Thalex-Account") != "ACC-001" {
			t.Errorf("expected account header ACC-001, got %q", r.Header.Get("X-Thalex-Account"))
		}

		body, _ := io.ReadAll(r.Body)
		var params types.InsertOrderParams
		if err := json.Unmarshal(body, &params); err != nil {
			t.Errorf("failed to unmarshal request body: %v", err)
		}
		if params.Direction != enums.DirectionBuy {
			t.Errorf("expected direction buy, got %s", params.Direction)
		}
		if params.InstrumentName != "BTC-PERPETUAL" {
			t.Errorf("expected instrument BTC-PERPETUAL, got %s", params.InstrumentName)
		}

		w.Write(wrapResult(t, expected))
	})

	params := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0).
		WithPrice(50000.0).
		WithOrderType(enums.OrderTypeLimit).
		WithTimeInForce(enums.TimeInForceGoodTillCancelled)

	result, err := c.Insert(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-001" {
		t.Errorf("expected OrderID=ord-001, got %s", result.OrderID)
	}
	if result.Direction != enums.DirectionBuy {
		t.Errorf("expected direction buy, got %s", result.Direction)
	}
	if result.Status != enums.OrderStatusOpen {
		t.Errorf("expected status open, got %s", result.Status)
	}
	if result.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected instrument BTC-PERPETUAL, got %s", result.InstrumentName)
	}
}

func TestInsert_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10010, "insufficient margin"))
	})

	params := types.NewBuyOrderParams("BTC-PERPETUAL", 1.0).WithPrice(50000.0)
	_, err := c.Insert(context.Background(), params)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "API error 10010: insufficient margin" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuy_Success(t *testing.T) {
	expected := types.OrderStatus{
		OrderID:        "ord-buy-001",
		OrderType:      enums.OrderTypeMarket,
		InstrumentName: "BTC-PERPETUAL",
		Direction:      enums.DirectionBuy,
		Amount:         0.5,
		Status:         enums.OrderStatusFilled,
		FilledAmount:   0.5,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/buy" {
			t.Errorf("expected path /private/buy, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.Buy(context.Background(), "BTC-PERPETUAL", 0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-buy-001" {
		t.Errorf("expected OrderID=ord-buy-001, got %s", result.OrderID)
	}
	if result.Status != enums.OrderStatusFilled {
		t.Errorf("expected status filled, got %s", result.Status)
	}
}

func TestSell_Success(t *testing.T) {
	expected := types.OrderStatus{
		OrderID:        "ord-sell-001",
		OrderType:      enums.OrderTypeMarket,
		InstrumentName: "BTC-PERPETUAL",
		Direction:      enums.DirectionSell,
		Amount:         0.5,
		Status:         enums.OrderStatusFilled,
		FilledAmount:   0.5,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/sell" {
			t.Errorf("expected path /private/sell, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	result, err := c.Sell(context.Background(), "BTC-PERPETUAL", 0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-sell-001" {
		t.Errorf("expected OrderID=ord-sell-001, got %s", result.OrderID)
	}
	if result.Direction != enums.DirectionSell {
		t.Errorf("expected direction sell, got %s", result.Direction)
	}
}

func TestAmend_Success(t *testing.T) {
	expected := types.OrderStatus{
		OrderID:   "ord-001",
		OrderType: enums.OrderTypeLimit,
		Direction: enums.DirectionBuy,
		Amount:    2.0,
		Status:    enums.OrderStatusOpen,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/amend" {
			t.Errorf("expected path /private/amend, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params types.AmendOrderParams
		if err := json.Unmarshal(body, &params); err != nil {
			t.Errorf("failed to unmarshal: %v", err)
		}
		if params.OrderID != "ord-001" {
			t.Errorf("expected OrderID=ord-001, got %s", params.OrderID)
		}
		if params.Price != 51000.0 {
			t.Errorf("expected price=51000, got %f", params.Price)
		}
		if params.Amount != 2.0 {
			t.Errorf("expected amount=2.0, got %f", params.Amount)
		}
		w.Write(wrapResult(t, expected))
	})

	params := types.NewAmendByOrderID("ord-001", 51000.0, 2.0)
	result, err := c.Amend(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Amount != 2.0 {
		t.Errorf("expected amount=2.0, got %f", result.Amount)
	}
}

func TestAmend_ByClientOrderID(t *testing.T) {
	expected := types.OrderStatus{
		OrderID:   "ord-002",
		Direction: enums.DirectionSell,
		Amount:    3.0,
		Status:    enums.OrderStatusOpen,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var params types.AmendOrderParams
		json.Unmarshal(body, &params)
		if params.ClientOrderID == nil || *params.ClientOrderID != 42 {
			t.Errorf("expected client_order_id=42")
		}
		w.Write(wrapResult(t, expected))
	})

	params := types.NewAmendByClientOrderID(42, 52000.0, 3.0)
	result, err := c.Amend(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-002" {
		t.Errorf("expected OrderID=ord-002, got %s", result.OrderID)
	}
}

func TestCancel_ByOrderID(t *testing.T) {
	expected := types.OrderStatus{
		OrderID:      "ord-001",
		Direction:    enums.DirectionBuy,
		Status:       enums.OrderStatusCancelled,
		ChangeReason: enums.ChangeReasonCancel,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/cancel" {
			t.Errorf("expected path /private/cancel, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params types.CancelOrderParams
		json.Unmarshal(body, &params)
		if params.OrderID != "ord-001" {
			t.Errorf("expected order_id=ord-001, got %s", params.OrderID)
		}
		w.Write(wrapResult(t, expected))
	})

	params := types.CancelByOrderID("ord-001")
	result, err := c.Cancel(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != enums.OrderStatusCancelled {
		t.Errorf("expected status cancelled, got %s", result.Status)
	}
}

func TestCancel_ByClientOrderID(t *testing.T) {
	expected := types.OrderStatus{
		OrderID: "ord-003",
		Status:  enums.OrderStatusCancelled,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var params types.CancelOrderParams
		json.Unmarshal(body, &params)
		if params.ClientOrderID == nil || *params.ClientOrderID != 99 {
			t.Errorf("expected client_order_id=99")
		}
		w.Write(wrapResult(t, expected))
	})

	params := types.CancelByClientOrderID(99)
	result, err := c.Cancel(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "ord-003" {
		t.Errorf("expected OrderID=ord-003, got %s", result.OrderID)
	}
}

func TestCancelAll_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/cancel_all" {
			t.Errorf("expected path /private/cancel_all, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, map[string]int{"n_cancelled": 5}))
	})

	n, err := c.CancelAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 cancelled, got %d", n)
	}
}

func TestCancelAll_Zero(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, map[string]int{"n_cancelled": 0}))
	})

	n, err := c.CancelAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 cancelled, got %d", n)
	}
}

func TestCancelAll_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(apiErrorJSON(10001, "unauthorized"))
	})

	_, err := c.CancelAll(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOpenOrders_Success(t *testing.T) {
	orders := []types.OrderStatus{
		{OrderID: "ord-001", InstrumentName: "BTC-PERPETUAL", Direction: enums.DirectionBuy, Status: enums.OrderStatusOpen},
		{OrderID: "ord-002", InstrumentName: "BTC-PERPETUAL", Direction: enums.DirectionSell, Status: enums.OrderStatusOpen},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/private/open_orders" {
			t.Errorf("expected path /private/open_orders, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("instrument_name") != "BTC-PERPETUAL" {
			t.Errorf("expected instrument_name query param")
		}
		w.Write(wrapResult(t, map[string]any{"orders": orders}))
	})

	result, err := c.OpenOrders(context.Background(), "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(result))
	}
	if result[0].OrderID != "ord-001" {
		t.Errorf("expected first order ID=ord-001, got %s", result[0].OrderID)
	}
}

func TestOpenOrders_EmptyInstrument(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("instrument_name") != "" {
			t.Errorf("expected no instrument_name query param when empty")
		}
		w.Write(wrapResult(t, map[string]any{"orders": []types.OrderStatus{}}))
	})

	result, err := c.OpenOrders(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 orders, got %d", len(result))
	}
}

func TestOpenOrders_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10002, "bad request"))
	})

	_, err := c.OpenOrders(context.Background(), "BTC-PERPETUAL")
	if err == nil {
		t.Fatal("expected error")
	}
}
