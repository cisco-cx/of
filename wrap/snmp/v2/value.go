package v2

import (
	"strconv"
	"strings"

	v1 "github.com/cisco-cx/of/pkg/v1"
	v2 "github.com/cisco-cx/of/pkg/v2"
	snmp "github.com/cisco-cx/of/pkg/v2/snmp"
	mib_registry "github.com/cisco-cx/of/wrap/mib/v1"
)

// Implements snmp.ValueGenerator
type Value struct {
	vars map[string]string
	mr   *mib_registry.MibRegistry
}

// Initialize Value. trapVars are converted into a map[oid]value
func NewValue(trapVars *[]v1.TrapVar, mr *mib_registry.MibRegistry) *Value {
	vars := make(map[string]string)
	for _, v := range *trapVars {
		vars[v.Oid] = v.Value
	}
	return &Value{vars: vars, mr: mr}
}

// Compute value as `As` for given OID.
func (v *Value) ValueAs(oid string, as snmp.As) (string, error) {

	var val string = ""
	var err error = nil
	// As constants
	switch as {
	case snmp.Value:
		val, err = v.Value(oid)
	case snmp.ValueStr:
		val, err = v.ValueStr(oid)
	case snmp.ValueStrShort:
		val, err = v.ValueStrShort(oid)
	case snmp.OidValue:
		val, err = v.OIDValue(oid)
	case snmp.OidValueStr:
		val, err = v.OIDValueStr(oid)
	case snmp.OidValueStrShort:
		val, err = v.OIDValueStrShort(oid)
	default:
		err = v2.ErrUnknownAs
	}
	return val, err
}

// Literal value for given OID.
func (v *Value) Value(oid string) (string, error) {
	var val string
	var ok bool
	if val, ok = v.vars[oid]; ok == false {
		return val, v2.ErrOIDNotFound
	}
	return val, nil
}

// String representation of the value, for given OID.
func (v *Value) ValueStr(oid string) (string, error) {
	val, err := v.numOid(oid)
	if err != nil {
		return val, err
	}

	return v.mr.String(val), nil
}

// Short Name of the value, for given OID.
func (v *Value) ValueStrShort(oid string) (string, error) {
	val, err := v.numOid(oid)
	if err != nil {
		return val, err
	}

	return v.mr.ShortString(val), nil
}

// Literal value for OID pointed by given OID,
func (v *Value) OIDValue(ptr string) (string, error) {
	oid, err := v.Value(ptr)
	if err != nil {
		return oid, err
	}
	return v.Value(oid)
}

// String representation of the value, for OID pointed by given OID.
func (v *Value) OIDValueStr(ptr string) (string, error) {
	oid, err := v.numOid(ptr)
	if err != nil {
		return oid, err
	}
	return v.ValueStr(oid)
}

// Short Name of the value, for OID pointed by given OID.
func (v *Value) OIDValueStrShort(ptr string) (string, error) {
	oid, err := v.numOid(ptr)
	if err != nil {
		return oid, err
	}
	return v.ValueStrShort(oid)
}

// Validate value for given OID is an OID.
func (v *Value) numOid(oid string) (string, error) {
	val, err := v.Value(oid)
	if err != nil {
		return val, err
	}
	nodes := strings.Split(val, ".")[1:]
	if len(nodes) <= 1 {
		return val, v2.ErrNoneNumericalOID
	}
	for _, n := range nodes {
		_, err = strconv.ParseInt(n, 10, 0)
		if err != nil {
			return val, v2.ErrNoneNumericalOID
		}
	}
	return val, nil
}
