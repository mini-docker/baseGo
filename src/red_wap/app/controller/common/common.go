package common

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/red_wap/app/middleware/validate"
	"baseGo/src/red_wap/app/po"
	"baseGo/src/red_wap/app/server"
	"fmt"
	"math/rand"
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
	//res.Data = string(byteStr)
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
	//res.Data = string(byteStr)
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
	//res.Data = string(byteStr)
	//res.Data = PswEncrypt(string(byteStr)) // 加密
	//return ctx.JSON(200, res)
	return ctx.JSON(200, hr)
}

func RandSeq(sess *xorm.Session, lineId, agencyId string, n int) string {
	//letters := []rune("abcdefghijklmnopqrstuvwxyz")
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	account := fmt.Sprintf("%s_%s_%s", lineId, agencyId, string(b))
	// 验证账号是否存在
	has, _ := new(bo.User).ExistUser(sess, account)
	if has {
		RandSeq(sess, lineId, agencyId, n)
	}
	return account
}

func RandPassword(n int) string {
	//letters := []rune("abcdefghijklmnopqrstuvwxyz")
	letters := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
