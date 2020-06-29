package services

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	registry_module "baseGo/src/fecho/registry/registry-module"
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_wap/app/middleware"
	"baseGo/src/red_wap/app/middleware/validate"
	"baseGo/src/red_wap/app/server"
	"baseGo/src/red_wap/conf"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserService struct{}

var (
	UserSessionService = new(middleware.UserSessionService)
)

// 登录im
func (UserService) LoginIm(lineId, agencyId string, sid, platform, server string, mid int64, userId int) (*structs.ConnectReply, error) {
	type params struct {
		Mid      int64   `json:"mid"`
		Key      string  `json:"key"`
		RoomID   string  `json:"roomId"`
		Platform string  `json:"platform"`
		Accepts  []int32 `json:"accepts"`
	}

	token := params{
		Mid:      mid,
		Key:      lineId + "-" + agencyId + "-" + fmt.Sprint(userId),
		RoomID:   "",
		Platform: platform,
		Accepts:  nil,
	}
	tokenByte, err := json.Marshal(token)
	if err != nil {
		golog.Error("AccountService", "LoginIm", "err:", err)
		return nil, err
	}

	type ConnectReq struct {
		UserId   int     `json:"userId"`
		Server   string  `json:"server"`
		LineId   string  `json:"lineId"`
		AgencyId string  `json:"agencyId"`
		Cookie   string  `json:"cookie"`
		Token    []byte  `json:"token"`
		Rooms    []int64 `json:"rooms"`
	}

	req := &ConnectReq{
		Server:   server,
		LineId:   lineId,
		AgencyId: agencyId,
		Cookie:   "",
		Token:    tokenByte,
		UserId:   userId,
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		golog.Error("AccountService", "LoginIm", "err:", err)
		return nil, err
	}
	if "" == registry_module.GetLogicHttpUrl() {
		return nil, &validate.Err{Code: code.IM_CONN_ERROR}
	}
	host := fmt.Sprintf("http://%v", registry_module.GetLogicHttpUrl())
	golog.Info("User", "loginIm", "logic地址：", host)
	client := &http.Client{} //客户端
	request, err := http.NewRequest("POST", host+"/goim/cachesession", ioutil.NopCloser(bytes.NewReader(reqByte)))
	if err != nil {
		golog.Error("AccountService", "LoginIm", "err:", err)
		return nil, err
	}
	//给一个key设定为响应的value.
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request) //发送请求
	if err != nil {
		golog.Error("User", "loginIm", "请求logic返回失败:", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		golog.Error("AccountService", "LoginIm", "err:", err)
		return nil, err
	}
	reply := new(structs.ConnectReply)
	if err := json.Unmarshal(respBody, reply); err != nil {
		golog.Error("AccountService", "LoginIm", "err:", err)
		return nil, err
	}
	return reply, nil
}

// 获取用户信息
func (ms UserService) GetUserInfo(lineId, agencyId, account string) (*structs.UserResp, error) {
	// 获取线路信息获取当前是钱包模式还是额度转换模式
	lineInfo, err := server.GetStytemLineInfo(lineId)
	if err != nil {
		golog.Error("UserService", "GetUserInfo", "error:", err)
		return nil, &validate.Err{Code: code.LINE_QUERY_FAILED}
	}

	sess := conf.GetXormSession()
	defer sess.Close()

	res := new(structs.UserResp)
	has, result := UserBo.GetOneByAccount(sess, lineId, agencyId, account)
	if !has {
		return nil, &validate.Err{Code: code.MEMBER_INFORMATION_QUERY_FAILED}
	}
	res = &structs.UserResp{
		Id:               result.Id,                                         // 主键id
		LineId:           result.LineId,                                     // 线路ID
		AgencyId:         result.AgencyId,                                   // 代理id
		Account:          result.Account,                                    // 账号
		Balance:          result.Balance,                                    // 会员余额
		CreateTime:       result.CreateTime,                                 // 创建时间
		EditTime:         result.EditTime,                                   // 修改时间
		Capital:          result.Capital,                                    // 红包押金
		AvailableBalance: common.DecimalSub(result.Balance, result.Capital), // 可用金额
		LastLoginIp:      result.LastLoginIp,                                // 上次登陆ip
		LastLoginTime:    result.LastLoginTime,                              // 上次登陆时间
	}
	if lineInfo.TransType == model.TRANS_TYPE_WALLET { // 钱包模式
		// 钱包模式的话需要获取线路那边余额
		respData, err := server.Wallet(lineInfo.ApiUrl, "getbalance", "account", conf.GetConfig().Listening.Md5key, conf.GetConfig().Listening.Deskey, &conf.ReqMember{
			Username: account,
			Currency: "CNY",
		}, "", 0)
		if err != nil {
			golog.Error("UserService", "GetUserInfo", "error:", err)
			return nil, &validate.Err{Code: code.MEMBER_INFORMATION_QUERY_FAILED}
		}
		if respData.Code != 1000 {
			sess.Rollback()
			return nil, &validate.Err{Code: code.MEMBER_INFORMATION_QUERY_FAILED}
		}
		res.Balance = respData.Member.Balance
		res.AvailableBalance = common.DecimalSub(result.Balance, result.Capital)
	}
	return res, nil
}
