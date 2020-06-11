package golog

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/mini-docker/baseGo/fecho/utility"
)

func log_fastcheck(level severity) bool {
	if fileLevelsActive.IsSet() {
		return true
	}
	if uint32(level) < atomic.LoadUint32(logLevel) {
		return false
	}
	return true
}

func logPrintRecord(level severity, from, module, method string, data interface{}) {
	// check if Level is enabled
	if !fileLevelsActive.IsSet() && uint32(level) < atomic.LoadUint32(logLevel) {
		return
	}
	// get Time
	now := utility.GetNowTime()
	// create log object
	log := &logLine{
		Time:   now,
		lvl:    level,
		Module: module,
		From:   from,
		Method: method,
		Data:   data,
	}
	// send log to processing
	select {
	case logBuffer <- log:
	default:
		forceEmptyingOfBuffer <- true
		logBuffer <- log
	}

	// wake up writer if necessary
	if logsWaitingFlag.SetToIf(false, true) {
		logsWaiting <- true
	}
}

func logPrint(level severity, from, module, method, msg string, args ...interface{}) {
	// check if Level is enabled
	if !fileLevelsActive.IsSet() && uint32(level) < atomic.LoadUint32(logLevel) {
		return
	}

	// get Time
	now := utility.GetNowTime()

	var (
		file string
		line int
		ok   bool
	)

	if level >= ErrorLevel {
		// get file and line
		_, file, line, ok = runtime.Caller(2)
		if !ok {
			file = "?"
			line = 0
		} else {
			fPartStart := strings.LastIndex(file, "/src/") + 5 // .Split(ll.file, "/src/")
			file = fmt.Sprintf("%s:%d", file[fPartStart:], line)
		}
	}

	argSize := len(args)
	num := len(args) / 2
	argsMap := make(map[string]interface{})
	if argSize%2 == 0 {
		argsMap["msg"] = msg
		for i := 0; i < num; i++ {
			argsMap[escape(fmt.Sprintf("%v", args[i*2]), false)] = escape(fmt.Sprintf("%v", args[i*2+1]), false)
		}
	} else {
		fmtStr := strings.Repeat("%v,", argSize)
		msg = fmt.Sprintf("%s "+strings.Trim(fmtStr, ", "), msg, args)
		argsMap["msg"] = msg
	}
	// create log object
	log := &logLine{
		Time:   now,
		lvl:    level,
		Line:   file,
		From:   from,
		Module: module,
		Method: method,
		Data:   argsMap,
	}
	// send log to processing
	select {
	case logBuffer <- log:
	default:
		forceEmptyingOfBuffer <- true
		logBuffer <- log
	}

	// wake up writer if necessary
	if logsWaitingFlag.SetToIf(false, true) {
		logsWaiting <- true
	}
}

func escape(s string, filterEqual bool) string {
	dest := make([]byte, 0, 2*len(s))
	for i := 0; i < len(s); i++ {
		r := s[i]
		switch r {
		case '|':
			continue
		case '%':
			dest = append(dest, '%', '%')
		case '=':
			if !filterEqual {
				dest = append(dest, '=')
			}
		default:
			dest = append(dest, r)
		}
	}

	return string(dest)
}

func Trace(module, method, msg string, args ...interface{}) {
	if log_fastcheck(TraceLevel) {
		logPrint(TraceLevel, Logger.From, module, method, msg, args...)
	}
}
func Debug(module, method, msg string, args ...interface{}) {
	if log_fastcheck(DebugLevel) {
		logPrint(DebugLevel, Logger.From, module, method, msg, args...)
	}
}
func Info(module, method, msg string, args ...interface{}) {
	if log_fastcheck(InfoLevel) {
		logPrint(InfoLevel, Logger.From, module, method, msg, args...)
	}
}

//siteId, indexId, module, method, format, err, args...)
func Warn(module, method, msg string, err error, args ...interface{}) {
	if log_fastcheck(WarningLevel) {
		if err != nil {
			if strings.Contains(msg, "%") {
				msg = fmt.Sprintf(msg, err)
			} else {
				msg = msg + fmt.Sprintf(" %v", err)
			}
		}
		logPrint(WarningLevel, Logger.From, module, method, msg, args...)
	}
}
func Error(module, method, msg string, err error, args ...interface{}) {
	if log_fastcheck(ErrorLevel) {
		if err != nil {
			if strings.Contains(msg, "%") {
				msg = fmt.Sprintf(msg, err)
			} else {
				msg = msg + fmt.Sprintf(" %v", err)
			}
		}
		logPrint(ErrorLevel, Logger.From, module, method, msg, args...)
	}
}
func Fatal(module, method, msg string, err error, args ...interface{}) {
	if log_fastcheck(CriticalLevel) {
		if err != nil {
			if strings.Contains(msg, "%") {
				msg = fmt.Sprintf(msg, err)
			} else {
				msg = msg + fmt.Sprintf(" %v", err)
			}
		}
		logPrint(CriticalLevel, Logger.From, module, method, msg, args...)
	}
}

func Record(module string, method string, data interface{}) {
	logPrintRecord(RecordLevel, Logger.From, module, method, data)
}

func Testf(things ...interface{}) {
	fmt.Printf(things[0].(string), things[1:]...)
}

func Test(msg string) {
	fmt.Println(msg)
}
