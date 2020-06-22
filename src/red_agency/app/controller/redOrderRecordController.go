package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
)

var (
	RedOrderService = new(services.RedOrderService)
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
		AgencyId  string `json:"agencyId"`  //代理
		RedId     int    `json:"redId"`     // 局号
		RoomId    int    `json:"roomId"`    // 群id
		IsRobot   int    `json:"isRobot"`   // 是否是机器人注单
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

	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	if session.IsAdmin != 1 {
		req.AgencyId = session.AgencyId()
	}

	// 查询注单列表
	pages, err := RedOrderService.QueryRedRecordList(session.LineId(), req.AgencyId, req.StartTime, req.EndTime, req.GameType, req.Status, req.OrderNo, req.Account, req.RedSender, req.PageIndex, req.PageSize, req.RedId, req.RoomId, req.IsRobot)

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
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	records, err := RedOrderService.GetRedInfo(session.LineId(), session.AgencyId(), req.RedId)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, records)
}
