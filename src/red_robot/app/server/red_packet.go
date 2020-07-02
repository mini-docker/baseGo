package server

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_robot/app/middleware/validate"
	"baseGo/src/red_robot/conf"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// 发红包
func CreateOrdinaryRedPacket(redInfo *structs.RedPacket, gameTime int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	sess.Begin()
	// 插入红包数据
	// 查询红包信息
	RedPacketBo := new(bo.RedPacket)
	has, redInfo := RedPacketBo.ByRoomIdGetRedInfo(sess, int(redInfo.Id))
	if !has {
		return &validate.Err{Code: code.RED_PACKET_QUERY_FAILED}
	} else if redInfo.DeleteTime != 0 {
		sess.Rollback()
		return nil
	}
	if redInfo.Status != model.RED_STATUS_NOT_START {
		// 状态不为未开始 不进行下序操作
		sess.Rollback()
		return nil
	}
	if redInfo.AutoTime > utility.GetNowTimestamp() { // 红包未开始
		sess.Rollback()
		return &validate.Err{Code: code.RED_PACKRT_SEND_TIME_NOT_UP}
	}
	if redInfo.EndTime <= utility.GetNowTimestamp() { // 红包已结束
		RedPacketBo.UpdateRedStatus(sess, int(redInfo.Id), model.RED_STATUS_OVER)
	}
	// 修改红包状态为进行中
	n, err := RedPacketBo.UpdateRedStatus(sess, int(redInfo.Id), model.RED_STATUS_NORMAL)
	if err != nil || n <= 0 {
		sess.Rollback()
		if err != nil {
			golog.Error("RedPlay", "invalidBetSlip", "err:", err)
		}
		return err
	}
	// 生成红包
	redLogList, err := redPacketAllocation(redInfo.RedEnvelopesNum, redInfo.RedEnvelopeAmount, 0.01)
	if err != nil {
		sess.Rollback()
		golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
		return err
	}
	// 将剩余的红包存入redis中
	logKey := fmt.Sprintf("%d_%d_redLog", redInfo.Id, redInfo.RoomId)
	redLogs := make(map[string]interface{}, 0)
	for k, v := range redLogList {
		log, _ := json.Marshal(v)
		redLogs[fmt.Sprint(k)] = string(log)
	}
	_, err = conf.GetRedis().Get().Do("HMSet", logKey, redLogs)
	if err != nil {
		// 红包存储失败
		sess.Rollback()
		golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
		return err
	}
	// 写入现金记录
	// 写现金记录
	recordInfo := new(structs.MemberCashRecord)
	recordInfo.LineId = redInfo.LineId
	recordInfo.AgencyId = redInfo.AgencyId
	recordInfo.GameType = redInfo.RedType
	recordInfo.GameName = model.RED_ENVELOPE_TYPE[redInfo.RedType]
	recordInfo.FlowType = model.MEMBER_RECORD_ORDINARY_CREATE
	recordInfo.Money = redInfo.RedEnvelopeAmount
	recordInfo.Remark = fmt.Sprintf("发普通红包,红包金额%v,扣除会员余额%v", recordInfo.Money, recordInfo.Money)
	recordInfo.CreateTime = utility.GetNowTimestamp()
	recordInfo.UserId = redInfo.UserId
	recordInfo.Account = redInfo.Account
	count, err := new(bo.MemberCashRecord).Inster(sess, recordInfo)
	if err != nil || count <= 0 {
		sess.Rollback()
		if err != nil {
			golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
		}
		return &validate.Err{Code: code.WRITING_CASH_RECORD_FAILED}
	}

	// 给会员推送红包信息
	//推送信息
	type RoomReq struct {
		Operation   int    `json:"operation"`
		RoomId      int    `json:"roomId"`
		LineId      string `json:"lineId"`
		AgencyId    string `json:"agencyId"`
		Msg         string `json:"msg"`
		MsgType     int    `json:"msgType"` // 1 发红包 2 红包结算
		SendId      int    `json:"sendId"`
		Key         string `json:"key"`
		ReceiveType int    `json:"receiveType"` // 1 私聊 2房间
	}
	msgData := make(map[string]interface{})
	msgData["redSender"] = redInfo.Account
	msgData["redMoney"] = redInfo.RedEnvelopeAmount
	msgData["redNum"] = redInfo.RedEnvelopesNum
	msgData["redId"] = int(redInfo.Id)
	msgData["gameType"] = redInfo.RedType
	msgData["gameTypeName"] = model.RED_ENVELOPE_TYPE[redInfo.RedType]
	msgData["gamePlay"] = redInfo.RedPlay
	msgData["gamePlayName"] = model.RED_ENVELOPE_TYPE_PLAY[redInfo.RedType][redInfo.RedPlay]
	msgData["gameTime"] = gameTime
	msgData["createTime"] = redInfo.CreateTime
	msgDataByte, _ := json.Marshal(msgData)
	data := &RoomReq{
		Operation:   4,
		RoomId:      redInfo.RoomId,
		LineId:      redInfo.LineId,
		AgencyId:    redInfo.AgencyId,
		Msg:         string(msgDataByte),
		MsgType:     1,
		SendId:      redInfo.UserId,
		Key:         "",
		ReceiveType: 2,
	}
	msgHis := &structs.MessageHistory{
		LineId:     redInfo.LineId,
		AgencyId:   redInfo.AgencyId,
		MsgType:    1,
		MsgContent: string(msgDataByte),
		SenderId:   redInfo.UserId,
		SenderName: redInfo.Account,
		Status:     1,
		SendTime:   utility.GetNowTimestamp(),
		RoomId:     redInfo.RoomId,
	}
	err = new(bo.MessageHistory).SaveMessageHistory(sess, msgHis)
	if err != nil {
		return err
	}
	sess.Commit()

	// 写入时间轮
	keyData := map[string]string{
		"key":            "redPacketSettle",
		"redId":          strconv.Itoa(int(redInfo.Id)),
		"roomId":         strconv.Itoa(redInfo.RoomId),
		"settlementTime": strconv.Itoa(utility.GetNowTimestamp() + gameTime*60*60),
		"lineId":         redInfo.LineId,
		"agencyId":       redInfo.AgencyId,
	}
	b, _ := json.Marshal(keyData)
	AddTimeWheel(string(b), 5)

	// 转发im发送信息
	SendRoomMessageFunc("/push/room", data)
	return nil
}

// 生成红包
// redPacketAllocation 分配红包份额
func redPacketAllocation(num int, money, minMoney float64) ([]*structs.OrderRecord, error) {
	// 红包列表
	redList := make([]*structs.OrderRecord, 0)
	// 去除最低金额后的总金额
	moneyInt := int(money*100) - (int(minMoney*100) * num)
	if moneyInt < 0 {
		return redList, &validate.Err{Code: code.RED_ENVELOPE_INSUFFICIENT_AMOUNT}
	}

	for i := num; i > 0; i-- {
		info := new(structs.OrderRecord)
		if moneyInt == 0 {
			info.ReceiveMoney = minMoney
		} else {
			if i > 1 {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				randNum := moneyInt / i * 2
				if randNum > 0 {
					redMoney := r.Intn(randNum)
					moneyInt = moneyInt - redMoney
					info.ReceiveMoney = common.DecimalSum(common.DecimalDivs(float64(redMoney), 100), minMoney)
				} else {
					info.ReceiveMoney = utility.Round(minMoney, 2)
				}
			} else {
				info.ReceiveMoney = common.DecimalSum(common.DecimalDivs(float64(moneyInt), 100), minMoney)
			}
		}
		redList = append(redList, info)
	}
	return redList, nil
}
