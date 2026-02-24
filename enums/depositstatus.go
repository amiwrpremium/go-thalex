package enums

// DepositStatus represents the status of a deposit.
type DepositStatus string

const (
	DepositStatusUnconfirmed DepositStatus = "unconfirmed"
	DepositStatusConfirmed   DepositStatus = "confirmed"
)

// String returns the string representation of the deposit status.
func (s DepositStatus) String() string {
	return string(s)
}

// IsValid returns true if the deposit status is a recognized value.
func (s DepositStatus) IsValid() bool {
	switch s {
	case DepositStatusUnconfirmed, DepositStatusConfirmed:
		return true
	}
	return false
}

// DepositStatusValues returns all valid DepositStatus values.
func DepositStatusValues() []DepositStatus {
	return []DepositStatus{DepositStatusUnconfirmed, DepositStatusConfirmed}
}

// IsPending returns true if the deposit has not yet been confirmed.
func (s DepositStatus) IsPending() bool {
	return s == DepositStatusUnconfirmed
}

// IsFinal returns true if the deposit has reached a terminal state.
func (s DepositStatus) IsFinal() bool {
	return s == DepositStatusConfirmed
}
