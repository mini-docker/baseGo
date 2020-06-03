package services

import (
	"fecho/utility"
	"model"
	"model/bo"
	"model/code"
	"model/structs"
	"red_admin/app/middleware/validate"
	"red_admin/conf"
	"regexp"
)

type SystemAgencyService struct{}

var (
	AgencyBo        = new(bo.SystemAgencyBo)
	RedPacketSiteBo = new(bo.RedPacketSite)
)

// 查询超管列表
func (SystemAgencyService) QueryAgencyAdminList(lineId, account string, isOnline int, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部超管信息
	count, agencys, err := AgencyBo.QuerySystemAgencyAdminList(sess, lineId, account, isOnline, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = agencys
	pageResp.Count = count
	return pageResp, nil
}

// 查询代理列表
func (SystemAgencyService) QueryAgencyList(lineId, account, agencyId string, isOnline, status int, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部超管信息
	count, agencys, err := AgencyBo.QuerySystemAgencyList(sess, lineId, account, agencyId, isOnline, status, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = agencys
	pageResp.Count = count
	return pageResp, nil
}

// 添加超管
func (SystemAgencyService) AddAgency(account, password, confirmPassword, lineId string, status int, whiteIpAddress string) error {
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

	// 判断账号是否存在
	_, has, _ := AgencyBo.QueryAgencyByAccount(sess, account)
	if has {
		return &validate.Err{Code: code.ACCOUNT_ALREADY_EXISTS}
	}

	Agency := new(structs.Agency)
	// 添加超管
	Agency.Account = account
	Agency.Password = utility.NewPasswordEncrypt(account, password)
	Agency.LineId = lineId
	Agency.Status = status
	Agency.IsAdmin = 1
	Agency.IsOnline = 2
	Agency.CreateTime = utility.GetNowTimestamp()
	Agency.WhiteIpAddress = whiteIpAddress
	_, err := AgencyBo.AddAgency(sess, Agency)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 根据id查询超管信息
func (SystemAgencyService) QueryAgencyOne(id int) (*structs.AgencyReq, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	agency, has, _ := AgencyBo.QueryAgencyById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	agencyReq := new(structs.AgencyReq)
	agencyReq.Id = agency.Id
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

// 修改超管
func (SystemAgencyService) EditAgency(id int, password, confirmPassword string, status int, whiteIpAddress string) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 两次密码不一致
	if password != confirmPassword {
		return &validate.Err{Code: code.PASSWORD_NOT_SAME}
	}

	// 判断超管是否存在
	Agency, has, _ := AgencyBo.QueryAgencyById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	if password != "" {
		Agency.Password = utility.NewPasswordEncrypt(Agency.Account, password)
	}
	Agency.Status = status
	Agency.EditTime = utility.GetNowTimestamp()
	Agency.WhiteIpAddress = whiteIpAddress
	err := AgencyBo.EditAgency(sess, Agency)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	if status == 2 {
		// 超管踢线
		AgencySessionService.DelSessionId(model.RED_AGENCY_SESSION_LIST_KEY, id)
		Agency.IsOnline = 2
		AgencyBo.EditAgencyOnlineStatus(sess, Agency)
	}
	return nil
}

// 修改超管
func (SystemAgencyService) EditAgencyStatus(id int, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断超管是否存在
	Agency, has, _ := AgencyBo.QueryAgencyById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	Agency.Status = status
	Agency.EditTime = utility.GetNowTimestamp()
	err := AgencyBo.EditAgency(sess, Agency)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	if status == 2 {
		// 停用踢线超管
		SessionService.DelSessionId(model.RED_AGENCY_SESSION_LIST_KEY, Agency.Id)
	}
	return nil
}

// 获取代理code
func (SystemAgencyService) SiteCode(lineId string) ([]*structs.AgencyCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	agencys, err := RedPacketSiteBo.SiteCode(sess, lineId)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return agencys, err
}
