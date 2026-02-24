package enums

// ChangeReason describes why an order status changed.
type ChangeReason string

const (
	ChangeReasonExisting ChangeReason = "existing"
	ChangeReasonInsert   ChangeReason = "insert"
	ChangeReasonAmend    ChangeReason = "amend"
	ChangeReasonCancel   ChangeReason = "cancel"
	ChangeReasonFill     ChangeReason = "fill"
)

// String returns the string representation of the change reason.
func (r ChangeReason) String() string {
	return string(r)
}

// IsValid returns true if the change reason is a recognized value.
func (r ChangeReason) IsValid() bool {
	switch r {
	case ChangeReasonExisting, ChangeReasonInsert, ChangeReasonAmend, ChangeReasonCancel, ChangeReasonFill:
		return true
	}
	return false
}

// ChangeReasonValues returns all valid ChangeReason values.
func ChangeReasonValues() []ChangeReason {
	return []ChangeReason{ChangeReasonExisting, ChangeReasonInsert, ChangeReasonAmend, ChangeReasonCancel, ChangeReasonFill}
}
