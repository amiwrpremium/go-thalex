package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// Bots retrieves all bots via WebSocket.
func (ws *Client) Bots(ctx context.Context, includeInactive bool) ([]types.Bot, error) {
	var result []types.Bot
	err := ws.call(ctx, "private/bots", map[string]any{
		"include_inactive": includeInactive,
	}, &result)
	return result, err
}

// CreateSGSLBot creates a new SGSL bot via WebSocket.
func (ws *Client) CreateSGSLBot(ctx context.Context, params *types.SGSLBotParams) (types.Bot, error) {
	var result types.Bot
	err := ws.call(ctx, "private/create_bot", params, &result)
	return result, err
}

// CreateOCQBot creates a new OCQ bot via WebSocket.
func (ws *Client) CreateOCQBot(ctx context.Context, params *types.OCQBotParams) (types.Bot, error) {
	var result types.Bot
	err := ws.call(ctx, "private/create_bot", params, &result)
	return result, err
}

// CreateLevelsBot creates a new Levels bot via WebSocket.
func (ws *Client) CreateLevelsBot(ctx context.Context, params *types.LevelsBotParams) (types.Bot, error) {
	var result types.Bot
	err := ws.call(ctx, "private/create_bot", params, &result)
	return result, err
}

// CreateGridBot creates a new Grid bot via WebSocket.
func (ws *Client) CreateGridBot(ctx context.Context, params *types.GridBotParams) (types.Bot, error) {
	var result types.Bot
	err := ws.call(ctx, "private/create_bot", params, &result)
	return result, err
}

// CreateDHedgeBot creates a new Delta Hedger bot via WebSocket.
func (ws *Client) CreateDHedgeBot(ctx context.Context, params *types.DHedgeBotParams) (types.Bot, error) {
	var result types.Bot
	err := ws.call(ctx, "private/create_bot", params, &result)
	return result, err
}

// CreateDFollowBot creates a new Delta Follower bot via WebSocket.
func (ws *Client) CreateDFollowBot(ctx context.Context, params *types.DFollowBotParams) (types.Bot, error) {
	var result types.Bot
	err := ws.call(ctx, "private/create_bot", params, &result)
	return result, err
}

// CancelBot cancels a running bot via WebSocket.
func (ws *Client) CancelBot(ctx context.Context, botID string) error {
	return ws.callNoResult(ctx, "private/cancel_bot", map[string]any{
		"bot_id": botID,
	})
}

// CancelAllBots cancels all running bots via WebSocket.
func (ws *Client) CancelAllBots(ctx context.Context) (int, error) {
	var result struct {
		NCancelled int `json:"n_cancelled"`
	}
	err := ws.call(ctx, "private/cancel_all_bots", nil, &result)
	return result.NCancelled, err
}
