package structs

// 会员现金记录
type MemberCashRecord struct {
	Id         int     `xorm:"'id' PK autoincr" json:"id"`    // ID
	LineId     string  `xorm:"line_id" json:"lineId"`         // 线路id
	AgencyId   string  `xorm:"agency_id" json:"agencyId"`     // 超管ID
	OrderNo    string  `xorm:"order_no" json:"orderNo"`       // 注单号
	GameType   int     `xorm:"game_type" json:"gameType"`     // 游戏类型
	GameName   string  `xorm:"game_name" json:"gameName"`     // 游戏名称
	FlowType   int     `xorm:"flow_type" json:"flowType"`     // 流水类型（发包、返还、赢利、亏损、普通红包发包、普通红包领取、普通红包返还）
	Money      float64 `xorm:"money" json:"money"`            // 金额
	Remark     string  `xorm:"remark" json:"remark"`          // 流水详情
	CreateTime int     `xorm:"create_time" json:"createTime"` // 创建时间
	UserId     int     `xorm:"user_id" json:"userId"`         // 会员ID
	Account    string  `xorm:"account" json:"account"`        // 会员帐号
}

func (*MemberCashRecord) TableName() string {
	return TABLE_MEMBER_CASH_RECORD
}
