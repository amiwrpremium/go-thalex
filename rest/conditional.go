package rest

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// ConditionalOrders retrieves all active conditional orders.
func (c *Client) ConditionalOrders(ctx context.Context) ([]types.ConditionalOrder, error) {
	var result []types.ConditionalOrder
	err := c.transport.DoPrivateGET(ctx, "/private/conditional_orders", nil, &result)
	return result, err
}

// CreateConditionalOrder creates a new conditional order.
func (c *Client) CreateConditionalOrder(ctx context.Context, params *types.CreateConditionalOrderParams) (types.ConditionalOrder, error) {
	var result types.ConditionalOrder
	err := c.transport.DoPrivatePOST(ctx, "/private/create_conditional_order", params, &result)
	return result, err
}

// CancelConditionalOrder cancels a conditional order by ID.
func (c *Client) CancelConditionalOrder(ctx context.Context, orderID string) error {
	body := struct {
		OrderID string `json:"order_id"`
	}{OrderID: orderID}
	return c.transport.DoPrivatePOST(ctx, "/private/cancel_conditional_order", body, nil)
}

// CancelAllConditionalOrders cancels all conditional orders.
func (c *Client) CancelAllConditionalOrders(ctx context.Context) (int, error) {
	var result cancelAllResult
	err := c.transport.DoPrivatePOST(ctx, "/private/cancel_all_conditional_orders", nil, &result)
	return result.NCancelled, err
}
