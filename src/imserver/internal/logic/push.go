package logic

import (
	log "baseGo/src/fecho/golog"
	comet "baseGo/src/imserver/api/comet/grpc"
	"baseGo/src/model/structs"
	"context"
	"encoding/json"
	"fmt"
)

// PushMids push a message by mid.
func (l *Logic) PushMid(op, senderId, receiverId, msgType int32, msg []byte, account string, receiveType int32) (info *structs.MsgResp, err error) {
	resp, err := l.dao.PushMsg(op, senderId, receiverId, msgType, msg, account, receiveType)
	if err != nil {
		log.Error("Logic", "PushMid", "err:", err)
	}
	return resp, err
}

// PushRoom push a message by room.
func (l *Logic) PushRoom(op, room, senderId, msgType int32, msg []byte, ReceiveType int32, lineId, agencyId string) (info *structs.MsgResp, err error) {
	return l.dao.BroadcastRoomMsg(op, room, senderId, msgType, msg, ReceiveType, lineId, agencyId)
}

// PushRoom push a message by room.
func (l *Logic) PushRoomSimulation(key string, op, room, senderId, msgType, userId int32, msg []byte, roomType int) (n int64, err error) {
	return l.dao.BroadcastRoomSimulationMsg(key, op, room, senderId, msgType, userId, msg, roomType)
}

// PushAll push a message to all.
func (l *Logic) PushAll(op, speed int32, msg []byte) (err error) {
	return l.dao.BroadcastMsg(op, speed, msg)
}

// PushAll push a message to all.
func (l *Logic) SendNotificationMessage(op, receiverId, msgType, SendTime, speed, pushCrowd int32, msg []byte) (err error) {
	return l.dao.SendNotificationMessage(op, receiverId, msgType, SendTime, pushCrowd, msg)
}
func (l *Logic) kick(userId, roomId int64) {
	key := fmt.Sprint(userId)
	c := context.Background()
	session, err := l.dao.SessionByKey(c, key)
	if err != nil {
		log.Error("Logic", "kick", "l.dao.AddMapping(%d,%s) error(%v)", nil, userId, key, err)
	}
	if session != nil && session.Server != "" {
		session.RoomId = 0
		var roomint []int64
		for _, sessroom := range session.Rooms {
			if sessroom != int64(roomId) {
				roomint = append(roomint, sessroom)
			}
		}
		session.Rooms = roomint
	}

	//当作Session存储
	logicSessionBytes, err := json.Marshal(session)
	if err != nil {
		log.Error("Logic", "kick", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
	}
	if err = l.dao.AddMapping(c, userId, key, string(logicSessionBytes)); err != nil {
		log.Error("Logic", "kick", "l.dao.AddMapping(%d,%s,%s) error(%v)", nil, userId, key, session.Server, err)
	}
	client := l.GetCometClient()
	if nil == client {
		return
	}
	_, err = client.LeaveRoom(c, &comet.LeaveRoomReq{
		RoomID: fmt.Sprint(roomId),
		Proto: &comet.Proto{
			SessionKey: fmt.Sprint(userId),
			UserId:     userId,
		},
	})
	if err != nil {
		log.Error("Logic", "kick", "did not connect: %v", err)
	}
}
