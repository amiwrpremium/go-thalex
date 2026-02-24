package enums

// Product represents a product group (e.g. "FBTCUSD", "OBTCUSD").
type Product string

// String returns the string representation of the product.
func (p Product) String() string {
	return string(p)
}
