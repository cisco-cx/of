package v1

// Represents Mib description.
type Mib struct {
	Name        string
	Description string
	Units       string
	strOID      []string
}

type MibRegistry interface {
	// Return Mib for given OID.
	Mib(string) Mib

	// Return Mib for given OID.
	// Translate each node in OID to its corresponding name, if MIB has its definition, else use the number.
	// Ex : 1.3.6.1.2.1.11.19 -> iso.org.dod.internet.mgmt.mib-2.snmp.snmpInTraps.
	//      1.3.6.1.2.1.11.19.54334 -> iso.org.dod.internet.mgmt.mib-2.snmp.snmpInTraps.54334.
	String(string) string

	// Translate the last node to its name. Ex: 1.3.6.1.2.1.11.19 -> snmpInTraps.
	ShortString(string) string

	// Load given map[oid]Mib into registry.
	Load(map[string]Mib)
}
