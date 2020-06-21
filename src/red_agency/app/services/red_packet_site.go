package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
	"fmt"
)

var (
	RedPacketSiteBo = new(bo.RedPacketSite)
)

type RedPacketSiteService struct {
}

// 站点列表
func (RedPacketSiteService) QueryPacketSiteList(lineId, siteName string, status, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	count, siteList, err := RedPacketSiteBo.QueryPacketSiteList(sess, lineId, siteName, status, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = siteList
	pageResp.Count = count
	return pageResp, nil
}

// 添加站点
func (RedPacketSiteService) AddPacketSite(lineId, agencyId, siteName string, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	site := new(structs.RedPacketSite)
	site.LineId = lineId
	site.AgencyId = agencyId
	site.SiteName = siteName
	site.Status = status
	site.CreateTime = utility.GetNowTimestamp()
	site.DeleteTime = 0
	_, err := RedPacketSiteBo.AddSite(sess, site)
	if err != nil {
		fmt.Printf("AddPacketSite error %v \n", err)
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 修改站点
func (RedPacketSiteService) EditPacketSite(id int, siteName string, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取站点信息
	has, site := RedPacketSiteBo.QuerySiteOne(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	site.SiteName = siteName
	site.Status = status
	err := RedPacketSiteBo.EditPacketSite(sess, site)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 删除站点
func (RedPacketSiteService) DelPacketSite(id int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取站点信息
	has, site := RedPacketSiteBo.QuerySiteOne(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	// 判断站点是否存在代理，存在代理不能删除
	agencys, _ := AgencyBo.FindAgencyByAgencyId(sess, site.LineId, site.AgencyId)
	if len(agencys) > 0 {
		return &validate.Err{Code: code.SITE_CAN_NOT_BE_DELETE}
	}
	site.DeleteTime = utility.GetNowTimestamp()
	err := RedPacketSiteBo.DelPacketSite(sess, site)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}
	return nil
}

// 修改站点转台
func (RedPacketSiteService) EditRedPacketSiteStatus(id int, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断站点是否存在
	has, site := RedPacketSiteBo.QuerySiteOne(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	site.Status = status
	err := RedPacketSiteBo.EditPacketSite(sess, site)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 获取代理code
func (RedPacketSiteService) SiteCode(lineId string) ([]*structs.AgencyCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	agencys, err := RedPacketSiteBo.SiteCode(sess, lineId)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return agencys, err
}
