package v1_test

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/lib/v1"
	acigo "github.com/cisco-cx/of/wrap/acigo/v1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
)

var fault = of.ACIFaultRaw{
	Ack:         "no",
	Cause:       "port-failure",
	ChangeSet:   "usage (New: epg)",
	ChildAction: "",
	Code:        "F1678",
	Created:     "2015-01-19T14:26:13.862+01:00",
	Desc: "TEST FAULT -- Port is down, reason:sfpAbsent(connected), " +
		"used by:EPG, lastLinkStChg:1970-01-01T01:00:00.000+01:00, operSt:down",
	DN:              "topology/pod-1/node-101/sys/phys-[eth1/25]/phys/fault-F1678",
	Domain:          "access",
	HighestSeverity: "critical",
	LastTransition:  "2015-01-19T14:28:41.668+01:00",
	LC:              "raised",
	Occur:           "1",
	OrigSeverity:    "critical",
	PrevSeverity:    "critical",
	//RN:              "fault-F1678",
	Rule:     "ethpm-if-port-down-infra-epg-test",
	Severity: "critical",
	Status:   "",
	Subject:  "port-down",
	Type:     "communications",
}

// Enforce interface implementation.
func TestInterface(t *testing.T) {
	var _ of.ACIFaultRawParser = &acigo.FaultParser{}
}

// Test parsing created time
func TestFaultCreated(t *testing.T) {
	fp := faultParser()
	createTime, err := fp.Created()
	require.NoError(t, err)
	require.EqualValues(t, "2015-01-19 13:26:13.862 +0000 UTC", createTime.String())
}

// Test parsing last transition time
func TestFaultLastTransition(t *testing.T) {
	fp := faultParser()
	lt, err := fp.LastTransition()
	require.NoError(t, err)
	require.EqualValues(t, "2015-01-19 13:28:41.668 +0000 UTC", lt.String())
}

// Test parsing fault's sub ID
func TestSubID(t *testing.T) {
	fp := faultParser()
	id, err := fp.SubID()
	require.NoError(t, err)
	fmt.Println(id)
	require.EqualValues(t, "topology/pod-1/node-101/sys/phys-[eth1/25]/phys", id)
}

// Test parsing fault's severity ID
func TestSeverityID(t *testing.T) {
	fp := faultParser()
	id, err := fp.SeverityID()
	require.NoError(t, err)
	fmt.Println(id)
	require.EqualValues(t, "5", fmt.Sprintf("%d", id))
}

// Wrapper to create a acigo.FaultParser with test fault data.
func faultParser() *acigo.FaultParser {
	log := logger.New()
	f := of.ACIFaultRaw{}
	mapstructure.Decode(fault, &f)
	return &acigo.FaultParser{f, log}
}
