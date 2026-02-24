package rest

import (
	"context"
	"net/url"

	"github.com/amiwrpremium/go-thalex/types"
)

// Insert places a new order.
func (c *Client) Insert(ctx context.Context, params *types.InsertOrderParams) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := c.transport.DoPrivatePOST(ctx, "/private/insert", params, &result)
	return result, err
}

// Buy places a market buy order.
func (c *Client) Buy(ctx context.Context, instrumentName string, amount float64) (types.OrderStatus, error) {
	var result types.OrderStatus
	body := struct {
		InstrumentName string  `json:"instrument_name"`
		Amount         float64 `json:"amount"`
	}{InstrumentName: instrumentName, Amount: amount}
	err := c.transport.DoPrivatePOST(ctx, "/private/buy", body, &result)
	return result, err
}

// Sell places a market sell order.
func (c *Client) Sell(ctx context.Context, instrumentName string, amount float64) (types.OrderStatus, error) {
	var result types.OrderStatus
	body := struct {
		InstrumentName string  `json:"instrument_name"`
		Amount         float64 `json:"amount"`
	}{InstrumentName: instrumentName, Amount: amount}
	err := c.transport.DoPrivatePOST(ctx, "/private/sell", body, &result)
	return result, err
}

// Amend modifies an existing order.
func (c *Client) Amend(ctx context.Context, params *types.AmendOrderParams) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := c.transport.DoPrivatePOST(ctx, "/private/amend", params, &result)
	return result, err
}

// Cancel cancels an existing order.
func (c *Client) Cancel(ctx context.Context, params *types.CancelOrderParams) (types.OrderStatus, error) {
	var result types.OrderStatus
	err := c.transport.DoPrivatePOST(ctx, "/private/cancel", params, &result)
	return result, err
}

type cancelAllResult struct {
	NCancelled int `json:"n_cancelled"`
}

// CancelAll cancels all orders, returning the number of orders cancelled.
func (c *Client) CancelAll(ctx context.Context) (int, error) {
	var result cancelAllResult
	err := c.transport.DoPrivatePOST(ctx, "/private/cancel_all", nil, &result)
	return result.NCancelled, err
}

type openOrdersResult struct {
	Orders []types.OrderStatus `json:"orders"`
}

// OpenOrders retrieves all open orders, optionally filtered by instrument.
func (c *Client) OpenOrders(ctx context.Context, instrumentName string) ([]types.OrderStatus, error) {
	q := url.Values{}
	if instrumentName != "" {
		q.Set("instrument_name", instrumentName)
	}
	var result openOrdersResult
	err := c.transport.DoPrivateGET(ctx, "/private/open_orders", q, &result)
	return result.Orders, err
}
