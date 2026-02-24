package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/apierr"
)

// Login authenticates the WebSocket session using the configured credentials.
func (ws *Client) Login(ctx context.Context) error {
	if ws.cfg.Credentials == nil {
		return &apierr.AuthError{Message: "no credentials configured"}
	}
	token, err := ws.cfg.Credentials.GenerateToken()
	if err != nil {
		return err
	}
	params := map[string]any{"token": token}
	if ws.cfg.AccountNumber != "" {
		params["account"] = ws.cfg.AccountNumber
	}
	return ws.callNoResult(ctx, "public/login", params)
}

// SetCancelOnDisconnect enables or disables cancel-on-disconnect for the session.
func (ws *Client) SetCancelOnDisconnect(ctx context.Context, enabled bool) error {
	return ws.callNoResult(ctx, "private/set_cancel_on_disconnect", map[string]any{
		"value": enabled,
	})
}

// CancelSession cancels all non-persistent orders in the current session.
func (ws *Client) CancelSession(ctx context.Context) (int, error) {
	var result struct {
		NCancelled int `json:"n_cancelled"`
	}
	err := ws.call(ctx, "private/cancel_session", nil, &result)
	return result.NCancelled, err
}
