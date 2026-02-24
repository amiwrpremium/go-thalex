package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ========== SGSL Bot ==========

func TestNewSGSLBotParams(t *testing.T) {
	p := types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetMark, 50000, 1.0, 45000, 0, 1700003600)
	if p == nil {
		t.Fatal("NewSGSLBotParams returned nil")
	}
	if p.Strategy != enums.BotStrategySGSL {
		t.Errorf("Strategy = %q, want %q", p.Strategy, enums.BotStrategySGSL)
	}
	if p.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "BTC-PERPETUAL")
	}
	if p.Signal != enums.TargetMark {
		t.Errorf("Signal = %q, want %q", p.Signal, enums.TargetMark)
	}
	if p.EntryPrice != 50000 {
		t.Errorf("EntryPrice = %v, want 50000", p.EntryPrice)
	}
	if p.TargetPosition != 1.0 {
		t.Errorf("TargetPosition = %v, want 1.0", p.TargetPosition)
	}
	if p.ExitPrice != 45000 {
		t.Errorf("ExitPrice = %v, want 45000", p.ExitPrice)
	}
	if p.ExitPosition != 0 {
		t.Errorf("ExitPosition = %v, want 0", p.ExitPosition)
	}
	if p.EndTime != 1700003600 {
		t.Errorf("EndTime = %v, want 1700003600", p.EndTime)
	}
	if p.MaxSlippage != nil {
		t.Errorf("MaxSlippage = %v, want nil", p.MaxSlippage)
	}
	if p.Label != "" {
		t.Errorf("Label = %q, want empty", p.Label)
	}
}

func TestSGSLBotParams_WithMaxSlippage(t *testing.T) {
	p := types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetMark, 50000, 1.0, 45000, 0, 1700003600)
	ret := p.WithMaxSlippage(100.0)
	if ret != p {
		t.Error("WithMaxSlippage should return the same pointer for chaining")
	}
	if p.MaxSlippage == nil || *p.MaxSlippage != 100.0 {
		t.Errorf("MaxSlippage = %v, want 100.0", p.MaxSlippage)
	}
}

func TestSGSLBotParams_WithLabel(t *testing.T) {
	p := types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetMark, 50000, 1.0, 45000, 0, 1700003600)
	ret := p.WithLabel("sgsl-bot")
	if ret != p {
		t.Error("WithLabel should return the same pointer for chaining")
	}
	if p.Label != "sgsl-bot" {
		t.Errorf("Label = %q, want %q", p.Label, "sgsl-bot")
	}
}

func TestSGSLBotParams_JSONRoundTrip(t *testing.T) {
	p := types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetLast, 50000, 1.0, 45000, 0, 1700003600).
		WithMaxSlippage(50.0).
		WithLabel("json-test")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.SGSLBotParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Strategy != enums.BotStrategySGSL {
		t.Errorf("Strategy = %q, want %q", got.Strategy, enums.BotStrategySGSL)
	}
	if got.Signal != enums.TargetLast {
		t.Errorf("Signal = %q, want %q", got.Signal, enums.TargetLast)
	}
	if got.MaxSlippage == nil || *got.MaxSlippage != 50.0 {
		t.Errorf("MaxSlippage = %v, want 50.0", got.MaxSlippage)
	}
	if got.Label != "json-test" {
		t.Errorf("Label = %q, want %q", got.Label, "json-test")
	}
}

// ========== OCQ Bot ==========

func TestNewOCQBotParams(t *testing.T) {
	p := types.NewOCQBotParams("ETH-PERPETUAL", enums.TargetIndex, 10, 15, 0.5, -5, 5, 1700003600)
	if p == nil {
		t.Fatal("NewOCQBotParams returned nil")
	}
	if p.Strategy != enums.BotStrategyOCQ {
		t.Errorf("Strategy = %q, want %q", p.Strategy, enums.BotStrategyOCQ)
	}
	if p.InstrumentName != "ETH-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "ETH-PERPETUAL")
	}
	if p.BidOffset != 10 {
		t.Errorf("BidOffset = %v, want 10", p.BidOffset)
	}
	if p.AskOffset != 15 {
		t.Errorf("AskOffset = %v, want 15", p.AskOffset)
	}
	if p.QuoteSize != 0.5 {
		t.Errorf("QuoteSize = %v, want 0.5", p.QuoteSize)
	}
	if p.MinPosition != -5 {
		t.Errorf("MinPosition = %v, want -5", p.MinPosition)
	}
	if p.MaxPosition != 5 {
		t.Errorf("MaxPosition = %v, want 5", p.MaxPosition)
	}
}

