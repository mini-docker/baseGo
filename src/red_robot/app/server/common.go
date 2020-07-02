package server

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/structs"
	"baseGo/src/red_robot/conf"
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
	// 获取线路信息

	b, err := conf.GetRedis().Get().Do("HGet", model.SYSTEM_LINE_REDIS_KEY, lineId)

	if err != nil {
		lineInfo, has, err := SystemLineBo.QueryLineBylineId(sess, lineId)
		if !has || err != nil {
			if err != nil {
				golog.Error("AccountService", "GetUserInfo", "err:", err)
			}
			return nil, err
		}
		lineByte, _ := json.Marshal(lineInfo)
		conf.GetRedis().Get().Do("HSet", model.SYSTEM_LINE_REDIS_KEY, lineId, string(lineByte))
	} else {
		fResult := string(b.([]byte))
		json.Unmarshal([]byte(fResult), lineInfo)
	}
	return lineInfo, err
}
