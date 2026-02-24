package types

import "github.com/amiwrpremium/go-thalex/enums"

// Bot represents a bot instance. Use the Strategy field to determine
// which strategy-specific fields are populated.
type Bot struct {
	BotID           string              `json:"bot_id"`
	Strategy        enums.BotStrategy   `json:"strategy"`
	Status          enums.BotStatus     `json:"status"`
	StopReason      enums.BotStopReason `json:"stop_reason,omitempty"`
	InstrumentName  string              `json:"instrument_name"`
	EndTime         *float64            `json:"end_time,omitempty"`
	StartTime       float64             `json:"start_time,omitempty"`
	StopTime        *float64            `json:"stop_time,omitempty"`
	Label           string              `json:"label,omitempty"`
	RealisedPnl     float64             `json:"realized_pnl,omitempty"`
	Fee             float64             `json:"fee,omitempty"`
	AveragePrice    *float64            `json:"average_price,omitempty"`
	PositionSize    *float64            `json:"position_size,omitempty"`
	MarkPriceAtStop *float64            `json:"mark_price_at_stop,omitempty"`

	// SGSL fields
	Signal         enums.Target `json:"signal,omitempty"`
	EntryPrice     *float64     `json:"entry_price,omitempty"`
	TargetPosition *float64     `json:"target_position,omitempty"`
	ExitPrice      *float64     `json:"exit_price,omitempty"`
	ExitPosition   *float64     `json:"exit_position,omitempty"`
	MaxSlippage    *float64     `json:"max_slippage,omitempty"`

	// OCQ fields
	BidOffset   *float64 `json:"bid_offset,omitempty"`
	AskOffset   *float64 `json:"ask_offset,omitempty"`
	ExitOffset  *float64 `json:"exit_offset,omitempty"`
	QuoteSize   *float64 `json:"quote_size,omitempty"`
	MinPosition *float64 `json:"min_position,omitempty"`
	MaxPosition *float64 `json:"max_position,omitempty"`

	// Levels/Grid fields
	Bids              []float64 `json:"bids,omitempty"`
	Asks              []float64 `json:"asks,omitempty"`
	Grid              []float64 `json:"grid,omitempty"`
	StepSize          *float64  `json:"step_size,omitempty"`
	BasePosition      *float64  `json:"base_position,omitempty"`
	TargetMeanPrice   *float64  `json:"target_mean_price,omitempty"`
	UpsideExitPrice   *float64  `json:"upside_exit_price,omitempty"`
	DownsideExitPrice *float64  `json:"downside_exit_price,omitempty"`

	// DHedge/DFollow fields
	Position         string   `json:"position,omitempty"`
	TargetDelta      *float64 `json:"target_delta,omitempty"`
	Threshold        *float64 `json:"threshold,omitempty"`
	Tolerance        *float64 `json:"tolerance,omitempty"`
	Period           *float64 `json:"period,omitempty"`
	TargetInstrument string   `json:"target_instrument,omitempty"`
	TargetAmount     *float64 `json:"target_amount,omitempty"`
}

// SGSLBotParams contains parameters for creating an SGSL bot.
type SGSLBotParams struct {
	Strategy       enums.BotStrategy `json:"strategy"`
	InstrumentName string            `json:"instrument_name"`
	Signal         enums.Target      `json:"signal"`
	EntryPrice     float64           `json:"entry_price"`
	TargetPosition float64           `json:"target_position"`
	ExitPrice      float64           `json:"exit_price"`
	ExitPosition   float64           `json:"exit_position"`
	EndTime        float64           `json:"end_time"`
	MaxSlippage    *float64          `json:"max_slippage,omitempty"`
	Label          string            `json:"label,omitempty"`
}

// NewSGSLBotParams creates SGSL bot parameters.
func NewSGSLBotParams(instrumentName string, signal enums.Target, entryPrice, targetPosition, exitPrice, exitPosition, endTime float64) *SGSLBotParams {
	return &SGSLBotParams{
		Strategy: enums.BotStrategySGSL, InstrumentName: instrumentName, Signal: signal,
		EntryPrice: entryPrice, TargetPosition: targetPosition,
		ExitPrice: exitPrice, ExitPosition: exitPosition, EndTime: endTime,
	}
}

// WithMaxSlippage sets the maximum slippage per trade.
func (p *SGSLBotParams) WithMaxSlippage(v float64) *SGSLBotParams { p.MaxSlippage = &v; return p }

// WithLabel sets a label for bot orders.
func (p *SGSLBotParams) WithLabel(v string) *SGSLBotParams { p.Label = v; return p }

