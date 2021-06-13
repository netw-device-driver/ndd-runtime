/*
Copyright 2020 Wim Henderickx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logging

import (
	"github.com/go-logr/logr"
)

// A Logger logs messages. Messages may be supplemented by structured data.
type Logger interface {
	// Info logs a message with optional structured data. Structured data must
	// be supplied as an array that alternates between string keys and values of
	// an arbitrary type. 
	Info(msg string, keysAndValues ...interface{})

	// Debug logs a message with optional structured data. Structured data must
	// be supplied as an array that alternates between string keys and values of
	// an arbitrary type. 
	Debug(msg string, keysAndValues ...interface{})

	// WithValues returns a Logger that will include the supplied structured
	// data with any subsequent messages it logs. Structured data must
	// be supplied as an array that alternates between string keys and values of
	// an arbitrary type.
	WithValues(keysAndValues ...interface{}) Logger
}

// NewNopLogger returns a Logger that does nothing.
func NewNopLogger() Logger { return nopLogger{} }

type nopLogger struct{}

func (l nopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l nopLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l nopLogger) WithValues(keysAndValues ...interface{}) Logger { return nopLogger{} }

// NewLogrLogger returns a Logger that is satisfied by the supplied logr.Logger,
// which may be satisfied in turn by various logging implementations (Zap, klog,
// etc).
func NewLogrLogger(l logr.Logger) Logger {
	return logrLogger{log: l}
}

type logrLogger struct {
	log logr.Logger
}

func (l logrLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Info(msg, keysAndValues...)
}

func (l logrLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log.V(1).Info(msg, keysAndValues...)
}

func (l logrLogger) WithValues(keysAndValues ...interface{}) Logger {
	return logrLogger{log: l.log.WithValues(keysAndValues...)}
}
