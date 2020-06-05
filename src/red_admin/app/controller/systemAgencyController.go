package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

// 超管控制器
type SystemAgencyController struct {
	AddAgencyAdminReq struct {
		LineId          string `json:"lineId" valid:"Must;ErrorCode(3030)`          // 线路id
		Account         string `json:"account" valid:"Must;ErrorCode(3026)`         // 账号
		Password        string `json:"password" valid:"Must;ErrorCode(3027)`        // 密码
		ConfirmPassword string `json:"confirmPassword,omitempty"`                   // 确认密码
		Status          int    `json:"status" valid:"Must;ErrorCode(3028)`          // 状态 1正常 2停用
		WhiteIpAddress  string `json:"whiteIpAddress" valid:"Must;ErrorCode(3064)"` // ip白名单
	}

	QueryAgencyAdminReq struct {
		LineId    string `json:"lineId"`    // 线路id
		Account   string `json:"account"`   // 账号
		IsOnline  int    `json:"isOnline"`  // 在线状态
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
	}

	QueryAgencyAdminOneReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)` // 超管id
	}

	EditAgencyAdminReq struct {
		Id              int    `json:"id" valid:"Must;ErrorCode(3031)`              // 超管id
		Password        string `json:"password"`                                    // 密码
		ConfirmPassword string `json:"confirmPassword,omitempty"`                   // 确认密码
		Status          int    `json:"status" valid:"Must;ErrorCode(3028)`          // 状态 1正常 2停用
		WhiteIpAddress  string `json:"whiteIpAddress" valid:"Must;ErrorCode(3064)"` // ip白名单
	}

	QueryAgencyReq struct {
		LineId    string `json:"lineId"`    // 线路id
		IsOnline  int    `json:"isOnline"`  // 在线状态
		Status    int    `json:"status"`    // 停用状态
		Account   string `json:"account"`   // 账号
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 站点id
	}

	EditAgencyStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)`     // 代理id
		Status int `json:"status" valid:"Must;ErrorCode(3028)` // 状态
	}

	AgencyCodeReq struct {
		LineId string `json:"lineId"` // 线路id
	}
}

var (
	SystemAgencyService = new(services.SystemAgencyService)
)

// 查询全部超管
func (m SystemAgencyController) QueryAgencyAdminList(ctx server.Context) error {
	req := &m.QueryAgencyAdminReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := SystemAgencyService.QueryAgencyAdminList(req.LineId, req.Account, req.IsOnline, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加超管
func (m SystemAgencyController) AddAgencyAdmin(ctx server.Context) error {
	req := &m.AddAgencyAdminReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemAgencyService.AddAgency(req.Account, req.Password, req.ConfirmPassword, req.LineId, req.Status, req.WhiteIpAddress)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询超管信息
func (m SystemAgencyController) QueryAgencyAdminOne(ctx server.Context) error {
	req := &m.QueryAgencyAdminOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	Agency, err := SystemAgencyService.QueryAgencyOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, Agency)
}

// 修改超管
func (m SystemAgencyController) EditAgencyAdmin(ctx server.Context) error {
	req := &m.EditAgencyAdminReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemAgencyService.EditAgency(req.Id, req.Password, req.ConfirmPassword, req.Status, req.WhiteIpAddress)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 查询全部代理
func (m SystemAgencyController) QueryAgencyList(ctx server.Context) error {
	req := &m.QueryAgencyReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	info, err := SystemAgencyService.QueryAgencyList(req.LineId, req.Account, req.AgencyId, req.IsOnline, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, info)
}

// 修改代理状态
func (m SystemAgencyController) EditAgencyStatus(ctx server.Context) error {
	req := &m.EditAgencyStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemAgencyService.EditAgencyStatus(req.Id, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 站点枚举
func (m SystemAgencyController) SiteCode(ctx server.Context) error {
	req := &m.AgencyCodeReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	agencys, err := SystemAgencyService.SiteCode(req.LineId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, agencys)
}
