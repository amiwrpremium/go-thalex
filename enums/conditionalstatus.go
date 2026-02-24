package enums

// ConditionalOrderStatus represents the status of a conditional order.
type ConditionalOrderStatus string

const (
	ConditionalOrderStatusCreated         ConditionalOrderStatus = "created"
	ConditionalOrderStatusActive          ConditionalOrderStatus = "active"
	ConditionalOrderStatusConverted       ConditionalOrderStatus = "converted"
	ConditionalOrderStatusRejected        ConditionalOrderStatus = "rejected"
	ConditionalOrderStatusCancelRequested ConditionalOrderStatus = "cancel requested"
	ConditionalOrderStatusCancelled       ConditionalOrderStatus = "cancelled"
)

// String returns the string representation of the conditional order status.
func (s ConditionalOrderStatus) String() string {
	return string(s)
}

// IsValid returns true if the conditional order status is a recognized value.
func (s ConditionalOrderStatus) IsValid() bool {
	switch s {
	case ConditionalOrderStatusCreated, ConditionalOrderStatusActive, ConditionalOrderStatusConverted, ConditionalOrderStatusRejected, ConditionalOrderStatusCancelRequested, ConditionalOrderStatusCancelled:
		return true
	}
	return false
}

// ConditionalOrderStatusValues returns all valid ConditionalOrderStatus values.
func ConditionalOrderStatusValues() []ConditionalOrderStatus {
	return []ConditionalOrderStatus{ConditionalOrderStatusCreated, ConditionalOrderStatusActive, ConditionalOrderStatusConverted, ConditionalOrderStatusRejected, ConditionalOrderStatusCancelRequested, ConditionalOrderStatusCancelled}
}

// IsActive returns true if the conditional order is still active.
func (s ConditionalOrderStatus) IsActive() bool {
	return s == ConditionalOrderStatusCreated || s == ConditionalOrderStatusActive
}
