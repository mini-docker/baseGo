package structs

// 站点管理
type RedPacketSite struct {
	Id         int64  `xorm:"'id' PK autoincr" json:"id"`    // ID
	LineId     string `xorm:"line_id" json:"lineId"`         // 线路id
	AgencyId   string `xorm:"agency_id" json:"agencyId"`     // 站点id
	SiteName   string `xorm:"site_name" json:"siteName"`     // 站点名称
	Status     int    `xorm:"status" json:"status"`          // 状态  1正常 2停用
	CreateTime int    `xorm:"create_time" json:"createTime"` // 添加时间
	DeleteTime int    `xorm:"delete_time" json:"deleteTime"` // 删除时间
}

func (RedPacketSite) TableName() string {
	return TABLE_PACKET_SITE
}

type AgencyCode struct {
	Account  string `xorm:"site_name" json:"account"`  // 账号
	AgencyId string `xorm:"agency_id" json:"agencyId"` // 超管id
}
