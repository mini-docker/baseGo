package services

import (
	"model"
	"model/bo"
	"model/code"
	"model/structs"
	"red_admin/app/middleware"
	"red_admin/app/middleware/validate"
	"red_admin/conf"
	"strconv"
	"strings"
)

type SystemUserService struct{}

var (
	SystemUserBo       = new(bo.User)
	UserSessionService = new(middleware.UserSessionService)
)

// 获取用户列表
func (SystemUserService) QueryUserList(lineId string, agencyId string, status, isOnline int, account string, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	count, list, err := SystemUserBo.QueryUserList(sess, lineId, agencyId, status, isOnline, account, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = list
	pageResp.Count = count
	return pageResp, nil
}

// 批量更新会员状态
func (SystemUserService) EditUsersStatus(ids string, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	var idlist []int
	idStrs := strings.Split(ids, ",")
	for _, v := range idStrs {
		id, _ := strconv.Atoi(v)
		if id != 0 {
			idlist = append(idlist, id)
		}
	}

	if len(idlist) > 0 {
		err := SystemUserBo.EditUsersStatus(sess, idlist, status)
		if err != nil {
			return &validate.Err{Code: code.UPDATE_FAILED}
		}
	}
	return nil
}

// 批量踢线
func (SystemUserService) KickUsers(ids string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	var idList []int
	if len(ids) > 0 {
		strs := strings.Split(ids, ",")
		for _, v := range strs {
			if "" != v {
				id, _ := strconv.Atoi(v)
				idList = append(idList, id)
			}
		}
	}
	// 批量修改会员状态
	err := SystemUserBo.KickUsers(sess, idList)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	for _, id := range idList {
		UserSessionService.DelSessionId(model.RED_API_SESSION_LIST_KEY, id)
	}
	return nil
}
