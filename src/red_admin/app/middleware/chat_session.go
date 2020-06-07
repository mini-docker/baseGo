package middleware

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/red_admin/conf"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type SessionService struct {
}
type SessionCache struct {
	Sess []string `json:"sess"`
}
type idIndex struct {
	Id          int    `json:"id"`
	SiteIndexId string `json:"siteIndexId"`
	Name        string `json:"name"`
}

// sessionIdFull:包含了uuid&Id
var (
	KEEP_ONLINE_TIME = time.Duration(8 * time.Hour) // 8小时超长在线时长(对应前端保持在线功能)
	FIVE_HOURS       = time.Duration(5 * time.Hour) // 5小时(android 和 ios 加5小时所用)
	ONE_HOURS        = time.Duration(time.Hour)     // 默认1小时
)

// GetSession 根据sessionIdFull获取用户session信息
func (ss *SessionService) GetSession(sessionIdFull string) (*model.AdminSession, error) {
	// 获取session数据
	redisClient := conf.GetRedis().Get()
	defer redisClient.Close()
	// 获取session数据
	result, err := redisClient.Do("GET", sessionIdFull)
	if err == redis.Nil {
		return nil, err
	} else if err != nil {
		golog.Error("SessionService", "GetSession", "err:", err)
		return nil, err
	}
	if result.(string) == "" {
		return nil, errors.New("Not found session info")
	}
	var adminSession model.AdminSession
	err = json.Unmarshal([]byte(result.(string)), &adminSession)
	if err != nil {
		golog.Error("ChatSessionService", "GetSession", "err:", err)
		return nil, err
	}
	return &adminSession, nil
}

// SaveSession 存储session信息
// listKey 用来存储所有sessionIdFull的redis list的key
func (ss *SessionService) SaveSession(listKey string, session *model.AdminSession) error {
	if session == nil {
		return errors.New("Cannot store empty information")
	}

	//设置长时间未操作的掉线时间(1小时)
	timeNow := utility.GetNowTime()
	var timeOut time.Time
	var setTime time.Duration
	if session.IsKeepOnline {
		// 8小时掉线
		timeOut = timeNow.Add(KEEP_ONLINE_TIME)
		setTime = KEEP_ONLINE_TIME
	} else {
		timeOut = timeNow.Add(ONE_HOURS)
		setTime = ONE_HOURS
	}

	// session 过期时间赋值
	session.TimeOut = int(timeOut.Unix())

	// json序列化
	sessionBs, err := json.Marshal(session)
	if err != nil {
		golog.Error("SessionService", "SaveSession", "err:%v", err)
		return err
	}

	// 存入redis,sessionKey不过期
	setTime = -1
	redisClient := conf.GetRedis().Get()
	defer redisClient.Close()
	redisClient.Do("Set", session.SessionId, string(sessionBs), setTime)
	// 将sessionIdFull存入redis list
	se := new(SessionCache)
	// 根据 listkey 和 session.User.Id 获取 对应结果
	sestr, err := redisClient.Do("HGet", listKey, fmt.Sprint(session.User.Id))
	if err != nil && err != redis.Nil {
		golog.Error("ChatSessionService", "SaveSession", "err:", err)
		return err
	}
	json.Unmarshal([]byte(sestr.(string)), se)

	// 如果有结果的 SessionId 和匹配的session.SessionId相同 那么 就不把 东西存进 redis list
	if se == nil || se.Sess == nil || len(se.Sess) == 0 {
		se.Sess = []string{session.SessionId}
	} else {
		if !common.Contain(se.Sess, session.SessionId) {
			se.Sess = append(se.Sess, session.SessionId)
		}

	}
	ksjson, err := json.Marshal(se)
	if err != nil && err != redis.Nil {
		golog.Error("SessionService", "SaveSession", "err:", err)
		return err
	}
	redisClient.Do("HSet", listKey, fmt.Sprint(session.User.Id), ksjson)
	return nil
}

// DelSessionId 删除sessionId:(挤线和主动下线使用)
func (ss *SessionService) DelSessionId(listkey string, id int) error {
	redisClient := conf.GetRedis().Get()
	defer redisClient.Close()
	sessionstr, err := redisClient.Do("HGet", listkey, fmt.Sprint(id))
	if err != nil && err != redis.Nil {
		golog.Error("ChatSessionService", "DelMemberSessionId", "err:", err)
		return err
	}
	se := new(SessionCache)
	json.Unmarshal([]byte(sessionstr.(string)), se)
	for _, ses := range se.Sess {
		// 删除session数据
		_, err = redisClient.Do("Del", ses)
		if err != nil && err != redis.Nil {
			golog.Error("ChatSessionService", "DelMemberSessionId", "err:", err)
			return err
		}
	}
	redisClient.Do("HDel", listkey, fmt.Sprint(id))
	return nil
}
