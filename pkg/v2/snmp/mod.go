package snmp

// Represents actions to be performed on snmp.Mod{}
type Modifier interface {
	Apply([]Mod)
	Copy(Mod)
	Set(Mod)
}
