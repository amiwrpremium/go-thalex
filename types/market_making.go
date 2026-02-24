package types

import "github.com/amiwrpremium/go-thalex/enums"

// QuoteLevel represents a single price level in a mass quote.
type QuoteLevel struct {
	Price  float64
	Amount float64
}

// SingleLevelQuote represents a single-level bid or ask quote.
type SingleLevelQuote struct {
	P float64 `json:"p"`
	A float64 `json:"a"`
}

// DoubleSidedQuote represents a mass quote for a single instrument.
type DoubleSidedQuote struct {
	I string `json:"i"`
	B any    `json:"b,omitempty"`
	A any    `json:"a,omitempty"`
}

// NewDoubleSidedQuote creates a mass quote for an instrument with multi-level quotes.
func NewDoubleSidedQuote(instrument string, bids, asks []QuoteLevel) DoubleSidedQuote {
	q := DoubleSidedQuote{I: instrument}
	if len(bids) > 0 {
		bidArr := make([][2]float64, len(bids))
		for i, b := range bids {
			bidArr[i] = [2]float64{b.Price, b.Amount}
		}
		q.B = bidArr
	}
	if len(asks) > 0 {
		askArr := make([][2]float64, len(asks))
		for i, a := range asks {
			askArr[i] = [2]float64{a.Price, a.Amount}
		}
		q.A = askArr
	}
	return q
}

// NewSingleLevelQuote creates a mass quote with single-level bid and ask.
func NewSingleLevelQuote(instrument string, bidPrice, bidAmount, askPrice, askAmount float64) DoubleSidedQuote {
	return DoubleSidedQuote{
		I: instrument,
		B: SingleLevelQuote{P: bidPrice, A: bidAmount},
		A: SingleLevelQuote{P: askPrice, A: askAmount},
	}
}

// DoubleSidedQuoteResult is the response from a mass quote operation.
type DoubleSidedQuoteResult struct {
	NSuccess int          `json:"n_success"`
	NFail    int          `json:"n_fail"`
	Errors   []QuoteError `json:"errors"`
}

// QuoteError describes a failure in a mass quote operation.
type QuoteError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Side    string   `json:"side,omitempty"`
	Price   *float64 `json:"price,omitempty"`
}

// MassQuoteParams contains parameters for a mass quote request.
type MassQuoteParams struct {
	Quotes         []DoubleSidedQuote `json:"quotes"`
	Label          string             `json:"label,omitempty"`
	PostOnly       *bool              `json:"post_only,omitempty"`
	RejectPostOnly *bool              `json:"reject_post_only,omitempty"`
	STPLevel       enums.STPLevel     `json:"stp_level,omitempty"`
	STPAction      enums.STPAction    `json:"stp_action,omitempty"`
}

// NewMassQuoteParams creates MassQuoteParams with the given quotes.
func NewMassQuoteParams(quotes []DoubleSidedQuote) *MassQuoteParams {
	return &MassQuoteParams{Quotes: quotes}
}

// WithLabel sets a label for all quotes.
func (p *MassQuoteParams) WithLabel(v string) *MassQuoteParams { p.Label = v; return p }

// WithPostOnly enables post-only mode for all quotes.
func (p *MassQuoteParams) WithPostOnly(v bool) *MassQuoteParams { p.PostOnly = &v; return p }

// WithRejectPostOnly enables reject-post-only mode.
func (p *MassQuoteParams) WithRejectPostOnly(v bool) *MassQuoteParams {
	p.RejectPostOnly = &v
	return p
}

// WithSTP sets self-trade prevention parameters.
func (p *MassQuoteParams) WithSTP(level enums.STPLevel, action enums.STPAction) *MassQuoteParams {
	p.STPLevel = level
	p.STPAction = action
	return p
}

// MMProtectionParams configures market maker protection for a product.
type MMProtectionParams struct {
	Product     enums.Product `json:"product"`
	TradeAmount float64       `json:"trade_amount"`
	QuoteAmount float64       `json:"quote_amount"`
}