func TestOCQBotParams_WithExitOffset(t *testing.T) {
	p := types.NewOCQBotParams("ETH-PERPETUAL", enums.TargetIndex, 10, 15, 0.5, -5, 5, 1700003600)
	ret := p.WithExitOffset(5.0)
	if ret != p {
		t.Error("WithExitOffset should return the same pointer for chaining")
	}
	if p.ExitOffset == nil || *p.ExitOffset != 5.0 {
		t.Errorf("ExitOffset = %v, want 5.0", p.ExitOffset)
	}
}

func TestOCQBotParams_WithTargetPosition(t *testing.T) {
	p := types.NewOCQBotParams("ETH-PERPETUAL", enums.TargetIndex, 10, 15, 0.5, -5, 5, 1700003600)
	ret := p.WithTargetPosition(2.0)
	if ret != p {
		t.Error("WithTargetPosition should return the same pointer for chaining")
	}
	if p.TargetPosition == nil || *p.TargetPosition != 2.0 {
		t.Errorf("TargetPosition = %v, want 2.0", p.TargetPosition)
	}
}

func TestOCQBotParams_WithLabel(t *testing.T) {
	p := types.NewOCQBotParams("ETH-PERPETUAL", enums.TargetIndex, 10, 15, 0.5, -5, 5, 1700003600)
	ret := p.WithLabel("ocq-bot")
	if ret != p {
		t.Error("WithLabel should return the same pointer for chaining")
	}
	if p.Label != "ocq-bot" {
		t.Errorf("Label = %q, want %q", p.Label, "ocq-bot")
	}
}

func TestOCQBotParams_JSONRoundTrip(t *testing.T) {
	p := types.NewOCQBotParams("ETH-PERPETUAL", enums.TargetIndex, 10, 15, 0.5, -5, 5, 1700003600).
		WithExitOffset(3.0).
		WithTargetPosition(1.0).
		WithLabel("ocq-json")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OCQBotParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Strategy != enums.BotStrategyOCQ {
		t.Errorf("Strategy = %q, want %q", got.Strategy, enums.BotStrategyOCQ)
	}
	if got.ExitOffset == nil || *got.ExitOffset != 3.0 {
		t.Errorf("ExitOffset = %v, want 3.0", got.ExitOffset)
	}
	if got.TargetPosition == nil || *got.TargetPosition != 1.0 {
		t.Errorf("TargetPosition = %v, want 1.0", got.TargetPosition)
	}
}

// ========== Levels Bot ==========

func TestNewLevelsBotParams(t *testing.T) {
	bids := []float64{49000, 48500, 48000}
	asks := []float64{51000, 51500, 52000}
	p := types.NewLevelsBotParams("BTC-PERPETUAL", bids, asks, 100, 1700003600)
	if p == nil {
		t.Fatal("NewLevelsBotParams returned nil")
	}
	if p.Strategy != enums.BotStrategyLevels {
		t.Errorf("Strategy = %q, want %q", p.Strategy, enums.BotStrategyLevels)
	}
	if len(p.Bids) != 3 {
		t.Errorf("len(Bids) = %d, want 3", len(p.Bids))
	}
	if len(p.Asks) != 3 {
		t.Errorf("len(Asks) = %d, want 3", len(p.Asks))
	}
	if p.StepSize != 100 {
		t.Errorf("StepSize = %v, want 100", p.StepSize)
	}
}

