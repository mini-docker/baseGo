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
	UserService = new(services.UserService)
)

type UserController struct {
	QueryUserListReq struct {
		Status    int    `json:"status"`    // 状态  1 启用  2 停用
		IsOnline  int    `json:"isOnline"`  // 在线状态 1 在线 2 离线
		Account   string `json:"account"`   // 账号
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 代理id
	}

	QueryKickUserReq struct {
		Ids string `json:"ids" valid:"Must;ErrorCode(3031)"` // 会员id以","分割的字符串
	}

	EditUserStatusReq struct {
		Ids    string `json:"ids" valid:"Must;ErrorCode(3031)"`    // 会员id以","分割的字符串
		Status int    `json:"status" valid:"Must;ErrorCode(3028)"` // 状态 1 启用  2 停用
	}
}

// 获取会员列表
func (ac UserController) QueryUserList(ctx server.Context) error {
	req := &ac.QueryUserListReq
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

	result, err := UserService.QueryUserList(session.User.LineId, req.AgencyId, req.Status, req.IsOnline, req.Account, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 批量踢线
func (ac UserController) KickUsers(ctx server.Context) error {
	req := &ac.QueryKickUserReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := UserService.KickUsers(req.Ids)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	return common.HttpResultJson(ctx, nil)
}

// 批量修改会员状态
func (ac UserController) EditUsersStatus(ctx server.Context) error {
	req := &ac.EditUserStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := UserService.EditUsersStatus(req.Ids, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
