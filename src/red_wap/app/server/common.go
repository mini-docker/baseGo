package server

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/conf"
	"encoding/json"
)

var (
	SystemLineBo = new(bo.SystemLineBo)
)

// 获取线路信息
func GetStytemLineInfo(lineId string) (*structs.SystemLine, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	lineInfo := new(structs.SystemLine)
	redisClient := conf.GetRedis().Get()
	// 获取线路信息
	b, err := redisClient.Do("HGet", model.SYSTEM_LINE_REDIS_KEY, lineId)
	if err != nil {
		lineInfo, has, err := SystemLineBo.QueryLineBylineId(sess, lineId)
		if !has || err != nil {
			if err != nil {
				golog.Error("AccountService", "GetUserInfo", "err:", err)
			}
			return nil, err
		}
		lineByte, _ := json.Marshal(lineInfo)
		redisClient.Do("HSet", model.SYSTEM_LINE_REDIS_KEY, lineId, string(lineByte))
	} else {
		c := string(b.([]byte))
		json.Unmarshal([]byte(c), lineInfo)
	}
	return lineInfo, err
}
