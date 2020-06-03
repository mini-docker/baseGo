package services

import (
	"model/bo"
	"model/code"
	"model/structs"
	"red_admin/app/middleware/validate"
	"red_admin/conf"
)

type SystemLineRoyaltyService struct {
}

var (
	SystemLineRoyaltyBo = new(bo.SystemLineRoyalty)
)

func (SystemLineRoyaltyService) QueryLineRoyaltyList(startTime, endTime int) ([]*structs.LineRoyaltyListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取线路提成信息
	royaltyList, err := SystemLineRoyaltyBo.QueryLineRoyaltyList(sess, startTime, endTime)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	for _, v := range royaltyList {
		v.Childrens = make([]structs.AgencyRoyaltyListResp, 0)
	}
	return royaltyList, nil
}

func (SystemLineRoyaltyService) QueryLineAgencyRoyaltyList(startTime, endTime int, lineId string) ([]*structs.AgencyRoyaltyListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取线路代理提成信息
	royaltyList, err := SystemLineRoyaltyBo.QueryLineAgencyRoyaltyList(sess, startTime, endTime, lineId)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}

	// 获取站点信息
	sites, err := new(bo.RedPacketSite).SiteCode(sess, lineId)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	var AgencyIds []string
	for _, v := range sites {
		AgencyIds = append(AgencyIds, v.AgencyId)
	}
	for _, r := range royaltyList {
		for _, a := range sites {
			if r.AgencyId == a.AgencyId {
				r.AgencyAccount = a.Account
			}
		}
	}
	return royaltyList, nil
}
