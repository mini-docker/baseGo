package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/fecho/xorm/builder"
	"baseGo/src/fecho/xorm/help"
	"baseGo/src/model"
	"baseGo/src/model/structs"
	"fmt"
	"strconv"
	"strings"
)

type RedPacket struct{}

// 查询红包详情
func (*RedPacket) GetOne(sess *xorm.Session, id int) (bool, *structs.RedPacket) {
	sess.Where("id = ? ", id)
	sess.Where("delete_time = ?", 0)
	//sess.Where("status = ?", model.RED_STATUS_NORMAL)
	info := new(structs.RedPacket)
	has, _ := sess.Get(info)
	return has, info
}

func (*RedPacket) ByRoomIdGetRedInfo(sess *xorm.Session, id int) (bool, *structs.RedPacket) {
	sess.Where("id = ? ", id)
	//sess.Where("delete_time = ?", 0)
	info := new(structs.RedPacket)
	has, _ := sess.Get(info)
	return has, info
}

// 查询红包详情
func (*RedPacket) GetRedPacketInfo(sess *xorm.Session, lineId string, agencyId string, userId, id int) (bool, *structs.RedPacket) {
	sess.Where("id = ? ", id)
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ? ", agencyId)
	if userId != 0 {
		sess.Where("user_id = ? ", userId)
	}

	sess.Where("delete_time = ?", 0)
	//sess.Where("status = ?", model.RED_STATUS_NORMAL)
	info := new(structs.RedPacket)
	has, _ := sess.Get(info)
	return has, info
}

// 查询红包详情
func (*RedPacket) GetRedInfo(sess *xorm.Session, lineId string, agencyId string, userId, id, roomId int) (bool, *structs.RedPacket) {
	sess.Where("id = ? ", id)
	sess.Where("line_id = ? ", lineId)
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	if userId != 0 {
		sess.Where("user_id = ? ", userId)
	}
	if roomId > 0 {
		sess.Where("room_id = ? ", roomId)
	}

	sess.Where("delete_time = ?", 0)
	//sess.Where("status = ?", model.RED_STATUS_NORMAL)
	info := new(structs.RedPacket)
	has, _ := sess.Get(info)
	return has, info
}

// 修改红包状态
func (*RedPacket) UpdateRedStatus(sess *xorm.Session, id, status int) (int, error) {
	info := new(structs.RedPacket)
	info.Status = status
	count, err := sess.ID(id).Cols("status").Update(info)
	return int(count), err
}

// 修改红包状态
func (*RedPacket) UpdateRed(sess *xorm.Session, req *structs.RedPacket, fields string) (int, error) {
	count, err := sess.ID(req.Id).Where("line_id = ? and agency_id = ? ", req.LineId, req.AgencyId).Cols(fields).Update(req)
	return int(count), err
}

// 添加红包
func (*RedPacket) InsertRedPacket(sess *xorm.Session, red *structs.RedPacket) (int, error) {
	count, err := sess.Insert(red)
	return int(count), err
}

// 查询红包记录
func (*RedPacket) GetUserRedList(sess *xorm.Session, lineId string, agencyId string, userId, startTime, endTime, gameType, gamePlay int, redIds []int, pageParams *help.PageParams) (int, []structs.RedPacketResp, error) {
	sess.Where("line_id = ? ", lineId)
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	if len(redIds) > 0 && userId != 0 {
		redIdsStr := make([]string, 0)
		for _, v := range redIds {
			redIdsStr = append(redIdsStr, strconv.Itoa(v))
		}
		inStr := fmt.Sprintf("id In (%v)", strings.Join(redIdsStr, ","))
		sess.Where("user_id = ? or "+inStr, userId)
	} else {
		if len(redIds) > 0 {
			sess.In("id", redIds)
		}
		if userId != 0 {
			sess.Where("user_id = ? ", userId)
		}
	}

	if startTime > 0 {
		sess.Where("create_time >= ?", startTime)
	}
	if endTime > 0 {
		sess.Where("create_time <= ?", endTime)
	}
	statusSql := "status IN (" + strings.Join([]string{strconv.Itoa(model.RED_STATUS_OVER), strconv.Itoa(model.RED_STATUS_INVALID)}, ",") + ")"
	if gameType > 0 && gameType != 99 {
		sess.Where("red_type = ? and "+statusSql, gameType)
	} else if gameType == 99 {
		sess.Where("red_type = 0")
	} else {
		sess.Where("(red_type = ?) or (red_type > ? and "+statusSql+")", model.ORDINARY_RED_ENVELOPE, model.ORDINARY_RED_ENVELOPE)
	}
	if gamePlay > 0 {
		sess.Where("red_play = ?", gamePlay)
	}
	pageParams.Make(sess)
	data := make([]structs.RedPacketResp, 0)
	count, err := sess.Table(new(structs.RedPacket).TableName()).OrderBy("id desc").FindAndCount(&data)
	return int(count), data, err
}

