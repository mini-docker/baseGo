package structs

// 红包表
type RedPacket struct {
	Id                int64   `xorm:"'id' PK autoincr" json:"id"`                   // ID
	LineId            string  `xorm:"line_id" json:"lineId"`                        // 线路ID
	AgencyId          string  `xorm:"agency_id" json:"agencyId"`                    // 超管ID
	RedEnvelopeAmount float64 `xorm:"red_envelope_amount" json:"redEnvelopeAmount"` // 红包金额
	RedEnvelopesNum   int     `xorm:"red_envelopes_num" json:"redEnvelopesNum"`     // 红包数量
	UserId            int     `xorm:"user_id" json:"userId"`                        // 会员ID
	Account           string  `xorm:"account" json:"account"`                       // 会员帐号
	CreateTime        int     `xorm:"create_time" json:"createTime"`                // 创建时间
	DeleteTime        int     `xorm:"delete_time" json:"deleteTime"`                // 删除时间
	RedType           int     `xorm:"red_type" json:"redType"`                      // 游戏类型 1牛牛 2扫雷
	RedPlay           int     `xorm:"red_play" json:"redPlay"`                      // 红包玩法
	RoomId            int     `xorm:"room_id" json:"roomId"`                        // 群ID
	RoomName          string  `xorm:"room_name" json:"roomName"`                    // 群名称
	Status            int     `xorm:"status" json:"status"`                         // 红包状态 1进行中 2已结束
	Mine              int     `xorm:"mine" json:"mine"`                             // 扫雷 炸弹数字
	Capital           float64 `xorm:"capital" json:"capital"`                       // 红包抵押本金
	Money             float64 `xorm:"money" json:"money"`                           // 输赢金额
	RealMoney         float64 `xorm:"real_money" json:"realMoney"`                  // 实际输赢金额
	RoyaltyMoney      float64 `xorm:"royalty_money" json:"royaltyMoney"`            // 抽水金额
	ReturnMoney       float64 `xorm:"return_money" json:"returnMoney"`              // 返还金额
	IsAuto            int     `xorm:"is_auto" json:"isAuto"`                        // 是否开启自动 1是 0否
	AutoTime          int     `xorm:"auto_time" json:"autoTime"`                    // 自动开始时间
	EndTime           int     `xorm:"end_time" json:"endTime"`                      // 红包结束时间
	IsRobot           int     `xorm:"is_robot" json:"isRobot"`                      // 是否是机器人发包
}

func (*RedPacket) TableName() string {
	return TABLE_RED_PACKET
}

// 红包列表返回
type RedPacketResp struct {
	Id                int64   `xorm:"'id' PK autoincr" json:"id"`                   // ID
	LineId            string  `xorm:"line_id" json:"lineId"`                        // 线路ID
	AgencyId          string  `xorm:"agency_id" json:"agencyId"`                    // 超管ID
	RedEnvelopeAmount float64 `xorm:"red_envelope_amount" json:"redEnvelopeAmount"` // 红包金额
	RedEnvelopesNum   int     `xorm:"red_envelopes_num" json:"redEnvelopesNum"`     // 红包数量
	UserId            int     `xorm:"user_id" json:"userId"`                        // 会员ID
	Account           string  `xorm:"account" json:"account"`                       // 会员帐号
	CreateTime        int     `xorm:"create_time" json:"createTime"`                // 创建时间
	RedType           int     `xorm:"red_type" json:"redType"`                      // 游戏类型 1牛牛 2扫雷
	RedPlay           int     `xorm:"red_play" json:"redPlay"`                      // 红包玩法
	RoomId            int     `xorm:"room_id" json:"roomId"`                        // 群ID
	RoomName          string  `xorm:"room_name" json:"roomName"`                    // 群名称
	Status            int     `xorm:"status" json:"status"`                         // 红包状态 1进行中 2已结束
	Mine              int     `xorm:"mine" json:"mine"`                             // 扫雷 炸弹数字
	Capital           float64 `xorm:"capital" json:"capital"`                       // 红包抵押本金
	Money             float64 `xorm:"money" json:"money"`                           // 输赢金额
	RealMoney         float64 `xorm:"real_money" json:"realMoney"`                  // 实际输赢金额
	RoomSort          int     `xorm:"-" json:"roomSort"`                            // 房间排序
	ReturnMoney       float64 `xorm:"return_money" json:"returnMoney"`              // 返还金额
	IsAdmin           int     `xorm:"-" json:"isAdmin"`                             // 是否是庄家  1是 2否
	EndTime           int     `xorm:"end_time" json:"endTime"`                      // 红包结束时间
}

