package enums

// TradeType represents the type of a trade.
type TradeType string

const (
	TradeTypeNormal           TradeType = "normal"
	TradeTypeBlock            TradeType = "block"
	TradeTypeCombo            TradeType = "combo"
	TradeTypeAmend            TradeType = "amend"
	TradeTypeDelete           TradeType = "delete"
	TradeTypeInternalTransfer TradeType = "internal_transfer"
	TradeTypeExpiration       TradeType = "expiration"
	TradeTypeDailyMark        TradeType = "daily_mark"
	TradeTypeRfq              TradeType = "rfq"
	TradeTypeLiquidation      TradeType = "liquidation"
)

// String returns the string representation of the trade type.
func (t TradeType) String() string {
	return string(t)
}

// IsValid returns true if the trade type is a recognized value.
func (t TradeType) IsValid() bool {
	switch t {
	case TradeTypeNormal, TradeTypeBlock, TradeTypeCombo, TradeTypeAmend, TradeTypeDelete, TradeTypeInternalTransfer, TradeTypeExpiration, TradeTypeDailyMark, TradeTypeRfq, TradeTypeLiquidation:
		return true
	}
	return false
}

// TradeTypeValues returns all valid TradeType values.
func TradeTypeValues() []TradeType {
	return []TradeType{TradeTypeNormal, TradeTypeBlock, TradeTypeCombo, TradeTypeAmend, TradeTypeDelete, TradeTypeInternalTransfer, TradeTypeExpiration, TradeTypeDailyMark, TradeTypeRfq, TradeTypeLiquidation}
}
