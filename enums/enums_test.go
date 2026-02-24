package enums_test

import (
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
)

// ---------------------------------------------------------------------------
// Direction
// ---------------------------------------------------------------------------

func TestDirection_String(t *testing.T) {
	tests := []struct {
		d    enums.Direction
		want string
	}{
		{enums.DirectionBuy, "buy"},
		{enums.DirectionSell, "sell"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirection_IsValid(t *testing.T) {
	for _, d := range enums.DirectionValues() {
		t.Run(d.String(), func(t *testing.T) {
			if !d.IsValid() {
				t.Errorf("%q should be valid", d)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.Direction("").IsValid() {
			t.Error("empty Direction should be invalid")
		}
		if enums.Direction("unknown").IsValid() {
			t.Error("unknown Direction should be invalid")
		}
	})
}

func TestDirectionValues(t *testing.T) {
	vals := enums.DirectionValues()
	if len(vals) != 2 {
		t.Errorf("DirectionValues() count = %d, want 2", len(vals))
	}
}

func TestDirection_Opposite(t *testing.T) {
	tests := []struct {
		d    enums.Direction
		want enums.Direction
	}{
		{enums.DirectionBuy, enums.DirectionSell},
		{enums.DirectionSell, enums.DirectionBuy},
	}
	for _, tt := range tests {
		t.Run(tt.d.String(), func(t *testing.T) {
			if got := tt.d.Opposite(); got != tt.want {
				t.Errorf("Opposite() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirection_Opposite_Invalid(t *testing.T) {
	// An invalid direction should return DirectionBuy (the default fallthrough)
	d := enums.Direction("invalid")
	got := d.Opposite()
	if got != enums.DirectionBuy {
		t.Errorf("Opposite() of invalid = %q, want %q", got, enums.DirectionBuy)
	}
}

// ---------------------------------------------------------------------------
// OrderType
// ---------------------------------------------------------------------------

func TestOrderType_String(t *testing.T) {
	tests := []struct {
		v    enums.OrderType
		want string
	}{
		{enums.OrderTypeLimit, "limit"},
		{enums.OrderTypeMarket, "market"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOrderType_IsValid(t *testing.T) {
	for _, v := range enums.OrderTypeValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.OrderType("").IsValid() {
			t.Error("empty should be invalid")
		}
		if enums.OrderType("stop").IsValid() {
			t.Error("unknown should be invalid")
		}
	})
}

func TestOrderTypeValues(t *testing.T) {
	vals := enums.OrderTypeValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

// ---------------------------------------------------------------------------
// TimeInForce
// ---------------------------------------------------------------------------

func TestTimeInForce_String(t *testing.T) {
	tests := []struct {
		v    enums.TimeInForce
		want string
	}{
		{enums.TimeInForceGoodTillCancelled, "good_till_cancelled"},
		{enums.TimeInForceImmediateOrCancel, "immediate_or_cancel"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTimeInForce_IsValid(t *testing.T) {
	for _, v := range enums.TimeInForceValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.TimeInForce("").IsValid() {
			t.Error("empty should be invalid")
		}
		if enums.TimeInForce("fill_or_kill").IsValid() {
			t.Error("unknown should be invalid")
		}
	})
}

func TestTimeInForceValues(t *testing.T) {
	vals := enums.TimeInForceValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

// ---------------------------------------------------------------------------
// Collar
// ---------------------------------------------------------------------------

func TestCollar_String(t *testing.T) {
	tests := []struct {
		v    enums.Collar
		want string
	}{
		{enums.CollarIgnore, "ignore"},
		{enums.CollarReject, "reject"},
		{enums.CollarClamp, "clamp"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCollar_IsValid(t *testing.T) {
	for _, v := range enums.CollarValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.Collar("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestCollarValues(t *testing.T) {
	vals := enums.CollarValues()
	if len(vals) != 3 {
		t.Errorf("count = %d, want 3", len(vals))
	}
}

// ---------------------------------------------------------------------------
// Target
// ---------------------------------------------------------------------------

func TestTarget_String(t *testing.T) {
	tests := []struct {
		v    enums.Target
		want string
	}{
		{enums.TargetLast, "last"},
		{enums.TargetMark, "mark"},
		{enums.TargetIndex, "index"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTarget_IsValid(t *testing.T) {
	for _, v := range enums.TargetValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.Target("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestTargetValues(t *testing.T) {
	vals := enums.TargetValues()
	if len(vals) != 3 {
		t.Errorf("count = %d, want 3", len(vals))
	}
}

// ---------------------------------------------------------------------------
// OrderStatusValue
// ---------------------------------------------------------------------------

func TestOrderStatusValue_String(t *testing.T) {
	tests := []struct {
		v    enums.OrderStatusValue
		want string
	}{
		{enums.OrderStatusOpen, "open"},
		{enums.OrderStatusPartiallyFilled, "partially_filled"},
		{enums.OrderStatusCancelled, "cancelled"},
		{enums.OrderStatusCancelledPartiallyFilled, "cancelled_partially_filled"},
		{enums.OrderStatusFilled, "filled"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOrderStatusValue_IsValid(t *testing.T) {
	for _, v := range enums.OrderStatusValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.OrderStatusValue("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestOrderStatusValues(t *testing.T) {
	vals := enums.OrderStatusValues()
	if len(vals) != 5 {
		t.Errorf("count = %d, want 5", len(vals))
	}
}

func TestOrderStatusValue_IsActive(t *testing.T) {
	tests := []struct {
		v    enums.OrderStatusValue
		want bool
	}{
		{enums.OrderStatusOpen, true},
		{enums.OrderStatusPartiallyFilled, true},
		{enums.OrderStatusCancelled, false},
		{enums.OrderStatusCancelledPartiallyFilled, false},
		{enums.OrderStatusFilled, false},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStatusValue_IsFinal(t *testing.T) {
	tests := []struct {
		v    enums.OrderStatusValue
		want bool
	}{
		{enums.OrderStatusOpen, false},
		{enums.OrderStatusPartiallyFilled, false},
		{enums.OrderStatusCancelled, true},
		{enums.OrderStatusCancelledPartiallyFilled, true},
		{enums.OrderStatusFilled, true},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsFinal(); got != tt.want {
				t.Errorf("IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ChangeReason
// ---------------------------------------------------------------------------

func TestChangeReason_String(t *testing.T) {
	tests := []struct {
		v    enums.ChangeReason
		want string
	}{
		{enums.ChangeReasonExisting, "existing"},
		{enums.ChangeReasonInsert, "insert"},
		{enums.ChangeReasonAmend, "amend"},
		{enums.ChangeReasonCancel, "cancel"},
		{enums.ChangeReasonFill, "fill"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChangeReason_IsValid(t *testing.T) {
	for _, v := range enums.ChangeReasonValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.ChangeReason("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestChangeReasonValues(t *testing.T) {
	vals := enums.ChangeReasonValues()
	if len(vals) != 5 {
		t.Errorf("count = %d, want 5", len(vals))
	}
}

// ---------------------------------------------------------------------------
// DeleteReason
// ---------------------------------------------------------------------------

func TestDeleteReason_String(t *testing.T) {
	tests := []struct {
		v    enums.DeleteReason
		want string
	}{
		{enums.DeleteReasonClientCancel, "client_cancel"},
		{enums.DeleteReasonClientBulkCancel, "client_bulk_cancel"},
		{enums.DeleteReasonSessionEnd, "session_end"},
		{enums.DeleteReasonInstrumentDeactivated, "instrument_deactivated"},
		{enums.DeleteReasonMMProtection, "mm_protection"},
		{enums.DeleteReasonFailover, "failover"},
		{enums.DeleteReasonMarginBreach, "margin_breach"},
		{enums.DeleteReasonFilled, "filled"},
		{enums.DeleteReasonImmediateCancel, "immediate_cancel"},
		{enums.DeleteReasonAdminCancel, "admin_cancel"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDeleteReason_IsValid(t *testing.T) {
	for _, v := range enums.DeleteReasonValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.DeleteReason("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestDeleteReasonValues(t *testing.T) {
	vals := enums.DeleteReasonValues()
	if len(vals) != 10 {
		t.Errorf("count = %d, want 10", len(vals))
	}
}

// ---------------------------------------------------------------------------
// InsertReason
// ---------------------------------------------------------------------------

func TestInsertReason_String(t *testing.T) {
	tests := []struct {
		v    enums.InsertReason
		want string
	}{
		{enums.InsertReasonClientRequest, "client_request"},
		{enums.InsertReasonConditionalOrder, "conditional_order"},
		{enums.InsertReasonLiquidation, "liquidation"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInsertReason_IsValid(t *testing.T) {
	for _, v := range enums.InsertReasonValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.InsertReason("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestInsertReasonValues(t *testing.T) {
	vals := enums.InsertReasonValues()
	if len(vals) != 3 {
		t.Errorf("count = %d, want 3", len(vals))
	}
}

// ---------------------------------------------------------------------------
// InstrumentType
// ---------------------------------------------------------------------------

func TestInstrumentType_String(t *testing.T) {
	tests := []struct {
		v    enums.InstrumentType
		want string
	}{
		{enums.InstrumentTypePerpetual, "perpetual"},
		{enums.InstrumentTypeFuture, "future"},
		{enums.InstrumentTypeOption, "option"},
		{enums.InstrumentTypeCombination, "combination"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInstrumentType_IsValid(t *testing.T) {
	for _, v := range enums.InstrumentTypeValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.InstrumentType("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestInstrumentTypeValues(t *testing.T) {
	vals := enums.InstrumentTypeValues()
	if len(vals) != 4 {
		t.Errorf("count = %d, want 4", len(vals))
	}
}

// ---------------------------------------------------------------------------
// OptionType
// ---------------------------------------------------------------------------

func TestOptionType_String(t *testing.T) {
	tests := []struct {
		v    enums.OptionType
		want string
	}{
		{enums.OptionTypeCall, "call"},
		{enums.OptionTypePut, "put"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOptionType_IsValid(t *testing.T) {
	for _, v := range enums.OptionTypeValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.OptionType("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestOptionTypeValues(t *testing.T) {
	vals := enums.OptionTypeValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

// ---------------------------------------------------------------------------
// MakerTaker
// ---------------------------------------------------------------------------

func TestMakerTaker_String(t *testing.T) {
	tests := []struct {
		v    enums.MakerTaker
		want string
	}{
		{enums.MakerTakerMaker, "maker"},
		{enums.MakerTakerTaker, "taker"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMakerTaker_IsValid(t *testing.T) {
	for _, v := range enums.MakerTakerValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.MakerTaker("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestMakerTakerValues(t *testing.T) {
	vals := enums.MakerTakerValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

// ---------------------------------------------------------------------------
// TradeType
// ---------------------------------------------------------------------------

func TestTradeType_String(t *testing.T) {
	tests := []struct {
		v    enums.TradeType
		want string
	}{
		{enums.TradeTypeNormal, "normal"},
		{enums.TradeTypeBlock, "block"},
		{enums.TradeTypeCombo, "combo"},
		{enums.TradeTypeAmend, "amend"},
		{enums.TradeTypeDelete, "delete"},
		{enums.TradeTypeInternalTransfer, "internal_transfer"},
		{enums.TradeTypeExpiration, "expiration"},
		{enums.TradeTypeDailyMark, "daily_mark"},
		{enums.TradeTypeRfq, "rfq"},
		{enums.TradeTypeLiquidation, "liquidation"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTradeType_IsValid(t *testing.T) {
	for _, v := range enums.TradeTypeValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.TradeType("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestTradeTypeValues(t *testing.T) {
	vals := enums.TradeTypeValues()
	if len(vals) != 10 {
		t.Errorf("count = %d, want 10", len(vals))
	}
}

// ---------------------------------------------------------------------------
// BotStrategy
// ---------------------------------------------------------------------------

func TestBotStrategy_String(t *testing.T) {
	tests := []struct {
		v    enums.BotStrategy
		want string
	}{
		{enums.BotStrategySGSL, "sgsl"},
		{enums.BotStrategyOCQ, "ocq"},
		{enums.BotStrategyLevels, "levels"},
		{enums.BotStrategyGrid, "grid"},
		{enums.BotStrategyDHedge, "dhedge"},
		{enums.BotStrategyDFollow, "dfollow"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBotStrategy_IsValid(t *testing.T) {
	for _, v := range enums.BotStrategyValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.BotStrategy("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestBotStrategyValues(t *testing.T) {
	vals := enums.BotStrategyValues()
	if len(vals) != 6 {
		t.Errorf("count = %d, want 6", len(vals))
	}
}

// ---------------------------------------------------------------------------
// BotStatus
// ---------------------------------------------------------------------------

func TestBotStatus_String(t *testing.T) {
	tests := []struct {
		v    enums.BotStatus
		want string
	}{
		{enums.BotStatusActive, "active"},
		{enums.BotStatusStopped, "stopped"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBotStatus_IsValid(t *testing.T) {
	for _, v := range enums.BotStatusValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.BotStatus("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestBotStatusValues(t *testing.T) {
	vals := enums.BotStatusValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

func TestBotStatus_IsActive(t *testing.T) {
	tests := []struct {
		v    enums.BotStatus
		want bool
	}{
		{enums.BotStatusActive, true},
		{enums.BotStatusStopped, false},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBotStatus_IsFinal(t *testing.T) {
	tests := []struct {
		v    enums.BotStatus
		want bool
	}{
		{enums.BotStatusActive, false},
		{enums.BotStatusStopped, true},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsFinal(); got != tt.want {
				t.Errorf("IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// BotStopReason
// ---------------------------------------------------------------------------

func TestBotStopReason_String(t *testing.T) {
	tests := []struct {
		v    enums.BotStopReason
		want string
	}{
		{enums.BotStopReasonClientCancel, "client_cancel"},
		{enums.BotStopReasonClientBulkCancel, "client_bulk_cancel"},
		{enums.BotStopReasonEndTime, "end_time"},
		{enums.BotStopReasonInstrumentDeactivated, "instrument_deactivated"},
		{enums.BotStopReasonMarginBreach, "margin_breach"},
		{enums.BotStopReasonAdminCancel, "admin_cancel"},
		{enums.BotStopReasonConflict, "conflict"},
		{enums.BotStopReasonStrategy, "strategy"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBotStopReason_IsValid(t *testing.T) {
	for _, v := range enums.BotStopReasonValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.BotStopReason("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestBotStopReasonValues(t *testing.T) {
	vals := enums.BotStopReasonValues()
	if len(vals) != 8 {
		t.Errorf("count = %d, want 8", len(vals))
	}
}

// ---------------------------------------------------------------------------
// ConditionalOrderStatus
// ---------------------------------------------------------------------------

func TestConditionalOrderStatus_String(t *testing.T) {
	tests := []struct {
		v    enums.ConditionalOrderStatus
		want string
	}{
		{enums.ConditionalOrderStatusCreated, "created"},
		{enums.ConditionalOrderStatusActive, "active"},
		{enums.ConditionalOrderStatusConverted, "converted"},
		{enums.ConditionalOrderStatusRejected, "rejected"},
		{enums.ConditionalOrderStatusCancelRequested, "cancel requested"},
		{enums.ConditionalOrderStatusCancelled, "cancelled"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConditionalOrderStatus_IsValid(t *testing.T) {
	for _, v := range enums.ConditionalOrderStatusValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.ConditionalOrderStatus("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestConditionalOrderStatusValues(t *testing.T) {
	vals := enums.ConditionalOrderStatusValues()
	if len(vals) != 6 {
		t.Errorf("count = %d, want 6", len(vals))
	}
}

func TestConditionalOrderStatus_IsActive(t *testing.T) {
	tests := []struct {
		v    enums.ConditionalOrderStatus
		want bool
	}{
		{enums.ConditionalOrderStatusCreated, true},
		{enums.ConditionalOrderStatusActive, true},
		{enums.ConditionalOrderStatusConverted, false},
		{enums.ConditionalOrderStatusRejected, false},
		{enums.ConditionalOrderStatusCancelRequested, false},
		{enums.ConditionalOrderStatusCancelled, false},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// DepositStatus
// ---------------------------------------------------------------------------

func TestDepositStatus_String(t *testing.T) {
	tests := []struct {
		v    enums.DepositStatus
		want string
	}{
		{enums.DepositStatusUnconfirmed, "unconfirmed"},
		{enums.DepositStatusConfirmed, "confirmed"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDepositStatus_IsValid(t *testing.T) {
	for _, v := range enums.DepositStatusValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.DepositStatus("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestDepositStatusValues(t *testing.T) {
	vals := enums.DepositStatusValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

func TestDepositStatus_IsPending(t *testing.T) {
	tests := []struct {
		v    enums.DepositStatus
		want bool
	}{
		{enums.DepositStatusUnconfirmed, true},
		{enums.DepositStatusConfirmed, false},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsPending(); got != tt.want {
				t.Errorf("IsPending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepositStatus_IsFinal(t *testing.T) {
	tests := []struct {
		v    enums.DepositStatus
		want bool
	}{
		{enums.DepositStatusUnconfirmed, false},
		{enums.DepositStatusConfirmed, true},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsFinal(); got != tt.want {
				t.Errorf("IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// WithdrawalStatus
// ---------------------------------------------------------------------------

func TestWithdrawalStatus_String(t *testing.T) {
	tests := []struct {
		v    enums.WithdrawalStatus
		want string
	}{
		{enums.WithdrawalStatusPending, "pending"},
		{enums.WithdrawalStatusAwaitingConfirmation, "awaiting_confirmation"},
		{enums.WithdrawalStatusExecuting, "executing"},
		{enums.WithdrawalStatusExecuted, "executed"},
		{enums.WithdrawalStatusRejected, "rejected"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWithdrawalStatus_IsValid(t *testing.T) {
	for _, v := range enums.WithdrawalStatusValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.WithdrawalStatus("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestWithdrawalStatusValues(t *testing.T) {
	vals := enums.WithdrawalStatusValues()
	if len(vals) != 5 {
		t.Errorf("count = %d, want 5", len(vals))
	}
}

func TestWithdrawalStatus_IsPending(t *testing.T) {
	tests := []struct {
		v    enums.WithdrawalStatus
		want bool
	}{
		{enums.WithdrawalStatusPending, true},
		{enums.WithdrawalStatusAwaitingConfirmation, true},
		{enums.WithdrawalStatusExecuting, true},
		{enums.WithdrawalStatusExecuted, false},
		{enums.WithdrawalStatusRejected, false},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsPending(); got != tt.want {
				t.Errorf("IsPending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithdrawalStatus_IsFinal(t *testing.T) {
	tests := []struct {
		v    enums.WithdrawalStatus
		want bool
	}{
		{enums.WithdrawalStatusPending, false},
		{enums.WithdrawalStatusAwaitingConfirmation, false},
		{enums.WithdrawalStatusExecuting, false},
		{enums.WithdrawalStatusExecuted, true},
		{enums.WithdrawalStatusRejected, true},
	}
	for _, tt := range tests {
		t.Run(tt.v.String(), func(t *testing.T) {
			if got := tt.v.IsFinal(); got != tt.want {
				t.Errorf("IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Resolution
// ---------------------------------------------------------------------------

func TestResolution_String(t *testing.T) {
	tests := []struct {
		v    enums.Resolution
		want string
	}{
		{enums.Resolution1m, "1m"},
		{enums.Resolution5m, "5m"},
		{enums.Resolution15m, "15m"},
		{enums.Resolution30m, "30m"},
		{enums.Resolution1h, "1h"},
		{enums.Resolution1d, "1d"},
		{enums.Resolution1w, "1w"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolution_IsValid(t *testing.T) {
	for _, v := range enums.ResolutionValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.Resolution("").IsValid() {
			t.Error("empty should be invalid")
		}
		if enums.Resolution("2h").IsValid() {
			t.Error("unknown should be invalid")
		}
	})
}

func TestResolutionValues(t *testing.T) {
	vals := enums.ResolutionValues()
	if len(vals) != 7 {
		t.Errorf("count = %d, want 7", len(vals))
	}
}

// ---------------------------------------------------------------------------
// Sort
// ---------------------------------------------------------------------------

func TestSort_String(t *testing.T) {
	tests := []struct {
		v    enums.Sort
		want string
	}{
		{enums.SortAsc, "asc"},
		{enums.SortDesc, "desc"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSort_IsValid(t *testing.T) {
	for _, v := range enums.SortValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.Sort("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestSortValues(t *testing.T) {
	vals := enums.SortValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

// ---------------------------------------------------------------------------
// Product (no IsValid or Values -- only String)
// ---------------------------------------------------------------------------

func TestProduct_String(t *testing.T) {
	p := enums.Product("FBTCUSD")
	if got := p.String(); got != "FBTCUSD" {
		t.Errorf("String() = %q, want %q", got, "FBTCUSD")
	}

	empty := enums.Product("")
	if got := empty.String(); got != "" {
		t.Errorf("String() of empty = %q, want empty", got)
	}
}

// ---------------------------------------------------------------------------
// Severity
// ---------------------------------------------------------------------------

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		v    enums.Severity
		want string
	}{
		{enums.SeverityInfo, "info"},
		{enums.SeverityWarning, "warning"},
		{enums.SeverityCritical, "critical"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSeverity_IsValid(t *testing.T) {
	for _, v := range enums.SeverityValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.Severity("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestSeverityValues(t *testing.T) {
	vals := enums.SeverityValues()
	if len(vals) != 3 {
		t.Errorf("count = %d, want 3", len(vals))
	}
}

// ---------------------------------------------------------------------------
// DisplayType
// ---------------------------------------------------------------------------

func TestDisplayType_String(t *testing.T) {
	tests := []struct {
		v    enums.DisplayType
		want string
	}{
		{enums.DisplayTypeSuccess, "success"},
		{enums.DisplayTypeFailure, "failure"},
		{enums.DisplayTypeInfo, "info"},
		{enums.DisplayTypeWarning, "warning"},
		{enums.DisplayTypeCritical, "critical"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDisplayType_IsValid(t *testing.T) {
	for _, v := range enums.DisplayTypeValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.DisplayType("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestDisplayTypeValues(t *testing.T) {
	vals := enums.DisplayTypeValues()
	if len(vals) != 5 {
		t.Errorf("count = %d, want 5", len(vals))
	}
}

// ---------------------------------------------------------------------------
// Delay
// ---------------------------------------------------------------------------

func TestDelay_String(t *testing.T) {
	tests := []struct {
		v    enums.Delay
		want string
	}{
		{enums.DelayNone, "raw"},
		{enums.Delay100ms, "100ms"},
		{enums.Delay1000ms, "1000ms"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDelay_IsValid(t *testing.T) {
	for _, v := range enums.DelayValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.Delay("").IsValid() {
			t.Error("empty should be invalid")
		}
		if enums.Delay("500ms").IsValid() {
			t.Error("unknown should be invalid")
		}
	})
}

func TestDelayValues(t *testing.T) {
	vals := enums.DelayValues()
	if len(vals) != 3 {
		t.Errorf("count = %d, want 3", len(vals))
	}
}

// ---------------------------------------------------------------------------
// STPLevel
// ---------------------------------------------------------------------------

func TestSTPLevel_String(t *testing.T) {
	tests := []struct {
		v    enums.STPLevel
		want string
	}{
		{enums.STPLevelAccount, "account"},
		{enums.STPLevelCustomer, "customer"},
		{enums.STPLevelSubaccount, "subaccount"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSTPLevel_IsValid(t *testing.T) {
	for _, v := range enums.STPLevelValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.STPLevel("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestSTPLevelValues(t *testing.T) {
	vals := enums.STPLevelValues()
	if len(vals) != 3 {
		t.Errorf("count = %d, want 3", len(vals))
	}
}

// ---------------------------------------------------------------------------
// STPAction
// ---------------------------------------------------------------------------

func TestSTPAction_String(t *testing.T) {
	tests := []struct {
		v    enums.STPAction
		want string
	}{
		{enums.STPActionCancelAggressor, "cancel_aggressor"},
		{enums.STPActionCancelPassive, "cancel_passive"},
		{enums.STPActionCancelBoth, "cancel_both"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSTPAction_IsValid(t *testing.T) {
	for _, v := range enums.STPActionValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.STPAction("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestSTPActionValues(t *testing.T) {
	vals := enums.STPActionValues()
	if len(vals) != 3 {
		t.Errorf("count = %d, want 3", len(vals))
	}
}

// ---------------------------------------------------------------------------
// RecentTradesCategory
// ---------------------------------------------------------------------------

func TestRecentTradesCategory_String(t *testing.T) {
	tests := []struct {
		v    enums.RecentTradesCategory
		want string
	}{
		{enums.RecentTradesCategoryAll, "all"},
		{enums.RecentTradesCategoryNormal, "normal"},
		{enums.RecentTradesCategoryBlock, "block"},
		{enums.RecentTradesCategoryCombo, "combo"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRecentTradesCategory_IsValid(t *testing.T) {
	for _, v := range enums.RecentTradesCategoryValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.RecentTradesCategory("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestRecentTradesCategoryValues(t *testing.T) {
	vals := enums.RecentTradesCategoryValues()
	if len(vals) != 4 {
		t.Errorf("count = %d, want 4", len(vals))
	}
}

// ---------------------------------------------------------------------------
// MMProtectionReason
// ---------------------------------------------------------------------------

func TestMMProtectionReason_String(t *testing.T) {
	tests := []struct {
		v    enums.MMProtectionReason
		want string
	}{
		{enums.MMProtectionReasonTriggered, "triggered"},
		{enums.MMProtectionReasonReset, "reset"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMMProtectionReason_IsValid(t *testing.T) {
	for _, v := range enums.MMProtectionReasonValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.MMProtectionReason("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestMMProtectionReasonValues(t *testing.T) {
	vals := enums.MMProtectionReasonValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

// ---------------------------------------------------------------------------
// RfqEvent
// ---------------------------------------------------------------------------

func TestRfqEvent_String(t *testing.T) {
	tests := []struct {
		v    enums.RfqEvent
		want string
	}{
		{enums.RfqEventCreated, "Created"},
		{enums.RfqEventCancelled, "Cancelled"},
		{enums.RfqEventTraded, "Traded"},
		{enums.RfqEventExisting, "Existing"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRfqEvent_IsValid(t *testing.T) {
	for _, v := range enums.RfqEventValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.RfqEvent("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestRfqEventValues(t *testing.T) {
	vals := enums.RfqEventValues()
	if len(vals) != 4 {
		t.Errorf("count = %d, want 4", len(vals))
	}
}

// ---------------------------------------------------------------------------
// RfqOrderEvent
// ---------------------------------------------------------------------------

func TestRfqOrderEvent_String(t *testing.T) {
	tests := []struct {
		v    enums.RfqOrderEvent
		want string
	}{
		{enums.RfqOrderEventInserted, "Inserted"},
		{enums.RfqOrderEventAmended, "Amended"},
		{enums.RfqOrderEventCancelled, "Cancelled"},
		{enums.RfqOrderEventFilled, "Filled"},
		{enums.RfqOrderEventExisting, "Existing"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRfqOrderEvent_IsValid(t *testing.T) {
	for _, v := range enums.RfqOrderEventValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.RfqOrderEvent("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestRfqOrderEventValues(t *testing.T) {
	vals := enums.RfqOrderEventValues()
	if len(vals) != 5 {
		t.Errorf("count = %d, want 5", len(vals))
	}
}

// ---------------------------------------------------------------------------
// RfqDeleteReason
// ---------------------------------------------------------------------------

func TestRfqDeleteReason_String(t *testing.T) {
	tests := []struct {
		v    enums.RfqDeleteReason
		want string
	}{
		{enums.RfqDeleteReasonClientCancel, "client_cancel"},
		{enums.RfqDeleteReasonSessionEnd, "session_end"},
		{enums.RfqDeleteReasonInstrumentDeactivated, "instrument_deactivated"},
		{enums.RfqDeleteReasonMMProtection, "mm_protection"},
		{enums.RfqDeleteReasonFailover, "failover"},
		{enums.RfqDeleteReasonMarginBreach, "margin_breach"},
		{enums.RfqDeleteReasonFilled, "filled"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRfqDeleteReason_IsValid(t *testing.T) {
	for _, v := range enums.RfqDeleteReasonValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.RfqDeleteReason("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestRfqDeleteReasonValues(t *testing.T) {
	vals := enums.RfqDeleteReasonValues()
	if len(vals) != 7 {
		t.Errorf("count = %d, want 7", len(vals))
	}
}

// ---------------------------------------------------------------------------
// RfqInsertReason
// ---------------------------------------------------------------------------

func TestRfqInsertReason_String(t *testing.T) {
	tests := []struct {
		v    enums.RfqInsertReason
		want string
	}{
		{enums.RfqInsertReasonClientRequest, "client_request"},
		{enums.RfqInsertReasonLiquidation, "liquidation"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRfqInsertReason_IsValid(t *testing.T) {
	for _, v := range enums.RfqInsertReasonValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.RfqInsertReason("").IsValid() {
			t.Error("empty should be invalid")
		}
	})
}

func TestRfqInsertReasonValues(t *testing.T) {
	vals := enums.RfqInsertReasonValues()
	if len(vals) != 2 {
		t.Errorf("count = %d, want 2", len(vals))
	}
}

// ---------------------------------------------------------------------------
// SystemEventType
// ---------------------------------------------------------------------------

func TestSystemEventType_String(t *testing.T) {
	tests := []struct {
		v    enums.SystemEventType
		want string
	}{
		{enums.SystemEventTypeReconnect, "reconnect"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSystemEventType_IsValid(t *testing.T) {
	for _, v := range enums.SystemEventTypeValues() {
		t.Run(v.String(), func(t *testing.T) {
			if !v.IsValid() {
				t.Errorf("%q should be valid", v)
			}
		})
	}
	t.Run("Invalid", func(t *testing.T) {
		if enums.SystemEventType("").IsValid() {
			t.Error("empty should be invalid")
		}
		if enums.SystemEventType("shutdown").IsValid() {
			t.Error("unknown should be invalid")
		}
	})
}

func TestSystemEventTypeValues(t *testing.T) {
	vals := enums.SystemEventTypeValues()
	if len(vals) != 1 {
		t.Errorf("count = %d, want 1", len(vals))
	}
}

// ---------------------------------------------------------------------------
// Cross-cutting: all Values() slices should match IsValid()
// ---------------------------------------------------------------------------

func TestAllValuesAreValid(t *testing.T) {
	t.Run("Direction", func(t *testing.T) {
		for _, v := range enums.DirectionValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("OrderType", func(t *testing.T) {
		for _, v := range enums.OrderTypeValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("TimeInForce", func(t *testing.T) {
		for _, v := range enums.TimeInForceValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("Collar", func(t *testing.T) {
		for _, v := range enums.CollarValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("Target", func(t *testing.T) {
		for _, v := range enums.TargetValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("OrderStatusValue", func(t *testing.T) {
		for _, v := range enums.OrderStatusValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("ChangeReason", func(t *testing.T) {
		for _, v := range enums.ChangeReasonValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("DeleteReason", func(t *testing.T) {
		for _, v := range enums.DeleteReasonValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("InsertReason", func(t *testing.T) {
		for _, v := range enums.InsertReasonValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("InstrumentType", func(t *testing.T) {
		for _, v := range enums.InstrumentTypeValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("OptionType", func(t *testing.T) {
		for _, v := range enums.OptionTypeValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("MakerTaker", func(t *testing.T) {
		for _, v := range enums.MakerTakerValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("TradeType", func(t *testing.T) {
		for _, v := range enums.TradeTypeValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("BotStrategy", func(t *testing.T) {
		for _, v := range enums.BotStrategyValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("BotStatus", func(t *testing.T) {
		for _, v := range enums.BotStatusValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("BotStopReason", func(t *testing.T) {
		for _, v := range enums.BotStopReasonValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("ConditionalOrderStatus", func(t *testing.T) {
		for _, v := range enums.ConditionalOrderStatusValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("DepositStatus", func(t *testing.T) {
		for _, v := range enums.DepositStatusValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("WithdrawalStatus", func(t *testing.T) {
		for _, v := range enums.WithdrawalStatusValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("Resolution", func(t *testing.T) {
		for _, v := range enums.ResolutionValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("Sort", func(t *testing.T) {
		for _, v := range enums.SortValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("Severity", func(t *testing.T) {
		for _, v := range enums.SeverityValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("DisplayType", func(t *testing.T) {
		for _, v := range enums.DisplayTypeValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("Delay", func(t *testing.T) {
		for _, v := range enums.DelayValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("STPLevel", func(t *testing.T) {
		for _, v := range enums.STPLevelValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("STPAction", func(t *testing.T) {
		for _, v := range enums.STPActionValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("RecentTradesCategory", func(t *testing.T) {
		for _, v := range enums.RecentTradesCategoryValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("MMProtectionReason", func(t *testing.T) {
		for _, v := range enums.MMProtectionReasonValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("RfqEvent", func(t *testing.T) {
		for _, v := range enums.RfqEventValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("RfqOrderEvent", func(t *testing.T) {
		for _, v := range enums.RfqOrderEventValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("RfqDeleteReason", func(t *testing.T) {
		for _, v := range enums.RfqDeleteReasonValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("RfqInsertReason", func(t *testing.T) {
		for _, v := range enums.RfqInsertReasonValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
	t.Run("SystemEventType", func(t *testing.T) {
		for _, v := range enums.SystemEventTypeValues() {
			if !v.IsValid() {
				t.Errorf("%q from Values() should be valid", v)
			}
		}
	})
}
