package rest

import (
	"context"
	"net/http"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

func TestTradeHistory_Success(t *testing.T) {
	trades := []types.Trade{
		{
			TradeID:        "trade-001",
			OrderID:        "ord-001",
			InstrumentName: "BTC-PERPETUAL",
			Direction:      enums.DirectionBuy,
			Price:          50000.0,
			Amount:         1.0,
			Time:           1700000000.0,
			Fee:            0.5,
			FeeRate:        0.0005,
		},
		{
			TradeID:        "trade-002",
			OrderID:        "ord-002",
			InstrumentName: "ETH-PERPETUAL",
			Direction:      enums.DirectionSell,
			Price:          3000.0,
			Amount:         10.0,
			Time:           1700000001.0,
			Fee:            0.3,
			FeeRate:        0.0003,
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/private/trade_history" {
			t.Errorf("expected path /private/trade_history, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, trades))
	})

	result, err := c.TradeHistory(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 trades, got %d", len(result))
	}
	if result[0].TradeID != "trade-001" {
		t.Errorf("expected TradeID=trade-001, got %s", result[0].TradeID)
	}
	if result[0].Price != 50000.0 {
		t.Errorf("expected price=50000, got %f", result[0].Price)
	}
}

func TestTradeHistory_WithParams(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("from") != "1700000000" {
			t.Errorf("expected from=1700000000, got %s", q.Get("from"))
		}
		if q.Get("to") != "1700100000" {
			t.Errorf("expected to=1700100000, got %s", q.Get("to"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("offset") != "10" {
			t.Errorf("expected offset=10, got %s", q.Get("offset"))
		}
		if q.Get("sort") != "desc" {
			t.Errorf("expected sort=desc, got %s", q.Get("sort"))
		}
		if q.Get("instrument_names") != "BTC-PERPETUAL,ETH-PERPETUAL" {
			t.Errorf("expected instrument_names=BTC-PERPETUAL,ETH-PERPETUAL, got %s", q.Get("instrument_names"))
		}
		w.Write(wrapResult(t, []types.Trade{}))
	})

	from := 1700000000.0
	to := 1700100000.0
	limit := 50
	offset := 10
	params := &types.TradeHistoryParams{
		From:            &from,
		To:              &to,
		Limit:           &limit,
		Offset:          &offset,
		Sort:            enums.SortDesc,
		InstrumentNames: []string{"BTC-PERPETUAL", "ETH-PERPETUAL"},
	}

	_, err := c.TradeHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTradeHistory_WithBotIDs(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("bot_ids") != "bot-001,bot-002" {
			t.Errorf("expected bot_ids=bot-001,bot-002, got %s", r.URL.Query().Get("bot_ids"))
		}
		w.Write(wrapResult(t, []types.Trade{}))
	})

	params := &types.TradeHistoryParams{
		BotIDs: []string{"bot-001", "bot-002"},
	}

	_, err := c.TradeHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTradeHistory_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10002, "invalid parameter"))
	})

	_, err := c.TradeHistory(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOrderHistory_Success(t *testing.T) {
	price := 50000.0
	orders := []types.OrderHistory{
		{
			OrderID:        "ord-hist-001",
			OrderType:      enums.OrderTypeLimit,
			InstrumentName: "BTC-PERPETUAL",
			Direction:      enums.DirectionBuy,
			Price:          &price,
			Amount:         1.0,
			FilledAmount:   1.0,
			Status:         enums.OrderStatusFilled,
			CreateTime:     1700000000.0,
			CloseTime:      1700000010.0,
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/order_history" {
			t.Errorf("expected path /private/order_history, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, orders))
	})

	result, err := c.OrderHistory(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 order, got %d", len(result))
	}
	if result[0].OrderID != "ord-hist-001" {
		t.Errorf("expected OrderID=ord-hist-001, got %s", result[0].OrderID)
	}
	if result[0].Status != enums.OrderStatusFilled {
		t.Errorf("expected status=filled, got %s", result[0].Status)
	}
}

func TestOrderHistory_WithParams(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("sort") != "asc" {
			t.Errorf("expected sort=asc, got %s", q.Get("sort"))
		}
		if q.Get("instrument_names") != "BTC-PERPETUAL" {
			t.Errorf("expected instrument_names=BTC-PERPETUAL, got %s", q.Get("instrument_names"))
		}
		w.Write(wrapResult(t, []types.OrderHistory{}))
	})

	params := &types.OrderHistoryParams{
		Sort:            enums.SortAsc,
		InstrumentNames: []string{"BTC-PERPETUAL"},
	}

	_, err := c.OrderHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDailyMarkHistory_Success(t *testing.T) {
	marks := []types.DailyMark{
		{
			Time:                1700000000.0,
			InstrumentName:      "BTC-PERPETUAL",
			Position:            1.0,
			MarkPrice:           50000.0,
			RealizedPositionPnl: 100.0,
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/daily_mark_history" {
			t.Errorf("expected path /private/daily_mark_history, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, marks))
	})

	result, err := c.DailyMarkHistory(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("expected BTC-PERPETUAL, got %s", result[0].InstrumentName)
	}
}

func TestDailyMarkHistory_WithParams(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("from") != "1700000000" {
			t.Errorf("expected from=1700000000, got %s", q.Get("from"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		w.Write(wrapResult(t, []types.DailyMark{}))
	})

	from := 1700000000.0
	limit := 100
	params := &types.DailyMarkHistoryParams{
		From:  &from,
		Limit: &limit,
	}

	_, err := c.DailyMarkHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransactionHistory_Success(t *testing.T) {
	txns := []types.Transaction{
		{
			TransactionID: "txn-001",
			Time:          1700000000.0,
			Type:          "trade",
			Amount:        0.5,
			Currency:      "BTC",
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/transaction_history" {
			t.Errorf("expected path /private/transaction_history, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, txns))
	})

	result, err := c.TransactionHistory(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(result))
	}
	if result[0].TransactionID != "txn-001" {
		t.Errorf("expected TransactionID=txn-001, got %s", result[0].TransactionID)
	}
}

func TestTransactionHistory_WithParams(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("sort") != "desc" {
			t.Errorf("expected sort=desc, got %s", q.Get("sort"))
		}
		w.Write(wrapResult(t, []types.Transaction{}))
	})

	params := &types.TransactionHistoryParams{
		Sort: enums.SortDesc,
	}

	_, err := c.TransactionHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
