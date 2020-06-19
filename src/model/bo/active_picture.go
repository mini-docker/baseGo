package bo

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

// 活动管理
type ActivePicture struct{}

// 活动列表查询
func (*ActivePicture) GetActiveList(sess *xorm.Session, lineId string, agencyId string) ([]*structs.ActivePicture, error) {
	sess.Where("start_time <= ?", utility.GetNowTimestamp())
	sess.Where("end_time >= ?", utility.GetNowTimestamp())
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("status = 1")
	sess.Where("delete_time = 0")
	data := make([]*structs.ActivePicture, 0)
	err := sess.OrderBy("sort asc").Find(&data)
	return data, err
}

// 查询代理活动
func (*ActivePicture) GetAgencyActiveList(sess *xorm.Session, lineId string, agencyId string, activeName string, status, page, pageSize int) (int64, []*structs.ActivePicture, error) {
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if activeName != "" {
		sess.Where("active_name like ? ", activeName+"%")
	}
	if status != 0 {
		sess.Where("status = ? ", status)
	}
	sess.Where("delete_time = 0")
	data := make([]*structs.ActivePicture, 0)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("sort asc").FindAndCount(&data)
	return count, data, err
}

// 添加活动
func (*ActivePicture) AddActive(sess *xorm.Session, active *structs.ActivePicture) error {
	_, err := sess.Insert(active)
	return err
}

// 根据id查询活动
func (*ActivePicture) QueryActiveById(sess *xorm.Session, id int) (*structs.ActivePicture, bool, error) {
	Active := new(structs.ActivePicture)
	has, err := sess.ID(id).Get(Active)
	return Active, has, err
}

// 修改活动信息
func (*ActivePicture) EditActive(sess *xorm.Session, Active *structs.ActivePicture) error {
	_, err := sess.ID(Active.Id).Cols("active_name", "start_time", "end_time", "status", "sort", "picture").Update(Active)
	return err
}

// 修改活动状态信息
func (*ActivePicture) EditActiveStatus(sess *xorm.Session, Active *structs.ActivePicture) error {
	_, err := sess.ID(Active.Id).Cols("status").Update(Active)
	return err
}

// 删除活动信息
func (*ActivePicture) DelActive(sess *xorm.Session, id int) error {
	_, err := sess.ID(id).Delete(new(structs.ActivePicture))
	return err
}

// 查询以启用的活动数量
func (*ActivePicture) CountActive(sess *xorm.Session, lineId, agencyId string) (int64, error) {
	return sess.Table(new(structs.ActivePicture).TableName()).Where("status = 1 and line_id = ? and agency_id = ? and delete_time = 0", lineId, agencyId).Count()
}
