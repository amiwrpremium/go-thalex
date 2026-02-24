package enums

// OrderType represents the type of an order.
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// String returns the string representation of the order type.
func (t OrderType) String() string {
	return string(t)
}

// IsValid returns true if the order type is a recognized value.
func (t OrderType) IsValid() bool {
	return t == OrderTypeLimit || t == OrderTypeMarket
}

// OrderTypeValues returns all valid OrderType values.
func OrderTypeValues() []OrderType {
	return []OrderType{OrderTypeLimit, OrderTypeMarket}
}
