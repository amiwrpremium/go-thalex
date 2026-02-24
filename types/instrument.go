package types

import "github.com/amiwrpremium/go-thalex/enums"

// Instrument represents a trading instrument on the exchange.
type Instrument struct {
	InstrumentName       string               `json:"instrument_name"`
	Product              string               `json:"product"`
	TickSize             float64              `json:"tick_size"`
	VolumeTickSize       float64              `json:"volume_tick_size"`
	MinOrderAmount       float64              `json:"min_order_amount"`
	Underlying           string               `json:"underlying"`
	Type                 enums.InstrumentType `json:"type"`
	OptionType           enums.OptionType     `json:"option_type,omitempty"`
	ExpiryDate           string               `json:"expiry_date,omitempty"`
	ExpirationTimestamp  *int64               `json:"expiration_timestamp,omitempty"`
	StrikePrice          *float64             `json:"strike_price,omitempty"`
	BaseCurrency         string               `json:"base_currency,omitempty"`
	Legs                 []Leg                `json:"legs,omitempty"`
	CreateTime           float64              `json:"create_time,omitempty"`
	SettlementPrice      *float64             `json:"settlement_price,omitempty"`
	SettlementIndexPrice *float64             `json:"settlement_index_price,omitempty"`
}

// IsOption returns true if the instrument is an option.
func (i *Instrument) IsOption() bool { return i.Type == enums.InstrumentTypeOption }

// IsFuture returns true if the instrument is a future.
func (i *Instrument) IsFuture() bool { return i.Type == enums.InstrumentTypeFuture }

// IsPerpetual returns true if the instrument is a perpetual.
func (i *Instrument) IsPerpetual() bool { return i.Type == enums.InstrumentTypePerpetual }

// IsCombination returns true if the instrument is a combination.
func (i *Instrument) IsCombination() bool { return i.Type == enums.InstrumentTypeCombination }

// Ticker represents a full ticker for an instrument.
type Ticker struct {
	BestBidPrice          *float64 `json:"best_bid_price"`
	BestBidAmount         *float64 `json:"best_bid_amount"`
	BestAskPrice          *float64 `json:"best_ask_price"`
	BestAskAmount         *float64 `json:"best_ask_amount"`
	LastPrice             *float64 `json:"last_price"`
	MarkPrice             float64  `json:"mark_price"`
	MarkTimestamp         float64  `json:"mark_timestamp"`
	IV                    *float64 `json:"iv,omitempty"`
	Delta                 *float64 `json:"delta,omitempty"`
	Index                 *float64 `json:"index,omitempty"`
	Forward               *float64 `json:"forward,omitempty"`
	Volume24h             *float64 `json:"volume_24h,omitempty"`
	Value24h              *float64 `json:"value_24h,omitempty"`
	LowPrice24h           *float64 `json:"low_price_24h,omitempty"`
	HighPrice24h          *float64 `json:"high_price_24h,omitempty"`
	Change24h             *float64 `json:"change_24h,omitempty"`
	CollarLow             *float64 `json:"collar_low,omitempty"`
	CollarHigh            *float64 `json:"collar_high,omitempty"`
	OpenInterest          *float64 `json:"open_interest,omitempty"`
	FundingRate           *float64 `json:"funding_rate,omitempty"`
	FundingMark           *float64 `json:"funding_mark,omitempty"`
	RealisedFunding24h    *float64 `json:"realised_funding_24h,omitempty"`
	AverageFundingRate24h *float64 `json:"average_funding_rate_24h,omitempty"`
}

// Spread returns the bid-ask spread, or nil if either side is empty.
func (t *Ticker) Spread() *float64 {
	if t.BestBidPrice == nil || t.BestAskPrice == nil {
		return nil
	}
	s := *t.BestAskPrice - *t.BestBidPrice
	return &s
}

// MidPrice returns the mid price, or nil if either side is empty.
func (t *Ticker) MidPrice() *float64 {
	if t.BestBidPrice == nil || t.BestAskPrice == nil {
		return nil
	}
	m := (*t.BestBidPrice + *t.BestAskPrice) / 2
	return &m
}

// IndexPrice represents index price information.
type IndexPrice struct {
	IndexName               string   `json:"index_name"`
	Price                   float64  `json:"price"`
	Timestamp               float64  `json:"timestamp"`
	ExpirationPrintAverage  *float64 `json:"expiration_print_average,omitempty"`
	ExpirationProgress      *float64 `json:"expiration_progress,omitempty"`
	ExpectedExpirationPrice *float64 `json:"expected_expiration_price,omitempty"`
	PreviousSettlementPrice *float64 `json:"previous_settlement_price,omitempty"`
}

// BookLevel represents a single price level [price, amount, outright_amount].
type BookLevel [3]float64

// Price returns the price of this book level.
func (l BookLevel) Price() float64 { return l[0] }

// Amount returns the total amount at this level.
func (l BookLevel) Amount() float64 { return l[1] }

// OutrightAmount returns the outright amount at this level.
func (l BookLevel) OutrightAmount() float64 { return l[2] }

// Book represents an order book snapshot.
type Book struct {
	Bids []BookLevel `json:"bids"`
	Asks []BookLevel `json:"asks"`
	Last *float64    `json:"last,omitempty"`
	Time float64     `json:"time"`
}
