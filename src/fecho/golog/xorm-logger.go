//兼容mysql的日志
package golog

import (
	"fecho/xorm/core"
	"errors"
	"sync/atomic"
)

type LoggingXorm struct {
	MShowSQL bool
}

func (l *LoggingXorm) Debug(v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Debug("xorm", "-", "", v...)
	}
}
func (l *LoggingXorm) Debugf(format string, v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Debug("xorm", "-", format, v...)
	}
}
func (l *LoggingXorm) Error(v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Error("xorm", "-", "", errors.New(""), v...)
	}
}
func (l *LoggingXorm) Errorf(format string, v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Error("xorm", "-", format, errors.New(""), v...)
	}
}
func (l *LoggingXorm) Info(v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Info("xorm", "-", "", v...)
	}
}
func (l *LoggingXorm) Infof(format string, v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Info("xorm", "-", format, v...)
	}
}
func (l *LoggingXorm) Warn(v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Warn("xorm", "-", "", errors.New(""), v...)
	}
}
func (l *LoggingXorm) Warnf(format string, v ...interface{}) {
	if l.IsShowSQL() {
		Logger.Warn("xorm", "-", format, errors.New(""), v...)
	}
}
func (*LoggingXorm) ShowSql(format string, v ...interface{}) {
}
func (*LoggingXorm) Level() core.LogLevel {
	return core.LogLevel(atomic.LoadUint32(logLevel))
}

//废弃
func (*LoggingXorm) SetLevel(l core.LogLevel) {

}

//废弃
func (l *LoggingXorm) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.MShowSQL = true
	} else {
		l.MShowSQL = show[0]
	}
}
func (l *LoggingXorm) IsShowSQL() bool {
	return l.MShowSQL
}
