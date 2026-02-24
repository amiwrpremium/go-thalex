package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// CryptoDeposits retrieves deposits via WebSocket.
func (ws *Client) CryptoDeposits(ctx context.Context) (types.DepositsResult, error) {
	var result types.DepositsResult
	err := ws.call(ctx, "private/crypto_deposits", nil, &result)
	return result, err
}

// CryptoWithdrawals retrieves withdrawals via WebSocket.
func (ws *Client) CryptoWithdrawals(ctx context.Context) ([]types.Withdrawal, error) {
	var result []types.Withdrawal
	err := ws.call(ctx, "private/crypto_withdrawals", nil, &result)
	return result, err
}

// BTCDepositAddress retrieves the BTC deposit address via WebSocket.
func (ws *Client) BTCDepositAddress(ctx context.Context) (string, error) {
	var result struct {
		Address string `json:"address"`
	}
	err := ws.call(ctx, "private/btc_deposit_address", nil, &result)
	return result.Address, err
}

// ETHDepositAddress retrieves the ETH deposit address via WebSocket.
func (ws *Client) ETHDepositAddress(ctx context.Context) (string, error) {
	var result struct {
		Address string `json:"address"`
	}
	err := ws.call(ctx, "private/eth_deposit_address", nil, &result)
	return result.Address, err
}

// VerifyWithdrawal verifies a withdrawal via WebSocket.
func (ws *Client) VerifyWithdrawal(ctx context.Context, params *types.WithdrawParams) (types.VerifyWithdrawalResult, error) {
	var result types.VerifyWithdrawalResult
	err := ws.call(ctx, "private/verify_withdrawal", params, &result)
	return result, err
}

// Withdraw initiates a withdrawal via WebSocket.
func (ws *Client) Withdraw(ctx context.Context, params *types.WithdrawParams) (types.Withdrawal, error) {
	var result types.Withdrawal
	err := ws.call(ctx, "private/withdraw", params, &result)
	return result, err
}

// VerifyInternalTransfer verifies a transfer via WebSocket.
func (ws *Client) VerifyInternalTransfer(ctx context.Context, params *types.InternalTransferParams) (types.VerifyInternalTransferResult, error) {
	var result types.VerifyInternalTransferResult
	err := ws.call(ctx, "private/verify_internal_transfer", params, &result)
	return result, err
}

// InternalTransfer executes a transfer via WebSocket.
func (ws *Client) InternalTransfer(ctx context.Context, params *types.InternalTransferParams) error {
	return ws.callNoResult(ctx, "private/internal_transfer", params)
}
