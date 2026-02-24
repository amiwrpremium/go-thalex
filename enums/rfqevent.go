package enums

// RfqEvent represents an RFQ lifecycle event.
type RfqEvent string

const (
	RfqEventCreated   RfqEvent = "Created"
	RfqEventCancelled RfqEvent = "Cancelled"
	RfqEventTraded    RfqEvent = "Traded"
	RfqEventExisting  RfqEvent = "Existing"
)

// String returns the string representation of the RFQ event.
func (e RfqEvent) String() string {
	return string(e)
}

// IsValid returns true if the RFQ event is a recognized value.
func (e RfqEvent) IsValid() bool {
	switch e {
	case RfqEventCreated, RfqEventCancelled, RfqEventTraded, RfqEventExisting:
		return true
	}
	return false
}

// RfqEventValues returns all valid RfqEvent values.
func RfqEventValues() []RfqEvent {
	return []RfqEvent{RfqEventCreated, RfqEventCancelled, RfqEventTraded, RfqEventExisting}
}
