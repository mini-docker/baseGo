package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
)

type RedPacketLog struct{}

// 查询红包领取记录
func (*RedPacketLog) GetRedLog(sess *xorm.Session, lineId string, agencyId string, redId int) ([]structs.OrderRecord, error) {
	sess.Where("red_id = ? ", redId)
	sess.Where("line_id = ?", lineId)
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	data := make([]structs.OrderRecord, 0)
	err := sess.Find(&data)
	return data, err
}

// 查询红包领取记录
func (*RedPacketLog) GetRedLogByRedIdList(sess *xorm.Session, redIds []int) ([]structs.OrderRecord, error) {
	sess.In("red_id", redIds)
	data := make([]structs.OrderRecord, 0)
	err := sess.Find(&data)
	return data, err
}

// 查询红包领取记录
func (*RedPacketLog) GetRedOrderRecord(sess *xorm.Session, linId string, agencyId string, redId int) ([]structs.RedPacketLogInfoResp, error) {
	sess.Where("red_id = ? ", redId)
	sess.Where("line_id = ?", linId)
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	data := make([]structs.RedPacketLogInfoResp, 0)
	err := sess.Table(new(structs.OrderRecord).TableName()).Find(&data)
	return data, err
}

// 查询会员领取过的
func (*RedPacketLog) GetUserRedLog(sess *xorm.Session, userId, redId int) (bool, *structs.OrderRecord, error) {
	sess.Where("red_id = ? ", redId)
	sess.Where("user_id = ? ", userId)
	info := new(structs.OrderRecord)
	has, err := sess.Get(info)
	return has, info, err
}

// 查询会员领取过的
func (*RedPacketLog) GetUserRedLogList(sess *xorm.Session, linId string, agencyId string, userId, startTime, endTime,
	gameType, gamePlay int) ([]structs.OrderRecord, error) {
	//sess.Where("red_id = ? ", redId)
	sess.Where("user_id = ? ", userId)
	sess.Where("line_id = ?", linId)
	sess.Where("agency_id = ?", agencyId)
	if startTime > 0 {
		sess.Where("red_start_time >= ?", startTime)
	}
	if endTime > 0 {
		sess.Where("red_start_time <= ?", endTime)
	}
	if gameType > 0 && gameType != 99 {
		sess.Where("game_type = ?", gameType)
	} else if gameType == 99 { // 99指普通红包
		sess.Where("game_type = 0")
	}
	if gamePlay > 0 {
		sess.Where("game_play = ?", gamePlay)
	}
	data := make([]structs.OrderRecord, 0)
	err := sess.Find(&data)
	return data, err
}

// 修改红包注单状态
func (*RedPacketLog) UpdateRedLogStatus(sess *xorm.Session, redId, status int) (int, error) {
	info := new(structs.OrderRecord)
	info.Status = status
	count, err := sess.Table(new(structs.OrderRecord).TableName()).ID(redId).Cols("status").Update(info)
	return int(count), err
}

// 修改红包状态
func (*RedPacketLog) UpdateRedLogListStatus(sess *xorm.Session, data structs.OrderRecord) (int64, error) {
	count, err := sess.Table(new(structs.OrderRecord).TableName()).ID(data.Id).Cols("money", "real_money", "royalty_money", "status", "extra", "valid_bet", "robot_win").Update(data)
	return count, err
}

// 修改发包人红包领取记录
func (*RedPacketLog) UpdateAdminRedLogInfo(sess *xorm.Session, logInfo *structs.OrderRecord, redId int) (int64, error) {
	sess.Where("user_id = ? ", logInfo.UserId)
	sess.Where("line_id = ?", logInfo.LineId)
	sess.Where("agency_id = ?", logInfo.AgencyId)
	sess.ID(logInfo.Id)
	sess.Cols("receive_money, real_money, money, status, extra,royalty_money,valid_bet,robot_win")
	count, err := sess.Table(logInfo.TableName()).Update(logInfo)
	return count, err
}

// 添加红包记录
func (*RedPacketLog) InsertRedPacketLog(sess *xorm.Session, red *structs.OrderRecord) (int64, error) {
	count, err := sess.Insert(red)
	return count, err
}

