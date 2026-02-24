package enums

// MMProtectionReason describes why MM protection was triggered.
type MMProtectionReason string

const (
	MMProtectionReasonTriggered MMProtectionReason = "triggered"
	MMProtectionReasonReset     MMProtectionReason = "reset"
)

// String returns the string representation of the MM protection reason.
func (r MMProtectionReason) String() string {
	return string(r)
}

// IsValid returns true if the MM protection reason is a recognized value.
func (r MMProtectionReason) IsValid() bool {
	switch r {
	case MMProtectionReasonTriggered, MMProtectionReasonReset:
		return true
	}
	return false
}

// MMProtectionReasonValues returns all valid MMProtectionReason values.
func MMProtectionReasonValues() []MMProtectionReason {
	return []MMProtectionReason{MMProtectionReasonTriggered, MMProtectionReasonReset}
}
