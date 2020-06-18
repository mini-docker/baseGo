package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/conf"
)

type SystemRoleService struct{}

var (
	SystemRoleBo = new(bo.SystemRoleBo)
)

// 查询角色列表
func (SystemRoleService) QuerySystemRoleList(roleName string, status int, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部角色信息
	count, roles, err := SystemRoleBo.QuerySystemRoleList(sess, roleName, status, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = roles
	pageResp.Count = count
	return pageResp, nil
}

// 添加角色
func (SystemRoleService) AddSystemRole(roleName, remark string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	systemRole := new(structs.SystemRole)
	// 添加角色
	systemRole.RoleName = roleName
	systemRole.Remark = remark
	systemRole.Status = 1
	systemRole.IsDefault = 2
	systemRole.CreateTime = utility.GetNowTimestamp()
	_, err := SystemRoleBo.AddRole(sess, systemRole)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 根据id查询角色信息
func (SystemRoleService) QueryRoleOne(id int) (*structs.SystemRole, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	role, has, _ := SystemRoleBo.QueryRoleById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	return role, nil
}

// 修改角色
func (SystemRoleService) EditSystemRole(id int, roleName, remark string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断角色是否存在
	role, has, _ := SystemRoleBo.QueryRoleById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	if role.IsDefault == 1 {
		return &validate.Err{Code: code.DEFAULT_ROLE_CAN_NOT_BE_UPDATE}
	}
	role.RoleName = roleName
	role.Remark = remark
	role.EditTime = utility.GetNowTimestamp()
	err := SystemRoleBo.EditRole(sess, role)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 修改角色状态
func (SystemRoleService) EditSystemRoleStatus(id, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断角色是否存在
	role, has, _ := SystemRoleBo.QueryRoleById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	if role.IsDefault == 1 {
		return &validate.Err{Code: code.DEFAULT_ROLE_CAN_NOT_BE_UPDATE}
	}
	role.Status = status
	role.EditTime = utility.GetNowTimestamp()
	err := SystemRoleBo.EditRoleStatus(sess, role)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 角色停用后踢线所有该角色用户
	admins, err := SystemAdminBo.QueryAdminByRoleId(sess, role.Id)
	if err != nil {
		return &validate.Err{Code: code.QUERY_FAILED}
	}
	for _, v := range admins {
		AgencySessionService.DelSessionId(model.RED_AGENCY_SESSION_LIST_KEY, v.Id)
	}
	return nil
}

func (SystemRoleService) DelSystemRole(id int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断角色是否存在
	role, has, _ := SystemRoleBo.QueryRoleById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	if role.IsDefault == 1 {
		return &validate.Err{Code: code.DEFAULT_ROLE_CAN_NOT_BE_DELETE}
	}
	role.DeleteTime = utility.GetNowTimestamp()
	err := SystemRoleBo.DelRole(sess, role)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 角色赋权
func (SystemRoleService) SetRolePermission(id int, menuIds []int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断角色是否存在
	_, has, _ := SystemRoleBo.QueryRoleById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	// 删除原来的权限
	err := SystemRoleBo.DelPermission(sess, id)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}
	// 保存新的角色权限
	permissions := make([]*structs.SystemRoleMenu, 0)
	for _, v := range menuIds {
		permission := new(structs.SystemRoleMenu)
		permission.RoleId = id
		permission.MenuId = v
		permissions = append(permissions, permission)
	}
	if len(permissions) > 0 {
		err := SystemRoleBo.AddPermission(sess, permissions)
		if err != nil {
			return &validate.Err{Code: code.INSET_ERROR}
		}
	}
	return nil
}

// 获取角色权限
func (SystemRoleService) QueryRolePermission(id int) ([]int, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据ID查询权限
	permissions, err := SystemRoleBo.QueryRolePermission(sess, id)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	var permissionList []int
	for _, p := range permissions {
		permissionList = append(permissionList, p.MenuId)
	}
	return permissionList, nil
}

// 获取角色权限菜单
func (SystemRoleService) QueryRoleMenu(id int) ([]*structs.SystemMenuResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据ID查询权限
	permissions, err := SystemRoleBo.QueryRolePermission(sess, id)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	// 获取菜单id
	var menuIds []int
	if len(permissions) > 0 {
		for _, v := range permissions {
			menuIds = append(menuIds, v.MenuId)
		}
	}
	// // 获取用户菜单
	// menus, err := SystemMenuBo.QueryMenusByIds(sess, menuIds)
	// if err != nil {
	// 	return nil, &validate.Err{Code: code.QUERY_FAILED}
	// }
	// // 整理菜单信息
	// var fristMenu, secondMenu, thirdMenu = []*structs.SystemMenuResp{}, []*structs.SystemMenuResp{}, []*structs.SystemMenuResp{}

	// // 菜单整理排序
	// if len(menus) > 0 {
	// 	for _, v := range menus {
	// 		// 整理一级菜单
	// 		if v.ParentId == 0 && v.Level == model.MENU_ONE {
	// 			menuResp := new(structs.SystemMenuResp)
	// 			menuResp.Id = v.Id
	// 			menuResp.ParentId = v.ParentId
	// 			menuResp.Level = v.Level
	// 			menuResp.Name = v.Name
	// 			menuResp.Icon = v.Icon
	// 			menuResp.Route = v.Route
	// 			menuResp.IsShow = v.IsShow
	// 			menuResp.Sort = v.Sort
	// 			menuResp.Status = v.Status
	// 			fristMenu = append(fristMenu, menuResp)
	// 		}
	// 		// 整理二级菜单
	// 		if v.ParentId != 0 && v.Level == model.MENU_TWO {
	// 			menuResp := new(structs.SystemMenuResp)
	// 			menuResp.Id = v.Id
	// 			menuResp.ParentId = v.ParentId
	// 			menuResp.Level = v.Level
	// 			menuResp.Name = v.Name
	// 			menuResp.Icon = v.Icon
	// 			menuResp.Route = v.Route
	// 			menuResp.IsShow = v.IsShow
	// 			menuResp.Sort = v.Sort
	// 			menuResp.Status = v.Status
	// 			secondMenu = append(secondMenu, menuResp)
	// 		}
	// 		// 整理三级菜单
	// 		if v.ParentId != 0 && v.Level == model.MENU_THREE {
	// 			menuResp := new(structs.SystemMenuResp)
	// 			menuResp.Id = v.Id
	// 			menuResp.ParentId = v.ParentId
	// 			menuResp.Level = v.Level
	// 			menuResp.Name = v.Name
	// 			menuResp.Icon = v.Icon
	// 			menuResp.Route = v.Route
	// 			menuResp.IsShow = v.IsShow
	// 			menuResp.Sort = v.Sort
	// 			menuResp.Status = v.Status
	// 			thirdMenu = append(thirdMenu, menuResp)
	// 		}
	// 	}

	// 	// 封装三级菜单到二级菜单子级
	// 	if len(thirdMenu) > 0 && len(secondMenu) > 0 {
	// 		for _, s := range secondMenu {
	// 			for _, t := range thirdMenu {
	// 				if s.Id == t.ParentId {
	// 					s.Children = append(s.Children, t)
	// 				}
	// 			}
	// 			if len(s.Children) > 0 {
	// 				// 三级菜单排序
	// 				sort.Slice(s.Children, func(i, j int) bool {
	// 					if s.Children[i].Sort == s.Children[j].Sort {
	// 						return s.Children[i].Id > s.Children[j].Id
	// 					}
	// 					return s.Children[i].Sort > s.Children[j].Sort
	// 				})
	// 			}
	// 		}
	// 	}

	// 	// 封装二级菜单到一级菜单
	// 	if len(fristMenu) > 0 && len(secondMenu) > 0 {
	// 		for _, f := range fristMenu {
	// 			for _, s := range secondMenu {
	// 				if f.Id == s.ParentId {
	// 					f.Children = append(f.Children, s)
	// 				}
	// 			}
	// 			if len(f.Children) > 0 {
	// 				// 三级菜单排序
	// 				sort.Slice(f.Children, func(i, j int) bool {
	// 					if f.Children[i].Sort == f.Children[j].Sort {
	// 						return f.Children[i].Id > f.Children[j].Id
	// 					}
	// 					return f.Children[i].Sort > f.Children[j].Sort
	// 				})
	// 			}
	// 		}
	// 	}
	// 	return fristMenu, nil
	// }
	return nil, nil
}

// 角色枚举
func (SystemRoleService) QuerySystemRoleCode() ([]*structs.SystemRoleCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	roles, err := SystemRoleBo.QuerySystemRoleCode(sess)
	return roles, err
}
