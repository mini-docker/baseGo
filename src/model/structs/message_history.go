package structs

//历史消息
type MessageHistory struct {
	Id         int    `xorm:"'id' PK autoincr" json:"id"`    // ID
	LineId     string `xorm:"line_id" json:"lineId"`         // 线路id
	AgencyId   string `xorm:"agency_id" json:"agencyId"`     // 超管ID
	MsgType    int    `xorm:"msg_type" json:"msgType"`       //消息类型 1 发红包 2 结算
	MsgContent string `xorm:"msg_content" json:"msgContent"` //消息
	SenderId   int    `xorm:"sender_id" json:"senderId"`     //
	SenderName string `xorm:"sender_name" json:"senderName"` //
	Status     int    `xorm:"status" json:"status"`          //
	SendTime   int    `xorm:"send_time" json:"sendTime"`     //
	RoomId     int    `xorm:"room_id" json:"roomId"`         //
}

func (MessageHistory) TableName() string {
	return TABLE_MESSAGE_HISTORY
}

type MessageHistoryResp struct {
	Id           int                    `xorm:"'id' PK autoincr" json:"id"`    // ID
	MsgType      int                    `xorm:"msg_type" json:"msgType"`       //消息类型 1 发红包 2 结算
	MsgContent   string                 `xorm:"msg_content" json:"msg"`        //消息
	SenderId     int                    `xorm:"sender_id" json:"senderId"`     //
	SenderName   string                 `xorm:"sender_name" json:"senderName"` //
	Status       int                    `xorm:"status" json:"status"`          //
	SendTime     int                    `xorm:"send_time" json:"sendTime"`     //
	RoomId       int                    `xorm:"room_id" json:"roomId"`         //
	RedId        int                    `xorm:"-" json:"redId"`                // 红包ID
	RedStatus    int                    `xorm:"-" json:"redStatus"`            // 红包状态 1进行中  2已结束 3 无效 4 已领取
	GameTime     int                    `xorm:"-" json:"gameTime"`             // 游戏时间
	CreateTime   int                    `xorm:"-" json:"createTime"`           // 创建时间
	Odds         float64                `xorm:"-" json:"odds"`                 // 赔率
	Mine         int                    `xorm:"-" json:"mine"`                 // 雷值
	IsAdmin      int                    `xorm:"-" json:"isAdmin"`              // 是否是庄家
	AdminWinNum  int                    `xorm:"-" json:"adminWinNum"`          // 庄赢数量
	MemberWinNum int                    `xorm:"-" json:"memberWinNum"`         // 闲赢数量
	EndTime      int                    `xorm:"-" json:"endTime"`              // 结束时间
	CurrentTime  int                    `xorm:"-" json:"currentTime"`          // 当前时间
	Data         []RedPacketLogInfoResp `xorm:"-" json:"data"`                 // 红包注单信息
	MsgId        int                    `xorm:"-" json:"msgId"`                // 消息id
}
