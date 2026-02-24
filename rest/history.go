package rest

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/amiwrpremium/go-thalex/types"
)

// TradeHistory retrieves trade history with optional filters.
func (c *Client) TradeHistory(ctx context.Context, params *types.TradeHistoryParams) ([]types.Trade, error) {
	q := url.Values{}
	if params != nil {
		if params.From != nil {
			q.Set("from", strconv.FormatFloat(*params.From, 'f', -1, 64))
		}
		if params.To != nil {
			q.Set("to", strconv.FormatFloat(*params.To, 'f', -1, 64))
		}
		if params.Offset != nil {
			q.Set("offset", strconv.Itoa(*params.Offset))
		}
		if params.Limit != nil {
			q.Set("limit", strconv.Itoa(*params.Limit))
		}
		if params.Sort != "" {
			q.Set("sort", string(params.Sort))
		}
		if len(params.InstrumentNames) > 0 {
			q.Set("instrument_names", strings.Join(params.InstrumentNames, ","))
		}
		if len(params.BotIDs) > 0 {
			q.Set("bot_ids", strings.Join(params.BotIDs, ","))
		}
	}
	var result []types.Trade
	err := c.transport.DoPrivateGET(ctx, "/private/trade_history", q, &result)
	return result, err
}

// OrderHistory retrieves order history with optional filters.
func (c *Client) OrderHistory(ctx context.Context, params *types.OrderHistoryParams) ([]types.OrderHistory, error) {
	q := url.Values{}
	if params != nil {
		if params.From != nil {
			q.Set("from", strconv.FormatFloat(*params.From, 'f', -1, 64))
		}
		if params.To != nil {
			q.Set("to", strconv.FormatFloat(*params.To, 'f', -1, 64))
		}
		if params.Offset != nil {
			q.Set("offset", strconv.Itoa(*params.Offset))
		}
		if params.Limit != nil {
			q.Set("limit", strconv.Itoa(*params.Limit))
		}
		if params.Sort != "" {
			q.Set("sort", string(params.Sort))
		}
		if len(params.InstrumentNames) > 0 {
			q.Set("instrument_names", strings.Join(params.InstrumentNames, ","))
		}
	}
	var result []types.OrderHistory
	err := c.transport.DoPrivateGET(ctx, "/private/order_history", q, &result)
	return result, err
}

// DailyMarkHistory retrieves daily mark history.
func (c *Client) DailyMarkHistory(ctx context.Context, params *types.DailyMarkHistoryParams) ([]types.DailyMark, error) {
	q := url.Values{}
	if params != nil {
		if params.From != nil {
			q.Set("from", strconv.FormatFloat(*params.From, 'f', -1, 64))
		}
		if params.To != nil {
			q.Set("to", strconv.FormatFloat(*params.To, 'f', -1, 64))
		}
		if params.Offset != nil {
			q.Set("offset", strconv.Itoa(*params.Offset))
		}
		if params.Limit != nil {
			q.Set("limit", strconv.Itoa(*params.Limit))
		}
	}
	var result []types.DailyMark
	err := c.transport.DoPrivateGET(ctx, "/private/daily_mark_history", q, &result)
	return result, err
}

// TransactionHistory retrieves transaction history.
func (c *Client) TransactionHistory(ctx context.Context, params *types.TransactionHistoryParams) ([]types.Transaction, error) {
	q := url.Values{}
	if params != nil {
		if params.From != nil {
			q.Set("from", strconv.FormatFloat(*params.From, 'f', -1, 64))
		}
		if params.To != nil {
			q.Set("to", strconv.FormatFloat(*params.To, 'f', -1, 64))
		}
		if params.Offset != nil {
			q.Set("offset", strconv.Itoa(*params.Offset))
		}
		if params.Limit != nil {
			q.Set("limit", strconv.Itoa(*params.Limit))
		}
		if params.Sort != "" {
			q.Set("sort", string(params.Sort))
		}
	}
	var result []types.Transaction
	err := c.transport.DoPrivateGET(ctx, "/private/transaction_history", q, &result)
	return result, err
}
