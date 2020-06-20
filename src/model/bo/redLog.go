package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type RedLog struct{}

// 查询操作日志
func (*RedLog) GetRedLogList(sess *xorm.Session, lineId, agencyId string, logType, startTime, endTime, page, pageSize int) (int64, []*structs.RedLog, error) {
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if logType != 0 {
		sess.Where("log_type = ? ", logType)
	}
	if startTime != 0 {
		sess.Where("create_time >= ? ", startTime)
	}
	if endTime != 0 {
		sess.Where("create_time <= ? ", endTime)
	}
	data := make([]*structs.RedLog, 0)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("create_time desc").FindAndCount(&data)
	return count, data, err
}

// 添加日志
func (*RedLog) AddLog(sess *xorm.Session, log *structs.RedLog) (int64, error) {
	return sess.Insert(log)
}
