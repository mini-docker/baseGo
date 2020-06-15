package structs

type SystemLine struct {
	Id         int     `json:"id" xorm:"'id' PK autoincr"`    // 主键
	LineId     string  `xorm:"line_id"  json:"lineId"`        // 线路id
	LineName   string  `xorm:"line_name" json:"lineName"`     // 线路名称
	LimitCost  float64 `xorm:"limit_cost" json:"limitCost"`   // 线路额度
	MealId     int     `xorm:"meal_id" json:"mealId"`         // 套餐id
	Domain     string  `xorm:"domain" json:"domain"`          // 域名
	Status     int     `xorm:"status" json:"status"`          // 状态 1启用 2停用 3维护
	TransType  int     `xorm:"trans_type" json:"transType"`   // 交易模式  1 钱包  2 额度转换
	ApiUrl     string  `xorm:"api_url" json:"apiUrl"`         // 钱包api地址
	Md5key     string  `xorm:"md5key" json:"md5key"`          // md5key
	RsaPubKey  string  `xorm:"rsa_pub_key" json:"rsaPubKey"`  // rsa共钥
	RsaPriKey  string  `xorm:"rsa_pri_key" json:"rsaPriKey"`  // rsa私钥
	CreateTime int     `xorm:"create_time" json:"createTime"` // 创建时间
	EditTime   int     `xorm:"edit_time" json:"editTime"`     // 修改时间
}

func (SystemLine) TableName() string {
	return TABLE_SYSTEM_LINE
}

type SystemLineCode struct {
	LineId string `xorm:"line_id"  json:"lineId"` // 线路id
}
