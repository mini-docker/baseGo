// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package modules

var logger Logger
var loggerRegistered chan bool

type Logger interface {
	Trace(module string, method string, msg string, args ...interface{})
	Debug(module string, method string, msg string, args ...interface{})
	Info(module string, method string, msg string, args ...interface{})
	Warn(module string, method string, format string, err error, args ...interface{})
	Error(module string, method string, format string, err error, args ...interface{})
	Fatal(module string, method string, format string, err error, args ...interface{})
	Record(module string, method string, args interface{})
}

func RegisterLogger(newLogger Logger) {
	if logger == nil {
		logger = newLogger
		loggerRegistered <- true
	}
}

func GetLogger() Logger {
	return logger
}
