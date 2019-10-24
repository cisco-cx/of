package v2

import (
	of_snmp "github.com/cisco-cx/of/pkg/v2/snmp"
)

type configs map[string]bool       // string -> config name. bool to be ignored, Not using array to avoid looping to dedup and search.
type valueMap map[string]configs   // string -> values under select.
type asMap map[of_snmp.As]valueMap // Different operations under of_snmp.As
type lookupMap map[string]asMap    // string -> OID from select.

type Lookup struct {
	lm      lookupMap        // Lookup map
	Configs of_snmp.V2Config // Concatenate list of configs
	V       *Value           // Gets value for given oid, based on the `of_snmp.As` type
}

// Build lookup to match with SNMP Trap var
//
//  Ex: Sample config.
//  epc:                          - key in configs
//    select:
//    - type: equals
//      oid: .1.3.6.1.6.3.1.1.4.1  - key in lookupMap
//      as: value                  - key in asMap
//      values:
//      - .1.3.6.1.4.1.8164.2.13   - key in valueMap
//      - .1.3.6.1.4.1.8164.2.4    - key in valueMap
//    - type: equals
//      oid: .1.3.6.1.6.3.1.1.4.1  - key in lookupMap
//      as: valueStr               - key in asMap
//      values:
//      - name_of_oid              - keys in valueMap
//  nso:                          - key in configs
//    - type: equals
//      oid: .1.3.6.1.6.3.1.1.4.1  - key in lookupMap
//      as: value                  - key in asMap
//      values:
//      - .1.3.6.1.4.1.8164.2.13   - key in valueMap
//
//  The lookup for above example config would be :
//
//  '.1.3.6.1.6.3.1.1.4.1' ->
//  		'value' ->
//    		'.1.3.6.1.4.1.8164.2.13' ->
//    			"epc" -> bool // bool value is irrelevant.
//    			"nso" -> bool // bool value is irrelevant.
//    		'.1.3.6.1.4.1.8164.2.4' ->
//    			"epc" -> bool // bool value is irrelevant.
//  		'valueStr' ->
//    		'name_of_oid' ->
//    			"epc" -> bool // bool value is irrelevant.
//
func (l *Lookup) Build() error {
	l.lm = make(lookupMap)

	// For each config
	for configName, config := range l.Configs {

		// For each Alert in config.Alert
		for _, alert := range config.Alerts {
			for _, selects := range alert.Firing {
				// For each select in config.Alerts.Firing
				l.buildFromSelects(configName, selects)
			}

			for _, selects := range alert.Clearing {
				// For each select in config.Alerts.Firing
				l.buildFromSelects(configName, selects)
			}
		}
	}
	return nil
}

// Build index from given selects.
func (l *Lookup) buildFromSelects(configName string, selects []of_snmp.Select) {
	// For each select in config.Alerts.Firing
	for _, s := range selects {
		// Create map the first time.
		if _, ok := l.lm[s.Oid]; ok == false {
			l.lm[s.Oid] = make(asMap)
		}
		if _, ok := l.lm[s.Oid][s.As]; ok == false {
			l.lm[s.Oid][s.As] = make(valueMap)
		}

		// For each value in config.Alerts.Firing
		for _, value := range s.Values {
			// Create map the first time.
			if _, ok := l.lm[s.Oid][s.As][value]; ok == false {
				l.lm[s.Oid][s.As][value] = make(configs)
			}

			// Add config name to lookup
			l.lm[s.Oid][s.As][value][configName] = false
		}
	}
}

// Lookup configs that are applicable for given oid.
func (l *Lookup) Find(oid string) ([]string, error) {
	var configList = make([]string, 0)
	var configNames = make(configs)
	var as asMap
	var ok bool

	// Checking if `oid` is present in lookupMap
	if as, ok = l.lm[oid]; ok == false {
		return configList, nil
	}

	// iterate through applicable of_snmp.As types for given oid
	// asType : of_snmp.As
	// values : values mentioned under select for given oid.
	for asType, values := range as {

		// Compute interested value of the oid based on of_snmp.As type.
		value, err := l.V.ValueAs(oid, asType)
		if err != nil {
			continue
		}

		// Check if configs are available where given oid is in select and computed value is among the values.
		var cfgs configs
		if cfgs, ok = values[value]; ok == false {
			continue
		}

		// If configs are available add them to the list.
		for cfgName, _ := range cfgs {
			if _, ok := configNames[cfgName]; ok == false {
				configNames[cfgName] = true
				configList = append(configList, cfgName)
			}
		}
	}

	// Return list of configs.
	return configList, nil
}
