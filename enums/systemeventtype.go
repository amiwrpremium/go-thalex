package enums

// SystemEventType represents a type of system event.
type SystemEventType string

const (
	SystemEventTypeReconnect SystemEventType = "reconnect"
)

// String returns the string representation of the system event type.
func (t SystemEventType) String() string {
	return string(t)
}

// IsValid returns true if the system event type is a recognized value.
func (t SystemEventType) IsValid() bool {
	return t == SystemEventTypeReconnect
}

// SystemEventTypeValues returns all valid SystemEventType values.
func SystemEventTypeValues() []SystemEventType {
	return []SystemEventType{SystemEventTypeReconnect}
}
