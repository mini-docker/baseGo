package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
)

type SystemAdminBo struct{}

// 返回所有角色列表
func (*SystemAdminBo) QuerySystemAdminList(sess *xorm.Session, startTime, endTime, roleId, isOnline int, account string, page, pageSize int) (int64, []*structs.SystemAdminReq, error) {
	// 返回所有权限
	rows := make([]*structs.SystemAdminReq, 0)
	sess.Where("delete_time = ?", model.UNDEL)
	if startTime != 0 {
		sess.Where("create_time >= ?", startTime)
	}
	if endTime != 0 {
		sess.Where("create_time <= ?", endTime)
	}
	if roleId != 0 {
		sess.Where("role_id = ?", roleId)
	}
	if isOnline != 0 {
		sess.Where("is_online = ? ", isOnline)
	}
	if account != "" {
		sess.Where("account like ?", account+"%")
	}
	count, err := sess.Table(new(structs.SystemAdmin).TableName()).Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&rows)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}

// 添加角色
func (*SystemAdminBo) AddAdmin(sess *xorm.Session, Admin *structs.SystemAdmin) (int64, error) {
	return sess.Insert(Admin)
}

// 修改角色信息
func (*SystemAdminBo) EditAdmin(sess *xorm.Session, Admin *structs.SystemAdmin) error {
	_, err := sess.Table(new(structs.SystemAdmin).TableName()).
		ID(Admin.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("password", "role_id", "update_time").
		Update(Admin)
	return err
}

// 修改角色状态信息
func (*SystemAdminBo) EditAdminStatus(sess *xorm.Session, Admin *structs.SystemAdmin) error {
	_, err := sess.Table(new(structs.SystemAdmin).TableName()).
		ID(Admin.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("status", "update_time").
		Update(Admin)
	return err
}

// 修改管理员登陆状态
func (*SystemAdminBo) EditAdminOnlineStatus(sess *xorm.Session, Admin *structs.SystemAdmin) error {
	_, err := sess.Table(new(structs.SystemAdmin).TableName()).
		ID(Admin.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("is_online", "last_ip", "last_login_time").
		Update(Admin)
	return err
}

// 根据id查询单个用户
func (*SystemAdminBo) QueryAdminById(sess *xorm.Session, id int) (*structs.SystemAdmin, bool, error) {
	Admin := new(structs.SystemAdmin)
	has, err := sess.Where("id = ? and delete_time = ?", id, model.UNDEL).Get(Admin)
	return Admin, has, err
}

// 根据id查询单个用户
func (*SystemAdminBo) QueryAdminByRoleId(sess *xorm.Session, id int) ([]*structs.SystemAdmin, error) {
	Admins := make([]*structs.SystemAdmin, 0)
	err := sess.Table(new(structs.SystemAdmin).TableName()).Where("role_id = ? and delete_time = ?", id, model.UNDEL).Find(&Admins)
	return Admins, err
}

// 根据account查询管理员
func (*SystemAdminBo) QueryAdminByAccount(sess *xorm.Session, account string) (*structs.SystemAdmin, bool, error) {
	admin := new(structs.SystemAdmin)
	has, err := sess.Where("account = ? and delete_time = ?", account, model.UNDEL).Get(admin)
	return admin, has, err
}

// 删除角色
func (*SystemAdminBo) DelAdmin(sess *xorm.Session, admin *structs.SystemAdmin) error {
	_, err := sess.Table(new(structs.SystemAdmin).TableName()).
		ID(admin.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("delete_time").
		Update(admin)
	return err
}

// 添加角色权限
func (*SystemAdminBo) AddPermission(sess *xorm.Session, permissions []*structs.SystemRoleMenu) error {
	_, err := sess.Insert(&permissions)
	return err
}

// 删除角色权限
func (*SystemAdminBo) DelPermission(sess *xorm.Session, id int) error {
	_, err := sess.Table(new(structs.SystemRoleMenu)).
		Where("role_id = ?", id).
		Delete(&structs.SystemRoleMenu{})
	return err
}
