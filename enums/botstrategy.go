package enums

// BotStrategy represents the type of bot strategy.
type BotStrategy string

const (
	BotStrategySGSL    BotStrategy = "sgsl"
	BotStrategyOCQ     BotStrategy = "ocq"
	BotStrategyLevels  BotStrategy = "levels"
	BotStrategyGrid    BotStrategy = "grid"
	BotStrategyDHedge  BotStrategy = "dhedge"
	BotStrategyDFollow BotStrategy = "dfollow"
)

// String returns the string representation of the bot strategy.
func (s BotStrategy) String() string {
	return string(s)
}

// IsValid returns true if the bot strategy is a recognized value.
func (s BotStrategy) IsValid() bool {
	switch s {
	case BotStrategySGSL, BotStrategyOCQ, BotStrategyLevels, BotStrategyGrid, BotStrategyDHedge, BotStrategyDFollow:
		return true
	}
	return false
}

// BotStrategyValues returns all valid BotStrategy values.
func BotStrategyValues() []BotStrategy {
	return []BotStrategy{BotStrategySGSL, BotStrategyOCQ, BotStrategyLevels, BotStrategyGrid, BotStrategyDHedge, BotStrategyDFollow}
}
