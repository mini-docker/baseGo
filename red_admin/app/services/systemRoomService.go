package services

import (
	"model/bo"
	"model/code"
	"model/structs"
	"red_admin/app/middleware/validate"
	"red_admin/conf"
)

type RoomService struct{}

var (
	RoomBo = new(bo.Room)
)

// 查询房间列表
func (RoomService) QueryRoomList(startTime, endTime, gameType, status int, page, pageSize int, agencyId, lineId string) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部房间信息
	count, rooms, err := RoomBo.QueryRoomList(sess, lineId, agencyId, startTime, endTime, gameType, status, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = rooms
	pageResp.Count = count
	return pageResp, nil
}

// 修改房间状态
func (RoomService) EditRoomStatus(id int, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断房间是否存在
	has, room := RoomBo.GetOne(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	room.Status = status
	err := RoomBo.ModifyRoomStatus(sess, room)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 群枚举
func (RoomService) RoomCode(lineId, agencyId string, gameType int) ([]*structs.RoomCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	roomCode, err := RoomBo.RoomCode(sess, lineId, agencyId, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return roomCode, err
}
