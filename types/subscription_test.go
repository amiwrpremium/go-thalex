package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

func TestBookChannel(t *testing.T) {
	tests := []struct {
		name       string
		instrument string
		grouping   int
		nlevels    int
		delay      enums.Delay
		want       string
	}{
		{
			name:       "perpetual_100ms",
			instrument: "BTC-PERPETUAL",
			grouping:   1,
			nlevels:    10,
			delay:      enums.Delay100ms,
			want:       "book.BTC-PERPETUAL.1.10.100ms",
		},
		{
			name:       "perpetual_raw",
			instrument: "ETH-PERPETUAL",
			grouping:   5,
			nlevels:    20,
			delay:      enums.DelayNone,
			want:       "book.ETH-PERPETUAL.5.20.raw",
		},
		{
			name:       "future_1000ms",
			instrument: "BTC-28MAR25",
			grouping:   10,
			nlevels:    5,
			delay:      enums.Delay1000ms,
			want:       "book.BTC-28MAR25.10.5.1000ms",
		},
		{
			name:       "option_instrument",
			instrument: "BTC-28MAR25-100000-C",
			grouping:   1,
			nlevels:    1,
			delay:      enums.Delay100ms,
			want:       "book.BTC-28MAR25-100000-C.1.1.100ms",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.BookChannel(tt.instrument, tt.grouping, tt.nlevels, tt.delay)
			if got != tt.want {
				t.Errorf("BookChannel(%q, %d, %d, %q) = %q, want %q",
					tt.instrument, tt.grouping, tt.nlevels, tt.delay, got, tt.want)
			}
		})
	}
}

func TestTickerChannel(t *testing.T) {
	tests := []struct {
		name       string
		instrument string
		delay      enums.Delay
		want       string
	}{
		{
			name:       "perpetual_100ms",
			instrument: "BTC-PERPETUAL",
			delay:      enums.Delay100ms,
			want:       "ticker.BTC-PERPETUAL.100ms",
		},
		{
			name:       "perpetual_raw",
			instrument: "ETH-PERPETUAL",
			delay:      enums.DelayNone,
			want:       "ticker.ETH-PERPETUAL.raw",
		},
		{
			name:       "future_1000ms",
			instrument: "BTC-28MAR25",
			delay:      enums.Delay1000ms,
			want:       "ticker.BTC-28MAR25.1000ms",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.TickerChannel(tt.instrument, tt.delay)
			if got != tt.want {
				t.Errorf("TickerChannel(%q, %q) = %q, want %q",
					tt.instrument, tt.delay, got, tt.want)
			}
		})
	}
}

func TestLWTChannel(t *testing.T) {
	tests := []struct {
		name       string
		instrument string
		delay      enums.Delay
		want       string
	}{
		{
			name:       "perpetual_100ms",
			instrument: "BTC-PERPETUAL",
			delay:      enums.Delay100ms,
			want:       "lwt.BTC-PERPETUAL.100ms",
		},
		{
			name:       "perpetual_raw",
			instrument: "ETH-PERPETUAL",
			delay:      enums.DelayNone,
			want:       "lwt.ETH-PERPETUAL.raw",
		},
		{
			name:       "future_1000ms",
			instrument: "BTC-28MAR25",
			delay:      enums.Delay1000ms,
			want:       "lwt.BTC-28MAR25.1000ms",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.LWTChannel(tt.instrument, tt.delay)
			if got != tt.want {
				t.Errorf("LWTChannel(%q, %q) = %q, want %q",
					tt.instrument, tt.delay, got, tt.want)
			}
		})
	}
}

func TestRecentTradesChannel(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		category enums.RecentTradesCategory
		want     string
	}{
		{
			name:     "btcusd_all",
			target:   "BTCUSD",
			category: enums.RecentTradesCategoryAll,
			want:     "recent_trades.BTCUSD.all",
		},
		{
			name:     "ethusd_normal",
			target:   "ETHUSD",
			category: enums.RecentTradesCategoryNormal,
			want:     "recent_trades.ETHUSD.normal",
		},
		{
			name:     "btcusd_block",
			target:   "BTCUSD",
			category: enums.RecentTradesCategoryBlock,
			want:     "recent_trades.BTCUSD.block",
		},
		{
			name:     "btcusd_combo",
			target:   "BTCUSD",
			category: enums.RecentTradesCategoryCombo,
			want:     "recent_trades.BTCUSD.combo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.RecentTradesChannel(tt.target, tt.category)
			if got != tt.want {
				t.Errorf("RecentTradesChannel(%q, %q) = %q, want %q",
					tt.target, tt.category, got, tt.want)
			}
		})
	}
}

