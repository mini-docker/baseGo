package services

import (
	"fecho/golog"
	"model/bo"
	"model/code"
	"model/structs"
	"red_admin/app/middleware/validate"
	"red_admin/conf"
)

type PostService struct{}

var (
	PostBo = new(bo.Post)
)

// 公告列表查询
func (*PostService) GetAgencyPostList(lienId string, agencyId string, title string, status, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	count, postList, err := PostBo.GetAgencyPostList(sess, lienId, agencyId, title, status, page, pageSize)
	if err != nil {
		golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = postList
	pageResp.Count = count
	return pageResp, nil
}

// 修改公告状态
func (*PostService) EditPostStatus(id, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询公告
	post, has, _ := PostBo.QueryPostById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}

	post.Status = status

	// 修改状态
	err := PostBo.EditPostStatus(sess, post)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 查询单个公告信息
func (*PostService) QueryPostOne(id int) (*structs.PostResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询公告
	post, has, _ := PostBo.QueryPostById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}

	// 查寻公告内容
	postContent, has, _ := PostBo.QueryPostContentByPid(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DELETE_FAILED}
	}

	postResp := new(structs.PostResp)
	postResp.Id = post.Id
	postResp.Title = post.Title
	postResp.StartTime = post.StartTime
	postResp.EndTime = post.EndTime
	postResp.Status = post.Status
	postResp.Content = postContent.Content

	return postResp, nil
}
