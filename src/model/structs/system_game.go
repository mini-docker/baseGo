package structs

type Game struct {
	Id       int    `xorm:"'id' PK autoincr" json:"id"` // 主键
	GameName string `xorm:"game_name" json:"gameName"`  // 游戏名称
	GameType int    `xorm:"game_type" json:"gameType"`  // 游戏类型  1 红包
	Status   int    `xorm:"status" json:"status"`       // 游戏状态  1 启用  2 停用
}

func (Game) TableName() string {
	return TABLE_SYSTEM_GAME
}
