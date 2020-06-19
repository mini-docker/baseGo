package structs

// 公告表
type Post struct {
	Id         int    `xorm:"'id' PK autoincr" json:"id"`    // 主键ID
	LineId     string `xorm:"line_id" json:"lineId"`         // 线路id
	AgencyId   string `xorm:"agency_id" json:"agencyId"`     // 超管id
	Title      string `xorm:"title" json:"title"`            // 公告标题
	StartTime  int    `xorm:"start_time" json:"startTime"`   // 开始时间
	EndTime    int    `xorm:"end_time" json:"endTime"`       // 结束时间
	Sort       int    `xorm:"sort" json:"sort"`              // 排序
	Status     int    `xorm:"status" json:"status"`          // 状态 1启用 2停用
	DeleteTime int    `xorm:"delete_time" json:"deleteTime"` // 删除时间
}

func (*Post) TableName() string {
	return TABLE_POST
}

// 公告详情表
type PostContent struct {
	Id      int    `xorm:"'id' PK autoincr" json:"id"` // 主键ID
	Pid     int    `xorm:"pid" json:"pid"`             // 公告id
	Content string `xorm:"content" json:"content"`     // 公告内容
}

func (*PostContent) TableName() string {
	return TABLE_POST_CONTENT
}

type PostResp struct {
	Id        int    `json:"id"`        // 主键ID
	Title     string `json:"title"`     // 公告标题
	StartTime int    `json:"startTime"` // 开始时间
	EndTime   int    `json:"endTime"`   // 结束时间
	Sort      int    `json:"sort"`      // 排序
	Content   string `json:"content"`   // 公告内容
	Status    int    `json:"status"`    // 状态
}
