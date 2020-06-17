package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
	"fmt"
)

type SystemAgencyBo struct{}

// 返回所有代理列表
func (*SystemAgencyBo) QuerySystemAgencyAdminList(sess *xorm.Session, lineId, account string, isOnline int, page, pageSize int) (int64, []*structs.AgencyReq, error) {
	rows := make([]*structs.AgencyReq, 0)
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if account != "" {
		sess.Where("account like ?", account+"%")
	}
	if isOnline != 0 {
		sess.Where("is_online = ?", isOnline)
	}
	sess.Where("delete_time = ? and is_admin = 1", model.UNDEL)
	count, err := sess.Table(new(structs.Agency).TableName()).Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&rows)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}

// 返回所有代理列表
func (*SystemAgencyBo) QuerySystemAgencyList(sess *xorm.Session, lineId, account string, agencyId string, isOnline, status int, page, pageSize int) (int64, []*structs.AgencyReq, error) {
	rows := make([]*structs.AgencyReq, 0)
	sess.Where("delete_time = ? and is_admin = 2", model.UNDEL)
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if account != "" {
		sess.Where("account like ? ", account+"%")
	}
	if isOnline != 0 {
		sess.Where("is_online = ? ", isOnline)
	}
	if status != 0 {
		sess.Where("status = ? ", status)
	}
	fmt.Println(lineId, agencyId, account, isOnline, status, "testAgencyAccount")
	count, err := sess.Table(new(structs.Agency).TableName()).Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&rows)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}

// 添加代理
func (*SystemAgencyBo) AddAgency(sess *xorm.Session, Agency *structs.Agency) (int64, error) {
	return sess.Insert(Agency)
}

// 修改代理信息
func (*SystemAgencyBo) EditAgency(sess *xorm.Session, Agency *structs.Agency) error {
	_, err := sess.Table(new(structs.Agency).TableName()).
		ID(Agency.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("limit", "white_ip_address", "password", "status", "edit_time").
		Update(Agency)
	return err
}

// 修改代理状态
func (*SystemAgencyBo) EditAgencyStatus(sess *xorm.Session, Agency *structs.Agency) error {
	_, err := sess.Table(new(structs.Agency).TableName()).
		ID(Agency.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("status", "edit_time").
		Update(Agency)
	return err
}

// 根据id查询单个代理
func (*SystemAgencyBo) QueryAgencyById(sess *xorm.Session, id int) (*structs.Agency, bool, error) {
	Agency := new(structs.Agency)
	has, err := sess.Where("id = ? and delete_time = ?", id, model.UNDEL).Get(Agency)
	return Agency, has, err
}

// 删除代理
func (*SystemAgencyBo) DelAgency(sess *xorm.Session, Agency *structs.Agency) error {
	_, err := sess.Table(new(structs.Agency).TableName()).
		ID(Agency.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("delete_time").
		Update(Agency)
	return err
}

// 根据账号查询代理信息
func (*SystemAgencyBo) QueryAgencyByAccount(sess *xorm.Session, account string) (*structs.Agency, bool, error) {
	agency := new(structs.Agency)
	has, err := sess.Where("account = ?", account).Get(agency)
	return agency, has, err
}

// 修改代理登陆状态
func (*SystemAgencyBo) EditAgencyOnlineStatus(sess *xorm.Session, agency *structs.Agency) error {
	_, err := sess.Table(new(structs.Agency).TableName()).
		ID(agency.Id).
		Where("delete_time = ?", model.UNDEL).
		Cols("is_online").
		Update(agency)
	return err
}

func (*SystemAgencyBo) QueryAgencysByIds(sess *xorm.Session, ids []int) ([]*structs.Agency, error) {
	agencys := make([]*structs.Agency, 0)
	sess.In("id", ids)
	sess.Where("delete_time = ?", model.UNDEL)
	err := sess.Find(&agencys)
	return agencys, err
}
