package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
)

type SystemMenuBo struct{}

// 返回所有菜单列表
func (*SystemMenuBo) QuerySystemMenuList(sess *xorm.Session) ([]*structs.SystemMenu, error) {
	// 返回所有权限
	rows := make([]*structs.SystemMenu, 0)
	sess.Where("delete_time = ?", model.UNDEL)
	err := sess.Find(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// 添加菜单
func (*SystemMenuBo) AddMenu(sess *xorm.Session, menu *structs.SystemMenu) (int64, error) {
	return sess.Insert(menu)
}

// 修改菜单信息
func (*SystemMenuBo) EditMenu(sess *xorm.Session, menu *structs.SystemMenu) error {
	_, err := sess.Table(new(structs.SystemMenu).TableName()).
		ID(menu.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("name", "route", "icon", "is_show", "sort", "status", "update_time").
		Update(menu)
	return err
}

// 根据id查询单个菜单
func (*SystemMenuBo) QueryMenuById(sess *xorm.Session, id int) (*structs.SystemMenu, bool, error) {
	menu := new(structs.SystemMenu)
	has, err := sess.Where("id = ?", id).Get(menu)
	return menu, has, err
}

// 根据id查询子菜单
func (*SystemMenuBo) QueryChildrenById(sess *xorm.Session, id int) ([]*structs.SystemMenuCode, error) {
	menus := make([]*structs.SystemMenuCode, 0)
	err := sess.Where("parent_id = ?", id).Find(&menus)
	return menus, err
}

// 查询用户菜单
func (*SystemMenuBo) QueryMenusByIds(sess *xorm.Session, menuIds []int) ([]*structs.SystemMenu, error) {
	menus := make([]*structs.SystemMenu, 0)
	sess.Where("status = ? ", 1)
	err := sess.In("id", menuIds).Find(&menus)
	return menus, err
}

// 根据名称查询单个菜单
func (*SystemMenuBo) QueryMenuByName(sess *xorm.Session, name string) (*structs.SystemMenu, bool, error) {
	menu := new(structs.SystemMenu)
	has, err := sess.Where("name = ?", name).Get(menu)
	return menu, has, err
}

// 根据路由查询单个菜单
func (*SystemMenuBo) QueryMenuByRoute(sess *xorm.Session, route string) (*structs.SystemMenu, bool, error) {
	menu := new(structs.SystemMenu)
	has, err := sess.Where("route = ?", route).Get(menu)
	return menu, has, err
}
