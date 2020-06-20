package services

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/utility/uuid"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
	"fmt"
	"regexp"
	"strings"
)

type SystemAgencyService struct{}

var (
	SystemAgencyBo = new(bo.SystemAgencyBo)
	SessionService = new(middleware.SessionService)
)

// 查询代理列表
func (SystemAgencyService) QueryAgencyList(lineId, account string, agencyId string, isOnline, status int, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部代理信息
	count, agencys, err := SystemAgencyBo.QuerySystemAgencyList(sess, lineId, account, agencyId, isOnline, status, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = agencys
	pageResp.Count = count
	return pageResp, nil
}

// 添加代理
func (SystemAgencyService) AddAgency(account, password, confirmPassword, lineId string, agencyId string, status int, whiteIpAddress string, creator string, creatorId int, ip string) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 数字+字母
	if !regexp.MustCompile(`^[0-9a-zA-Z]{5,16}$`).MatchString(account) {
		return &validate.Err{Code: code.USER_NAME_NOT_RIGHT}
	}

	// 两次密码不一致
	if password != confirmPassword {
		return &validate.Err{Code: code.PASSWORD_NOT_SAME}
	}

	// 根据账号查询管理员信息
	_, has, _ := SystemAgencyBo.QueryAgencyByAccount(sess, account)
	if has {
		return &validate.Err{Code: code.ACCOUNT_ALREADY_EXISTS}
	}

	Agency := new(structs.Agency)
	// 添加代理
	Agency.Account = account
	Agency.Password = utility.NewPasswordEncrypt(account, password)
	Agency.LineId = lineId
	Agency.Status = status
	Agency.AgencyId = agencyId
	Agency.IsAdmin = 2
	Agency.IsOnline = 2
	Agency.WhiteIpAddress = whiteIpAddress
	Agency.CreateTime = utility.GetNowTimestamp()
	_, err := SystemAgencyBo.AddAgency(sess, Agency)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = lineId
	log.AgencyId = agencyId
	log.CreatorIp = ip
	log.Remark = fmt.Sprintf("超管%v添加了代理%v", creator, account)
	RedLogBo.AddLog(sess, log)
	return nil
}

// 根据id查询超管信息
func (SystemAgencyService) QueryAgencyOne(id int) (*structs.AgencyReq, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	agency, has, _ := SystemAgencyBo.QueryAgencyById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	agencyReq := new(structs.AgencyReq)
	agencyReq.Id = agency.Id
	agencyReq.AgencyId = agency.AgencyId
	agencyReq.Account = agency.Account
	agencyReq.LineId = agency.LineId
	agencyReq.IsOnline = agency.IsOnline
	agencyReq.IsAdmin = agency.IsAdmin
	agencyReq.Status = agency.Status
	agencyReq.CreateTime = agency.CreateTime
	agencyReq.EditTime = agency.EditTime
	agencyReq.WhiteIpAddress = agency.WhiteIpAddress

	return agencyReq, nil
}

// 修改代理
func (SystemAgencyService) EditAgency(id int, password, confirmPassword string, status int, whiteIpAddress string, creator string, creatorId int, ip string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	if password != confirmPassword {
		return &validate.Err{Code: code.PASSWORD_NOT_SAME}
	}
	// 判断代理是否存在
	Agency, has, _ := SystemAgencyBo.QueryAgencyById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	var remark string
	if password != "" {
		if Agency.Password != utility.NewPasswordEncrypt(Agency.Account, password) {
			remark += "将密码修改了"
		}
		Agency.Password = utility.NewPasswordEncrypt(Agency.Account, password)
	}
	if Agency.Status != status {
		Agency.Status = status
		switch status {
		case 1:
			remark += fmt.Sprintf("将代理状态修改为启用️;")
		case 2:
			remark += fmt.Sprintf("将代理状态修改为停用;")
		}
	}
	if Agency.WhiteIpAddress != whiteIpAddress {
		Agency.WhiteIpAddress = whiteIpAddress
		remark += fmt.Sprintf("将ip白名单修改为%v;", whiteIpAddress)
	}

	Agency.EditTime = utility.GetNowTimestamp()
	err := SystemAgencyBo.EditAgency(sess, Agency)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = Agency.LineId
	log.AgencyId = Agency.AgencyId
	log.CreatorIp = ip
	log.Remark = fmt.Sprintf("超管%v修改了代理%v信息:", creator, Agency.Account) + remark
	RedLogBo.AddLog(sess, log)
	return nil
}

// 修改超管
func (SystemAgencyService) EditAgencyStatus(id int, status int, creator string, creatorId int, ip string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断超管是否存在
	Agency, has, _ := SystemAgencyBo.QueryAgencyById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	var remark string
	if Agency.Status != status {
		Agency.Status = status
		switch status {
		case 1:
			remark += fmt.Sprintf("将代理状态修改为启用;")
		case 2:
			remark += fmt.Sprintf("将代理状态修改为停用;")
		}
	}
	Agency.EditTime = utility.GetNowTimestamp()
	err := SystemAgencyBo.EditAgency(sess, Agency)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = Agency.LineId
	log.AgencyId = Agency.AgencyId
	log.CreatorIp = ip
	log.Remark = fmt.Sprintf("超管%v修改了代理%v信息:", creator, Agency.Account) + remark
	RedLogBo.AddLog(sess, log)
	return nil
}

