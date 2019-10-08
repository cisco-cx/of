package v1

import (
	"fmt"
	"strings"
	"time"

	of "github.com/cisco-cx/of/lib/v1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
)

type FaultParser struct {
	Fault of.ACIFaultRaw
	Log   *logger.Logger
}

const timeLayout = time.RFC3339

// Created returns the ACI API's created time for the fault in UTC timezone
// and as RFC3339 time format.
func (f *FaultParser) Created() (time.Time, error) {
	t, err := time.Parse(timeLayout, f.Fault.Created)
	if err != nil {
		err := fmt.Errorf("APIC fault's 'created' field contains an unrecognized time string: %s\n", err)
		return time.Time{}, err
	}

	zone, offset := t.Zone()
	f.Log.Debugf("Before forcing UTC time zone, APIC fault's 'created' timezone and offset were: %s, %d\n", zone, offset)
	return t.UTC(), nil
}

// LastTransition returns the ACI API's last transition time for the fault
// in UTC timezone and as RFC3339 time format.
func (f *FaultParser) LastTransition() (time.Time, error) {
	t, err := time.Parse(timeLayout, f.Fault.LastTransition)
	if err != nil {
		err := fmt.Errorf("APIC fault's 'lastTransition' field contains an unrecognized time string: %s\n", err)
		return time.Time{}, err
	}
	zone, offset := t.Zone()
	f.Log.Debugf("Before forcing UTC time zone, APIC fault's 'lastTransition' timezone and offset were: %s, %d\n", zone, offset)
	return t.UTC(), nil
}

// SubID returns the fault's `sub_id`. The return value is result of
// pruning pattern `/fault-.*` from the fault's Distinguished Name (or DN).
func (f *FaultParser) SubID() (string, error) {
	// subID() returns f.DN but without "/fault-.*"
	s := strings.Split(f.Fault.DN, "/fault-")[0]
	return s, nil
}

// ServerityID returns a numerical severity for the fault based on
// the return value from ACIFaultRawSeverityIDParser.
func (f *FaultParser) SeverityID() (of.ACIFaultSeverityID, error) {
	r, err := NewACIFaultSeverityRaw(f.Fault.Severity)
	if err != nil {
		return -1, err
	}
	return r.ID(), nil
}
