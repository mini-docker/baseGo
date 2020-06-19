package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

type OrderStatistical struct {
	OrderStatisticalReq struct {
		LineId    string `json:"lineId"`    // 线路
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
	result, err := orderStatistical.QueryOrderStatistical(req.LineId, req.AgencyId, req.StartTime, req.EndTime, req.GameType)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}
