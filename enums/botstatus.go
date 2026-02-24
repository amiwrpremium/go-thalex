package enums

// BotStatus represents the status of a bot.
type BotStatus string

const (
	BotStatusActive  BotStatus = "active"
	BotStatusStopped BotStatus = "stopped"
)

// String returns the string representation of the bot status.
func (s BotStatus) String() string {
	return string(s)
}

// IsValid returns true if the bot status is a recognized value.
func (s BotStatus) IsValid() bool {
	switch s {
	case BotStatusActive, BotStatusStopped:
		return true
	}
	return false
}

// BotStatusValues returns all valid BotStatus values.
func BotStatusValues() []BotStatus {
	return []BotStatus{BotStatusActive, BotStatusStopped}
}

// IsActive returns true if the bot is currently running.
func (s BotStatus) IsActive() bool {
	return s == BotStatusActive
}

// IsFinal returns true if the bot has reached a terminal state.
func (s BotStatus) IsFinal() bool {
	return s == BotStatusStopped
}
