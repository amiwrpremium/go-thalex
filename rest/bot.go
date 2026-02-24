package rest

import (
	"context"
	"net/url"

	"github.com/amiwrpremium/go-thalex/types"
)

// Bots retrieves all bots, optionally including inactive ones.
func (c *Client) Bots(ctx context.Context, includeInactive bool) ([]types.Bot, error) {
	q := url.Values{}
	if includeInactive {
		q.Set("include_inactive", "true")
	}
	var result []types.Bot
	err := c.transport.DoPrivateGET(ctx, "/private/bots", q, &result)
	return result, err
}

// CreateSGSLBot creates a new SGSL bot.
func (c *Client) CreateSGSLBot(ctx context.Context, params *types.SGSLBotParams) (types.Bot, error) {
	var result types.Bot
	err := c.transport.DoPrivatePOST(ctx, "/private/create_bot", params, &result)
	return result, err
}

// CreateOCQBot creates a new OCQ bot.
func (c *Client) CreateOCQBot(ctx context.Context, params *types.OCQBotParams) (types.Bot, error) {
	var result types.Bot
	err := c.transport.DoPrivatePOST(ctx, "/private/create_bot", params, &result)
	return result, err
}

// CreateLevelsBot creates a new Levels bot.
func (c *Client) CreateLevelsBot(ctx context.Context, params *types.LevelsBotParams) (types.Bot, error) {
	var result types.Bot
	err := c.transport.DoPrivatePOST(ctx, "/private/create_bot", params, &result)
	return result, err
}

// CreateGridBot creates a new Grid bot.
func (c *Client) CreateGridBot(ctx context.Context, params *types.GridBotParams) (types.Bot, error) {
	var result types.Bot
	err := c.transport.DoPrivatePOST(ctx, "/private/create_bot", params, &result)
	return result, err
}

// CreateDHedgeBot creates a new Delta Hedger bot.
func (c *Client) CreateDHedgeBot(ctx context.Context, params *types.DHedgeBotParams) (types.Bot, error) {
	var result types.Bot
	err := c.transport.DoPrivatePOST(ctx, "/private/create_bot", params, &result)
	return result, err
}

// CreateDFollowBot creates a new Delta Follower bot.
func (c *Client) CreateDFollowBot(ctx context.Context, params *types.DFollowBotParams) (types.Bot, error) {
	var result types.Bot
	err := c.transport.DoPrivatePOST(ctx, "/private/create_bot", params, &result)
	return result, err
}

// CancelBot cancels a running bot.
func (c *Client) CancelBot(ctx context.Context, botID string) error {
	body := struct {
		BotID string `json:"bot_id"`
	}{BotID: botID}
	return c.transport.DoPrivatePOST(ctx, "/private/cancel_bot", body, nil)
}

// CancelAllBots cancels all running bots.
func (c *Client) CancelAllBots(ctx context.Context) (int, error) {
	var result cancelAllResult
	err := c.transport.DoPrivatePOST(ctx, "/private/cancel_all_bots", nil, &result)
	return result.NCancelled, err
}
