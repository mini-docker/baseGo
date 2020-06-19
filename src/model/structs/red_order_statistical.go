package structs

// 红包表
type RedOrderStatistical struct {
	Id              int64   `xorm:"'id' PK autoincr" json:"id"`              // ID
	LineId          string  `xorm:"line_id" json:"lineId"`                   // 线路ID
	AgencyId        string  `xorm:"agency_id" json:"agencyId"`               // 站点ID
	StatisticalDate string  `xorm:"statistical_date" json:"statisticalDate"` // 统计日期
	ValidBet        float64 `xorm:"valid_bet" json:"validBet"`               // 有效投注
	RedNum          int64   `xorm:"red_num" json:"redNum"`                   // 总局数
	OrderNum        int64   `xorm:"order_num" json:"orderNum"`               // 总注单数
	RoyaltyMoney    float64 `xorm:"royalty_money" json:"royaltyMoney"`       //总抽水金额
	FreeDeathWin    float64 `xorm:"free_death_win" json:"freeDeathWin"`      // 免死号盈利
	RobotWin        float64 `xorm:"robot_win" json:"robotWin"`               // 机器人盈利
	TotalWin        float64 `xorm:"total_win" json:"totalWin"`               // 总盈利
	GameType        int     `xorm:"game_type" json:"gameType"`               // 游戏类型
	ValidMember     int     `xorm:"-" json:"validMember"`					 // 有效会员人数
}

func (*RedOrderStatistical) TableName() string {
	return TABLE_ORDER_STATISTICAL
}

// 统计返回
type OrderStatisticalResp struct {
	TotalData        *TotalData             `json:"totalData"`        // 总计
	OrderSeries      *OrderSeries           `json:"orderSeries"`      // 折线图
	OrderStatistical []*RedOrderStatistical `json:"orderStatistical"` // 站点盈利列表
}

// 总计
type TotalData struct {
	ValidMember  int64   `json:"validMember"`  // 有效人数
	ValidBet     float64 `json:"validBet"`     // 有效投注
	RedNum       int64   `json:"redNum"`       // 总局数
	OrderNum     int64   `json:"orderNum"`     // 总注单数
	RoyaltyMoney float64 `json:"royaltyMoney"` //总抽水金额
	FreeDeathWin float64 `json:"freeDeathWin"` // 免死号盈利
	RobotWin     float64 `json:"robotWin"`     // 机器人盈利
	TotalWin     float64 `json:"totalWin"`     // 总盈利
}

// 趋势图
type OrderSeries struct {
	Data     []string  `json:"data"`     // x轴名称列表
	NameData []string  `json:"nameData"` // 头部显示
	Series   []*Series `json:"series"`   // 对应Data中数据
}

// 折线图Series子项
type Series struct {
	Name  string    `json:"name"`  // 鼠标hover名称
	Type  string    `json:"type"`  // 固定值'line'
	Data  []float64 `json:"data"`  // 对应Data中数据
}

// 有效游戏人数统计
type ValidMemberCount struct {
	Total    int    `xorm:"total" json:"total"`
	GameTime string `xorm:"gameTime" json:"gameTime"`
}

// 有效游戏人数统计
type SiteValidMemberCount struct {
	Total    int    `xorm:"total" json:"total"`
	AgencyId string `xorm:"agency_id" json:"agencyId"`
}