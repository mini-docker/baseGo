package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

var (
	PostService = new(services.PostService)
)

type SystemPostController struct {
	QueryPostReq struct {
		LineId    string `json:"lineId"`    // 线路id
		Title     string `json:"title"`     // 公告标题
		Status    int    `json:"status"`    // 状态
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 站点id
	}
	QueryPostOneReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"`
	}
	EditPostStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)"`
		Status int `json:"status" valid:"Must;ErrorCode(3028)"` // 状态
	}
}

// 查询公告信息
func (m SystemPostController) GetAgencyPostList(ctx server.Context) error {
	req := &m.QueryPostReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	result, err := PostService.GetAgencyPostList(req.LineId, req.AgencyId, req.Title, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)

}

// 根据id查询公告信息
func (m SystemPostController) QueryPostById(ctx server.Context) error {
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

// 修改公告状态
func (m SystemPostController) EditPostStatus(ctx server.Context) error {
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
