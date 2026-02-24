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

func TestBots_Success(t *testing.T) {
	bots := []types.Bot{
		{
			BotID:          "bot-001",
			Strategy:       enums.BotStrategySGSL,
			Status:         enums.BotStatusActive,
			InstrumentName: "BTC-PERPETUAL",
			StartTime:      1700000000.0,
		},
		{
			BotID:          "bot-002",
			Strategy:       enums.BotStrategyGrid,
			Status:         enums.BotStatusStopped,
			InstrumentName: "ETH-PERPETUAL",
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/bots" {
			t.Errorf("expected path /private/bots, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("include_inactive") != "true" {
			t.Errorf("expected include_inactive=true")
		}
		w.Write(wrapResult(t, bots))
	})

	result, err := c.Bots(context.Background(), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 bots, got %d", len(result))
	}
	if result[0].BotID != "bot-001" {
		t.Errorf("expected BotID=bot-001, got %s", result[0].BotID)
	}
	if result[0].Strategy != enums.BotStrategySGSL {
		t.Errorf("expected strategy=sgsl, got %s", result[0].Strategy)
	}
	if result[0].Status != enums.BotStatusActive {
		t.Errorf("expected status=active, got %s", result[0].Status)
	}
}

func TestBots_ActiveOnly(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("include_inactive") != "" {
			t.Errorf("expected no include_inactive param when false")
		}
		w.Write(wrapResult(t, []types.Bot{}))
	})

	_, err := c.Bots(context.Background(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBots_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(apiErrorJSON(10001, "unauthorized"))
	})

	_, err := c.Bots(context.Background(), false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateSGSLBot_Success(t *testing.T) {
	expected := types.Bot{
		BotID:          "bot-sgsl-001",
		Strategy:       enums.BotStrategySGSL,
		Status:         enums.BotStatusActive,
		InstrumentName: "BTC-PERPETUAL",
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/create_bot" {
			t.Errorf("expected path /private/create_bot, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params types.SGSLBotParams
		json.Unmarshal(body, &params)
		if params.Strategy != enums.BotStrategySGSL {
			t.Errorf("expected strategy=sgsl, got %s", params.Strategy)
		}
		if params.InstrumentName != "BTC-PERPETUAL" {
			t.Errorf("expected instrument=BTC-PERPETUAL, got %s", params.InstrumentName)
		}
		w.Write(wrapResult(t, expected))
	})

	params := types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetMark, 49000, 1.0, 51000, 0, 1700100000).
		WithMaxSlippage(100).
		WithLabel("test-sgsl")
	result, err := c.CreateSGSLBot(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-sgsl-001" {
		t.Errorf("expected BotID=bot-sgsl-001, got %s", result.BotID)
	}
}

func TestCreateGridBot_Success(t *testing.T) {
	expected := types.Bot{
		BotID:          "bot-grid-001",
		Strategy:       enums.BotStrategyGrid,
		Status:         enums.BotStatusActive,
		InstrumentName: "BTC-PERPETUAL",
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var params types.GridBotParams
		json.Unmarshal(body, &params)
		if params.Strategy != enums.BotStrategyGrid {
			t.Errorf("expected strategy=grid, got %s", params.Strategy)
		}
		w.Write(wrapResult(t, expected))
	})

	params := types.NewGridBotParams("BTC-PERPETUAL", []float64{49000, 50000, 51000}, 0.1, 1700100000).
		WithBasePosition(0).
		WithLabel("test-grid")
	result, err := c.CreateGridBot(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Strategy != enums.BotStrategyGrid {
		t.Errorf("expected strategy=grid, got %s", result.Strategy)
	}
}

func TestCreateOCQBot_Success(t *testing.T) {
	expected := types.Bot{
		BotID:    "bot-ocq-001",
		Strategy: enums.BotStrategyOCQ,
		Status:   enums.BotStatusActive,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, expected))
	})

	params := types.NewOCQBotParams("BTC-PERPETUAL", enums.TargetMark, 10, 10, 0.1, -1, 1, 1700100000)
	result, err := c.CreateOCQBot(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-ocq-001" {
		t.Errorf("expected BotID=bot-ocq-001, got %s", result.BotID)
	}
}

func TestCreateLevelsBot_Success(t *testing.T) {
	expected := types.Bot{
		BotID:    "bot-levels-001",
		Strategy: enums.BotStrategyLevels,
		Status:   enums.BotStatusActive,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, expected))
	})

	params := types.NewLevelsBotParams("BTC-PERPETUAL", []float64{49000, 48000}, []float64{51000, 52000}, 0.1, 1700100000)
	result, err := c.CreateLevelsBot(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Strategy != enums.BotStrategyLevels {
		t.Errorf("expected strategy=levels, got %s", result.Strategy)
	}
}

func TestCreateDHedgeBot_Success(t *testing.T) {
	expected := types.Bot{
		BotID:    "bot-dhedge-001",
		Strategy: enums.BotStrategyDHedge,
		Status:   enums.BotStatusActive,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, expected))
	})

	params := types.NewDHedgeBotParams("BTC-PERPETUAL", 60).
		WithTargetDelta(0).
		WithThreshold(0.1).
		WithLabel("test-dhedge")
	result, err := c.CreateDHedgeBot(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Strategy != enums.BotStrategyDHedge {
		t.Errorf("expected strategy=dhedge, got %s", result.Strategy)
	}
}

func TestCreateDFollowBot_Success(t *testing.T) {
	expected := types.Bot{
		BotID:    "bot-dfollow-001",
		Strategy: enums.BotStrategyDFollow,
		Status:   enums.BotStatusActive,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, expected))
	})

	params := types.NewDFollowBotParams("BTC-PERPETUAL", "ETH-PERPETUAL", 1.0, 60, 1700100000)
	result, err := c.CreateDFollowBot(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Strategy != enums.BotStrategyDFollow {
		t.Errorf("expected strategy=dfollow, got %s", result.Strategy)
	}
}

func TestCancelBot_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/cancel_bot" {
			t.Errorf("expected path /private/cancel_bot, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params struct {
			BotID string `json:"bot_id"`
		}
		json.Unmarshal(body, &params)
		if params.BotID != "bot-001" {
			t.Errorf("expected bot_id=bot-001, got %s", params.BotID)
		}
		w.Write([]byte(`{"result":null}`))
	})

	err := c.CancelBot(context.Background(), "bot-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelBot_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10020, "bot not found"))
	})

	err := c.CancelBot(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCancelAllBots_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/cancel_all_bots" {
			t.Errorf("expected path /private/cancel_all_bots, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, map[string]int{"n_cancelled": 3}))
	})

	n, err := c.CancelAllBots(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 cancelled, got %d", n)
	}
}

func TestCancelAllBots_Zero(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, map[string]int{"n_cancelled": 0}))
	})

	n, err := c.CancelAllBots(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 cancelled, got %d", n)
	}
}
