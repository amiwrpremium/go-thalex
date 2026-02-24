package types

import "github.com/amiwrpremium/go-thalex/enums"

// RfqLeg represents a leg in an RFQ.
type RfqLeg struct {
	InstrumentName string  `json:"instrument_name"`
	Quantity       float64 `json:"quantity"`
	FeeQuantity    float64 `json:"fee_quantity"`
}

// RfqQuotedSide represents the quoted price/amount on one side of an RFQ.
type RfqQuotedSide struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

// Rfq represents a Request for Quote.
type Rfq struct {
	RfqID          string                `json:"rfq_id"`
	Legs           []RfqLeg              `json:"legs"`
	Amount         float64               `json:"amount"`
	CreateTime     float64               `json:"create_time"`
	ValidUntil     *float64              `json:"valid_until,omitempty"`
	Label          string                `json:"label,omitempty"`
	InsertReason   enums.RfqInsertReason `json:"insert_reason,omitempty"`
	DeleteReason   string                `json:"delete_reason,omitempty"`
	VolumeTickSize *float64              `json:"volume_tick_size,omitempty"`
	QuotedBid      *RfqQuotedSide        `json:"quoted_bid,omitempty"`
	QuotedAsk      *RfqQuotedSide        `json:"quoted_ask,omitempty"`
	TradePrice     *float64              `json:"trade_price,omitempty"`
	TradeAmount    *float64              `json:"trade_amount,omitempty"`
	CloseTime      *float64              `json:"close_time,omitempty"`
	Event          enums.RfqEvent        `json:"event,omitempty"`
}

// RfqOrder represents a quote on an RFQ.
type RfqOrder struct {
	RfqID         string                `json:"rfq_id"`
	OrderID       string                `json:"order_id"`
	ClientOrderID *uint64               `json:"client_order_id,omitempty"`
	Direction     enums.Direction       `json:"direction"`
	Price         float64               `json:"price"`
	Amount        float64               `json:"amount"`
	Label         string                `json:"label,omitempty"`
	TradePrice    *float64              `json:"trade_price,omitempty"`
	TradeAmount   *float64              `json:"trade_amount,omitempty"`
	DeleteReason  enums.RfqDeleteReason `json:"delete_reason,omitempty"`
	Event         enums.RfqOrderEvent   `json:"event,omitempty"`
}

// CreateRfqParams contains parameters for creating an RFQ.
type CreateRfqParams struct {
	Legs   []InsertLeg `json:"legs"`
	Amount float64     `json:"amount"`
	Label  string      `json:"label,omitempty"`
}

// TradeRfqParams contains parameters for trading on an RFQ.
type TradeRfqParams struct {
	RfqID     string          `json:"rfq_id"`
	Direction enums.Direction `json:"direction"`
	Price     float64         `json:"price"`
	Amount    float64         `json:"amount"`
}

// RfqQuoteInsertParams contains parameters for inserting an RFQ quote.
type RfqQuoteInsertParams struct {
	RfqID         string          `json:"rfq_id"`
	Direction     enums.Direction `json:"direction"`
	Amount        float64         `json:"amount"`
	Price         float64         `json:"price"`
	ClientOrderID *uint64         `json:"client_order_id,omitempty"`
	Label         string          `json:"label,omitempty"`
}

// RfqQuoteAmendParams contains parameters for amending an RFQ quote.
type RfqQuoteAmendParams struct {
	OrderID       string  `json:"order_id,omitempty"`
	ClientOrderID *uint64 `json:"client_order_id,omitempty"`
	Amount        float64 `json:"amount"`
	Price         float64 `json:"price"`
}

// RfqQuoteDeleteParams identifies an RFQ quote to delete.
type RfqQuoteDeleteParams struct {
	OrderID       string  `json:"order_id,omitempty"`
	ClientOrderID *uint64 `json:"client_order_id,omitempty"`
}
