package server

// var (
// 	SystemLineBo = new(bo.SystemLineBo)
// )

// 获取线路信息
// func GetStytemLineInfo(lineId string) (*structs.SystemLine, error) {
// 	sess := conf.GetXormSession()
// 	defer sess.Close()
// 	lineInfo := new(structs.SystemLine)

// 	// 获取线路信息
// 	b, err := conf.GetRedis().HGet(model.SYSTEM_LINE_REDIS_KEY, lineId).Result()
// 	if err != nil {
// 		lineInfo, has, err := SystemLineBo.QueryLineBylineId(sess, lineId)
// 		if !has || err != nil {
// 			if err != nil {
// 				golog.Error("AccountService", "GetUserInfo", "err:", err)
// 			}
// 			return nil, err
// 		}
// 		lineByte, _ := json.Marshal(lineInfo)
// 		conf.GetRedis().HSet(model.SYSTEM_LINE_REDIS_KEY, lineId, string(lineByte))
// 	} else {
// 		json.Unmarshal([]byte(b), lineInfo)
// 	}
// 	return lineInfo, err
// }
