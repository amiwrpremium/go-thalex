package enums

// MakerTaker indicates whether a trade was on the maker or taker side.
type MakerTaker string

const (
	MakerTakerMaker MakerTaker = "maker"
	MakerTakerTaker MakerTaker = "taker"
)

// String returns the string representation of the maker/taker side.
func (m MakerTaker) String() string {
	return string(m)
}

// IsValid returns true if the maker/taker value is a recognized value.
func (m MakerTaker) IsValid() bool {
	switch m {
	case MakerTakerMaker, MakerTakerTaker:
		return true
	}
	return false
}

// MakerTakerValues returns all valid MakerTaker values.
func MakerTakerValues() []MakerTaker {
	return []MakerTaker{MakerTakerMaker, MakerTakerTaker}
}
