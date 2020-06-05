// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package golog

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"fecho/golog/filelog"
	"fecho/modules"
	"fecho/utility"
	"fecho/utility/atomicbool"
	"fmt"
)

// concept
/*
- Logging function:
  - check if file-based levelling enabled
    - if yes, check if Level is active on this file
  - check if Level is active
  - send Data to backend via big buffered channel
- Backend:
  - wait until there is Time for writing logs
  - write logs
  - configurable if logged to folder (buffer + rollingFileAppender) and/or console
  - console: log everything above INFO to stderr
- Channel overbuffering protection:
  - if buffer is full, trigger write
- Anti-Importing-Loop:
  - everything imports logging
  - logging is configured by main Module and is supplied access to configuration and taskmanager
*/

type severity uint32

func (s severity) String() string {
	switch s {
	case TraceLevel:
		return "TRAC"
	case DebugLevel:
		return "DEBU"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARN"
	case ErrorLevel:
		return "ERRO"
	case CriticalLevel:
		return "CRIT"
	case RecordLevel:
		return "RECORD"
	case SqlLevel:
		return "SQL"
	default:
		return "NONE"
	}
}

type logLine struct {
	Time   time.Time   `json:"time"`
	lvl    severity    `json:"-"`
	Level  string      `json:"level"`
	Line   string      `json:"line"`
	From   string      `json:"from"`
	Module string      `json:"module"`
	Method string      `json:"method"`
	msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func (l *logLine) MarshalJSON() ([]byte, error) {
	inInterface := make(map[string]interface{})
	switch l.Data.(type) {
	case map[string]string:
		inInterface = l.Data.(map[string]interface{})
	case interface{}:
		j, _ := json.Marshal(l.Data)
		err := json.Unmarshal(j, &inInterface)
		if err != nil {
			fmt.Println("log MarshalJSON error:", err)
		}
	default:
		inInterface["msg"] = l.Data
	}
	//inInterface["time"] = l.Time.Format("2006-01-02 15:04:05.999999")
	inInterface["time"] = l.Time.Format(time.RFC3339Nano)
	if v, ok := inInterface["dayType"]; ok {
		if dayType, ok := v.(float64); ok {
			inInterface["time"] = time.Unix(int64(dayType), 0).In(utility.ChinaZone).Format(time.RFC3339Nano)
		}
	}
	inInterface["level"] = l.Level
	if l.Line != "" {
		inInterface["line"] = l.Line
	}
	//TODO 排序
	inInterface["module"] = l.Module
	inInterface["method"] = l.Method
	inInterface["from"] = l.From

	return json.Marshal(inInterface)
}

func (ll *logLine) String() string {
	ll.Level = ll.lvl.String()

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.Encode(ll)

	return buf.String()
}

func (ll *logLine) ColorString() string {
	ll.Level = ll.lvl.String()
	//
	//buf := bytes.NewBuffer(nil)
	//encoder := json.NewEncoder(buf)
	//encoder.Encode(ll)
	//Time    time.Time   `json:"time"`
	//lvl     severity    `json:"-"`
	//Level   string      `json:"level"`
	//Line    string      `json:"line"`
	//SiteId  string      `json:"site_id"`
	//IndexId string      `json:"index_id"`
	//Module  string      `json:"module"`
	//Method  string      `json:"method"`
	//msg     string      `json:"msg"`
	//Data    interface{} `json:"data"`
	msg := fmt.Sprintf("%s:[%s] %s Module=%s,Method=%s\n\t%+v", ll.Time.Format(time.RFC3339Nano), ll.Level, ll.Line, ll.Module, ll.Method, ll.Data)
	color := 37
	if ll.lvl == WarningLevel {
		color = 35
	} else if ll.lvl == DebugLevel {
		color = 32
	} else if ll.lvl == ErrorLevel {
		color = 31
	} else if ll.lvl == CriticalLevel {
		color = 36
	}
	/*
			 30: 黑
			 31: 红 \n\x1b[0m", 31);
		     32: 绿 \n\x1b[0m", 32);
		     33: 黄 \n\x1b[0m", 33);
		     34: 蓝 \n\x1b[0m", 34);
		     35: 紫 \n\x1b[0m", 35);
		     36: 深绿 \n\x1b[0m", 36);
		     37: 白色 \n\x1b[0m", 37);
	*/
	pd := fmt.Sprintf("%c[0;%dm%s%c[0m", 0x1b, color, msg, 0x1b)
	pd = strings.TrimRight(pd, "\n")
	return pd
}

const (
	TraceLevel    severity = 1
	DebugLevel    severity = 2
	InfoLevel     severity = 3
	WarningLevel  severity = 4
	ErrorLevel    severity = 5
	CriticalLevel severity = 6

	RecordLevel severity = 11
	SqlLevel    severity = 12

	LogPath = "/home/www/logs/applogs" //${APP_NAME}/${POD_NAME}/"

)

var (
	module *modules.Module

	fileLog *GoLogger //*filelog.Logger

	logPath        string
	RecordFilePath string

	logBuffer             chan *logLine
	forceEmptyingOfBuffer chan bool

	logLevel *uint32

	fileLevelsActive *atomicbool.AtomicBool
	fileLevels       map[string]severity
	fileLevelsLock   sync.RWMutex

	logsWaiting     chan bool
	logsWaitingFlag *atomicbool.AtomicBool
)

func SetFileLevels(levels map[string]severity) {
	fileLevelsLock.Lock()
	fileLevels = levels
	fileLevelsLock.Unlock()
	fileLevelsActive.Set()
}

func UnSetFileLevels() {
	fileLevelsActive.UnSet()
}

func SetLogLevel(level severity) {
	atomic.StoreUint32(logLevel, uint32(level))
}

func ParseLevel(level string) severity {
	switch strings.ToLower(level) {
	case "trace":
		return 1
	case "debug":
		return 2
	case "info":
		return 3
	case "warning":
		return 4
	case "error":
		return 5
	case "critical":
		return 6
	}
	return 0
}

func init() {
	module = modules.Register("Logging", 0)
	logBuffer = make(chan *logLine, 1024)
	forceEmptyingOfBuffer = make(chan bool, 4)

	fileLevelsActive = atomicbool.NewBool(false)
	fileLevels = nil

	logsWaiting = make(chan bool, 1)
	logsWaitingFlag = atomicbool.NewBool(false)
	//${APP_NAME}/${POD_NAME}/"
	logPath = path.Join(LogPath, os.Getenv("APP_NAME"), os.Getenv("POD_NAME"))
	os.MkdirAll(logPath, os.ModePerm)
	log.SetFlags(0)
	initLogLevel := uint32(3)
	logLevel = &initLogLevel
	go writer() //写入
}

func LogInit(leve, filename string, maxsize, MaxBackups, MaxAge int, Compress bool) error {
	initialLogLevel := ParseLevel(leve)
	if initialLogLevel == 0 {
		initialLogLevel = 3
	}
	logLevelInt := uint32(initialLogLevel)
	logLevel = &logLevelInt

	if !strings.HasPrefix(filename, "/") {
		filename = path.Join(logPath, filename)
	} else {
		//开发使用
		//exePath,_ := utility.ExecPath()
		//filename = path.Join(exePath,filename)
	}
	dir := filepath.Dir(filename)
	if utility.IsDir(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("warn log file can not create dir write logs in tty Output:", filename)
			return nil
		}
	}
	if !utility.IsExist(filename) {
		f, err := os.Create(filename)
		if err != nil {
			log.Println("warn log file can not create write logs in tty Output:", filename)
			return nil
		} else {
			f.Close()
		}
	}

	logger := &GoLogger{}
	if utility.IsExist(filename) {
		logger.fileLogger = &filelog.Logger{
			Filename:   filename,
			MaxSize:    maxsize, // megabytes
			MaxBackups: MaxBackups,
			MaxAge:     MaxAge,   //days
			Compress:   Compress, // disabled by default
		}
	}

	//log.Println("log file ", filename)
	//if utility.IsExist(filename) {
	//	fileLog = &filelog.Logger{
	//		Filename:   filename,
	//		MaxSize:    maxsize, // megabytes
	//		MaxBackups: MaxBackups,
	//		MaxAge:     MaxAge,   //days
	//		Compress:   Compress, // disabled by default
	//	}
	//	log.SetOutput(fileLog)
	//}

	// --------- record logger start ------------
	if RecordFilePath != "" {
		if !strings.HasPrefix(RecordFilePath, "/") {
			RecordFilePath = path.Join(logPath, RecordFilePath)
		} else {
			//开发使用
			//exePath,_ := utility.ExecPath()
			//recordFileName = path.Join(exePath,recordFileName)
		}
		dir = filepath.Dir(RecordFilePath)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("warn log file can not create write logs in tty Output:", RecordFilePath)
			return nil
		}
		if !utility.IsExist(RecordFilePath) {
			f, err := os.Create(RecordFilePath)
			if err != nil {
				log.Println("warn log file can not create write logs in tty Output:", RecordFilePath)
				return nil
			} else {
				f.Close()
			}
		}

		if utility.IsExist(RecordFilePath) {
			logger.recordLogger = &filelog.Logger{
				Filename:   RecordFilePath,
				MaxSize:    maxsize, // megabytes
				MaxBackups: MaxBackups,
				MaxAge:     MaxAge,   //days
				Compress:   Compress, // disabled by default
			}
		}
	}

	if logger.fileLogger != nil || logger.recordLogger != nil {
		fileLog = logger
		log.SetOutput(fileLog)
	}

	// --------- record logger end ------------

	return nil
}
