package changelog

// ChangeType is a type of changes
type ChangeType string

const (
	ChangeTypeAdded      = "Added"
	ChangeTypeChanged    = "Changed"
	ChangeTypeDeprecated = "Deprecated"
	ChangeTypeRemoved    = "Removed"
	ChangeTypeFixed      = "Fixed"
	ChangeTypeSecurity   = "Security"
	ChangeTypeUnknown    = "Unknown"
)

// ChangeTypeFromString creates a type based on its string name.
func ChangeTypeFromString(s string) ChangeType {
	switch s {
	case ChangeTypeAdded:
		return ChangeTypeAdded
	case ChangeTypeChanged:
		return ChangeTypeChanged
	case ChangeTypeDeprecated:
		return ChangeTypeDeprecated
	case ChangeTypeRemoved:
		return ChangeTypeRemoved
	case ChangeTypeFixed:
		return ChangeTypeFixed
	case ChangeTypeSecurity:
		return ChangeTypeSecurity
	default:
		return ChangeTypeUnknown
	}
}
