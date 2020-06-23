package services

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
)

type PostService struct{}

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

// 添加公告
func (*PostService) AddPost(lineId string, agencyId string, title string, startTime, endTime, status int, content string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	post := new(structs.Post)
	post.LineId = lineId
	post.AgencyId = agencyId
	post.Title = title
	post.StartTime = startTime
	post.EndTime = endTime
	post.Status = status
	// 保存主表
	err := PostBo.AddPost(sess, post)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	postContent := new(structs.PostContent)
	postContent.Pid = post.Id
	postContent.Content = content
	// 保存从表
	err = PostBo.AddPostContent(sess, postContent)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 修改公告
func (*PostService) EditPost(id int, title string, startTime, endTime, status int, content string, sort int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询公告
	post, has, _ := PostBo.QueryPostById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	if title != post.Title {
		post.Title = title
	}
	if startTime != post.StartTime {
		post.StartTime = startTime
	}
	if endTime != post.EndTime {
		post.EndTime = endTime
	}
	if status != post.Status {
		post.Status = status
	}
	if sort != post.Sort && sort != 0 {
		post.Sort = sort
	}
	// 修改主表
	err := PostBo.EditPost(sess, post)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	if content != "" {
		postContent := new(structs.PostContent)
		postContent.Pid = post.Id
		postContent.Content = content
		// 修改从表
		err = PostBo.EditPostContent(sess, postContent)
		if err != nil {
			return &validate.Err{Code: code.INSET_ERROR}
		}
	}
	return nil
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

// 删除公告
func (*PostService) DelPost(id int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 根据id查询公告
	_, has, _ := PostBo.QueryPostById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}

	// 删除公告内容
	err := PostBo.DelPostContent(sess, id)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}
	// 删除公告
	err = PostBo.DelPost(sess, id)
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
		return nil, &validate.Err{Code: code.QUERY_FAILED}
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
