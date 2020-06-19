package structs

type ActivePicture struct {
	Id         int64  `xorm:"'id' PK autoincr" json:"id"`    // 主键ID
	LineId     string `xorm:"line_id" json:"lineId"`         // 线路id
	AgencyId   string `xorm:"agency_id" json:"agencyId"`     // 超管id
	ActiveName string `xorm:"active_name" json:"activeName"` // 活动名称
	Picture    string `xorm:"picture" json:"picture"`        // 图片地址
	StartTime  int    `xorm:"start_time" json:"startTime"`   // 开始时间
	EndTime    int    `xorm:"end_time" json:"endTime"`       // 结束时间
	Sort       int    `xorm:"sort" json:"sort"`              // 排序
	Status     int    `xorm:"status" json:"status"`          // 状态 1启用 2停用
	DeleteTime int    `xorm:"delete_time" json:"deleteTime"` // 删除时间
}

func (*ActivePicture) TableName() string {
	return TABLE_ACTIVE_PICYURE
}

type ActivePictureResp struct {
	Id         int    `json:"id"`         // 主键ID
	ActiveName string `json:"activeName"` // 活动名称
	Picture    string `json:"picture"`    // 图片地址
	StartTime  int    `json:"startTime"`  // 开始时间
	EndTime    int    `json:"endTime"`    // 结束时间
	Sort       int    `json:"sort"`       // 排序
}
