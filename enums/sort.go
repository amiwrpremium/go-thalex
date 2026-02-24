package enums

// Sort represents the sort order for paginated results.
type Sort string

const (
	SortAsc  Sort = "asc"
	SortDesc Sort = "desc"
)

// String returns the string representation of the sort order.
func (s Sort) String() string {
	return string(s)
}

// IsValid returns true if the sort order is a recognized value.
func (s Sort) IsValid() bool {
	switch s {
	case SortAsc, SortDesc:
		return true
	}
	return false
}

// SortValues returns all valid Sort values.
func SortValues() []Sort {
	return []Sort{SortAsc, SortDesc}
}
