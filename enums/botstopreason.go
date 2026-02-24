package enums

// BotStopReason describes why a bot was stopped.
type BotStopReason string

const (
	BotStopReasonClientCancel          BotStopReason = "client_cancel"
	BotStopReasonClientBulkCancel      BotStopReason = "client_bulk_cancel"
	BotStopReasonEndTime               BotStopReason = "end_time"
	BotStopReasonInstrumentDeactivated BotStopReason = "instrument_deactivated"
	BotStopReasonMarginBreach          BotStopReason = "margin_breach"
	BotStopReasonAdminCancel           BotStopReason = "admin_cancel"
	BotStopReasonConflict              BotStopReason = "conflict"
	BotStopReasonStrategy              BotStopReason = "strategy"
)

// String returns the string representation of the bot stop reason.
func (r BotStopReason) String() string {
	return string(r)
}

// IsValid returns true if the bot stop reason is a recognized value.
func (r BotStopReason) IsValid() bool {
	switch r {
	case BotStopReasonClientCancel, BotStopReasonClientBulkCancel, BotStopReasonEndTime, BotStopReasonInstrumentDeactivated, BotStopReasonMarginBreach, BotStopReasonAdminCancel, BotStopReasonConflict, BotStopReasonStrategy:
		return true
	}
	return false
}

// BotStopReasonValues returns all valid BotStopReason values.
func BotStopReasonValues() []BotStopReason {
	return []BotStopReason{BotStopReasonClientCancel, BotStopReasonClientBulkCancel, BotStopReasonEndTime, BotStopReasonInstrumentDeactivated, BotStopReasonMarginBreach, BotStopReasonAdminCancel, BotStopReasonConflict, BotStopReasonStrategy}
}
