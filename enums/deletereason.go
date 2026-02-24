package enums

// DeleteReason describes why an order was deleted.
type DeleteReason string

const (
	DeleteReasonClientCancel          DeleteReason = "client_cancel"
	DeleteReasonClientBulkCancel      DeleteReason = "client_bulk_cancel"
	DeleteReasonSessionEnd            DeleteReason = "session_end"
	DeleteReasonInstrumentDeactivated DeleteReason = "instrument_deactivated"
	DeleteReasonMMProtection          DeleteReason = "mm_protection"
	DeleteReasonFailover              DeleteReason = "failover"
	DeleteReasonMarginBreach          DeleteReason = "margin_breach"
	DeleteReasonFilled                DeleteReason = "filled"
	DeleteReasonImmediateCancel       DeleteReason = "immediate_cancel"
	DeleteReasonAdminCancel           DeleteReason = "admin_cancel"
)

// String returns the string representation of the delete reason.
func (r DeleteReason) String() string {
	return string(r)
}

// IsValid returns true if the delete reason is a recognized value.
func (r DeleteReason) IsValid() bool {
	switch r {
	case DeleteReasonClientCancel, DeleteReasonClientBulkCancel, DeleteReasonSessionEnd, DeleteReasonInstrumentDeactivated, DeleteReasonMMProtection, DeleteReasonFailover, DeleteReasonMarginBreach, DeleteReasonFilled, DeleteReasonImmediateCancel, DeleteReasonAdminCancel:
		return true
	}
	return false
}

// DeleteReasonValues returns all valid DeleteReason values.
func DeleteReasonValues() []DeleteReason {
	return []DeleteReason{DeleteReasonClientCancel, DeleteReasonClientBulkCancel, DeleteReasonSessionEnd, DeleteReasonInstrumentDeactivated, DeleteReasonMMProtection, DeleteReasonFailover, DeleteReasonMarginBreach, DeleteReasonFilled, DeleteReasonImmediateCancel, DeleteReasonAdminCancel}
}
