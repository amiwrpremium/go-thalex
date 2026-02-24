package types

import "github.com/amiwrpremium/go-thalex/enums"

// Trade represents a single trade.
type Trade struct {
	TradeType            enums.TradeType  `json:"trade_type"`
	TradeID              string           `json:"trade_id"`
	OrderID              string           `json:"order_id"`
	InstrumentName       string           `json:"instrument_name"`
	Direction            enums.Direction  `json:"direction"`
	Price                float64          `json:"price"`
	Amount               float64          `json:"amount"`
	Label                string           `json:"label,omitempty"`
	Time                 float64          `json:"time"`
	PositionAfter        float64          `json:"position_after"`
	SessionRealisedAfter float64          `json:"session_realised_after,omitempty"`
	PositionPnl          *float64         `json:"position_pnl,omitempty"`
	PerpetualFundingPnl  *float64         `json:"perpetual_funding_pnl,omitempty"`
	Fee                  float64          `json:"fee"`
	Index                *float64         `json:"index,omitempty"`
	FeeRate              float64          `json:"fee_rate"`
	FeeBasis             float64          `json:"fee_basis"`
	FundingMark          *float64         `json:"funding_mark,omitempty"`
	LiquidationFee       *float64         `json:"liquidation_fee,omitempty"`
	ClientOrderID        *uint64          `json:"client_order_id,omitempty"`
	MakerTaker           enums.MakerTaker `json:"maker_taker,omitempty"`
	BotID                string           `json:"bot_id,omitempty"`
	LegIndex             int              `json:"leg_index"`
}

// DailyMark represents a daily mark settlement entry.
type DailyMark struct {
	Time                float64  `json:"time"`
	InstrumentName      string   `json:"instrument_name"`
	Position            float64  `json:"position"`
	MarkPrice           float64  `json:"mark_price"`
	RealizedPositionPnl float64  `json:"realized_position_pnl"`
	RealizedFundingPnl  *float64 `json:"realized_funding_pnl,omitempty"`
}

// Transaction represents a transaction in the account history.
type Transaction struct {
	TransactionID  string  `json:"transaction_id"`
	Time           float64 `json:"time"`
	Type           string  `json:"type"`
	InstrumentName string  `json:"instrument_name,omitempty"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	TradeID        string  `json:"trade_id,omitempty"`
}

// TradeHistoryParams configures a trade history query.
type TradeHistoryParams struct {
	From            *float64   `json:"from,omitempty"`
	To              *float64   `json:"to,omitempty"`
	Offset          *int       `json:"offset,omitempty"`
	Limit           *int       `json:"limit,omitempty"`
	Sort            enums.Sort `json:"sort,omitempty"`
	InstrumentNames []string   `json:"instrument_names,omitempty"`
	BotIDs          []string   `json:"bot_ids,omitempty"`
}

// OrderHistoryParams configures an order history query.
type OrderHistoryParams struct {
	From            *float64   `json:"from,omitempty"`
	To              *float64   `json:"to,omitempty"`
	Offset          *int       `json:"offset,omitempty"`
	Limit           *int       `json:"limit,omitempty"`
	Sort            enums.Sort `json:"sort,omitempty"`
	InstrumentNames []string   `json:"instrument_names,omitempty"`
}

// DailyMarkHistoryParams configures a daily mark history query.
type DailyMarkHistoryParams struct {
	From   *float64 `json:"from,omitempty"`
	To     *float64 `json:"to,omitempty"`
	Offset *int     `json:"offset,omitempty"`
	Limit  *int     `json:"limit,omitempty"`
}

// TransactionHistoryParams configures a transaction history query.
type TransactionHistoryParams struct {
	From   *float64   `json:"from,omitempty"`
	To     *float64   `json:"to,omitempty"`
	Offset *int       `json:"offset,omitempty"`
	Limit  *int       `json:"limit,omitempty"`
	Sort   enums.Sort `json:"sort,omitempty"`
}
