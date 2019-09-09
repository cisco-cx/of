// Copyright 2019 Cisco Systems, Inc.
//
// This work incorporates works covered by the following notice:
//
// Copyright 2018 The Prometheus Authors
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

package v1alpha1_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/juju/loggo"
	"github.com/stretchr/testify/require"
	"github.com/cisco-cx/of/wrap/loggo/v1alpha1"
	log "github.com/cisco-cx/of/wrap/loggo/v1alpha1"
)

// Enforce interface implementation.
func TestInterface(t *testing.T) {
	var _ v1alpha1.Logger = log.Logger{}
}

// Test Error without any fields
func TestLogError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := customLogger("test.errorf", buf)
	err := errors.New("This is an error message.")
	logger.Errorf("%s", err.Error())
	// time="2019-09-09 17:21:02" level="ERROR" tag="test.errorf" location="log_test.go:43" msg="This is an error message."
	require.Contains(t, string(buf.Bytes()), fmt.Sprintf("level=\"%s\" tag=\"%s\" location=\"%s\" msg=\"%s\"", "ERROR", "test.errorf", "log_test.go:43", "This is an error message."))
}

// Test logger with WithError method.
func TestWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := customLogger("test.with.error", buf)
	err := errors.New("This is a custom error.")
	logger.WithError(err).Errorf("Encountered an error.")

	// time="2019-09-09 17:48:40" level="ERROR" tag="test.with.error" location="log_test.go:53" msg="Encountered an error."	error="This is a custom error."
	require.Contains(t, string(buf.Bytes()), fmt.Sprintf("level=\"%s\" tag=\"%s\" location=\"%s\" msg=\"%s\"\terror=\"%s\"", "ERROR", "test.with.error", "log_test.go:53", "Encountered an error.", "This is a custom error."))
}

// Test logger with WithField method.
func TestWithField(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := customLogger("test.with.field", buf)
	logger.WithField("key", "value").WithField("key2", "value2").Criticalf("Errors with custom field.")

	// time="2019-09-09 18:38:23" level="CRITC" tag="test.with.field" location="log_test.go:63" msg="Errors with custom field."		key="value" key2="value2"
	require.Contains(t, string(buf.Bytes()), fmt.Sprintf("level=\"%s\" tag=\"%s\" location=\"%s\" msg=\"%s\"\t\tkey=\"value\" key2=\"value2\"", "CRITC", "test.with.field", "log_test.go:63", "Errors with custom field."))
}

// Test logger with WithFields method.
func TestWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := customLogger("test.with.fields", buf)
	logger.SetLogLevel("critical")
	logger.WithFields(map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}).
		WithField("key3", "val3").
		Criticalf("Errors with custom fields.")

	// time="2019-09-09 18:38:23" level="CRITC" tag="test.with.field" location="log_test.go:79" msg="Errors with custom field."		key1="val1" key2="val2" key3="val3"
	require.Contains(t, string(buf.Bytes()), fmt.Sprintf("level=\"%s\" tag=\"%s\" location=\"%s\" msg=\"%s\"\t\tkey1=\"val1\" key2=\"val2\" key3=\"val3\"", "CRITC", "test.with.fields", "log_test.go:79", "Errors with custom fields."))
}

// Test debug log level enabled
func TestWithDebugEnabled(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := customLogger("test.debug.enabled", buf)
	logger.SetLogLevel("deBug")
	logger.Debugf("Debug log enabled.")

	// "time="2019-09-09 19:12:11" level="DEBUG" tag="test.with.fields" location="log_test.go:90" msg="Debug log enabled."
	require.Contains(t, string(buf.Bytes()), fmt.Sprintf("level=\"%s\" tag=\"%s\" location=\"%s\" msg=\"%s\"", "DEBUG", "test.debug.enabled", "log_test.go:90", "Debug log enabled."))
}

// Test debug log level disabled
func TestWithDebugDisabled(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := customLogger("test.debug.enabled", buf)
	logger.SetLogLevel("deBug")
	logger.Debugf("Debug log enabled.")

	// "time="2019-09-09 19:12:11" level="DEBUG" tag="test.with.fields" location="log_test.go:90" msg="Debug log enabled."
	require.Contains(t, string(buf.Bytes()), "")
}

// Custom logger that writes to a buffer for testing, instead of os.Stderr
func customLogger(name string, output io.Writer) *log.Logger {

	logger := log.New(name)
	w := log.NewCustomWriter(output)
	_, err := loggo.ReplaceDefaultWriter(w)
	if err != nil {
		fmt.Errorf("Failed to replace default writer, %s", err.Error())
	}
	return logger
}
