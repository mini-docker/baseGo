package bo

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type Post struct{}

// 查询进行中的公告
func (*Post) GetPostList(sess *xorm.Session, lineId string, agencyId string) ([]structs.Post, error) {
	sess.Where("start_time <= ?", utility.GetNowTimestamp())
	sess.Where("end_time >= ?", utility.GetNowTimestamp())
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("status = 1")
	sess.Where("delete_time = 0")
	data := make([]structs.Post, 0)
	err := sess.OrderBy("sort desc").Find(&data)
	return data, err
}

// 查询公告内容
func (*Post) GetPostContentList(sess *xorm.Session, pids []int) ([]structs.PostContent, error) {
	sess.In("pid", pids)
	data := make([]structs.PostContent, 0)
	err := sess.Find(&data)
	return data, err
}

// 查询代理公告
func (*Post) GetAgencyPostList(sess *xorm.Session, lineId string, agencyId string, title string, status, page, pageSize int) (int64, []*structs.Post, error) {
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if title != "" {
		sess.Where("title like ? ", title+"%")
	}
	if status != 0 {
		sess.Where("status = ? ", status)
	}
	sess.Where("delete_time = 0")
	data := make([]*structs.Post, 0)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("sort asc").FindAndCount(&data)
	return count, data, err
}

// 添加公告
func (*Post) AddPost(sess *xorm.Session, post *structs.Post) error {
	_, err := sess.Insert(post)
	return err
}

// 添加公告内容
func (*Post) AddPostContent(sess *xorm.Session, postContent *structs.PostContent) error {
	_, err := sess.Insert(postContent)
	return err
}

// 根据id查询公告
func (*Post) QueryPostById(sess *xorm.Session, id int) (*structs.Post, bool, error) {
	post := new(structs.Post)
	has, err := sess.ID(id).Get(post)
	return post, has, err
}

// 修改公告信息
func (*Post) EditPost(sess *xorm.Session, post *structs.Post) error {
	_, err := sess.ID(post.Id).Cols("title", "start_time", "end_time", "status", "sort").Update(post)
	return err
}

// 根据pid查询公告内容
func (*Post) QueryPostContentByPid(sess *xorm.Session, pid int) (*structs.PostContent, bool, error) {
	postContent := new(structs.PostContent)
	has, err := sess.Where("pid = ? ", pid).Get(postContent)
	return postContent, has, err
}

// 修改公告内容
func (*Post) EditPostContent(sess *xorm.Session, postContent *structs.PostContent) error {
	_, err := sess.Where("pid = ?", postContent.Pid).Cols("content").Update(postContent)
	return err
}

// 修改公告状态信息
func (*Post) EditPostStatus(sess *xorm.Session, post *structs.Post) error {
	_, err := sess.ID(post.Id).Cols("status").Update(post)
	return err
}

// 删除公告信息
func (*Post) DelPost(sess *xorm.Session, id int) error {
	_, err := sess.ID(id).Delete(new(structs.Post))
	return err
}

// 删除公告内容
func (*Post) DelPostContent(sess *xorm.Session, pid int) error {
	_, err := sess.Where("pid = ?", pid).Delete(new(structs.PostContent))
	return err
}
