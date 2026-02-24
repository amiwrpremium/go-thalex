package types

import "github.com/amiwrpremium/go-thalex/enums"

// OrderStatus represents the full status of an order.
type OrderStatus struct {
	OrderID            string                 `json:"order_id"`
	OrderType          enums.OrderType        `json:"order_type"`
	TimeInForce        enums.TimeInForce      `json:"time_in_force"`
	InstrumentName     string                 `json:"instrument_name,omitempty"`
	Legs               []Leg                  `json:"legs,omitempty"`
	Direction          enums.Direction        `json:"direction"`
	Price              *float64               `json:"price,omitempty"`
	Amount             float64                `json:"amount"`
	FilledAmount       float64                `json:"filled_amount"`
	RemainingAmount    float64                `json:"remaining_amount"`
	Label              string                 `json:"label,omitempty"`
	ClientOrderID      *uint64                `json:"client_order_id,omitempty"`
	Status             enums.OrderStatusValue `json:"status"`
	Fills              []OrderFill            `json:"fills"`
	ChangeReason       enums.ChangeReason     `json:"change_reason"`
	DeleteReason       enums.DeleteReason     `json:"delete_reason,omitempty"`
	InsertReason       enums.InsertReason     `json:"insert_reason"`
	ConditionalOrderID string                 `json:"conditional_order_id,omitempty"`
	BotID              string                 `json:"bot_id,omitempty"`
	CreateTime         float64                `json:"create_time"`
	CloseTime          *float64               `json:"close_time,omitempty"`
	ReduceOnly         bool                   `json:"reduce_only,omitempty"`
	Persistent         bool                   `json:"persistent"`
}

// OrderFill represents a single fill on an order.
type OrderFill struct {
	TradeID    string           `json:"trade_id"`
	Price      float64          `json:"price"`
	Amount     float64          `json:"amount"`
	Time       float64          `json:"time,omitempty"`
	MakerTaker enums.MakerTaker `json:"maker_taker"`
	LegIndex   int              `json:"leg_index"`
}

// OrderHistory represents a historical order.
type OrderHistory struct {
	OrderID            string                 `json:"order_id"`
	OrderType          enums.OrderType        `json:"order_type"`
	InstrumentName     string                 `json:"instrument_name,omitempty"`
	Legs               []Leg                  `json:"legs,omitempty"`
	Direction          enums.Direction        `json:"direction"`
	Price              *float64               `json:"price,omitempty"`
	Amount             float64                `json:"amount"`
	FilledAmount       float64                `json:"filled_amount"`
	Label              string                 `json:"label,omitempty"`
	ClientOrderID      *uint64                `json:"client_order_id,omitempty"`
	Status             enums.OrderStatusValue `json:"status"`
	Fills              []OrderFill            `json:"fills"`
	DeleteReason       enums.DeleteReason     `json:"delete_reason,omitempty"`
	InsertReason       enums.InsertReason     `json:"insert_reason"`
	ConditionalOrderID string                 `json:"conditional_order_id,omitempty"`
	BotID              string                 `json:"bot_id,omitempty"`
	CreateTime         float64                `json:"create_time"`
	CloseTime          float64                `json:"close_time"`
	ReduceOnly         bool                   `json:"reduce_only,omitempty"`
}

// InsertOrderParams contains parameters for inserting an order.
// Use NewInsertOrderParams or NewComboInsertOrderParams to create.
type InsertOrderParams struct {
	Direction      enums.Direction   `json:"direction"`
	InstrumentName string            `json:"instrument_name,omitempty"`
	Legs           []InsertLeg       `json:"legs,omitempty"`
	Amount         float64           `json:"amount"`
	Price          *float64          `json:"price,omitempty"`
	OrderType      enums.OrderType   `json:"order_type,omitempty"`
	TimeInForce    enums.TimeInForce `json:"time_in_force,omitempty"`
	PostOnly       *bool             `json:"post_only,omitempty"`
	RejectPostOnly *bool             `json:"reject_post_only,omitempty"`
	ReduceOnly     *bool             `json:"reduce_only,omitempty"`
	Collar         enums.Collar      `json:"collar,omitempty"`
	Label          string            `json:"label,omitempty"`
	ClientOrderID  *uint64           `json:"client_order_id,omitempty"`
	STPLevel       enums.STPLevel    `json:"stp_level,omitempty"`
	STPAction      enums.STPAction   `json:"stp_action,omitempty"`
}

