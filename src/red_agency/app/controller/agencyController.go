package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
)

var (
	SessionService      = new(middleware.SessionService)
	SystemAgencyService = new(services.SystemAgencyService)
)

type AgencyController struct {
	QueryAgencyReq struct {
		Account   string `json:"account"`   // 账号
		IsOnline  int    `json:"isOnline"`  // 在线状态
		Status    int    `json:"status"`    // 停用状态
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 代理id
	}

	AddAgencyReq struct {
		Account         string `json:"account" valid:"Must;ErrorCode(3026)"`        // 账号
		AgencyId        string `json:"agencyId" valid:"Must;ErrorCode(3026)"`       // 站点
		Limit           string `json:"limit"`                                       // 额度
		Password        string `json:"password" valid:"Must;ErrorCode(3027)"`       // 密码
		ConfirmPassword string `json:"confirmPassword,omitempty"`                   // 确认密码
		Status          int    `json:"status" valid:"Must;ErrorCode(3028)"`         // 状态
		WhiteIpAddress  string `json:"whiteIpAddress" valid:"Must;ErrorCode(3064)"` // ip白名单
	}

	QueryAgencyOneReq struct {
		Id int `json:"id"` // 代理id
	}

	ResetPasswordReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"` // 代理id
	}

	EditAgencyReq struct {
		Id              int    `json:"id" valid:"Must;ErrorCode(3031)"`             // 代理id
		Password        string `json:"password"`                                    // 密码
		ConfirmPassword string `json:"confirmPassword,omitempty"`                   // 确认密码
		Status          int    `json:"status" valid:"Must;ErrorCode(3028)"`         // 状态
		WhiteIpAddress  string `json:"whiteIpAddress" valid:"Must;ErrorCode(3064)"` // ip白名单
	}

	EditAgencyStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)"`     // 代理id
		Status int `json:"status" valid:"Must;ErrorCode(3028)"` // 状态
	}

	DelAgencyReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"` // 代理id
	}

	EditPasswordReq struct {
		OldPassword     string `json:"oldPassword" valid:"Must;ErrorCode(3027)"` // 原密码
		Password        string `json:"password" valid:"Must;ErrorCode(3027)"`    // 新密码
		ConfirmPassword string `json:"confirmPassword,omitempty"`                // 确认密码
	}
}

// 代理列表
func (m AgencyController) QueryAgencyList(ctx server.Context) error {
	req := &m.QueryAgencyReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	if user.IsAdmin != 1 {
		req.AgencyId = user.AgencyId()
	}
	info, err := SystemAgencyService.QueryAgencyList(user.LineId(), req.Account, req.AgencyId, req.IsOnline, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, info)
}

// 添加代理
func (m AgencyController) AddAgency(ctx server.Context) error {
	req := &m.AddAgencyReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	if user.IsAdmin != 1 {
		req.AgencyId = user.AgencyId()
	}
	ip := ctx.RealIP()
	err = SystemAgencyService.AddAgency(req.Account, req.Password, req.ConfirmPassword, user.LineId(), req.AgencyId, req.Status, req.WhiteIpAddress, user.Account(), user.Uid(), ip)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 查询单个代理
func (m AgencyController) QueryAgencyOne(ctx server.Context) error {
	req := &m.QueryAgencyOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	Agency, err := SystemAgencyService.QueryAgencyOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, Agency)
}

// 修改代理
func (m AgencyController) EditAgency(ctx server.Context) error {
	req := &m.EditAgencyReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	err = SystemAgencyService.EditAgency(req.Id, req.Password, req.ConfirmPassword, req.Status, req.WhiteIpAddress, user.Account(), user.Uid(), ip)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改代理状态
func (m AgencyController) EditAgencyStatus(ctx server.Context) error {
	req := &m.EditAgencyStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	err = SystemAgencyService.EditAgencyStatus(req.Id, req.Status, user.Account(), user.Uid(), ip)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 重置密码
func (m AgencyController) ResetAgencyPassword(ctx server.Context) error {
	req := &m.ResetPasswordReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	err = SystemAgencyService.ResetPassword(req.Id, user.Account(), user.Uid(), ip)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改密码
func (m AgencyController) EditPassword(ctx server.Context) error {
	req := &m.EditPasswordReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	err = SystemAgencyService.EditPassword(user.Account(), req.OldPassword, req.Password, req.ConfirmPassword, user.Account(), user.Uid(), ip)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 删除代理
func (m AgencyController) DelAgency(ctx server.Context) error {
	req := &m.DelAgencyReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	err = SystemAgencyService.DelAgency(req.Id, user.Account(), user.Uid(), ip)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
