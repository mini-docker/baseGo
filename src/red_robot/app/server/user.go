package server

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/red_robot/app/middleware/validate"
	"baseGo/src/red_robot/conf"
)

type User interface {
	Info() interface{}
}

type UserServer struct{}

var (
	UserBo = new(bo.User)
)

// 会员金额变更
func (UserServer) ChangeAmount(sess *xorm.Session, lineId, agencyId, account, apiUrl string, transType, userId int, redAmount float64, remark string, capital float64) error {
	// 直接扣会员余额
	if transType == model.TRANS_TYPE_CONVERSION {
		err := UserBo.UpdateUserBalanceIncr(sess, lineId, agencyId, userId, redAmount, capital)
		if err != nil {
			sess.Rollback()
			golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
			return &validate.Err{Code: code.MEMBER_BALANCE_UPDATE_FAILED}
		}
	} else {
		respData, err := Wallet(apiUrl, "TRANSFER", "account", conf.GetConfig().Listening.Md5key, conf.GetConfig().Listening.Deskey, &conf.ReqMember{
			Username: account,
			Currency: "CNY",
			Amount:   redAmount,
		}, remark, 0)
		if err != nil {
			sess.Rollback()
			golog.Error("UserService", "ChangeAmount", "error:", err)
			return &validate.Err{Code: code.MEMBER_BALANCE_UPDATE_FAILED}
		}
		if respData.Code != 1000 {
			sess.Rollback()
			golog.Error("UserService", "ChangeAmount", "error:", nil, respData.Code, respData, account, redAmount, respData.Msg)
			return &validate.Err{Code: code.MEMBER_BALANCE_UPDATE_FAILED}
		}
		err = UserBo.UpdateUserBalanceIncr(sess, lineId, agencyId, userId, 0, capital)
		if err != nil {
			sess.Rollback()
			golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
			return &validate.Err{Code: code.MEMBER_BALANCE_UPDATE_FAILED}
		}
	}
	return nil
}
