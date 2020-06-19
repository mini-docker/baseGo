package structs

type LineRoyaltyListResp struct {
	LineId    string                  `xorm:"lineId" json:"lineId"`       // 线路id
	NnRoyalty float64                 `xorm:"nnRoyalty" json:"nnRoyalty"` // 牛牛红包盈利
	SlRoyalty float64                 `xorm:"slRoyalty" json:"slRoyalty"` // 扫雷红包盈利
	Childrens []AgencyRoyaltyListResp `xorm:"-" json:"childrens"`
}

type AgencyRoyaltyListResp struct {
	AgencyId      string  `xorm:"agencyId" json:"agencyId"`   // 代理Id
	NnRoyalty     float64 `xorm:"nnRoyalty" json:"nnRoyalty"` // 牛牛红包盈利
	SlRoyalty     float64 `xorm:"slRoyalty" json:"slRoyalty"` // 扫雷红包盈利
	AgencyAccount string  `xorm:"-" json:"agencyAccount"`     // 代理账号
}

func (LineRoyaltyListResp) TableName() string {
	return TABLE_ORDER_RECORD
}

func (AgencyRoyaltyListResp) TableName() string {
	return TABLE_ORDER_RECORD
}
