package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_api/app/middleware/validate"
	"baseGo/src/red_api/conf"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type FinanceService struct{}

type FinanceResp struct {
	Amount  float64 `json:"amount"`
	OrderNo string  `json:"orderNo"`
}

// 转入金额
func (ms FinanceService) TransferredIn(lineId string, agencyId string, userAccount string, amount float64) (*FinanceResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	//userAccount = GetLineAccount(lineId, fmt.Sprint(agencyId), userAccount)
	// check exist
	has, user := UserBo.GetOneByAccount(sess, lineId, agencyId, userAccount)
	fr := new(FinanceResp)
	if !has {
		return fr, &validate.Err{Code: code.ACCOUNT_DOES_NOT_EXIST}
	}

	if user.Status != 1 {
		return fr, &validate.Err{Code: code.ACCOUNT_DISABLED}
	}

	sess.Begin()
	err := UserBo.UpdateUserBalanceIncr(sess, lineId, agencyId, user.Id, amount, 0)
	if err != nil {
		sess.Rollback()
		return fr, &validate.Err{Code: code.UPDATE_FAILED}
	}
	orderNo := OderNo()
	recordInfo := new(structs.MemberCashRecord)
	recordInfo.LineId = lineId
	recordInfo.AgencyId = agencyId
	recordInfo.GameType = 0
	recordInfo.GameName = ""
	recordInfo.OrderNo = orderNo
	recordInfo.Money = amount
	recordInfo.FlowType = model.TRANSFERRED_IN
	recordInfo.Remark = fmt.Sprintf("转入金额")
	recordInfo.CreateTime = utility.GetNowTimestamp()
	_, err = new(bo.MemberCashRecord).Inster(sess, recordInfo)
	if err != nil {
		sess.Rollback()
		return fr, err
	}
	sess.Commit()
	fr.OrderNo = orderNo
	fr.Amount = user.Balance + amount
	return fr, nil
}

// 转出金额
func (ms FinanceService) TransferredOut(lineId string, agencyId string, userAccount string, amount float64) (*FinanceResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	//userAccount = GetLineAccount(lineId, fmt.Sprint(agencyId), userAccount)
	// check exist
	has, user := UserBo.GetOneByAccount(sess, lineId, agencyId, userAccount)
	fr := new(FinanceResp)
	if !has {
		return fr, &validate.Err{Code: code.ACCOUNT_DOES_NOT_EXIST}
	}

	if user.Status != 1 {
		return fr, &validate.Err{Code: code.ACCOUNT_DISABLED}
	}

	if user.Balance < amount {
		return fr, &validate.Err{Code: code.DOESNT_HAS_ENOUGH_AMOUNT}
	}
	sess.Begin()
	err := UserBo.UpdateUserBalance(sess, lineId, agencyId, user.Id, amount)
	if err != nil {
		sess.Rollback()
		return fr, &validate.Err{Code: code.UPDATE_FAILED}
	}
	orderNo := OderNo()
	recordInfo := new(structs.MemberCashRecord)
	recordInfo.LineId = lineId
	recordInfo.AgencyId = agencyId
	recordInfo.GameType = 0
	recordInfo.GameName = ""
	recordInfo.OrderNo = orderNo
	recordInfo.Money = -amount
	recordInfo.FlowType = model.TRANSFERRED_OUT
	recordInfo.Remark = fmt.Sprintf("转出金额")
	recordInfo.CreateTime = utility.GetNowTimestamp()
	_, err = new(bo.MemberCashRecord).Inster(sess, recordInfo)
	if err != nil {
		sess.Rollback()
		return fr, err
	}
	sess.Commit()
	fr.OrderNo = orderNo
	fr.Amount = user.Balance + amount
	return fr, nil
}

// 订单号生成
func OderNo() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	order := "edzh"
	order += time.Now().Format("20060102150405") + strconv.Itoa(100000+r.Intn(899999))
	return order
}
