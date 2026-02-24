package enums

// DisplayType specifies notification display style.
type DisplayType string

const (
	DisplayTypeSuccess  DisplayType = "success"
	DisplayTypeFailure  DisplayType = "failure"
	DisplayTypeInfo     DisplayType = "info"
	DisplayTypeWarning  DisplayType = "warning"
	DisplayTypeCritical DisplayType = "critical"
)

// String returns the string representation of the display type.
func (d DisplayType) String() string {
	return string(d)
}

// IsValid returns true if the display type is a recognized value.
func (d DisplayType) IsValid() bool {
	switch d {
	case DisplayTypeSuccess, DisplayTypeFailure, DisplayTypeInfo, DisplayTypeWarning, DisplayTypeCritical:
		return true
	}
	return false
}

// DisplayTypeValues returns all valid DisplayType values.
func DisplayTypeValues() []DisplayType {
	return []DisplayType{DisplayTypeSuccess, DisplayTypeFailure, DisplayTypeInfo, DisplayTypeWarning, DisplayTypeCritical}
}
