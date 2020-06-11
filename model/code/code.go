package code

import "fmt"

const (
	ZH      = "zh"
	EN      = "en"
	TC      = "tc"
	LangKey = "lang"

	JSON_MARSHAL_ERROR   = 11086 // JSON序列化失败
	JSON_UNMARSHAL_ERROR = 11087 // JSON反序列化失败

	DEFAULT_CODE         = 1000 // 默认错误
	SYSTEM_ERROR         = 1001 // 系统错误
	OPERATION_FAILED     = 1005 // 操作失败
	RESOURCE_NOT_FOUND   = 1006 // 记录不存在
	LOGIN_INFO_GET_FAIL  = 1008 // 登陆信息获取失败
	UPDATE_FAILED        = 1009 // 更新失败
	NO_UPPERCASE_LETTER  = 1010
	NO_LOWERCASE_LETTER  = 1011
	NO_FIGURES           = 1012
	ILLEGAL_CHARACTERS   = 1013
	NEED_LOGIN           = 1014
	INSET_ERROR          = 1015
	USER_DOES_NOT_EXISTS = 1016
	IM_CONN_ERROR        = 1017
	QUERY_FAILED         = 1018
	DELETE_FAILED        = 1019 // 删除失败

	ACCOUNT_ALREADY_EXISTS               = 5000 // 该账号已经被注册
	ACCOUNT_DOES_NOT_EXIST               = 5001 // 账号不存在
	PASSWORD_ERRORS                      = 5002 // 密码不正确
	VERIFICATION_CODE_ERRORS_UPPER_LIMIT = 5003 // 验证码错误次数太多，请稍后再试
	VERIFICATION_FAILED                  = 5004 // 验证失败
	VERIFY_IMAGE_ACQUISTION_FAILED       = 5005 // 验证图片获取失败
	CAPTCHA_ERROR                        = 5006 // 验证码错误
	USER_NAME_NOT_RIGHT                  = 5007 // 用户名非法
	READ_AGREE_ERR                       = 5008 // 必须阅读并同意注册协议
	THE_PASSWORD_YOU_ENTERED_IS_WEAK     = 5009 // 您输入的密码较弱,请尝试数字+字母的组合
	REGISTER_ERROR                       = 5010 // 注册失败
	ACCOUNT_PASSWORD_ERR                 = 5011 // 账号或密码不正确
	SEND_MEMBER_ERROR                    = 3345 // 发送会员通知失败
	RELATION_ALREADY_EXIST               = 5060 // 关系已存在
	ADD_FAILED                           = 5015 // 添加失败
	PASSWORD_NOT_SAME                    = 5012 // 两次密码输入不一致
	DATA_NOT_EXIST                       = 5013 // 数据不存在
	ANSWER_NOT_SAME                      = 5014 // 原安全问题答案不匹配
	SET_PASSWORD_ERROR                   = 5016 // 设置密码失败
	DATA_EXIST                           = 5017 // 数据已存在
	NO_PERMISSION                        = 5061 // 权限不足
	SECRET_CHAT_IS_REQUIRED              = 5018 // 密聊对象不能为空
	GROUP_CHAT_ROOM_REQUIRED             = 5019 // 群聊房间不能为空
	MESSAGE_FAILED_TO_BE_SENT            = 5020 // 消息发送失败
	ROOM_NOT_EXSIT                       = 5021 // 群聊不存在
	FRIEND_RELATION_NOT_EXIST            = 5022 // 好友关系不存在
	SHARE_TIME_OUT                       = 5023 // 分享已过期
	APPLY_EXIST                          = 5024 // 正在申请中，请耐心等待
	ROOM_RELATION_NOT_EXIST              = 5025 // 群关系不存在
	ROOM_IS_FROBIDDEN                    = 5026 // 该群已被停用
	USER_IS_FROBIDDEN                    = 5027 // 你已被禁言

	GAME_GROUP_QUERY_FAILED                 = 3001 // 游戏群查询失败
	RED_PACKET_QUERY_FAILED                 = 3002 // 红包查询失败
	DEALER_RESULTS_GET_FAILED               = 3003 // 庄家结果获取失败
	RED_ENVELOPE_TYPE_IS_INCORRECT          = 3004 // 红包类型不正确
	MEMBER_BALANCE_UPDATE_FAILED            = 3005 // 会员余额更新失败
	WRITING_CASH_RECORD_FAILED              = 3006 // 写入现金记录失败
	MEMBER_INFORMATION_QUERY_FAILED         = 3007 // 会员信息查询失败
	INSUFFICIENNT_MEMBER_BALANCE            = 3008 // 会员余额不足
	RED_ENVELOPE_ADDITION_FAILED            = 3009 // 红包添加失败
	UPDATE_MEMBER_CAPITAL_FAILED            = 3010 // 更新会员红包保证金失败
	RED_ENVELOPE_INSUFFICIENT_AMOUNT        = 3011 // 红包金额不足
	RED_ENVELOPE_CLAIM_FAILED               = 3012 // 红包领取失败
	MEMBER_DEPOSIT_REFUND_FAILED            = 3013 // 会员押金返还失败
	ROOM_FAILURE_TO_INCREASE_PUMPING_AMOUNT = 3014 // 游戏群抽水金额增加失败
	RED_ENVELOPE_HAS_ENDED                  = 3015 // 来晚了 本轮红包已经结束
	RED_ENVELOPE_HAS_BEEN_STOLEN            = 3016 // 来晚了 本轮红包已抢光
	MARGIN_UPDATE_FAILED                    = 3017 // 保证金更新失败
	GROUP_PUMPING_AMOUNT_UPDATE_FAILED      = 3018 // 群抽水金额更新失败
	RED_PACKET_AMOUNT_CANNOT_BE_EMPTY       = 3019 // 红包金额不能为空
	RED_PACKET_NUMBER_CANNOT_BE_EMPTY       = 3020 // 红包数量不能为空
	ROOM_ID_CANNOT_BE_EMPTY                 = 3021 // 房间内ID不能为空
	RED_ENVELOPE_ID_CANNOT_BE_EMPTY         = 3022 // 红包ID不能为空
	RED_ENVELOPE_TYPE_CANNOT_BE_EMPTY       = 3023 // 红包类型不能为空
	RED_ENVELOPE_PLAY_CANNOT_BE_EMPTY       = 3024 // 红包玩法不能为空
	RED_ENVELOPE_NUM_NOT_ENOUGH             = 3025 // 红包数量不能小于最小数量
	ACCOUNT_CAN_NOT_BYE_EMPTY               = 3026 // 账号不能为空
	PASSWORD_CAN_NOT_BYE_EMPTY              = 3027 // 密码不能为空
	STATUS_CAN_NOT_BYE_EMPTY                = 3028 // 状态不能为空
	ROLE_CAN_NOT_BE_EMPTY                   = 3029 // 角色不能为空
	LINE_CAN_NOT_BE_EMPTY                   = 3030 // 线路不能为空
	ID_CAN_NOT_BE_EMPTY                     = 3031 // id不能为空
	GAME_TYPE_CAN_NOT_BE_EMPTY              = 3032 // 游戏类型不能为空
	GAME_NAME_CAN_NOT_BE_EMPTY              = 3033 // 游戏名称不能为空
	LINE_NAME_CAN_NOT_BE_EMPTY              = 3034 // 线路名称不能为空
	LINE_LIMIT_CAN_NOT_BE_EMPTY             = 3035 // 额度不能为空
	LINE_MEAL_CAN_NOT_BE_EMPTY              = 3036 // 线路套餐不能为空
	DOMAIN_CAN_NOT_BE_EMPTY                 = 3037 // 域名不能为空
	TRANS_TYPE_CAN_NOT_BE_EMPTY             = 3038 // 交易模式不能为空
	MD5_KEY_CAN_NOT_BE_EMPTY                = 3039 // md5Key不能为空
	RSA_PUB_KEY_CAN_NOT_BE_EMPTY            = 3040 // rsaPubKey不能为空
	RSA_PRI_KEY_CAN_NOT_BE_EMPTY            = 3041 // rsaPriKey不能为空
	MEAL_NAME_CAN_NOT_BE_EMPTY              = 3042 // 套餐名称不能为空
	NN_ROYALTY_CAN_NOT_BE_EMPTY             = 3043 // 牛牛红包提成不能为空
	SL_ROYALTY_CAN_NOT_BE_EMPTY             = 3044 // 扫雷红包提成不能为空
	PARENT_ID_CAN_NOT_BE_EMPTY              = 3045 // 父级菜单不能为空
	MENU_NAME_CAN_NOT_BE_EMPTY              = 3046 // 菜单名称不能为空
	MENU_ROUTE_CAN_NOT_BE_EMPTY             = 3047 // 菜单路由不能为空
	ROLE_NAME_CAN_NOT_BE_EMPTY              = 3048 // 角色名称不能为空
	ROOM_NAME_CAN_NOT_BE_EMPTY              = 3049 // 群名称不能为空
	MAX_MONEY_CAN_NOT_BE_EMPTY              = 3050 // 最大红包金额不能为空
	MIN_MONEY_CAN_NOT_BE_EMPTY              = 3051 // 最小红包金额不能为空
	GAME_PLAY_CAN_NOT_BE_EMPTY              = 3052 // 群游戏玩法不能为空
	ODDS_CAN_NOT_BE_EMPTY                   = 3053 // 赔率不能为空
	RED_NUM_CAN_NOT_BE_EMPTY                = 3054 // 红包个数不能为空
	ROYALTY_CAN_NOT_BE_EMPTY                = 3055 // 抽水比例不能为空
	GAME_TIME_CAN_NOT_BE_EMPTY              = 3056 // 游戏时间不能为空
	ACCOUNT_CAN_NOT_BE_LOGIN                = 3057 // 账号被管理员停用
	POST_TITLE_CAN_NOT_BE_EMPTY             = 3058 // 公告标题不能为空
	START_TIME_CAN_NOT_BE_EMPTY             = 3059 // 开始时间不能为空
	END_TIME_CAN_NOT_BE_EMPTY               = 3060 // 结束时间不能为空
	POST_CONTENT_CAN_NOT_BE_EMPTY           = 3061 // 公告内容不能为空
	ACTIVE_NAME_CAN_NOT_BE_EMPTY            = 3062 // 活动标题不能为空
	PICTURE_CAN_NOT_BE_EMPTY                = 3063 // 活动图片不能为空
	WHITE_IP_CAN_NOT_BE_EMPTY               = 3064 // ip白名单不能为空
	FILE_CAN_NOT_BE_EMPTY                   = 3065 // 没有获取到上传文件
	LINE_QUERY_FAILED                       = 3066 // 线路信息查询失败
	ROOM_TYPE_NOT_BE_EMPTY                  = 3067 // 群属性不能为空
	FREE_DEATH_TYPE_NOT_BE_EMPTY            = 3068 // 是否开启免死号不能为空
	ROBOT_SEND_PACKET_TIME_CAN_NOT_BE_EMPTY = 3069 // 开启机器人自动发包，自动发包时间不能为空
	AGENCY_ID_CAN_NOT_BE_EMPTY              = 3070 // 站点不能为空
	RED_PACKRT_STATUS_IS_INCORRECT          = 3071 // 红包状态不正确
	ROOM_NO_CAN_NOT_BE_EMPTY                = 3072 // 房间id不能为空
	RED_PACKRT_SEND_TIME_NOT_UP             = 3073 // 红包发送时间未到
	RED_PACKRT_THE_AMOUNT_IS_INCORRECT      = 3074 // 红包金额不正确
	ROBOT_NUM_IS_TOO_LARGE                  = 3075 // 单次生成机器人不能超过100个
	SITE_CAN_NOT_BE_DELETE                  = 3076 // 站点仍存在代理，不能删除
	SITE_NAME_CAN_NOT_BE_EMPTY              = 3077 // 站点名称不能为空
	CAN_NOT_GRAB_SELF_BAG                   = 3078 // 自己不能抢自己的红包
	CAN_NOT_GRAB_ONE_BAG_AGAIN              = 3079 // 您已领取过该红包，不能再次领取
	THE_BAG_IS_TIME_OUT                     = 3080 // 您的手慢了，红包已过期
	ACTIVE_LIMIT                            = 3081 // 最多只能启用5个活动图片
	MENU_NAME_EXIST                         = 3082 // 菜单名称已存在
	MENU_ROUTE_EXIST                        = 3083 // 菜单路由已存在

	PARENT_MENU_NOT_EXIST           = 5028 // 父节点菜单不存在
	PARENT_MENU_NOT_RIGHT           = 5029 // 父节点菜单不能是三级菜单
	ADMIN_ROLE_NOT_EXIST            = 5030 // 管理员角色不存在
	DEFAULT_ROLE_CAN_NOT_BE_DELETE  = 5031 // 默认角色不能被删除
	DEFAULT_ROLE_CAN_NOT_BE_UPDATE  = 5032 // 默认角色不能被修改
	APIURL_CAN_NOT_BE_EMPTY         = 5033 // 钱包模式，钱包api地址不能为空
	ACCOUNT_DISABLED                = 6000 // 账号停用
	LOGIN_FAIL                      = 6001 // 登陆失败
	CONNECT_FAIL                    = 6002 // 连接失败
	DOESNT_HAS_ENOUGH_AMOUNT        = 6003 // 用户余额不足
	NO_DATA                         = 6004 // 无数据
	TIME_TOO_SHORT                  = 6005 // 时间间隔太短
	LINE_ID_EXIST                   = 5034 // 线路id已存在
	LOGINIP_NOT_IN_WHITE_IP_ADDRESS = 5035 // 登陆ip不是常用ip,请联系管理员添加
	ADMIN_ROLE_HAS_BEEN_STOPED      = 5036 // 管理员角色被停用
	GROUP_HAS_ORDER_CAN_NOT_BE_DEL  = 5037 // 群存在未结算的注单不能删除
	// 5038，5039已用
	GET_VALIDMEMBER_NUM_ERROR       = 5040 // 获取有效游戏人数失败
	GET_STATISTICAL_DATA_ERROR      = 5041 // 获取统计数据失败
	GET_SITE_STATISTICAL_DATA_ERROR = 5042 // 获取统计数据失败
)

