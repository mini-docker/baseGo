package controllers

import (
	"baseGo/src/fecho/xorm/help"
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
	"strconv"
)

type OrdinaryRedPacketController struct {
	ListReq struct {
		StartTime int    `json:"startTime"` // 开始时间
		EndTime   int    `json:"endTime"`   // 结束时间
		Status    int    `json:"status"`
		RoomName  string `json:"roomName"`
		help.PageParams
		AgencyId string `json:"agencyId"` // 站点id
	}
	GetRedInfoReq struct {
		RedId int `json:"redId"` // 红包ID
	}
	DelRedInfoReq struct {
		Id int `json:"id"` // 红包ID
	}
	EditRedInfoReq struct {
		Id        int    `json:"Id"`        // 红包id
		RedAmount string `json:"redAmount"` // 红包金额
		RedNum    int    `json:"redNum"`    // 红包最大个数
		RoomId    int    `json:"roomId"`    // 群ID
		IsAuto    int    `json:"isAuto"`    // 是否开启自动 1是 0否
		AutoTime  int    `json:"autoTime"`  // 自动开始时间
		GameTime  int    `json:"gameTime"`  // 游戏时间
	}
}

// 普通红包列表查询
func (c OrdinaryRedPacketController) GetRedList(ctx server.Context) error {
	// 接收参数 红包金额 红包数量 房间ID 红包类型 红包玩法
	req := &c.ListReq
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

	count, res, err := new(services.RedPacketService).GetOrdinaryRedList(session.User.LineId, req.AgencyId, req.StartTime, req.EndTime, req.Status, req.RoomName, &req.PageParams)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, &map[string]interface{}{
		"data":  res,
		"count": count,
	})
}

// 普通红包详情查询
func (c OrdinaryRedPacketController) GetRedInfo(ctx server.Context) error {
	// 接收参数 红包金额 红包数量 房间ID 红包类型 红包玩法
	req := &c.GetRedInfoReq
	err := ctx.Validate(req)
	if err != nil {
		return err
	}
	// 获取登陆session
	session, err := new(middleware.SessionService).GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	res, err := new(services.RedPacketService).GetOrdinaryRedInfo(session.User.LineId, session.User.AgencyId, req.RedId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, res)
}

// 普通红包修改
func (c OrdinaryRedPacketController) EditRedInfo(ctx server.Context) error {
	// 接收参数 红包金额 红包数量 房间ID 红包类型 红包玩法
	req := &c.EditRedInfoReq
	err := ctx.Validate(req)
	if err != nil {
		return err
	}
	// 获取登陆session
	session, err := new(middleware.SessionService).GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	RedAmount, _ := strconv.ParseFloat(req.RedAmount, 64)
	err = new(services.RedPacketService).EditOrdinaryRedInfo(session.User.LineId, session.User.AgencyId, req.Id,
		RedAmount,
		req.RedNum,
		req.RoomId,
		req.IsAuto,
		req.AutoTime,
		req.GameTime)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 普通红包删除
func (c OrdinaryRedPacketController) DelRedInfo(ctx server.Context) error {
	// 接收参数 红包金额 红包数量 房间ID 红包类型 红包玩法
	req := &c.DelRedInfoReq
	err := ctx.Validate(req)
	if err != nil {
		return err
	}
	// 获取登陆session
	session, err := new(middleware.SessionService).GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	err = new(services.RedPacketService).DelOrdinaryRedInfo(session.User.LineId, session.User.AgencyId, req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

func (c OrdinaryRedPacketController) Orders(ctx server.Context) error {
	rs, err := new(services.RedPacketService).Orders("aaa", "a", 183)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, rs)
}
