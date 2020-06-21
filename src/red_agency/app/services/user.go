package services

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
	"strconv"
	"strings"
)

type UserService struct{}

var (
	UserSessionService = new(middleware.UserSessionService)
)

// 获取用户列表
func (UserService) QueryUserList(lineId string, agencyId string, status, isOnline int, account string, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	count, list, err := UserBo.QueryUserList(sess, lineId, agencyId, status, isOnline, account, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = list
	pageResp.Count = count
	return pageResp, nil
}

// 批量更新会员状态
func (UserService) EditUsersStatus(ids string, status int) error {
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
		err := UserBo.EditUsersStatus(sess, idlist, status)
		if err != nil {
			return &validate.Err{Code: code.UPDATE_FAILED}
		}
		// 停用踢线
		if status == model.MENU_TWO {
			for _, id := range idlist {
				UserSessionService.DelSessionId(model.RED_API_SESSION_LIST_KEY, id)
			}
		}
	}
	return nil
}

// 批量踢线
func (UserService) KickUsers(ids string) error {
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
	err := UserBo.KickUsers(sess, idList)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	for _, id := range idList {
		UserSessionService.DelSessionId(model.RED_API_SESSION_LIST_KEY, id)
	}
	return nil
}
