package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
)

var RobotService = new(services.RobotService)

type RobotController struct {
	QueryRobotListReq struct {
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 站点id
	}
	CreatRobotAccountsReq struct {
		Num      int    `json:"num"`      // 数量
		AgencyId string `json:"agencyId"` // 超管选择的代理id
	}
	InsertRobotsReq struct {
		Accounts string `json:"accounts"` // 账号","分割的字符串
		AgencyId string `json:"agencyId"` // 超管选择的代理id
	}
	DelRobotsReq struct {
		Ids string `json:"ids"` // 机器人id","分割的字符串
	}
}

// 查询机器人列表
func (m RobotController) QueryRobotList(ctx server.Context) error {
	req := &m.QueryRobotListReq
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
	pages, err := RobotService.QueryRobotList(session.LineId(), req.AgencyId, req.PageIndex, req.PageSize)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, pages)
}

// 批量生成机器人账号
func (m RobotController) CreatRobotAccounts(ctx server.Context) error {
	req := &m.CreatRobotAccountsReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	if session.IsAdmin == 1 {
		if req.AgencyId == "" {
			return common.HttpResultJsonError(ctx, &validate.Err{Code: code.AGENCY_ID_CAN_NOT_BE_EMPTY})
		}
	} else {
		req.AgencyId = session.AgencyId()
	}
	if req.Num > 100 {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.ROBOT_NUM_IS_TOO_LARGE})
	}
	robots := RobotService.CreatRobotAccounts(session.LineId(), req.AgencyId, req.Num)
	return common.HttpResultJson(ctx, robots)
}

// 批量保存机器人
func (m RobotController) InsertRobots(ctx server.Context) error {
	req := &m.InsertRobotsReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	if session.IsAdmin == 1 {
		if req.AgencyId == "" {
			return common.HttpResultJsonError(ctx, &validate.Err{Code: code.AGENCY_ID_CAN_NOT_BE_EMPTY})
		}
	} else {
		req.AgencyId = session.AgencyId()
	}
	err = RobotService.InsertRobots(session.LineId(), req.AgencyId, req.Accounts)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 批量删除机器人
func (m RobotController) DelRobots(ctx server.Context) error {
	req := &m.DelRobotsReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := RobotService.DelRobots(req.Ids)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
