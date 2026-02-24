package enums

// RfqOrderEvent represents an RFQ order lifecycle event.
type RfqOrderEvent string

const (
	RfqOrderEventInserted  RfqOrderEvent = "Inserted"
	RfqOrderEventAmended   RfqOrderEvent = "Amended"
	RfqOrderEventCancelled RfqOrderEvent = "Cancelled"
	RfqOrderEventFilled    RfqOrderEvent = "Filled"
	RfqOrderEventExisting  RfqOrderEvent = "Existing"
)

// String returns the string representation of the RFQ order event.
func (e RfqOrderEvent) String() string {
	return string(e)
}

// IsValid returns true if the RFQ order event is a recognized value.
func (e RfqOrderEvent) IsValid() bool {
	switch e {
	case RfqOrderEventInserted, RfqOrderEventAmended, RfqOrderEventCancelled, RfqOrderEventFilled, RfqOrderEventExisting:
		return true
	}
	return false
}

// RfqOrderEventValues returns all valid RfqOrderEvent values.
func RfqOrderEventValues() []RfqOrderEvent {
	return []RfqOrderEvent{RfqOrderEventInserted, RfqOrderEventAmended, RfqOrderEventCancelled, RfqOrderEventFilled, RfqOrderEventExisting}
}
