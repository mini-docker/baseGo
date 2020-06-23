package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
	"strings"
)

var RobotBo = new(bo.RobotBo)

type RobotService struct {
}

// 机器人列表
func (RobotService) QueryRobotList(lineId string, agencyId string, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取注单列表
	count, orders, err := RobotBo.QueryRobotList(sess, lineId, agencyId, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = orders
	pageResp.Count = count
	return pageResp, nil
}

// 批量生成机器人账号
func (RobotService) CreatRobotAccounts(lineId string, agencyId string, num int) []*structs.RobotAccounts {
	sess := conf.GetXormSession()
	defer sess.Close()
	robots := make([]*structs.RobotAccounts, 0)
	for i := 0; i < num; i++ {
		robot := new(structs.RobotAccounts)
		m := randInt(5, 8)
		if m < 5 {
			m = 6
		}
		robot.Account = common.RandSeq(sess, lineId, agencyId, m)
		robots = append(robots, robot)
	}
	return robots
}

// 批量生成机器人
func (RobotService) InsertRobots(lineId string, agencyId string, accounts string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	accountlist := strings.Split(accounts, ",")
	robots := make([]*structs.User, 0)
	for _, account := range accountlist {
		if account != "" {
			robot := new(structs.User)
			robot.AgencyId = agencyId
			robot.LineId = lineId
			robot.Status = 1
			robot.CreateTime = utility.GetNowTimestamp()
			robot.IsGroupOwner = 2
			robot.IsRobot = 1
			robot.IsOnline = 2
			robot.Account = account
			robot.Password = utility.NewPasswordEncrypt(account, common.RandPassword(8))
			robots = append(robots, robot)
		}
	}
	err := UserBo.SaveRobots(sess, robots)
	return err
}

// 批量删除机器人
func (RobotService) DelRobots(ids string) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	idList := strings.Split(ids, ",")
	err := UserBo.DelUserByIds(sess, idList)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}
	return nil
}