// OCQBotParams contains parameters for creating an OCQ bot.
type OCQBotParams struct {
	Strategy       enums.BotStrategy `json:"strategy"`
	InstrumentName string            `json:"instrument_name"`
	Signal         enums.Target      `json:"signal"`
	BidOffset      float64           `json:"bid_offset"`
	AskOffset      float64           `json:"ask_offset"`
	QuoteSize      float64           `json:"quote_size"`
	MinPosition    float64           `json:"min_position"`
	MaxPosition    float64           `json:"max_position"`
	EndTime        float64           `json:"end_time"`
	ExitOffset     *float64          `json:"exit_offset,omitempty"`
	TargetPosition *float64          `json:"target_position,omitempty"`
	Label          string            `json:"label,omitempty"`
}

// NewOCQBotParams creates OCQ bot parameters.
func NewOCQBotParams(instrumentName string, signal enums.Target, bidOffset, askOffset, quoteSize, minPos, maxPos, endTime float64) *OCQBotParams {
	return &OCQBotParams{
		Strategy: enums.BotStrategyOCQ, InstrumentName: instrumentName, Signal: signal,
		BidOffset: bidOffset, AskOffset: askOffset, QuoteSize: quoteSize,
		MinPosition: minPos, MaxPosition: maxPos, EndTime: endTime,
	}
}

// WithExitOffset sets the exit offset.
func (p *OCQBotParams) WithExitOffset(v float64) *OCQBotParams { p.ExitOffset = &v; return p }

// WithTargetPosition sets the target position.
func (p *OCQBotParams) WithTargetPosition(v float64) *OCQBotParams { p.TargetPosition = &v; return p }

// WithLabel sets a label for bot orders.
func (p *OCQBotParams) WithLabel(v string) *OCQBotParams { p.Label = v; return p }

// LevelsBotParams contains parameters for creating a Levels bot.
type LevelsBotParams struct {
	Strategy          enums.BotStrategy `json:"strategy"`
	InstrumentName    string            `json:"instrument_name"`
	Bids              []float64         `json:"bids"`
	Asks              []float64         `json:"asks"`
	StepSize          float64           `json:"step_size"`
	EndTime           float64           `json:"end_time"`
	BasePosition      *float64          `json:"base_position,omitempty"`
	TargetMeanPrice   *float64          `json:"target_mean_price,omitempty"`
	UpsideExitPrice   *float64          `json:"upside_exit_price,omitempty"`
	DownsideExitPrice *float64          `json:"downside_exit_price,omitempty"`
	MaxSlippage       *float64          `json:"max_slippage,omitempty"`
	Label             string            `json:"label,omitempty"`
}

// NewLevelsBotParams creates Levels bot parameters.
func NewLevelsBotParams(instrumentName string, bids, asks []float64, stepSize, endTime float64) *LevelsBotParams {
	return &LevelsBotParams{
		Strategy: enums.BotStrategyLevels, InstrumentName: instrumentName,
		Bids: bids, Asks: asks, StepSize: stepSize, EndTime: endTime,
	}
}

func (p *LevelsBotParams) WithBasePosition(v float64) *LevelsBotParams { p.BasePosition = &v; return p }
func (p *LevelsBotParams) WithTargetMeanPrice(v float64) *LevelsBotParams {
	p.TargetMeanPrice = &v
	return p
}
func (p *LevelsBotParams) WithUpsideExitPrice(v float64) *LevelsBotParams {
	p.UpsideExitPrice = &v
	return p
}
func (p *LevelsBotParams) WithDownsideExitPrice(v float64) *LevelsBotParams {
	p.DownsideExitPrice = &v
	return p
}
func (p *LevelsBotParams) WithMaxSlippage(v float64) *LevelsBotParams { p.MaxSlippage = &v; return p }
func (p *LevelsBotParams) WithLabel(v string) *LevelsBotParams        { p.Label = v; return p }

// GridBotParams contains parameters for creating a Grid bot.
type GridBotParams struct {
	Strategy          enums.BotStrategy `json:"strategy"`
	InstrumentName    string            `json:"instrument_name"`
	Grid              []float64         `json:"grid"`
	StepSize          float64           `json:"step_size"`
	EndTime           float64           `json:"end_time"`
	BasePosition      *float64          `json:"base_position,omitempty"`
	TargetMeanPrice   *float64          `json:"target_mean_price,omitempty"`
	UpsideExitPrice   *float64          `json:"upside_exit_price,omitempty"`
	DownsideExitPrice *float64          `json:"downside_exit_price,omitempty"`
	MaxSlippage       *float64          `json:"max_slippage,omitempty"`
	Label             string            `json:"label,omitempty"`
}