// 红包注单详情查询查询
type RedPacketOrderRecordResp struct {
	Id                int64                  `xorm:"'id' PK autoincr" json:"id"`                   // ID
	RedEnvelopeAmount float64                `xorm:"red_envelope_amount" json:"redEnvelopeAmount"` // 红包金额
	RedEnvelopesNum   int                    `xorm:"red_envelopes_num" json:"redEnvelopesNum"`     // 红包数量
	UserId            int                    `xorm:"user_id" json:"userId"`                        // 会员ID
	Account           string                 `xorm:"account" json:"account"`                       // 会员帐号
	CreateTime        int                    `xorm:"create_time" json:"createTime"`                // 创建时间
	RedType           int                    `xorm:"red_type" json:"redType"`                      // 游戏类型 1牛牛 2扫雷
	RedPlay           int                    `xorm:"red_play" json:"redPlay"`                      // 红包玩法
	RoomId            int                    `xorm:"room_id" json:"roomId"`                        // 群ID
	RoomName          string                 `xorm:"room_name" json:"roomName"`                    // 群名称
	Status            int                    `xorm:"status" json:"status"`                         // 红包状态 1进行中 2已结束
	Mine              int                    `xorm:"mine" json:"mine"`                             // 扫雷 炸弹数字
	Capital           float64                `xorm:"capital" json:"capital"`                       // 红包抵押本金
	Money             float64                `xorm:"money" json:"money"`                           // 输赢金额
	RealMoney         float64                `xorm:"real_money" json:"realMoney"`                  // 实际输赢金额
	AdminWinNum       int                    `xorm:"-" json:"adminWinNum"`                         // 庄赢数量
	MemberWinNum      int                    `xorm:"-" json:"memberWinNum"`                        // 闲赢数量
	ReceiveNum        int                    `xorm:"-" json:"receiveNum"`                          // 红包领取数量
	ReceiveCountMoney float64                `xorm:"-" json:"receiveCountMoney"`                   // 红包领取总金额
	ReturnMoney       float64                `xorm:"return_money" json:"returnMoney"`              // 返还金额
	IsAdmin           int                    `xorm:"-" json:"isAdmin"`                             // 是否是庄家  1是 2否
	GameTime          int                    `xorm:"-" json:"gameTime"`                            // 游戏时间
	RoomSort          int                    `xorm:"-" json:"roomSort"`                            // 房间排序
	RedNum            int                    `xorm:"-" json:"redNum"`                              // 最大红包个数
	RedMinNum         int                    `xorm:"-" json:"redMinNum"`                           // 红包最小个数
	MaxMoney          float64                `xorm:"-" json:"maxMoney"`                            //  最大红包金额
	MinMoney          float64                `xorm:"-" json:"minMoney"`                            // 最小红包金额
	Odds              float64                `xorm:"-" json:"odds"`                                // 赔率
	Data              []RedPacketLogInfoResp `xorm:"-" json:"data"`                                // 红包注单信息
	AutoTime          int                    `xorm:"auto_time" json:"autoTime"`                    // 自动开始时间
	EndTime           int                    `xorm:"-" json:"endTime"`                             // 结束时间
	CurrentTime       int                    `xorm:"-" json:"currentTime"`                         // 当前时间
	OrderNo           string                 `xorm:"-" json:"orderNo"`                             // 订单号
}

// 红包列表返回
type OrdinaryRedPacketResp struct {
	Id                int64   `xorm:"'id' PK autoincr" json:"id"`                   // ID
	LineId            string  `xorm:"line_id" json:"lineId"`                        // 线路ID
	AgencyId          string  `xorm:"agency_id" json:"agencyId"`                    // 超管ID
	RedEnvelopeAmount float64 `xorm:"red_envelope_amount" json:"redEnvelopeAmount"` // 红包金额
	RedEnvelopesNum   int     `xorm:"red_envelopes_num" json:"redEnvelopesNum"`     // 红包数量
	UserId            int     `xorm:"user_id" json:"userId"`                        // 会员ID
	Account           string  `xorm:"account" json:"account"`                       // 会员帐号
	CreateTime        int     `xorm:"create_time" json:"createTime"`                // 创建时间
	RedType           int     `xorm:"red_type" json:"redType"`                      // 游戏类型 1牛牛 2扫雷 0 普通红包
	RoomId            int     `xorm:"room_id" json:"roomId"`                        // 群ID
	RoomName          string  `xorm:"room_name" json:"roomName"`                    // 群名称
	Status            int     `xorm:"status" json:"status"`                         // 红包状态 1进行中 2已结束
	ReturnMoney       float64 `xorm:"return_money" json:"returnMoney"`              // 返还金额
	ReturnNum         int     `xorm:"-" json:"returnNum"`                           // 剩余数量
	ReceiveMoney      float64 `xorm:"-" json:"receiveMoney"`                        // 领取金额
	ReceiveNum        int     `xorm:"-" json:"receiveNum"`                          // 领取数量
	IsAuto            int     `xorm:"is_auto" json:"isAuto"`                        // 是否开启自动 1是 0否
	AutoTime          int     `xorm:"auto_time" json:"autoTime"`                    // 自动开始时间
	EndTime           int     `xorm:"end_time" json:"endTime"`                      // 红包结束时间
}
