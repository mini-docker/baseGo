package controller

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
	"strconv"
	"strings"
)

// 角色控制器
type SystemRoleController struct {
	AddRoleReq struct {
		RoleName string `json:"roleName" valid:"Must;ErrorCode(3048)` // 角色名称
		Remark   string `json:"remark"`                               // 角色备注
	}

	QueryRoleOneReq struct {
		Id int `json:"id"` // 角色id
	}

	QueryRoleReq struct {
		RoleName  string `json:"roleName"`  // 角色名称
		Status    int    `json:"status"`    // 状态 1 启用 2 停用
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
	}

	EditRoleReq struct {
		Id       int    `json:"id" valid:"Must;ErrorCode(3031)`       // 角色id
		RoleName string `json:"roleName" valid:"Must;ErrorCode(3045)` // 角色名称
		Remark   string `json:"remark"`                               // 角色备注
	}

	EditRoleStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)`     // 角色id
		Status int `json:"status" valid:"Must;ErrorCode(3028)` // 状态 1 启用 2 停用
	}

	DelRoleReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)` // 角色id
	}

	SetRolePermissionReq struct {
		Id      int    `json:"id" valid:"Must;ErrorCode(3031)`      // 角色id
		MenuIds string `json:"menuIds" valid:"Must;ErrorCode(3031)` // 菜单ids
	}
}

var (
	SystemRoleService = new(services.SystemRoleService)
)

// 查询全部角色
func (m SystemRoleController) QueryRoleList(ctx server.Context) error {
	req := &m.QueryRoleReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := SystemRoleService.QuerySystemRoleList(req.RoleName, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加角色
func (m SystemRoleController) AddRole(ctx server.Context) error {
	req := &m.AddRoleReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemRoleService.AddSystemRole(req.RoleName, req.Remark)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询角色信息
func (m SystemRoleController) QueryRoleOne(ctx server.Context) error {
	req := &m.QueryRoleOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	Role, err := SystemRoleService.QueryRoleOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, Role)
}

// 修改角色
func (m SystemRoleController) EidtRole(ctx server.Context) error {
	req := &m.EditRoleReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemRoleService.EditSystemRole(req.Id, req.RoleName, req.Remark)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改角色状态
func (m SystemRoleController) EidtRoleStatus(ctx server.Context) error {
	req := &m.EditRoleStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemRoleService.EditSystemRoleStatus(req.Id, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 删除角色
func (m SystemRoleController) DelRole(ctx server.Context) error {
	req := &m.DelRoleReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemRoleService.DelSystemRole(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 角色赋权
func (m SystemRoleController) SetRolePermission(ctx server.Context) error {
	req := &m.SetRolePermissionReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	var menuIds []int
	if len(req.MenuIds) > 0 {
		str := strings.Split(req.MenuIds, ",")
		for _, v := range str {
			if "" != v {
				id, _ := strconv.Atoi(v)
				menuIds = append(menuIds, id)
			}
		}

	}
	if len(menuIds) > 0 {
		err := SystemRoleService.SetRolePermission(req.Id, menuIds)
		if err != nil {
			return common.HttpResultJsonError(ctx, err)
		}
	}
	return common.HttpResultJson(ctx, nil)
}

// 获取角色菜单
func (m SystemRoleController) QueryRoleMenu(ctx server.Context) error {
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	// if user.IsAdmin == 1 {
	// 	// 超级管理员，获取全部菜单
	// 	menus, err := SystemMenuService.QuerySystemMenuList()
	// 	if err != nil {
	// 		return common.HttpResultJsonError(ctx, err)
	// 	}
	// 	return common.HttpResultJson(ctx, menus)
	// }
	// 获取用户菜单
	menus, err := SystemRoleService.QueryRoleMenu(user.User.RoleId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, menus)
}

// 获取角色菜单
func (m SystemRoleController) QueryRolePermission(ctx server.Context) error {
	req := &m.QueryRoleOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户权限
	menus, err := SystemRoleService.QueryRolePermission(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, menus)
}

// 角色枚举
func (SystemRoleController) QuerySystemRoleCode(ctx server.Context) error {
	roles, err := SystemRoleService.QuerySystemRoleCode()
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, roles)
}
