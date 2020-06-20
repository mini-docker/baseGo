package middleware

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/red_admin/conf"
	"encoding/json"
	"fmt"
	"strconv"

	"baseGo/src/github.com/jasonlvhit/gocron"
)

var (
	adminBo = new(bo.SystemAdminBo)
)

func InitOnlineCheck() {
	s := gocron.NewScheduler() // 创建一个定时器任务
	// 每一分钟执行执行一次 onlineCheck
	s.Every(1).Minutes().Do(onlineCheck)
	// s.Every(5).Seconds().Do(onlineCheck)
	<-s.Start()
}

func onlineCheck() {
	redisClient := conf.GetRedis().Get()
	// 获取所有session
	sessionCmd, err := redisClient.Do("HKeys", model.GetAdminListKey())
	fmt.Println(sessionCmd, "sessionCmd")
	if err != nil {
		golog.Error("OnlineCheck", "onlineCheck", "err:", err)
		return
	}
	sess := conf.GetXormSession()
	defer sess.Close()

	// var keys = make([]string, 0)
	if sessionCmd != nil {
		keys := sessionCmd.([]interface{})
		fmt.Println(keys, "keyskeys")
		// for _, v := range keys {
		// 	fv := string(v.([]byte))
		// 	fmt.Println(fv, "vvvvv")
		// }

		// _, err := json.Unmarshal([]byte(sessionCmd), &keys)
		// if err != nil {
		// 	golog.Error("OnlineCheck", "onlineCheck", "err:", err)
		// 	return
		// }
		// fmt.Println(keysf)
		// var dat SessionCache
		// if err := json.Unmarshal([]byte(keysf), &keys); err == nil {
		// 	fmt.Println("==============json str 转map=======================")
		// 	fmt.Println(keys, "keysf_keys")
		// 	// fmt.Println(dat["sess"])
		// }

		for _, fv := range keys {
			v := string(fv.([]byte))
			isDel := false
			// 遍历解析所有sessionKey
			id, _ := strconv.Atoi(v)
			sessionKeys, err := redisClient.Do("HGet", model.GetAdminListKey(), v)
			if err != nil {
				golog.Error("OnlineCheck", "onlineCheck", "err:", err)
				continue
			}
			se := new(SessionCache)

			// json.Unmarshal([]byte(sessionKeys.Val()), se)
			fResults := string(sessionKeys.([]byte))
			var dat SessionCache
			if err := json.Unmarshal([]byte(fResults), &dat); err == nil {
				fmt.Println("==============json str 转map=======================")
				fmt.Println(dat, "fResults_dat")
				// fmt.Println(dat["sess"])
			}
			se = &dat

			for _, ses := range se.Sess {
				// 遍历所有session,检测是否过期
				session, err := redisClient.Do("Get", ses)
				if err != nil {
					golog.Error("OnlineCheck", "onlineCheck", "err:", err)
					redisClient.Do("HDel", model.GetAdminListKey(), v)
					continue
				}
				sessionBo := new(model.AdminSession)
				fResults := string(session.([]byte))
				// json.Unmarshal([]byte(session.Val()), &sessionBo)
				// var dat SessionCache
				json.Unmarshal([]byte(fResults), &sessionBo)

				if sessionBo.TimeOut < utility.GetNowTimestamp() {
					// session过期，更新在线状态
					admin, has, _ := adminBo.QueryAdminById(sess, id)
					if has {
						// 修改在线状态为离线
						admin.IsOnline = model.OFFLINE
						adminBo.EditAdminOnlineStatus(sess, admin)
						// 删除session
						redisClient.Do("Del", ses)
						isDel = true
					}
					break
				}
			}
			if isDel {
				// 删除sessionListKey
				redisClient.Do("HDel", model.GetAgencyListKey(), v)
			}
		}
	}

}
