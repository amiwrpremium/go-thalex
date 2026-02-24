package ws

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/internal/jsonrpc"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------------------------------------------------------------------------
// MassQuote
// ---------------------------------------------------------------------------

func TestMassQuote_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/mass_quote" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"n_success":2,"n_fail":0,"errors":[]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.MassQuote(ctx, &types.MassQuoteParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NSuccess != 2 {
		t.Errorf("expected n_success=2, got %d", result.NSuccess)
	}
	if result.NFail != 0 {
		t.Errorf("expected n_fail=0, got %d", result.NFail)
	}
}

func TestMassQuote_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10001, Message: "insufficient margin"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.MassQuote(ctx, &types.MassQuoteParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// CancelMassQuote
// ---------------------------------------------------------------------------

func TestCancelMassQuote_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/cancel_mass_quote" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.CancelMassQuote(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// SetMMProtection
// ---------------------------------------------------------------------------

func TestSetMMProtection_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/set_mm_protection" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.SetMMProtection(ctx, &types.MMProtectionParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetMMProtection_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "not authenticated"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.SetMMProtection(ctx, &types.MMProtectionParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
