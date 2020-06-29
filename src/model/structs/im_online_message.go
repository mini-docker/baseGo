package structs

// 实时消息表
type ImOnlineMessage struct {
	Id            int    `xorm:"'id' PK autoincr" json:"id"`             // 主键id
	Message       string `xorm:"message" json:"message"`                 // 消息内容
	SenderId      int    `xorm:"sender_id" json:"sender_id"`             // 消息发送人id
	Sender        string `xorm:"sender" json:"sender"`                   // 消息发送人昵称
	SendTime      int    `xorm:"send_time" json:"send_time"`             // 发送时间
	ReceiverId    int    `xorm:"receiver_id" json:"receiver_id"`         // 接收人id
	ReceiveRoomId int    `xorm:"receive_room_id" json:"receive_room_id"` // 接收房间id
	MessageType   int    `xorm:"message_type" json:"message_type"`       // 消息类型（1.文本；2.图片；3.视频；4.语音;5.红包；6.投注计划；7.通知；）
	NoticeType    int    `xorm:"notice_type" json:"notice_type"`         // 通知消息类型
	DeleteTime    int    `xorm:"delete_time" json:"delete_time"`         // 删除时间
}

// 历史消息返回
type MsgResp struct {
	RoomId      int    `json:"room_id"`      // 房间id （群聊必填）
	Mid         int    `json:"mid"`          // 发送人ID  （密聊必填）
	SenderName  string `json:"sender_name"`  // 发送人昵称
	SenderHead  string `json:"sender_head"`  // 发送人头像
	SendTime    int    `json:"send_time"`    // 发送时间 （非必填，取服务器时间）
	MsgType     int    `json:"msg_type"`     // 消息类型（1.文本；2.图片；3.视频；4.语音;）（必填）
	Msg         string `json:"msg"`          // 消息内容(json)  （必填）
	MsgId       int    `json:"msg_id"`       // 消息id
	ReceiverId  int    `json:"receiver_id"`  // 接收人
	ReceiveType int    `json:"receive_type"` // 消息类型（2群聊，1私聊）
}

func (*ImOnlineMessage) TableName() string {
	return TABLE_ONLINE_MESSAGE
}

// im连接返回对象
type ConnectReply struct {
	Data struct {
		Mid       int64   `json:"mid"`
		Key       string  `json:"key"`
		RoomID    string  `json:"room_id"`
		Accepts   []int32 `json:"accepts"`
		Heartbeat int64   `json:"heartbeat"`
		Server    string  `json:"server"`
	} `json:"data"`
}

// 会员连接信息缓存
type ChannelObj struct {
	Online  bool   `json:"online"`
	Account string `json:"account"`
	UserId  int    `json:"user_id"`
}

// 通知消息结构体
type Notice struct {
	NoticeType int    `json:"notice_type"` // 操作类型
	MuteTime   int    `json:"mute_int"`    //禁言时间
	Message    string `json:"message"`     // 群内通知备注
}

// 批量拉群消息结构体
type IntiveOrKickRoomReq struct {
	UserKeys   []string `json:"user_keys"`
	RoomId     int      `json:"roomId"`
	NoticeType int      `json:"notice_type"`
	SenderId   int      `json:"sender_id"`
	IsInvite   bool     `json:"is_invite"`
}

type DelRoomReq struct {
	RoomId int `json:"roomId"`
}
