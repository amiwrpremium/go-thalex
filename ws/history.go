package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// TradeHistory retrieves trade history via WebSocket.
func (ws *Client) TradeHistory(ctx context.Context, params *types.TradeHistoryParams) ([]types.Trade, error) {
	var result []types.Trade
	err := ws.call(ctx, "private/trade_history", params, &result)
	return result, err
}

// OrderHistory retrieves order history via WebSocket.
func (ws *Client) OrderHistory(ctx context.Context, params *types.OrderHistoryParams) ([]types.OrderHistory, error) {
	var result []types.OrderHistory
	err := ws.call(ctx, "private/order_history", params, &result)
	return result, err
}

// DailyMarkHistory retrieves daily mark history via WebSocket.
func (ws *Client) DailyMarkHistory(ctx context.Context, params *types.DailyMarkHistoryParams) ([]types.DailyMark, error) {
	var result []types.DailyMark
	err := ws.call(ctx, "private/daily_mark_history", params, &result)
	return result, err
}

// TransactionHistory retrieves transaction history via WebSocket.
func (ws *Client) TransactionHistory(ctx context.Context, params *types.TransactionHistoryParams) ([]types.Transaction, error) {
	var result []types.Transaction
	err := ws.call(ctx, "private/transaction_history", params, &result)
	return result, err
}
