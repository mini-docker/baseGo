package model

type CacheSessionReq struct {
	UserId int     `json:"userId"`
	Server string  `json:"server"`
	Cookie string  `json:"cookie"`
	Token  []byte  `json:"token"`
	Rooms  []int64 `json:"rooms"`
}

type ConnectReply struct {
	Mid       int64   `json:"mid"`
	Key       string  `json:"key"`
	RoomID    string  `json:"roomId"`
	Accepts   []int32 `json:"accepts"`
	Heartbeat int64   `json:"heartbeat"`
	Server    string  `json:"server"`
}
type LogicSession struct {
	Server string  `json:"server"`
	RoomId int     `json:"roomId"`
	Rooms  []int64 `json:"rooms"`
	Online bool    `json:"online"`
}
type ChangeRoomReq struct {
	LineId      string `json:"lineId"`
	AgencyId    string `json:"agencyId"`
	UserId      []int  `json:"userId"`
	RoomId      int    `json:"roomId"`
	NoticeType  int    `json:"noticeType"`
	SenderId    int    `json:"senderId"`
	Content     string `json:"content"`
	Title       string `json:"title"`
	PublishTime int    `json:"publishTime"`
	RoomType    int    `json:"roomType"`
}

type OutRoomReq struct {
	LineId     string `json:"lineId"`
	AgencyId   string `json:"agencyId"`
	UserId     int    `json:"userId"`
	RoomId     int    `json:"roomId"`
	NoticeType int    `json:"notice_type"`
	RoomType   int    `json:"roomType"`
}
type DelRoomReq struct {
	RoomId int `json:"roomId"`
}

type ChannelObj struct {
	Online bool  `json:"online"`
	UserId int64 `json:"userId"`
}

type IntiveOrKickRoomReq struct {
	UserKeys   []string `json:"user_keys"`
	RoomId     int      `json:"roomId"`
	NoticeType int      `json:"notice_type"`
	SenderId   int      `json:"sender_id"`
	IsInvite   bool     `json:"is_invite"`
	LineId     string   `json:"lineId"`
	AgencyId   string   `json:"AgencyId"`
}
