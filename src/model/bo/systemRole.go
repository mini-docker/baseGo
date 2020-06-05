package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
)

type SystemRoleBo struct{}

// 返回所有角色列表
func (*SystemRoleBo) QuerySystemRoleList(sess *xorm.Session, roleName string, status int, page, pageSize int) (int64, []*structs.SystemRole, error) {
	// 返回所有权限
	rows := make([]*structs.SystemRole, 0)
	if roleName != "" {
		sess.Where("role_name like ?", roleName+"%")
	}
	if status != 0 {
		sess.Where("status = ?", status)
	}
	sess.Where("delete_time = ?", model.UNDEL)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&rows)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}

// 添加角色
func (*SystemRoleBo) AddRole(sess *xorm.Session, role *structs.SystemRole) (int64, error) {
	return sess.Insert(role)
}

// 修改角色信息
func (*SystemRoleBo) EditRole(sess *xorm.Session, role *structs.SystemRole) error {
	_, err := sess.Table(new(structs.SystemRole).TableName()).
		ID(role.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("role_name", "remark", "edit_time").
		Update(role)
	return err
}

// 修改角色状态信息
func (*SystemRoleBo) EditRoleStatus(sess *xorm.Session, role *structs.SystemRole) error {
	_, err := sess.Table(new(structs.SystemRole).TableName()).
		ID(role.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("status", "edit_time").
		Update(role)
	return err
}

// 根据id查询单个角色
func (*SystemRoleBo) QueryRoleById(sess *xorm.Session, id int) (*structs.SystemRole, bool, error) {
	role := new(structs.SystemRole)
	has, err := sess.Where("id = ? and delete_time = ?", id, model.UNDEL).Get(role)
	return role, has, err
}

// 删除角色
func (*SystemRoleBo) DelRole(sess *xorm.Session, role *structs.SystemRole) error {
	_, err := sess.Table(new(structs.SystemRole).TableName()).
		ID(role.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("delete_time").
		Update(role)
	return err
}

// 添加角色权限
func (*SystemRoleBo) AddPermission(sess *xorm.Session, permissions []*structs.SystemRoleMenu) error {
	_, err := sess.Insert(&permissions)
	return err
}

// 查询角色权限
func (*SystemRoleBo) QueryRolePermission(sess *xorm.Session, id int) ([]*structs.SystemRoleMenu, error) {
	permissions := make([]*structs.SystemRoleMenu, 0)
	err := sess.Where("role_id = ?", id).Find(&permissions)
	return permissions, err
}

// 删除角色权限
func (*SystemRoleBo) DelPermission(sess *xorm.Session, id int) error {
	_, err := sess.Table(new(structs.SystemRoleMenu)).
		Where("role_id = ?", id).
		Delete(&structs.SystemRoleMenu{})
	return err
}

// 获取角色枚举(不包括停用)
func (*SystemRoleBo) QuerySystemRoleCode(sess *xorm.Session) ([]*structs.SystemRoleCode, error) {
	roles := make([]*structs.SystemRoleCode, 0)
	err := sess.Table(new(structs.SystemRole).TableName()).Where("delete_time = ? and status = 1", model.UNDEL).Find(&roles)
	return roles, err
}

// 获取全部角色(包括停用)
func (*SystemRoleBo) QuerySystemRoleAll(sess *xorm.Session) ([]*structs.SystemRoleCode, error) {
	roles := make([]*structs.SystemRoleCode, 0)
	err := sess.Table(new(structs.SystemRole).TableName()).Where("delete_time = ?", model.UNDEL).Find(&roles)
	return roles, err
}
