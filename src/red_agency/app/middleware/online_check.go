package middleware

// var (
// 	agencyBo = new(bo.SystemAgencyBo)
// )

// func InitOnlineCheck() {
// 	s := gocron.NewScheduler()
// 	s.Every(1).Minutes().Do(onlineCheck)
// 	<-s.Start()
// }

// func onlineCheck() {
// 	// 获取所有session
// 	sessionCmd := conf.GetRedis().HKeys(model.GetAgencyListKey())
// 	if sessionCmd.Err() != nil {
// 		golog.Error("OnlineCheck", "onlineCheck", "err:", sessionCmd.Err())
// 		return
// 	}
// 	sess := conf.GetXormSession()
// 	defer sess.Close()
// 	keys := sessionCmd.Val()
// 	for _, v := range keys {
// 		isDel := false
// 		// 遍历解析所有sessionKey
// 		id, _ := strconv.Atoi(v)
// 		sessionKeys := conf.GetRedis().HGet(model.GetAgencyListKey(), v)
// 		if sessionKeys.Err() != nil {
// 			golog.Error("OnlineCheck", "onlineCheck", "err:", sessionKeys.Err())
// 			continue
// 		}
// 		se := new(SessionCache)
// 		json.Unmarshal([]byte(sessionKeys.Val()), se)
// 		for _, ses := range se.Sess {
// 			// 遍历所有session,检测是否过期
// 			session := conf.GetRedis().Get(ses)
// 			if session.Err() != nil {
// 				golog.Error("OnlineCheck", "onlineCheck", "err:", session.Err())
// 				continue
// 			}
// 			sessionBo := new(model.AgencySession)
// 			json.Unmarshal([]byte(session.Val()), &sessionBo)
// 			if sessionBo.TimeOut < utility.GetNowTimestamp() {
// 				// session过期，更新在线状态
// 				agency, has, _ := agencyBo.QueryAgencyById(sess, id)
// 				if has {
// 					// 修改在线状态为离线
// 					agency.IsOnline = model.OFFLINE
// 					agencyBo.EditAgencyOnlineStatus(sess, agency)
// 					// 删除session
// 					conf.GetRedis().Del(ses)
// 					isDel = true
// 				}
// 				break
// 			}
// 		}
// 		if isDel {
// 			// 删除sessionListKey
// 			conf.GetRedis().HDel(model.GetAgencyListKey(), v)
// 		}
// 	}
// }
