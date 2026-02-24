package rest

import (
	"context"
	"net/url"
	"strconv"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// MarkPriceHistoricalData retrieves mark price historical data in OHLC format.
func (c *Client) MarkPriceHistoricalData(ctx context.Context, instrumentName string, from, to float64, resolution enums.Resolution) (types.MarkPriceHistoricalResult, error) {
	q := url.Values{}
	q.Set("instrument_name", instrumentName)
	q.Set("from", strconv.FormatFloat(from, 'f', -1, 64))
	q.Set("to", strconv.FormatFloat(to, 'f', -1, 64))
	q.Set("resolution", string(resolution))
	var result types.MarkPriceHistoricalResult
	err := c.transport.DoPublic(ctx, "/public/mark_price_historical_data", q, &result)
	return result, err
}

// IndexPriceHistoricalData retrieves index price historical data in OHLC format.
func (c *Client) IndexPriceHistoricalData(ctx context.Context, indexName string, from, to float64, resolution enums.Resolution) (types.IndexPriceHistoricalResult, error) {
	q := url.Values{}
	q.Set("index_name", indexName)
	q.Set("from", strconv.FormatFloat(from, 'f', -1, 64))
	q.Set("to", strconv.FormatFloat(to, 'f', -1, 64))
	q.Set("resolution", string(resolution))
	var result types.IndexPriceHistoricalResult
	err := c.transport.DoPublic(ctx, "/public/index_price_historical_data", q, &result)
	return result, err
}
