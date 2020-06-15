package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type SystemLineBo struct{}

// 返回所有线路列表
func (*SystemLineBo) QuerySystemLineList(sess *xorm.Session, lineId, lineName string, status, transType int, page, pageSize int) (int64, []*structs.SystemLine, error) {
	rows := make([]*structs.SystemLine, 0)
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if lineName != "" {
		sess.Where("line_name like ?", lineName+"%")
	}
	if status != 0 {
		sess.Where("status = ? ", status)
	}
	if transType != 0 {
		sess.Where("trans_type = ?", transType)
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&rows)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}

// 返回所有线路Id枚举
func (*SystemLineBo) QuerySystemLineIdList(sess *xorm.Session) ([]*structs.SystemLineCode, error) {
	rows := make([]*structs.SystemLineCode, 0)
	sess.Where("status != 2")
	err := sess.Table(new(structs.SystemLine).TableName()).Find(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// 添加线路
func (*SystemLineBo) AddLine(sess *xorm.Session, line *structs.SystemLine) (int64, error) {
	return sess.Insert(line)
}

// 修改线路信息
func (*SystemLineBo) EditLine(sess *xorm.Session, line *structs.SystemLine) error {
	_, err := sess.Table(new(structs.SystemLine).TableName()).
		ID(line.Id).
		Cols("line_name", "limit", "meal_id", "domain", "trans_type", "api_url", "md5key", "deskey", "status", "edit_time", "rsa_pub_key", "rsa_pri_key").
		Update(line)
	return err
}

// 修改线路状态
func (*SystemLineBo) EditLineStatus(sess *xorm.Session, line *structs.SystemLine) error {
	_, err := sess.Table(new(structs.SystemLine).TableName()).
		ID(line.Id).
		Cols("status", "edit_time").
		Update(line)
	return err
}

// 根据id查询单个线路
func (*SystemLineBo) QueryLineById(sess *xorm.Session, id int) (*structs.SystemLine, bool, error) {
	line := new(structs.SystemLine)
	has, err := sess.Where("id = ?", id).Get(line)
	return line, has, err
}

// 根据lineId查询单个线路
func (*SystemLineBo) QueryLineBylineId(sess *xorm.Session, lineId string) (*structs.SystemLine, bool, error) {
	line := new(structs.SystemLine)
	has, err := sess.Where("line_id = ?", lineId).Get(line)
	return line, has, err
}

func (*SystemLineBo) QuerySystemLineWallet(sess *xorm.Session) ([]*structs.SystemLine, error) {
	rows := make([]*structs.SystemLine, 0)
	sess.Where("trans_type = 1")
	err := sess.Table(new(structs.SystemLine).TableName()).Find(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
