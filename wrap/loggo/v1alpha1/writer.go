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

package v1alpha1

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/juju/loggo"
)

// TimeFormat is the time format used for the default writer.
// This can be set with the environment variable LOGGO_TIME_FORMAT.
var TimeFormat = initTimeFormat()

// Colors mapping for log levels
var SeverityColor = map[loggo.Level][]color.Attribute{
	loggo.TRACE:    []color.Attribute{color.FgWhite},
	loggo.DEBUG:    []color.Attribute{color.FgGreen},
	loggo.INFO:     []color.Attribute{color.FgBlue},
	loggo.WARNING:  []color.Attribute{color.FgYellow},
	loggo.ERROR:    []color.Attribute{color.FgRed},
	loggo.CRITICAL: []color.Attribute{color.FgWhite, color.BgRed},
}

// Represents loggo.{riter
type CustomWriter struct {
	writer io.Writer
}

// Init new custom writer.
func NewCustomWriter(w io.Writer) *CustomWriter {
	return &CustomWriter{w}
}

// Fix the format of log output.
// Current format: <timestamp> <log_level> <tag> <file_name and line_number> <message>.
func (cw *CustomWriter) Write(entry loggo.Entry) {
	ts := entry.Timestamp.Format(TimeFormat)
	// Just get the basename from the filename
	filename := filepath.Base(entry.Filename)

	// timestamp
	cw.writeColor([]color.Attribute{color.FgGreen}, "time=\"%s\" ", ts)

	// Log level
	cw.writeColor(SeverityColor[entry.Level], "level=\"%s\" ", entry.Level.Short())

	// Tag/name of logger
	cw.writeColor([]color.Attribute{color.FgWhite}, "tag=\"%s\" ", entry.Module)

	// Filename and line number
	cw.writeColor([]color.Attribute{color.FgMagenta}, "location=\"%s:%d\" ", filename, entry.Line)

	// If the logger was called with WithField(s) or WithError method, entry.Message will have
	// those fields too.
	cw.writeColor([]color.Attribute{color.FgBlue}, "%s\n", entry.Message)
}

// Output text according to given attribute
func (cw *CustomWriter) writeColor(attr []color.Attribute, str string, args ...interface{}) {
	color.Set(attr...)
	defer color.Unset()
	fmt.Fprintf(cw.writer, str, args...)
}

// Override default time format using ENV.
func initTimeFormat() string {
	format := os.Getenv("LOGGO_TIME_FORMAT")
	if format != "" {
		return format
	}
	return "2006-01-02 15:04:05"
}
