package server

import (
	"baseGo/src/fecho/golog"
	registry_module "baseGo/src/fecho/registry/registry-module"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_robot/app/middleware/validate"
	"baseGo/src/red_robot/conf"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LogicSession struct {
	Server string  `json:"server"`
	RoomId int     `json:"roomId"`
	Rooms  []int64 `json:"rooms"`
	Online bool    `json:"online"`
}

type MsgReq struct {
	Data Data `json:"data"`
}

type Data struct {
	Code int              `json:"code"`
	Info *structs.MsgResp `json:"info"`
}

// 查询房间在线会员mids
func SendRoomMessageFunc(url string, data interface{}) (int, error) {
	reqByte, err := json.Marshal(data)
	if err != nil {
		golog.Error("MemberMsgService", "sendRoomMessageFunc", "err:", err)
		return 0, err
	}
	if registry_module.GetLogicHttpUrl() == "" {
		return 0, nil
	}
	host := fmt.Sprintf("http://%v", registry_module.GetLogicHttpUrl())
	client := &http.Client{} //客户端
	request, err := http.NewRequest("POST", host+"/goim"+url, ioutil.NopCloser(bytes.NewReader(reqByte)))
	if err != nil {
		golog.Error("MemberMsgService", "sendRoomMessageFunc", "err:", err)
		return 0, err
	}
	//给一个key设定为响应的value.
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request) //发送请求
	if err != nil {
		golog.Error("MemberMsgService", "sendRoomMessageFunc", "", err)
		return 0, &validate.Err{Code: code.MESSAGE_FAILED_TO_BE_SENT}
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		golog.Error("MemberMsgService", "sendRoomMessageFunc", "err:", err)
		return 0, &validate.Err{Code: code.MESSAGE_FAILED_TO_BE_SENT}
	}
	msgReq := new(MsgReq)
	err = json.Unmarshal(result, &msgReq)
	if err != nil || msgReq.Data.Code != 200 {
		if err != nil {
			golog.Error("MemberMsgService", "sendRoomMessageFunc", "err:", err)
		}
		return 0, &validate.Err{Code: code.JSON_UNMARSHAL_ERROR}
	}
	if nil == msgReq.Data.Info {
		return 0, &validate.Err{Code: code.MESSAGE_FAILED_TO_BE_SENT}
	}

	return msgReq.Data.Info.MsgId, nil
}

func SessionByKey(key string) (session *LogicSession, err error) {
	conn := conf.GetRedis()
	redisClient := conn.Get()
	back, err := redisClient.Do("Get", key)
	if err != nil {
		return
	}
	session = new(LogicSession)
	if err = json.Unmarshal(back.([]byte), session); err != nil {
		return
	}
	return
}
