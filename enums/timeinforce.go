package enums

// TimeInForce controls the lifetime of an order.
type TimeInForce string

const (
	TimeInForceGoodTillCancelled TimeInForce = "good_till_cancelled"
	TimeInForceImmediateOrCancel TimeInForce = "immediate_or_cancel"
)

// String returns the string representation of the time in force.
func (t TimeInForce) String() string {
	return string(t)
}

// IsValid returns true if the time in force is a recognized value.
func (t TimeInForce) IsValid() bool {
	return t == TimeInForceGoodTillCancelled || t == TimeInForceImmediateOrCancel
}

// TimeInForceValues returns all valid TimeInForce values.
func TimeInForceValues() []TimeInForce {
	return []TimeInForce{TimeInForceGoodTillCancelled, TimeInForceImmediateOrCancel}
}
