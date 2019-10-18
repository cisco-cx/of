package snmp

type Event struct {
	EventName string `json:"event_name" yaml:"event_name"`
}

type AlertEntry struct {
	AlertSeverity string           `json:"alert_severity" yaml:"alert_severity"`
	ErrorEvents   map[string]Event `json:"error_events" yaml:"error_events"`
	ClearEvents   map[string]Event `json:"clear_events" yaml:"clear_events"`
}

type AlertsConfig struct {
	Alerts        map[string]*AlertEntry `json:"alerts" yaml:"alerts"`
	DroppedEvents map[string]Event       `json:"dropped_events" yaml:"dropped_events"`
}
