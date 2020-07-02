package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type ImOfflineCount struct{}

// 添加
func (ImOfflineCount) Add(sess *xorm.Session, Info ...*structs.OfflineMessageCount) (int64, error) {
	return sess.Insert(Info)
}

// 添加离线消息数量
func (ImOfflineCount) AddOfflineCount(sess *xorm.Session, offlineCount []*structs.OfflineMessageCount) (int64, error) {
	return sess.Table(structs.TABLE_OFFLINE_MESSAFE_COUNT).Insert(offlineCount)
}

// 获取用户离线消息数量
func (ImOfflineCount) GetOfflineMessageCount(sess *xorm.Session, userId int32, senderId int32, roomId int32) ([]*structs.OfflineMessageCount, error) {
	if roomId != 0 {
		sess.Where("room_id = ?", roomId)
	}
	if userId != 0 {
		sess.Where("user_id = ?", userId)
	}
	if senderId != 0 {
		sess.Where("sender_id = ?", senderId)
	}
	messgaeCount := make([]*structs.OfflineMessageCount, 0)
	err := sess.Table(structs.TABLE_OFFLINE_MESSAFE_COUNT).Find(&messgaeCount)
	return messgaeCount, err
}

// 修改离线消息数量
func (ImOfflineCount) UpdateOfflineCount(sess *xorm.Session, offlineCount []*structs.OfflineMessageCount) {
	for _, v := range offlineCount {
		sess.Where("id = ?", v.Id).Update(v)
	}
}

// 获取用户离线消息数量
func (ImOfflineCount) GetMemberOfflineMessage(sess *xorm.Session, userId int32) (*structs.OfflineMessageCount, error) {
	sess.Where("user_id = ? ", userId)
	messgaeCount := make([]*structs.OfflineMessageCount, 0)
	err := sess.Table(structs.TABLE_OFFLINE_MESSAFE_COUNT).OrderBy("offline_message_id asc").Find(&messgaeCount)
	if len(messgaeCount) > 0 {
		return messgaeCount[0], err
	}
	return nil, err
}

// 删除离线消息条数记录
func (ImOfflineCount) DelOfflineCount(sess *xorm.Session, userId int) (int64, error) {
	return sess.Table(structs.TABLE_OFFLINE_MESSAFE_COUNT).Where("user_id = ? ", userId).Delete(new(structs.OfflineMessageCount))
}
