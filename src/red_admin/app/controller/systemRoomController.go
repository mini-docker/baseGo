package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

var (
	RoomService = new(services.RoomService)
)

type SystemRoomController struct {
	QueryRoomListReq struct {
		StartTime int    `json:"startTime"` // 开始时间
		EndTime   int    `json:"endTime"`   // 结束时间
		GameType  int    `json:"gameType"`  // 房间类型
		Status    int    `json:"status"`    // 状态  1 启用  2 停用
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 站点id
		LineId    string `json:"lineId"`
	}

	EditRoomStatusReq struct {
		Id     int `json:"id"`     // 房间id
		Status int `json:"status"` // 状态
	}

	RoomCodeReq struct {
		LineId   string `json:"lineId"`   // 线路id
		AgencyId string `json:"agencyId"` // 站点id
		GameType int    `json:"gameType"` // 游戏类型
	}
}

// 查询全部房间
func (m SystemRoomController) QueryRoomList(ctx server.Context) error {
	req := &m.QueryRoomListReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	result, err := RoomService.QueryRoomList(req.StartTime, req.EndTime, req.GameType, req.Status, req.PageIndex, req.PageSize, req.AgencyId, req.LineId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 修改房间状态
func (m SystemRoomController) EditRoomStatus(ctx server.Context) error {
	req := &m.EditRoomStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := RoomService.EditRoomStatus(req.Id, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 群枚举
func (c *SystemRoomController) RoomCode(ctx server.Context) error {
	req := &c.RoomCodeReq
	err := ctx.Validate(req)
	if err != nil {
		return err
	}
	roomCode, err := RoomService.RoomCode(req.LineId, req.AgencyId, req.GameType)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, roomCode)
}
