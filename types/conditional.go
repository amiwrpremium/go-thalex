package types

import "github.com/amiwrpremium/go-thalex/enums"

// ConditionalOrder represents a conditional (stop/bracket/trailing) order.
type ConditionalOrder struct {
	OrderID                  string                       `json:"order_id"`
	InstrumentName           string                       `json:"instrument_name"`
	Direction                enums.Direction              `json:"direction"`
	Amount                   float64                      `json:"amount"`
	Target                   enums.Target                 `json:"target"`
	StopPrice                float64                      `json:"stop_price"`
	LimitPrice               *float64                     `json:"limit_price,omitempty"`
	BracketPrice             *float64                     `json:"bracket_price,omitempty"`
	TrailingStopCallbackRate *float64                     `json:"trailing_stop_callback_rate,omitempty"`
	Label                    string                       `json:"label,omitempty"`
	Status                   enums.ConditionalOrderStatus `json:"status"`
	CreateTime               float64                      `json:"create_time"`
	UpdateTime               float64                      `json:"update_time"`
	ConvertTime              *float64                     `json:"convert_time,omitempty"`
	ConvertedOrderID         string                       `json:"converted_order_id,omitempty"`
	RejectReason             string                       `json:"reject_reason"`
	ReduceOnly               bool                         `json:"reduce_only"`
}

// IsStopLimit returns true if this is a stop limit order.
func (o *ConditionalOrder) IsStopLimit() bool { return o.LimitPrice != nil }

// IsBracket returns true if this is a bracket order.
func (o *ConditionalOrder) IsBracket() bool { return o.BracketPrice != nil }

// IsTrailingStop returns true if this is a trailing stop order.
func (o *ConditionalOrder) IsTrailingStop() bool { return o.TrailingStopCallbackRate != nil }

// CreateConditionalOrderParams contains parameters for creating a conditional order.
type CreateConditionalOrderParams struct {
	Direction                enums.Direction `json:"direction"`
	InstrumentName           string          `json:"instrument_name"`
	Amount                   float64         `json:"amount"`
	StopPrice                float64         `json:"stop_price"`
	LimitPrice               *float64        `json:"limit_price,omitempty"`
	BracketPrice             *float64        `json:"bracket_price,omitempty"`
	TrailingStopCallbackRate *float64        `json:"trailing_stop_callback_rate,omitempty"`
	Target                   enums.Target    `json:"target,omitempty"`
	Label                    string          `json:"label,omitempty"`
	ReduceOnly               *bool           `json:"reduce_only,omitempty"`
}

// NewStopOrder creates a simple stop order.
func NewStopOrder(direction enums.Direction, instrumentName string, amount, stopPrice float64) *CreateConditionalOrderParams {
	return &CreateConditionalOrderParams{
		Direction: direction, InstrumentName: instrumentName,
		Amount: amount, StopPrice: stopPrice,
	}
}

// NewStopLimitOrder creates a stop limit order.
func NewStopLimitOrder(direction enums.Direction, instrumentName string, amount, stopPrice, limitPrice float64) *CreateConditionalOrderParams {
	return &CreateConditionalOrderParams{
		Direction: direction, InstrumentName: instrumentName,
		Amount: amount, StopPrice: stopPrice, LimitPrice: &limitPrice,
	}
}

// NewBracketOrder creates a bracket order.
func NewBracketOrder(direction enums.Direction, instrumentName string, amount, stopPrice, bracketPrice float64) *CreateConditionalOrderParams {
	return &CreateConditionalOrderParams{
		Direction: direction, InstrumentName: instrumentName,
		Amount: amount, StopPrice: stopPrice, BracketPrice: &bracketPrice,
	}
}

// NewTrailingStopOrder creates a trailing stop loss order.
func NewTrailingStopOrder(direction enums.Direction, instrumentName string, amount, stopPrice, callbackRate float64) *CreateConditionalOrderParams {
	return &CreateConditionalOrderParams{
		Direction: direction, InstrumentName: instrumentName,
		Amount: amount, StopPrice: stopPrice, TrailingStopCallbackRate: &callbackRate,
	}
}

// WithTarget sets the trigger target (last, mark, or index).
func (p *CreateConditionalOrderParams) WithTarget(v enums.Target) *CreateConditionalOrderParams {
	p.Target = v
	return p
}

// WithLabel sets a user label.
func (p *CreateConditionalOrderParams) WithLabel(v string) *CreateConditionalOrderParams {
	p.Label = v
	return p
}

// WithReduceOnly marks the activated order as reduce-only.
func (p *CreateConditionalOrderParams) WithReduceOnly(v bool) *CreateConditionalOrderParams {
	p.ReduceOnly = &v
	return p
}
