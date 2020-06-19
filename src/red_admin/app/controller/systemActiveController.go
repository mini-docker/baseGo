package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

var (
	ActivePictureService = new(services.ActivePictureService)
)

type SystemActivePictureController struct {
	QueryActiveReq struct {
		LineId     string `json:"lineId"`     // 线路id
		ActiveName string `json:"activeName"` // 活动标题
		Status     int    `json:"status"`     // 状态
		PageIndex  int    `json:"pageIndex"`  // 页码
		PageSize   int    `json:"pageSize"`   // 每页条数
		AgencyId   string `json:"agencyId"`   // 站点id
	}

	QueryActiveOneReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"`
	}

	EditActiveStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)"`
		Status int `json:"status" valid:"Must;ErrorCode(3028)"` // 状态
	}
}

// 查询活动列表
func (m SystemActivePictureController) GetAgencyActiveList(ctx server.Context) error {
	req := &m.QueryActiveReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := ActivePictureService.GetAgencyActiveList(req.LineId, req.AgencyId, req.ActiveName, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)

}

// 根据id查询活动信息
func (m SystemActivePictureController) QueryActiveById(ctx server.Context) error {
	req := &m.QueryActiveOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	Active, err := ActivePictureService.QueryActiveOne(req.Id)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, Active)
}

// 修改活动状态
func (m SystemActivePictureController) EditActiveStatus(ctx server.Context) error {
	req := &m.EditActiveStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := ActivePictureService.EditActiveStatus(req.Id, req.Status)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
