package structs

// 群信息
type Room struct {
	Id                  int64   `xorm:"'id' PK autoincr" json:"id"`                        // 主键ID
	LineId              string  `xorm:"line_id" json:"lineId"`                             // 线路id
	AgencyId            string  `xorm:"agency_id" json:"agencyId"`                         // 线路id
	RoomName            string  `xorm:"room_name" json:"roomName"`                         // 群名称
	GameType            int     `xorm:"game_type" json:"gameType"`                         // 游戏 2扫雷 1牛牛
	MaxMoney            float64 `xorm:"max_money" json:"maxMoney"`                         //  最大红包金额
	MinMoney            float64 `xorm:"min_money" json:"minMoney"`                         // 最小红包金额
	GamePlay            int     `xorm:"game_play" json:"gamePlay"`                         // 游戏玩法 牛牛 1经典牛牛 2平倍牛牛 3超倍牛牛 扫雷 1固定赔率 2不固定赔率
	Odds                float64 `xorm:"odds" json:"odds"`                                  // 赔率
	RedNum              int     `xorm:"red_num" json:"redNum"`                             // 红包个数
	RedMinNum           int     `xorm:"red_min_num" json:"redMinNum"`                      // 红包最小个数
	Royalty             float64 `xorm:"royalty" json:"royalty"`                            // 抽水比例
	RoyaltyMoney        float64 `xorm:"royalty_money" json:"royaltyMoney"`                 // 抽水金额
	GameTime            int     `xorm:"game_time" json:"gameTime"`                         // 游戏时间
	RoomSort            int     `xorm:"room_sort" json:"roomSort"`                         // 排序
	Status              int     `xorm:"status" json:"status"`                              // 状态 1启用 2禁用
	CreateTime          int     `xorm:"create_time" json:"createTime"`                     // 创建时间
	DeleteTime          int     `xorm:"delete_time" json:"deleteTime"`                     // 删除时间
	RoomType            int     `xorm:"room_type" json:"roomType"`                         // 群类型  1公群  2私群
	FreeFromDeath       int     `xorm:"free_from_death" json:"freeFromDeath"`              // 是否开启免死号  1开启  2关闭
	RobotSendPacket     int     `xorm:"robot_send_packet" json:"robotSendPacket"`          // 机器人发包状态 1开启 2关闭
	RobotSendPacketTime int     `xorm:"robot_send_packet_time" json:"robotSendPacketTime"` // 机器人发包时间
	RobotGrabPacket     int     `xorm:"robot_grab_packet" json:"robotGrabPacket"`          // 机器人抢包开关 1开启 2关闭
	RoomNo              int     `xorm:"room_no" json:"roomNo"`                             // 房间号
	RobotId             int     `xorm:"robot_id" json:"robotId"`                           // 群主id
	ControlKill         int     `xorm:"control_kill" json:"control_kill"`                  // 是否开启控杀  1 开启  2 关闭
	RobotAccount        string  `xorm:"-" json:"robotAccount"`                             // 群主账号
	LastTime            int     `xorm:"-" json:"lastTime"`                                 // 上次发送红包时间
}

