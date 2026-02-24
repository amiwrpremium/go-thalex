package types

import (
	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/enums"
)

// Deposit represents a cryptocurrency deposit.
type Deposit struct {
	Currency             string              `json:"currency"`
	Amount               float64             `json:"amount"`
	Blockchain           string              `json:"blockchain"`
	TransactionHash      string              `json:"transaction_hash"`
	TransactionTimestamp float64             `json:"transaction_timestamp"`
	Status               enums.DepositStatus `json:"status"`
	Confirmations        *int                `json:"confirmations,omitempty"`
}

// DepositsResult contains confirmed and unconfirmed deposits.
type DepositsResult struct {
	Confirmed   []Deposit `json:"confirmed"`
	Unconfirmed []Deposit `json:"unconfirmed"`
}

// Withdrawal represents a cryptocurrency withdrawal.
type Withdrawal struct {
	Currency        string                 `json:"currency"`
	Amount          float64                `json:"amount"`
	TargetAddress   string                 `json:"target_address"`
	Blockchain      string                 `json:"blockchain,omitempty"`
	TransactionHash string                 `json:"transaction_hash,omitempty"`
	CreateTime      float64                `json:"create_time"`
	Label           string                 `json:"label,omitempty"`
	State           enums.WithdrawalStatus `json:"state"`
	Remark          string                 `json:"remark,omitempty"`
	Fee             *float64               `json:"fee,omitempty"`
	FeeAsset        string                 `json:"fee_asset,omitempty"`
}

// WithdrawParams contains parameters for a withdrawal request.
type WithdrawParams struct {
	AssetName     string  `json:"asset_name"`
	Amount        float64 `json:"amount"`
	TargetAddress string  `json:"target_address"`
	Label         string  `json:"label,omitempty"`
}

// VerifyWithdrawalResult contains the result of a withdrawal verification.
type VerifyWithdrawalResult struct {
	AvailableMargin float64          `json:"available_margin"`
	RequiredMargin  float64          `json:"required_margin"`
	Fee             *float64         `json:"fee,omitempty"`
	FeeAsset        string           `json:"fee_asset,omitempty"`
	Error           *apierr.APIError `json:"error"`
}

// InternalTransferParams contains parameters for an internal transfer.
type InternalTransferParams struct {
	DestinationAccountNumber string             `json:"destination_account_number"`
	Assets                   []Asset            `json:"assets,omitempty"`
	Positions                []PositionTransfer `json:"positions,omitempty"`
	Label                    string             `json:"label,omitempty"`
}

// VerifyInternalTransferResult contains the result of a transfer verification.
type VerifyInternalTransferResult struct {
	SourceAvailableMargin      float64          `json:"source_available_margin"`
	SourceRequiredMargin       float64          `json:"source_required_margin"`
	DestinationAvailableMargin float64          `json:"destination_available_margin"`
	DestinationRequiredMargin  float64          `json:"destination_required_margin"`
	Error                      *apierr.APIError `json:"error"`
}
