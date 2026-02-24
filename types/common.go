package types

// Ptr returns a pointer to the given value. Useful for optional fields.
func Ptr[T any](v T) *T {
	return &v
}

// Leg represents a leg in a combination order status.
type Leg struct {
	InstrumentName  string  `json:"instrument_name"`
	Quantity        float64 `json:"quantity"`
	FilledAmount    float64 `json:"filled_amount"`
	RemainingAmount float64 `json:"remaining_amount,omitempty"`
}

// InsertLeg represents a leg in a combination order insertion request.
type InsertLeg struct {
	InstrumentName string  `json:"instrument_name"`
	Quantity       float64 `json:"quantity"`
}

// Asset represents an asset in an internal transfer.
type Asset struct {
	AssetName string  `json:"asset_name"`
	Amount    float64 `json:"amount"`
}

// PositionTransfer represents a position in an internal transfer.
type PositionTransfer struct {
	InstrumentName string  `json:"instrument_name"`
	Amount         float64 `json:"amount"`
}
