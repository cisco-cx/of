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
	"os"
	"sort"

	"github.com/juju/loggo"
)

// Represents loggo logger and fields to support structured logging.
type Logger struct {
	logger     loggo.Logger
	fields     map[string]interface{}
	fieldsText string
	errorText  string
}

// Initiate logger.
func New(name string) *Logger {
	logger := loggo.GetLogger(name)
	f := make(map[string]interface{})
	l := Logger{logger: logger, fields: f}

	w := NewCustomWriter(os.Stderr)
	_, err := loggo.ReplaceDefaultWriter(w)
	if err != nil {
		fmt.Errorf("Failed to replace default writer, %s", err.Error())
	}

	return &l
}

// Log at error level.
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.log(loggo.ERROR, fmt.Sprintf(msg, args...))
}

// Log at info level.
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.log(loggo.INFO, fmt.Sprintf(msg, args...))
}

// Log at critical level.
func (l *Logger) Criticalf(msg string, args ...interface{}) {
	l.log(loggo.CRITICAL, fmt.Sprintf(msg, args...))
}

// Log at debug level.
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.log(loggo.DEBUG, fmt.Sprintf(msg, args...))
}

// Log at trace level.
func (l *Logger) Tracef(msg string, args ...interface{}) {
	l.log(loggo.TRACE, fmt.Sprintf(msg, args...))
}

// Log at warning level.
func (l *Logger) Warningf(msg string, args ...interface{}) {
	l.log(loggo.WARNING, fmt.Sprintf(msg, args...))
}

// Logs at the given loggo.Level. If WithErrors or WithFields have been used,
// the same is appended to the message.
func (l *Logger) log(level loggo.Level, msg string) {
	errorStr := ""
	fieldsStr := ""
	if l.errorText != "" {
		errorStr = fmt.Sprintf("\terror=\"%s\"", l.errorText)
	}
	if l.fieldsText != "" {
		fieldsStr = fmt.Sprintf("\t%s", l.fieldsText)
	}

	// Call depth is set to 2, to get actual filename and line number in file
	// where Errorf, Infof,... methods are called
	l.logger.LogCallf(2, level, "msg=\"%s\"%s\t%s\n", msg, errorStr, fieldsStr)
}

// Log the given error as a seperate field.
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		l.errorText = ""
	} else {
		l.errorText = err.Error()
	}
	return l
}

// Add given key, value as custom field and value in log.
func (l *Logger) WithField(k string, v interface{}) *Logger {
	if k == "" {
		return l
	}
	l.fields[k] = v
	l.populateFieldsText()
	return l
}

// Add given key, value pairs as custom fields and values in log.
func (l *Logger) WithFields(kv map[string]interface{}) *Logger {
	for k, v := range kv {
		l.fields[k] = v
	}
	l.populateFieldsText()
	return l
}

// Sort key value pair by field and convert to string.
func (l *Logger) populateFieldsText() {
	kk := make([]string, len(l.fields))
	ctr := 0
	for k, _ := range l.fields {
		kk[ctr] = k
		ctr += 1
	}

	sort.Strings(kk)
	l.fieldsText = ""

	for _, k := range kk {
		l.fieldsText += fmt.Sprintf("%s=\"%+v\" ", k, l.fields[k])
	}

}

// Set log level
func (l *Logger) SetLogLevel(level_str string) {
	if level, ok := loggo.ParseLevel(level_str); ok == true {
		l.logger.SetLogLevel(level)
	}
}

// Current log level
func (l *Logger) LogLevel() string {
	return l.logger.LogLevel().String()
}
