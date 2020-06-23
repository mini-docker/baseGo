package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
)

type OrderStatistical struct {
	OrderStatisticalReq struct {
		AgencyId  string `json:"agencyId"`  // 站点
		StartTime int64  `json:"startTime"` // 开始时间
		EndTime   int64  `json:"endTime"`   // 结束时间
		GameType  int64  `json:"gameType"`  // 游戏类型
	}
}

var (
	orderStatistical = new(services.OrderStatistical)
)

// 统计查询
func (m *OrderStatistical) QueryOrderStatistical(ctx server.Context) error {
	req := &m.OrderStatisticalReq
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
	result, err := orderStatistical.QueryOrderStatistical(session.LineId(), req.AgencyId, req.StartTime, req.EndTime, req.GameType)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}