// NewInsertOrderParams creates parameters for a single-instrument order.
func NewInsertOrderParams(direction enums.Direction, instrumentName string, amount float64) *InsertOrderParams {
	return &InsertOrderParams{
		Direction:      direction,
		InstrumentName: instrumentName,
		Amount:         amount,
	}
}

// NewBuyOrderParams creates parameters for a buy order.
func NewBuyOrderParams(instrumentName string, amount float64) *InsertOrderParams {
	return NewInsertOrderParams(enums.DirectionBuy, instrumentName, amount)
}

// NewSellOrderParams creates parameters for a sell order.
func NewSellOrderParams(instrumentName string, amount float64) *InsertOrderParams {
	return NewInsertOrderParams(enums.DirectionSell, instrumentName, amount)
}

// NewComboInsertOrderParams creates parameters for a combination order.
func NewComboInsertOrderParams(direction enums.Direction, legs []InsertLeg, amount float64) *InsertOrderParams {
	return &InsertOrderParams{
		Direction: direction,
		Legs:      legs,
		Amount:    amount,
	}
}

// WithPrice sets the limit price.
func (p *InsertOrderParams) WithPrice(v float64) *InsertOrderParams { p.Price = &v; return p }

// WithOrderType sets the order type (limit or market).
func (p *InsertOrderParams) WithOrderType(v enums.OrderType) *InsertOrderParams {
	p.OrderType = v
	return p
}

// WithTimeInForce sets the time in force policy.
func (p *InsertOrderParams) WithTimeInForce(v enums.TimeInForce) *InsertOrderParams {
	p.TimeInForce = v
	return p
}

// WithPostOnly enables post-only mode.
func (p *InsertOrderParams) WithPostOnly(v bool) *InsertOrderParams { p.PostOnly = &v; return p }

// WithRejectPostOnly enables reject-post-only mode (book-or-cancel with PostOnly).
func (p *InsertOrderParams) WithRejectPostOnly(v bool) *InsertOrderParams {
	p.RejectPostOnly = &v
	return p
}

// WithReduceOnly marks the order as reduce-only.
func (p *InsertOrderParams) WithReduceOnly(v bool) *InsertOrderParams { p.ReduceOnly = &v; return p }

// WithCollar sets the collar handling mode.
func (p *InsertOrderParams) WithCollar(v enums.Collar) *InsertOrderParams { p.Collar = v; return p }

// WithLabel sets a user label for the order.
func (p *InsertOrderParams) WithLabel(v string) *InsertOrderParams { p.Label = v; return p }

// WithClientOrderID sets a client-supplied order ID.
func (p *InsertOrderParams) WithClientOrderID(v uint64) *InsertOrderParams {
	p.ClientOrderID = &v
	return p
}

// WithSTP sets self-trade prevention parameters.
func (p *InsertOrderParams) WithSTP(level enums.STPLevel, action enums.STPAction) *InsertOrderParams {
	p.STPLevel = level
	p.STPAction = action
	return p
}

// AmendOrderParams contains parameters for amending an order.
type AmendOrderParams struct {
	OrderID       string       `json:"order_id,omitempty"`
	ClientOrderID *uint64      `json:"client_order_id,omitempty"`
	Price         float64      `json:"price"`
	Amount        float64      `json:"amount"`
	Collar        enums.Collar `json:"collar,omitempty"`
}

// NewAmendByOrderID creates amend parameters using a system order ID.
func NewAmendByOrderID(orderID string, price, amount float64) *AmendOrderParams {
	return &AmendOrderParams{OrderID: orderID, Price: price, Amount: amount}
}

// NewAmendByClientOrderID creates amend parameters using a client order ID.
func NewAmendByClientOrderID(clientOrderID uint64, price, amount float64) *AmendOrderParams {
	return &AmendOrderParams{ClientOrderID: &clientOrderID, Price: price, Amount: amount}
}

// WithCollar sets the collar handling mode for the amend.
func (p *AmendOrderParams) WithCollar(v enums.Collar) *AmendOrderParams { p.Collar = v; return p }

// CancelOrderParams identifies an order to cancel.
type CancelOrderParams struct {
	OrderID       string  `json:"order_id,omitempty"`
	ClientOrderID *uint64 `json:"client_order_id,omitempty"`
}

// CancelByOrderID creates cancel parameters using a system order ID.
func CancelByOrderID(orderID string) *CancelOrderParams {
	return &CancelOrderParams{OrderID: orderID}
}

// CancelByClientOrderID creates cancel parameters using a client order ID.
func CancelByClientOrderID(clientOrderID uint64) *CancelOrderParams {
	return &CancelOrderParams{ClientOrderID: &clientOrderID}
}
