package enums

// RfqInsertReason describes why an RFQ was created.
type RfqInsertReason string

const (
	RfqInsertReasonClientRequest RfqInsertReason = "client_request"
	RfqInsertReasonLiquidation   RfqInsertReason = "liquidation"
)

// String returns the string representation of the RFQ insert reason.
func (r RfqInsertReason) String() string {
	return string(r)
}

// IsValid returns true if the RFQ insert reason is a recognized value.
func (r RfqInsertReason) IsValid() bool {
	switch r {
	case RfqInsertReasonClientRequest, RfqInsertReasonLiquidation:
		return true
	}
	return false
}

// RfqInsertReasonValues returns all valid RfqInsertReason values.
func RfqInsertReasonValues() []RfqInsertReason {
	return []RfqInsertReason{RfqInsertReasonClientRequest, RfqInsertReasonLiquidation}
}
