package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
	"strconv"
)

var (
	RoomService      = new(services.RoomService)
	RedPacketService = new(services.RedPacketService)
)

type RoomController struct {
	AddRoomReq struct {
		RoomName      string `json:"roomName"`                                   // 群名称
		GameType      int    `json:"gameType" valid:"Must;ErrorCode(3032)"`      // 房间 2扫雷 1牛牛
		AgencyId      string `json:"agencyId"`                                   // 代理id
		MaxMoney      string `json:"maxMoney" valid:"Must;ErrorCode(3050)"`      // 最大红包金额
		MinMoney      string `json:"minMoney" valid:"Must;ErrorCode(3051)"`      // 最小红包金额
		GamePlay      int    `json:"gamePlay" valid:"Must;ErrorCode(3052)"`      // 房间玩法 牛牛 1经典牛牛 2平倍牛牛 3超倍牛牛 扫雷 1固定赔率 2不固定赔率
		Odds          string `json:"odds" valid:"Must;ErrorCode(3053)"`          // 赔率
		RedNum        int    `json:"redNum" valid:"Must;ErrorCode(3054)"`        // 红包个数
		RedMinNum     int    `json:"redMinNum"`                                  // 红包最小个数
		Royalty       string `json:"royalty" valid:"Must;ErrorCode(3055)"`       // 抽水比例
		GameTime      int    `json:"gameTime" valid:"Must;ErrorCode(3056)"`      // 房间时间
		RoomType      int    `json:"roomType" valid:"Must;ErrorCode(3056)"`      // 群类型  1公群  2私群
		FreeFromDeath int    `json:"freeFromDeath" valid:"Must;ErrorCode(3056)"` // 是否开启免死号  1开启  2关闭
	}

	QueryRoomListReq struct {
		StartTime int    `json:"startTime"` // 开始时间
		EndTime   int    `json:"endTime"`   // 结束时间
		GameType  int    `json:"gameType"`  // 房间类型
		Status    int    `json:"status"`    // 状态  1 启用  2 停用
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
		AgencyId  string `json:"agencyId"`  // 站点id
	}

	QueryRoomReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)"` // 房间id
	}

	EditRoomReq struct {
		Id                  int    `json:"id" valid:"Must;ErrorCode(3031)"`            // 房间id
		RoomName            string `json:"roomName"`                                   // 群名称
		GameType            int    `json:"gameType" valid:"Must;ErrorCode(3032)"`      // 房间 2扫雷 1牛牛
		MaxMoney            string `json:"maxMoney" valid:"Must;ErrorCode(3050)"`      //  最大红包金额
		MinMoney            string `json:"minMoney" valid:"Must;ErrorCode(3051)"`      // 最小红包金额
		GamePlay            int    `json:"gamePlay" valid:"Must;ErrorCode(3052)"`      // 房间玩法 牛牛 1经典牛牛 2平倍牛牛 3超倍牛牛 扫雷 1固定赔率 2不固定赔率
		Odds                string `json:"odds" valid:"Must;ErrorCode(3053)"`          // 赔率
		RedNum              int    `json:"redNum" valid:"Must;ErrorCode(3054)"`        // 红包个数
		RedMinNum           int    `json:"redMinNum"`                                  // 红包最小个数
		Royalty             string `json:"royalty" valid:"Must;ErrorCode(3055)"`       // 抽水比例
		GameTime            int    `json:"gameTime" valid:"Must;ErrorCode(3056)"`      // 房间时间
		RoomSort            int    `json:"roomSort"`                                   // 排序
		RoomType            int    `json:"roomType" valid:"Must;ErrorCode(3056)"`      // 群类型  1公群  2私群
		FreeFromDeath       int    `json:"freeFromDeath" valid:"Must;ErrorCode(3056)"` // 是否开启免死号  1开启  2关闭
		RobotSendPacket     int    `json:"robotSendPacket"`                            // 机器人发包状态 1开启 2关闭
		RobotSendPacketTime int    `json:"robotSendPacketTime"`                        // 机器人发包时间
		RobotGrabPacket     int    `json:"robotGrabPacket"`                            // 机器人抢包开关 1开启 2关闭
		ControlKill         int    `json:"controlKill"`                                // 控杀开关  1 开启 2 关闭
	}

	EditRoomStatusReq struct {
		Id     int `json:"id" valid:"Must;ErrorCode(3031)"`     // 房间id
		Status int `json:"status" valid:"Must;ErrorCode(3028)"` // 状态
	}

	DelRoomReq struct {
		Id int `json:"id"` // 房间id
	}
	RoomListReq struct {
		GameType int `json:"gameType" valid:"Must;ErrorCode(3023)"`
		GamePlay int `json:"gamePlay" valid:"Must;ErrorCode(3024)"`
	}
	AddRedReq struct {
		AgencyId  string `json:"agencyId"`  // 代理id
		RedAmount string `json:"redAmount"` // 红包金额
		RedNum    int    `json:"redNum"`    // 红包最大个数
		RoomId    int    `json:"roomId"`    // 群ID
		IsAuto    int    `json:"isAuto"`    // 是否开启自动 1是 0否
		AutoTime  int    `json:"autoTime"`  // 自动开始时间
		GameTime  int    `json:"gameTime"`  // 游戏时间
	}

	RoomCodeReq struct {
		AgencyId string `json:"agencyId"` // 站点id
		GameType int    `json:"gameType"` // 游戏类型
	}
}

