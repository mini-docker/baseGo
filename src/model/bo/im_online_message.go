package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type ImOnlineMessage struct{}

type UserInfos struct {
	Id         int    // id
	SenderName string // 昵称
	SenderHead string // 头像
	SenderRole int    // 角色（1.成员，2.管理员，3.群主）
}

// 添加
func (ImOnlineMessage) AddOne(sess *xorm.Session, Info *structs.ImOnlineMessage) (int64, error) {
	return sess.Insert(Info)
}

// 查询会员历史消息列表
func (ImOnlineMessage) GetSenderHisList(sess *xorm.Session, userId, receiverId int, roomIds string) (int64, []*structs.ImOnlineMessage, error) {
	data := make([]*structs.ImOnlineMessage, 0)
	sess.Where("delete_time = 0 and (receiver_id = ? or find_in_set(receive_room_id,?) > 0 or (receiver_id = 0 and receive_room_id = 0 and sender_id = ? ) or (receiver_id = 0 and receive_room_id = 0 and sender_id = 0))", userId, roomIds, userId)
	if receiverId != 0 {
		// 会员查询全部离线消息
		sess.Where("id >= ?", receiverId)
	} else {
		return 0, nil, nil
	}
	count, err := sess.FindAndCount(&data)
	return count, data, err
}