func TestLevelsBotParams_Builders(t *testing.T) {
	p := types.NewLevelsBotParams("BTC-PERPETUAL", nil, nil, 100, 1700003600)

	ret := p.WithBasePosition(1.0)
	if ret != p {
		t.Error("WithBasePosition should return the same pointer")
	}
	if p.BasePosition == nil || *p.BasePosition != 1.0 {
		t.Errorf("BasePosition = %v, want 1.0", p.BasePosition)
	}

	p.WithTargetMeanPrice(50000.0)
	if p.TargetMeanPrice == nil || *p.TargetMeanPrice != 50000.0 {
		t.Errorf("TargetMeanPrice = %v, want 50000.0", p.TargetMeanPrice)
	}

	p.WithUpsideExitPrice(55000.0)
	if p.UpsideExitPrice == nil || *p.UpsideExitPrice != 55000.0 {
		t.Errorf("UpsideExitPrice = %v, want 55000.0", p.UpsideExitPrice)
	}

	p.WithDownsideExitPrice(40000.0)
	if p.DownsideExitPrice == nil || *p.DownsideExitPrice != 40000.0 {
		t.Errorf("DownsideExitPrice = %v, want 40000.0", p.DownsideExitPrice)
	}

	p.WithMaxSlippage(200.0)
	if p.MaxSlippage == nil || *p.MaxSlippage != 200.0 {
		t.Errorf("MaxSlippage = %v, want 200.0", p.MaxSlippage)
	}

	p.WithLabel("levels-bot")
	if p.Label != "levels-bot" {
		t.Errorf("Label = %q, want %q", p.Label, "levels-bot")
	}
}

func TestLevelsBotParams_JSONRoundTrip(t *testing.T) {
	p := types.NewLevelsBotParams("BTC-PERPETUAL", []float64{49000}, []float64{51000}, 100, 1700003600).
		WithBasePosition(0.5).
		WithLabel("levels-json")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.LevelsBotParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Strategy != enums.BotStrategyLevels {
		t.Errorf("Strategy = %q, want %q", got.Strategy, enums.BotStrategyLevels)
	}
	if got.BasePosition == nil || *got.BasePosition != 0.5 {
		t.Errorf("BasePosition = %v, want 0.5", got.BasePosition)
	}
}

// ========== Grid Bot ==========

func TestNewGridBotParams(t *testing.T) {
	grid := []float64{48000, 49000, 50000, 51000, 52000}
	p := types.NewGridBotParams("BTC-PERPETUAL", grid, 50, 1700003600)
	if p == nil {
		t.Fatal("NewGridBotParams returned nil")
	}
	if p.Strategy != enums.BotStrategyGrid {
		t.Errorf("Strategy = %q, want %q", p.Strategy, enums.BotStrategyGrid)
	}
	if len(p.Grid) != 5 {
		t.Errorf("len(Grid) = %d, want 5", len(p.Grid))
	}
	if p.StepSize != 50 {
		t.Errorf("StepSize = %v, want 50", p.StepSize)
	}
}

func TestGridBotParams_Builders(t *testing.T) {
	p := types.NewGridBotParams("BTC-PERPETUAL", []float64{50000}, 50, 1700003600)

	ret := p.WithBasePosition(2.0)
	if ret != p {
		t.Error("WithBasePosition should return the same pointer")
	}
	if p.BasePosition == nil || *p.BasePosition != 2.0 {
		t.Errorf("BasePosition = %v, want 2.0", p.BasePosition)
	}

	p.WithTargetMeanPrice(50000.0)
	if p.TargetMeanPrice == nil || *p.TargetMeanPrice != 50000.0 {
		t.Errorf("TargetMeanPrice = %v, want 50000.0", p.TargetMeanPrice)
	}

	p.WithUpsideExitPrice(60000.0)
	if p.UpsideExitPrice == nil || *p.UpsideExitPrice != 60000.0 {
		t.Errorf("UpsideExitPrice = %v, want 60000.0", p.UpsideExitPrice)
	}

	p.WithDownsideExitPrice(40000.0)
	if p.DownsideExitPrice == nil || *p.DownsideExitPrice != 40000.0 {
		t.Errorf("DownsideExitPrice = %v, want 40000.0", p.DownsideExitPrice)
	}

	p.WithMaxSlippage(150.0)
	if p.MaxSlippage == nil || *p.MaxSlippage != 150.0 {
		t.Errorf("MaxSlippage = %v, want 150.0", p.MaxSlippage)
	}

	p.WithLabel("grid-bot")
	if p.Label != "grid-bot" {
		t.Errorf("Label = %q, want %q", p.Label, "grid-bot")
	}
}

func TestGridBotParams_JSONRoundTrip(t *testing.T) {
	p := types.NewGridBotParams("BTC-PERPETUAL", []float64{48000, 50000, 52000}, 100, 1700003600).
		WithLabel("grid-json")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.GridBotParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Strategy != enums.BotStrategyGrid {
		t.Errorf("Strategy = %q, want %q", got.Strategy, enums.BotStrategyGrid)
	}
	if len(got.Grid) != 3 {
		t.Errorf("len(Grid) = %d, want 3", len(got.Grid))
	}
}

