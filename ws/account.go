package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// Portfolio retrieves portfolio positions via WebSocket.
func (ws *Client) Portfolio(ctx context.Context) ([]types.PortfolioEntry, error) {
	var result []types.PortfolioEntry
	err := ws.call(ctx, "private/portfolio", nil, &result)
	return result, err
}

// AccountSummary retrieves the account summary via WebSocket.
func (ws *Client) AccountSummary(ctx context.Context) (types.AccountSummary, error) {
	var result types.AccountSummary
	err := ws.call(ctx, "private/account_summary", nil, &result)
	return result, err
}

// AccountBreakdown retrieves a detailed account breakdown via WebSocket.
func (ws *Client) AccountBreakdown(ctx context.Context) (types.AccountBreakdown, error) {
	var result types.AccountBreakdown
	err := ws.call(ctx, "private/account_breakdown", nil, &result)
	return result, err
}

// RequiredMarginBreakdown retrieves the margin breakdown via WebSocket.
func (ws *Client) RequiredMarginBreakdown(ctx context.Context) (types.PortfolioMarginBreakdown, error) {
	var result types.PortfolioMarginBreakdown
	err := ws.call(ctx, "private/required_margin_breakdown", nil, &result)
	return result, err
}

// RequiredMarginForOrder checks margin impact of a hypothetical order via WebSocket.
func (ws *Client) RequiredMarginForOrder(ctx context.Context, instrumentName string, price, amount float64) (types.MarginForOrderResult, error) {
	var result types.MarginForOrderResult
	err := ws.call(ctx, "private/required_margin_for_order", map[string]any{
		"instrument_name": instrumentName,
		"price":           price,
		"amount":          amount,
	}, &result)
	return result, err
}
