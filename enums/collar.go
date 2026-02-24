package enums

// Collar determines how price collar violations are handled.
type Collar string

const (
	CollarIgnore Collar = "ignore"
	CollarReject Collar = "reject"
	CollarClamp  Collar = "clamp"
)

// String returns the string representation of the collar.
func (c Collar) String() string {
	return string(c)
}

// IsValid returns true if the collar is a recognized value.
func (c Collar) IsValid() bool {
	return c == CollarIgnore || c == CollarReject || c == CollarClamp
}

// CollarValues returns all valid Collar values.
func CollarValues() []Collar {
	return []Collar{CollarIgnore, CollarReject, CollarClamp}
}
