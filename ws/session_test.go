package ws

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
)

// ---------------------------------------------------------------------------
// SetCancelOnDisconnect
// ---------------------------------------------------------------------------

func TestSetCancelOnDisconnect_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/set_cancel_on_disconnect" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.SetCancelOnDisconnect(ctx, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetCancelOnDisconnect_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "not authenticated"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.SetCancelOnDisconnect(ctx, false)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// CancelSession
// ---------------------------------------------------------------------------

func TestCancelSession_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/cancel_session" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method"}
		}
		return json.RawMessage(`{"n_cancelled":5}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := c.CancelSession(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 cancelled, got %d", n)
	}
}

func TestCancelSession_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "not authenticated"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := c.CancelSession(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	if n != 0 {
		t.Errorf("expected 0 cancelled on error, got %d", n)
	}
}