// ========== DHedge Bot ==========

func TestNewDHedgeBotParams(t *testing.T) {
	p := types.NewDHedgeBotParams("BTC-PERPETUAL", 60.0)
	if p == nil {
		t.Fatal("NewDHedgeBotParams returned nil")
	}
	if p.Strategy != enums.BotStrategyDHedge {
		t.Errorf("Strategy = %q, want %q", p.Strategy, enums.BotStrategyDHedge)
	}
	if p.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "BTC-PERPETUAL")
	}
	if p.Period != 60.0 {
		t.Errorf("Period = %v, want 60.0", p.Period)
	}
}

func TestDHedgeBotParams_Builders(t *testing.T) {
	p := types.NewDHedgeBotParams("BTC-PERPETUAL", 60.0)

	ret := p.WithPosition("BTCUSD")
	if ret != p {
		t.Error("WithPosition should return the same pointer")
	}
	if p.Position != "BTCUSD" {
		t.Errorf("Position = %q, want %q", p.Position, "BTCUSD")
	}

	p.WithTargetDelta(0.5)
	if p.TargetDelta == nil || *p.TargetDelta != 0.5 {
		t.Errorf("TargetDelta = %v, want 0.5", p.TargetDelta)
	}

	p.WithThreshold(0.1)
	if p.Threshold == nil || *p.Threshold != 0.1 {
		t.Errorf("Threshold = %v, want 0.1", p.Threshold)
	}

	p.WithTolerance(0.01)
	if p.Tolerance == nil || *p.Tolerance != 0.01 {
		t.Errorf("Tolerance = %v, want 0.01", p.Tolerance)
	}

	p.WithMaxSlippage(50.0)
	if p.MaxSlippage == nil || *p.MaxSlippage != 50.0 {
		t.Errorf("MaxSlippage = %v, want 50.0", p.MaxSlippage)
	}

	p.WithEndTime(1700003600.0)
	if p.EndTime == nil || *p.EndTime != 1700003600.0 {
		t.Errorf("EndTime = %v, want 1700003600.0", p.EndTime)
	}

	p.WithLabel("dhedge-bot")
	if p.Label != "dhedge-bot" {
		t.Errorf("Label = %q, want %q", p.Label, "dhedge-bot")
	}
}

func TestDHedgeBotParams_JSONRoundTrip(t *testing.T) {
	p := types.NewDHedgeBotParams("BTC-PERPETUAL", 60.0).
		WithTargetDelta(0.5).
		WithLabel("dhedge-json")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.DHedgeBotParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Strategy != enums.BotStrategyDHedge {
		t.Errorf("Strategy = %q, want %q", got.Strategy, enums.BotStrategyDHedge)
	}
	if got.TargetDelta == nil || *got.TargetDelta != 0.5 {
		t.Errorf("TargetDelta = %v, want 0.5", got.TargetDelta)
	}
}

// ========== DFollow Bot ==========

func TestNewDFollowBotParams(t *testing.T) {
	p := types.NewDFollowBotParams("BTC-PERPETUAL", "ETH-PERPETUAL", 10.0, 30.0, 1700003600)
	if p == nil {
		t.Fatal("NewDFollowBotParams returned nil")
	}
	if p.Strategy != enums.BotStrategyDFollow {
		t.Errorf("Strategy = %q, want %q", p.Strategy, enums.BotStrategyDFollow)
	}
	if p.InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("InstrumentName = %q, want %q", p.InstrumentName, "BTC-PERPETUAL")
	}
	if p.TargetInstrument != "ETH-PERPETUAL" {
		t.Errorf("TargetInstrument = %q, want %q", p.TargetInstrument, "ETH-PERPETUAL")
	}
	if p.TargetAmount != 10.0 {
		t.Errorf("TargetAmount = %v, want 10.0", p.TargetAmount)
	}
	if p.Period != 30.0 {
		t.Errorf("Period = %v, want 30.0", p.Period)
	}
	if p.EndTime != 1700003600 {
		t.Errorf("EndTime = %v, want 1700003600", p.EndTime)
	}
}

