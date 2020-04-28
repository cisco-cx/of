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
	"io"
	"log"
	"runtime"

	"github.com/sirupsen/logrus"
	of "github.com/cisco-cx/of/pkg/v2"
)

// Represents loggo logger and fields to support structured logging.
type Logger struct {
	entry *logrus.Entry
}

type WithLogger struct {
	logger *Logger
	fields map[string]interface{}
	err    error
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
	return &l
}

// Log correct file name and line number from where Logger call was invoked.
func prettyfier(r *runtime.Frame) (string, string) {
	return "", ""
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
	wl := WithLogger{
		logger: l,
		err:    err,
	}
	return &wl
}

// Add given key, value as custom field and value in log.
func (l *Logger) WithField(k string, v interface{}) of.Logger {
	wl := WithLogger{
		logger: l,
		fields: map[string]interface{}{
			k: v,
		},
	}
	return &wl
}

// Add given key, value pairs as custom fields and values in log.
func (l *Logger) WithFields(kv map[string]interface{}) of.Logger {
	wl := WithLogger{
		logger: l,
		fields: kv,
	}
	return &wl
}

// Log at given log level.
func (l *Logger) Logf(level logrus.Level, msg string, args ...interface{}) {
	l.entry.Logf(level, msg, args...)
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

// Log at given log level.
func (wl *WithLogger) Logf(level logrus.Level, msg string, args ...interface{}) {
	e := wl.logger.entry
	if wl.err != nil {
		e = e.WithError(wl.err)
	}

	if len(wl.fields) != 0 {
		e = e.WithFields(wl.fields)
	}
	e.Logf(level, msg, args...)
}

// Log at error level.
func (wl *WithLogger) Errorf(msg string, args ...interface{}) {
	wl.Logf(logrus.ErrorLevel, msg, args...)
}

// Log at info level.
func (wl *WithLogger) Infof(msg string, args ...interface{}) {
	wl.Logf(logrus.InfoLevel, msg, args...)
}

// Log at fatal level.
func (wl *WithLogger) Fatalf(msg string, args ...interface{}) {
	if wl.logger.entry.Logger.IsLevelEnabled(logrus.FatalLevel) {
		wl.Logf(logrus.FatalLevel, msg, args...)
		wl.logger.entry.Logger.Exit(1)
	}
}

// Log at debug level.
func (wl *WithLogger) Debugf(msg string, args ...interface{}) {
	wl.Logf(logrus.DebugLevel, msg, args...)
}

// Log at trace level.
func (wl *WithLogger) Tracef(msg string, args ...interface{}) {
	wl.Logf(logrus.TraceLevel, msg, args...)
}

// Log at warning level.
func (wl *WithLogger) Warningf(msg string, args ...interface{}) {
	wl.Logf(logrus.WarnLevel, msg, args...)
}

func (l *WithLogger) Panicf(msg string, args ...interface{}) {
	l.logger.entry.Logger.Panicf(msg, args...)
}

// Log the given error as a seperate field.
func (wl *WithLogger) WithError(err error) of.Logger {
	wl.err = err
	return wl
}

// Add given key, value as custom field and value in log.
func (wl *WithLogger) WithField(k string, v interface{}) of.Logger {
	wl.fields = map[string]interface{}{
		k: v,
	}
	return wl
}

// Add given key, value pairs as custom fields and values in log.
func (wl *WithLogger) WithFields(kv map[string]interface{}) of.Logger {
	wl.fields = kv
	return wl
}

// Set log level
func (wl *WithLogger) SetLevel(level_str string) {
	wl.logger.Panicf("Unsupported method.")
}

// Current log level
func (wl *WithLogger) LogLevel() string {
	wl.logger.Panicf("Unsupported method.")
	return ""
}

// Change output. Default output is os.Stderr.
func (wl *WithLogger) SetOutput(w io.Writer) {
	wl.logger.Panicf("Unsupported method.")
}
