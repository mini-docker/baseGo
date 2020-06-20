package structs

type RedLog struct {
	Id         int64  `xorm:"'id' PK autoincr" json:"id"`    // 主键ID
	LineId     string `xorm:"line_id" json:"lineId"`         // 线路id
	AgencyId   string `xorm:"agency_id" json:"agencyId"`     // 站点id
	LogType    int    `xorm:"log_type" json:"logType"`       // 日志类型（1登录日志，2操作日志）
	Remark     string `xorm:"remark" json:"remark"`          // 日志
	Creator    string `xorm:"creator" json:"creator"`        // 操作人
	CreatorId  int    `xorm:"creator_id" json:"creatorId"`   // 操作人id
	CreatorIp  string `xorm:"creator_ip" json:"creatorIp"`   // 操作人ip
	CreateTime int    `xorm:"create_time" json:"createTime"` // 创建时间
}

func (*RedLog) TableName() string {
	return TABLE_LOG
}
