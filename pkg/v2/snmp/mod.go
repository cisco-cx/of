package snmp

// Represents actions to be performed on snmp.Mod{}
type Modifier interface {
	Apply([]Mod) error
	Copy(Mod) error
	Set(Mod) error
}
