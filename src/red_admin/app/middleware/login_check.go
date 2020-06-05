package middleware

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/red_admin/app/po"
	"baseGo/src/red_admin/app/server"
	"net/http"
)

var sessionService = new(SessionService)

func NotAuthInit(next server.HandlerFunc) server.HandlerFunc {
	return func(c server.Context) error {
		hd, err := getHd(c)
		if err != nil || hd == nil {
			return NeedLoginJsonMsg(c, code.OPERATION_FAILED)
		}
		c.Set(model.DEVICE, hd.Terminal)

		//// 读取body
		//s, _ := ioutil.ReadAll(c.Request().Body)
		//// 解密
		//if len(s) > 0 {
		//	str := common.PswDecrypt(string(s))
		//	// 重新写入
		//	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer([]byte(str)))
		//}
		return next(c)
	}
}

func AuthInit(next server.HandlerFunc) server.HandlerFunc {
	return func(c server.Context) error {
		// 获取cookie信息
		sid := c.Request().Header.Get(model.SessionKey)
		if sid == "" {
			return NeedLoginJsonMsg(c, code.LOGIN_INFO_GET_FAIL)
		}
		hd, err := getHd(c)
		if err != nil || hd == nil {
			return NeedLoginJsonMsg(c, code.LOGIN_INFO_GET_FAIL)
		}
		// 获取session信息
		session, err := sessionService.GetSession(sid)
		if err == nil {
			if session.TimeOut < utility.GetNowTimestamp() {
				return NeedLoginJsonMsg(c, code.LOGIN_INFO_GET_FAIL)
			}
			err = sessionService.SaveSession(model.RED_ADMIN_SESSION_LIST_KEY, session)

			if err != nil {
				golog.Error("middleware", "AuthInit", "err:%v", err)
				return server.NewHTTPError(http.StatusInternalServerError, "save user info error")
			}
			// session存在直接通过写入ctx
			c.Set(model.SessionKey, session.SessionId)
			c.Set(model.DEVICE, hd.Terminal)
			if c.Path() == "/api/goimfiles/upload" {
				return next(c)
			}
			// 读取body
			//s, _ := ioutil.ReadAll(c.Request().Body)
			//if len(s) > 0 {
			//	// 解密
			//	str := common.PswDecrypt(string(s))
			//	// 重新写入
			//	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer([]byte(str)))
			//}
		} else {
			return NeedLoginJsonMsg(c, code.LOGIN_INFO_GET_FAIL)
		}
		return next(c)
	}
}

func NeedLoginJsonMsg(ctx server.Context, code int) error {
	hr := new(po.HttpResult)
	hr.Code = code
	hr.Message = "您还未登录,请重新登录" //TODO 多语言
	hr.Success = false
	hr.Version = ""
	//byteStr, _ := json.Marshal(hr)
	//res := new(po.HttpRes)
	//res.Data = common.PswEncrypt(string(byteStr)) // 加密
	//return ctx.JSON(200, res)
	return ctx.JSON(200, hr)
}

func getHd(c server.Context) (*po.HttpHeaderData, error) {
	// 得到客户设备类型
	platform := c.Request().Header.Get("platform")
	hd := new(po.HttpHeaderData)
	switch platform {
	case model.ANDROID:
		hd.Terminal = model.IS_ANDROID
	case model.IOS:
		hd.Terminal = model.IS_IOS
	case model.WAP:
		hd.Terminal = model.IS_WAP
	case model.PC:
		hd.Terminal = model.IS_PC
	}
	return hd, nil
}
