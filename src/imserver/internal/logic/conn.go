package logic

import (
	log "baseGo/src/fecho/golog"
	registry_module "baseGo/src/fecho/registry/registry-module"
	comet "baseGo/src/imserver/api/comet/grpc"
	"baseGo/src/imserver/internal/logic/conf"
	"baseGo/src/imserver/internal/logic/model"
	"baseGo/src/model/bo"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Connect connected a conn.
func (l *Logic) Connect(c context.Context, server, cookie string, token []byte) (mid int64, key, roomID string, accepts []int32, hb int64, err error) {
	var params struct {
		Mid      int64   `json:"mid"`
		Key      string  `json:"key"`
		RoomID   string  `json:"room_id"`
		Platform string  `json:"platform"`
		Accepts  []int32 `json:"accepts"`
	}
	if err = json.Unmarshal(token, &params); err != nil {
		log.Error("Logic", "Connect", "error:", err)
		return
	}
	mid = params.Mid
	roomID = params.RoomID
	accepts = params.Accepts
	hb = int64(l.c.Node.Heartbeat) * int64(l.c.Node.HeartbeatMax)
	if key = params.Key; key != "" {
		logicSession := &model.LogicSession{}
		//当作Session存储
		logicSessionBytes, err := json.Marshal(logicSession)
		if err = l.dao.AddMapping(c, mid, key, string(logicSessionBytes)); err != nil {
			log.Error("Logic", "Connect", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
		}
	}
	//log.Info( "Logic", "Connect", "conn connected key:%s server:%s mid:%d token:%s", key, server, mid, token)
	return
}

// Disconnect disconnect a conn.
func (l *Logic) Disconnect(c context.Context, mid int64, key, server string) (has bool, err error) {

	if sess, err := l.dao.SessionByKey(c, key); err != nil {
		log.Error("Logic", "Disconnect", "l.dao.SessionByKey(%d,%s) error(%v)", nil, mid, key, server)

	} else {
		strs := strings.Split(key, "-")
		sess.Online = false
		logicSessionBytes, err := json.Marshal(sess)
		if err = l.dao.AddMapping(c, mid, key, string(logicSessionBytes)); err != nil {
			log.Error("Logic", "Disconnect", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
		}

		// 查询所有房间信息
		db := conf.GetXormSession()
		var ids []int
		for _, v := range sess.Rooms {
			ids = append(ids, int(v))
		}
		fmt.Println("ids-------------:", ids)
		res, err := RoomBo.FindRoomsByIds(db, ids)
		if err != nil {
			log.Error("Logic", "Disconnect", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)

		} else {
			if len(strs) > 2 {
				for _, room := range res {
					err = l.outRoomCache(c, fmt.Sprint(room.Id), mid, false, strs[0], false)
					if err != nil {
						log.Error("Logic", "Disconnect", "error(%v)", nil, err)
					}

				}
			}

		}

	}
	// 修改会员在线状态
	sess := conf.GetXormSession()
	new(bo.User).EditUserOnlineStatu(sess, mid, 2)
	sess.Close()
	return
}

// Heartbeat heartbeat a conn.
func (l *Logic) Heartbeat(c context.Context, mid int64, key, server string) (err error) {
	_, err = l.dao.ExpireMapping(c, mid, key)
	if err != nil {
		log.Error("Logic", "Heartbeat", "l.dao.ExpireMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
		return
	}

	//当作Session存储
	//logicSessionBytes, err := json.Marshal(logicSession)
	//if err != nil {
	//	log.Error( "Logic", "Heartbeat", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
	//	return
	//}
	//if !has {
	//	if err = l.dao.AddMapping(c, mid, key, string(logicSessionBytes)); err != nil {
	//		log.Error( "Logic", "Heartbeat", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
	//		return
	//	}
	//}
	//log.Info( "Logic", "Heartbeat", "conn heartbeat key:%s server:%s mid:%d", key, server, mid)
	return
}

// RenewOnline renew a server online.
func (l *Logic) RenewOnline(c context.Context, server string, roomCount map[string]int32) (map[string]int32, error) {
	online := &model.Online{
		Server:    server,
		RoomCount: roomCount,
		Updated:   time.Now().Unix(),
	}
	if err := l.dao.AddServerOnline(context.Background(), server, online); err != nil {
		return nil, err
	}
	return l.roomCount, nil
}

// Receive receive a message.
func (l *Logic) Receive(c context.Context, mid int64, proto *comet.Proto) (err error) {
	//log.Info( "Logic", "Receive", "receive mid:%d message:%+v", mid, proto)
	//fmt.Println("接收到消息：", mid, &proto)
	var message struct {
		OpType  int    `json:"op_type"`  // 操作类型（1.发送消息；2.删除消息）
		MsgType int    `json:"msg_type"` // 消息类型（1.普通消息；2.通知消息；3.系统消息；4.账务通知；5.图片消息；6.视频消息；7.语音消息）
		Msg     string `json:"msg"`      // 消息内容
		RoomId  int    `json:"room_id"`  // 接收房间id
		Mid     int    `json:"mid"`      // 接收用户id
		Keys    string `json:"keys"`     // 接收keys
	}
	if err = json.Unmarshal(proto.Body, &message); err != nil {
		//fmt.Println("==============err:", err.Error())
		log.Error("Logic", "Receive", "error:", err)
		return
	}
	//fmt.Println(message)
	return
}

// CashSession
func (l *Logic) CacheSession(c context.Context, cookie string, token []byte, userId int, rooms []int64) (mid int64, key, roomID string, accepts []int32, hb int64, server string, err error) {
	var params struct {
		Mid      int64   `json:"mid"`
		Key      string  `json:"key"`
		RoomID   string  `json:"room_id"`
		Platform string  `json:"platform"`
		Accepts  []int32 `json:"accepts"`
	}

	if err = json.Unmarshal(token, &params); err != nil {
		log.Error("Logic", "CacheSession", "json.Unmarshal(%s) error(%v)", nil, token, err)
		return
	}
	mid = params.Mid
	roomID = params.RoomID
	accepts = params.Accepts
	hb = int64(l.c.Node.Heartbeat) * int64(l.c.Node.HeartbeatMax)
	if key = params.Key; key == "" {
		key = fmt.Sprint(params.Mid)
	}

	switch params.Platform {
	case "ios", "android":
		server = registry_module.GetCometTcpUrl()
	case "wap", "pc", "web":
		server = registry_module.GetCometWsUrl()
	}
	logicSession := &model.LogicSession{
		Rooms:  rooms,
		Online: false,
	}
	//当作Session存储
	logicSessionBytes, err := json.Marshal(logicSession)
	if err != nil {
		log.Error("Logic", "CacheSession", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
		return
	}
	if err = l.dao.AddMapping(c, mid, key, string(logicSessionBytes)); err != nil {
		return
	}

	//log.Info( "Logic", "CacheSession", "conn connected key:%s server:%s mid:%d token:%s", key, server, mid, token)
	return
}

func (l *Logic) EstablishConn(c context.Context, userId int64, server string, sessionKey string, lineId, agencyId string) (mid int64, key, roomID string, hb int64, has bool, roomids []string, err error) {
	key = lineId + "-" + agencyId + "-" + fmt.Sprint(userId)
	session, err := l.dao.SessionByKey(c, key)
	if err != nil {
		log.Error("Logic", "EstablishConn", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, server, err)
		return
	}
	if session == nil || err == nil {
		return
	}
	has = true
	mid = userId
	//当作Session存储
	session.Server = server
	session.Online = true
	if session.Rooms != nil && len(session.Rooms) > 0 {
		for _, rid := range session.Rooms {
			roomids = append(roomids, fmt.Sprint(rid))

			err = l.joinRoomCache(c, fmt.Sprint(rid), userId, true, lineId)
			if err != nil {
				log.Error("Logic", "EstablishConn", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
			}
		}
	}
	hb = int64(l.c.Node.Heartbeat) * int64(l.c.Node.HeartbeatMax)
	logicSessionBytes, err := json.Marshal(session)
	if err != nil {
		log.Error("Logic", "EstablishConn", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
		return
	}
	model.GetRedis().Get().Do("Set", key, string(logicSessionBytes), 0)
	//fmt.Println(cmd.Err())
	if err = l.dao.AddMapping(c, userId, key, string(logicSessionBytes)); err != nil {
		log.Error("Logic", "EstablishConn", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, mid, key, server, err)
		return
	}
	// 修改会员在线状态
	sess := conf.GetXormSession()
	new(bo.User).EditUserOnlineStatu(sess, userId, 1)
	sess.Close()
	//log.Info( "Logic", "EstablishConn", "conn connected key:%s server:%s mid:%d ", key, server, userId)
	return
}

// CometGetJoinRoom
func (l *Logic) CometGetJoinRoom(c context.Context, userId int64, server string) (string, bool, error) {

	key := fmt.Sprint(userId)
	session, err := l.dao.SessionByKey(c, key)
	if err != nil {
		log.Error("Logic", "CometGetJoinRoom", "l.dao.AddMapping(%d,%s) error(%v)", nil, userId, key, err)
	}
	if session != nil && session.Server != "" {
		session.Server = server
	}

	//当作Session存储
	logicSessionBytes, err := json.Marshal(session)
	if err != nil {
		log.Error("Logic", "CometGetJoinRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
	}
	if err = l.dao.AddMapping(c, userId, key, string(logicSessionBytes)); err != nil {
		log.Error("Logic", "CometGetJoinRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
	}

	//log.Info( "Logic", "CometGetJoinRoom", "conn connected key:%s server:%s mid:%d", key, session.Server, userId)
	if session.RoomId == 0 {
		return "", false, err
	}
	return fmt.Sprint(session.RoomId), true, err
}

// DelRoom
func (l *Logic) DelRoom(c context.Context, roomId int) error {
	client := l.GetCometClient()
	if nil == client {
		return nil
	}
	_, err := client.CancelRoom(context.Background(), &comet.CancelRoomReq{
		RoomID: fmt.Sprint(roomId),
	})
	if err != nil {
		log.Error("Logic", "DelRoom", "did not connect: %v", err)
	}
	return err
}

// RoomCount
func (l *Logic) RoomCount(c context.Context, lineId, agencyId string) (*comet.RoomUserCountResp, error) {
	client := l.GetCometClient()
	if nil == client {
		return nil, nil
	}
	//defer conn.Close()

	//client := comet.NewCometClient(conn)
	//resp, err := client.RoomUserCount(context.Background(), &comet.RoomUserCountReq{

	//})
	var resp = new(comet.RoomUserCountResp)
	channelObjs, err := l.dao.GetGroupMemberAll(c, lineId)
	if err != nil {
		log.Error("Logic", "RoomCount", " %v", err)
	}
	for roomId, channelObj := range channelObjs {
		objInt, err := strconv.Atoi(channelObj)
		if err != nil {
			log.Error("Logic", "RoomCount", " %v", err)
		}
		resp.RoomUserCount = append(resp.RoomUserCount, &comet.RoomUserCount{
			RoomId: roomId,
			Count:  int64(objInt),
		})
	}
	resp.Result = true
	return resp, err
}

// VeeGebruikerSessieUit 判断是否在线 不在线 删除用户session
func (l *Logic) VeeGebruikerSessieUit(userId int64) (*comet.RoomUserCountResp, error) {
	client := l.GetCometClient()
	if nil == client {
		return nil, nil
	}
	resp, err := client.VeeGebruikerSessieUit(context.Background(), &comet.RoomUserCountReq{
		Proto: &comet.Proto{
			UserId: userId,
		},
	})

	return resp, err
}

// RoomCount
func (l *Logic) RoomInfo(c context.Context, roomId string, lineId, agencyId string) (*comet.RoomInfoResp, error) {
	client := l.GetCometClient()
	if nil == client {
		return nil, nil
	}
	//defer conn.Close()

	//client := comet.NewCometClient(conn)
	//resp, err := client.RoomInfolRoom(context.Background(), &comet.RoomInfoReq{
	//	RoomID: roomId,
	//	Proto: &comet.Proto{
	//	},
	//})
	var resp = new(comet.RoomInfoResp)
	channelObjs, err := l.dao.GetMemberCount(c, lineId, roomId)
	for key, _ := range channelObjs {
		objInt, err := strconv.Atoi(key)
		if err != nil {
			log.Error("Logic", "RoomInfo", " %v", err)
		}
		resp.Mid = append(resp.Mid, int64(objInt))
	}

	return resp, err
}

// IntiveRoom
func (l *Logic) IntiveRoom(c context.Context, userKeys []string, roomId int, noticeType int, senderId int, lineId string) error {
	client := l.GetCometClient()
	if nil == client {
		return nil
	}
	for _, key := range userKeys {
		strs := strings.Split(key, "-")
		if len(strs) != 2 {
			continue
		}
		userId, _ := strconv.ParseInt(strs[1], 10, 64)
		// 获取session
		session, err := l.dao.SessionByKey(c, key)
		if err != nil {
			log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s) error(%v)", nil, userId, key, err)
			return err
		}
		// session存在
		if nil != session {
			// 判断是否在线
			if session.Online == true {
				// 在线绑定comet内存缓存
				_, err = client.JoinRoom(context.Background(), &comet.JoinRoomReq{
					RoomID: fmt.Sprint(roomId),
					Proto: &comet.Proto{
						UserId:     userId,
						SessionKey: key,
					},
				})
			}
			// 判断是否已经监听该房间信息
			has := false
			for _, v := range session.Rooms {
				if int64(roomId) == v {
					has = true
				}
			}
			if !has {
				session.Rooms = append(session.Rooms, int64(roomId))
			}
			//更新Session存储
			logicSessionBytes, err := json.Marshal(session)
			if err != nil {
				log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				continue
			}
			if err = l.dao.AddMapping(c, userId, key, string(logicSessionBytes)); err != nil {
				log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				continue
			}
			err = l.joinRoomCache(c, fmt.Sprint(roomId), userId, session.Online, lineId)
			if err != nil {
				log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				continue
			}
		} else {
			// session不存在
			err = l.joinRoomCache(c, fmt.Sprint(roomId), userId, false, lineId)
			if err != nil {
				log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				continue
			}
		}
	}
	return nil
}

// kickRoom
func (l *Logic) KickRoom(c context.Context, userKeys []string, roomId int, noticeType int, senderId int, lineId string) error {
	client := l.GetCometClient()
	if nil == client {
		return nil
	}
	for _, key := range userKeys {
		strs := strings.Split(key, "-")
		if len(strs) != 2 {
			continue
		}
		userId, _ := strconv.ParseInt(strs[1], 10, 64)

		session, err := l.dao.SessionByKey(c, key)
		if err != nil {
			log.Error("Logic", "KickRoom", "l.dao.AddMapping(%d,%s) error(%v)", nil, userId, key, err)
			return err
		}
		redisNil := false
		if session != nil && session.Server != "" {
			if session.Online == false {
				return nil
			}
			session.RoomId = roomId
			roomIds := []int64{}
			// 删除监听该房间信息
			for _, v := range session.Rooms {
				if int64(roomId) != v {
					roomIds = append(roomIds, v)
				}
			}
			session.Rooms = roomIds
		}
		online := false
		if redisNil == false && session != nil {
			if session.Online == true {
				online = true
				_, err = client.LeaveRoom(context.Background(), &comet.LeaveRoomReq{
					RoomID: fmt.Sprint(roomId),
					Proto: &comet.Proto{
						UserId:     userId,
						SessionKey: key,
					},
				})
			}
			//当作Session存储
			logicSessionBytes, err := json.Marshal(session)
			if err != nil {
				log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				return err
			}
			if err = l.dao.AddMapping(c, userId, key, string(logicSessionBytes)); err != nil {
				log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				return err
			}
		}
		/**
		roommembercount
		*/
		err = l.kickRoomCache(c, fmt.Sprint(roomId), userId, online, lineId)
		if err != nil {
			log.Error("Logic", "IntiveRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
			return err
		}
	}
	return nil
}

func (l *Logic) joinRoomCache(c context.Context, roomId string, userId int64, online bool, lineId string) error {
	channelObj := &model.ChannelObj{
		Online: online,
		UserId: userId,
	}
	channelObjstr, err := json.Marshal(channelObj)
	if err != nil {
		log.Error("Logic", "joinRoomCache", "json.Marshal(%v) error(%v)", nil, channelObj, err)
		return err
	}
	if err := l.dao.UpdateGroupMemberCount(c, lineId, fmt.Sprint(roomId), fmt.Sprint(userId), online, string(channelObjstr), online, false); err != nil {
		return err
	}
	return nil
}
func (l *Logic) outRoomCache(c context.Context, roomId string, userId int64, isDel bool, lineId string, online bool) error {
	channelObj := &model.ChannelObj{
		Online: online,
		UserId: userId,
	}
	channelObjstr, err := json.Marshal(channelObj)
	if err != nil {
		log.Error("Logic", "outRoomCache", "json.Marshal(%v) error(%v)", nil, channelObj, err)
		return err
	}
	if err := l.dao.UpdateGroupMemberCount(c, lineId, fmt.Sprint(roomId), fmt.Sprint(userId), false, string(channelObjstr), online, isDel); err != nil {
		return err
	}
	return nil
}
func (l *Logic) kickRoomCache(c context.Context, roomId string, userId int64, online bool, lineId string) error {

	if err := l.dao.KickGroupMemberCount(c, lineId, fmt.Sprint(roomId), fmt.Sprint(userId), online); err != nil {
		return err
	}
	return nil
}

// JoinRoom
func (l *Logic) JoinRoom(c context.Context, lineId, agencyId string, userId int64, roomId int, noticeType int, senderId int, content string, roomType int) error {
	key := lineId + "-" + agencyId + "-" + fmt.Sprint(userId)
	session, err := l.dao.SessionByKey(c, key)
	if err != nil {
		return err
	}

	if session != nil && session.Server != "" {
		log.Info("Dao", "BroadcastRoomMsg", "进入聊天室会员session：", session.Server, session.Online)
		if session.Online == false {
			return errors.New("user not online")
		} else {
			if session.RoomId != 0 {
				// 从之前的房间退出
				err = l.outRoomCache(c, fmt.Sprint(session.RoomId), userId, true, lineId, session.Online)
				if err != nil {
					log.Error("Logic", "JoinRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
					return err
				}
			}
			clients := l.GetCometClinets()
			if len(clients) == 0 {
				log.Error("Dao", "JoinRoom", "获取comet连接失败：", err, "")
				return nil
			}
			for _, client := range clients {
				req, err := client.JoinRoom(context.Background(), &comet.JoinRoomReq{
					RoomID: fmt.Sprint(roomId),
					Proto: &comet.Proto{
						LineId:     lineId,
						AgencyId:   agencyId,
						UserId:     userId,
						SessionKey: lineId + "-" + agencyId + "-" + fmt.Sprint(userId),
					},
				})
				log.Info("conn", "JoinRoom", "进入聊天室返回:", req, err)
			}
			session.RoomId = roomId
			//当作Session存储
			logicSessionBytes, err := json.Marshal(session)
			if err != nil {
				log.Error("Logic", "JoinRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				return err
			}
			if err = l.dao.AddMapping(c, userId, key, string(logicSessionBytes)); err != nil {
				log.Error("Logic", "JoinRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				return err
			}
			/**
			roommembercount
			*/
			err = l.joinRoomCache(c, fmt.Sprint(roomId), userId, true, lineId)
			if err != nil {
				log.Error("Logic", "JoinRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
				return err
			}

		}
	}

	//log.Info("Logic", "JoinRoom", "conn connected key:%s server:%s mid:%d", key, session.Server, userId)
	return err
}

// OutRoom
func (l *Logic) OutRoom(c context.Context, lineId, agencyId string, userId int64, roomId int, noticeType int, roomType int) error {
	key := lineId + "-" + agencyId + "-" + fmt.Sprint(userId)
	session, err := l.dao.SessionByKey(c, key)
	if err != nil {
		log.Error("Logic", "OutRoom", "l.dao.AddMapping(%d,%s) error(%v)", nil, userId, key, err)
	}
	session.RoomId = 0
	err = l.outRoomCache(c, fmt.Sprint(roomId), userId, true, lineId, session.Online)
	if err != nil {
		return err
	}
	//当作Session存储
	logicSessionBytes, err := json.Marshal(session)
	if err != nil {
		log.Error("Logic", "OutRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
	}
	if err = l.dao.AddMapping(c, userId, key, string(logicSessionBytes)); err != nil {
		log.Error("Logic", "OutRoom", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
	}
	clients := l.GetCometClinets()
	if len(clients) == 0 {
		log.Error("Dao", "JoinRoom", "获取comet连接失败：", err, "")
		return nil
	}
	for _, client := range clients {
		_, err = client.LeaveRoom(context.Background(), &comet.LeaveRoomReq{
			RoomID: fmt.Sprint(roomId),
			Proto: &comet.Proto{
				LineId:     lineId,
				AgencyId:   agencyId,
				SessionKey: lineId + "-" + agencyId + "-" + fmt.Sprint(userId),
				UserId:     userId,
			},
		})
	}
	if err != nil {
		log.Error("Logic", "OutRoom", "did not connect: %v", err)
	}

	//log.Info("Logic", "OutRoom", "conn connected key:%s server:%s mid:%d", key, session.Server, userId)
	return err
}
