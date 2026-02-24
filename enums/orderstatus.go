package enums

// OrderStatusValue represents the status of an order.
type OrderStatusValue string

const (
	OrderStatusOpen                     OrderStatusValue = "open"
	OrderStatusPartiallyFilled          OrderStatusValue = "partially_filled"
	OrderStatusCancelled                OrderStatusValue = "cancelled"
	OrderStatusCancelledPartiallyFilled OrderStatusValue = "cancelled_partially_filled"
	OrderStatusFilled                   OrderStatusValue = "filled"
)

// String returns the string representation of the order status.
func (s OrderStatusValue) String() string {
	return string(s)
}

// IsValid returns true if the order status is a recognized value.
func (s OrderStatusValue) IsValid() bool {
	switch s {
	case OrderStatusOpen, OrderStatusPartiallyFilled, OrderStatusCancelled, OrderStatusCancelledPartiallyFilled, OrderStatusFilled:
		return true
	}
	return false
}

// OrderStatusValues returns all valid OrderStatusValue values.
func OrderStatusValues() []OrderStatusValue {
	return []OrderStatusValue{OrderStatusOpen, OrderStatusPartiallyFilled, OrderStatusCancelled, OrderStatusCancelledPartiallyFilled, OrderStatusFilled}
}

// IsActive returns true if the order is still active in the book.
func (s OrderStatusValue) IsActive() bool {
	return s == OrderStatusOpen || s == OrderStatusPartiallyFilled
}

// IsFinal returns true if the order has reached a terminal state.
func (s OrderStatusValue) IsFinal() bool {
	return !s.IsActive()
}
