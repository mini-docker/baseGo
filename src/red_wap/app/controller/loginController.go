package controllers

// import (
// 	"baseGo/src/model"
// 	"baseGo/src/model/code"
// 	"baseGo/src/red_wap/app/controller/common"
// 	"baseGo/src/red_wap/app/middleware/validate"
// 	"baseGo/src/red_wap/app/server"
// )

// type LoginController struct {
// 	LoginReq struct {
// 		Account  string `json:"account" valid:"Must;ErrorCode(3026)"`  // 账号
// 		Password string `json:"password" valid:"Must;ErrorCode(3027)"` // 密码
// 	}
// }

// // 登录
// func (m LoginController) Login(ctx server.Context) error {
// 	req := &m.LoginReq
// 	if err := ctx.Validate(req); err != nil {
// 		return err
// 	}
// 	// 获取登陆ip
// 	ip := ctx.RealIP()
// 	device := ctx.Get(model.DEVICE).(int)
// 	info, err := SystemAgencyService.Login(req.Account, req.Password, ip, device)
// 	if err != nil {
// 		return common.HttpResultJsonError(ctx, err)
// 	}
// 	return common.HttpResultJson(ctx, info)
// }

// // 登出
// func (m LoginController) Logout(ctx server.Context) error {
// 	// 获取用户登陆信息
// 	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
// 	if err != nil {
// 		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
// 	}
// 	err = SystemAgencyService.Logout(user)
// 	if err != nil {
// 		return common.HttpResultJsonError(ctx, err)
// 	}
// 	return common.HttpResultJson(ctx, nil)
// }
