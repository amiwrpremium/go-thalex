package types_test

import (
	"encoding/json"
	"testing"

	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
)

// ---------- NewDoubleSidedQuote ----------

func TestNewDoubleSidedQuote(t *testing.T) {
	t.Run("with_bids_and_asks", func(t *testing.T) {
		bids := []types.QuoteLevel{
			{Price: 49000.0, Amount: 1.0},
			{Price: 48900.0, Amount: 2.0},
		}
		asks := []types.QuoteLevel{
			{Price: 51000.0, Amount: 1.5},
		}
		q := types.NewDoubleSidedQuote("BTC-PERPETUAL", bids, asks)
		if q.I != "BTC-PERPETUAL" {
			t.Errorf("I = %q, want %q", q.I, "BTC-PERPETUAL")
		}
		if q.B == nil {
			t.Fatal("B should not be nil when bids are provided")
		}
		if q.A == nil {
			t.Fatal("A should not be nil when asks are provided")
		}

		// Verify the bid array is [][2]float64
		bidArr, ok := q.B.([][2]float64)
		if !ok {
			t.Fatalf("B type = %T, want [][2]float64", q.B)
		}
		if len(bidArr) != 2 {
			t.Fatalf("len(B) = %d, want 2", len(bidArr))
		}
		if bidArr[0][0] != 49000.0 || bidArr[0][1] != 1.0 {
			t.Errorf("B[0] = %v, want [49000 1]", bidArr[0])
		}
		if bidArr[1][0] != 48900.0 || bidArr[1][1] != 2.0 {
			t.Errorf("B[1] = %v, want [48900 2]", bidArr[1])
		}

		askArr, ok := q.A.([][2]float64)
		if !ok {
			t.Fatalf("A type = %T, want [][2]float64", q.A)
		}
		if len(askArr) != 1 {
			t.Fatalf("len(A) = %d, want 1", len(askArr))
		}
		if askArr[0][0] != 51000.0 || askArr[0][1] != 1.5 {
			t.Errorf("A[0] = %v, want [51000 1.5]", askArr[0])
		}
	})

	t.Run("empty_bids", func(t *testing.T) {
		asks := []types.QuoteLevel{{Price: 51000.0, Amount: 1.0}}
		q := types.NewDoubleSidedQuote("ETH-PERPETUAL", nil, asks)
		if q.B != nil {
			t.Errorf("B = %v, want nil when no bids", q.B)
		}
		if q.A == nil {
			t.Error("A should not be nil when asks are provided")
		}
	})

	t.Run("empty_asks", func(t *testing.T) {
		bids := []types.QuoteLevel{{Price: 49000.0, Amount: 1.0}}
		q := types.NewDoubleSidedQuote("ETH-PERPETUAL", bids, nil)
		if q.A != nil {
			t.Errorf("A = %v, want nil when no asks", q.A)
		}
		if q.B == nil {
			t.Error("B should not be nil when bids are provided")
		}
	})

	t.Run("both_empty", func(t *testing.T) {
		q := types.NewDoubleSidedQuote("ETH-PERPETUAL", nil, nil)
		if q.B != nil {
			t.Errorf("B = %v, want nil", q.B)
		}
		if q.A != nil {
			t.Errorf("A = %v, want nil", q.A)
		}
	})

	t.Run("empty_slices", func(t *testing.T) {
		q := types.NewDoubleSidedQuote("ETH-PERPETUAL", []types.QuoteLevel{}, []types.QuoteLevel{})
		if q.B != nil {
			t.Errorf("B = %v, want nil for empty slice", q.B)
		}
		if q.A != nil {
			t.Errorf("A = %v, want nil for empty slice", q.A)
		}
	})
}

// ---------- NewSingleLevelQuote ----------

