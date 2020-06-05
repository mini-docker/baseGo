package golog

var Logger *LoggingInterface

type LoggingInterface struct {
	From string
}

func (l *LoggingInterface) Trace(module string, method string, msg string, args ...interface{}) {
	Trace(module, method, msg, args...)
}
func (l *LoggingInterface) Debug(module string, method string, msg string, args ...interface{}) {
	Debug(module, method, msg, args...)
}
func (l *LoggingInterface) Info(module string, method string, msg string, args ...interface{}) {
	Info(module, method, msg, args...)
}
func (l *LoggingInterface) Warn(module string, method string, format string, err error, args ...interface{}) {
	Warn(module, method, format, err, args...)
}
func (l *LoggingInterface) Error(module string, method string, format string, err error, args ...interface{}) {
	Error(module, method, format, err, args...)
}
func (l *LoggingInterface) Fatal(module string, method string, format string, err error, args ...interface{}) {
	Fatal(module, method, format, err, args...)
}

//记录一些重要日志
func (l *LoggingInterface) Record(module string, method string, args interface{}) {
	Record(module, method, args)
}
func init() {
	Logger = &LoggingInterface{}
}
