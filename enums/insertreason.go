package enums

// InsertReason describes why an order was inserted.
type InsertReason string

const (
	InsertReasonClientRequest    InsertReason = "client_request"
	InsertReasonConditionalOrder InsertReason = "conditional_order"
	InsertReasonLiquidation      InsertReason = "liquidation"
)

// String returns the string representation of the insert reason.
func (r InsertReason) String() string {
	return string(r)
}

// IsValid returns true if the insert reason is a recognized value.
func (r InsertReason) IsValid() bool {
	switch r {
	case InsertReasonClientRequest, InsertReasonConditionalOrder, InsertReasonLiquidation:
		return true
	}
	return false
}

// InsertReasonValues returns all valid InsertReason values.
func InsertReasonValues() []InsertReason {
	return []InsertReason{InsertReasonClientRequest, InsertReasonConditionalOrder, InsertReasonLiquidation}
}
