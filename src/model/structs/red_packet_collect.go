package structs

// 红包采集表
type RedPacketCollect struct {
	Id             int64  `xorm:"'id' PK autoincr" json:"id"`            // ID
	LineId         string `xorm:"line_id" json:"lineId"`                 // 线路ID
	AgencyId       string `xorm:"agency_id" json:"agencyId"`             // 超管ID
	SettlementInfo string `xorm:"settlement_info" json:"settlementInfo"` // 结算信息
	CollectStatus  int    `xorm:"collect_status" json:"collectStatus"`   // 采集状态 1未采集 2 已采集
	CreateTime     int    `xorm:"create_time" json:"createTime"`         //时间
}

func (*RedPacketCollect) TableName() string {
	return TABLE_RED_PACKET_COLLECT
}