func TestNewSingleLevelQuote(t *testing.T) {
	q := types.NewSingleLevelQuote("BTC-PERPETUAL", 49000.0, 1.0, 51000.0, 2.0)
	if q.I != "BTC-PERPETUAL" {
		t.Errorf("I = %q, want %q", q.I, "BTC-PERPETUAL")
	}

	bid, ok := q.B.(types.SingleLevelQuote)
	if !ok {
		t.Fatalf("B type = %T, want SingleLevelQuote", q.B)
	}
	if bid.P != 49000.0 || bid.A != 1.0 {
		t.Errorf("B = %+v, want P=49000 A=1", bid)
	}

	ask, ok := q.A.(types.SingleLevelQuote)
	if !ok {
		t.Fatalf("A type = %T, want SingleLevelQuote", q.A)
	}
	if ask.P != 51000.0 || ask.A != 2.0 {
		t.Errorf("A = %+v, want P=51000 A=2", ask)
	}
}

// ---------- DoubleSidedQuote JSON ----------

func TestDoubleSidedQuote_JSONRoundTrip(t *testing.T) {
	q := types.NewSingleLevelQuote("BTC-PERPETUAL", 49000.0, 1.0, 51000.0, 2.0)

	data, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify the JSON contains expected fields
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}
	if raw["i"] != "BTC-PERPETUAL" {
		t.Errorf("JSON i = %v, want BTC-PERPETUAL", raw["i"])
	}
}

// ---------- DoubleSidedQuoteResult JSON ----------

func TestDoubleSidedQuoteResult_JSONRoundTrip(t *testing.T) {
	price := 50000.0
	r := types.DoubleSidedQuoteResult{
		NSuccess: 5,
		NFail:    1,
		Errors: []types.QuoteError{
			{Code: 100, Message: "price_too_far", Side: "buy", Price: &price},
		},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.DoubleSidedQuoteResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.NSuccess != 5 {
		t.Errorf("NSuccess = %d, want 5", got.NSuccess)
	}
	if got.NFail != 1 {
		t.Errorf("NFail = %d, want 1", got.NFail)
	}
	if len(got.Errors) != 1 {
		t.Fatalf("len(Errors) = %d, want 1", len(got.Errors))
	}
	if got.Errors[0].Code != 100 {
		t.Errorf("Errors[0].Code = %d, want 100", got.Errors[0].Code)
	}
	if got.Errors[0].Side != "buy" {
		t.Errorf("Errors[0].Side = %q, want %q", got.Errors[0].Side, "buy")
	}
	if got.Errors[0].Price == nil || *got.Errors[0].Price != price {
		t.Errorf("Errors[0].Price = %v, want %v", got.Errors[0].Price, price)
	}
}

func TestQuoteError_NilPrice(t *testing.T) {
	qe := types.QuoteError{Code: 200, Message: "error"}
	data, err := json.Marshal(qe)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}
	if _, ok := raw["price"]; ok {
		t.Error("expected price to be omitted when nil")
	}
}

// ---------- MassQuoteParams ----------

func TestNewMassQuoteParams(t *testing.T) {
	quotes := []types.DoubleSidedQuote{
		types.NewSingleLevelQuote("BTC-PERPETUAL", 49000, 1, 51000, 1),
	}
	p := types.NewMassQuoteParams(quotes)
	if p == nil {
		t.Fatal("NewMassQuoteParams returned nil")
	}
	if len(p.Quotes) != 1 {
		t.Errorf("len(Quotes) = %d, want 1", len(p.Quotes))
	}
	if p.Label != "" {
		t.Errorf("Label = %q, want empty", p.Label)
	}
	if p.PostOnly != nil {
		t.Errorf("PostOnly = %v, want nil", p.PostOnly)
	}
}

func TestMassQuoteParams_WithLabel(t *testing.T) {
	p := types.NewMassQuoteParams(nil)
	ret := p.WithLabel("mm-label")
	if ret != p {
		t.Error("WithLabel should return the same pointer for chaining")
	}
	if p.Label != "mm-label" {
		t.Errorf("Label = %q, want %q", p.Label, "mm-label")
	}
}

