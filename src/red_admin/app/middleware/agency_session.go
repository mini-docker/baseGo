package middleware

import (
	"baseGo/src/red_admin/conf"
	"encoding/json"
	"fmt"

	"baseGo/src/fecho/golog"

	"github.com/go-redis/redis"
)

type AdminSessionService struct {
}
type AdminSessionCache struct {
	Sess []string `json:"sess"`
}

// GetSession 根据sessionIdFull获取用户session信息
// func (ss *AdminSessionService) GetSession(sessionIdFull string) (*model.AdminSession, error) {
// 	// 获取session数据
// 	result, err := conf.GetRedis().Get(sessionIdFull).Result()
// 	if err == redis.Nil {
// 		return nil, err
// 	} else if err != nil {
// 		golog.Error("SessionService", "GetSession", "err:", err)
// 		return nil, err
// 	}
// 	if result == "" {
// 		return nil, errors.New("Not found session info")
// 	}
// 	var adminSession model.AdminSession
// 	err = json.Unmarshal([]byte(result), &adminSession)
// 	if err != nil {
// 		golog.Error("SessionService", "GetSession", "err:", err)
// 		return nil, err
// 	}
// 	return &adminSession, nil
// }

// // SaveSession 存储session信息
// // listKey 用来存储所有sessionIdFull的redis list的key
// func (ss *AdminSessionService) SaveSession(listKey string, session *model.AdminSession) error {
// 	if session == nil {
// 		return errors.New("Cannot store empty information")
// 	}

// 	//设置长时间未操作的掉线时间(1小时)
// 	timeNow := utility.GetNowTime()
// 	var timeOut time.Time
// 	var setTime time.Duration
// 	if session.IsKeepOnline {
// 		// 8小时掉线
// 		timeOut = timeNow.Add(KEEP_ONLINE_TIME)
// 		setTime = KEEP_ONLINE_TIME
// 	} else {
// 		timeOut = timeNow.Add(ONE_HOURS)
// 		setTime = KEEP_ONLINE_TIME
// 	}

// 	// session 过期时间赋值
// 	session.TimeOut = int(timeOut.Unix())

// 	// json序列化
// 	sessionBs, err := json.Marshal(session)
// 	if err != nil {
// 		golog.Error("SessionService", "SaveSession", "err:%v", err)
// 		return err
// 	}

// 	// 存入redis
// 	err = conf.GetRedis().Set(session.SessionId, string(sessionBs), setTime).Err()

// 	// 将sessionIdFull存入redis list
// 	se := new(SessionCache)
// 	sestr, err := conf.GetRedis().HGet(listKey, fmt.Sprint(session.User.Id)).Result()
// 	if err != nil && err != redis.Nil {
// 		golog.Error("SessionService", "SaveSession", "err:", err)
// 		return err
// 	}
// 	json.Unmarshal([]byte(sestr), se)

// 	if se == nil || se.Sess == nil || len(se.Sess) == 0 {
// 		se.Sess = []string{session.SessionId}
// 	} else {
// 		if !common.Contain(se.Sess, session.SessionId) {
// 			se.Sess = append(se.Sess, session.SessionId)
// 		}

// 	}
// 	ksjson, err := json.Marshal(se)
// 	if err != nil && err != redis.Nil {
// 		golog.Error("SessionService", "SaveSession", "err:", err)
// 		return err
// 	}
// 	conf.GetRedis().HSet(listKey, fmt.Sprint(session.User.Id), ksjson)
// 	return nil
// }

// // DelSessionId 删除sessionId:(挤线和主动下线使用)
// func (ss *AdminSessionService) DelSessionId(listkey string, id int) error {
// 	sessionstr, err := conf.GetRedis().HGet(listkey, fmt.Sprint(id)).Result()
// 	if err != nil && err != redis.Nil {
// 		golog.Error("SessionService", "DelMemberSessionId", "err:", err)
// 		return err
// 	}
// 	se := new(SessionCache)
// 	json.Unmarshal([]byte(sessionstr), se)
// 	for _, ses := range se.Sess {
// 		// 删除session数据
// 		_, err = conf.GetRedis().Del(ses).Result()
// 		if err != nil && err != redis.Nil {
// 			golog.Error("SessionService", "DelMemberSessionId", "err:", err)
// 			return err
// 		}
// 	}

// 	return nil
// }

func (ss *AdminSessionService) DelSessionId(listkey string, id int) error {
	redisClient := conf.GetRedis().Get()
	defer redisClient.Close()
	sessionstr, err := redisClient.Do("HGet", listkey, fmt.Sprint(id))
	if err != nil && err != redis.Nil {
		golog.Error("SessionService", "DelMemberSessionId", "err:", err)
		return err
	}
	se := new(SessionCache)
	json.Unmarshal([]byte(sessionstr.(string)), se)
	for _, ses := range se.Sess {
		// 删除session数据
		_, err = redisClient.Do("Del", ses)
		if err != nil && err != redis.Nil {
			golog.Error("SessionService", "DelMemberSessionId", "err:", err)
			return err
		}
	}

	return nil
}
