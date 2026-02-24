package enums

// Delay represents the minimum interval between subscription feeds.
type Delay string

const (
	DelayNone   Delay = "raw"
	Delay100ms  Delay = "100ms"
	Delay1000ms Delay = "1000ms"
)

// String returns the string representation of the delay.
func (d Delay) String() string {
	return string(d)
}

// IsValid returns true if the delay is a recognized value.
func (d Delay) IsValid() bool {
	switch d {
	case DelayNone, Delay100ms, Delay1000ms:
		return true
	}
	return false
}

// DelayValues returns all valid Delay values.
func DelayValues() []Delay {
	return []Delay{DelayNone, Delay100ms, Delay1000ms}
}
