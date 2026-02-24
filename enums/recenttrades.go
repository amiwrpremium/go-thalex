package enums

// RecentTradesCategory categorizes recent trade subscriptions.
type RecentTradesCategory string

const (
	RecentTradesCategoryAll    RecentTradesCategory = "all"
	RecentTradesCategoryNormal RecentTradesCategory = "normal"
	RecentTradesCategoryBlock  RecentTradesCategory = "block"
	RecentTradesCategoryCombo  RecentTradesCategory = "combo"
)

// String returns the string representation of the recent trades category.
func (c RecentTradesCategory) String() string {
	return string(c)
}

// IsValid returns true if the recent trades category is a recognized value.
func (c RecentTradesCategory) IsValid() bool {
	switch c {
	case RecentTradesCategoryAll, RecentTradesCategoryNormal, RecentTradesCategoryBlock, RecentTradesCategoryCombo:
		return true
	}
	return false
}

// RecentTradesCategoryValues returns all valid RecentTradesCategory values.
func RecentTradesCategoryValues() []RecentTradesCategory {
	return []RecentTradesCategory{RecentTradesCategoryAll, RecentTradesCategoryNormal, RecentTradesCategoryBlock, RecentTradesCategoryCombo}
}