type RoomResp struct {
	Id                  int64   `xorm:"'id' PK autoincr" json:"id"`                        // 主键ID
	LineId              string  `xorm:"line_id" json:"lineId"`                             // 线路id
	AgencyId            string  `xorm:"agency_id" json:"agencyId"`                         // 线路id
	RoomName            string  `xorm:"room_name" json:"roomName"`                         // 群名称
	GameType            int     `xorm:"game_type" json:"gameType"`                         // 游戏 2扫雷 1牛牛
	MaxMoney            float64 `xorm:"max_money" json:"maxMoney"`                         //  最大红包金额
	MinMoney            float64 `xorm:"min_money" json:"minMoney"`                         // 最小红包金额
	GamePlay            int     `xorm:"game_play" json:"gamePlay"`                         // 游戏玩法 牛牛 1经典牛牛 2平倍牛牛 3超倍牛牛 扫雷 1固定赔率 2不固定赔率
	Odds                float64 `xorm:"odds" json:"odds"`                                  // 赔率
	RedNum              int     `xorm:"red_num" json:"redNum"`                             // 红包个数
	RedMinNum           int     `xorm:"red_min_num" json:"redMinNum"`                      // 红包最小个数
	Royalty             float64 `xorm:"royalty" json:"royalty"`                            // 抽水比例
	RoyaltyMoney        float64 `xorm:"royalty_money" json:"royaltyMoney"`                 // 抽水金额
	GameTime            int     `xorm:"game_time" json:"gameTime"`                         // 游戏时间
	RoomSort            int     `xorm:"room_sort" json:"roomSort"`                         // 排序
	Status              int     `xorm:"status" json:"status"`                              // 状态 1启用 2禁用
	CreateTime          int     `xorm:"create_time" json:"createTime"`                     // 创建时间
	RoomType            int     `xorm:"room_type" json:"roomType"`                         // 群类型  1公群  2私群
	FreeFromDeath       int     `xorm:"free_from_death" json:"freeFromDeath"`              // 是否开启免死号  1开启  2关闭
	RobotSendPacket     int     `xorm:"robot_send_packet" json:"robotSendPacket"`          // 机器人发包状态 1开启 2关闭
	RobotSendPacketTime int     `xorm:"robot_send_packet_time" json:"robotSendPacketTime"` // 机器人发包时间
	RobotGrabPacket     int     `xorm:"robot_grab_packet" json:"robotGrabPacket"`          // 机器人抢包开关 1开启 2关闭
	RoomNo              int     `xorm:"room_no" json:"roomNo"`                             // 房间号
	RobotId             int     `xorm:"robot_id" json:"robotId"`                           // 群主id
	ControlKill         int     `xorm:"control_kill" json:"controlKill"`                   // 是否开启控杀  1 开启  2 关闭
}

func (*Room) TableName() string {
	return TABLE_ROOM
}

type RoomInfoResp struct {
	Id                  int64   `json:"id"`                  // 主键ID
	RoomName            string  `json:"roomName"`            // 群名称
	GameType            int     `json:"gameType"`            // 游戏 2扫雷 1牛牛
	MaxMoney            float64 `json:"maxMoney"`            // 最大红包金额
	MinMoney            float64 `json:"minMoney"`            // 最小红包金额
	GamePlay            int     `json:"gamePlay"`            // 游戏玩法 牛牛 1经典牛牛 2平倍牛牛 3超倍牛牛 扫雷 1固定赔率 2不固定赔率
	Odds                float64 `json:"odds"`                // 赔率
	RedNum              int     `json:"redNum"`              // 红包个数
	Royalty             float64 `json:"royalty"`             // 抽水比例
	RoyaltyMoney        float64 `json:"royaltyMoney"`        // 抽水金额
	GameTime            int     `json:"gameTime"`            // 游戏时间
	RoomSort            int     `json:"roomSort"`            // 排序
	Status              int     `json:"status"`              // 状态 1启用 2禁用
	CreateTime          int     `json:"createTime"`          // 状态 1启用 2禁用
	RoomType            int     `json:"roomType"`            // 群类型  1公群  2私群
	FreeFromDeath       int     `json:"freeFromDeath"`       // 是否开启免死号  1开启  2关闭
	RobotSendPacket     int     `json:"robotSendPacket"`     // 机器人发包状态 1开启 2关闭
	RobotSendPacketTime int     `json:"robotSendPacketTime"` // 机器人发包时间
	RobotGrabPacket     int     `json:"robotGrabPacket"`     // 机器人抢包开关 1开启 2关闭
	RoomNo              int     `json:"roomNo"`              // 房间号
	RobotId             int     `json:"robotId"`             // 群主id
}

// 群枚举
type RoomCode struct {
	Id       int    `xorm:"'id' PK autoincr" json:"id"`
	RoomName string `xorm:"room_name" json:"roomName"`
}
