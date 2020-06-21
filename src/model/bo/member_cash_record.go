package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

// 会员现金记录查询
type MemberCashRecord struct{}

// 添加现金流水
func (MemberCashRecord) Inster(sess *xorm.Session, data ...*structs.MemberCashRecord) (int, error) {
	count, err := sess.Insert(data)
	return int(count), err
}

// 查询会员现金流水
func (*MemberCashRecord) GetUserCashList(sess *xorm.Session, lineId string, agencyId, userId int, account string) (int, []structs.MemberCashRecord, error) {
	return 0, nil, nil
}
