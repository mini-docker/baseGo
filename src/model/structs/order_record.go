package structs

type OrderRecord struct {
	Id           int     `xorm:"'id' PK autoincr" json:"id"`         // ID
	LineId       string  `xorm:"line_id" json:"lineId"`              // 线路id
	AgencyId     string  `xorm:"agency_id" json:"agencyId"`          // 超管ID
	UserId       int     `xorm:"user_id" json:"userId"`              // 会员ID
	Account      string  `xorm:"account" json:"account"`             // 账号
	RedSender    string  `xorm:"red_sender" json:"redSender"`        // 发包者
	GameType     int     `xorm:"game_type" json:"gameType"`          // 游戏类型 1牛牛 2扫雷
	GamePlay     int     `xorm:"game_play" json:"gamePlay"`          // 游戏玩法
	RoomId       int     `xorm:"room_id" json:"roomId"`              // 群id
	RoomName     string  `xorm:"room_name" json:"roomName"`          // 群名称
	OrderNo      string  `xorm:"order_no" json:"orderNo"`            // 注单号
	RedId        int     `xorm:"red_id" json:"redId"`                // 红包ID
	RedMoney     float64 `xorm:"red_money" json:"redMoney"`          // 发包金额
	RedNum       int     `xorm:"red_num" json:"redNum"`              // 红包个数
	ReceiveMoney float64 `xorm:"receive_money" json:"receiveMoney"`  // 领取金额
	Royalty      float64 `xorm:"royalty" json:"royalty"`             // 抽水比例
	RoyaltyMoney float64 `xorm:"royalty_money" json:"royaltyMoney"`  // 抽水金额
	Money        float64 `xorm:"money" json:"money"`                 // 输赢金额
	RealMoney    float64 `xorm:"real_money" json:"realMoney"`        // 实际输赢金额
	GameTime     int     `xorm:"game_time" json:"gameTime"`          // 游戏时间
	ReceiveTime  int     `xorm:"receive_time" json:"receiveTime"`    // 红包领取时间
	RedStartTime int     `xorm:"red_start_time" json:"redStartTime"` // 红包开始时间
	Status       int     `xorm:"status" json:"status"`               // 状态 0未结算，1赢，2输，3无效
	Extra        string  `xorm:"extra" json:"extra"`                 // 特性json
	IsRobot      int     `xorm:"is_robot" json:"isRobot"`            // 是否是机器人 1是 2不是
	IsFreeDeath  int     `xorm:"is_free_death" json:"isFreeDeath"`   // 是否免死  1是 2不是
	RobotWin     float64 `xorm:"robot_win" json:"robotWin"`          // 机器人盈利
	ValidBet     float64 `xorm:"valid_bet" json:"validBet"`          // 有效投注
}

func (*OrderRecord) TableName() string {
	return TABLE_ORDER_RECORD
}

// 红包详情返回数据
type RedPacketLogInfoResp struct {
	Id           int     `xorm:"'id' PK autoincr" json:"id"`         // ID
	LineId       string  `xorm:"line_id" json:"lineId"`              // 线路id
	AgencyId     string  `xorm:"agency_id" json:"agencyId"`          // 超管ID
	Account      string  `xorm:"account" json:"account"`             // 账号
	RedSender    string  `xorm:"red_sender" json:"redSender"`        // 发包者
	GameType     int     `xorm:"game_type" json:"gameType"`          // 游戏类型 1牛牛 2扫雷
	GamePlay     int     `xorm:"game_play" json:"gamePlay"`          // 游戏玩法
	RoomId       int     `xorm:"room_id" json:"roomId"`              // 群id
	RoomName     string  `xorm:"room_name" json:"roomName"`          // 群名称
	OrderNo      string  `xorm:"order_no" json:"orderNo"`            // 注单号
	RedId        int     `xorm:"red_id" json:"redId"`                // 红包ID
	RedMoney     float64 `xorm:"red_money" json:"redMoney"`          // 发包金额
	RedNum       int     `xorm:"red_num" json:"redNum"`              // 红包个数
	ReceiveMoney float64 `xorm:"receive_money" json:"receiveMoney"`  // 领取金额
	Royalty      float64 `xorm:"royalty" json:"royalty"`             // 抽水比例
	RoyaltyMoney float64 `xorm:"royalty_money" json:"royaltyMoney"`  // 抽水金额
	Odds         float64 `xorm:"-" json:"odds"`                      // 赔率
	ThunderNum   int     `xorm:"-" json:"thunderNum"`                // 雷值
	AdminNum     int     `xorm:"-" json:"adminNum"`                  // 庄家牛数
	MemberNum    int     `xorm:"-" json:"memberNum"`                 // 玩家牛数
	MemberMine   int     `xorm:"-" json:"memberMine"`                // 闲家结果（尾数）
	Money        float64 `xorm:"money" json:"money"`                 // 输赢金额
	RealMoney    float64 `xorm:"real_money" json:"realMoney"`        // 实际输赢金额
	GameTime     int     `xorm:"game_time" json:"gameTime"`          // 游戏时间
	ReceiveTime  int     `xorm:"receive_time" json:"receiveTime"`    // 红包领取时间
	RedStartTime int     `xorm:"red_start_time" json:"redStartTime"` // 红包开始时间
	Status       int     `xorm:"status" json:"status"`               // 状态 0未结算，1赢，2输，3无效
	IsRobot      int     `xorm:"is_robot" json:"isRobot"`            // 是否是机器人 1是 2不是
	IsFreeDeath  int     `xorm:"is_free_death" json:"isFreeDeath"`   // 是否免死  1是 2不是
	Extra        string  `xorm:"extra" json:"extra"`                 // 特性json
}