func TestMassQuoteParams_WithPostOnly(t *testing.T) {
	p := types.NewMassQuoteParams(nil)
	ret := p.WithPostOnly(true)
	if ret != p {
		t.Error("WithPostOnly should return the same pointer for chaining")
	}
	if p.PostOnly == nil || *p.PostOnly != true {
		t.Errorf("PostOnly = %v, want true", p.PostOnly)
	}
}

func TestMassQuoteParams_WithRejectPostOnly(t *testing.T) {
	p := types.NewMassQuoteParams(nil)
	ret := p.WithRejectPostOnly(true)
	if ret != p {
		t.Error("WithRejectPostOnly should return the same pointer for chaining")
	}
	if p.RejectPostOnly == nil || *p.RejectPostOnly != true {
		t.Errorf("RejectPostOnly = %v, want true", p.RejectPostOnly)
	}
}

func TestMassQuoteParams_WithSTP(t *testing.T) {
	p := types.NewMassQuoteParams(nil)
	ret := p.WithSTP(enums.STPLevelAccount, enums.STPActionCancelBoth)
	if ret != p {
		t.Error("WithSTP should return the same pointer for chaining")
	}
	if p.STPLevel != enums.STPLevelAccount {
		t.Errorf("STPLevel = %q, want %q", p.STPLevel, enums.STPLevelAccount)
	}
	if p.STPAction != enums.STPActionCancelBoth {
		t.Errorf("STPAction = %q, want %q", p.STPAction, enums.STPActionCancelBoth)
	}
}

func TestMassQuoteParams_Chaining(t *testing.T) {
	p := types.NewMassQuoteParams(nil).
		WithLabel("chain").
		WithPostOnly(true).
		WithRejectPostOnly(false).
		WithSTP(enums.STPLevelCustomer, enums.STPActionCancelPassive)

	if p.Label != "chain" {
		t.Errorf("Label = %q, want %q", p.Label, "chain")
	}
	if p.PostOnly == nil || *p.PostOnly != true {
		t.Errorf("PostOnly = %v, want true", p.PostOnly)
	}
	if p.RejectPostOnly == nil || *p.RejectPostOnly != false {
		t.Errorf("RejectPostOnly = %v, want false", p.RejectPostOnly)
	}
	if p.STPLevel != enums.STPLevelCustomer {
		t.Errorf("STPLevel = %q, want %q", p.STPLevel, enums.STPLevelCustomer)
	}
}

func TestMassQuoteParams_JSONRoundTrip(t *testing.T) {
	p := types.NewMassQuoteParams([]types.DoubleSidedQuote{
		types.NewSingleLevelQuote("BTC-PERPETUAL", 49000, 1, 51000, 1),
	}).WithLabel("test").WithPostOnly(true)

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	if raw["label"] != "test" {
		t.Errorf("label = %v, want test", raw["label"])
	}
}

// ---------- MMProtectionParams JSON ----------

func TestMMProtectionParams_JSONRoundTrip(t *testing.T) {
	p := types.MMProtectionParams{
		Product:     enums.Product("FBTCUSD"),
		TradeAmount: 100.0,
		QuoteAmount: 200.0,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.MMProtectionParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.Product != enums.Product("FBTCUSD") {
		t.Errorf("Product = %q, want %q", got.Product, "FBTCUSD")
	}
	if got.TradeAmount != 100.0 {
		t.Errorf("TradeAmount = %v, want 100", got.TradeAmount)
	}
	if got.QuoteAmount != 200.0 {
		t.Errorf("QuoteAmount = %v, want 200", got.QuoteAmount)
	}
}

// ---------- SingleLevelQuote JSON ----------

func TestSingleLevelQuote_JSONRoundTrip(t *testing.T) {
	q := types.SingleLevelQuote{P: 50000.0, A: 1.5}

	data, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var got types.SingleLevelQuote
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if got.P != 50000.0 {
		t.Errorf("P = %v, want 50000", got.P)
	}
	if got.A != 1.5 {
		t.Errorf("A = %v, want 1.5", got.A)
	}
}
