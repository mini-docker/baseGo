package bo

import (
	"github.com/mini-docker/baseGo/src/fecho/xorm"

	"github.com/mini-docker/baseGo/src/model/structs"
)

type RedPacketSite struct{}

func (RedPacketSite) QueryPacketSiteList(sess *xorm.Session, lineId string, siteName string, status, page, pageSize int) (int64, []*structs.RedPacketSite, error) {
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if siteName != "" {
		sess.Where("site_name like ? ", siteName+"%")
	}
	if status != 0 {
		sess.Where("status = ? ", status)
	}
	sess.Where("delete_time = 0")
	data := make([]*structs.RedPacketSite, 0)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&data)
	return count, data, err
}

func (RedPacketSite) AddSite(sess *xorm.Session, site *structs.RedPacketSite) (int64, error) {
	return sess.Insert(site)
}

func (RedPacketSite) QuerySiteOne(sess *xorm.Session, id int) (bool, *structs.RedPacketSite) {
	site := new(structs.RedPacketSite)
	has, _ := sess.Where("id = ? and delete_time = ?", id, 0).Get(site)
	return has, site
}

func (RedPacketSite) EditPacketSite(sess *xorm.Session, site *structs.RedPacketSite) error {
	_, err := sess.ID(site.Id).Cols("site_name", "status").Update(site)
	return err
}

func (RedPacketSite) DelPacketSite(sess *xorm.Session, site *structs.RedPacketSite) error {
	_, err := sess.ID(site.Id).Cols("delete_time").Update(site)
	return err
}

// 查询站点枚举
func (RedPacketSite) SiteCode(sess *xorm.Session, lineId string) ([]*structs.AgencyCode, error) {
	agencys := make([]*structs.AgencyCode, 0)
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	sess.Where("delete_time = 0")
	err := sess.Table(new(structs.RedPacketSite).TableName()).Find(&agencys)
	return agencys, err
}

// 查询没有机器人的站点
func (RedPacketSite) GetNoRobotsSite(sess *xorm.Session, lineId string) ([]*structs.RedPacketSite, error) {
	site := make([]*structs.RedPacketSite, 0)
	err := sess.SQL("select * from red_packet_site where agency_id not in (select agency_id from red_user where is_robot = 1 and is_group_owner = 2) and line_id = ?", lineId).Find(&site)
	return site, err
}

// 验证站点是否存在
func (RedPacketSite) CheckSite(sess *xorm.Session, lineId string, agencyId string) (bool, *structs.RedPacketSite) {
	site := new(structs.RedPacketSite)
	has, _ := sess.Where("line_id = ? and agency_id = ? and delete_time = ?", lineId, agencyId, 0).Get(site)
	return has, site
}
