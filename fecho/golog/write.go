package golog

import (
	"log"
	"time"

	"github.com/mini-docker/baseGo/fecho/modules/taskmanager"
	"github.com/mini-docker/baseGo/fecho/utility"
)

func writeLine(line *logLine) {
	if fileLog == nil {
		log.Print(line.ColorString())
	} else {
		log.Print(line)
	}
	//写入文件
	// TODO: implement file logging and setting console/file logging
	// TODO: use https://github.com/natefinch/lumberjack
}

func writer() {
	var line *logLine
	startedTask := false

	for {

		// wait until logs need to be processed
		select {
		case <-logsWaiting:
			logsWaitingFlag.UnSet()
		case <-module.Stop:
		}

		// wait for timeslot to log, or when buffer is full
		select {
		case <-taskmanager.StartVeryLowPriorityMicroTask():
			startedTask = true
		case <-forceEmptyingOfBuffer:
		case <-module.Stop:
			select {
			case line = <-logBuffer:
				writeLine(line)
			case <-time.After(200 * time.Millisecond): //延时1s关闭 给logBuffer时间
				writeLine(&logLine{
					lvl:  WarningLevel,
					Time: utility.GetNowTime(),
					Data: map[string]interface{}{
						"msg": "===== LOGGING STOPPED =====",
					},
				})
				if fileLog != nil {
					fileLog.Rotate()
				}
				module.StopComplete()
				return
			}
		}

		// write all the logs!
	writeLoop:
		for {
			select {
			case line = <-logBuffer:
				writeLine(line)
			default:
				if startedTask {
					taskmanager.EndMicroTask()
					startedTask = false
				}
				break writeLoop
			}
		}

	}
}
