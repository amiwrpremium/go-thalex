package enums

// STPAction represents the self-trade prevention action.
type STPAction string

const (
	STPActionCancelAggressor STPAction = "cancel_aggressor"
	STPActionCancelPassive   STPAction = "cancel_passive"
	STPActionCancelBoth      STPAction = "cancel_both"
)

// String returns the string representation of the STP action.
func (a STPAction) String() string {
	return string(a)
}

// IsValid returns true if the STP action is a recognized value.
func (a STPAction) IsValid() bool {
	switch a {
	case STPActionCancelAggressor, STPActionCancelPassive, STPActionCancelBoth:
		return true
	}
	return false
}

// STPActionValues returns all valid STPAction values.
func STPActionValues() []STPAction {
	return []STPAction{STPActionCancelAggressor, STPActionCancelPassive, STPActionCancelBoth}
}
