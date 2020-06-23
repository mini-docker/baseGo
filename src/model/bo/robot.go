package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type RobotBo struct{}

func (RobotBo) QueryRobotList(sess *xorm.Session, lineId string, agencyId string, page, pageSize int) (int64, []*structs.User, error) {
	robots := make([]*structs.User, 0)
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	sess.Where("is_robot = ? and delete_time = ?", 1, 0)
	count, err := sess.Table(new(structs.User).TableName()).Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&robots)
	if err != nil {
		return 0, nil, err
	}
	return count, robots, err
}
