package common

import (
	"baseGo/src/model/code"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/app/po"
	"baseGo/src/red_admin/app/server"
)

//data := `{"code":"0","data":"GMT+08:00","message":"请求成功","success":true,"title":"","version":"app_01"}`
func HttpResultJson(ctx server.Context, data interface{}) error {
	hr := new(po.HttpResult)
	hr.Code = 0
	hr.Data = data
	hr.Message = "success"
	hr.Success = true
	hr.Version = "app_01"
	//byteStr, _ := json.Marshal(hr)
	//res := new(po.HttpRes)
	//res.Data = PswEncrypt(string(byteStr)) // 加密
	//err := ctx.JSON(200, res)
	err := ctx.JSON(200, hr)
	if err != nil {
		return HttpResultJsonError(ctx, err) // 9001 json解析错误
	}
	return nil
}

//{"code":"1001","data":null,"message":"您还未登录,请重新登录","success":false,"title":"","version":"app_01"}
func HttpResultJsonMsg(ctx server.Context, code int, data string) error {
	hr := new(po.HttpResult)
	hr.Code = code
	hr.Message = data
	hr.Success = false
	hr.Version = "app_01"
	//byteStr, _ := json.Marshal(hr)
	//res := new(po.HttpRes)
	//res.Data = PswEncrypt(string(byteStr)) // 加密
	//return ctx.JSON(200, res)
	return ctx.JSON(200, hr)
}

func HttpResultJsonError(ctx server.Context, data error) error {

	hr := new(po.HttpResult)
	hr.Success = false
	hr.Version = "app_01"

	validateErr, ok := data.(*validate.Err)
	if ok {
		hr.Code = validateErr.Code
		lang, ok := ctx.Get(code.LangKey).(string)
		if !ok {
			lang = validate.ZH
		}
		if validateErr.Msg == "" {
			hr.Message = validate.Find(validateErr.Code, lang)
		} else {
			hr.Message = validate.Find(validateErr.Code, lang) + " " + validateErr.Msg
		}
	} else {
		hr.Code = 9001
		hr.Message = data.Error()
	}
	//byteStr, _ := json.Marshal(hr)
	//res := new(po.HttpRes)
	//res.Data = PswEncrypt(string(byteStr)) // 加密
	//return ctx.JSON(200, res)
	return ctx.JSON(200, hr)
}
