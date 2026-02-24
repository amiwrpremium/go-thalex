package rest

import (
	"context"
	"net/url"
	"strconv"

	"github.com/amiwrpremium/go-thalex/types"
)

// Portfolio retrieves the current portfolio positions.
func (c *Client) Portfolio(ctx context.Context) ([]types.PortfolioEntry, error) {
	var result []types.PortfolioEntry
	err := c.transport.DoPrivateGET(ctx, "/private/portfolio", nil, &result)
	return result, err
}

// AccountSummary retrieves the account financial summary.
func (c *Client) AccountSummary(ctx context.Context) (types.AccountSummary, error) {
	var result types.AccountSummary
	err := c.transport.DoPrivateGET(ctx, "/private/account_summary", nil, &result)
	return result, err
}

// AccountBreakdown retrieves a detailed account breakdown.
func (c *Client) AccountBreakdown(ctx context.Context) (types.AccountBreakdown, error) {
	var result types.AccountBreakdown
	err := c.transport.DoPrivateGET(ctx, "/private/account_breakdown", nil, &result)
	return result, err
}

// RequiredMarginBreakdown retrieves the portfolio margin breakdown.
func (c *Client) RequiredMarginBreakdown(ctx context.Context) (types.PortfolioMarginBreakdown, error) {
	var result types.PortfolioMarginBreakdown
	err := c.transport.DoPrivateGET(ctx, "/private/required_margin_breakdown", nil, &result)
	return result, err
}

// RequiredMarginForOrder checks margin impact of a hypothetical order.
func (c *Client) RequiredMarginForOrder(ctx context.Context, instrumentName string, price, amount float64) (types.MarginForOrderResult, error) {
	q := url.Values{}
	q.Set("instrument_name", instrumentName)
	q.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
	q.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	var result types.MarginForOrderResult
	err := c.transport.DoPrivateGET(ctx, "/private/required_margin_for_order", q, &result)
	return result, err
}
