package dao

import (
	"baseGo/src/fecho/echo"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	pb "baseGo/src/imserver/api/logic/grpc"
	"baseGo/src/imserver/internal/logic/conf"
	"baseGo/src/imserver/internal/logic/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"context"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/gogo/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/Shopify/sarama.v1"
)

var (
	ImOnlineMessage = new(bo.ImOnlineMessage)
	ImOfflineCount  = new(bo.ImOfflineCount)
	OffMessageMap   = make(map[string]int, 0)
	OffMessageCount = 0   // 缓存消息条数统计字段（条）
	TimeCount       = 0   // 缓存消息存储时间统计字段（秒）
	RunCount        = 500 // 触发极光推送缓存消息条数（条）
	RunTimeCount    = 10  // 触发极光推送缓存消息时间（秒）
)

// PushMsg push a message to databus. 修改 存储消息 单人
/**
* @ parameter senderId int 消息发送人id
* @ parameter sender string 消息发送人昵称
**/
func (d *Dao) PushMsg(op, senderId, receiverId, msgType int32, msg []byte, account string, receiveType int32) (*structs.MsgResp, error) {
	// 连接数据库
	sess := conf.GetXormSession()
	defer sess.Close()

	// 保存消息信息
	info := new(structs.ImOnlineMessage)
	info.Message = string(msg[:])
	info.SenderId = int(senderId)
	info.ReceiverId = int(receiverId)
	info.SendTime = utility.GetNowTimestamp()
	info.MessageType = int(msgType)
	_, err := ImOnlineMessage.AddOne(sess, info)
	if err != nil {
		golog.Error("Dao", "PushMsg", "error:", err)
		return nil, &echo.Err{Code: code.ADD_FAILED}
	}
	msgId := int32(info.Id)

	var keyStr = fmt.Sprintf("%s-%d", account, receiverId) // 接收人key
	session, err := d.SessionByKey(context.Background(), keyStr)
	if err != nil || (nil != session && (session.Server == "" || !session.Online)) {
		// todo 接收人不在线,离线推送

		// 查询该用户是否存在离线消息
		offlineMessageCount, err := ImOfflineCount.GetOfflineMessageCount(sess, receiverId, senderId, 0)
		offlineCounts := make([]*structs.OfflineMessageCount, 0)
		if err != nil {
			golog.Error("Dao", "PushMsg", "error:", err)
		}
		if len(offlineMessageCount) > 0 {
			// 存在离线消息数量，修改数量
			offlineCount := offlineMessageCount[0]
			offlineCount.MessageCount += 1
			if int(msgId) < offlineCount.OfflineMessageId {
				offlineCount.OfflineMessageId = int(msgId)
			}
			offlineCounts = append(offlineCounts, offlineCount)
			ImOfflineCount.UpdateOfflineCount(sess, offlineCounts)
		} else {
			// 不存在离线消息数量，新增离线消息数量
			offlineCount := new(structs.OfflineMessageCount)
			offlineCount.UserId = int(receiverId)
			offlineCount.SenderId = int(senderId)
			offlineCount.OfflineMessageId = int(msgId)
			offlineCount.MessageCount = 1
			offlineCounts = append(offlineCounts, offlineCount)
			n, err := ImOfflineCount.AddOfflineCount(sess, offlineCounts)
			if err != nil || n != 1 {
				golog.Error("Dao", "PushMsg", "error:", err)
			}
		}

		res := &structs.MsgResp{
			RoomId:     0,
			Mid:        int(senderId),
			SendTime:   utility.GetNowTimestamp(),
			MsgType:    int(msgType),
			Msg:        string(msg),
			MsgId:      int(msgId),
			ReceiverId: int(receiverId),
		}
		return res, nil
	}

	// 获取发送人信息

	pushMsg := &pb.PushMsg{
		Type:        pb.PushMsg_PUSH,
		Operation:   op,
		Server:      session.Server,
		Msg:         msg,
		MsgType:     msgType,
		MsgId:       msgId,
		Mid:         senderId,
		SenderName:  "",
		SenderHead:  "",
		ReceiverId:  receiverId,
		SendTime:    int32(utility.GetNowTimestamp()),
		ReceiveType: receiveType,
		Account:     account,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		logs.Error("Dao", "PushMsg", "error:", err)
		return nil, &echo.Err{Code: code.SEND_MEMBER_ERROR}
	}
	// m的key 是mid或者roodID 需要加密
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(receiverId),
		Topic: d.c.Kafka.Topic,
		//Topic: TOPIC,
		Value: sarama.ByteEncoder(b),
	}
	// todo 推送失败的离线消息暂时不做
	if _, _, err := d.kafkaPub.SendMessage(m); err != nil {
		fmt.Println("kafka消息发送失败：", err)
		logs.Error("Dao", "PushMsg", "PushMsg.send(push pushMsg:%v) error(%v)", nil, pushMsg, err)
		return nil, &echo.Err{Code: code.SEND_MEMBER_ERROR}
	}
	fmt.Println("kafka消息发送成功")
	//fmt.Println("kafka success:", &m)
	res := &structs.MsgResp{
		RoomId:     0,
		Mid:        int(senderId),
		SendTime:   utility.GetNowTimestamp(),
		MsgType:    int(msgType),
		Msg:        string(msg),
		MsgId:      int(msgId),
		ReceiverId: int(receiverId),
	}
	return res, nil
}