func TestPriceIndexChannel(t *testing.T) {
	tests := []struct {
		name       string
		underlying string
		want       string
	}{
		{"btcusd", "BTCUSD", "price_index.BTCUSD"},
		{"ethusd", "ETHUSD", "price_index.ETHUSD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.PriceIndexChannel(tt.underlying)
			if got != tt.want {
				t.Errorf("PriceIndexChannel(%q) = %q, want %q", tt.underlying, got, tt.want)
			}
		})
	}
}

func TestUnderlyingStatisticsChannel(t *testing.T) {
	tests := []struct {
		name       string
		underlying string
		want       string
	}{
		{"btcusd", "BTCUSD", "underlying_statistics.BTCUSD"},
		{"ethusd", "ETHUSD", "underlying_statistics.ETHUSD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.UnderlyingStatisticsChannel(tt.underlying)
			if got != tt.want {
				t.Errorf("UnderlyingStatisticsChannel(%q) = %q, want %q", tt.underlying, got, tt.want)
			}
		})
	}
}

func TestBasePriceChannel(t *testing.T) {
	tests := []struct {
		name       string
		underlying string
		expiration string
		want       string
	}{
		{
			name:       "btcusd_march",
			underlying: "BTCUSD",
			expiration: "2025-03-28",
			want:       "base_price.BTCUSD.2025-03-28",
		},
		{
			name:       "ethusd_june",
			underlying: "ETHUSD",
			expiration: "2025-06-27",
			want:       "base_price.ETHUSD.2025-06-27",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.BasePriceChannel(tt.underlying, tt.expiration)
			if got != tt.want {
				t.Errorf("BasePriceChannel(%q, %q) = %q, want %q",
					tt.underlying, tt.expiration, got, tt.want)
			}
		})
	}
}

func TestIndexComponentsChannel(t *testing.T) {
	tests := []struct {
		name       string
		underlying string
		want       string
	}{
		{"btcusd", "BTCUSD", "index_components.BTCUSD"},
		{"ethusd", "ETHUSD", "index_components.ETHUSD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.IndexComponentsChannel(tt.underlying)
			if got != tt.want {
				t.Errorf("IndexComponentsChannel(%q) = %q, want %q", tt.underlying, got, tt.want)
			}
		})
	}
}

func TestChannelConstants(t *testing.T) {
	// Verify all constant channel names have the expected values.
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"ChannelInstruments", types.ChannelInstruments, "instruments"},
		{"ChannelRfqs", types.ChannelRfqs, "rfqs"},
		{"ChannelSystem", types.ChannelSystem, "system"},
		{"ChannelBanners", types.ChannelBanners, "banners"},
		{"ChannelAccountOrders", types.ChannelAccountOrders, "account.orders"},
		{"ChannelAccountPersistent", types.ChannelAccountPersistent, "account.persistent_orders"},
		{"ChannelSessionOrders", types.ChannelSessionOrders, "session.orders"},
		{"ChannelAccountTradeHistory", types.ChannelAccountTradeHistory, "account.trade_history"},
		{"ChannelAccountOrderHistory", types.ChannelAccountOrderHistory, "account.order_history"},
		{"ChannelAccountPortfolio", types.ChannelAccountPortfolio, "account.portfolio"},
		{"ChannelAccountSummary", types.ChannelAccountSummary, "account.summary"},
		{"ChannelAccountRfqs", types.ChannelAccountRfqs, "account.rfqs"},
		{"ChannelAccountRfqHistory", types.ChannelAccountRfqHistory, "account.rfq_history"},
		{"ChannelAccountConditional", types.ChannelAccountConditional, "account.conditional_orders"},
		{"ChannelAccountBots", types.ChannelAccountBots, "account.bots"},
		{"ChannelUserNotifications", types.ChannelUserNotifications, "user.inbox_notifications"},
		{"ChannelSessionMMProtection", types.ChannelSessionMMProtection, "session.mm_protection"},
		{"ChannelMMRfqs", types.ChannelMMRfqs, "mm.rfqs"},
		{"ChannelMMRfqQuotes", types.ChannelMMRfqQuotes, "mm.rfq_quotes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.want)
			}
		})
	}
}

// ---------- BookUpdate JSON round-trip ----------

