package rest

import (
	"context"
	"net/url"
	"strconv"

	"github.com/amiwrpremium/go-thalex/types"
)

// CreateRfq creates a new Request for Quote.
func (c *Client) CreateRfq(ctx context.Context, params *types.CreateRfqParams) (types.Rfq, error) {
	var result types.Rfq
	err := c.transport.DoPrivatePOST(ctx, "/private/create_rfq", params, &result)
	return result, err
}

// CancelRfq cancels an open RFQ.
func (c *Client) CancelRfq(ctx context.Context, rfqID string) error {
	body := struct {
		RfqID string `json:"rfq_id"`
	}{RfqID: rfqID}
	return c.transport.DoPrivatePOST(ctx, "/private/cancel_rfq", body, nil)
}

// TradeRfq executes a trade on an RFQ.
func (c *Client) TradeRfq(ctx context.Context, params *types.TradeRfqParams) ([]types.Trade, error) {
	var result []types.Trade
	err := c.transport.DoPrivatePOST(ctx, "/private/trade_rfq", params, &result)
	return result, err
}

// OpenRfqs retrieves all open RFQs created by this account.
func (c *Client) OpenRfqs(ctx context.Context) ([]types.Rfq, error) {
	var result []types.Rfq
	err := c.transport.DoPrivateGET(ctx, "/private/open_rfqs", nil, &result)
	return result, err
}

// RfqHistory retrieves historical RFQs.
func (c *Client) RfqHistory(ctx context.Context, from, to *float64, offset, limit *int) ([]types.Rfq, error) {
	q := url.Values{}
	if from != nil {
		q.Set("from", strconv.FormatFloat(*from, 'f', -1, 64))
	}
	if to != nil {
		q.Set("to", strconv.FormatFloat(*to, 'f', -1, 64))
	}
	if offset != nil {
		q.Set("offset", strconv.Itoa(*offset))
	}
	if limit != nil {
		q.Set("limit", strconv.Itoa(*limit))
	}
	var result []types.Rfq
	err := c.transport.DoPrivateGET(ctx, "/private/rfq_history", q, &result)
	return result, err
}

// MMRfqs retrieves all market maker RFQ opportunities.
func (c *Client) MMRfqs(ctx context.Context) ([]types.Rfq, error) {
	var result []types.Rfq
	err := c.transport.DoPrivateGET(ctx, "/private/mm_rfqs", nil, &result)
	return result, err
}

// MMRfqInsertQuote inserts a quote on an RFQ.
func (c *Client) MMRfqInsertQuote(ctx context.Context, params *types.RfqQuoteInsertParams) (types.RfqOrder, error) {
	var result types.RfqOrder
	err := c.transport.DoPrivatePOST(ctx, "/private/mm_rfq_insert_quote", params, &result)
	return result, err
}

// MMRfqAmendQuote amends a quote on an RFQ.
func (c *Client) MMRfqAmendQuote(ctx context.Context, params *types.RfqQuoteAmendParams) (types.RfqOrder, error) {
	var result types.RfqOrder
	err := c.transport.DoPrivatePOST(ctx, "/private/mm_rfq_amend_quote", params, &result)
	return result, err
}

// MMRfqDeleteQuote deletes a quote from an RFQ.
func (c *Client) MMRfqDeleteQuote(ctx context.Context, params *types.RfqQuoteDeleteParams) error {
	return c.transport.DoPrivatePOST(ctx, "/private/mm_rfq_delete_quote", params, nil)
}

// MMRfqQuotes retrieves all active RFQ quotes.
func (c *Client) MMRfqQuotes(ctx context.Context) ([]types.RfqOrder, error) {
	var result []types.RfqOrder
	err := c.transport.DoPrivateGET(ctx, "/private/mm_rfq_quotes", nil, &result)
	return result, err
}