// BroadcastRoomMsg push a message to databus.
func (d *Dao) BroadcastRoomMsg(op, room, senderId, msgType int32, msg []byte, receiveType int32, lineId, agencyId string) (*structs.MsgResp, error) {
	fmt.Println("发送群聊消息-----------：", string(msg), msgType, room)
	// 获取数据库连接
	sess := conf.GetXormSession()
	defer sess.Close()

	// 开启事物
	sess.Begin()

	info := new(structs.ImOnlineMessage)
	info.Message = string(msg[:])
	info.SenderId = int(senderId)
	info.ReceiveRoomId = int(room)
	info.SendTime = utility.GetNowTimestamp()
	info.MessageType = int(msgType)
	_, err := ImOnlineMessage.AddOne(sess, info)
	if err != nil {
		golog.Error("Dao", "BroadcastRoomMsg", "error", err)
		sess.Rollback()
		//fmt.Println("数据库保存失败：", err)
		return nil, &echo.Err{Code: code.ADD_FAILED}
	}
	msgId := int32(info.Id)

	pushMsg := &pb.PushMsg{
		Type:        pb.PushMsg_ROOM,
		Operation:   op,
		RoomId:      room,
		Msg:         msg,
		MsgType:     msgType,
		MsgId:       msgId,
		Mid:         senderId,
		LineId:      lineId,
		AgencyId:    agencyId,
		SenderName:  "",
		SenderHead:  "",
		SenderLevel: 0,
		SendTime:    int32(utility.GetNowTimestamp()),
		ReceiveType: receiveType,
	}

	b, err := proto.Marshal(pushMsg)
	if err != nil {
		logs.Error("Dao", "BroadcastRoomMsg", "error", err)
		//fmt.Println("转protobuf协议失败")
		return nil, &echo.Err{Code: code.SEND_MEMBER_ERROR}
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(room),
		Topic: d.c.Kafka.Topic,
		//Topic: TOPIC,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err = d.kafkaPub.SendMessage(m); err != nil {
		fmt.Println("kafka消息发送失败：", err)
		sess.Rollback()
		logs.Error("Dao", "BroadcastRoomMsg", "error", err)
		return nil, &echo.Err{Code: code.SEND_MEMBER_ERROR}
	}
	fmt.Println("kafka消息发送成功")
	sess.Commit()

	// 处理离线会员消息逻辑
	go func() {
		// 从redis获取群信息
		redisClient := model.GetRedis().Get()
		// 获取群信息
		memberListResult, err := redis.Values(redisClient.Do("HGetAll", fmt.Sprintf("ChatRoom_%d", room)))
		session := conf.GetXormSession()
		if err != nil {
			golog.Error("Dao", "BroadcastRoomMsg", "error", err)
		}
		updateOfflineCounts := make([]*structs.OfflineMessageCount, 0)
		addOfflineCounts := make([]*structs.OfflineMessageCount, 0)
		var memberSessionKey string
		var memberId int

		for _, srm := range memberListResult {
			member := new(structs.ChannelObj)
			fResult := string(srm.(byte))
			json.Unmarshal([]byte(fResult), &member)

			memberSessionKey = fmt.Sprintf("%s_%d", member.Account, member.UserId)
			memberId = member.UserId
			memberSession, err := d.SessionByKey(context.Background(), memberSessionKey)
			if err != nil {
				golog.Error("Dao", "BroadcastRoomMsg", "error", err)
				sess.Rollback()
				//fmt.Println("获取session失败：", keyStr)
				continue
			}
			if err == nil || nil == memberSession || !memberSession.Online {
				// todo 接收人不在线,离线推送

				// 添加会员离线消息数量
				// 查询该用户是否存在离线消息
				offlineMessageCount, err := ImOfflineCount.GetOfflineMessageCount(session, int32(memberId), 0, room)
				if err != nil {
					sess.Rollback()
					golog.Error("Dao", "BroadcastRoomMsg", "error:", err)
					continue
				}
				// 组装离线会员
				if len(offlineMessageCount) > 0 {
					// 存在离线消息数量，修改数量
					offlineCount := offlineMessageCount[0]
					offlineCount.MessageCount += 1
					if int(msgId) < offlineCount.OfflineMessageId {
						offlineCount.OfflineMessageId = int(msgId)
					}
					updateOfflineCounts = append(updateOfflineCounts, offlineCount)
				} else {
					// 不存在离线消息数量，新增离线消息数量
					offlineCount := new(structs.OfflineMessageCount)
					offlineCount.UserId = memberId
					offlineCount.RoomId = int(room)
					offlineCount.OfflineMessageId = int(msgId)
					offlineCount.MessageCount = 1
					addOfflineCounts = append(addOfflineCounts, offlineCount)
				}
			}
		}
		if len(addOfflineCounts) > 0 {
			// 批量添加
			ImOfflineCount.AddOfflineCount(session, addOfflineCounts)
		}
		if len(updateOfflineCounts) > 0 {
			// 批量修改
			ImOfflineCount.UpdateOfflineCount(session, updateOfflineCounts)
		}
		session.Close()
	}()

	res := &structs.MsgResp{
		RoomId:   int(room),
		Mid:      int(senderId),
		SendTime: utility.GetNowTimestamp(),
		MsgType:  int(msgType),
		Msg:      string(msg),
		MsgId:    int(msgId),
	}
	return res, nil
}

// Push simulation data
func (d *Dao) BroadcastRoomSimulationMsg(key string, op, room, senderId, msgType, userId int32, msg []byte, roomType int) (int64, error) {
	return 0, nil
}

// BroadcastMsg push a message to databus.
func (d *Dao) BroadcastMsg(op, speed int32, msg []byte) (err error) {
	return
}

// 推送通知消息
func (d *Dao) SendNotificationMessage(op, receiverId, msgType, SendTime, pushCrowd int32, msg []byte) error {
	return nil
}

// 发送broadcast
func (d *Dao) sendBroadcastMsg(op int32, b []byte) error {
	return nil
}

func (d *Dao) PushRoomNoticyMsg(op, room, msgType, userId int32, msg []byte, roomType int) error {
	return nil
}
