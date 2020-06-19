package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

var (
	RedOrderService = new(services.SystemRedOrderService)
)

type RedOrderRecordController struct {
	OrderRecordReq struct {
		StartTime int    `json:"startTime"` // 开始时间
		EndTime   int    `json:"endTime"`   // 结束时间
		GameType  int    `json:"gameType"`  // 游戏名称 2扫雷 1牛牛
		Status    int    `json:"status"`    // 注单状态   1进行中 2已结算 3已返还
		OrderNo   string `json:"orderNo"`   // 注单号
		Account   string `json:"account"`   // 账号
		RedSender string `json:"redSender"` // 发包人账号
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 站点id
		RedId     int    `json:"redId"`     // 局号
		LineId    string `json:"lineId"`    // 线路id
		RoomId    int    `json:"roomId"`    // 群id
		IsRobot   int    `json:"isRobot"`   // 是否机器人注单  1 是  2 否
	}

	GetRedInfoReq struct {
		RedId int `json:"redId"` // 红包id
	}
}

// 查询红包注单列表
func (m RedOrderRecordController) QueryRedRecordList(ctx server.Context) error {
	req := &m.OrderRecordReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 查询注单列表
	pages, err := RedOrderService.QueryRedRecordList(req.StartTime, req.EndTime, req.GameType, req.Status, req.OrderNo, req.Account, req.RedSender, req.PageIndex, req.PageSize, req.AgencyId, req.RedId, req.RoomId, req.IsRobot, req.LineId)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, pages)
}

// 查询红包领取记录
func (m RedOrderRecordController) GetRedInfo(ctx server.Context) error {
	req := &m.GetRedInfoReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	records, err := RedOrderService.GetRedInfo(req.RedId)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, records)
}
