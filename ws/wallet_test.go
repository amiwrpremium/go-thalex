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
// CryptoDeposits
// ---------------------------------------------------------------------------

func TestCryptoDeposits_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/crypto_deposits" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"confirmed":[{"currency":"BTC","amount":1.0,"blockchain":"bitcoin","transaction_hash":"abc","transaction_timestamp":1.0,"status":"completed"}],"unconfirmed":[]}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.CryptoDeposits(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Confirmed) != 1 {
		t.Fatalf("expected 1 confirmed deposit, got %d", len(result.Confirmed))
	}
	if result.Confirmed[0].Currency != "BTC" {
		t.Errorf("expected currency=BTC, got %q", result.Confirmed[0].Currency)
	}
}

func TestCryptoDeposits_APIError(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		return nil, &jsonrpc.Error{Code: 10000, Message: "not authenticated"}
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.CryptoDeposits(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

// ---------------------------------------------------------------------------
// CryptoWithdrawals
// ---------------------------------------------------------------------------

func TestCryptoWithdrawals_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/crypto_withdrawals" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`[{"currency":"BTC","amount":0.5,"target_address":"1abc","create_time":1.0,"state":"completed"}]`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	withdrawals, err := c.CryptoWithdrawals(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(withdrawals) != 1 {
		t.Fatalf("expected 1 withdrawal, got %d", len(withdrawals))
	}
	if withdrawals[0].Amount != 0.5 {
		t.Errorf("expected amount=0.5, got %f", withdrawals[0].Amount)
	}
}

// ---------------------------------------------------------------------------
// BTCDepositAddress
// ---------------------------------------------------------------------------

func TestBTCDepositAddress_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/btc_deposit_address" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"address":"bc1qtest123"}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr, err := c.BTCDepositAddress(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "bc1qtest123" {
		t.Errorf("expected bc1qtest123, got %q", addr)
	}
}

// ---------------------------------------------------------------------------
// ETHDepositAddress
// ---------------------------------------------------------------------------

func TestETHDepositAddress_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/eth_deposit_address" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"address":"0xtest456"}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr, err := c.ETHDepositAddress(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "0xtest456" {
		t.Errorf("expected 0xtest456, got %q", addr)
	}
}

// ---------------------------------------------------------------------------
// VerifyWithdrawal
// ---------------------------------------------------------------------------

func TestVerifyWithdrawal_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/verify_withdrawal" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"available_margin":5000,"required_margin":1000,"error":null}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.VerifyWithdrawal(ctx, &types.WithdrawParams{
		AssetName:     "BTC",
		Amount:        0.1,
		TargetAddress: "bc1qtest",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AvailableMargin != 5000 {
		t.Errorf("expected available_margin=5000, got %f", result.AvailableMargin)
	}
}

// ---------------------------------------------------------------------------
// Withdraw
// ---------------------------------------------------------------------------

func TestWithdraw_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/withdraw" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"currency":"BTC","amount":0.1,"target_address":"bc1qtest","create_time":1.0,"state":"pending"}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.Withdraw(ctx, &types.WithdrawParams{
		AssetName:     "BTC",
		Amount:        0.1,
		TargetAddress: "bc1qtest",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Currency != "BTC" {
		t.Errorf("expected currency=BTC, got %q", result.Currency)
	}
}

// ---------------------------------------------------------------------------
// VerifyInternalTransfer
// ---------------------------------------------------------------------------

func TestVerifyInternalTransfer_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/verify_internal_transfer" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`{"source_available_margin":5000,"source_required_margin":1000,"destination_available_margin":3000,"destination_required_margin":500,"error":null}`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := c.VerifyInternalTransfer(ctx, &types.InternalTransferParams{
		DestinationAccountNumber: "acct-2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SourceAvailableMargin != 5000 {
		t.Errorf("expected source_available_margin=5000, got %f", result.SourceAvailableMargin)
	}
}

// ---------------------------------------------------------------------------
// InternalTransfer
// ---------------------------------------------------------------------------

func TestInternalTransfer_Success(t *testing.T) {
	handler := func(req *jsonrpc.Request) (json.RawMessage, *jsonrpc.Error) {
		if req.Method != "private/internal_transfer" {
			return nil, &jsonrpc.Error{Code: -32601, Message: "unexpected method: " + req.Method}
		}
		return json.RawMessage(`null`), nil
	}
	c := newConnectedClient(t, handler)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.InternalTransfer(ctx, &types.InternalTransferParams{
		DestinationAccountNumber: "acct-2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
