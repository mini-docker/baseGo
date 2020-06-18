package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/conf"
	"fmt"

	"sort"
)

type SystemMenuService struct{}

var (
	SystemMenuBo = new(bo.SystemMenuBo)
)

// 查询菜单列表
func (SystemMenuService) QuerySystemMenuList() ([]*structs.SystemMenuResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部菜单信息
	menus, err := SystemMenuBo.QuerySystemMenuList(sess)
	// todo 直接获取该值
	fmt.Println("menus", menus)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	var fristMenu, secondMenu, thirdMenu = []*structs.SystemMenuResp{}, []*structs.SystemMenuResp{}, []*structs.SystemMenuResp{}
	// 菜单整理排序
	if len(menus) > 0 {
		for _, v := range menus {
			// 整理一级菜单
			if v.ParentId == 0 && v.Level == model.MENU_ONE {
				menuResp := new(structs.SystemMenuResp)
				menuResp.Id = v.Id
				menuResp.ParentId = v.ParentId
				menuResp.Level = v.Level
				menuResp.Name = v.Name
				menuResp.Icon = v.Icon
				menuResp.Route = v.Route
				menuResp.IsShow = v.IsShow
				menuResp.Sort = v.Sort
				menuResp.Status = v.Status
				fristMenu = append(fristMenu, menuResp)
			}
			// 整理二级菜单
			if v.ParentId != 0 && v.Level == model.MENU_TWO {
				menuResp := new(structs.SystemMenuResp)
				menuResp.Id = v.Id
				menuResp.ParentId = v.ParentId
				menuResp.Level = v.Level
				menuResp.Name = v.Name
				menuResp.Icon = v.Icon
				menuResp.Route = v.Route
				menuResp.IsShow = v.IsShow
				menuResp.Sort = v.Sort
				menuResp.Status = v.Status
				secondMenu = append(secondMenu, menuResp)
			}
			// 整理三级菜单
			if v.ParentId != 0 && v.Level == model.MENU_THREE {
				menuResp := new(structs.SystemMenuResp)
				menuResp.Id = v.Id
				menuResp.ParentId = v.ParentId
				menuResp.Level = v.Level
				menuResp.Name = v.Name
				menuResp.Icon = v.Icon
				menuResp.Route = v.Route
				menuResp.IsShow = v.IsShow
				menuResp.Sort = v.Sort
				menuResp.Status = v.Status
				thirdMenu = append(thirdMenu, menuResp)
			}
		}

		// 封装三级菜单到二级菜单子级
		if len(thirdMenu) > 0 && len(secondMenu) > 0 {
			for _, s := range secondMenu {
				for _, t := range thirdMenu {
					if s.Id == t.ParentId {
						s.Children = append(s.Children, t)
					}
				}
				if len(s.Children) > 0 {
					// 三级菜单排序
					sort.Slice(s.Children, func(i, j int) bool {
						if s.Children[i].Sort == s.Children[j].Sort {
							return s.Children[i].Id > s.Children[j].Id
						}
						return s.Children[i].Sort > s.Children[j].Sort
					})
				}
			}
		}

		// 封装二级菜单到一级菜单
		if len(fristMenu) > 0 && len(secondMenu) > 0 {
			for _, f := range fristMenu {
				for _, s := range secondMenu {
					if f.Id == s.ParentId {
						f.Children = append(f.Children, s)
					}
				}
				if len(f.Children) > 0 {
					// 菜单排序
					sort.Slice(f.Children, func(i, j int) bool {
						if f.Children[i].Sort == f.Children[j].Sort {
							return f.Children[i].Id > f.Children[j].Id
						}
						return f.Children[i].Sort > f.Children[j].Sort
					})
				}
			}
		}
		fmt.Println(fristMenu, "fristMenu")
		return fristMenu, nil
	}
	return nil, nil
}

// 添加菜单
func (SystemMenuService) AddSystemMenu(parentId int, name, route, icon string, isShow, status, sort int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	systemMenu := new(structs.SystemMenu)
	if parentId != 0 {
		// 判断父节点是否存在
		parent, has, _ := SystemMenuBo.QueryMenuById(sess, parentId)
		if !has {
			return &validate.Err{Code: code.PARENT_MENU_NOT_EXIST}
		}
		// 父节点不能为三级
		if parent.Level == model.MENU_THREE {
			return &validate.Err{Code: code.PARENT_MENU_NOT_RIGHT}
		}
		// 设置本节点层级
		systemMenu.Level = parent.Level + 1
	}
	// 判断菜单名称是否存在
	_, has, _ := SystemMenuBo.QueryMenuByName(sess, name)
	if has {
		return &validate.Err{Code: code.MENU_NAME_EXIST}
	}
	// 判断菜单路由是否存在
	if route != "#" {
		_, has, _ = SystemMenuBo.QueryMenuByRoute(sess, route)
		if has {
			return &validate.Err{Code: code.MENU_ROUTE_EXIST}
		}
	}
	// 添加菜单
	systemMenu.ParentId = parentId
	systemMenu.Name = name
	systemMenu.Route = route
	systemMenu.Icon = icon
	systemMenu.IsShow = isShow
	systemMenu.Status = status
	systemMenu.Sort = sort
	systemMenu.CreateTime = utility.GetNowTimestamp()
	_, err := SystemMenuBo.AddMenu(sess, systemMenu)
	if err != nil {
		fmt.Println(err, "error")
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 根据id查询菜单信息
func (SystemMenuService) QueryMenuOne(id int) (*structs.SystemMenu, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	menu, has, _ := SystemMenuBo.QueryMenuById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	return menu, nil
}

// 根据id查询子菜单
func (SystemMenuService) QueryChildrenById(id int) ([]*structs.SystemMenuCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	menus, err := SystemMenuBo.QueryChildrenById(sess, id)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return menus, nil
}

// 修改菜单
func (SystemMenuService) EditSystemMenu(id int, name, route, icon string, isShow, status, sort int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断菜单是否存在
	menu, has, _ := SystemMenuBo.QueryMenuById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	// 判断菜单名称是否存在
	if menu.Name != name {
		_, has, _ := SystemMenuBo.QueryMenuByName(sess, name)
		if has {
			return &validate.Err{Code: code.MENU_NAME_EXIST}
		}
	}
	if menu.Route != route && menu.Route != "#" {
		// 判断菜单路由是否存在
		_, has, _ := SystemMenuBo.QueryMenuByRoute(sess, route)
		if has {
			return &validate.Err{Code: code.MENU_ROUTE_EXIST}
		}
	}

	menu.Name = name
	menu.Route = route
	menu.Icon = icon
	menu.IsShow = isShow
	menu.Status = status
	menu.Sort = sort
	menu.UpdateTime = utility.GetNowTimestamp()
	err := SystemMenuBo.EditMenu(sess, menu)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}
