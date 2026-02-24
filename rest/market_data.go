package rest

import (
	"context"
	"net/url"

	"github.com/amiwrpremium/go-thalex/types"
)

// Instruments retrieves all currently active instruments.
func (c *Client) Instruments(ctx context.Context) ([]types.Instrument, error) {
	var result []types.Instrument
	err := c.transport.DoPublic(ctx, "/public/instruments", nil, &result)
	return result, err
}

// AllInstruments retrieves all instruments including inactive ones.
func (c *Client) AllInstruments(ctx context.Context) ([]types.Instrument, error) {
	var result []types.Instrument
	err := c.transport.DoPublic(ctx, "/public/all_instruments", nil, &result)
	return result, err
}

// Instrument retrieves a single instrument by name.
func (c *Client) Instrument(ctx context.Context, instrumentName string) (types.Instrument, error) {
	q := url.Values{}
	q.Set("instrument_name", instrumentName)
	var result types.Instrument
	err := c.transport.DoPublic(ctx, "/public/instrument", q, &result)
	return result, err
}

// Ticker retrieves the ticker for a single instrument.
func (c *Client) Ticker(ctx context.Context, instrumentName string) (types.Ticker, error) {
	q := url.Values{}
	q.Set("instrument_name", instrumentName)
	var result types.Ticker
	err := c.transport.DoPublic(ctx, "/public/ticker", q, &result)
	return result, err
}

// Index retrieves the index price for an underlying.
func (c *Client) Index(ctx context.Context, underlying string) (types.IndexPrice, error) {
	q := url.Values{}
	q.Set("underlying", underlying)
	var result types.IndexPrice
	err := c.transport.DoPublic(ctx, "/public/index", q, &result)
	return result, err
}

// Book retrieves the order book for a single instrument.
func (c *Client) Book(ctx context.Context, instrumentName string) (types.Book, error) {
	q := url.Values{}
	q.Set("instrument_name", instrumentName)
	var result types.Book
	err := c.transport.DoPublic(ctx, "/public/book", q, &result)
	return result, err
}

// SystemInfo retrieves system status information.
func (c *Client) SystemInfo(ctx context.Context) (types.SystemInfo, error) {
	var result types.SystemInfo
	err := c.transport.DoPublic(ctx, "/public/system_info", nil, &result)
	return result, err
}
