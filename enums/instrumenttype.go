package enums

// InstrumentType represents the type of a trading instrument.
type InstrumentType string

const (
	InstrumentTypePerpetual   InstrumentType = "perpetual"
	InstrumentTypeFuture      InstrumentType = "future"
	InstrumentTypeOption      InstrumentType = "option"
	InstrumentTypeCombination InstrumentType = "combination"
)

// String returns the string representation of the instrument type.
func (t InstrumentType) String() string {
	return string(t)
}

// IsValid returns true if the instrument type is a recognized value.
func (t InstrumentType) IsValid() bool {
	switch t {
	case InstrumentTypePerpetual, InstrumentTypeFuture, InstrumentTypeOption, InstrumentTypeCombination:
		return true
	}
	return false
}

// InstrumentTypeValues returns all valid InstrumentType values.
func InstrumentTypeValues() []InstrumentType {
	return []InstrumentType{InstrumentTypePerpetual, InstrumentTypeFuture, InstrumentTypeOption, InstrumentTypeCombination}
}