// 初始化代理密码
func (SystemAgencyService) ResetPassword(id int, creator string, creatorId int, ip string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断系统管理员是否存在
	admin, has, _ := SystemAgencyBo.QueryAgencyById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	admin.Password = utility.NewPasswordEncrypt(admin.Account, "123456")
	err := SystemAgencyBo.EditAgency(sess, admin)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = admin.LineId
	log.AgencyId = admin.AgencyId
	log.CreatorIp = ip
	log.Remark = fmt.Sprintf("超管%v初始化了代理%v的密码", creator, admin.Account)
	RedLogBo.AddLog(sess, log)
	return nil
}

// 修改密码
func (SystemAgencyService) EditPassword(account, oldPassword, password, confirmPassword string, creator string, creatorId int, ip string) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 两次密码不一致
	if password != confirmPassword {
		return &validate.Err{Code: code.PASSWORD_NOT_SAME}
	}

	// 根据账号查询管理员信息
	agency, has, _ := SystemAgencyBo.QueryAgencyByAccount(sess, account)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}

	if agency.Password != utility.NewPasswordEncrypt(agency.Account, oldPassword) {
		return &validate.Err{Code: code.PASSWORD_ERRORS}
	}

	agency.Password = utility.NewPasswordEncrypt(agency.Account, password)
	err := SystemAgencyBo.EditAgency(sess, agency)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = agency.LineId
	log.AgencyId = agency.AgencyId
	log.CreatorIp = ip
	log.Remark = fmt.Sprintf("代理%v修改了登录密码", creator)
	RedLogBo.AddLog(sess, log)
	return nil
}

// 删除代理
func (SystemAgencyService) DelAgency(id int, creator string, creatorId int, ip string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断系统管理员是否存在
	admin, has, _ := SystemAgencyBo.QueryAgencyById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	admin.DeleteTime = utility.GetNowTimestamp()
	err := SystemAgencyBo.DelAgency(sess, admin)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = admin.LineId
	log.AgencyId = admin.AgencyId
	log.CreatorIp = ip
	log.Remark = fmt.Sprintf("超管%v删除了代理%v", creator, admin.Account)
	RedLogBo.AddLog(sess, log)
	return nil
}

// 登陆
func (SystemAgencyService) Login(account, password, ip string, device int) (*model.AgencySession, error) {

	// 数据库连接
	sess := conf.GetXormSession()
	defer sess.Close()
	// 根据账号查询管理员信息
	agency, has, _ := SystemAgencyBo.QueryAgencyByAccount(sess, account)
	if !has {
		return nil, &validate.Err{Code: code.ACCOUNT_DOES_NOT_EXIST}
	}
	if agency.DeleteTime != 0 {
		return nil, &validate.Err{Code: code.ACCOUNT_DOES_NOT_EXIST}
	}

	// 验证白名单
	if !strings.Contains(agency.WhiteIpAddress, ip) {
		return nil, &validate.Err{Code: code.LOGINIP_NOT_IN_WHITE_IP_ADDRESS}
	}

	// 验证状态
	if agency.Status == model.MENU_TWO {
		return nil, &validate.Err{Code: code.ACCOUNT_CAN_NOT_BE_LOGIN}
	}

	// 验证密码
	if utility.NewPasswordEncrypt(agency.Account, password) != agency.Password {
		return nil, &validate.Err{Code: code.PASSWORD_ERRORS}
	}

	// 挤线操作
	SessionService.DelSessionId(model.RED_AGENCY_SESSION_LIST_KEY, agency.Id)

	session := new(model.AgencySession)
	session.SessionId = fmt.Sprintf("agency_%s_%d_%d", uuid.NewV4().String(), device, agency.Id)
	user := new(model.AgencyUser)
	user.Id = agency.Id
	user.Account = agency.Account
	user.LineId = agency.LineId
	if agency.IsAdmin == 1 {
		user.AgencyId = ""
	} else {
		user.AgencyId = agency.AgencyId
	}
	session.User = user
	session.IsAdmin = agency.IsAdmin
	err := SessionService.SaveSession(model.GetAgencyListKey(), session)
	if err != nil {
		golog.Error("SystemAdminService", "login", "err:%v", err)
		return nil, err
	}
	// 更新管理员在线状态
	agency.IsOnline = model.ONLINE
	err = SystemAgencyBo.EditAgencyOnlineStatus(sess, agency)
	if err != nil {
		golog.Error("SystemAdminService", "login", "err:%v", err)
		return nil, &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入登录日志
	log := new(structs.RedLog)
	log.Creator = agency.Account
	log.CreatorId = agency.Id
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_LOGIN
	log.LineId = agency.LineId
	log.AgencyId = agency.AgencyId
	log.CreatorIp = ip
	if agency.IsAdmin == 1 {
		log.Remark = fmt.Sprintf("超管%v登录了系统", agency.Account)
	} else {
		log.Remark = fmt.Sprintf("代理%v登录了系统", agency.Account)
	}
	RedLogBo.AddLog(sess, log)
	return session, nil
}

// 注销
func (SystemAgencyService) Logout(user *model.AgencySession) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 验证管理员和是否存在
	agency, has, _ := SystemAgencyBo.QueryAgencyById(sess, user.Uid())
	if !has {
		return &validate.Err{Code: code.USER_DOES_NOT_EXISTS}
	}

	// 删除session信息
	err := SessionService.DelSessionId(model.GetAdminListKey(), agency.Id)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}

	// 更新管理员在线状态
	agency.IsOnline = model.OFFLINE
	err = SystemAgencyBo.EditAgencyOnlineStatus(sess, agency)
	if err != nil {
		golog.Error("SystemAdminService", "login", "err:%v", err)
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}
