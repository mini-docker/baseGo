package help

import "fecho/xorm"

type Times struct {
	StartTime int `json:"startTime"` //开始时间
	EndTime   int `json:"endTime"`   //结束时间
}

// Make根据时间查询 时间格式：YYYY-MM-DD HH:II:SS
func (t *Times) Make(timeParam string, model *xorm.Session) {

	if t.StartTime > 0 && t.EndTime > 0 {
		model.Where(timeParam+">=?", t.StartTime).And(timeParam+"<=?", t.EndTime)
	} else if t.StartTime > 0 && t.EndTime == 0 {
		model.Where(timeParam+">=?", t.StartTime)
	} else if t.StartTime == 0 && t.EndTime > 0 {
		model.Where(timeParam+"<=?", t.EndTime)
	}
}
