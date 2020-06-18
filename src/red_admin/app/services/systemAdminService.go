package services

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/utility/uuid"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_admin/app/middleware"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/conf"
	"fmt"
	"regexp"
)

type SystemAdminService struct{}

var (
	SystemAdminBo = new(bo.SystemAdminBo)
	// SystemRoleBo         = new(bo.SystemRoleBo)
	SessionService       = new(middleware.SessionService)
	AgencySessionService = new(middleware.AdminSessionService)
)

// 查询系统管理员列表
func (SystemAdminService) QuerySystemAdminList(startTime, endTime, roleId, isOnline int, account string, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部系统管理员信息
	count, Admins, err := SystemAdminBo.QuerySystemAdminList(sess, startTime, endTime, roleId, isOnline, account, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	// 获取全部角色信息
	roles, err := SystemRoleBo.QuerySystemRoleAll(sess)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	for _, v := range Admins {
		for _, r := range roles {
			if v.RoleId == r.Id {
				v.RoleName = r.RoleName
			}
		}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = Admins
	pageResp.Count = count
	return pageResp, nil
}

// 添加系统管理员
func (SystemAdminService) AddSystemAdmin(account, password, confirmPassword string, roleId int) error {
	if password != confirmPassword {
		return &validate.Err{Code: code.PASSWORD_NOT_SAME}
	}
	// 数字+字母
	if !regexp.MustCompile(`^[0-9a-zA-Z]{5,16}$`).MatchString(account) {
		return &validate.Err{Code: code.USER_NAME_NOT_RIGHT}
	}
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据账号查询管理员信息
	_, has, _ := SystemAdminBo.QueryAdminByAccount(sess, account)
	if has {
		return &validate.Err{Code: code.ACCOUNT_ALREADY_EXISTS}
	}

	systemAdmin := new(structs.SystemAdmin)
	// 添加系统管理员
	systemAdmin.Account = account
	systemAdmin.Password = utility.NewPasswordEncrypt(account, password)
	systemAdmin.RoleId = roleId
	systemAdmin.IsOnline = 2
	systemAdmin.CreateTime = utility.GetNowTimestamp()
	_, err := SystemAdminBo.AddAdmin(sess, systemAdmin)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 根据id查询系统管理员信息
func (SystemAdminService) QueryAdminOne(id int) (*structs.SystemAdminReq, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	admin, has, _ := SystemAdminBo.QueryAdminById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	Admin := new(structs.SystemAdminReq)
	Admin.Id = admin.Id
	Admin.Account = admin.Account
	Admin.RoleId = admin.RoleId
	Admin.IsOnline = admin.IsOnline
	Admin.CreateTime = admin.CreateTime
	Admin.LastIp = admin.LastIp
	Admin.LastLoginTime = admin.LastLoginTime

	return Admin, nil
}

// 修改系统管理员
func (SystemAdminService) EditSystemAdmin(id int, password, confirmPassword string, roleId int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	if password != confirmPassword {
		return &validate.Err{Code: code.PASSWORD_NOT_SAME}
	}
	// 判断系统管理员是否存在
	admin, has, _ := SystemAdminBo.QueryAdminById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	if password != "" {
		admin.Password = utility.NewPasswordEncrypt(admin.Account, password)
	}
	admin.RoleId = roleId
	err := SystemAdminBo.EditAdmin(sess, admin)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 踢线，删除session信息
	err = SessionService.DelSessionId(model.GetAdminListKey(), admin.Id)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}

	// 更新管理员在线状态
	admin.IsOnline = model.OFFLINE
	err = SystemAdminBo.EditAdminOnlineStatus(sess, admin)
	if err != nil {
		golog.Error("SystemAdminService", "login", "err:%v", err)
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 删除管理员
func (SystemAdminService) DelSystemAdmin(id int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断系统管理员是否存在
	Admin, has, _ := SystemAdminBo.QueryAdminById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	Admin.DeleteTime = utility.GetNowTimestamp()
	err := SystemAdminBo.DelAdmin(sess, Admin)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 踢线，删除session信息
	err = SessionService.DelSessionId(model.GetAdminListKey(), id)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}
	return nil
}

// 重置密码
func (SystemAdminService) ResetPassword(id int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断系统管理员是否存在
	admin, has, _ := SystemAdminBo.QueryAdminById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	admin.Password = utility.NewPasswordEncrypt(admin.Account, "123456")
	err := SystemAdminBo.EditAdmin(sess, admin)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}

	// 踢线，删除session信息
	err = SessionService.DelSessionId(model.GetAdminListKey(), id)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}

	// 更新管理员在线状态
	admin.IsOnline = model.OFFLINE
	err = SystemAdminBo.EditAdminOnlineStatus(sess, admin)
	if err != nil {
		golog.Error("SystemAdminService", "login", "err:%v", err)
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 登陆
func (SystemAdminService) Login(account, password, ip string, device int) (*model.AdminSession, error) {
	if account == "pkplusadmin" && password == "#123456" {
		// 超级管理员
		session := new(model.AdminSession)
		user := new(model.AdminUser)
		user.Account = account
		user.Id = -1
		user.RoleId = -1
		user.RoleName = "超级管理员"
		session.User = user
		session.SessionId = fmt.Sprintf("admin_%s_%d_super", uuid.NewV4().String(), device)
		session.IsAdmin = 1
		// 挤线操作
		SessionService.DelSessionId(model.RED_ADMIN_SESSION_LIST_KEY, -1)
		err := SessionService.SaveSession(model.GetAdminListKey(), session)
		if err != nil {
			golog.Error("SystemAdminService", "login", "err:%v", err)
			return nil, err
		}
		return session, nil
	} else {
		// 系统管理员
		sess := conf.GetXormSession()
		defer sess.Close()
		// 根据账号查询管理员信息
		admin, has, _ := SystemAdminBo.QueryAdminByAccount(sess, account)
		if !has {
			return nil, &validate.Err{Code: code.ACCOUNT_DOES_NOT_EXIST}
		}
		// 验证密码
		if utility.NewPasswordEncrypt(admin.Account, password) != admin.Password {
			return nil, &validate.Err{Code: code.PASSWORD_ERRORS}
		}
		// 获取角色信息
		role, has, _ := SystemRoleBo.QueryRoleById(sess, admin.RoleId)
		if !has {
			return nil, &validate.Err{Code: code.ADMIN_ROLE_NOT_EXIST}
		}

		if role.Status == 2 {
			return nil, &validate.Err{Code: code.ADMIN_ROLE_HAS_BEEN_STOPED}
		}

		// 挤线操作
		SessionService.DelSessionId(model.RED_ADMIN_SESSION_LIST_KEY, admin.Id)

		session := new(model.AdminSession)
		session.SessionId = fmt.Sprintf("admin_%s_%d_%d", uuid.NewV4().String(), device, admin.Id)
		user := new(model.AdminUser)
		user.Id = admin.Id
		user.Account = admin.Account
		user.RoleId = admin.RoleId
		user.RoleName = role.RoleName
		user.LastIp = admin.LastIp
		user.LastLoginTime = admin.LastLoginTime
		session.User = user
		session.IsAdmin = 2
		err := SessionService.SaveSession(model.GetAdminListKey(), session)
		if err != nil {
			golog.Error("SystemAdminService", "login", "err:%v", err)
			return nil, err
		}
		// 更新管理员在线状态
		admin.LastIp = ip
		admin.LastLoginTime = utility.GetNowTimestamp()
		admin.IsOnline = model.ONLINE
		err = SystemAdminBo.EditAdminOnlineStatus(sess, admin)
		if err != nil {
			golog.Error("SystemAdminService", "login", "err:%v", err)
			return nil, &validate.Err{Code: code.UPDATE_FAILED}
		}
		return session, nil
	}
	// return nil, nil
}

// 注销
func (SystemAdminService) Logout(user *model.AdminSession) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	if user.IsAdmin == 1 {
		// 删除session信息
		err := SessionService.DelSessionId(model.GetAdminListKey(), user.Uid())
		if err != nil {
			return &validate.Err{Code: code.DELETE_FAILED}
		}
		return nil
	}
	// 验证管理员和是否存在
	admin, has, _ := SystemAdminBo.QueryAdminById(sess, user.Uid())
	if !has {
		return &validate.Err{Code: code.USER_DOES_NOT_EXISTS}
	}

	// 删除session信息
	err := SessionService.DelSessionId(model.GetAdminListKey(), admin.Id)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}

	// 更新管理员在线状态
	admin.IsOnline = model.OFFLINE
	err = SystemAdminBo.EditAdminOnlineStatus(sess, admin)
	if err != nil {
		golog.Error("SystemAdminService", "login", "err:%v", err)
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 修改密码
func (SystemAdminService) EditPassword(account, oldPassword, password, confirmPassword string) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 两次密码不一致
	if password != confirmPassword {
		return &validate.Err{Code: code.PASSWORD_NOT_SAME}
	}

	// 根据账号查询管理员信息
	admin, has, _ := SystemAdminBo.QueryAdminByAccount(sess, account)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}

	if admin.Password != utility.NewPasswordEncrypt(admin.Account, oldPassword) {
		return &validate.Err{Code: code.PASSWORD_ERRORS}
	}

	admin.Password = utility.NewPasswordEncrypt(admin.Account, password)
	err := SystemAdminBo.EditAdmin(sess, admin)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}
