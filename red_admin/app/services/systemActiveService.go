package services

import (
	"fecho/golog"
	"model/bo"
	"model/code"
	"model/structs"
	"red_admin/app/middleware/validate"
	"red_admin/conf"
)

type ActivePictureService struct{}

var (
	ActivePictureBo = new(bo.ActivePicture)
)

// 活动列表
func (*ActivePictureService) GetAgencyActiveList(lienId string, agencyId string, activeName string, status, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	count, acList, err := ActivePictureBo.GetAgencyActiveList(sess, lienId, agencyId, activeName, status, page, pageSize)
	if err != nil {
		golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	for _, v := range acList {
		v.Picture = conf.GetCDNConfig().Host + v.Picture
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = acList
	pageResp.Count = count
	return pageResp, nil
	return pageResp, nil
}

// 修改活动状态
func (*ActivePictureService) EditActiveStatus(id, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询活动
	Active, has, _ := ActivePictureBo.QueryActiveById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}

	Active.Status = status

	// 修改状态
	err := ActivePictureBo.EditActiveStatus(sess, Active)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 查询单个活动信息
func (*ActivePictureService) QueryActiveOne(id int) (*structs.ActivePicture, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询活动
	Active, has, _ := ActivePictureBo.QueryActiveById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	Active.Picture = conf.GetCDNConfig().Host + Active.Picture
	return Active, nil
}
