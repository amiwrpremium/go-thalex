package types

import "github.com/amiwrpremium/go-thalex/enums"

// SystemInfo represents system status information.
type SystemInfo struct {
	Environment string   `json:"environment"`
	APIVersion  string   `json:"api_version,omitempty"`
	Banners     []Banner `json:"banners"`
}

// Banner represents a system banner/announcement.
type Banner struct {
	ID       *int           `json:"id,omitempty"`
	Time     float64        `json:"time"`
	Severity enums.Severity `json:"severity"`
	Title    string         `json:"title,omitempty"`
	Message  string         `json:"message"`
}

// SystemEvent represents a system event notification.
type SystemEvent struct {
	Event enums.SystemEventType `json:"event"`
}
