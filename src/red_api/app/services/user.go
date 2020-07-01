package services

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/utility/uuid"
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_api/app/middleware"
	"baseGo/src/red_api/app/middleware/validate"
	"baseGo/src/red_api/conf"
	"fmt"
)

type UserService struct{}

var (
	sessionService = new(middleware.SessionService)
)

//func GetLineAccount(lineId string, agencyId string, account string) string {
//	return fmt.Sprintf("%s_%s_%s", lineId, agencyId, account)
//}

// 登陆
func (ms UserService) AttemptLogin(lineId string, agencyId string, userAccount, password, ip string, device int) (*model.UserSession, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	//userAccount = GetLineAccount(lineId, fmt.Sprint(agencyId), userAccount)
	// check exist
	password = utility.NewPasswordEncrypt(userAccount, password)
	agencys, err := AgencyBo.FindSiteByAgencyId(sess, lineId, agencyId)
	if err != nil || agencys == nil || len(agencys) == 0 {
		return nil, &validate.Err{Code: code.LOGIN_FAIL}
	}
	has, user := UserBo.GetOneByAccount(sess, lineId, agencyId, userAccount)

	if !has {
		return nil, &validate.Err{Code: code.ACCOUNT_DOES_NOT_EXIST}
	}
	// check password
	if user.Password != password {
		return nil, &validate.Err{Code: code.ACCOUNT_PASSWORD_ERR}
	}
	if user.Status != 1 {
		return nil, &validate.Err{Code: code.ACCOUNT_DISABLED}
	}

	// store new session
	sessionIdFull := fmt.Sprintf("%s_%d_%d", uuid.NewV4().String(), user.Id, device)
	session := &model.UserSession{
		SessionId: sessionIdFull,
		User: &model.User{
			Id:            user.Id,
			Account:       userAccount,
			LastIp:        user.LastLoginIp,
			LastLoginTime: user.LastLoginTime,
			LineId:        lineId,
			AgencyId:      user.AgencyId,
		},
		TimeOut:      0,
		IsKeepOnline: true,
	}
	err = sessionService.SaveSession(model.GetMemberListKey(), session)
	if err != nil {
		return nil, &validate.Err{Code: code.LOGIN_FAIL}
	}

	user.IsOnline = 2
	user.LastLoginIp = ip
	user.LastLoginTime = utility.GetNowTimestamp()
	err = UserBo.UpdateLoginIp(sess, user)
	if err != nil {
		return nil, &validate.Err{Code: code.LOGIN_FAIL}
	}

	return session, nil
}

// 注册
func (ms UserService) Register(lineId string, agencyId string, userAccount, password, ip string, device int) (*model.UserSession, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	//userAccount = GetLineAccount(lineId, agencyId, userAccount)
	// check exist、
	password = utility.NewPasswordEncrypt(userAccount, password)
	agencys, err := AgencyBo.FindSiteByAgencyId(sess, lineId, agencyId)
	if err != nil || agencys == nil || len(agencys) == 0 {
		return nil, &validate.Err{Code: code.REGISTER_ERROR}
	}

	has, _ := UserBo.GetOneByAccount(sess, lineId, agencyId, userAccount)

	if has {
		return nil, &validate.Err{Code: code.ACCOUNT_ALREADY_EXISTS}
	}
	nowTime := utility.GetNowTimestamp()
	user := &structs.User{
		LineId:        lineId,
		AgencyId:      agencyId,
		Account:       userAccount,
		Password:      password,
		IsOnline:      2,
		Ip:            ip,
		Status:        1,
		CreateTime:    nowTime,
		LastLoginIp:   ip,
		LastLoginTime: nowTime,
	}
	err = UserBo.SaveUser(sess, user)

	if err != nil {
		return nil, &validate.Err{Code: code.ADD_FAILED}
	}

	// store new session
	sessionIdFull := fmt.Sprintf("%s_%d_%d", uuid.NewV4().String(), user.Id, device)
	session := &model.UserSession{
		SessionId: sessionIdFull,
		User: &model.User{
			Id:            user.Id,
			Account:       userAccount,
			LastIp:        user.LastLoginIp,
			LastLoginTime: user.LastLoginTime,
			LineId:        lineId,
			AgencyId:      user.AgencyId,
		},
		TimeOut:      0,
		IsKeepOnline: true,
	}
	err = sessionService.SaveSession(model.GetMemberListKey(), session)
	if err != nil {
		return nil, &validate.Err{Code: code.LOGIN_FAIL}
	}

	return session, nil
}

func (ms UserService) GetUserInfo(lineId string, agencyId string, account string) (*structs.UserResp, error) {
	//account = GetLineAccount(lineId, agencyId, account)
	sess := conf.GetXormSession()
	defer sess.Close()

	has, result := UserBo.GetOneByAccount(sess, lineId, agencyId, account)
	if !has {
		return nil, &validate.Err{Code: code.MEMBER_INFORMATION_QUERY_FAILED}
	}
	res := &structs.UserResp{
		Id:               result.Id,                                         // 主键id
		Account:          result.Account,                                    // 账号
		Balance:          result.Balance,                                    // 会员余额
		CreateTime:       result.CreateTime,                                 // 创建时间
		EditTime:         result.EditTime,                                   // 修改时间
		Capital:          result.Capital,                                    // 红包押金
		AvailableBalance: common.DecimalSub(result.Balance, result.Capital), // 可用金额
		LastLoginIp:      result.LastLoginIp,                                // 上次登陆ip
		LastLoginTime:    result.LastLoginTime,                              // 上次登陆时间
	}
	return res, nil
}
