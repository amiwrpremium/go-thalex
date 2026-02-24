package enums

// RfqDeleteReason describes why an RFQ order was deleted.
type RfqDeleteReason string

const (
	RfqDeleteReasonClientCancel          RfqDeleteReason = "client_cancel"
	RfqDeleteReasonSessionEnd            RfqDeleteReason = "session_end"
	RfqDeleteReasonInstrumentDeactivated RfqDeleteReason = "instrument_deactivated"
	RfqDeleteReasonMMProtection          RfqDeleteReason = "mm_protection"
	RfqDeleteReasonFailover              RfqDeleteReason = "failover"
	RfqDeleteReasonMarginBreach          RfqDeleteReason = "margin_breach"
	RfqDeleteReasonFilled                RfqDeleteReason = "filled"
)

// String returns the string representation of the RFQ delete reason.
func (r RfqDeleteReason) String() string {
	return string(r)
}

// IsValid returns true if the RFQ delete reason is a recognized value.
func (r RfqDeleteReason) IsValid() bool {
	switch r {
	case RfqDeleteReasonClientCancel, RfqDeleteReasonSessionEnd, RfqDeleteReasonInstrumentDeactivated, RfqDeleteReasonMMProtection, RfqDeleteReasonFailover, RfqDeleteReasonMarginBreach, RfqDeleteReasonFilled:
		return true
	}
	return false
}

// RfqDeleteReasonValues returns all valid RfqDeleteReason values.
func RfqDeleteReasonValues() []RfqDeleteReason {
	return []RfqDeleteReason{RfqDeleteReasonClientCancel, RfqDeleteReasonSessionEnd, RfqDeleteReasonInstrumentDeactivated, RfqDeleteReasonMMProtection, RfqDeleteReasonFailover, RfqDeleteReasonMarginBreach, RfqDeleteReasonFilled}
}
