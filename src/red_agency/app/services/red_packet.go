package services

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm/help"
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/conf"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type RedPacketService struct{}

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 生成红包
// redPacketAllocation 分配红包份额
func redPacketAllocation(num int, money, minMoney float64) ([]structs.OrderRecord, error) {
	// 红包列表
	redList := make([]structs.OrderRecord, 0)
	// 去除最低金额后的总金额
	moneyInt := int(money*100) - (int(minMoney*100) * num)
	if moneyInt < 0 {
		return redList, &validate.Err{Code: code.RED_ENVELOPE_INSUFFICIENT_AMOUNT}
	}

	for i := num; i > 0; i-- {
		info := structs.OrderRecord{}
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

// 创建普通红包 CreateOrdinaryRedPacket
func (RedPacketService) CreateOrdinaryRedPacket(redAmount float64, redNum, roomId, isAuto, autoTime, gameTime int) error {
	// 开启自动之后写入定时任务  到时间之后自动发红包
	// 未开启 即时发红包 所以红包自动发送时间不写入红包记录中
	// 前面的处理是一样的 插入红包 需要根据红包自动开启时间来
	// 获取线路信息获取当前是钱包模式还是额度转换模式
	//lineInfo, err := server.GetStytemLineInfo(lineId)
	//if err != nil {
	//	golog.Error("UserService", "GetUserInfo", "error:", err)
	//	return &validate.Err{Code: code.LINE_QUERY_FAILED}
	//}

	sess := conf.GetXormSession()
	defer sess.Close()
	// 查询当前房间信息 // 发送普通红包不受聊天室红包限制
	has, roomInfo := RoomBo.GetOne(sess, roomId)
	if !has {
		return &validate.Err{Code: code.GAME_GROUP_QUERY_FAILED}
	}

	// 查询会员余额
	has, userInfo := UserBo.GetById(sess, roomInfo.LineId, roomInfo.AgencyId, roomInfo.RobotId)
	if !has {
		return &validate.Err{Code: code.MEMBER_INFORMATION_QUERY_FAILED}
	}

	// 组装数据
	//红包数据组装
	redInfo := new(structs.RedPacket)
	redInfo.RedEnvelopeAmount = redAmount
	redInfo.RedEnvelopesNum = redNum
	redInfo.AgencyId = roomInfo.AgencyId
	redInfo.LineId = roomInfo.LineId
	redInfo.UserId = roomInfo.RobotId // 群主ID
	redInfo.Account = userInfo.Account
	redInfo.CreateTime = utility.GetNowTimestamp()
	redInfo.RedType = model.ORDINARY_RED_ENVELOPE
	redInfo.RedPlay = 0
	redInfo.RoomId = roomId
	redInfo.RoomName = roomInfo.RoomName
	redInfo.Status = model.RED_STATUS_NOT_START
	redInfo.IsAuto = isAuto
	redInfo.AutoTime = autoTime
	redInfo.EndTime = redInfo.AutoTime + (gameTime * 60 * 60)

	// 插入红包数据
	sess.Begin()

	// 根据是否自动进入创建红包流程
	if isAuto == model.IS_AUTO_YES && autoTime > utility.GetNowTimestamp() {
		count, err := RedPacketBo.InsertRedPacket(sess, redInfo)

		if err != nil || count == 0 {
			sess.Rollback()
			if err != nil {
				golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
			}
			return &validate.Err{Code: code.RED_ENVELOPE_ADDITION_FAILED}
		}
		sess.Commit()
		redbyte, _ := json.Marshal(redInfo)
		// 写入时间轮
		keyData := map[string]string{
			"key":      "OrdinaryRedPacket",
			"gameTime": strconv.Itoa(gameTime),
			"autoTime": strconv.Itoa(autoTime),
			"redInfo":  string(redbyte),
		}
		b, _ := json.Marshal(keyData)
		// 创建红包格子
		server.AddTimeWheel(string(b), (autoTime-utility.GetNowTimestamp())%600)
	} else {
		redInfo.Status = model.RED_STATUS_NORMAL
		if autoTime == 0 {
			redInfo.AutoTime = utility.GetNowTimestamp()
		}
		count, err := RedPacketBo.InsertRedPacket(sess, redInfo)
		if err != nil || count == 0 {
			sess.Rollback()
			if err != nil {
				golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
			}
			return &validate.Err{Code: code.RED_ENVELOPE_ADDITION_FAILED}
		}
		// 立即发送
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
		fmt.Printf("forPrint logKey: %T, redLogs: %t", logKey, redLogs)
		// redis.hmset("key", {"id": 5, somekey: "someval"}); //will work
		// redis.hmset(1, {"id": 5, somekey: "someval"}); //will not work
		// _, err = conf.GetRedis().Get().Do("HMSet", logKey, redLogs)
		// if err != nil {
		// 	// 红包存储失败
		// 	sess.Rollback()
		// 	golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
		// 	return err
		// }
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
		count, err = MemberCashRecordBo.Inster(sess, recordInfo)
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
		msgData["redStatus"] = 1
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
		err = MessageHistoryBo.SaveMessageHistory(sess, msgHis)
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
		server.AddTimeWheel(string(b), gameTime*60*60)
		// 转发im发送信息
		fmt.Println(data, "data")
		// server.SendRoomMessageFunc("/push/room", data)
	}
	return nil
}

// 获取普通红包列表
func (RedPacketService) GetOrdinaryRedList(lineId string, agencyId string, startTime, endTime, status int, roomName string, pageParams *help.PageParams) (int, []structs.OrdinaryRedPacketResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 查询红包记录 获取发的红包
	count, data, err := RedPacketBo.GetOrdinaryRedList(sess, lineId, agencyId, startTime, endTime, status, roomName, pageParams)
	if err != nil {
		golog.Error("RedPacketService", "GetOrdinaryRedList", "err:", err)
		return 0, nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	result := make([]structs.OrdinaryRedPacketResp, 0)
	if len(data) > 0 {
		redId := make([]int, 0)
		for _, v := range data {
			redId = append(redId, int(v.Id))
		}
		// 查询红包记录 获取红包的领取金额和剩余金额
		logCount, err := RedPacketLogBo.GetOrderReceiveCountResp(sess, lineId, agencyId, redId)
		if err != nil {
			golog.Error("RedPacketService", "GetOrdinaryRedList", "err:", err)
			return 0, nil, &validate.Err{Code: code.QUERY_FAILED}
		}
		redCount := make(map[int]int, 0)
		redMoney := make(map[int]float64, 0)
		for _, v := range logCount {
			redCount[v.RedId] = v.Count
			redMoney[v.RedId] = v.ReceiveMoney
		}
		for _, v := range data {
			info := structs.OrdinaryRedPacketResp{
				Id:                v.Id,                // ID
				RedEnvelopeAmount: v.RedEnvelopeAmount, // 红包金额
				RedEnvelopesNum:   v.RedEnvelopesNum,   // 红包数量
				UserId:            v.UserId,            // 会员ID
				Account:           v.Account,           // 会员帐号
				CreateTime:        v.CreateTime,        // 创建时间
				RedType:           v.RedType,           // 游戏类型 1牛牛 2扫雷 0 普通红包
				RoomId:            v.RoomId,            // 群ID
				RoomName:          v.RoomName,          // 群名称
				Status:            v.Status,            // 红包状态 1进行中 2已结束
				ReturnMoney:       v.ReturnMoney,       // 返还金额
				IsAuto:            v.IsAuto,            // 是否开启自动 1是 0否
				AutoTime:          v.AutoTime,          // 自动开始时间
				EndTime:           v.EndTime,           // 红包结束时间
			}
			if _, ok := redCount[int(v.Id)]; ok {
				info.ReceiveNum = redCount[int(v.Id)]
			}
			if _, ok := redMoney[int(v.Id)]; ok {
				info.ReceiveMoney = redMoney[int(v.Id)]
			}
			if info.EndTime < utility.GetNowTimestamp() {
				info.Status = 2
			}
			info.ReturnMoney = common.DecimalSub(info.RedEnvelopeAmount, info.ReceiveMoney)
			info.ReturnNum = info.RedEnvelopesNum - info.ReceiveNum
			result = append(result, info)
		}
	}
	return count, result, err
}

// 查看普通红包领取详情
func (RedPacketService) GetOrdinaryRedInfo(lineId string, agencyId string, redId int) ([]*structs.RedOrderResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 查询红包记录 获取发的红包
	has, _ := RedPacketBo.GetRedInfo(sess, lineId, agencyId, 0, redId, 0)
	if !has {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}

	// 查询红包记录 获取红包的领取金额和剩余金额
	logList, err := RedPacketLogBo.GetRedInfo(sess, lineId, agencyId, redId)
	if err != nil {
		golog.Error("RedPacketService", "GetOrdinaryRedList", "err:", err)
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return logList, err
}

// 修改普通红包
func (RedPacketService) EditOrdinaryRedInfo(lineId string, agencyId string, redId int, redAmount float64, redNum, roomId, isAuto, autoTime, gameTime int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 查询红包记录 获取发的红包
	has, redInfo := RedPacketBo.GetRedInfo(sess, lineId, agencyId, 0, redId, roomId)
	if !has {
		return &validate.Err{Code: code.RED_PACKET_QUERY_FAILED}
	}
	if redInfo.Status != model.RED_STATUS_NOT_START {
		return &validate.Err{Code: code.RED_PACKRT_STATUS_IS_INCORRECT}
	}
	redInfo.RedEnvelopeAmount = redAmount
	redInfo.RedEnvelopesNum = redNum
	redInfo.IsAuto = isAuto
	redInfo.AutoTime = autoTime
	redInfo.EndTime = autoTime + (gameTime * 60 * 60)

	count, err := RedPacketBo.UpdateRed(sess, redInfo, "red_envelope_amount, red_envelopes_num, is_auto, auto_time, end_time")
	if err != nil && count > 0 {
		if err != nil {
			golog.Error("RedPacketService", "EditOrdinaryRedInfo", "err:", err)
		}
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 根据是否自动进入创建红包流程
	if isAuto == model.IS_AUTO_YES && autoTime < utility.GetNowTimestamp() {
		redbyte, _ := json.Marshal(redInfo)
		// 写入时间轮
		keyData := map[string]string{
			"key":      "OrdinaryRedPacket",
			"gameTime": strconv.Itoa(gameTime),
			"autoTime": strconv.Itoa(autoTime),
			"redInfo":  string(redbyte),
		}
		b, _ := json.Marshal(keyData)
		// 创建红包格子
		server.AddTimeWheel(string(b), (autoTime-utility.GetNowTimestamp())%600)
	} else {
		// 立即发送
		err = server.CreateOrdinaryRedPacket(redInfo, gameTime)
		if err != nil {
			return err
		}
	}
	return err
}

// 删除普通红包
func (RedPacketService) DelOrdinaryRedInfo(lineId string, agencyId string, redId int) error {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 查询红包记录 获取发的红包
	has, redInfo := RedPacketBo.GetRedInfo(sess, lineId, agencyId, 0, redId, 0)
	if !has {
		return &validate.Err{Code: code.RED_PACKET_QUERY_FAILED}
	}
	if redInfo.Status != model.RED_STATUS_NOT_START {
		return &validate.Err{Code: code.RED_PACKRT_STATUS_IS_INCORRECT}
	}
	redInfo.DeleteTime = utility.GetNowTimestamp()

	count, err := RedPacketBo.UpdateRed(sess, redInfo, "delete_time")
	if err != nil && count > 0 {
		if err != nil {
			golog.Error("RedPacketService", "EditOrdinaryRedInfo", "err:", err)
		}
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return err
}

func (RedPacketService) Orders(lineId string, agencyId string, redId int) ([]structs.OrderCollectResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	//插入注单采集
	rps, err := RedPacketLogBo.GetRedOrderRecordByRedIdWithStatus12(sess, lineId, agencyId, redId)
	if err != nil {
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		return nil, &validate.Err{Code: code.WRITING_CASH_RECORD_FAILED}
	}
	return rps, nil
}
