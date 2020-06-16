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
	fmt.Println(sessionIdFull, "sessionIdFull")
	// if err == redis.Nil {
	// 	return nil, err
	// } else
	if err != nil {
		golog.Error("SessionService", "GetSession", "err:", err)
		return nil, err
	}
	fmt.Println(result, "result")
	fResult := string(result.([]byte))
	fmt.Println(fResult, "fResult")
	if fResult == "" {
		return nil, errors.New("Not found session info")
	}
	var adminSession model.AdminSession
	err = json.Unmarshal([]byte(fResult), &adminSession)
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
	fmt.Println(sessionBs, "sessionBs") //对象
	// 存入redis,sessionKey不过期
	setTime = -1
	fmt.Println(setTime, "setTime")
	redisClient := conf.GetRedis().Get()
	defer redisClient.Close()
	// 设置 SessionId 为 对象转为的字符串
	fmt.Println(string(sessionBs), "stringBs")
	fmt.Println(listKey, "listKey")
	// redisClient.Do("Set", session.SessionId, string(sessionBs), setTime)

	redisClient.Do("Set", session.SessionId, string(sessionBs))
	redisClient.Do("expire", session.SessionId, setTime)

	resultf, err := redisClient.Do("Get", session.SessionId)
	if err != nil {
		fmt.Println(err, "err")
	}
	fmt.Println(resultf, "resultf")
	// 将sessionIdFull存入redis list
	se := new(SessionCache)
	// 根据 listkey 和 session.User.Id 获取 对应结果
	sestr, err := redisClient.Do("HGet", listKey, fmt.Sprint(session.User.Id))
	fmt.Println(fmt.Sprint(session.User.Id), "HGetId", listKey)
	if err != nil {
		golog.Error("ChatSessionService", "SaveSession", "err:", err)
		return err
	}
	fmt.Println(sestr, "sestr233333")
	// fmt.Println(sestr.(map[string]string), "bytesss")
	if sestr != nil {

		fResults := string(sestr.([]byte))
		fmt.Println(fResults, "fResults")

		var dat SessionCache
		if err := json.Unmarshal([]byte(fResults), &dat); err == nil {
			fmt.Println("==============json str 转map=======================")
			fmt.Println(dat, "fResults_dat")
			// fmt.Println(dat["sess"])
		}
		se = &dat
		// json.Unmarshal([]byte(sestr.(string)), se)

		fmt.Println(se, "sesese")

		// 如果有结果的 SessionId 和匹配的session.SessionId相同 那么 就不把 东西存进 redis list
		if se == nil || se.Sess == nil || len(se.Sess) == 0 {
			se.Sess = []string{session.SessionId}
		} else {
			if !common.Contain(se.Sess, session.SessionId) {
				se.Sess = append(se.Sess, session.SessionId)
			}

		}
		ksjson, err := json.Marshal(se)
		if err != nil {
			golog.Error("SessionService", "SaveSession", "err:", err)
			return err
		}
		redisClient.Do("HSet", listKey, fmt.Sprint(session.User.Id), ksjson)
	} else {
		se.Sess = []string{session.SessionId}
		ksjson, err := json.Marshal(se)
		if err != nil {
			golog.Error("SessionService", "SaveSession", "err:", err)
			return err
		}
		fmt.Println(ksjson, "ksjsonksjson", fmt.Sprint(session.User.Id), listKey)
		redisClient.Do("HSet", listKey, fmt.Sprint(session.User.Id), ksjson)
	}
	// else {

	// 	sestr = ""
	// 	json.Unmarshal([]byte(sestr.(string)), se)

	// 	// 如果有结果的 SessionId 和匹配的session.SessionId相同 那么 就不把 东西存进 redis list
	// 	if se == nil || se.Sess == nil || len(se.Sess) == 0 {
	// 		se.Sess = []string{session.SessionId}
	// 	} else {
	// 		if !common.Contain(se.Sess, session.SessionId) {
	// 			se.Sess = append(se.Sess, session.SessionId)
	// 		}

	// 	}
	// 	fmt.Println(se, ".....se")
	// 	ksjson, err := json.Marshal(se)
	// 	if err != nil {
	// 		golog.Error("SessionService", "SaveSession", "err:", err)
	// 		return err
	// 	}
	// 	redisClient.Do("HSet", listKey, fmt.Sprint(session.User.Id), ksjson)

	// 	// 无返回 需要判断 默认账号 将session信息设置进去
	// 	// redisClient.Do("HSet", listKey, fmt.Sprint(session.User.Id), )
	// }

	return nil
}

// DelSessionId 删除sessionId:(挤线和主动下线使用)
func (ss *SessionService) DelSessionId(listkey string, id int) error {
	redisClient := conf.GetRedis().Get()
	defer redisClient.Close()
	sessionstr, err := redisClient.Do("HGet", listkey, fmt.Sprint(id))
	if err != nil {
		golog.Error("ChatSessionService", "DelMemberSessionId", "err:", err)
		return err
	}
	se := new(SessionCache)
	fmt.Println(se, "sesese")
	fmt.Println(sessionstr, "sessionstr")

	if sessionstr != nil {
		// json.Unmarshal([]byte(strsf.(string)), se)
		strsf := string(sessionstr.([]byte))
		var dat SessionCache
		if err := json.Unmarshal([]byte(strsf), &dat); err == nil {
			fmt.Println("==============json str 转map=======================")
			fmt.Println(dat, "datttt")
			// fmt.Println(dat["sess"])
		}
		se = &dat

		for _, ses := range se.Sess {
			// 删除session数据
			fmt.Println(ses, "sessesses")
			_, err = redisClient.Do("Del", ses)
			if err != nil {
				golog.Error("ChatSessionService", "DelMemberSessionId", "err:", err)
				return err
			}
		}
		fmt.Println(id, "ididid")
		redisClient.Do("HDel", listkey, fmt.Sprint(id))
	}

	return nil
}
