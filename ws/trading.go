package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// Insert places a new order via WebSocket.
func (ws *Client) Insert(ctx context.Context, params *types.InsertOrderParams) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := ws.call(ctx, "private/insert", params, &result)
	return result, err
}

// Buy places a market buy order via WebSocket.
func (ws *Client) Buy(ctx context.Context, instrumentName string, amount float64) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := ws.call(ctx, "private/buy", map[string]any{
		"instrument_name": instrumentName,
		"amount":          amount,
	}, &result)
	return result, err
}

// Sell places a market sell order via WebSocket.
func (ws *Client) Sell(ctx context.Context, instrumentName string, amount float64) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := ws.call(ctx, "private/sell", map[string]any{
		"instrument_name": instrumentName,
		"amount":          amount,
	}, &result)
	return result, err
}

// Amend modifies an existing order via WebSocket.
func (ws *Client) Amend(ctx context.Context, params *types.AmendOrderParams) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := ws.call(ctx, "private/amend", params, &result)
	return result, err
}

// Cancel cancels an existing order via WebSocket.
func (ws *Client) Cancel(ctx context.Context, params *types.CancelOrderParams) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := ws.call(ctx, "private/cancel", params, &result)
	return result, err
}

// CancelAll cancels all orders, returning the number cancelled.
func (ws *Client) CancelAll(ctx context.Context) (int, error) {
	var result struct {
		NCancelled int `json:"n_cancelled"`
	}
	err := ws.call(ctx, "private/cancel_all", nil, &result)
	return result.NCancelled, err
}

// OpenOrders retrieves open orders, optionally filtered by instrument.
func (ws *Client) OpenOrders(ctx context.Context, instrumentName string) ([]types.OrderStatus, error) {
	params := map[string]any{}
	if instrumentName != "" {
		params["instrument_name"] = instrumentName
	}
	var result struct {
		Orders []types.OrderStatus `json:"orders"`
	}
	err := ws.call(ctx, "private/open_orders", params, &result)
	return result.Orders, err
}
