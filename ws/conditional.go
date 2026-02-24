package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// ConditionalOrders retrieves all active conditional orders via WebSocket.
func (ws *Client) ConditionalOrders(ctx context.Context) ([]types.ConditionalOrder, error) {
	var result []types.ConditionalOrder
	err := ws.call(ctx, "private/conditional_orders", nil, &result)
	return result, err
}

// CreateConditionalOrder creates a new conditional order via WebSocket.
func (ws *Client) CreateConditionalOrder(ctx context.Context, params *types.CreateConditionalOrderParams) (types.ConditionalOrder, error) {
	var result types.ConditionalOrder
	err := ws.call(ctx, "private/create_conditional_order", params, &result)
	return result, err
}

// CancelConditionalOrder cancels a conditional order via WebSocket.
func (ws *Client) CancelConditionalOrder(ctx context.Context, orderID string) error {
	return ws.callNoResult(ctx, "private/cancel_conditional_order", map[string]any{
		"order_id": orderID,
	})
}

// CancelAllConditionalOrders cancels all conditional orders via WebSocket.
func (ws *Client) CancelAllConditionalOrders(ctx context.Context) (int, error) {
	var result struct {
		NCancelled int `json:"n_cancelled"`
	}
	err := ws.call(ctx, "private/cancel_all_conditional_orders", nil, &result)
	return result.NCancelled, err
}
