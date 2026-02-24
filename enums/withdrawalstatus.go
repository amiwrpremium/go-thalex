package enums

// WithdrawalStatus represents the status of a withdrawal.
type WithdrawalStatus string

const (
	WithdrawalStatusPending              WithdrawalStatus = "pending"
	WithdrawalStatusAwaitingConfirmation WithdrawalStatus = "awaiting_confirmation"
	WithdrawalStatusExecuting            WithdrawalStatus = "executing"
	WithdrawalStatusExecuted             WithdrawalStatus = "executed"
	WithdrawalStatusRejected             WithdrawalStatus = "rejected"
)

// String returns the string representation of the withdrawal status.
func (s WithdrawalStatus) String() string {
	return string(s)
}

// IsValid returns true if the withdrawal status is a recognized value.
func (s WithdrawalStatus) IsValid() bool {
	switch s {
	case WithdrawalStatusPending, WithdrawalStatusAwaitingConfirmation, WithdrawalStatusExecuting, WithdrawalStatusExecuted, WithdrawalStatusRejected:
		return true
	}
	return false
}

// WithdrawalStatusValues returns all valid WithdrawalStatus values.
func WithdrawalStatusValues() []WithdrawalStatus {
	return []WithdrawalStatus{WithdrawalStatusPending, WithdrawalStatusAwaitingConfirmation, WithdrawalStatusExecuting, WithdrawalStatusExecuted, WithdrawalStatusRejected}
}

// IsPending returns true if the withdrawal is still in progress.
func (s WithdrawalStatus) IsPending() bool {
	return s == WithdrawalStatusPending || s == WithdrawalStatusAwaitingConfirmation || s == WithdrawalStatusExecuting
}

// IsFinal returns true if the withdrawal has reached a terminal state.
func (s WithdrawalStatus) IsFinal() bool {
	return s == WithdrawalStatusExecuted || s == WithdrawalStatusRejected
}
