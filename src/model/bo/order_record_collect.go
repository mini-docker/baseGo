package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type RedPacketCollect struct{}

// 添加红包记录
func (*RedPacketCollect) InsertRedPacketCollects(sess *xorm.Session, red ...*structs.RedPacketCollect) (int64, error) {
	count, err := sess.Insert(red)
	return count, err
}

// 查询注单列表
func (*RedPacketCollect) QueryPacketCollects(sess *xorm.Session, lineId string, startTime, endTime int) ([]*structs.RedPacketCollect, error) {
	orders := make([]*structs.RedPacketCollect, 0)
	if startTime != 0 {
		sess.Where("create_time >= ? ", startTime)
	}
	if endTime != 0 {
		sess.Where("create_time <= ? ", endTime)
	}

	sess.Where("collect_status = ? ", 1)

	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	err := sess.Table(new(structs.RedPacketCollect).TableName()).OrderBy("id desc").Find(&orders)
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (*RedPacketCollect) UpdatePacketCollects(sess *xorm.Session, lineId string, startTime, endTime int) error {
	if startTime != 0 {
		sess.Where("create_time >= ? ", startTime)
	}
	if endTime != 0 {
		sess.Where("create_time <= ? ", endTime)
	}
	sess.Where("collect_status = ? ", 1)
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	st := new(structs.RedPacketCollect)
	st.CollectStatus = 2
	_, err := sess.Table(new(structs.RedPacketCollect).TableName()).Cols("collect_status").Update(st)
	if err != nil {
		return err
	}
	return err
}

func (*RedPacketCollect) GetInfoById(sess *xorm.Session, lineId string, status, gameType int) ([]structs.OrderCollectResp, error) {
	sess.Where("line_id = ?", lineId)
	sess.Where("status = ?", status)
	sess.Where("game_type = ?", gameType)
	sess.Where("receive_time <= 1581556881")
	sess.Where("account = red_sender")
	data := make([]structs.OrderCollectResp, 0)
	err := sess.Table(new(structs.OrderRecord).TableName()).Find(&data)
	return data, err
}

func (*RedPacketCollect) SettleOrder(sess *xorm.Session, redId int) ([]structs.OrderRecord, error) {
	sess.Where("red_id = ?", redId)
	data := make([]structs.OrderRecord, 0)
	err := sess.Find(&data)
	return data, err
}
