package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

var (
	SystemUserService = new(services.SystemUserService)
)

type SystemUserController struct {
	QueryUserListReq struct {
		Status    int    `json:"status"`    // 状态  1 启用  2 停用
		IsOnline  int    `json:"isOnline"`  // 在线状态 1 在线 2 离线
		Account   string `json:"account"`   // 账号
		LineId    string `json:"lineId"`    // 线路id
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 站点id
	}

	QueryKickUserReq struct {
		Ids string `json:"ids" valid:"Must;ErrorCode(3031)` // 会员id以","分割的字符串
	}

	EditUserStatusReq struct {
		Ids    string `json:"ids" valid:"Must;ErrorCode(3031)`    // 会员id以","分割的字符串
		Status int    `json:"status" valid:"Must;ErrorCode(3028)` // 状态 1 启用  2 停用
	}
}

// 获取会员列表
func (ac SystemUserController) QueryUserList(ctx server.Context) error {
	req := &ac.QueryUserListReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := SystemUserService.QueryUserList(req.LineId, req.AgencyId, req.Status, req.IsOnline, req.Account, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 批量踢线
func (ac SystemUserController) KickUsers(ctx server.Context) error {
	req := &ac.QueryKickUserReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := SystemUserService.KickUsers(req.Ids)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	return common.HttpResultJson(ctx, nil)
}

// 批量修改会员状态
func (ac SystemUserController) EditUsersStatus(ctx server.Context) error {
	req := &ac.EditUserStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemUserService.EditUsersStatus(req.Ids, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
