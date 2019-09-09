// Copyright 2019 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import "time"

// ACIFaultRaw represents an ACI fault as scraped from the ACI API.
// It is designed to be used in the receiver of ACIFaultRawParser.
type ACIFaultRaw struct {
	Ack             string `json:"ack,omitempty"`
	Cause           string `json:"cause,omitempty"`
	ChangeSet       string `json:"changeSet,omitempty"`
	ChildAction     string `json:"childAction,omitempty"`
	Code            string `json:"code,omitempty"`
	Created         string `json:"created,omitempty"`
	DN              string `json:"dn,omitempty"`
	Delegated       string `json:"delegated,omitempty"`
	Desc            string `json:"desc,omitempty"`
	Domain          string `json:"domain,omitempty"`
	HighestSeverity string `json:"highestSeverity,omitempty"`
	LC              string `json:"lc,omitempty"`
	LastTransition  string `json:"lastTransition,omitempty"`
	Occur           string `json:"occur,omitempty"`
	OrigSeverity    string `json:"origSeverity,omitempty"`
	PrevSeverity    string `json:"prevSeverity,omitempty"`
	Rule            string `json:"rule,omitempty"`
	Severity        string `json:"severity,omitempty"`
	Status          string `json:"status,omitempty"`
	Subject         string `json:"subject,omitempty"`
	Type            string `json:"type,omitempty"`
}

// ACIFaultRawParser represents the ability to parse an ACI fault as scraped
// from the ACI API.
type ACIFaultRawParser interface {

	// Created returns the ACI API's created time for the fault in UTC timezone
	// and as RFC3339 time format.
	Created() (time.Time, error)

	// LastTransition returns the ACI API's last transition time for the fault
	// in UTC timezone and as RFC3339 time format.
	LastTransition() (time.Time, error)

	// ServerityID returns a numerical severity for the fault based on
	// the return value from ACIFaultRawSeverityIDParser.
	SeverityID() (ACIFaultSeverityID, error)

	// SubID returns the fault's `sub_id`. The return value is result of
	// pruning pattern `/fault-.*` from the fault's Distinguished Name (or DN).
	SubID() (string, error)
}
