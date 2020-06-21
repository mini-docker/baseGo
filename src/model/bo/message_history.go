package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

const (
	MESSAGE_INDEX_COUNT = 100
)

type MessageHistory struct{}

func (*MessageHistory) SaveMessageHistory(sess *xorm.Session, messageHistory *structs.MessageHistory) error {
	_, err := sess.Insert(messageHistory)
	return err
}

func (*MessageHistory) FindByIds(sess *xorm.Session, roomId int, lineId string, agencyId string, index int) ([]*structs.MessageHistory, error) {
	sess.Where("room_id = ?", roomId)
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.OrderBy("send_time desc")
	sess.Limit(MESSAGE_INDEX_COUNT, index*MESSAGE_INDEX_COUNT)
	messageHistorys := make([]*structs.MessageHistory, 0)
	err := sess.Table(new(structs.MessageHistory).TableName()).Find(&messageHistorys)
	return messageHistorys, err
}