// 查询全部房间
func (m RoomController) QueryRoomList(ctx server.Context) error {
	req := &m.QueryRoomListReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	if session.IsAdmin != 1 {
		req.AgencyId = session.User.AgencyId
	}

	result, err := RoomService.QueryRoomList(session.User.LineId, req.AgencyId, req.StartTime, req.EndTime, req.GameType, req.Status, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加房间
func (m RoomController) AddRoom(ctx server.Context) error {
	req := &m.AddRoomReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	if session.IsAdmin == 1 {
		if req.AgencyId == "" {
			return common.HttpResultJsonError(ctx, &validate.Err{Code: code.AGENCY_ID_CAN_NOT_BE_EMPTY})
		}
	} else {
		req.AgencyId = session.User.AgencyId
	}
	//err = RoomService.AddRoom(session.User.LineId, session.User.AgencyId, req.RoomName, req.GameType, req.GamePlay, req.MaxMoney, req.MinMoney, req.Odds, req.RedNum, req.Royalty, req.GameTime, req.RedMinNum)
	err = RoomService.AddRoom(session.User.LineId, req.AgencyId, req.RoomName, req.GameType, req.GamePlay, req.MaxMoney, req.MinMoney,
		req.Odds, req.RedNum, req.Royalty, req.GameTime, req.RedMinNum, req.RoomType, req.FreeFromDeath, session.Uid(),
		session.Account(), ip, session.IsAdmin)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询房间信息
func (m RoomController) QueryRoomOne(ctx server.Context) error {
	req := &m.QueryRoomReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	_, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	Room, err := RoomService.QueryRoomOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, Room)
}

// 修改房间
func (m RoomController) EditRoom(ctx server.Context) error {
	req := &m.EditRoomReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	//err = RoomService.EditRoom(session.User.LineId, session.User.AgencyId, req.Id, req.RoomName, req.GameType, req.GamePlay, req.MaxMoney, req.MinMoney, req.Odds, req.RedNum, req.Royalty, req.GameTime, req.RoomSort, req.RedMinNum)
	err = RoomService.EditRoom(req.Id, req.RoomName, req.GameType,
		req.GamePlay, req.MaxMoney, req.MinMoney, req.Odds, req.RedNum, req.Royalty, req.GameTime,
		req.RoomSort, req.RedMinNum, req.RoomType, req.FreeFromDeath, req.RobotSendPacket, req.RobotSendPacketTime,
		req.RobotGrabPacket, req.ControlKill, session.Uid(), session.Account(), ip, session.IsAdmin)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改房间状态
func (m RoomController) EditRoomStatus(ctx server.Context) error {
	req := &m.EditRoomStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	err = RoomService.EditRoomStatus(req.Id, req.Status, session.Uid(), session.Account(), ip, session.IsAdmin)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 删除房间
func (m RoomController) DelRoom(ctx server.Context) error {
	req := &m.DelRoomReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取登陆信息
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	ip := ctx.RealIP()
	err = RoomService.DelRoom(req.Id, session.Uid(), session.Account(), ip, session.IsAdmin)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 发送普通红包
func (c *RoomController) AddRed(ctx server.Context) error {
	//ctx // 获取线路信息和超管信息
	req := &c.AddRedReq
	err := ctx.Validate(req)
	if err != nil {
		return err
	}
	// 获取登陆session
	session, err := new(middleware.SessionService).GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	if session.IsAdmin != 1 {
		req.AgencyId = session.AgencyId()
	} else {
		if req.AgencyId == "" {
			return common.HttpResultJsonError(ctx, &validate.Err{Code: code.AGENCY_ID_CAN_NOT_BE_EMPTY})
		}
	}

	// 开启自动之后写入定时任务  到时间之后自动发红包
	// 未开启 即时发红包 所以红包自动发送时间不写入红包记录中
	// 前面的处理是一样的 插入红包 需要根据红包自动开启时间来
	RedAmount, _ := strconv.ParseFloat(req.RedAmount, 64)
	err = RedPacketService.CreateOrdinaryRedPacket(RedAmount, req.RedNum, req.RoomId, req.IsAuto, req.AutoTime, req.GameTime)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 群枚举
func (c *RoomController) RoomCode(ctx server.Context) error {
	req := &c.RoomCodeReq
	err := ctx.Validate(req)
	if err != nil {
		return err
	}
	// 获取登陆session
	session, err := new(middleware.SessionService).GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	if session.IsAdmin != 1 {
		req.AgencyId = session.AgencyId()
	}
	roomCode, err := RoomService.RoomCode(session.LineId(), req.AgencyId, req.GameType)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, roomCode)
}
