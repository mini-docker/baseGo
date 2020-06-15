package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

// 线路控制器
type SystemLineController struct {
	AddLineReq struct {
		LineId    string  `json:"lineId" valid:"Must;ErrorCode(3030)`    // 线路id
		LineName  string  `json:"lineName" valid:"Must;ErrorCode(3034)`  // 线路名称
		LimitCost float64 `json:"limitCost" valid:"Must;ErrorCode(3035)` // 线路额度
		MealId    int     `json:"mealId" valid:"Must;ErrorCode(3036)`    // 套餐id
		Domain    string  `json:"domain" valid:"Must;ErrorCode(3037)`    // 域名
		TransType int     `json:"transType" valid:"Must;ErrorCode(3038)` // 交易模式  1 钱包  2 额度转换
		ApiUrl    string  `json:"apiUrl"`                                // 钱包api地址
		Md5key    string  `json:"md5key" valid:"Must;ErrorCode(3039)`    // md5key
		RsaPubKey string  `json:"rsaPubKey" valid:"Must;ErrorCode(3040)` //
		RsaPriKey string  `json:"rsaPriKey" valid:"Must;ErrorCode(3041)` //
		Status    int     `json:"status" valid:"Must;ErrorCode(3028)`    // 状态
	}

	QueryLineOneReq struct {
		Id int `json:"id"` // 线路id
	}

	QueryLineReq struct {
		LineId    string `json:"lineId"`    // 线路id
		LineName  string `json:"lineName"`  // 线路名称
		Status    int    `json:"status"`    // 状态 1 启用 2 停用 3 维护
		TransType int    `json:"transType"` // 交易模式  1 钱包  2 额度转换
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
	}

	EditLineReq struct {
		Id        int     `json:"id" valid:"Must;ErrorCode(3031)`        // 线路id
		LineId    string  `json:"lineId" valid:"Must;ErrorCode(3030)`    // 线路id
		LineName  string  `json:"lineName" valid:"Must;ErrorCode(3034)`  // 线路名称
		LimitCost float64 `json:"limitCost" valid:"Must;ErrorCode(3035)` // 线路额度
		MealId    int     `json:"mealId" valid:"Must;ErrorCode(3036)`    // 套餐id
		Domain    string  `json:"domain" valid:"Must;ErrorCode(3037)`    // 域名
		TransType int     `json:"transType" valid:"Must;ErrorCode(3038)` // 交易模式  1 钱包  2 额度转换
		ApiUrl    string  `json:"apiUrl"`                                // 钱包api地址
		Md5key    string  `json:"md5key" valid:"Must;ErrorCode(3039)`    // md5key
		RsaPubKey string  `json:"rsaPubKey" valid:"Must;ErrorCode(3040)` //
		RsaPriKey string  `json:"rsaPriKey" valid:"Must;ErrorCode(3041)` //
		Status    int     `json:"status" valid:"Must;ErrorCode(3028)`    // 状态
	}

	EditLineStatusReq struct {
		Id     int `json:"id"`     // 线路id
		Status int `json:"status"` // 状态
	}
}

var (
	SystemLineService = new(services.SystemLineService)
)

// 查询全部线路
func (m SystemLineController) QueryLineList(ctx server.Context) error {
	req := &m.QueryLineReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := SystemLineService.QuerySystemLineList(req.LineId, req.LineName, req.Status, req.TransType, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加线路
func (m SystemLineController) AddLine(ctx server.Context) error {
	req := &m.AddLineReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemLineService.AddSystemLine(req.LineId, req.LineName, req.LimitCost, req.MealId, req.Domain, req.TransType, req.ApiUrl, req.Md5key, req.RsaPubKey, req.RsaPriKey, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询线路信息
func (m SystemLineController) QueryLineOne(ctx server.Context) error {
	req := &m.QueryLineOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	Line, err := SystemLineService.QueryLineOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, Line)
}

// 修改线路
func (m SystemLineController) EidtLine(ctx server.Context) error {
	req := &m.EditLineReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemLineService.EditSystemLine(req.Id, req.LineName, req.LimitCost, req.MealId, req.Domain, req.TransType, req.ApiUrl, req.Md5key, req.RsaPubKey, req.RsaPriKey, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改角色状态
func (m SystemLineController) EidtLineStatus(ctx server.Context) error {
	req := &m.EditLineStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemLineService.EditSystemLineStatus(req.Id, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 查询全部线路id
func (m SystemLineController) QueryAllLineId(ctx server.Context) error {
	lineIds, err := SystemLineService.QueryAllLineId()
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, lineIds)
}
