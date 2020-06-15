package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

// 套餐控制器
type SystemLineMealController struct {
	AddLineMealReq struct {
		MealName  string  `json:"mealName" valid:"Must;ErrorCode(3042)`  // 套餐名称
		NnRoyalty float64 `json:"nnRoyalty" valid:"Must;ErrorCode(3043)` // 牛牛红包抽成
		SlRoyalty float64 `json:"slRoyalty" valid:"Must;ErrorCode(3044)` // 扫雷红包抽成
	}

	QueryLineMealReq struct {
		PageIndex int `json:"pageIndex"` // 页码
		PageSize  int `json:"pageSize"`  // 每页条数
	}

	QueryLineMealOneReq struct {
		Id int `json:"id"` // 套餐id
	}

	EditLineMealReq struct {
		Id        int     `json:"id" valid:"Must;ErrorCode(3031)`        // 套餐id
		MealName  string  `json:"mealName" valid:"Must;ErrorCode(3042)`  // 套餐名称
		NnRoyalty float64 `json:"nnRoyalty" valid:"Must;ErrorCode(3043)` // 牛牛红包抽成
		SlRoyalty float64 `json:"slRoyalty" valid:"Must;ErrorCode(3044)` // 扫雷红包抽成
	}
}

var (
	SystemLineMealService = new(services.SystemLineMealService)
)

// 查询全部套餐
func (m SystemLineMealController) QueryLineMealList(ctx server.Context) error {
	req := &m.QueryLineMealReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := SystemLineMealService.QuerySystemLineMealList(req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加套餐
func (m SystemLineMealController) AddLineMeal(ctx server.Context) error {
	req := &m.AddLineMealReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemLineMealService.AddSystemLineMeal(req.MealName, req.SlRoyalty, req.NnRoyalty)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询套餐信息
func (m SystemLineMealController) QueryLineMealOne(ctx server.Context) error {
	req := &m.QueryLineMealOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	LineMeal, err := SystemLineMealService.QueryLineMealOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, LineMeal)
}

// 修改套餐
func (m SystemLineMealController) EidtLineMeal(ctx server.Context) error {
	req := &m.EditLineMealReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemLineMealService.EditSystemLineMeal(req.Id, req.MealName, req.SlRoyalty, req.NnRoyalty)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 查询全部套餐id
func (m SystemLineMealController) QueryAllLineMealCode(ctx server.Context) error {
	lineIds, err := SystemLineMealService.QueryAllLineMealCode()
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, lineIds)
}
