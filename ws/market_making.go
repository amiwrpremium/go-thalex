package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// MassQuote sends a mass quote (WebSocket-only).
func (ws *Client) MassQuote(ctx context.Context, params *types.MassQuoteParams) (types.DoubleSidedQuoteResult, error) {
	var result types.DoubleSidedQuoteResult
	err := ws.call(ctx, "private/mass_quote", params, &result)
	return result, err
}

// CancelMassQuote cancels all mass quotes (WebSocket-only).
func (ws *Client) CancelMassQuote(ctx context.Context) error {
	return ws.callNoResult(ctx, "private/cancel_mass_quote", nil)
}

// SetMMProtection configures market maker protection (WebSocket-only).
func (ws *Client) SetMMProtection(ctx context.Context, params *types.MMProtectionParams) error {
	return ws.callNoResult(ctx, "private/set_mm_protection", params)
}