// 添加红包记录
func (*RedPacketLog) InsertRedPacketLogList(sess *xorm.Session, red ...*structs.OrderRecord) (int64, error) {
	count, err := sess.Insert(red)
	return count, err
}

// 查询注单列表
func (*RedPacketLog) QueryRedRecordList(sess *xorm.Session, lineId string, agencyId string, startTime,
	endTime, gameType, status int, orderNo, account, redSender string, page, pageSize, redId, roomId, isRobot int) (int64, []*structs.RedPacketLogInfoResp, error) {
	orders := make([]*structs.RedPacketLogInfoResp, 0)
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	if startTime != 0 {
		sess.Where("red_start_time >= ? ", startTime)
	}
	if endTime != 0 {
		sess.Where("red_start_time <= ? ", endTime)
	}
	if gameType != 0 {
		sess.Where("game_type = ? ", gameType)
	}
	if status != 10 {
		sess.Where("status = ? ", status)
	}
	if orderNo != "" {
		sess.Where("order_no = ? ", orderNo)
	}
	if account != "" {
		sess.Where("account = ? ", account)
	}
	if redSender != "" {
		sess.Where("red_sender = ? ", redSender)
	}
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if redId != 0 {
		sess.Where("red_id = ?", redId)
	}
	if roomId != 0 {
		sess.Where("room_id = ?", roomId)
	}
	if isRobot == 1 {
		sess.Where("is_robot = ?", model.USER_IS_ROBOT_YES)
	}
	if isRobot == 2 {
		sess.Where("is_robot = ?", 0)
	}
	count, err := sess.Table(new(structs.OrderRecord).TableName()).Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&orders)
	if err != nil {
		return 0, nil, err
	}
	return count, orders, err
}

func (*RedPacketLog) GetRedInfo(sess *xorm.Session, lineId string, agencyId string, redId int) ([]*structs.RedOrderResp, error) {
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	sess.Where("red_id = ? ", redId)
	orders := make([]*structs.RedOrderResp, 0)
	err := sess.Table(new(structs.OrderRecord).TableName()).Find(&orders)
	return orders, err
}

// 查询红包领取记录
func (*RedPacketLog) GetRedOrderRecordByRedIdWithStatus12(sess *xorm.Session, linId string, agencyId string, redId int) ([]structs.OrderCollectResp, error) {
	sess.Where("red_id = ? ", redId)
	sess.Where("line_id = ?", linId)
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	data := make([]structs.OrderCollectResp, 0)
	err := sess.Table(new(structs.OrderRecord).TableName()).Find(&data)
	return data, err
}

// 查询红包领取数量和金额
func (*RedPacketLog) GetOrderReceiveCountResp(sess *xorm.Session, linId string, agencyId string, redIds []int) ([]structs.OrderReceiveCountResp, error) {
	sess.Where("line_id = ?", linId)
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if len(redIds) > 0 {
		sess.In("red_id", redIds)
	}
	sess.GroupBy("red_id")
	sess.Select("red_id, count(id) as count, sum(receive_money) as receiveMoney")
	data := make([]structs.OrderReceiveCountResp, 0)
	err := sess.Table(new(structs.OrderRecord).TableName()).Find(&data)
	return data, err
}

// 根据房间id查询未结算注单
func (*RedPacketLog) GetOrdersByRoomId(sess *xorm.Session, roomId int, lineId, agencyId string) (int64, error) {
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("room_id = ?", roomId)
	sess.Where("status = 0")
	orders := make([]*structs.OrderRecord, 0)
	count, err := sess.Table(new(structs.OrderRecord).TableName()).FindAndCount(&orders)
	return count, err
}

// 查询有会员参与的红包
func (*RedPacketLog) GetMemberReds(sess *xorm.Session) ([]structs.RedPacket, error) {
	reds := make([]structs.RedPacket, 0)
	err := sess.SQL("select * from red_packet where id in (select distinct(red_id) from red_order_record where is_robot = 0 ) and status = 2").Find(&reds)
	return reds, err
}
