package conf

import (
	"fecho/golog"
	"path/filepath"
	"strings"

	"fecho/utility"
	"os"
)

// initLog 系统运行初始化日志
func InitLog(config LogConfig) error {
	// --------- record logger start ------------

	recordLogPath := config.RecordPath
	if !strings.HasPrefix(recordLogPath, "/") && os.Getenv("RECORD_LOG_OUTPUT") != "" {
		path, err := utility.ExecPath()
		if err != nil {
			return err
		}
		//如果是test的路径
		if strings.Contains(path, "go-build") {
			goPaths := utility.GetGOPATHs()
			if len(goPaths) > 0 && len(goPaths[0]) > 0 {
				path = filepath.Join(goPaths[0], "bin")
			}
		}
		recordLogPath = filepath.Join(path, recordLogPath)
	}
	//golog.RecordFilePath = recordLogPath

	err := golog.LogInit(config.Level, recordLogPath, config.MaxLogSize, 6, 30, true)
	if err != nil {
		return err
	}

	// --------- record logger end ------------

	return nil
}
