package structs

// 离线消息条数信息
type OfflineMessageCount struct {
	Id               int64 `xorm:"'id' PK autoincr" json:"id"`                   // 主键id
	UserId           int   `xorm:"user_id" json:"user_id"`                       // 用户id
	SenderId         int   `xorm:"sender_id" json:"sender_id"`                   // 发送人id
	RoomId           int   `xorm:"room_id" json:"room_id"`                       // 群id
	MessageCount     int   `xorm:"message_count" json:"message_count"`           // 离线消息条数
	OfflineMessageId int   `xorm:"offline_message_id" json:"offline_message_id"` // 第一条离线消息id
}

func (*OfflineMessageCount) TableName() string {
	return TABLE_OFFLINE_MESSAFE_COUNT
}
