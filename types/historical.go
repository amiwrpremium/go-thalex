package types

import "github.com/amiwrpremium/go-thalex/enums"

// OHLC represents a standard Open-High-Low-Close data point.
type OHLC struct {
	Time  float64 `json:"time"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

// TopOfBook represents best bid/ask at a point in time.
type TopOfBook struct {
	BidPrice *float64 `json:"bid_price"`
	BidSize  *float64 `json:"bid_size"`
	AskPrice *float64 `json:"ask_price"`
	AskSize  *float64 `json:"ask_size"`
}

// PerpetualMarkData represents a mark price data point for perpetuals.
// Format: [time, open, high, low, close, funding_payment, [bid_price, bid_size, ask_price, ask_size]]
type PerpetualMarkData struct {
	OHLC
	FundingPayment float64    `json:"funding_payment"`
	TopOfBook      *TopOfBook `json:"top_of_book,omitempty"`
}

// FutureMarkData represents a mark price data point for futures/combinations.
// Format: [time, open, high, low, close, [bid_price, bid_size, ask_price, ask_size]]
type FutureMarkData struct {
	OHLC
	TopOfBook *TopOfBook `json:"top_of_book,omitempty"`
}

// OptionMarkData represents a mark price data point for options.
// Format: [time, open, high, low, close, iv_open, iv_high, iv_low, iv_close, [bid, bid_size, ask, ask_size]]
type OptionMarkData struct {
	OHLC
	IVOpen    float64    `json:"iv_open"`
	IVHigh    float64    `json:"iv_high"`
	IVLow     float64    `json:"iv_low"`
	IVClose   float64    `json:"iv_close"`
	TopOfBook *TopOfBook `json:"top_of_book,omitempty"`
}

// MarkPriceHistoricalResult contains mark price historical data.
// The raw Mark field contains arrays whose format depends on InstrumentType.
// Use the typed accessor methods to parse the data.
type MarkPriceHistoricalResult struct {
	InstrumentType enums.InstrumentType `json:"instrument_type"`
	Mark           [][]any              `json:"mark"`
	NoData         bool                 `json:"no_data,omitempty"`
}

// PerpetualData parses the Mark field as perpetual data points.
func (r *MarkPriceHistoricalResult) PerpetualData() []PerpetualMarkData {
	out := make([]PerpetualMarkData, 0, len(r.Mark))
	for _, row := range r.Mark {
		if len(row) < 6 {
			continue
		}
		d := PerpetualMarkData{
			OHLC: OHLC{
				Time:  toFloat(row[0]),
				Open:  toFloat(row[1]),
				High:  toFloat(row[2]),
				Low:   toFloat(row[3]),
				Close: toFloat(row[4]),
			},
			FundingPayment: toFloat(row[5]),
		}
		if len(row) >= 7 {
			d.TopOfBook = parseTopOfBook(row[6])
		}
		out = append(out, d)
	}
	return out
}

// FutureData parses the Mark field as future/combination data points.
func (r *MarkPriceHistoricalResult) FutureData() []FutureMarkData {
	out := make([]FutureMarkData, 0, len(r.Mark))
	for _, row := range r.Mark {
		if len(row) < 5 {
			continue
		}
		d := FutureMarkData{
			OHLC: OHLC{
				Time:  toFloat(row[0]),
				Open:  toFloat(row[1]),
				High:  toFloat(row[2]),
				Low:   toFloat(row[3]),
				Close: toFloat(row[4]),
			},
		}
		if len(row) >= 6 {
			d.TopOfBook = parseTopOfBook(row[5])
		}
		out = append(out, d)
	}
	return out
}

// OptionData parses the Mark field as option data points.
func (r *MarkPriceHistoricalResult) OptionData() []OptionMarkData {
	out := make([]OptionMarkData, 0, len(r.Mark))
	for _, row := range r.Mark {
		if len(row) < 9 {
			continue
		}
		d := OptionMarkData{
			OHLC: OHLC{
				Time:  toFloat(row[0]),
				Open:  toFloat(row[1]),
				High:  toFloat(row[2]),
				Low:   toFloat(row[3]),
				Close: toFloat(row[4]),
			},
			IVOpen:  toFloat(row[5]),
			IVHigh:  toFloat(row[6]),
			IVLow:   toFloat(row[7]),
			IVClose: toFloat(row[8]),
		}
		if len(row) >= 10 {
			d.TopOfBook = parseTopOfBook(row[9])
		}
		out = append(out, d)
	}
	return out
}

// IndexPriceHistoricalResult contains index price historical data.
type IndexPriceHistoricalResult struct {
	Index  [][]any `json:"index"`
	NoData bool    `json:"no_data,omitempty"`
}

// Data parses the Index field as OHLC data points.
func (r *IndexPriceHistoricalResult) Data() []OHLC {
	out := make([]OHLC, 0, len(r.Index))
	for _, row := range r.Index {
		if len(row) < 5 {
			continue
		}
		out = append(out, OHLC{
			Time:  toFloat(row[0]),
			Open:  toFloat(row[1]),
			High:  toFloat(row[2]),
			Low:   toFloat(row[3]),
			Close: toFloat(row[4]),
		})
	}
	return out
}

// HistoricalDataParams contains parameters for historical data queries.
type HistoricalDataParams struct {
	InstrumentName string           `json:"instrument_name,omitempty"`
	IndexName      string           `json:"index_name,omitempty"`
	From           float64          `json:"from"`
	To             float64          `json:"to"`
	Resolution     enums.Resolution `json:"resolution"`
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}

func parseTopOfBook(v any) *TopOfBook {
	arr, ok := v.([]any)
	if !ok || len(arr) < 4 {
		return nil
	}
	tob := &TopOfBook{}
	if arr[0] != nil {
		f := toFloat(arr[0])
		tob.BidPrice = &f
	}
	if arr[1] != nil {
		f := toFloat(arr[1])
		tob.BidSize = &f
	}
	if arr[2] != nil {
		f := toFloat(arr[2])
		tob.AskPrice = &f
	}
	if arr[3] != nil {
		f := toFloat(arr[3])
		tob.AskSize = &f
	}
	return tob
}
