package server

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model/structs"
	"baseGo/src/red_robot/conf"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"baseGo/src/github.com/jasonlvhit/gocron"
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

// 红包时间轮
func InitTimeWheel() {
	gocron.Every(0).Seconds().Do(timeWheel)
	gocron.NextRun()
	<-gocron.Start()
	s := gocron.NewScheduler()
	s.Every(0).Seconds().Do(timeWheel)
	<-s.Start()
}

func timeWheel() {
	redis := conf.GetRedis().Get()
	// 获取当前需要清除的信息格子
	currentCmd, err := redis.Do("Get", CURRENT)
	if err != nil {
		// 创建新的时间轮
		redis.Do("Set", CURRENT, START, 0)
	} else {
		fResult := string(currentCmd.([]byte))
		START, _ = strconv.Atoi(fResult)
	}
	// 获取当前需要清理的数据
	currentDataCmd, err := redis.Do("Get", fmt.Sprintf("%v_%v", CURRENT_DATA, START))
	redis.Do("Del", fmt.Sprintf("%v_%v", CURRENT_DATA, START))
	if err == nil {

		keys := string(currentDataCmd.([]byte))

		if err == nil && keys != "" {
			data := make([]map[string]string, 0)
			spkeys := strings.Split(keys, ";")
			item := make(map[string]bool)
			for _, v := range spkeys {
				if !item[v] {
					if len(v) > 0 {
						info := make(map[string]string)
						err = json.Unmarshal([]byte(v), &info)
						data = append(data, info)
						item[v] = true
					}
				}
			}
			go func(data []map[string]string) {
				for _, v := range data {
					switch v["key"] { // 红包结算
					case "redPacketSettle":
						go func(order map[string]string) {
							has := false
							if _, ok := order["settlementTime"]; ok {
								settlementTime, _ := strconv.Atoi(order["settlementTime"])
								if utility.GetNowTimestamp() < settlementTime {
									b, _ := json.Marshal(order)
									AddTimeWheel(string(b), 60)
									has = true
								}
							}
							if !has {
								// 结算红包
								redId, _ := strconv.Atoi(order["redId"])
								roomId, _ := strconv.Atoi(order["roomId"])
								agencyId, _ := order["agencyId"]
								_, err = new(RedPlay).RedEnvelopeAmountCalculation(order["lineId"], agencyId, redId, roomId)
								if err != nil {
									b, _ := json.Marshal(order)
									AddTimeWheel(string(b), 30)
								} else {
									// 删除红包缓存
									logKey := fmt.Sprintf("%v_%v_redLog", redId, roomId)
									conf.GetRedis().Get().Do("Del", logKey)
								}
							}
						}(v)
					case "OrdinaryRedPacket": // 普通红包自动创建
						go func(order map[string]string) {
							autoTime, _ := strconv.Atoi(order["autoTime"])
							if utility.GetNowTimestamp() < autoTime { // 时间没到 进入下一轮
								b, _ := json.Marshal(order)
								AddTimeWheel(string(b), 30)
							} else { // 时间到了 创建红包
								// 创建红包
								redInfo := new(structs.RedPacket)
								json.Unmarshal([]byte(order["redInfo"]), redInfo)
								gameTime, _ := strconv.Atoi(order["gameTime"])
								err = CreateOrdinaryRedPacket(redInfo, gameTime)
								if err != nil {
									b, _ := json.Marshal(order)
									AddTimeWheel(string(b), 30)
								}
							}
						}(v)
					case "robotGrabRedPacket":
						go func(order map[string]string) {
							redId, _ := strconv.Atoi(order["redId"])
							roomId, _ := strconv.Atoi(order["roomId"])
							new(RobotGrabRed).GrabRed(order["lineId"], order["agencyId"], roomId, redId)
						}(v)
					}
				}
			}(data)
		}
	}
	START++
	// 判断是否越界，越界重置时间轮
	if START > TOTAL {
		START -= TOTAL
	}
	// 设置下次需要清除的格子编号
	redis.Do("Set", CURRENT, START, 0)
}

// 添加到时间轮
func AddTimeWheel(key string, times int) {
	redis := conf.GetRedis().Get()
	// 获取当前执行格子信息
	currentCmd, err := redis.Do("Get", CURRENT)
	if err == nil {
		start := string(currentCmd.([]byte))

		times1, err := strconv.Atoi(start)
		if err == nil {
			total := times1 + times
			if total > TOTAL {
				total = total % TOTAL
			}
			// 将key信息写入redisFF
			redis.Do("Append", fmt.Sprintf("%v_%v", CURRENT_DATA, total), key+";")
		}
	}
}