type Code struct {
	Zh string //中文
	En string //英文
	Tc string //繁体
}

type ErrCode = map[int]Code

var TianliangCode = ErrCode{
	DEFAULT_CODE:                            {"参数错误", "Parameter error", "參數錯誤"},
	SYSTEM_ERROR:                            {"系统错误", "system error", "系統錯誤"},
	LOGIN_INFO_GET_FAIL:                     {"登录信息获取失败", "login info get failed", "登錄信息獲取失敗"},
	UPDATE_FAILED:                           {"更新失败", "update failed", "更新失敗"},
	OPERATION_FAILED:                        {"操作失败", "operation error", "操作失敗"},
	ACCOUNT_ALREADY_EXISTS:                  {"该账号已经被注册", "Account already exists", "該賬號已經被註冊"},
	ACCOUNT_DOES_NOT_EXIST:                  {"账号不存在", "Account does not exist", "賬號不存在"},
	PASSWORD_ERRORS:                         {"密码错误", "password error", "密碼錯誤"},
	VERIFICATION_CODE_ERRORS_UPPER_LIMIT:    {"验证码错误次数太多，请稍后再试", "There are too many verification code errors. Please try again later.", "驗證碼錯誤次數太多，請稍後再試"},
	VERIFICATION_FAILED:                     {"验证失败", "verification failed.", "驗證失敗"},
	VERIFY_IMAGE_ACQUISTION_FAILED:          {"验证图片获取失败", "Verify image acquisition failed.", "驗證圖片獲取失敗"},
	CAPTCHA_ERROR:                           {"验证码错误", "captcha error", "驗證碼錯誤"},
	USER_NAME_NOT_RIGHT:                     {"用户名非法,6-12位字母或数字", "Illegal username", "用戶名非法,6-12位字母或數字"},
	READ_AGREE_ERR:                          {"必须阅读并同意注册协议", "must read and agree register protocol", "必須閱讀並同意註冊協議"},
	THE_PASSWORD_YOU_ENTERED_IS_WEAK:        {"您输入的密码较弱,请尝试6-12位数字+字母的组合", "The password you entered is weak, please try the combination of number + letter", "您輸入的密碼較弱,請嘗試6-12位數字+字母的組合"},
	REGISTER_ERROR:                          {"注册失败", "register error", "註冊失敗"},
	ACCOUNT_PASSWORD_ERR:                    {"密码不正确", "Incorrect password", "密碼不正確"},
	NO_UPPERCASE_LETTER:                     {"未包含大写字母", "No uppercase letter", "未包含大寫字母"},
	NO_LOWERCASE_LETTER:                     {"未包含小写字母", "No lowercase letter", "未包含小寫字母"},
	NO_FIGURES:                              {"未包含数字", "No figures", "未包含數字"},
	ILLEGAL_CHARACTERS:                      {"未含有非法字符", "Illegal characters", "未含有非法字符"},
	NEED_LOGIN:                              {"您还未登陆，请重新登陆", "please relogin", "您還未登陸，請重新登陸"},
	SEND_MEMBER_ERROR:                       {"发送消息失败", "Sending member notifications", "發送消息失敗"},
	ADD_FAILED:                              {"添加失败", "add failed", "添加失敗"},
	RELATION_ALREADY_EXIST:                  {"关系已经存在", "realation already exist", "關系已經存在"},
	PASSWORD_NOT_SAME:                       {"两次密码输入不一致", "password error", "兩次密碼輸入不一致"},
	INSET_ERROR:                             {"添加失败", "insert error", "添加失敗"},
	DATA_NOT_EXIST:                          {"数据不存在", "Data not exist", "數據不存在"},
	ANSWER_NOT_SAME:                         {"原安全问题答案不匹配", "Answer not same", "原安全問題答案不匹配"},
	USER_DOES_NOT_EXISTS:                    {"用户不存在", "User does not exist", "用戶不存在"},
	SET_PASSWORD_ERROR:                      {"设置密码失败", "set password error", "設置密碼失敗"},
	DATA_EXIST:                              {"数据已存在", "data exist", "數據已存在"},
	SECRET_CHAT_IS_REQUIRED:                 {"密聊对象不能为空", "Secret chat object cannot be empty.", "密聊對象不能為空"},
	GROUP_CHAT_ROOM_REQUIRED:                {"群聊房间不能为空", "Group chat room cannot be empty.", "群聊房間不能為空"},
	MESSAGE_FAILED_TO_BE_SENT:               {"消息发送失败", "Message failed to be sent.", "消息發送失敗"},
	NO_PERMISSION:                           {"权限不足", "no permission", "權限不足"},
	ROOM_NOT_EXSIT:                          {"群聊不存在", "room not exist", "群聊不存在"},
	FRIEND_RELATION_NOT_EXIST:               {"好友关系不存在", "friend relation not exist", "好友關系不存在"},
	IM_CONN_ERROR:                           {"连接失败", "im conn error", "連接失敗"},
	QUERY_FAILED:                            {"查询失败", "query error", "查詢失敗"},
	SHARE_TIME_OUT:                          {"分享已过期", "share is timeout", "分享已過期"},
	APPLY_EXIST:                             {"正在申请中，请耐心等待", "apply exist", "正在申請中，請耐心等待"},
	ROOM_RELATION_NOT_EXIST:                 {"群关系不存在", "room relation not exist", "群關系不存在"},
	ROOM_IS_FROBIDDEN:                       {"该群已被停用", "room is forbidden speak", "該群已被停用"},
	USER_IS_FROBIDDEN:                       {"你已被禁言", "you have been forbidden speak", "你已被禁言"},
	GAME_GROUP_QUERY_FAILED:                 {"游戏群查询失败", "Game group query failed", "遊戲群查詢失敗"},
	RED_PACKET_QUERY_FAILED:                 {"红包查询失败", "Red envelope query failed", "紅包查詢失敗"},
	DEALER_RESULTS_GET_FAILED:               {"庄家结果获取失败", "Dealer Results Get Failed", "莊家結果獲取失敗"},
	RED_ENVELOPE_TYPE_IS_INCORRECT:          {"红包类型不正确", "Red envelope type is incorrect", "紅包類型不正確"},
	MEMBER_BALANCE_UPDATE_FAILED:            {"会员余额更新失败", "Member balance update failed", "會員余額更新失敗"},
	WRITING_CASH_RECORD_FAILED:              {"写入现金记录失败", "Writing cash record failed", "寫現金記錄失敗"},
	MEMBER_INFORMATION_QUERY_FAILED:         {"会员信息查询失败", "Member information query failed", "會員信息查詢失敗"},
	INSUFFICIENNT_MEMBER_BALANCE:            {"会员余额不足", "Insufficient member balance", "會員余額不足"},
	RED_ENVELOPE_ADDITION_FAILED:            {"红包添加失败", "Red envelope addition failed", "紅包添加失敗"},
	UPDATE_MEMBER_CAPITAL_FAILED:            {"更新会员红包保证金失败", "Failed to update member red envelope margin", "更新會員紅包保證金失敗"},
	RED_ENVELOPE_INSUFFICIENT_AMOUNT:        {"红包金额不足", "Insufficient amount of red envelopes", "紅包金額不足"},
	RED_ENVELOPE_CLAIM_FAILED:               {"红包领取失败", "Red envelope collection failed", "紅包領取失敗"},
	MEMBER_DEPOSIT_REFUND_FAILED:            {"会员押金返还失败", "Member deposit refund failed", "會員押金返還失敗"},
	ROOM_FAILURE_TO_INCREASE_PUMPING_AMOUNT: {"游戏群抽水金额增加失败", "Game group pumping amount failed to increase", "遊戲群抽水金額增加失敗"},
	RED_ENVELOPE_HAS_ENDED:                  {"来晚了,本轮红包已经结束", "Coming late This round of red envelopes has ended", "來晚了,本輪紅包已經結束"},
	RED_ENVELOPE_HAS_BEEN_STOLEN:            {"来晚了,本轮红包已抢光", "It's late. This round of red envelopes has been exhausted.", "來晚了,本輪紅包已搶光"},
	MARGIN_UPDATE_FAILED:                    {"保证金更新失败", "Margin update failed.", "保證金更新失敗"},
	GROUP_PUMPING_AMOUNT_UPDATE_FAILED:      {"群抽水金额更新失败", "Group pumping amount update failed.", "群抽水金額更新失敗"},
	RED_PACKET_AMOUNT_CANNOT_BE_EMPTY:       {"红包金额不能为空", "Red envelope amount cannot be empty.", "紅包金額不能為空"},
	RED_PACKET_NUMBER_CANNOT_BE_EMPTY:       {"红包数量不能为空", "The number of red envelopes cannot be empty.", "紅包數量不能為空"},
	ROOM_ID_CANNOT_BE_EMPTY:                 {"房间ID不能为空", "Room ID cannot be empty.", "房間ID不能為空"},
	RED_ENVELOPE_ID_CANNOT_BE_EMPTY:         {"红包ID不能为空", "Red envelope ID cannot be empty.", "紅包ID不能為空"},
	RED_ENVELOPE_TYPE_CANNOT_BE_EMPTY:       {"红包类型不能为空", "Red envelope type cannot be empty.", "紅包類型不能為空"},
	RED_ENVELOPE_PLAY_CANNOT_BE_EMPTY:       {"红包玩法不能为空", "Red envelope play cannot be empty.", "紅包玩法不能為空"},
	RED_ENVELOPE_NUM_NOT_ENOUGH:             {"红包数量不能小于最小数量", "The number of red envelopes cannot be less than the minimum number.", "紅包數量不能小於最小數量"},
	PARENT_MENU_NOT_EXIST:                   {"父节点菜单不存在", "parent not exist", ""},
	PARENT_MENU_NOT_RIGHT:                   {"父节点菜单不能是三级菜单", "parent menu not right", ""},
	DELETE_FAILED:                           {"删除失败", "delete error", ""},
	ADMIN_ROLE_NOT_EXIST:                    {"管理员角色不存在", "admin role not exist", ""},
	ACCOUNT_DISABLED:                        {"账号被停用", "account disabled", ""},
	LOGIN_FAIL:                              {"登陆失败", "login fail", ""},
	CONNECT_FAIL:                            {"连接失败", "connect fail", ""},
	DOESNT_HAS_ENOUGH_AMOUNT:                {"用户余额不足", "doesnt has enough amount", ""},
	NO_DATA:                                 {"无数据", "no data", ""},
	TIME_TOO_SHORT:                          {"时间间隔太短", "time too short", ""},
	DEFAULT_ROLE_CAN_NOT_BE_DELETE:          {"默认角色不能被删除", "default role can not be delete", ""},
	DEFAULT_ROLE_CAN_NOT_BE_UPDATE:          {"默认角色不能被修改", "default role can not be update", ""},
	JSON_MARSHAL_ERROR:                      {"json序列化失败", "json marshal failed", "json序列化失敗"},
	JSON_UNMARSHAL_ERROR:                    {"json反序列化失败", "json unmarshal failed", "json反序列化失敗"},
	APIURL_CAN_NOT_BE_EMPTY:                 {"钱包模式,钱包api地址不能为空", "apiUrl can not be empty", ""},
	LINE_ID_EXIST:                           {"线路id已存在", "lineId is exist", ""},
	LOGINIP_NOT_IN_WHITE_IP_ADDRESS:         {"登陆ip不是常用ip,请联系管理员添加", "login ip not in whiteIpAddress", ""},
	ACCOUNT_CAN_NOT_BYE_EMPTY:               {"账号不能为空", "account can not be empty", ""},
	PASSWORD_CAN_NOT_BYE_EMPTY:              {"密码不能为空", "password can not be empty", ""},
	STATUS_CAN_NOT_BYE_EMPTY:                {"状态不能为空", "status can not be empty", ""},
	ROLE_CAN_NOT_BE_EMPTY:                   {"角色不能为空", "role can not be empty", ""},
	LINE_CAN_NOT_BE_EMPTY:                   {"线路不能为空", "line can not be empty", ""},
	ID_CAN_NOT_BE_EMPTY:                     {"id不能为空", "id can not be empty", ""},
	GAME_TYPE_CAN_NOT_BE_EMPTY:              {"游戏类型不能为空", "game type can not be empty", ""},
	GAME_NAME_CAN_NOT_BE_EMPTY:              {"游戏名称不能为空", "game name can not be empty", ""},
	LINE_NAME_CAN_NOT_BE_EMPTY:              {"线路名称不能为空", "line name can not be empty", ""},
	LINE_LIMIT_CAN_NOT_BE_EMPTY:             {"线路额度不能为空", "line limit can not be empty", ""},
	LINE_MEAL_CAN_NOT_BE_EMPTY:              {"线路套餐不能为空", "line meal can not be empty", ""},
	DOMAIN_CAN_NOT_BE_EMPTY:                 {"域名不能为空", "domain can not be empty", ""},
	TRANS_TYPE_CAN_NOT_BE_EMPTY:             {"交易模式不能为空", "trans type can not be empty", ""},
	MD5_KEY_CAN_NOT_BE_EMPTY:                {"md5Key不能为空", "md5Key can not be empty", ""},
	RSA_PUB_KEY_CAN_NOT_BE_EMPTY:            {"rsaPubKey不能为空", "aesPubKey can not be empty", ""},
	RSA_PRI_KEY_CAN_NOT_BE_EMPTY:            {"rsaPriKey不能为空", "aesPriKey can not be empty", ""},
	MEAL_NAME_CAN_NOT_BE_EMPTY:              {"套餐名称不能为空", "meal name can not be empty", ""},
	NN_ROYALTY_CAN_NOT_BE_EMPTY:             {"牛牛红包提成不能为空", "nnRoyalty can not be empty", ""},
	SL_ROYALTY_CAN_NOT_BE_EMPTY:             {"扫雷红包提成不能为空", "slRoyalty can not be empty", ""},
	PARENT_ID_CAN_NOT_BE_EMPTY:              {"父级菜单不能为空", "parent can not be empty", ""},
	MENU_NAME_CAN_NOT_BE_EMPTY:              {"菜单名称不能为空", "menu name can not be empty", ""},
	MENU_ROUTE_CAN_NOT_BE_EMPTY:             {"菜单路由不能为空", "menu route can not be empty", ""},
	ROLE_NAME_CAN_NOT_BE_EMPTY:              {"角色名称不能为空", "roleName can not be empty", ""},
	ROOM_NAME_CAN_NOT_BE_EMPTY:              {"群名称不能为空", "room name can not be empty", ""},
	MAX_MONEY_CAN_NOT_BE_EMPTY:              {"最大红包金额不能为空", "max money can not be empty", ""},
	MIN_MONEY_CAN_NOT_BE_EMPTY:              {"最小红包金额不能为空", "min money can not be empty", ""},
	GAME_PLAY_CAN_NOT_BE_EMPTY:              {"群游戏玩法不能为空", "game play can not be empty", ""},
	ODDS_CAN_NOT_BE_EMPTY:                   {"赔率不能为空", "odds can not be empty", ""},
	RED_NUM_CAN_NOT_BE_EMPTY:                {"红包个数不能为空", "red num can not be empty", ""},
	ROYALTY_CAN_NOT_BE_EMPTY:                {"抽水比例不能为空", "royalty can not be empty", ""},
	GAME_TIME_CAN_NOT_BE_EMPTY:              {"游戏时间不能为空", "game time can not be empty", ""},
	ACCOUNT_CAN_NOT_BE_LOGIN:                {"账号被管理员停用", "account is stoped by manager", ""},
	POST_TITLE_CAN_NOT_BE_EMPTY:             {"公告标题不能为空", "post title can not be empty", ""},
	START_TIME_CAN_NOT_BE_EMPTY:             {"开始时间不能为空", "start time can not be empty", ""},
	END_TIME_CAN_NOT_BE_EMPTY:               {"结束时间不能为空", "end time can not be empty", ""},
	POST_CONTENT_CAN_NOT_BE_EMPTY:           {"公告内容不能为空", "post content can not be empty", ""},
	ACTIVE_NAME_CAN_NOT_BE_EMPTY:            {"活动标题不能为空", "active name can not be empty", ""},
	PICTURE_CAN_NOT_BE_EMPTY:                {"活动图片不能为空", "active picture can not be empty", ""},
	WHITE_IP_CAN_NOT_BE_EMPTY:               {"ip白名单不能为空", "white ip can not be empty", ""},
	FILE_CAN_NOT_BE_EMPTY:                   {"没有获取到上传文件", "file can not be empty", ""},
	LINE_QUERY_FAILED:                       {"线路信息查询失败", "Line information query failed", "線路信息查詢失敗"},
	ROOM_TYPE_NOT_BE_EMPTY:                  {"群属性不能为空", "roomType can not be empty", ""},
	FREE_DEATH_TYPE_NOT_BE_EMPTY:            {"是否开启免死号不能为空", "free from death can not be empty", ""},
	ROBOT_SEND_PACKET_TIME_CAN_NOT_BE_EMPTY: {"开启机器人自动发包，自动发包时间不能为空", "robot send packet time can not be empty", ""},
	AGENCY_ID_CAN_NOT_BE_EMPTY:              {"站点不能为空", "site can not be empty", ""},
	RED_PACKRT_STATUS_IS_INCORRECT:          {"红包状态不正确", "Red envelope status is incorrect", "紅包狀態不正確"},
	ROOM_NO_CAN_NOT_BE_EMPTY:                {"房间id不能为空", "room no can not be empty", ""},
	RED_PACKRT_SEND_TIME_NOT_UP:             {"红包发送时间未到", "Red envelope sending time has not arrived", "紅包發送時間未到"},
	RED_PACKRT_THE_AMOUNT_IS_INCORRECT:      {"红包金额不正确", "Red envelope amount is incorrect", "紅包金額不正確"},
	ROBOT_NUM_IS_TOO_LARGE:                  {"单次生成机器人不能超过100个", "robot num can not set more than 100", ""},
	SITE_CAN_NOT_BE_DELETE:                  {"站点仍存在代理，不能删除", "site can not be deleted", ""},
	SITE_NAME_CAN_NOT_BE_EMPTY:              {"站点名称不能为空", "site name can not be empty", ""},
	CAN_NOT_GRAB_SELF_BAG:                   {"自己不能抢自己的红包", "can not grab self bag", ""},
	CAN_NOT_GRAB_ONE_BAG_AGAIN:              {"你已领取过该红包，不能再次领取", "can not grab one bag again", ""},
	THE_BAG_IS_TIME_OUT:                     {"您的手慢了，红包已过期", "the bag is time out", ""},
	ACTIVE_LIMIT:                            {"最多只能启用5个活动图片", "active limit", ""},
	ADMIN_ROLE_HAS_BEEN_STOPED:              {"管理员角色被停用", "system role has been stoped", ""},
	MENU_NAME_EXIST:                         {"菜单名称已存在", "menu name exist", ""},
	MENU_ROUTE_EXIST:                        {"菜单路由已存在", "menu route exist", ""},
	GROUP_HAS_ORDER_CAN_NOT_BE_DEL:          {"群存在未结算的注单不能删除", "group has order can not be delete", ""},
	GET_VALIDMEMBER_NUM_ERROR:               {"获取有效游戏人数失败", "get valid member error", ""},
	GET_STATISTICAL_DATA_ERROR:              {"获取统计数据失败", "get statistical data error", ""},
	GET_SITE_STATISTICAL_DATA_ERROR:         {"获取统计数据失败", "get site statustucal data error", ""},
}

// ErrCodes 获取所有的code
func ErrCodes() ErrCode {
	// 所有人的code码
	var code4People = []ErrCode{
		TianliangCode,
	}
	allCode := make(ErrCode, 0)
	for _, codes := range code4People {
		for k, code := range codes {
			_, ok := allCode[k]
			if ok {
				panic("code码重复:" + fmt.Sprintf("%d", k))
			}
			allCode[k] = code
		}
	}
	return allCode
}
