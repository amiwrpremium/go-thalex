package enums

// Target represents the trigger type for conditional orders.
type Target string

const (
	TargetLast  Target = "last"
	TargetMark  Target = "mark"
	TargetIndex Target = "index"
)

// String returns the string representation of the target.
func (t Target) String() string {
	return string(t)
}

// IsValid returns true if the target is a recognized value.
func (t Target) IsValid() bool {
	return t == TargetLast || t == TargetMark || t == TargetIndex
}

// TargetValues returns all valid Target values.
func TargetValues() []Target {
	return []Target{TargetLast, TargetMark, TargetIndex}
}
