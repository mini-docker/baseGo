package server

import (
	"baseGo/src/red_agency/conf"
	"fmt"
	"strconv"
)

var (
	total = 600
	//i     = 1
)

const (
	CURRENT      = "red_time_wheel_current"      // 当前需要清除的数据格子
	CURRENT_DATA = "red_time_wheel_current_data" // 当前需要清除格子数据
)

var START = 1   // 初始格子
var TOTAL = 600 // 总格子数

// 添加到时间轮
func AddTimeWheel(key string, times int) {
	redisClient := conf.GetRedis().Get()
	// 获取当前执行格子信息
	// currentCmd := conf.GetRedis().Get(CURRENT)
	currentCmd, err := redisClient.Do("Get", CURRENT)
	if err == nil && currentCmd != nil {
		// statrt := string(currentCmd.([]byte))
		statrt := string(currentCmd.([]byte))
		times1, err := strconv.Atoi(statrt)
		START = times1
		if err == nil {
			total := times1 + times
			if total > TOTAL {
				total = total % TOTAL
			}
			// 将key信息写入redisFF
			redisClient.Do("Append", fmt.Sprintf("%v_%v", CURRENT_DATA, total), key+";")
		}
	}
}
