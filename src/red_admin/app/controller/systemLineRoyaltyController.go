package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

type SystemLineRoyaltyController struct {
	LineRoyaltyQueryReq struct {
		StartTime int `json:"startTime"` // 开始时间
		EndTime   int `json:"endTime"`   // 结束时间
	}

	LineAgencyRoyaltyQueryReq struct {
		LineId    string `json:"lineId"  valid:"Must;ErrorCode(3030)` // 线路id
		StartTime int    `json:"startTime"`                           // 开始时间
		EndTime   int    `json:"endTime"`                             // 结束时间
	}
}

var (
	SystemLineRoyaltyService = new(services.SystemLineRoyaltyService)
)

func (m SystemLineRoyaltyController) QueryLineRoyaltyList(ctx server.Context) error {
	req := &m.LineRoyaltyQueryReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	infos, err := SystemLineRoyaltyService.QueryLineRoyaltyList(req.StartTime, req.EndTime)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, infos)
}

func (m SystemLineRoyaltyController) QueryLineAgencyRoyaltyList(ctx server.Context) error {
	req := &m.LineAgencyRoyaltyQueryReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	infos, err := SystemLineRoyaltyService.QueryLineAgencyRoyaltyList(req.StartTime, req.EndTime, req.LineId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, infos)
}
