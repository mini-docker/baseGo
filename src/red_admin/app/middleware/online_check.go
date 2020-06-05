package middleware

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/github.com/jasonlvhit/gocron"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/red_admin/conf"
	"encoding/json"

	"strconv"
)

var (
	adminBo = new(bo.SystemAdminBo)
)

func InitOnlineCheck() {
	s := gocron.NewScheduler() // 创建一个定时器任务
	// 每一分钟执行执行一次 onlineCheck
	s.Every(1).Minutes().Do(onlineCheck)
	<-s.Start()
}

func onlineCheck() {
	// 获取所有session
	sessionCmd := conf.GetRedis().HKeys(model.GetAdminListKey())
	if sessionCmd.Err() != nil {
		golog.Error("OnlineCheck", "onlineCheck", "err:", sessionCmd.Err())
		return
	}
	sess := conf.GetXormSession()
	defer sess.Close()
	keys := sessionCmd.Val()
	for _, v := range keys {
		isDel := false
		// 遍历解析所有sessionKey
		id, _ := strconv.Atoi(v)
		sessionKeys := conf.GetRedis().HGet(model.GetAdminListKey(), v)
		if sessionKeys.Err() != nil {
			golog.Error("OnlineCheck", "onlineCheck", "err:", sessionKeys.Err())
			continue
		}
		se := new(SessionCache)
		json.Unmarshal([]byte(sessionKeys.Val()), se)
		for _, ses := range se.Sess {
			// 遍历所有session,检测是否过期
			session := conf.GetRedis().Get(ses)
			if session.Err() != nil {
				golog.Error("OnlineCheck", "onlineCheck", "err:", session.Err())
				conf.GetRedis().HDel(model.GetAdminListKey(), v)
				continue
			}
			sessionBo := new(model.AdminSession)
			json.Unmarshal([]byte(session.Val()), &sessionBo)
			if sessionBo.TimeOut < utility.GetNowTimestamp() {
				// session过期，更新在线状态
				admin, has, _ := adminBo.QueryAdminById(sess, id)
				if has {
					// 修改在线状态为离线
					admin.IsOnline = model.OFFLINE
					adminBo.EditAdminOnlineStatus(sess, admin)
					// 删除session
					conf.GetRedis().Del(ses)
					isDel = true
				}
				break
			}
		}
		if isDel {
			// 删除sessionListKey
			conf.GetRedis().HDel(model.GetAgencyListKey(), v)
		}
	}
}
