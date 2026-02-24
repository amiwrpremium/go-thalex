package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- Deposit JSON round-trip ----------

func TestDeposit_JSONRoundTrip(t *testing.T) {
	confirmations := 6
	d := types.Deposit{
		Currency:             "BTC",
		Amount:               1.5,
		Blockchain:           "bitcoin",
		TransactionHash:      "abc123",
		TransactionTimestamp: 1700000000.0,
		Status:               enums.DepositStatusConfirmed,
		Confirmations:        &confirmations,
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Deposit
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Currency != "BTC" {
		t.Errorf("Currency = %q, want %q", got.Currency, "BTC")
	}
	if got.Amount != 1.5 {
		t.Errorf("Amount = %v, want 1.5", got.Amount)
	}
	if got.Status != enums.DepositStatusConfirmed {
		t.Errorf("Status = %q, want %q", got.Status, enums.DepositStatusConfirmed)
	}
	if got.Confirmations == nil || *got.Confirmations != 6 {
		t.Errorf("Confirmations = %v, want 6", got.Confirmations)
	}
}

func TestDeposit_NilConfirmations(t *testing.T) {
	d := types.Deposit{
		Currency: "ETH",
		Amount:   10.0,
		Status:   enums.DepositStatusUnconfirmed,
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	if _, ok := raw["confirmations"]; ok {
		t.Error("expected confirmations to be omitted when nil")
	}
}

// ---------- DepositsResult JSON round-trip ----------

func TestDepositsResult_JSONRoundTrip(t *testing.T) {
	r := types.DepositsResult{
		Confirmed: []types.Deposit{
			{Currency: "BTC", Amount: 1.0, Status: enums.DepositStatusConfirmed},
		},
		Unconfirmed: []types.Deposit{
			{Currency: "ETH", Amount: 5.0, Status: enums.DepositStatusUnconfirmed},
		},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.DepositsResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Confirmed) != 1 {
		t.Fatalf("len(Confirmed) = %d, want 1", len(got.Confirmed))
	}
	if len(got.Unconfirmed) != 1 {
		t.Fatalf("len(Unconfirmed) = %d, want 1", len(got.Unconfirmed))
	}
	if got.Confirmed[0].Currency != "BTC" {
		t.Errorf("Confirmed[0].Currency = %q, want %q", got.Confirmed[0].Currency, "BTC")
	}
}

// ---------- Withdrawal JSON round-trip ----------

func TestWithdrawal_JSONRoundTrip(t *testing.T) {
	fee := 0.001
	w := types.Withdrawal{
		Currency:        "BTC",
		Amount:          0.5,
		TargetAddress:   "bc1qaddr",
		Blockchain:      "bitcoin",
		TransactionHash: "txhash123",
		CreateTime:      1700000000.0,
		Label:           "my-withdrawal",
		State:           enums.WithdrawalStatusExecuted,
		Remark:          "done",
		Fee:             &fee,
		FeeAsset:        "BTC",
	}

	data, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Withdrawal
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Currency != "BTC" {
		t.Errorf("Currency = %q, want %q", got.Currency, "BTC")
	}
	if got.TargetAddress != "bc1qaddr" {
		t.Errorf("TargetAddress = %q, want %q", got.TargetAddress, "bc1qaddr")
	}
	if got.State != enums.WithdrawalStatusExecuted {
		t.Errorf("State = %q, want %q", got.State, enums.WithdrawalStatusExecuted)
	}
	if got.Fee == nil || *got.Fee != fee {
		t.Errorf("Fee = %v, want %v", got.Fee, fee)
	}
	if got.FeeAsset != "BTC" {
		t.Errorf("FeeAsset = %q, want %q", got.FeeAsset, "BTC")
	}
}

// ---------- WithdrawParams JSON round-trip ----------

func TestWithdrawParams_JSONRoundTrip(t *testing.T) {
	p := types.WithdrawParams{
		AssetName:     "BTC",
		Amount:        0.5,
		TargetAddress: "bc1qaddr",
		Label:         "my-withdrawal",
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.WithdrawParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.AssetName != "BTC" {
		t.Errorf("AssetName = %q, want %q", got.AssetName, "BTC")
	}
	if got.Amount != 0.5 {
		t.Errorf("Amount = %v, want 0.5", got.Amount)
	}
	if got.TargetAddress != "bc1qaddr" {
		t.Errorf("TargetAddress = %q, want %q", got.TargetAddress, "bc1qaddr")
	}
	if got.Label != "my-withdrawal" {
		t.Errorf("Label = %q, want %q", got.Label, "my-withdrawal")
	}
}

// ---------- VerifyWithdrawalResult JSON round-trip ----------

func TestVerifyWithdrawalResult_JSONRoundTrip(t *testing.T) {
	fee := 0.001
	r := types.VerifyWithdrawalResult{
		AvailableMargin: 10000.0,
		RequiredMargin:  5000.0,
		Fee:             &fee,
		FeeAsset:        "BTC",
		Error:           nil,
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.VerifyWithdrawalResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.AvailableMargin != 10000.0 {
		t.Errorf("AvailableMargin = %v, want 10000.0", got.AvailableMargin)
	}
	if got.Fee == nil || *got.Fee != fee {
		t.Errorf("Fee = %v, want %v", got.Fee, fee)
	}
}

func TestVerifyWithdrawalResult_WithError(t *testing.T) {
	r := types.VerifyWithdrawalResult{
		AvailableMargin: 100.0,
		RequiredMargin:  5000.0,
		Error:           &apierr.APIError{Code: 1000, Message: "insufficient_margin"},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.VerifyWithdrawalResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Error == nil {
		t.Fatal("Error should not be nil")
	}
	if got.Error.Code != 1000 {
		t.Errorf("Error.Code = %d, want 1000", got.Error.Code)
	}
	if got.Error.Message != "insufficient_margin" {
		t.Errorf("Error.Message = %q, want %q", got.Error.Message, "insufficient_margin")
	}
}

// ---------- InternalTransferParams JSON round-trip ----------

func TestInternalTransferParams_JSONRoundTrip(t *testing.T) {
	p := types.InternalTransferParams{
		DestinationAccountNumber: "ACC-001",
		Assets: []types.Asset{
			{AssetName: "BTC", Amount: 1.0},
			{AssetName: "ETH", Amount: 10.0},
		},
		Positions: []types.PositionTransfer{
			{InstrumentName: "BTC-PERPETUAL", Amount: 0.5},
		},
		Label: "internal-xfer",
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.InternalTransferParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.DestinationAccountNumber != "ACC-001" {
		t.Errorf("DestinationAccountNumber = %q, want %q", got.DestinationAccountNumber, "ACC-001")
	}
	if len(got.Assets) != 2 {
		t.Fatalf("len(Assets) = %d, want 2", len(got.Assets))
	}
	if got.Assets[0].AssetName != "BTC" {
		t.Errorf("Assets[0].AssetName = %q, want %q", got.Assets[0].AssetName, "BTC")
	}
	if len(got.Positions) != 1 {
		t.Fatalf("len(Positions) = %d, want 1", len(got.Positions))
	}
	if got.Positions[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("Positions[0].InstrumentName = %q, want %q", got.Positions[0].InstrumentName, "BTC-PERPETUAL")
	}
}

func TestInternalTransferParams_EmptyLists(t *testing.T) {
	p := types.InternalTransferParams{
		DestinationAccountNumber: "ACC-002",
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Assets and positions with omitempty should be omitted
	if _, ok := raw["assets"]; ok {
		t.Error("expected assets to be omitted when nil")
	}
	if _, ok := raw["positions"]; ok {
		t.Error("expected positions to be omitted when nil")
	}
}

// ---------- VerifyInternalTransferResult JSON ----------

func TestVerifyInternalTransferResult_JSONRoundTrip(t *testing.T) {
	r := types.VerifyInternalTransferResult{
		SourceAvailableMargin:      10000.0,
		SourceRequiredMargin:       5000.0,
		DestinationAvailableMargin: 20000.0,
		DestinationRequiredMargin:  8000.0,
		Error:                      nil,
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.VerifyInternalTransferResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.SourceAvailableMargin != 10000.0 {
		t.Errorf("SourceAvailableMargin = %v, want 10000.0", got.SourceAvailableMargin)
	}
	if got.DestinationRequiredMargin != 8000.0 {
		t.Errorf("DestinationRequiredMargin = %v, want 8000.0", got.DestinationRequiredMargin)
	}
	if got.Error != nil {
		t.Errorf("Error = %v, want nil", got.Error)
	}
}
