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

package v2

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
	of "github.com/cisco-cx/of/pkg/v2"
)

var mutex = &sync.Mutex{}

// Represents loggo logger and fields to support structured logging.
type Logger struct {
	entry *logrus.Entry
	// if true , Will clear any field that is set using WithField(s) call after a log line is logged/printed.
	autoClearFields bool // Default: true
}

// Initiate logger.
func New() *Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier:       prettyfier,
		DisableLevelTruncation: true,
	})
	e := logrus.NewEntry(logger)
	l := Logger{entry: e}
	l.AutoClearFields(true)
	return &l
}

// Log correct file name and line number from where Logger call was invoked.
func prettyfier(r *runtime.Frame) (string, string) {

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, 25)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {

		file := filepath.Base(f.File)
		// If the caller isn't part of logrus files, we're done
		if file != "log.go" && file != "entry.go" {
			return "", fmt.Sprintf("%s:%d", file, f.Line)
		}
	}

	return "", ""
}

// if true , Will clear any field that is set using WithField(s) call after a log line is logged/printed.
func (l *Logger) AutoClearFields(enabled bool) {
	l.autoClearFields = enabled
}

// Reset all fields set by WithField(s) method.
func (l *Logger) ClearFields() {
	mutex.Lock()
	defer mutex.Unlock()
	l.entry.Data = make(logrus.Fields)
}

// Log at error level.
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, msg, args...)
}

// Log at info level.
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, msg, args...)
}

// Log at fatal level.
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.FatalLevel) {
		l.Logf(logrus.FatalLevel, msg, args...)
		l.entry.Logger.Exit(1)
	}
}

// Log at panic level.
func (l *Logger) Panicf(msg string, args ...interface{}) {
	i := make([]interface{}, 0)
	i = append(i, msg)
	i = append(i, args...)
	log.Panic(i)
}

// Log at debug level.
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, msg, args...)
}

// Log at trace level.
func (l *Logger) Tracef(msg string, args ...interface{}) {
	l.Logf(logrus.TraceLevel, msg, args...)
}

// Log at warning level.
func (l *Logger) Warningf(msg string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, msg, args...)
}

// Log the given error as a seperate field.
func (l *Logger) WithError(err error) of.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	l.entry = l.entry.WithError(err)
	return l
}

// Add given key, value as custom field and value in log.
func (l *Logger) WithField(k string, v interface{}) of.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	l.entry = l.entry.WithField(k, v)
	return l
}

// Add given key, value pairs as custom fields and values in log.
func (l *Logger) WithFields(kv map[string]interface{}) of.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	l.entry = l.entry.WithFields(logrus.Fields(kv))
	return l
}

// Log at given log level.
func (l *Logger) Logf(level logrus.Level, msg string, args ...interface{}) {
	mutex.Lock()
	l.entry.Logf(level, msg, args...)
	delete(l.entry.Data, logrus.ErrorKey)
	mutex.Unlock()
	if l.autoClearFields == true {
		l.ClearFields()
	}
}

// Set log level
func (l *Logger) SetLevel(level_str string) {
	if level, err := logrus.ParseLevel(level_str); err == nil {
		l.entry.Logger.SetLevel(level)
	}
}

// Current log level
func (l *Logger) LogLevel() string {
	return l.entry.Logger.GetLevel().String()
}

// Change output. Default output is os.Stderr.
func (l *Logger) SetOutput(w io.Writer) {
	l.entry.Logger.SetOutput(w)
}

// Change output. Default output is os.Stderr.
func (l *Logger) Logger() *logrus.Logger {
	return l.entry.Logger
}
