package controllers

import (
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
	"model"
	"model/code"
)

var (
	RedPacketSiteService = new(services.RedPacketSiteService)
)

type RedPacketSiteController struct {
	QueryPacketSiteReq struct {
		SiteName  string `json:"siteName"`  // 站点名称
		Status    int    `json:"status"`    // 状态
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
	}

	AddRedPacketSiteReq struct {
		AgencyId string `json:"agencyId" valid:"Must;ErrorCode(3070)"` // 站点id
		SiteName string `json:"siteName" valid:"Must;ErrorCode(3077)"` // 站点名称
		Status   int    `json:"status" valid:"Must;ErrorCode(3028)"`   // 状态 1 正常 2 停用
	}

	EditRedPacketSiteReq struct {
		Id       int    `json:"id" valid:"Must;ErrorCode(3031)"`       // 主键id
		SiteName string `json:"siteName" valid:"Must;ErrorCode(3077)"` // 站点名称
		Status   int    `json:"status" valid:"Must;ErrorCode(3028)"`   // 状态 1 正常 2 停用
	}

	DelRedPacketSiteReq struct {
		Id int `json:"id"` // 主键id
	}

	EditSiteStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)"`
		Status int `json:"status" valid:"Must;ErrorCode(3028)"` // 状态
	}
}

// 查询站点列表
func (m RedPacketSiteController) QueryPacketSiteList(ctx server.Context) error {
	req := &m.QueryPacketSiteReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	// 判断权限
	if session.IsAdmin != 1 {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.NO_PERMISSION})
	}

	result, err := RedPacketSiteService.QueryPacketSiteList(session.User.LineId, req.SiteName, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)

}

// 添加站点
func (m RedPacketSiteController) AddPacketSite(ctx server.Context) error {
	req := &m.AddRedPacketSiteReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	// 判断权限
	if session.IsAdmin != 1 {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.NO_PERMISSION})
	}

	err = RedPacketSiteService.AddPacketSite(session.LineId(), req.AgencyId, req.SiteName, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	return common.HttpResultJson(ctx, nil)
}

// 修改站点
func (m RedPacketSiteController) EditPacketSite(ctx server.Context) error {
	req := &m.EditRedPacketSiteReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	// 判断权限
	if session.IsAdmin != 1 {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.NO_PERMISSION})
	}

	err = RedPacketSiteService.EditPacketSite(req.Id, req.SiteName, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	return common.HttpResultJson(ctx, nil)
}

// 删除站点
func (m RedPacketSiteController) DelPacketSite(ctx server.Context) error {
	req := &m.DelRedPacketSiteReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	// 判断权限
	if session.IsAdmin != 1 {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.NO_PERMISSION})
	}

	err = RedPacketSiteService.DelPacketSite(req.Id)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改站点状态
func (m RedPacketSiteController) EditRedPacketSiteStatus(ctx server.Context) error {
	req := &m.EditSiteStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	// 判断权限
	if session.IsAdmin != 1 {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.NO_PERMISSION})
	}

	err = RedPacketSiteService.EditRedPacketSiteStatus(req.Id, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 站点枚举
func (m RedPacketSiteController) SiteCode(ctx server.Context) error {
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	if user.IsAdmin != 1 {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.NO_PERMISSION})
	}
	agencys, err := RedPacketSiteService.SiteCode(user.LineId())
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, agencys)
}
