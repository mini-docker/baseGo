package middleware

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/red_agency/conf"
	"encoding/json"
	"fmt"
	"strconv"

	"baseGo/src/github.com/jasonlvhit/gocron"
)

var (
	agencyBo    = new(bo.SystemAgencyBo)
	redisClient = conf.GetRedis().Get()
)

func InitOnlineCheck() {
	s := gocron.NewScheduler()
	s.Every(1).Minutes().Do(onlineCheck)
	<-s.Start()
}

func onlineCheck() {
	// 获取所有session
	sessionCmd, err := redisClient.Do("HKeys", model.GetAgencyListKey())
	if err != nil {
		golog.Error("OnlineCheck", "onlineCheck", "err:", err)
		return
	}
	sess := conf.GetXormSession()
	defer sess.Close()

	keys := sessionCmd.([]interface{})

	for _, fv := range keys {
		v := string(fv.([]byte))
		isDel := false
		// 遍历解析所有sessionKey
		id, _ := strconv.Atoi(v)
		sessionKeys, err := redisClient.Do("HGet", model.GetAgencyListKey(), v)
		if err != nil {
			golog.Error("OnlineCheck", "onlineCheck", "err:", err)
			continue
		}
		se := new(SessionCache)

		// json.Unmarshal([]byte(sessionKeys.Val()), se)
		fResult := string(sessionKeys.([]byte))
		var dat SessionCache
		if err := json.Unmarshal([]byte(fResult), &dat); err == nil {
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
				continue
			}
			sessionBo := new(model.AgencySession)

			fResult := string(session.([]byte))
			json.Unmarshal([]byte(fResult), &agencyBo)
			// json.Unmarshal([]byte(session.Val()), &sessionBo)

			if sessionBo.TimeOut < utility.GetNowTimestamp() {
				// session过期，更新在线状态
				agency, has, _ := agencyBo.QueryAgencyById(sess, id)
				if has {
					// 修改在线状态为离线
					agency.IsOnline = model.OFFLINE
					agencyBo.EditAgencyOnlineStatus(sess, agency)
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
