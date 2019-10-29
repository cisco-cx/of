package v2

// Represents of.UUIDGen
type FixedUUID struct {
}

// UUID generated to be used in test cases, always returns a fixed string.
func (u *FixedUUID) UUID() string {
	return "9dcc77fc-dda5-4edf-a683-64f2589036d6"
}
