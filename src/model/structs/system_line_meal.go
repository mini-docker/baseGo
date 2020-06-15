package structs

type SystemLineMeal struct {
	Id         int     `xorm:"'id' PK autoincr" json:"id" `   // 主键
	MealName   string  `xorm:"meal_name" json:"mealName"`     // 套餐名称
	NnRoyalty  float64 `xorm:"nn_royalty" json:"nnRoyalty"`   // 牛牛红包抽成
	SlRoyalty  float64 `xorm:"sl_royalty" json:"slRoyalty"`   // 扫雷红包抽成
	CreateTime int     `xorm:"create_time" json:"createTime"` // 创建时间
	EditTime   int     `xorm:"edit_time" json:"editTime"`     // 修改时间
}

func (SystemLineMeal) TableName() string {
	return TABLE_SYSTEM_LINE_MEAL
}

type SystemLineMealCode struct {
	Id       int    `xorm:"'id' PK autoincr" json:"id" ` // 主键
	MealName string `xorm:"meal_name" json:"mealName"`   // 套餐名称
}
