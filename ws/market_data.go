package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// Instruments retrieves active instruments via WebSocket.
func (ws *Client) Instruments(ctx context.Context) ([]types.Instrument, error) {
	var result []types.Instrument
	err := ws.call(ctx, "public/instruments", nil, &result)
	return result, err
}

// AllInstruments retrieves all instruments including inactive via WebSocket.
func (ws *Client) AllInstruments(ctx context.Context) ([]types.Instrument, error) {
	var result []types.Instrument
	err := ws.call(ctx, "public/all_instruments", nil, &result)
	return result, err
}

// Instrument retrieves a single instrument by name via WebSocket.
func (ws *Client) Instrument(ctx context.Context, instrumentName string) (types.Instrument, error) {
	var result types.Instrument
	err := ws.call(ctx, "public/instrument", map[string]any{
		"instrument_name": instrumentName,
	}, &result)
	return result, err
}

// Ticker retrieves a ticker via WebSocket.
func (ws *Client) Ticker(ctx context.Context, instrumentName string) (types.Ticker, error) {
	var result types.Ticker
	err := ws.call(ctx, "public/ticker", map[string]any{
		"instrument_name": instrumentName,
	}, &result)
	return result, err
}

// Index retrieves an index price via WebSocket.
func (ws *Client) Index(ctx context.Context, underlying string) (types.IndexPrice, error) {
	var result types.IndexPrice
	err := ws.call(ctx, "public/index", map[string]any{
		"underlying": underlying,
	}, &result)
	return result, err
}

// Book retrieves the order book via WebSocket.
func (ws *Client) Book(ctx context.Context, instrumentName string) (types.Book, error) {
	var result types.Book
	err := ws.call(ctx, "public/book", map[string]any{
		"instrument_name": instrumentName,
	}, &result)
	return result, err
}

// SystemInfo retrieves system status information via WebSocket.
func (ws *Client) SystemInfo(ctx context.Context) (types.SystemInfo, error) {
	var result types.SystemInfo
	err := ws.call(ctx, "public/system_info", nil, &result)
	return result, err
}

// MarkPriceHistoricalData retrieves mark price historical data via WebSocket.
func (ws *Client) MarkPriceHistoricalData(ctx context.Context, instrumentName string, from, to float64, resolution enums.Resolution) (types.MarkPriceHistoricalResult, error) {
	var result types.MarkPriceHistoricalResult
	err := ws.call(ctx, "public/mark_price_historical_data", map[string]any{
		"instrument_name": instrumentName,
		"from":            from,
		"to":              to,
		"resolution":      resolution,
	}, &result)
	return result, err
}

// IndexPriceHistoricalData retrieves index price historical data via WebSocket.
func (ws *Client) IndexPriceHistoricalData(ctx context.Context, indexName string, from, to float64, resolution enums.Resolution) (types.IndexPriceHistoricalResult, error) {
	var result types.IndexPriceHistoricalResult
	err := ws.call(ctx, "public/index_price_historical_data", map[string]any{
		"index_name": indexName,
		"from":       from,
		"to":         to,
		"resolution": resolution,
	}, &result)
	return result, err
}
