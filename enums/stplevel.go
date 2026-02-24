package enums

// STPLevel represents the self-trade prevention level.
type STPLevel string

const (
	STPLevelAccount    STPLevel = "account"
	STPLevelCustomer   STPLevel = "customer"
	STPLevelSubaccount STPLevel = "subaccount"
)

// String returns the string representation of the STP level.
func (l STPLevel) String() string {
	return string(l)
}

// IsValid returns true if the STP level is a recognized value.
func (l STPLevel) IsValid() bool {
	switch l {
	case STPLevelAccount, STPLevelCustomer, STPLevelSubaccount:
		return true
	}
	return false
}

// STPLevelValues returns all valid STPLevel values.
func STPLevelValues() []STPLevel {
	return []STPLevel{STPLevelAccount, STPLevelCustomer, STPLevelSubaccount}
}
