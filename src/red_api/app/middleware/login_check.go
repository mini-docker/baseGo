package middleware

import (
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_api/app/controller/common"
	"baseGo/src/red_api/app/po"
	"baseGo/src/red_api/app/server"
	"baseGo/src/red_api/conf"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fecho/golog"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var sessionService = new(SessionService)

var (
	SystemLineBo = new(bo.SystemLineBo)
)

func NotAuthInit(next server.HandlerFunc) server.HandlerFunc {
	return func(c server.Context) error {

		hd, err := getHd(c)
		if err != nil || hd == nil {
			golog.Error("loginCheck", "AuthInit", "header获取失败", err)
			return NeedLoginJsonMsg(c, code.OPERATION_FAILED)
		}
		c.Set(model.DEVICE, hd.Terminal)

		// 读取body
		s, _ := ioutil.ReadAll(c.Request().Body)
		//lineId
		lineId := c.Request().Header.Get(model.LineId)
		//agencyId
		agencyId := c.Request().Header.Get(model.AgencyId)
		//sign
		sign := c.Request().Header.Get(model.Sign)
		golog.Info("loginCheck", "AuthInit", "获取参数：", fmt.Sprintf("body:%s,lineId:%s,agencyId:%s,sign:%s", string(s), lineId, agencyId, sign))

		c.Set(model.LineId, lineId)
		c.Set(model.AgencyId, agencyId)

		line, err := getAesPriKeyAndMd5Key(lineId)
		if err != nil {
			golog.Error("loginCheck", "AuthInit", "获取线路失败", err)
			return err
		}
		priKey := getPrikey(line.RsaPriKey)
		golog.Info("loginCheck", "AuthInit", "线路私钥,盐：", priKey, line.Md5key)
		if len(s) > 0 {
			decBody, checkResult := DecRequest(s, line.Md5key, priKey, sign)
			if !checkResult {
				golog.Error("loginCheck", "AuthInit", "解密失败", nil)
				return server.NewHTTPError(http.StatusInternalServerError, "dec error")
			}
			golog.Info("loginCheck", "AuthInit", "解析参数", string(decBody))
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(decBody))
		}
		// 解密

		return next(c)
	}
}

func NeedLoginJsonMsg(ctx server.Context, code int) error {
	hr := new(po.HttpResult)
	hr.Code = code
	hr.Message = "您还未登录,请重新登录" //TODO 多语言
	byteStr, _ := json.Marshal(hr)
	res := new(po.HttpRes)
	//res.Data = string(byteStr)
	res.Data = common.PswEncrypt(string(byteStr)) // 加密
	return ctx.JSON(200, res)
	//return ctx.JSON(200, hr)
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
	}
	return hd, nil
}

func getAesPriKeyAndMd5Key(lineId string) (*structs.SystemLine, error) {

	conn := conf.GetRedis()
	redisClient := conn.Get()
	line := new(structs.SystemLine)
	back, err := redisClient.Do("HGet", model.SYSTEM_LINE_REDIS_KEY, lineId)
	// back = back.(byte)
	fResults := string(back.([]byte))
	if err != nil {
		// if err != redis.Nil {
		// 	golog.Error("loginLock", "getAesPriKeyAndMd5Key", "查询redis失败", err)
		// 	return nil, server.NewHTTPError(http.StatusInternalServerError, "redis get err")
		// } else {
		// 	golog.Error("loginLock", "getAesPriKeyAndMd5Key", "查询redis失败", err)
		// 	goto REDISNILOUT
		// }
		golog.Error("loginLock", "getAesPriKeyAndMd5Key", "查询redis失败", err)
		goto REDISNILOUT

	} else {
		if err = json.Unmarshal([]byte(fResults), &line); err != nil {
			return nil, server.NewHTTPError(http.StatusInternalServerError, "redis unmarshal err")
		}
		return line, nil
	}
REDISNILOUT:
	xormSess := conf.GetXormSession()
	defer xormSess.Close()

	line, has, _ := SystemLineBo.QueryLineBylineId(xormSess, lineId)
	if !has {
		golog.Error("loginLock", "getAesPriKeyAndMd5Key", "查询线路信息失败", err)
		return nil, server.NewHTTPError(http.StatusInternalServerError, "line doesnt exist")
	}

	go func() {
		if lineStr, err := json.Marshal(line); err != nil {
			_, err = redisClient.Do("HSet", model.SYSTEM_LINE_REDIS_KEY, lineId, lineStr)
			if err != nil {
				return
			}
		}

	}()
	return line, nil
}

func DecRequest(requerstBody []byte, md5Key, rsaPriKey, sign string) (decBody []byte, checkResult bool) {

	s := NewPasswordEncrypt(string(requerstBody), md5Key)
	if s == sign {
		checkResult = true
	} else {
		return
	}

	bodyBase64code, err := base64.StdEncoding.DecodeString(string(requerstBody))

	if err != nil {
		return
	}
	block, _ := pem.Decode([]byte(rsaPriKey))

	if block == nil {
		return
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return
	}
	decBody, err = rsa.DecryptPKCS1v15(rand.Reader, priv, bodyBase64code)

	return

}

func NewPasswordEncrypt(password, salt string) string {
	return Md5(password + salt)
}
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func getPrikey(key string) string {
	key = strings.Replace(key, "-----BEGIN RSA PRIVATE KEY-----", "", -1)
	key = strings.Replace(key, "-----END RSA PRIVATE KEY-----", "", -1)
	newKey := `
-----BEGIN RSA PRIVATE KEY-----
` + key + `
-----END RSA PRIVATE KEY-----`
	return newKey
}
