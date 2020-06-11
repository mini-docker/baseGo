package golog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mini-docker/baseGo/fecho/golog/filelog"
)

// 日志
// 持有两个filelog中的Logger
// 一个作为正常日志输出
// 一个作为统计日志输出
// 这个结构实现filelog中的Logger所实现的方法
type GoLogger struct {
	// 正常日志输出
	fileLogger *filelog.Logger
	// 统计日志输出
	recordLogger *filelog.Logger
}

func (l *GoLogger) Write(p []byte) (n int, err error) {

	// 反序列化日志数据
	ll := &logLine{}
	json.Unmarshal(p, ll)

	// 解析日志是否是统计日志
	switch ll.Level {
	case fmt.Sprint(RecordLevel):
		if l.recordLogger != nil {
			return l.recordLogger.Write(p)
		}
	default:
		if l.fileLogger != nil {
			return l.fileLogger.Write(p)
		}
	}

	return 0, nil
}

func (l *GoLogger) Close() error {

	errs := make([]string, 0)

	if l.fileLogger != nil {
		fErr := l.fileLogger.Close()
		if fErr != nil {
			errs = append(errs, "file logger close error : "+fErr.Error())
		}
	}

	if l.recordLogger != nil {
		rErr := l.recordLogger.Close()
		if rErr != nil {
			errs = append(errs, "record logger close error : "+rErr.Error())
		}
	}

	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ","))
	}

	return nil
}

func (l *GoLogger) Rotate() error {
	errs := make([]string, 0)

	if l.fileLogger != nil {
		fErr := l.fileLogger.Rotate()
		if fErr != nil {
			errs = append(errs, "file logger rotate error : "+fErr.Error())
		}
	}

	if l.recordLogger != nil {
		rErr := l.recordLogger.Rotate()
		if rErr != nil {
			errs = append(errs, "record logger rotate error : "+rErr.Error())
		}
	}

	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ","))
	}

	return nil
}
