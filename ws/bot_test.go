package ws

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------------------------------------------------------------------------
// Bots
// ---------------------------------------------------------------------------

func TestBots_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/bots" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"bot_id":"bot-1","type":"sgsl","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bots, err := c.Bots(ctx, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bots) != 1 {
		t.Fatalf("expected 1 bot, got %d", len(bots))
	}
}

func TestBots_IncludeInactive(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`[{"bot_id":"bot-1","type":"sgsl","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0},{"bot_id":"bot-2","type":"grid","status":"cancelled","instrument_name":"ETH-PERPETUAL","create_time":2.0}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bots, err := c.Bots(ctx, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bots) != 2 {
		t.Errorf("expected 2 bots, got %d", len(bots))
	}
}

// ---------------------------------------------------------------------------
// CreateSGSLBot
// ---------------------------------------------------------------------------

func TestCreateSGSLBot_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/create_bot" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"bot_id":"bot-new","type":"sgsl","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CreateSGSLBot(ctx, &types.SGSLBotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-new" {
		t.Errorf("expected bot_id=bot-new, got %q", result.BotID)
	}
}

// ---------------------------------------------------------------------------
// CreateOCQBot
// ---------------------------------------------------------------------------

func TestCreateOCQBot_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"bot_id":"bot-ocq","type":"ocq","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CreateOCQBot(ctx, &types.OCQBotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-ocq" {
		t.Errorf("expected bot_id=bot-ocq, got %q", result.BotID)
	}
}

// ---------------------------------------------------------------------------
// CreateLevelsBot
// ---------------------------------------------------------------------------

func TestCreateLevelsBot_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"bot_id":"bot-lvl","type":"levels","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CreateLevelsBot(ctx, &types.LevelsBotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-lvl" {
		t.Errorf("expected bot_id=bot-lvl, got %q", result.BotID)
	}
}

// ---------------------------------------------------------------------------
// CreateGridBot
// ---------------------------------------------------------------------------

func TestCreateGridBot_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"bot_id":"bot-grid","type":"grid","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CreateGridBot(ctx, &types.GridBotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-grid" {
		t.Errorf("expected bot_id=bot-grid, got %q", result.BotID)
	}
}

// ---------------------------------------------------------------------------
// CreateDHedgeBot
// ---------------------------------------------------------------------------

func TestCreateDHedgeBot_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"bot_id":"bot-dh","type":"dhedge","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CreateDHedgeBot(ctx, &types.DHedgeBotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-dh" {
		t.Errorf("expected bot_id=bot-dh, got %q", result.BotID)
	}
}

// ---------------------------------------------------------------------------
// CreateDFollowBot
// ---------------------------------------------------------------------------

func TestCreateDFollowBot_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return json.RawMessage(`{"bot_id":"bot-df","type":"dfollow","status":"active","instrument_name":"BTC-PERPETUAL","create_time":1.0}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CreateDFollowBot(ctx, &types.DFollowBotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.BotID != "bot-df" {
		t.Errorf("expected bot_id=bot-df, got %q", result.BotID)
	}
}

// ---------------------------------------------------------------------------
// CancelBot
// ---------------------------------------------------------------------------

func TestCancelBot_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/cancel_bot" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.CancelBot(ctx, "bot-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// CancelAllBots
// ---------------------------------------------------------------------------

func TestCancelAllBots_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/cancel_all_bots" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"n_cancelled":4}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := c.CancelAllBots(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Errorf("expected 4 cancelled, got %d", n)
	}
}

func TestCancelAllBots_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "not authenticated"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := c.CancelAllBots(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if n != 0 {
		t.Errorf("expected 0 on error, got %d", n)
	}
}
