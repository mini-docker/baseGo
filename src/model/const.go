package model

// 存放有意义的常量
const (
	// 滑块验证图片缓存key
	UNDEL = 0 //未删除的时间
	// 保留小数点后位数
	ROUND_TWO  = 2 // 两位
	ROUND_FIVE = 5 // 五位
	//pc,wap端
	IS_ANDROID                  = 1         //android端
	IS_IOS                      = 2         //ios端
	IS_WAP                      = 3         //wap
	IS_PC                       = 4         // pc
	ANDROID                     = "ANDROID" //android
	IOS                         = "IOS"     //ios
	WAP                         = "WAP"     // wap
	PC                          = "PC"      //pc
	DEVICE                      = "device"
	SessionKey                  = "sid"
	LineId                      = "lineId"
	AgencyId                    = "agencyId"
	Sign                        = "sign"
	RED_ADMIN_SESSION_LIST_KEY  = "red_admin_session_list_key"
	RED_AGENCY_SESSION_LIST_KEY = "red_agency_session_list_key"
	RED_API_SESSION_LIST_KEY    = "red_api_session_list_key"
	RED_WAP_SESSION_LIST_KEY    = "red_api_session_list_key"

	// 滑块验证图片缓存key
	SLIDER_VERIFICATION_CODE_BIG_KEY   = "slider_verification_code_big_key"
	SLIDER_VERIFICATION_CODE_SMALL_KEY = "slider_verification_code_small_key"

	// 线路redis
	SYSTEM_LINE_REDIS_KEY = "red_system_line_key"

	// 菜单层级
	MENU_ONE   = 1
	MENU_TWO   = 2
	MENU_THREE = 3

	// 在线状态
	ONLINE  = 1
	OFFLINE = 2

	// 聊天类型
	ROOM_TYPE_SECRET = 1 // 私聊
	ROOM_TYPE_GRUOUP = 2 // 群聊

	// 分享类型
	SHARE_TYPE_USER = 1 // 用户
	SHARE_TYPE_ROOM = 2 // 群

	// im消息类型
	IM_MSG_TYPE_TEXT         = 1 // 文本
	IM_MSG_TYPE_PICTURE      = 2 // 图片
	IM_MSG_TYPE_VIDEO        = 3 // 视频
	IM_MSG_TYPE_VOICE        = 4 // 语音
	IM_MSG_TYPE_RED          = 5 // 红包
	IM_MSG_TYPE_ROOM         = 7 // 群内通知
	IM_MSG_TYPE_NOTIFICSTION = 7 // 通知

	//系统消息推送类型
	USER_RELATION_ADD    = 1  //添加好友请求
	USER_RELATION_ACCEPT = 2  //接受好友添加请求
	USER_RELATION_REFUSE = 3  //拒绝好友添加请求
	ROOM_RELATION_ADD    = 4  //申请入群请求
	ROOM_RELATION_ACCEPT = 5  //接受群添加请求
	ROOM_RELATION_REFUSE = 6  //拒绝群添加请求
	ROOM_USER_KICK       = 7  //被踢
	ROOM_ADMIN_ADD       = 8  //被委任管理权
	ROOM_ADMIN_DEL       = 9  //被取消管理权
	ROOM_USER_ADD        = 10 //被拉群通知
	ROOM_USER_MUTE       = 11 //被禁言

	ROOM_KICK_USER   = 12 //踢人
	ROOM_ADD_ADMIN   = 13 //委任管理权
	ROOM_DEL_ADMIN   = 14 //取消管理权
	ROOM_ADD_USER    = 15 //拉进群通知
	ROOM_ACCEPT_USER = 16 //同意入群通知
	ROOM_INFO_UPDATE = 17 //群信息修改
	ROOM_USER_OUT    = 18 //退出房间通知
	ROOM_MUTE_USER   = 19 //禁言
	ROOM_MUTE        = 20 //房间禁言
	ROOM_UNMUTE      = 21 //房间解除禁言
	ROOM_DEL         = 22 //群解散
	USER_UNMUTE      = 23 //会员解除禁言

	// 红包类型
	ORDINARY_RED_ENVELOPE  = 0 // 普通红包
	NIUNIU_RED_ENVELOPE    = 1 // 牛牛红包
	MINESWEEPER_RED_PACKET = 2 // 扫雷红包

	// 红包玩法 牛牛红包
	NIUNIU_CLASSIC_PLAY = 1 // 经典玩法
	NIUNIU_FLAT_PLAY    = 2 // 平倍玩法
	NIUNIU_SUPER_PLAY   = 3 // 超倍玩法

	// 红包玩法 扫雷红包
	MINESWEEPER_FIXED_ODDS   = 1 // 固定赔率
	MINESWEEPER_UNFIXED_ODDS = 2 // 不固定赔率

	// 红包状态
	RED_STATUS_NORMAL    = 1 // 进行中
	RED_STATUS_OVER      = 2 // 已结束
	RED_STATUS_INVALID   = 3 // 无效
	RED_STATUS_FINISHED  = 4 // 已领完
	RED_STATUS_NOT_START = 5 // 未开始

	// 红包结果 输 赢 无效
	RED_RESULT_WIN     = 1 // 赢
	RED_RESULT_LOSE    = 2 // 输
	RED_RESULT_INVALID = 3 // 无效

	// 现金流水数据来源类型
	MEMBER_CASH_NIUNIU      = 1 // 牛牛红包
	MEMBER_CASH_MINESWEEPER = 2 // 扫雷红包

	// 现金流水数据类型
	// 1发包、2返还、3赢利、4亏损
	MEMBER_RECORD_CREATE           = 1 // 发包
	MEMBER_RECORD_RETURBN          = 2 // 返还
	MEMBER_RECORD_WIN              = 3 // 赢利
	MEMBER_RECORD_LOSE             = 4 // 亏损
	TRANSFERRED_IN                 = 5 // 转入金额
	TRANSFERRED_OUT                = 6 // 转出金额
	MEMBER_RECORD_ORDINARY_CREATE  = 7 // 普通红包发包
	MEMBER_RECORD_ORDINARY_WIN     = 8 // 普通红包领取
	MEMBER_RECORD_ORDINARY_RETURBN = 9 // 普通红包返还

	//ordinary

	IS_THE_QUERY = 1 //是查询接口
	NOT_QUERY    = 2 //不是查询接口

	TRANS_TYPE_WALLET     = 1 // 钱包
	TRANS_TYPE_CONVERSION = 2 // 额度转换

	// 红包是否自动发送
	IS_AUTO_YES = 1 // 红包自动发送

	// 是否是普通红包
	IS_ORDINAY_YES = 1 // 红包自动发送

	// 聊天群机器人发红包开关
	ROBOT_SEND_PACKET_OFF = 1 // 开启
	ROBOT_SEND_PACKET_NO  = 2 // 关闭

	// 聊天室机器人抢红包开关
	ROBOT_GRAB_PACKET_OFF = 1 // 开启
	ROBOT_GRAB_PACKET_NO  = 2 // 关闭

	// 会员是否是机器人
	USER_IS_ROBOT_YES = 1 // 是
	USER_IS_ROBOT_NO  = 0 // 不是

	// 会员是否是群主
	USER_IS_GROUP_OWNER_YES = 1 // 是
	USER_IS_GROUP_OWNER_NO  = 2 // 不是

	// 群是否开启免死
	ROOM_FREE_FROM_DEATH_OFF = 1 // 开启
	ROOM_FREE_FROM_DEATH_NO  = 2 // 关闭

	// 日志类型
	LOG_TYPE_LOGIN = 1 // 登录日志
	LOG_TYPE_OTHER = 2 // 操作日志

)