func TestBookUpdate_JSONRoundTrip(t *testing.T) {
	lastPrice := 50000.0
	bu := types.BookUpdate{
		Bids: []types.BookLevel{
			{49000.0, 2.0, 1.5},
		},
		Asks: []types.BookLevel{
			{51000.0, 1.0, 0.8},
		},
		Last: &lastPrice,
		Time: 1700000000.0,
		Trades: []types.BookTrade{
			{Direction: enums.DirectionBuy, Price: 50000.0, Amount: 0.5, Time: 1700000000.0},
		},
	}

	data, err := json.Marshal(bu)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.BookUpdate
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Bids) != 1 {
		t.Fatalf("len(Bids) = %d, want 1", len(got.Bids))
	}
	if len(got.Asks) != 1 {
		t.Fatalf("len(Asks) = %d, want 1", len(got.Asks))
	}
	if got.Last == nil || *got.Last != lastPrice {
		t.Errorf("Last = %v, want %v", got.Last, lastPrice)
	}
	if len(got.Trades) != 1 {
		t.Fatalf("len(Trades) = %d, want 1", len(got.Trades))
	}
	if got.Trades[0].Direction != enums.DirectionBuy {
		t.Errorf("Trades[0].Direction = %q, want %q", got.Trades[0].Direction, enums.DirectionBuy)
	}
}

// ---------- BookTrade JSON round-trip ----------

