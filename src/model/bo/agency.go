package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type Agency struct{}

// 添加超管额度
func (*Agency) UpdatelimitDecr(sess *xorm.Session, lineId string, agencyId string, money float64) (int, error) {
	sess.Where("id = ? and delete_time = 0", agencyId)
	sess.Where("line_id = ?", lineId)
	sess.Incr("limit", money)
	count, err := sess.Cols("capital").Update(&structs.Agency{})
	return int(count), err
}
func (*Agency) FindAgencyByAgencyId(sess *xorm.Session, lineId string, agencyId string) ([]*structs.Agency, error) {
	sess.Where("delete_time  = 0")
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ? ", agencyId)
	agency := make([]*structs.Agency, 0)
	err := sess.Table(new(structs.Agency).TableName()).Find(&agency)
	return agency, err
}

func (*Agency) FindSiteByAgencyId(sess *xorm.Session, lineId string, agencyId string) ([]*structs.RedPacketSite, error) {
	sess.Where("delete_time  = 0")
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ? ", agencyId)
	sess.Where("status = ?", 1)
	agency := make([]*structs.RedPacketSite, 0)
	err := sess.Table(new(structs.RedPacketSite).TableName()).Find(&agency)
	return agency, err
}
