package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

func TestCryptoDeposits_Success(t *testing.T) {
	deposits := types.DepositsResult{
		Confirmed: []types.Deposit{
			{
				Currency:        "BTC",
				Amount:          1.0,
				Blockchain:      "bitcoin",
				TransactionHash: "abc123",
				Status:          enums.DepositStatusConfirmed,
			},
		},
		Unconfirmed: []types.Deposit{
			{
				Currency:        "ETH",
				Amount:          10.0,
				Blockchain:      "ethereum",
				TransactionHash: "def456",
				Status:          enums.DepositStatusUnconfirmed,
			},
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/crypto_deposits" {
			t.Errorf("expected path /private/crypto_deposits, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, deposits))
	})

	result, err := c.CryptoDeposits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Confirmed) != 1 {
		t.Fatalf("expected 1 confirmed deposit, got %d", len(result.Confirmed))
	}
	if len(result.Unconfirmed) != 1 {
		t.Fatalf("expected 1 unconfirmed deposit, got %d", len(result.Unconfirmed))
	}
	if result.Confirmed[0].Currency != "BTC" {
		t.Errorf("expected BTC, got %s", result.Confirmed[0].Currency)
	}
	if result.Confirmed[0].Amount != 1.0 {
		t.Errorf("expected amount=1.0, got %f", result.Confirmed[0].Amount)
	}
}

func TestCryptoDeposits_Empty(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, types.DepositsResult{
			Confirmed:   []types.Deposit{},
			Unconfirmed: []types.Deposit{},
		}))
	})

	result, err := c.CryptoDeposits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Confirmed) != 0 {
		t.Errorf("expected 0 confirmed, got %d", len(result.Confirmed))
	}
}

func TestCryptoDeposits_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(apiErrorJSON(10001, "unauthorized"))
	})

	_, err := c.CryptoDeposits(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCryptoWithdrawals_Success(t *testing.T) {
	withdrawals := []types.Withdrawal{
		{
			Currency:      "BTC",
			Amount:        0.5,
			TargetAddress: "bc1q...",
			CreateTime:    1700000000.0,
			State:         enums.WithdrawalStatusExecuted,
		},
		{
			Currency:      "ETH",
			Amount:        5.0,
			TargetAddress: "0x...",
			CreateTime:    1700000001.0,
			State:         enums.WithdrawalStatusPending,
		},
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/crypto_withdrawals" {
			t.Errorf("expected path /private/crypto_withdrawals, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, withdrawals))
	})

	result, err := c.CryptoWithdrawals(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 withdrawals, got %d", len(result))
	}
	if result[0].Currency != "BTC" {
		t.Errorf("expected BTC, got %s", result[0].Currency)
	}
	if result[0].State != enums.WithdrawalStatusExecuted {
		t.Errorf("expected state=executed, got %s", result[0].State)
	}
}

func TestCryptoWithdrawals_Empty(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(wrapResult(t, []types.Withdrawal{}))
	})

	result, err := c.CryptoWithdrawals(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}

func TestBTCDepositAddress_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/btc_deposit_address" {
			t.Errorf("expected path /private/btc_deposit_address, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, map[string]string{"address": "bc1qtest123"}))
	})

	addr, err := c.BTCDepositAddress(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "bc1qtest123" {
		t.Errorf("expected bc1qtest123, got %s", addr)
	}
}

func TestETHDepositAddress_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/eth_deposit_address" {
			t.Errorf("expected path /private/eth_deposit_address, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, map[string]string{"address": "0xtest456"}))
	})

	addr, err := c.ETHDepositAddress(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "0xtest456" {
		t.Errorf("expected 0xtest456, got %s", addr)
	}
}

func TestVerifyWithdrawal_Success(t *testing.T) {
	expected := types.VerifyWithdrawalResult{
		AvailableMargin: 50000,
		RequiredMargin:  10000,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/verify_withdrawal" {
			t.Errorf("expected path /private/verify_withdrawal, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var params types.WithdrawParams
		json.Unmarshal(body, &params)
		if params.AssetName != "BTC" {
			t.Errorf("expected asset_name=BTC, got %s", params.AssetName)
		}
		w.Write(wrapResult(t, expected))
	})

	params := &types.WithdrawParams{
		AssetName:     "BTC",
		Amount:        0.1,
		TargetAddress: "bc1qtest",
	}
	result, err := c.VerifyWithdrawal(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AvailableMargin != 50000 {
		t.Errorf("expected available_margin=50000, got %f", result.AvailableMargin)
	}
}

func TestWithdraw_Success(t *testing.T) {
	expected := types.Withdrawal{
		Currency:      "BTC",
		Amount:        0.1,
		TargetAddress: "bc1qtest",
		State:         enums.WithdrawalStatusPending,
		CreateTime:    1700000000.0,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/withdraw" {
			t.Errorf("expected path /private/withdraw, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	params := &types.WithdrawParams{
		AssetName:     "BTC",
		Amount:        0.1,
		TargetAddress: "bc1qtest",
	}
	result, err := c.Withdraw(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.State != enums.WithdrawalStatusPending {
		t.Errorf("expected state=pending, got %s", result.State)
	}
}

func TestVerifyInternalTransfer_Success(t *testing.T) {
	expected := types.VerifyInternalTransferResult{
		SourceAvailableMargin:      40000,
		SourceRequiredMargin:       10000,
		DestinationAvailableMargin: 20000,
		DestinationRequiredMargin:  5000,
	}

	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/verify_internal_transfer" {
			t.Errorf("expected path /private/verify_internal_transfer, got %s", r.URL.Path)
		}
		w.Write(wrapResult(t, expected))
	})

	params := &types.InternalTransferParams{
		DestinationAccountNumber: "ACC-999",
		Assets:                   []types.Asset{{AssetName: "BTC", Amount: 0.1}},
	}
	result, err := c.VerifyInternalTransfer(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SourceAvailableMargin != 40000 {
		t.Errorf("expected source_available_margin=40000, got %f", result.SourceAvailableMargin)
	}
}

func TestInternalTransfer_Success(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private/internal_transfer" {
			t.Errorf("expected path /private/internal_transfer, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"result":null}`))
	})

	params := &types.InternalTransferParams{
		DestinationAccountNumber: "ACC-999",
		Assets:                   []types.Asset{{AssetName: "BTC", Amount: 0.1}},
	}
	err := c.InternalTransfer(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInternalTransfer_APIError(t *testing.T) {
	c := newTestClientWithAuth(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(apiErrorJSON(10030, "insufficient balance"))
	})

	params := &types.InternalTransferParams{
		DestinationAccountNumber: "ACC-999",
	}
	err := c.InternalTransfer(context.Background(), params)
	if err == nil {
		t.Fatal("expected error")
	}
}
