package rest

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// CryptoDeposits retrieves confirmed and unconfirmed deposits.
func (c *Client) CryptoDeposits(ctx context.Context) (types.DepositsResult, error) {
	var result types.DepositsResult
	err := c.transport.DoPrivateGET(ctx, "/private/crypto_deposits", nil, &result)
	return result, err
}

// CryptoWithdrawals retrieves all withdrawals.
func (c *Client) CryptoWithdrawals(ctx context.Context) ([]types.Withdrawal, error) {
	var result []types.Withdrawal
	err := c.transport.DoPrivateGET(ctx, "/private/crypto_withdrawals", nil, &result)
	return result, err
}

type addressResult struct {
	Address string `json:"address"`
}

// BTCDepositAddress retrieves the BTC deposit address.
func (c *Client) BTCDepositAddress(ctx context.Context) (string, error) {
	var result addressResult
	err := c.transport.DoPrivateGET(ctx, "/private/btc_deposit_address", nil, &result)
	return result.Address, err
}

// ETHDepositAddress retrieves the ETH deposit address.
func (c *Client) ETHDepositAddress(ctx context.Context) (string, error) {
	var result addressResult
	err := c.transport.DoPrivateGET(ctx, "/private/eth_deposit_address", nil, &result)
	return result.Address, err
}

// VerifyWithdrawal verifies a withdrawal without executing it.
func (c *Client) VerifyWithdrawal(ctx context.Context, params *types.WithdrawParams) (types.VerifyWithdrawalResult, error) {
	var result types.VerifyWithdrawalResult
	err := c.transport.DoPrivatePOST(ctx, "/private/verify_withdrawal", params, &result)
	return result, err
}

// Withdraw initiates a cryptocurrency withdrawal.
func (c *Client) Withdraw(ctx context.Context, params *types.WithdrawParams) (types.Withdrawal, error) {
	var result types.Withdrawal
	err := c.transport.DoPrivatePOST(ctx, "/private/withdraw", params, &result)
	return result, err
}

// VerifyInternalTransfer verifies an internal transfer without executing it.
func (c *Client) VerifyInternalTransfer(ctx context.Context, params *types.InternalTransferParams) (types.VerifyInternalTransferResult, error) {
	var result types.VerifyInternalTransferResult
	err := c.transport.DoPrivatePOST(ctx, "/private/verify_internal_transfer", params, &result)
	return result, err
}

// InternalTransfer executes an internal transfer to another account.
func (c *Client) InternalTransfer(ctx context.Context, params *types.InternalTransferParams) error {
	return c.transport.DoPrivatePOST(ctx, "/private/internal_transfer", params, nil)
}
