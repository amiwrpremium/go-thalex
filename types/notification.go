package types

import "github.com/amiwrpremium/go-thalex/enums"

// Notification represents a user notification.
type Notification struct {
	ID            string            `json:"id"`
	Time          float64           `json:"time"`
	Category      string            `json:"category"`
	Title         string            `json:"title"`
	Message       string            `json:"message"`
	DisplayType   enums.DisplayType `json:"display_type"`
	Read          bool              `json:"read"`
	AccountName   string            `json:"account_name,omitempty"`
	AccountNumber string            `json:"account_number,omitempty"`
	Popup         bool              `json:"popup"`
}

// NotificationsResult wraps a list of notifications.
type NotificationsResult struct {
	Notifications []Notification `json:"notifications"`
}

// NotificationPreferences maps category names to preference settings.
type NotificationPreferences map[string]NotificationPreference

// NotificationPreference holds preferences for a single notification category.
type NotificationPreference struct {
	Email bool `json:"email,omitempty"`
	Inbox bool `json:"inbox"`
	Popup bool `json:"popup,omitempty"`
	SMS   bool `json:"sms,omitempty"`
}
