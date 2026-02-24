package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// CreateRfq creates a new RFQ via WebSocket.
func (ws *Client) CreateRfq(ctx context.Context, params *types.CreateRfqParams) (types.Rfq, error) {
	var result types.Rfq
	err := ws.call(ctx, "private/create_rfq", params, &result)
	return result, err
}

// CancelRfq cancels an open RFQ via WebSocket.
func (ws *Client) CancelRfq(ctx context.Context, rfqID string) error {
	return ws.callNoResult(ctx, "private/cancel_rfq", map[string]any{
		"rfq_id": rfqID,
	})
}

// TradeRfq executes a trade on an RFQ via WebSocket.
func (ws *Client) TradeRfq(ctx context.Context, params *types.TradeRfqParams) ([]types.Trade, error) {
	var result []types.Trade
	err := ws.call(ctx, "private/trade_rfq", params, &result)
	return result, err
}

// OpenRfqs retrieves open RFQs via WebSocket.
func (ws *Client) OpenRfqs(ctx context.Context) ([]types.Rfq, error) {
	var result []types.Rfq
	err := ws.call(ctx, "private/open_rfqs", nil, &result)
	return result, err
}

// RfqHistory retrieves historical RFQs via WebSocket.
func (ws *Client) RfqHistory(ctx context.Context, from, to *float64, offset, limit *int) ([]types.Rfq, error) {
	params := map[string]any{}
	if from != nil {
		params["from"] = *from
	}
	if to != nil {
		params["to"] = *to
	}
	if offset != nil {
		params["offset"] = *offset
	}
	if limit != nil {
		params["limit"] = *limit
	}
	var result []types.Rfq
	err := ws.call(ctx, "private/rfq_history", params, &result)
	return result, err
}

// MMRfqs retrieves market maker RFQ opportunities via WebSocket.
func (ws *Client) MMRfqs(ctx context.Context) ([]types.Rfq, error) {
	var result []types.Rfq
	err := ws.call(ctx, "private/mm_rfqs", nil, &result)
	return result, err
}

// MMRfqInsertQuote inserts a quote on an RFQ via WebSocket.
func (ws *Client) MMRfqInsertQuote(ctx context.Context, params *types.RfqQuoteInsertParams) (types.RfqOrder, error) {
	var result types.RfqOrder
	err := ws.call(ctx, "private/mm_rfq_insert_quote", params, &result)
	return result, err
}

// MMRfqAmendQuote amends a quote on an RFQ via WebSocket.
func (ws *Client) MMRfqAmendQuote(ctx context.Context, params *types.RfqQuoteAmendParams) (types.RfqOrder, error) {
	var result types.RfqOrder
	err := ws.call(ctx, "private/mm_rfq_amend_quote", params, &result)
	return result, err
}

// MMRfqDeleteQuote deletes a quote from an RFQ via WebSocket.
func (ws *Client) MMRfqDeleteQuote(ctx context.Context, params *types.RfqQuoteDeleteParams) error {
	return ws.callNoResult(ctx, "private/mm_rfq_delete_quote", params)
}

// MMRfqQuotes retrieves active RFQ quotes via WebSocket.
func (ws *Client) MMRfqQuotes(ctx context.Context) ([]types.RfqOrder, error) {
	var result []types.RfqOrder
	err := ws.call(ctx, "private/mm_rfq_quotes", nil, &result)
	return result, err
}
