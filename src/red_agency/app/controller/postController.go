package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
)

var (
	PostService        = new(services.PostService)
	UserSessionService = new(middleware.UserSessionService)
)

type PostController struct {
	QueryPostReq struct {
		Title     string `json:"title"`     // 公告标题
		Status    int    `json:"status"`    // 状态
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 代理id
	}

	AddPostReq struct {
		AgencyId  string `json:"agencyId"`                               // 代理id
		Title     string `json:"title" valid:"Must;ErrorCode(3062)"`     // 公告标题
		StartTime int    `json:"startTime" valid:"Must;ErrorCode(3059)"` // 开始时间
		EndTime   int    `json:"endTime" valid:"Must;ErrorCode(3060)"`   // 结束时间
		Content   string `json:"content" valid:"Must;ErrorCode(3061)"`   // 公告内容
		Status    int    `json:"status" valid:"Must;ErrorCode(3028)"`    // 状态
	}

	QueryPostOneReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"`
	}

	EditPostReq struct {
		Id        int    `json:"id" valid:"Must;ErrorCode(3031)"`
		Title     string `json:"title" valid:"Must;ErrorCode(3062)"`     // 公告标题
		StartTime int    `json:"startTime" valid:"Must;ErrorCode(3059)"` // 开始时间
		EndTime   int    `json:"endTime" valid:"Must;ErrorCode(3060)"`   // 结束时间
		Content   string `json:"content"`                                // 公告内容
		Status    int    `json:"status" valid:"Must;ErrorCode(3028)"`    // 状态
		Sort      int    `json:"sort"`
	}

	EditPostStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)"`
		Status int `json:"status" valid:"Must;ErrorCode(3028)"` // 状态
	}

	DelPostStatusReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"`
	}
}

// 查询公告信息
func (m PostController) GetAgencyPostList(ctx server.Context) error {
	req := &m.QueryPostReq
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
	result, err := PostService.GetAgencyPostList(session.User.LineId, req.AgencyId, req.Title, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)

}

// 添加公告信息
func (m PostController) AddPost(ctx server.Context) error {
	req := &m.AddPostReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	if session.IsAdmin == 1 {
		if req.AgencyId == "" {
			return common.HttpResultJsonError(ctx, &validate.Err{Code: code.AGENCY_ID_CAN_NOT_BE_EMPTY})
		}
	} else {
		req.AgencyId = session.AgencyId()
	}

	err = PostService.AddPost(session.LineId(), req.AgencyId, req.Title, req.StartTime, req.EndTime, req.Status, req.Content)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询公告信息
func (m PostController) QueryPostById(ctx server.Context) error {
	req := &m.QueryPostOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	post, err := PostService.QueryPostOne(req.Id)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, post)
}

// 修改公告信息
func (m PostController) EditPost(ctx server.Context) error {
	req := &m.EditPostReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	if req.Sort == 0 {
		if req.Content == "" {
			return common.HttpResultJsonError(ctx, &validate.Err{Code: code.POST_CONTENT_CAN_NOT_BE_EMPTY})
		}
	}
	err := PostService.EditPost(req.Id, req.Title, req.StartTime, req.EndTime, req.Status, req.Content, req.Sort)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改公告状态
func (m PostController) EditPostStatus(ctx server.Context) error {
	req := &m.EditPostStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := PostService.EditPostStatus(req.Id, req.Status)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 删除公告
func (m PostController) DelPost(ctx server.Context) error {
	req := &m.DelPostStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	err := PostService.DelPost(req.Id)

	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
