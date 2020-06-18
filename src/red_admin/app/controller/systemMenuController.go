package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

// 菜单控制器
type SystemMenuController struct {
	AddMenuReq struct {
		ParentId int    `json:"parentId" valid:"Must;ErrorCode(3045)` // 父id
		Name     string `json:"name" valid:"Must;ErrorCode(3046)`     // 菜单名称
		Route    string `json:"route" valid:"Must;ErrorCode(3047)`    // 菜单路由
		Icon     string `json:"icon"`                                 // 菜单图标
		Status   int    `json:"status" valid:"Must;ErrorCode(3028)`   // 状态 1 启用 2 停用
		IsShow   int    `json:"isShow"`                               // 是否可见 1可见 2不可见
		Sort     int    `json:"sort"`                                 // 排列序号
	}

	QueryMenuReq struct {
		Id int `json:"id"` // 菜单id
	}

	EditMenuReq struct {
		Id     int    `json:"id"  valid:"Must;ErrorCode(3031)`    // 菜单id
		Name   string `json:"name" valid:"Must;ErrorCode(3046)`   // 菜单名称
		Route  string `json:"route" valid:"Must;ErrorCode(3047)`  // 菜单路由
		Icon   string `json:"icon"`                               // 菜单图标
		Status int    `json:"status" valid:"Must;ErrorCode(3028)` // 状态 1 启用 2 停用
		IsShow int    `json:"isShow"`                             // 是否可见 1可见 2不可见
		Sort   int    `json:"sort"`                               // 排列序号
	}
}

var (
	SystemMenuService = new(services.SystemMenuService)
)

// 查询全部菜单
func (SystemMenuController) QueryMenuList(ctx server.Context) error {
	result, err := SystemMenuService.QuerySystemMenuList()
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加菜单
func (m SystemMenuController) AddMenu(ctx server.Context) error {
	req := &m.AddMenuReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemMenuService.AddSystemMenu(req.ParentId, req.Name, req.Route, req.Icon, req.IsShow, req.Status, req.Sort)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询菜单信息
func (m SystemMenuController) QueryMenuOne(ctx server.Context) error {
	req := &m.QueryMenuReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	menu, err := SystemMenuService.QueryMenuOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, menu)
}

// 修改菜单
func (m SystemMenuController) EidtMenu(ctx server.Context) error {
	req := &m.EditMenuReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemMenuService.EditSystemMenu(req.Id, req.Name, req.Route, req.Icon, req.IsShow, req.Status, req.Sort)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据父级id查询子级菜单
func (m SystemMenuController) QueryChildrenById(ctx server.Context) error {
	req := &m.QueryMenuReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	menus, err := SystemMenuService.QueryChildrenById(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, menus)
}
