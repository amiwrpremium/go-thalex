package enums

// OptionType represents the type of an option.
type OptionType string

const (
	OptionTypeCall OptionType = "call"
	OptionTypePut  OptionType = "put"
)

// String returns the string representation of the option type.
func (t OptionType) String() string {
	return string(t)
}

// IsValid returns true if the option type is a recognized value.
func (t OptionType) IsValid() bool {
	switch t {
	case OptionTypeCall, OptionTypePut:
		return true
	}
	return false
}

// OptionTypeValues returns all valid OptionType values.
func OptionTypeValues() []OptionType {
	return []OptionType{OptionTypeCall, OptionTypePut}
}
