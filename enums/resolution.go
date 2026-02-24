package enums

// Resolution represents the time resolution for historical data.
type Resolution string

const (
	Resolution1m  Resolution = "1m"
	Resolution5m  Resolution = "5m"
	Resolution15m Resolution = "15m"
	Resolution30m Resolution = "30m"
	Resolution1h  Resolution = "1h"
	Resolution1d  Resolution = "1d"
	Resolution1w  Resolution = "1w"
)

// String returns the string representation of the resolution.
func (r Resolution) String() string {
	return string(r)
}

// IsValid returns true if the resolution is a recognized value.
func (r Resolution) IsValid() bool {
	switch r {
	case Resolution1m, Resolution5m, Resolution15m, Resolution30m, Resolution1h, Resolution1d, Resolution1w:
		return true
	}
	return false
}

// ResolutionValues returns all valid Resolution values.
func ResolutionValues() []Resolution {
	return []Resolution{Resolution1m, Resolution5m, Resolution15m, Resolution30m, Resolution1h, Resolution1d, Resolution1w}
}