func TestBookTrade_JSONRoundTrip(t *testing.T) {
	bt := types.BookTrade{
		Direction: enums.DirectionSell,
		Price:     50000.0,
		Amount:    1.5,
		Time:      1700000000.0,
	}

	data, err := json.Marshal(bt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.BookTrade
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Direction != enums.DirectionSell {
		t.Errorf("Direction = %q, want %q", got.Direction, enums.DirectionSell)
	}
	if got.Price != 50000.0 {
		t.Errorf("Price = %v, want 50000.0", got.Price)
	}
}

// ---------- LightweightTicker JSON round-trip ----------

func TestLightweightTicker_JSONRoundTrip(t *testing.T) {
	bid := 49000.0
	ask := 51000.0
	iv := 0.65
	last := 50000.0
	bidAmt := 5.0
	askAmt := 3.0

	lwt := types.LightweightTicker{
		BestBidPrice:  &bid,
		BestAskPrice:  &ask,
		MarkPrice:     50050.0,
		IV:            &iv,
		LastPrice:     &last,
		BestBidAmount: &bidAmt,
		BestAskAmount: &askAmt,
	}

	data, err := json.Marshal(lwt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.LightweightTicker
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.BestBidPrice == nil || *got.BestBidPrice != bid {
		t.Errorf("BestBidPrice = %v, want %v", got.BestBidPrice, bid)
	}
	if got.MarkPrice != 50050.0 {
		t.Errorf("MarkPrice = %v, want 50050.0", got.MarkPrice)
	}
	if got.IV == nil || *got.IV != iv {
		t.Errorf("IV = %v, want %v", got.IV, iv)
	}
}

func TestLightweightTicker_NilFields(t *testing.T) {
	lwt := types.LightweightTicker{
		MarkPrice: 50000.0,
	}

	data, err := json.Marshal(lwt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.LightweightTicker
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.BestBidPrice != nil {
		t.Errorf("BestBidPrice = %v, want nil", got.BestBidPrice)
	}
	if got.IV != nil {
		t.Errorf("IV = %v, want nil", got.IV)
	}
}

// ---------- RecentTrade JSON round-trip ----------

func TestRecentTrade_JSONRoundTrip(t *testing.T) {
	idx := 50000.0
	rt := types.RecentTrade{
		TradeID:        "rt-123",
		InstrumentName: "BTC-PERPETUAL",
		Direction:      enums.DirectionBuy,
		Price:          50000.0,
		Amount:         1.0,
		Time:           1700000000.0,
		TradeType:      enums.TradeTypeNormal,
		Index:          &idx,
	}

	data, err := json.Marshal(rt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.RecentTrade
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.TradeID != "rt-123" {
		t.Errorf("TradeID = %q, want %q", got.TradeID, "rt-123")
	}
	if got.TradeType != enums.TradeTypeNormal {
		t.Errorf("TradeType = %q, want %q", got.TradeType, enums.TradeTypeNormal)
	}
	if got.Index == nil || *got.Index != idx {
		t.Errorf("Index = %v, want %v", got.Index, idx)
	}
}

// ---------- UnderlyingStatistics JSON round-trip ----------

func TestUnderlyingStatistics_JSONRoundTrip(t *testing.T) {
	us := types.UnderlyingStatistics{
		Underlying: "BTCUSD",
		Time:       1700000000.0,
		OpenInterest: types.OpenInterest{
			Futures: 100.0,
			Calls:   200.0,
			Puts:    150.0,
		},
	}

	data, err := json.Marshal(us)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.UnderlyingStatistics
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Underlying != "BTCUSD" {
		t.Errorf("Underlying = %q, want %q", got.Underlying, "BTCUSD")
	}
	if got.Futures != 100.0 {
		t.Errorf("Futures = %v, want 100.0", got.Futures)
	}
	if got.Calls != 200.0 {
		t.Errorf("Calls = %v, want 200.0", got.Calls)
	}
}

// ---------- OpenInterestSummary JSON round-trip ----------

func TestOpenInterestSummary_JSONRoundTrip(t *testing.T) {
	ois := types.OpenInterestSummary{
		TotalOI:   450.0,
		PutOI:     150.0,
		CallOI:    200.0,
		FuturesOI: 100.0,
		Expirations: []types.ExpirationOI{
			{ExpiryDate: "2025-03-28", TotalOI: 300.0, PutOI: 100.0, CallOI: 150.0, FuturesOI: 50.0},
		},
	}

	data, err := json.Marshal(ois)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.OpenInterestSummary
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.TotalOI != 450.0 {
		t.Errorf("TotalOI = %v, want 450.0", got.TotalOI)
	}
	if len(got.Expirations) != 1 {
		t.Fatalf("len(Expirations) = %d, want 1", len(got.Expirations))
	}
	if got.Expirations[0].ExpiryDate != "2025-03-28" {
		t.Errorf("Expirations[0].ExpiryDate = %q, want %q", got.Expirations[0].ExpiryDate, "2025-03-28")
	}
}

// ---------- BasePrice JSON round-trip ----------

func TestBasePrice_JSONRoundTrip(t *testing.T) {
	bp := types.BasePrice{
		Underlying: "BTCUSD",
		Expiration: "2025-03-28",
		Price:      51000.0,
		Time:       1700000000.0,
	}

	data, err := json.Marshal(bp)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.BasePrice
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Underlying != "BTCUSD" {
		t.Errorf("Underlying = %q, want %q", got.Underlying, "BTCUSD")
	}
	if got.Price != 51000.0 {
		t.Errorf("Price = %v, want 51000.0", got.Price)
	}
}

// ---------- InstrumentChange JSON round-trip ----------

func TestInstrumentChange_JSONRoundTrip(t *testing.T) {
	ic := types.InstrumentChange{
		Instruments: []types.Instrument{
			{InstrumentName: "BTC-PERPETUAL", Type: enums.InstrumentTypePerpetual},
		},
	}

	data, err := json.Marshal(ic)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.InstrumentChange
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(got.Instruments) != 1 {
		t.Fatalf("len(Instruments) = %d, want 1", len(got.Instruments))
	}
	if got.Instruments[0].InstrumentName != "BTC-PERPETUAL" {
		t.Errorf("Instruments[0].InstrumentName = %q, want %q", got.Instruments[0].InstrumentName, "BTC-PERPETUAL")
	}
}

// ---------- IndexComponents JSON round-trip ----------

func TestIndexComponents_JSONRoundTrip(t *testing.T) {
	ic := types.IndexComponents{
		IndexName: "BTCUSD",
		Price:     50000.0,
		Time:      1700000000.0,
		Components: []types.IndexComponent{
			{Exchange: "exchange1", Price: 50000.0, Weight: 0.5, Stale: false},
			{Exchange: "exchange2", Price: 50100.0, Weight: 0.5, Stale: true},
		},
	}

	data, err := json.Marshal(ic)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.IndexComponents
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.IndexName != "BTCUSD" {
		t.Errorf("IndexName = %q, want %q", got.IndexName, "BTCUSD")
	}
	if len(got.Components) != 2 {
		t.Fatalf("len(Components) = %d, want 2", len(got.Components))
	}
	if !got.Components[1].Stale {
		t.Error("Components[1].Stale = false, want true")
	}
}

// ---------- MMProtectionUpdate JSON round-trip ----------

func TestMMProtectionUpdate_JSONRoundTrip(t *testing.T) {
	mmp := types.MMProtectionUpdate{
		Product: enums.Product("FBTCUSD"),
		Reason:  enums.MMProtectionReasonTriggered,
		Time:    1700000000.0,
	}

	data, err := json.Marshal(mmp)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.MMProtectionUpdate
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Reason != enums.MMProtectionReasonTriggered {
		t.Errorf("Reason = %q, want %q", got.Reason, enums.MMProtectionReasonTriggered)
	}
}

// ---------- SubscriptionNotification JSON round-trip ----------

func TestSubscriptionNotification_JSONRoundTrip(t *testing.T) {
	sn := types.SubscriptionNotification{
		Channel: "ticker.BTC-PERPETUAL.100ms",
		Data:    "test-data",
	}

	data, err := json.Marshal(sn)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.SubscriptionNotification
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Channel != "ticker.BTC-PERPETUAL.100ms" {
		t.Errorf("Channel = %q, want %q", got.Channel, "ticker.BTC-PERPETUAL.100ms")
	}
}
