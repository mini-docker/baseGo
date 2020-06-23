package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
)

var (
	ActivePictureService = new(services.ActivePictureService)
)

type ActivePictureController struct {
	QueryActiveReq struct {
		ActiveName string `json:"activeName"` // 活动标题
		Status     int    `json:"status"`     // 状态
		PageIndex  int    `json:"pageIndex"`  // 页码
		PageSize   int    `json:"pageSize"`   // 每页条数
		AgencyId   string `json:"agencyId"`   // 站点id
	}

	AddActiveReq struct {
		ActiveName string `json:"activeName" valid:"Must;ErrorCode(3062)"` // 活动标题
		StartTime  int    `json:"startTime" valid:"Must;ErrorCode(3026)"`  // 开始时间
		EndTime    int    `json:"endTime" valid:"Must;ErrorCode(3026)"`    // 结束时间
		Status     int    `json:"status" valid:"Must;ErrorCode(3028)"`     // 状态
		Picture    string `json:"picture" valid:"Must;ErrorCode(3063)"`    // 图片地址
		AgencyId   string `json:"agencyId"`
	}

	QueryActiveOneReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"`
	}

	EditActiveReq struct {
		Id         int    `json:"id" valid:"Must;ErrorCode(3031)"`
		ActiveName string `json:"activeName" valid:"Must;ErrorCode(3062)"` // 活动标题
		StartTime  int    `json:"startTime" valid:"Must;ErrorCode(3026)"`  // 开始时间
		EndTime    int    `json:"endTime" valid:"Must;ErrorCode(3026)"`    // 结束时间
		Picture    string `json:"picture" valid:"Must;ErrorCode(3063)"`    // 图片地址
		Status     int    `json:"status" valid:"Must;ErrorCode(3028)"`     // 状态
		Sort       int    `json:"sort"`                                    // 排序
	}

	EditActiveStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)"`
		Status int `json:"status" valid:"Must;ErrorCode(3028)"` // 状态
	}

	DelActiveStatusReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"`
	}
}

// 查询活动列表
func (m ActivePictureController) GetAgencyActiveList(ctx server.Context) error {
	req := &m.QueryActiveReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	if session.IsAdmin != 1 {
		req.AgencyId = session.AgencyId()
	}

	result, err := ActivePictureService.GetAgencyActiveList(session.User.LineId, req.AgencyId, req.ActiveName, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)

}

// 添加活动信息
func (m ActivePictureController) AddActive(ctx server.Context) error {
	req := &m.AddActiveReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	if session.IsAdmin != 1 {
		req.AgencyId = session.AgencyId()
	}

	err = ActivePictureService.AddActive(session.LineId(), req.AgencyId, req.ActiveName, req.StartTime, req.EndTime, req.Status, req.Picture)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询活动信息
func (m ActivePictureController) QueryActiveById(ctx server.Context) error {
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

// 修改活动信息
func (m ActivePictureController) EditActive(ctx server.Context) error {
	req := &m.EditActiveReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := ActivePictureService.EditActive(req.Id, req.ActiveName, req.StartTime, req.EndTime, req.Status, req.Picture, req.Sort)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改活动状态
func (m ActivePictureController) EditActiveStatus(ctx server.Context) error {
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

// 删除活动
func (m ActivePictureController) DelActive(ctx server.Context) error {
	req := &m.DelActiveStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := ActivePictureService.DelActive(req.Id)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
