package v2

import (
	of "github.com/cisco-cx/of/pkg/v2"
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
)

type Modifier struct {
	Map *map[string]string // Pointer to map to be modified
	V   *Value             // Gets value for given oid, based on the `of_snmp.As` type
}

// Apply given mods to Map.
func (m *Modifier) Apply(mods []of_snmp.Mod) error {
	var err error = nil

	for _, mod := range mods {
		// Redirect to Copy or Set based on type of operation.
		switch mod.Type {
		case of_snmp.Set:
			err = m.Set(mod)
		case of_snmp.Copy:
			err = m.Copy(mod)
		}
		if err != nil {
			return err
		}
	}
	return err
}

// Set given mod to Map.
func (m *Modifier) Set(mod of_snmp.Mod) error {

	// Confirm mode is for Set.
	if mod.Type != of_snmp.Set {
		return of.ErrInvalidOperation
	}

	// Check if key is empty.
	// Not checking if value is empty. There could be a scenario, where user wants to clear the value of this key.
	if mod.Key == "" {
		return of.ErrKeyMissing
	}

	// Set operation.
	(*m.Map)[mod.Key] = mod.Value
	return nil
}

// Copy given mod to Map.
func (m *Modifier) Copy(mod of_snmp.Mod) error {

	// Confirm mode is for Copy.
	if mod.Type != of_snmp.Copy {
		return of.ErrInvalidOperation
	}

	// Check if keys are empty.
	if mod.Oid == "" || mod.ToKey == "" {
		return of.ErrKeyMissing
	}

	// Get value for key based on of_snmp.As type.
	value, err := m.V.ValueAs(mod.Oid, mod.As)
	if err != nil {
		// Bubble up error if drop on error is set.
		if mod.OnError == of_snmp.Drop {
			return err
		}
		return nil
	}

	// If map is not present copy value to map
	if len(mod.Map) == 0 {
		(*m.Map)[mod.ToKey] = value
		return nil
	}

	// Check if value is a key in map,
	if v, ok := mod.Map[value]; ok == true {
		(*m.Map)[mod.ToKey] = v
		return nil
	}

	return nil
}