// 红包类型
var RED_ENVELOPE_TYPE = map[int]string{
	ORDINARY_RED_ENVELOPE:  "普通红包",
	NIUNIU_RED_ENVELOPE:    "牛牛红包",
	MINESWEEPER_RED_PACKET: "扫雷红包",
}

// 牛牛红包玩法
var NIUNIU_PLAY = map[int]string{
	NIUNIU_CLASSIC_PLAY: "经典玩法",
	NIUNIU_FLAT_PLAY:    "平倍玩法",
	NIUNIU_SUPER_PLAY:   "超倍玩法",
}

// 扫雷红包玩法
var MINESWEEPER_PLAY = map[int]string{
	MINESWEEPER_FIXED_ODDS:   "固定赔率",
	MINESWEEPER_UNFIXED_ODDS: "不固定赔率",
}

var RED_ENVELOPE_TYPE_PLAY = map[int]map[int]string{
	NIUNIU_RED_ENVELOPE:    NIUNIU_PLAY,
	MINESWEEPER_RED_PACKET: MINESWEEPER_PLAY,
}

var NIUNIU_PLAY_ODDS = map[int]map[int]int{
	NIUNIU_CLASSIC_PLAY: map[int]int{
		1: 1,
		2: 1,
		3: 1,
		4: 1,
		5: 1,
		6: 1,
		7: 2,
		8: 2,
		9: 2,
		0: 3,
	},
	NIUNIU_FLAT_PLAY: map[int]int{
		1: 1,
		2: 1,
		3: 1,
		4: 1,
		5: 1,
		6: 1,
		7: 1,
		8: 1,
		9: 1,
		0: 1,
	},
	NIUNIU_SUPER_PLAY: map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
		8: 8,
		9: 9,
		0: 10,
	},
}

var MINESWEEPER_UNFIXED_PLAY = map[int]float64{
	2:  10,
	3:  6,
	4:  4,
	5:  3,
	6:  2.5,
	7:  2.1,
	8:  1.8,
	9:  1.6,
	10: 1.5,
}

func GetMemberListKey() string {
	return "red_member_session_list_key"
}

func GetAgencyListKey() string {
	return "red_agency_session_list_key"
}

func GetAdminListKey() string {
	return "red_admin_session_list_key"
}