type RedOrderResp struct {
	Account      string  `xorm:"account" json:"account"`            // 账号
	GameType     int     `xorm:"game_type" json:"gameType"`         // 游戏类型 1牛牛 2扫雷
	RedSender    string  `xorm:"red_sender" json:"redSender"`       // 发包者
	ReceiveMoney float64 `xorm:"receive_money" json:"receiveMoney"` // 领取金额
	MemberNum    int     `xorm:"-" json:"memberNum"`                // 玩家牛数
	Money        float64 `xorm:"money" json:"money"`                // 输赢金额
	ReceiveTime  int     `xorm:"receive_time" json:"receiveTime"`   // 红包领取时间
	RoyaltyMoney float64 `xorm:"royalty_money" json:"royaltyMoney"` // 抽水金额
	RealMoney    float64 `xorm:"real_money" json:"realMoney"`       // 实际输赢金额
	IsRobot      int     `xorm:"is_robot" json:"isRobot"`           // 是否是机器人 1是 2不是
	IsFreeDeath  int     `xorm:"is_free_death" json:"isFreeDeath"`  // 是否免死  1是 2不是
	InfoType     int     `xorm:"-" json:"infoType"`                 // 类型
	Extra        string  `xorm:"extra" json:"extra"`                // 特性json
}

type OrderRoyalty struct {
	LineId   string  `json:"lineId"`   // 线路id
	GameType int     `json:"gameType"` // 游戏类型
	WinMoney float64 `json:"winMoney"` // 盈利金额
}

// 红包领取数量/金额
type OrderReceiveCountResp struct {
	RedId        int     `xorm:"red_id" json:"redId"`              // ID
	Count        int     `xorm:"count" json:"count"`               // 领取数量
	ReceiveMoney float64 `xorm:"receiveMoney" json:"receiveMoney"` // 领取金额
}

// 注单采集返回
type OrderCollectResp struct {
	Id           int     `xorm:"'id' PK autoincr" json:"id"`         // ID
	LineId       string  `xorm:"line_id" json:"lineId"`              // 线路id
	AgencyId     string  `xorm:"agency_id" json:"agencyId"`          // 超管ID
	UserId       int     `xorm:"user_id" json:"userId"`              // 会员ID
	Account      string  `xorm:"account" json:"account"`             // 账号
	RedSender    string  `xorm:"red_sender" json:"redSender"`        // 发包者
	GameType     int     `xorm:"game_type" json:"gameType"`          // 游戏类型 1牛牛 2扫雷
	GamePlay     int     `xorm:"game_play" json:"gamePlay"`          // 游戏玩法
	RoomId       int     `xorm:"room_id" json:"roomId"`              // 群id
	RoomName     string  `xorm:"room_name" json:"roomName"`          // 群名称
	OrderNo      string  `xorm:"order_no" json:"orderNo"`            // 注单号
	RedId        int     `xorm:"red_id" json:"redId"`                // 红包ID
	RedMoney     float64 `xorm:"red_money" json:"redMoney"`          // 发包金额
	RedNum       int     `xorm:"red_num" json:"redNum"`              // 红包个数
	ReceiveMoney float64 `xorm:"receive_money" json:"receiveMoney"`  // 领取金额
	Royalty      float64 `xorm:"royalty" json:"royalty"`             // 抽水比例
	RoyaltyMoney float64 `xorm:"royalty_money" json:"royaltyMoney"`  // 抽水金额
	Money        float64 `xorm:"money" json:"money"`                 // 输赢金额
	RealMoney    float64 `xorm:"real_money" json:"realMoney"`        // 实际输赢金额
	GameTime     int     `xorm:"game_time" json:"gameTime"`          // 游戏时间
	ReceiveTime  int     `xorm:"receive_time" json:"receiveTime"`    // 红包领取时间
	RedStartTime int     `xorm:"red_start_time" json:"redStartTime"` // 红包开始时间
	Status       int     `xorm:"status" json:"status"`               // 状态 0未结算，1赢，2输，3无效
	Extra        string  `xorm:"extra" json:"extra"`                 // 特性json
	IsRobot      int     `xorm:"is_robot" json:"isRobot"`            // 是否是机器人 1是 2不是
	IsFreeDeath  int     `xorm:"is_free_death" json:"isFreeDeath"`   // 是否免死  1是 2不是
	ValidBet     float64 `xorm:"valid_bet" json:"validBet"`          // 有效投注金额
	Remark       string  `xorm:"-" json:"remark"`                    // 备注
}