// 查询红包记录
func (*RedPacket) GetRedByIdList(sess *xorm.Session, lineId string, agencyId string, roomId, gameType, gamePlay int, redIds []int) ([]structs.RedPacketResp, error) {
	sess.Where("line_id = ? ", lineId)
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	if gameType > 0 {
		sess.Where("red_type = ?", gameType)
	}
	if gamePlay > 0 {
		sess.Where("red_play = ?", gamePlay)
	}
	if roomId > 0 {
		sess.Where("room_id = ?", roomId)
	}
	sess.Where("status > ?", model.ORDINARY_RED_ENVELOPE)
	sess.In("id", redIds)
	data := make([]structs.RedPacketResp, 0)
	err := sess.Table(new(structs.RedPacket).TableName()).Find(&data)
	return data, err
}

// 查询红包记录
func (*RedPacket) GetOrdinaryRedList(sess *xorm.Session, lineId string, agencyId string, startTime, endTime, status int, roomName string, pageParams *help.PageParams) (int, []structs.RedPacket, error) {
	sess.Where("line_id = ? ", lineId)
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}

	if startTime > 0 {
		sess.Where("create_time >= ?", startTime)
	}
	if endTime > 0 {
		sess.Where("create_time <= ?", endTime)
	}
	if status > 0 {
		sess.Where("status = ?", status)
	}
	if roomName != "" {
		sess.Where(builder.Like{"room_name", roomName})
	}
	sess.Where("red_type = 0")
	sess.Where("delete_time = 0")
	pageParams.Make(sess)
	data := make([]structs.RedPacket, 0)
	count, err := sess.Table(new(structs.RedPacket).TableName()).FindAndCount(&data)
	return int(count), data, err
}

// 查询未结算红包
func (*RedPacket) GetAllNeedSettlementPacket(sess *xorm.Session, gameType int) ([]structs.OrderRecord, error) {
	reds := make([]structs.OrderRecord, 0)
	sess.Where("game_type = ? ", gameType)
	if gameType == model.ORDINARY_RED_ENVELOPE {
		sess.Where("status = ?", 1)
	} else {
		sess.Where("status = ?", 0)
	}

	err := sess.Table(new(structs.OrderRecord).TableName()).Find(&reds)
	return reds, err
}

// 查询发包人注单
func (*RedPacket) GetSenderOrder(sess *xorm.Session, redId, senderId int) (*structs.OrderRecord, error) {
	order := new(structs.OrderRecord)
	sess.Where("red_id = ?", redId)
	sess.Where("user_id = ?", senderId)
	_, err := sess.Get(order)
	return order, err
}

// 查询扫雷异常注单
func (*RedPacket) GetCheckOrders(sess *xorm.Session, lineId, agencyId string) ([]*structs.OrderRecord, error) {
	reds := make([]*structs.OrderRecord, 0)
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	sess.Where("game_type = ? and account = red_sender and real_money < ? and is_robot = ?", 2, 0, 0)
	err := sess.Find(&reds)
	return reds, err
}
