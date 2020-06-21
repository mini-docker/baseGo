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
	RedLogService = new(services.RedLogService)
)

type LogController struct {
	QueryLogReq struct {
		AgencyId  string `json:"agencyId"`  // 站点id
		LogType   int    `json:"logType"`   // 日志类型
		StartTime int    `json:"startTime"` // 开始时间
		EndTime   int    `json:"endTime"`   // 结束时间
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
	}
}

func (r LogController) QueryLogs(ctx server.Context) error {
	req := &r.QueryLogReq
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
	result, err := RedLogService.GetRedLogList(session.User.LineId, req.AgencyId, req.LogType, req.StartTime, req.EndTime, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}
