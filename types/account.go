package types

// CashHolding represents a single currency balance in an account.
type CashHolding struct {
	Currency             string   `json:"currency"`
	Balance              float64  `json:"balance"`
	CollateralFactor     float64  `json:"collateral_factor"`
	CollateralIndexPrice *float64 `json:"collateral_index_price"`
	Transactable         bool     `json:"transactable"`
}

// AccountSummary represents the overall account financial summary.
type AccountSummary struct {
	Cash               []CashHolding `json:"cash"`
	UnrealisedPnl      float64       `json:"unrealised_pnl"`
	CashCollateral     float64       `json:"cash_collateral"`
	Margin             float64       `json:"margin"`
	RequiredMargin     float64       `json:"required_margin"`
	RemainingMargin    float64       `json:"remaining_margin"`
	SessionRealisedPnl float64       `json:"session_realised_pnl"`
}

// MarginUtilization returns the fraction of margin in use (0.0 to 1.0+).
func (s *AccountSummary) MarginUtilization() float64 {
	if s.Margin == 0 {
		return 0
	}
	return s.RequiredMargin / s.Margin
}

// PortfolioEntry represents a single position in the portfolio.
type PortfolioEntry struct {
	InstrumentName             string   `json:"instrument_name"`
	Position                   float64  `json:"position"`
	MarkPrice                  float64  `json:"mark_price"`
	IV                         *float64 `json:"iv,omitempty"`
	Index                      *float64 `json:"index,omitempty"`
	StartPrice                 float64  `json:"start_price"`
	AveragePrice               float64  `json:"average_price"`
	UnrealisedPnl              float64  `json:"unrealised_pnl"`
	RealisedPnl                float64  `json:"realised_pnl"`
	EntryValue                 float64  `json:"entry_value"`
	PerpetualFundingEntryValue *float64 `json:"perpetual_funding_entry_value,omitempty"`
	UnrealisedPerpetualFunding *float64 `json:"unrealised_perpetual_funding,omitempty"`
}

// IsLong returns true if the position is positive (long).
func (p *PortfolioEntry) IsLong() bool { return p.Position > 0 }

// IsShort returns true if the position is negative (short).
func (p *PortfolioEntry) IsShort() bool { return p.Position < 0 }

// MarginState represents margin information for a single underlying.
type MarginState struct {
	Underlying                   string  `json:"underlying"`
	RequiredMargin               float64 `json:"required_margin"`
	LossMargin                   float64 `json:"loss_margin"`
	RollContingencyMargin        float64 `json:"roll_contingency_margin"`
	D1RollContingencyMargin      float64 `json:"d1_roll_contingency_margin"`
	OptionsRollContingencyMargin float64 `json:"options_roll_contingency_margin"`
	OptionsContingencyMargin     float64 `json:"options_contingency_margin"`
}

// MarginForOrderResult contains margin information with and without a hypothetical order.
type MarginForOrderResult struct {
	Current  MarginBreakdownSide `json:"current"`
	WithBuy  MarginBreakdownSide `json:"with_buy"`
	WithSell MarginBreakdownSide `json:"with_sell"`
}

// MarginBreakdownSide holds margin totals and underlying detail.
type MarginBreakdownSide struct {
	RequiredMargin float64     `json:"required_margin"`
	Underlying     MarginState `json:"underlying"`
}

// ScenarioPosition holds per-instrument scenario simulation results.
type ScenarioPosition struct {
	InstrumentName string  `json:"instrument_name"`
	Position       float64 `json:"position"`
	InstrumentPnl  float64 `json:"instrument_pnl"`
	Pnl            float64 `json:"pnl"`
	CurrentPrice   float64 `json:"current_price"`
	ScenarioPrice  float64 `json:"scenario_price"`
	OpenOrder      bool    `json:"open_order,omitempty"`
	AssumedFilled  *bool   `json:"assumed_filled,omitempty"`
}

// ScenarioAsset holds per-asset scenario simulation results.
type ScenarioAsset struct {
	AssetName     string  `json:"asset_name"`
	Position      float64 `json:"position"`
	UnderlyingPnl float64 `json:"underlying_pnl"`
	Pnl           float64 `json:"pnl"`
	CurrentPrice  float64 `json:"current_price"`
	ScenarioPrice float64 `json:"scenario_price"`
}

// MarginScenario holds a single margin scenario simulation.
type MarginScenario struct {
	UnderlyingChangePct          float64            `json:"underlying_change_pct"`
	VolChangePctPoint            float64            `json:"vol_change_pct_point"`
	Pnl                          float64            `json:"pnl"`
	CoverageFactor               float64            `json:"coverage_factor"`
	RequiredMargin               float64            `json:"required_margin"`
	LossMargin                   float64            `json:"loss_margin"`
	RollContingencyMargin        float64            `json:"roll_contingency_margin"`
	RollCashPosition             float64            `json:"roll_cash_position"`
	D1RollContingencyMargin      float64            `json:"d1_roll_contingency_margin"`
	D1RollCashPosition           float64            `json:"d1_roll_cash_position"`
	OptionsRollContingencyMargin float64            `json:"options_roll_contingency_margin"`
	OptionsRollCashPosition      float64            `json:"options_roll_cash_position"`
	OptionsContingencyMargin     float64            `json:"options_contingency_margin"`
	OptionsShortCashPosition     float64            `json:"options_short_cash_position"`
	Positions                    []ScenarioPosition `json:"positions,omitempty"`
	Assets                       []ScenarioAsset    `json:"assets,omitempty"`
	Highlight                    bool               `json:"highlight,omitempty"`
}

// UnderlyingMarginDetail holds detailed margin info for a single underlying.
type UnderlyingMarginDetail struct {
	Underlying                   string           `json:"underlying"`
	RequiredMargin               float64          `json:"required_margin"`
	LossMargin                   float64          `json:"loss_margin"`
	D1RollCashPosition           float64          `json:"d1_roll_cash_position"`
	OptionsRollCashPosition      float64          `json:"options_roll_cash_position"`
	RollCashPosition             float64          `json:"roll_cash_position"`
	D1RollContingencyMargin      float64          `json:"d1_roll_contingency_margin"`
	OptionsRollContingencyMargin float64          `json:"options_roll_contingency_margin"`
	RollContingencyMargin        float64          `json:"roll_contingency_margin"`
	OptionsShortCashPosition     float64          `json:"options_short_cash_position"`
	OptionsContingencyMargin     float64          `json:"options_contingency_margin"`
	ScenarioUsed                 int              `json:"scenario_used"`
	Scenarios                    []MarginScenario `json:"scenarios,omitempty"`
}

// PortfolioMarginBreakdown contains the full portfolio margin breakdown.
type PortfolioMarginBreakdown struct {
	Portfolio struct {
		RequiredMargin float64                  `json:"required_margin"`
		Underlyings    []UnderlyingMarginDetail `json:"underlyings"`
	} `json:"portfolio"`
}

// AccountBreakdown represents a detailed account breakdown.
type AccountBreakdown struct {
	Cash               []CashHolding    `json:"cash"`
	Portfolio          []PortfolioEntry `json:"portfolio"`
	UnrealisedPnl      float64          `json:"unrealised_pnl"`
	CashCollateral     float64          `json:"cash_collateral"`
	Margin             float64          `json:"margin"`
	RequiredMargin     float64          `json:"required_margin"`
	RemainingMargin    float64          `json:"remaining_margin"`
	SessionRealisedPnl float64          `json:"session_realised_pnl"`
}
