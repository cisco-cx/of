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

package v1

import (
	"fmt"
	"time"
)

// Time is based on time.Time.
//
// FYI: OF's domain type in this case is literally "time.Time".
type Time struct {
	time time.Time
}

// NewTime returns a new instance of Time in UTC timezone.
func NewTime(t time.Time) Time {
	return Time{time: t.UTC()}
}

// MarshalJSON implements the Marshaler interface in encoding/json.
// It leverages the String method of Time.
func (t Time) MarshalJSON() ([]byte, error) {
	// https://stackoverflow.com/a/23695774
	return []byte(fmt.Sprintf("\"%s\"", t)), nil
}

// String implements the fmt.Stringer interface.
//
// String returns Time as a string in UTC timezone and RFC3339 format.
func (t Time) String() string {
	return time.Time(t.time).UTC().Format(time.RFC3339)
}
