package controller

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/model"
	"baseGo/src/red_api/app/controller/common"
	"baseGo/src/red_api/app/server"
	"baseGo/src/red_api/app/services"
	"baseGo/src/red_api/conf"
)

type UserController struct {
	LoginReq struct {
		// 账号
		Account string `json:"account"`
		// 密码
		Password string `json:"password"`
		//上次登陆ip
		Ip string `json:"ip"`
	}
	SaveRegisterReq struct {
		// 账号
		Account string `json:"account"`
		// 密码
		Password string `json:"password"`
		//上次登陆ip
		Ip string `json:"ip"`
	}
	UserInfoReq struct {
		// 账号
		Account string `json:"account"`
	}
}
type UserLoginResp struct {
	SessionId string `json:"sessionId"`
	GameUrl   string `json:"gameUrl"`
}

func (ac UserController) Login(ctx server.Context) error {
	req := &ac.LoginReq
	err := ctx.Validate(req)
	if err != nil {
		golog.Error("User", "Login", "登陆参数解析失败", err)
		return common.HttpResultJsonError(ctx, err)
	}

	lineId := ctx.Get(model.LineId).(string)
	agencyId := ctx.Get(model.AgencyId).(string)
	device := ctx.Get(model.DEVICE).(int)

	session, err := new(services.UserService).AttemptLogin(lineId, agencyId, req.Account, req.Password, req.Ip, device)
	if err != nil {
		golog.Error("User", "Login", "登陆失败", err)
		return common.HttpResultJsonError(ctx, err)
	}
	rs := new(UserLoginResp)
	rs.SessionId = session.SessionId
	rs.GameUrl = conf.GetUrlConfig().Game
	return common.HttpResultJson(ctx, rs)
}
func (ac UserController) Register(ctx server.Context) error {
	req := &ac.SaveRegisterReq
	err := ctx.Validate(req)
	if err != nil {
		golog.Error("User", "Register", "参数解析失败", err)
		return common.HttpResultJsonError(ctx, err)
	}
	lineId := ctx.Get(model.LineId).(string)
	agencyId := ctx.Get(model.AgencyId).(string)
	device := ctx.Get(model.DEVICE).(int)

	session, err := new(services.UserService).Register(lineId, agencyId, req.Account, req.Password, req.Ip, device)
	if err != nil {
		golog.Error("User", "Register", "注册失败", err)
		return common.HttpResultJsonError(ctx, err)
	}
	rs := new(UserLoginResp)
	rs.SessionId = session.SessionId
	rs.GameUrl = conf.GetUrlConfig().Game
	return common.HttpResultJson(ctx, rs)
}

// 获取会员详情
func (ac UserController) GetUserInfo(ctx server.Context) error {
	req := &ac.UserInfoReq
	err := ctx.Validate(req)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	lineId := ctx.Get(model.LineId).(string)
	agencyId := ctx.Get(model.AgencyId).(string)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	res, err := new(services.UserService).GetUserInfo(lineId, agencyId, req.Account)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, res)
}