// NewGridBotParams creates Grid bot parameters.
func NewGridBotParams(instrumentName string, grid []float64, stepSize, endTime float64) *GridBotParams {
	return &GridBotParams{
		Strategy: enums.BotStrategyGrid, InstrumentName: instrumentName,
		Grid: grid, StepSize: stepSize, EndTime: endTime,
	}
}

func (p *GridBotParams) WithBasePosition(v float64) *GridBotParams { p.BasePosition = &v; return p }
func (p *GridBotParams) WithTargetMeanPrice(v float64) *GridBotParams {
	p.TargetMeanPrice = &v
	return p
}
func (p *GridBotParams) WithUpsideExitPrice(v float64) *GridBotParams {
	p.UpsideExitPrice = &v
	return p
}
func (p *GridBotParams) WithDownsideExitPrice(v float64) *GridBotParams {
	p.DownsideExitPrice = &v
	return p
}
func (p *GridBotParams) WithMaxSlippage(v float64) *GridBotParams { p.MaxSlippage = &v; return p }
func (p *GridBotParams) WithLabel(v string) *GridBotParams        { p.Label = v; return p }

// DHedgeBotParams contains parameters for creating a Delta Hedger bot.
type DHedgeBotParams struct {
	Strategy       enums.BotStrategy `json:"strategy"`
	InstrumentName string            `json:"instrument_name"`
	Period         float64           `json:"period"`
	Position       string            `json:"position,omitempty"`
	TargetDelta    *float64          `json:"target_delta,omitempty"`
	Threshold      *float64          `json:"threshold,omitempty"`
	Tolerance      *float64          `json:"tolerance,omitempty"`
	MaxSlippage    *float64          `json:"max_slippage,omitempty"`
	EndTime        *float64          `json:"end_time,omitempty"`
	Label          string            `json:"label,omitempty"`
}

// NewDHedgeBotParams creates DHedge bot parameters.
func NewDHedgeBotParams(instrumentName string, period float64) *DHedgeBotParams {
	return &DHedgeBotParams{
		Strategy: enums.BotStrategyDHedge, InstrumentName: instrumentName, Period: period,
	}
}

func (p *DHedgeBotParams) WithPosition(v string) *DHedgeBotParams     { p.Position = v; return p }
func (p *DHedgeBotParams) WithTargetDelta(v float64) *DHedgeBotParams { p.TargetDelta = &v; return p }
func (p *DHedgeBotParams) WithThreshold(v float64) *DHedgeBotParams   { p.Threshold = &v; return p }
func (p *DHedgeBotParams) WithTolerance(v float64) *DHedgeBotParams   { p.Tolerance = &v; return p }
func (p *DHedgeBotParams) WithMaxSlippage(v float64) *DHedgeBotParams { p.MaxSlippage = &v; return p }
func (p *DHedgeBotParams) WithEndTime(v float64) *DHedgeBotParams     { p.EndTime = &v; return p }
func (p *DHedgeBotParams) WithLabel(v string) *DHedgeBotParams        { p.Label = v; return p }

// DFollowBotParams contains parameters for creating a Delta Follower bot.
type DFollowBotParams struct {
	Strategy         enums.BotStrategy `json:"strategy"`
	InstrumentName   string            `json:"instrument_name"`
	TargetInstrument string            `json:"target_instrument"`
	TargetAmount     float64           `json:"target_amount"`
	Period           float64           `json:"period"`
	EndTime          float64           `json:"end_time"`
	Threshold        *float64          `json:"threshold,omitempty"`
	Tolerance        *float64          `json:"tolerance,omitempty"`
	MaxSlippage      *float64          `json:"max_slippage,omitempty"`
	Label            string            `json:"label,omitempty"`
}

// NewDFollowBotParams creates DFollow bot parameters.
func NewDFollowBotParams(instrumentName, targetInstrument string, targetAmount, period, endTime float64) *DFollowBotParams {
	return &DFollowBotParams{
		Strategy: enums.BotStrategyDFollow, InstrumentName: instrumentName,
		TargetInstrument: targetInstrument, TargetAmount: targetAmount,
		Period: period, EndTime: endTime,
	}
}

func (p *DFollowBotParams) WithThreshold(v float64) *DFollowBotParams   { p.Threshold = &v; return p }
func (p *DFollowBotParams) WithTolerance(v float64) *DFollowBotParams   { p.Tolerance = &v; return p }
func (p *DFollowBotParams) WithMaxSlippage(v float64) *DFollowBotParams { p.MaxSlippage = &v; return p }
func (p *DFollowBotParams) WithLabel(v string) *DFollowBotParams        { p.Label = v; return p }
