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

func TestConditionalOrders_Success(t *testing.T) {
	stopPrice := 48000.0
	orders := []types.ConditionalOrder{
		{
			OrderID:        "cond-001",
			InstrumentName: "BTC-PERPETUAL",
			Direction:      enums.DirectionSell,
			Amount:         1.0,
			Target:         enums.TargetMark,
			StopPrice:      stopPrice,
			Status:         enums.ConditionalOrderStatusActive,
			CreateTime:     1700000000.0,
			UpdateTime:     1700000000.0,
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/private/conditional_orders" {
			t.Errorf("expected path /private/conditional_orders, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, orders))
	})

	result, err := c.ConditionalOrders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 conditional order, got %d", len(result))
	}
	if result[0].OrderID != "cond-001" {
		t.Errorf("expected OrderID=cond-001, got %s", result[0].OrderID)
	}
	if result[0].Status != enums.ConditionalOrderStatusActive {
		t.Errorf("expected status=active, got %s", result[0].Status)
	}
	if result[0].Direction != enums.DirectionSell {
		t.Errorf("expected direction=sell, got %s", result[0].Direction)
	}
}

func TestConditionalOrders_Empty(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, []types.ConditionalOrder{}))
	})

	result, err := c.ConditionalOrders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}

func TestConditionalOrders_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(apiErrorJSON(10001, "unauthorized"))
	})

	_, err := c.ConditionalOrders(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateConditionalOrder_StopOrder(t *testing.T) {
	expected := types.ConditionalOrder{
		OrderID:        "cond-002",
		InstrumentName: "BTC-PERPETUAL",
		Direction:      enums.DirectionSell,
		Amount:         1.0,
		StopPrice:      48000.0,
		Status:         enums.ConditionalOrderStatusCreated,
		CreateTime:     1700000000.0,
		UpdateTime:     1700000000.0,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/create_conditional_order" {
			t.Errorf("expected path /private/create_conditional_order, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params types.CreateConditionalOrderParams
		json.Unmarshal(body, &params)
		if params.Direction != enums.DirectionSell {
			t.Errorf("expected direction=sell, got %s", params.Direction)
		}
		if params.StopPrice != 48000.0 {
			t.Errorf("expected stop_price=48000, got %f", params.StopPrice)
		}
		w.Write(wrapResult(t, expected))
	})

	params := types.NewStopOrder(enums.DirectionSell, "BTC-PERPETUAL", 1.0, 48000.0).
		WithTarget(enums.TargetMark).
		WithLabel("test-stop")
	result, err := c.CreateConditionalOrder(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OrderID != "cond-002" {
		t.Errorf("expected OrderID=cond-002, got %s", result.OrderID)
	}
}

func TestCreateConditionalOrder_StopLimit(t *testing.T) {
	limitPrice := 47500.0
	expected := types.ConditionalOrder{
		OrderID:        "cond-003",
		InstrumentName: "BTC-PERPETUAL",
		Direction:      enums.DirectionSell,
		Amount:         1.0,
		StopPrice:      48000.0,
		LimitPrice:     &limitPrice,
		Status:         enums.ConditionalOrderStatusCreated,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, expected))
	})

	params := types.NewStopLimitOrder(enums.DirectionSell, "BTC-PERPETUAL", 1.0, 48000.0, 47500.0)
	result, err := c.CreateConditionalOrder(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsStopLimit() {
		t.Error("expected IsStopLimit() to be true")
	}
}

func TestCreateConditionalOrder_Bracket(t *testing.T) {
	bracketPrice := 52000.0
	expected := types.ConditionalOrder{
		OrderID:      "cond-004",
		StopPrice:    48000.0,
		BracketPrice: &bracketPrice,
		Status:       enums.ConditionalOrderStatusCreated,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, expected))
	})

	params := types.NewBracketOrder(enums.DirectionSell, "BTC-PERPETUAL", 1.0, 48000.0, 52000.0)
	result, err := c.CreateConditionalOrder(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsBracket() {
		t.Error("expected IsBracket() to be true")
	}
}

func TestCreateConditionalOrder_TrailingStop(t *testing.T) {
	callbackRate := 0.05
	expected := types.ConditionalOrder{
		OrderID:                  "cond-005",
		StopPrice:                48000.0,
		TrailingStopCallbackRate: &callbackRate,
		Status:                   enums.ConditionalOrderStatusCreated,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, expected))
	})

	params := types.NewTrailingStopOrder(enums.DirectionSell, "BTC-PERPETUAL", 1.0, 48000.0, 0.05).
		WithReduceOnly(true)
	result, err := c.CreateConditionalOrder(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsTrailingStop() {
		t.Error("expected IsTrailingStop() to be true")
	}
}

func TestCreateConditionalOrder_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10015, "invalid stop price"))
	})

	params := types.NewStopOrder(enums.DirectionSell, "BTC-PERPETUAL", 1.0, 0)
	_, err := c.CreateConditionalOrder(context.Background(), params)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCancelConditionalOrder_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/cancel_conditional_order" {
			t.Errorf("expected path /private/cancel_conditional_order, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params struct {
			OrderID string `json:"order_id"`
		}
		json.Unmarshal(body, &params)
		if params.OrderID != "cond-001" {
			t.Errorf("expected order_id=cond-001, got %s", params.OrderID)
		}
		w.Write([]byte(`{"result":null}`))
	})

	err := c.CancelConditionalOrder(context.Background(), "cond-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelConditionalOrder_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10021, "conditional order not found"))
	})

	err := c.CancelConditionalOrder(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCancelAllConditionalOrders_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/cancel_all_conditional_orders" {
			t.Errorf("expected path /private/cancel_all_conditional_orders, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, map[string]int{"n_cancelled": 4}))
	})

	n, err := c.CancelAllConditionalOrders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 cancelled, got %d", n)
	}
}

func TestCancelAllConditionalOrders_Zero(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, map[string]int{"n_cancelled": 0}))
	})

	n, err := c.CancelAllConditionalOrders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}
