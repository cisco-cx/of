package snmp

type SourceType string
type ModType string
type SelectType string
type As string
type OnError string
type EventType string

type Enabled *bool // Need *bool to differentiate between Enabled being set to false and not being defined in config.
type URLPrefix string

const (
	// SourceType constants
	HostType    SourceType = "host"
	ClusterType SourceType = "cluster"

	// ModType constants
	Copy ModType = "copy"
	Set  ModType = "set"

	// SelectType constants
	Equals SelectType = "equals"

	// As constants
	Value         As = "value"
	ValueStr      As = "value-str"
	ValueStrShort As = "value-str-short"

	OidValue         As = "oid.value"
	OidValueStr      As = "oid.value-str"
	OidValueStrShort As = "oid.value-str-short"

	// OnError constants
	Send OnError = "send"
	Drop OnError = "drop"

	// Alert related constants
	Firing          EventType = "error"
	Clearing        EventType = "clear"
	EventTypeText   string    = "event_type"
	FingerprintText string    = "alert_fingerprint"
	SNMPTrapOID     string    = ".1.3.6.1.6.3.1.1.4.1.0"
)

// Represents map of configs from different files in conf.d
// key in the map is the top-level name of the respective configs.
type V2Config map[string]Config

// Represents version v2 of SNMP config
type Config struct {
	Defaults Default `yaml:"defaults,omitempty"`
	Alerts   []Alert `yaml:"alerts,omitempty"`
}

// Represents Default attributes to be addded to Labels and Annotations.
type Default struct {
	Enabled            Enabled            `yaml:"enabled,omitempty"`
	SourceType         SourceType         `yaml:"source_type,omitempty"`
	DeviceIdentifiers  []string           `yaml:"device_identifiers,omitempty"`
	Clusters           map[string]Cluster `yaml:"clusters,omitempty"`
	GeneratorUrlPrefix URLPrefix          `yaml:"generator_url_prefix,omitempty"`
	LabelMods          []Mod              `yaml:"label_mods,omitempty"`
	AnnotationMods     []Mod              `yaml:"annotation_mods,omitempty"`
	EndsAt             int                `yaml:"ends_at,omitempty"`
}

// Maps IPaddresses to there cluster name.
type Cluster struct {
	SourceAddresses []string `yaml:"source_addresses,omitempty"`
}

// Represents a mod operation to be performed on labels and annotations.
type Mod struct {
	Type ModType `yaml:"type,omitempty"`

	// Set specific keys
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`

	// Copy specific keys
	Oid     string            `yaml:"oid,omitempty"` // OID can be replaced with Key while removing SNMP specific details from v2 config.
	As      As                `yaml:"as,omitempty"`
	ToKey   string            `yaml:"to_key,omitempty"`
	OnError OnError           `yaml:"on_error,omitempty"`
	Map     map[string]string `yaml:"map,omitempty"`
}

// Represents an alert group under v2 config
type Alert struct {
	Name               string              `yaml:"name,omitempty"`
	Enabled            Enabled             `yaml:"enabled,omitempty"`
	GeneratorUrlPrefix URLPrefix           `yaml:"generator_url_prefix,omitempty"`
	LabelMods          []Mod               `yaml:"label_mods,omitempty"`
	AnnotationMods     []Mod               `yaml:"annotation_mods,omitempty"`
	Firing             map[string][]Select `yaml:"firing,omitempty"`
	Clearing           map[string][]Select `yaml:"clearing,omitempty"`
	EndsAt             int                 `yaml:"ends_at,omitempty"`
}

// Represents the alert selection criteria.
type Select struct {
	Type           SelectType `yaml:"type,omitempty"`
	Oid            string     `yaml:"oid,omitempty"` // OID can be replaced with Key while removing SNMP specific details from v2 config.
	As             As         `yaml:"as,omitempty"`
	Values         []string   `yaml:"values,omitempty"`
	AnnotationMods []Mod      `yaml:"annotation_mods,omitempty"`
}
