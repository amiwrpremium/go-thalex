package enums

// Direction represents the side of a trade or order.
type Direction string

const (
	DirectionBuy  Direction = "buy"
	DirectionSell Direction = "sell"
)

// IsValid returns true if the direction is a recognized value.
func (d Direction) IsValid() bool {
	return d == DirectionBuy || d == DirectionSell
}

// String returns the string representation of the direction.
func (d Direction) String() string {
	return string(d)
}

// DirectionValues returns all valid Direction values.
func DirectionValues() []Direction {
	return []Direction{DirectionBuy, DirectionSell}
}

// Opposite returns the opposite direction.
func (d Direction) Opposite() Direction {
	if d == DirectionBuy {
		return DirectionSell
	}
	return DirectionBuy
}
