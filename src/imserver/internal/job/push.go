package job

import (
	"context"
	"fmt"

	comet "baseGo/src/imserver/api/comet/grpc"
	pb "baseGo/src/imserver/api/logic/grpc"
	"baseGo/src/imserver/pkg/bytes"
	"encoding/json"
	log "fecho/golog"
	"strconv"
)

type MsgResp struct {
	RoomId      int32  `json:"roomId"`      // 房间id （群聊必填）
	SenderId    int32  `json:"senderId"`    // 发送id  （密聊必填）
	SendTime    int32  `json:"sendTime"`    // 发送时间 （非必填，取服务器时间）
	MsgType     int32  `json:"msgType"`     // 消息类型（1.文本；2.图片；3.视频；4.语音;）（必填）
	Msg         string `json:"msg"`         // 消息内容(json)  （必填）
	MsgId       int32  `json:"msgId"`       // 消息id
	ReceiverId  int32  `json:"receiverId"`  // 接收人
	SenderName  string `json:"senderName"`  // 发送人
	SenderHead  string `json:"senderHead"`  // 头像
	SenderLevel int64  `json:"senderLevel"` // 标签（1普通用户，2管理员，3群主）
	ReceiveType int32  `json:"receiveType"` // 消息类型（1私聊2群聊）
	Account     string `json:"account"`     // 会员账号
}

func (j *Job) push(ctx context.Context, pushMsg *pb.PushMsg) (err error) {
	info := MsgResp{
		RoomId:      pushMsg.RoomId,
		SenderId:    pushMsg.Mid, // 发送人id
		SendTime:    pushMsg.SendTime,
		MsgType:     pushMsg.MsgType,
		Msg:         string(pushMsg.Msg[:]),
		MsgId:       pushMsg.MsgId,
		SenderName:  pushMsg.SenderName,
		SenderHead:  pushMsg.SenderHead,
		SenderLevel: pushMsg.SenderLevel,
		ReceiverId:  pushMsg.ReceiverId,
		ReceiveType: pushMsg.ReceiveType,
	}
	msgByte, _ := json.Marshal(info)
	var keys []string
	keys = append(keys, fmt.Sprintf("%s-%d", pushMsg.Account, pushMsg.ReceiverId)) // 接收人key

	switch pushMsg.Type {
	case pb.PushMsg_PUSH:
		err = j.pushKeys(pushMsg.Operation, pushMsg.Server, keys, msgByte, pushMsg.ReceiverId)
		log.Info("job", "push", "PushMsg_PUSH send message success:", string(msgByte))
		return
	case pb.PushMsg_ROOM:
		err = j.getRoom(strconv.Itoa(int(pushMsg.RoomId))).Push(pushMsg.Operation, msgByte, pushMsg.LineId, pushMsg.AgencyId)
		log.Info("job", "push", "PushMsg_ROOM send message success:", string(msgByte))
		return
	case pb.PushMsg_BROADCAST:
		err = j.broadcast(pushMsg.Operation, msgByte, pushMsg.Speed, pushMsg.LineId, pushMsg.AgencyId)
		log.Info("job", "push", "PushMsg_BROADCAST send message success:", string(msgByte))
		return
	default:
		err = fmt.Errorf("no match push type: %s", pushMsg.Type)
	}
	return
}

// pushKeys push a message to a batch of subkeys.
func (j *Job) pushKeys(operation int32, serverID string, keys []string, body []byte, userId int32) (err error) {
	buf := bytes.NewWriterSize(len(body) + 64)
	p := &comet.Proto{
		Ver:    1,
		Op:     operation,
		Body:   body,
		UserId: int64(userId),
	}
	p.WriteTo(buf)
	p.Body = buf.Buffer()
	p.Op = comet.OpRaw
	var args = comet.PushMsgReq{
		Keys:    keys,
		ProtoOp: operation,
		Proto:   p,
	}
	comets := j.cometServers
	if c, ok := comets[serverID]; ok {
		if err = c.Push(&args); err != nil {
			log.Error("Job", "pushKeys", "c.Push(%v) serverID:%s error(%v)", nil, args, serverID, err)
		}
		log.Info("Job", "pushKeys", "pushKey:%s comets:%d", serverID, len(j.cometServers))
	}
	return
}

// broadcast broadcast a message to all.
func (j *Job) broadcast(operation int32, body []byte, speed int32, lineId, agencyId string) (err error) {
	buf := bytes.NewWriterSize(len(body) + 64)
	p := &comet.Proto{
		Ver:  1,
		Op:   operation,
		Body: body,
	}
	p.WriteTo(buf)
	p.Body = buf.Buffer()
	p.Op = comet.OpRaw
	p.LineId = lineId
	p.AgencyId = agencyId
	comets := j.cometServers
	fmt.Println("-----------注册到job的comets:", comets)
	if len(comets) > 0 {
		speed /= int32(len(comets))
	}
	var args = comet.BroadcastReq{
		ProtoOp: operation,
		Proto:   p,
		Speed:   speed,
	}
	for serverID, c := range comets {
		if err = c.Broadcast(&args); err != nil {
			log.Error("Job", "broadcast", "c.Broadcast(%v) serverID:%s error(%v)", nil, args, serverID, err)
		}
	}
	//log.Info( "Job", "broadcast", "broadcast comets:%d", len(comets))
	return
}

// broadcastRoomRawBytes broadcast aggregation messages to room.
func (j *Job) broadcastRoomRawBytes(roomID string, body []byte, lineId, agencyId string) (err error) {
	args := comet.BroadcastRoomReq{
		RoomID: roomID,
		Proto: &comet.Proto{
			Ver:      1,
			Op:       comet.OpRaw,
			Body:     body,
			LineId:   lineId,
			AgencyId: agencyId,
		},
	}
	comets := j.cometServers
	fmt.Println("-----------注册到job的comets:", comets)
	for serverID, c := range comets {
		if err = c.BroadcastRoom(&args); err != nil {
			log.Error("Job", "broadcastRoomRawBytes", "c.BroadcastRoom(%v) roomID:%s serverID:%s error(%v)", nil, args, roomID, serverID, err)
		}
	}
	//log.Info( "Job", "broadcastRoomRawBytes", "broadcastRoom comets:%d", len(comets))
	return
}
