package enums

// Severity represents the severity level of a banner or notification.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// String returns the string representation of the severity.
func (s Severity) String() string {
	return string(s)
}

// IsValid returns true if the severity is a recognized value.
func (s Severity) IsValid() bool {
	switch s {
	case SeverityInfo, SeverityWarning, SeverityCritical:
		return true
	}
	return false
}

// SeverityValues returns all valid Severity values.
func SeverityValues() []Severity {
	return []Severity{SeverityInfo, SeverityWarning, SeverityCritical}
}