func TestDFollowBotParams_Builders(t *testing.T) {
	p := types.NewDFollowBotParams("BTC-PERPETUAL", "ETH-PERPETUAL", 10.0, 30.0, 1700003600)

	ret := p.WithThreshold(0.2)
	if ret != p {
		t.Error("WithThreshold should return the same pointer")
	}
	if p.Threshold == nil || *p.Threshold != 0.2 {
		t.Errorf("Threshold = %v, want 0.2", p.Threshold)
	}

	p.WithTolerance(0.05)
	if p.Tolerance == nil || *p.Tolerance != 0.05 {
		t.Errorf("Tolerance = %v, want 0.05", p.Tolerance)
	}

	p.WithMaxSlippage(25.0)
	if p.MaxSlippage == nil || *p.MaxSlippage != 25.0 {
		t.Errorf("MaxSlippage = %v, want 25.0", p.MaxSlippage)
	}

	p.WithLabel("dfollow-bot")
	if p.Label != "dfollow-bot" {
		t.Errorf("Label = %q, want %q", p.Label, "dfollow-bot")
	}
}

func TestDFollowBotParams_JSONRoundTrip(t *testing.T) {
	p := types.NewDFollowBotParams("BTC-PERPETUAL", "ETH-PERPETUAL", 10.0, 30.0, 1700003600).
		WithThreshold(0.1).
		WithTolerance(0.01).
		WithMaxSlippage(50.0).
		WithLabel("dfollow-json")

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.DFollowBotParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Strategy != enums.BotStrategyDFollow {
		t.Errorf("Strategy = %q, want %q", got.Strategy, enums.BotStrategyDFollow)
	}
	if got.Threshold == nil || *got.Threshold != 0.1 {
		t.Errorf("Threshold = %v, want 0.1", got.Threshold)
	}
	if got.Tolerance == nil || *got.Tolerance != 0.01 {
		t.Errorf("Tolerance = %v, want 0.01", got.Tolerance)
	}
	if got.MaxSlippage == nil || *got.MaxSlippage != 50.0 {
		t.Errorf("MaxSlippage = %v, want 50.0", got.MaxSlippage)
	}
	if got.Label != "dfollow-json" {
		t.Errorf("Label = %q, want %q", got.Label, "dfollow-json")
	}
}

// ========== Bot JSON round-trip ==========

func TestBot_JSONRoundTrip(t *testing.T) {
	endTime := 1700003600.0
	stopTime := 1700002000.0
	avgPrice := 50000.0
	posSize := 2.5
	markPriceStop := 49500.0
	entryPrice := 50000.0
	targetPos := 1.0
	exitPrice := 45000.0
	exitPos := 0.0
	maxSlip := 100.0

	bot := types.Bot{
		BotID:           "bot-abc",
		Strategy:        enums.BotStrategySGSL,
		Status:          enums.BotStatusActive,
		StopReason:      enums.BotStopReason(""),
		InstrumentName:  "BTC-PERPETUAL",
		EndTime:         &endTime,
		StartTime:       1700000000.0,
		StopTime:        &stopTime,
		Label:           "my-bot",
		RealisedPnl:     50.0,
		Fee:             2.0,
		AveragePrice:    &avgPrice,
		PositionSize:    &posSize,
		MarkPriceAtStop: &markPriceStop,
		Signal:          enums.TargetMark,
		EntryPrice:      &entryPrice,
		TargetPosition:  &targetPos,
		ExitPrice:       &exitPrice,
		ExitPosition:    &exitPos,
		MaxSlippage:     &maxSlip,
	}

	data, err := json.Marshal(bot)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.Bot
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.BotID != bot.BotID {
		t.Errorf("BotID = %q, want %q", got.BotID, bot.BotID)
	}
	if got.Strategy != bot.Strategy {
		t.Errorf("Strategy = %q, want %q", got.Strategy, bot.Strategy)
	}
	if got.Status != bot.Status {
		t.Errorf("Status = %q, want %q", got.Status, bot.Status)
	}
	if got.EndTime == nil || *got.EndTime != endTime {
		t.Errorf("EndTime = %v, want %v", got.EndTime, endTime)
	}
	if got.AveragePrice == nil || *got.AveragePrice != avgPrice {
		t.Errorf("AveragePrice = %v, want %v", got.AveragePrice, avgPrice)
	}
	if got.EntryPrice == nil || *got.EntryPrice != entryPrice {
		t.Errorf("EntryPrice = %v, want %v", got.EntryPrice, entryPrice)
	}
	if got.MaxSlippage == nil || *got.MaxSlippage != maxSlip {
		t.Errorf("MaxSlippage = %v, want %v", got.MaxSlippage, maxSlip)
	}
}
