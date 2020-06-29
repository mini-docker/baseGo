package controllers

import (
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_wap/app/controller/common"
	"baseGo/src/red_wap/app/middleware"
	"baseGo/src/red_wap/app/middleware/validate"
	"baseGo/src/red_wap/app/server"
	"baseGo/src/red_wap/app/services"
	"fmt"
)

var (
	SessionService = new(middleware.UserSessionService)
	UserService    = new(services.UserService)
)

type UserController struct{}

func (ac UserController) ImCoon(ctx server.Context) error {
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}

	sid := ctx.Request().Header.Get(model.SessionKey)
	device := ctx.Get(model.DEVICE).(int)
	platform := ""
	if device == model.IS_ANDROID {
		platform = "android"
	} else if device == model.IS_IOS {
		platform = "ios"
	} else if device == model.IS_WAP {
		platform = "wap"
	}

	data, err := UserService.LoginIm(session.User.LineId, fmt.Sprint(session.User.AgencyId), sid, platform, ctx.Request().Host, int64(session.Uid()), session.User.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	return common.HttpResultJson(ctx, data.Data)
}

// 获取会员详情
func (ac UserController) GetUserInfo(ctx server.Context) error {
	// 获取登陆session
	session, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	res, err := UserService.GetUserInfo(session.User.LineId, session.User.AgencyId, session.Account())
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, res)
}
