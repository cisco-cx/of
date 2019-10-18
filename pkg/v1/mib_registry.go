package v1

type Mib struct {
	Name        string
	Description string
	Units       string
	strOID      []string
}

type MibRegistry interface {
	Mib(string) Mib       // Return Mib for given OID.
	String(string) string // Return display string for given OID.
	Load(map[string]Mib)  // Load given map[oid]Mib into registry.
}
