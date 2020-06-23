package services

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
	"strings"
)

type ActivePictureService struct{}

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
}

// 添加活动
func (*ActivePictureService) AddActive(lineId string, agencyId string, activeName string, startTime, endTime, status int, picture string) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 判断启用活动数量
	count, _ := ActivePictureBo.CountActive(sess, lineId, agencyId)
	if count >= 5 {
		if status == 1 {
			return &validate.Err{Code: code.ACTIVE_LIMIT}
		}
	}
	activePicture := new(structs.ActivePicture)
	activePicture.LineId = lineId
	activePicture.AgencyId = agencyId
	activePicture.ActiveName = activeName
	activePicture.StartTime = startTime
	activePicture.EndTime = endTime
	activePicture.Status = status
	activePicture.Picture = picture
	// 保存主表
	err := ActivePictureBo.AddActive(sess, activePicture)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 修改活动
func (*ActivePictureService) EditActive(id int, activeName string, startTime, endTime, status int, picture string, sort int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询活动
	Active, has, _ := ActivePictureBo.QueryActiveById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	if activeName != Active.ActiveName {
		Active.ActiveName = activeName
	}
	if startTime != Active.StartTime {
		Active.StartTime = startTime
	}
	if endTime != Active.EndTime {
		Active.EndTime = endTime
	}
	if status != Active.Status {
		// 判断启用活动数量
		count, _ := ActivePictureBo.CountActive(sess, Active.LineId, Active.AgencyId)
		if count >= 5 {
			if status == 1 {
				return &validate.Err{Code: code.ACTIVE_LIMIT}
			}
		}
		Active.Status = status
	}
	if sort != Active.Sort {
		Active.Sort = sort
	}
	if picture != Active.Picture {
		Active.Picture = strings.Replace(picture, conf.GetCDNConfig().Host, "", -1)
	}
	// 修改活动
	err := ActivePictureBo.EditActive(sess, Active)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}

	return nil
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
	// 判断启用活动数量
	count, _ := ActivePictureBo.CountActive(sess, Active.LineId, Active.AgencyId)
	if count >= 5 {
		if status == 1 {
			return &validate.Err{Code: code.ACTIVE_LIMIT}
		}
	}

	Active.Status = status

	// 修改状态
	err := ActivePictureBo.EditActiveStatus(sess, Active)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 删除活动
func (*ActivePictureService) DelActive(id int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询活动
	_, has, _ := ActivePictureBo.QueryActiveById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	// 删除活动
	err := ActivePictureBo.DelActive(sess, id)
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
