package server

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/echo"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm"
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

// 红包玩法计算
type RedPlay struct{}

//
type RedSettlement struct {
	RedId   int // 红包ID
	roomId  int // 房间ID
	RedType int // 红包类型 1 牛牛红包 2
	RedPlay int // 红包玩法 1
	Data    []UserRed
}

// 红包结果结构体
type UserRed struct {
	UserId    int     // 会员ID
	Identity  int     //	会员身份 1庄家 2闲家
	RedAmount float64 //	红包金额
	Win       float64 // 结果
}

// 根据房间ID获取房间信息
var (
	roomBo             = new(bo.Room)
	userBo             = new(bo.User)
	memberCashRecordBo = new(bo.MemberCashRecord)
	redPacketBo        = new(bo.RedPacket)
	redPacketLogBo     = new(bo.RedPacketLog)
)

// 红包金额计算
func (RedPlay) RedEnvelopeAmountCalculation(lineId, agencyId string, redID, roomId int) (*RedSettlement, error) {
	redisClient := conf.GetRedis().Get()
	// 检查结算红包是否加锁
	res, err := redisClient.Do("Exists", fmt.Sprintf("settleBag_%v", redID))
	fResult, _ := strconv.Atoi(string(res.([]byte)))
	if fResult == 1 {
		return nil, nil
	}
	// 加锁
	redisClient.Do("Set", fmt.Sprintf("settleBag_%v", redID), 1, time.Second*10)
	// 获取线路信息获取当前是钱包模式还是额度转换模式
	lineInfo, err := GetStytemLineInfo(lineId)
	if err != nil {
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "error:", err)
		return nil, &validate.Err{Code: code.LINE_QUERY_FAILED}
	}
	// 根据房间ID和红包ID 获取红包类型和红包玩法 和会员下注金额等数据
	sess := conf.GetXormSession()

	has, room := roomBo.GetOne(sess, roomId)
	if !has {
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "error:游戏群查询失败", nil, roomId)
		sess.Close()
		return nil, &validate.Err{Code: code.GAME_GROUP_QUERY_FAILED}
	}
	// 查询红包信息
	has, redInfo := redPacketBo.ByRoomIdGetRedInfo(sess, redID)
	if !has && redInfo.DeleteTime != 0 {
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "error:红包信息查询失败", nil)
		sess.Close()
		return nil, &validate.Err{Code: code.RED_PACKET_QUERY_FAILED}
	}
	if redInfo.Status == model.RED_STATUS_INVALID || redInfo.Status == model.RED_STATUS_OVER {
		golog.Info("RedPlay", "RedEnvelopeAmountCalculation", "红包信息已结算", nil)
		sess.Close()
		return nil, nil
	}
	has, userInfo := UserBo.GetOne(sess, redInfo.UserId)
	if !has {
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "用户信息不存在:", err, "userId:", redInfo.UserId, "account:", userInfo.Account)
		sess.Close()
		return nil, err
	}

	// 获取红包领取记录（下注记录）
	RedPacketLogBo := new(bo.RedPacketLog)
	logList, err := RedPacketLogBo.GetRedLog(sess, lineId, agencyId, redID)
	if err != nil {
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "获取下注记录失败", nil)
		sess.Close()
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	if len(logList) <= 1 {
		// 抢红包人数太少 流局
		// 将红包状态改为无效
		sess.Begin()
		err = invalidBetSlip(sess, redID)
		if err != nil {
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "修改注单无效失败:", err)
			sess.Close()
			return nil, err
		}
		if len(logList) == 1 {
			recordInfo := new(structs.MemberCashRecord)
			switch redInfo.RedType {
			case model.NIUNIU_RED_ENVELOPE:
				// 返还红包押金
				remark := fmt.Sprintf("注单号:%v,牛牛红包流局返还保证金%v元", logList[0].OrderNo, redInfo.Capital)
				if userInfo.IsRobot != model.USER_IS_ROBOT_YES {
					err = new(UserServer).ChangeAmount(sess, lineId, agencyId, redInfo.Account, lineInfo.ApiUrl, lineInfo.TransType, redInfo.UserId, redInfo.Capital, remark, common.DecimalSub(0, redInfo.Capital))
					if err != nil {
						sess.Rollback()
						sess.Close()
						golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "流局牛牛会员金额变动失败:", nil, err.Error(), logList[0].LineId, logList[0].AgencyId)
						return nil, err
					}
				}

				recordInfo.LineId = redInfo.LineId
				recordInfo.AgencyId = redInfo.AgencyId
				recordInfo.GameType = redInfo.RedType
				recordInfo.GameName = model.RED_ENVELOPE_TYPE[redInfo.RedType]
				recordInfo.Money = redInfo.RedEnvelopeAmount
				recordInfo.FlowType = model.MEMBER_RECORD_RETURBN
				recordInfo.Remark = fmt.Sprintf("返还红包保证金,返还%v,增加会员余额%v", redInfo.Capital, redInfo.Capital)
				recordInfo.CreateTime = utility.GetNowTimestamp()
			case model.MINESWEEPER_RED_PACKET:
				// 返还扣除的红包金额
				// 因为是流局 所以返还金额=红包总金额
				// 增加会员余额
				remark := fmt.Sprintf("注单号:%v,扫雷红包流局返还发包金额%v元", logList[0].OrderNo, redInfo.RedEnvelopeAmount)
				if userInfo.IsRobot != model.USER_IS_ROBOT_YES {
					err = new(UserServer).ChangeAmount(sess, lineId, agencyId, redInfo.Account, lineInfo.ApiUrl, lineInfo.TransType, redInfo.UserId, redInfo.RedEnvelopeAmount, remark, 0)
					if err != nil {
						sess.Rollback()
						sess.Close()
						golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "流局扫雷会员金额变动失败:", err, userInfo.IsRobot, userInfo.Account)
						return nil, err
					}
				}

				recordInfo.LineId = redInfo.LineId
				recordInfo.AgencyId = redInfo.AgencyId
				recordInfo.GameType = redInfo.RedType
				recordInfo.GameName = model.RED_ENVELOPE_TYPE[redInfo.RedType]
				recordInfo.Money = redInfo.RedEnvelopeAmount
				recordInfo.FlowType = model.MEMBER_RECORD_RETURBN
				recordInfo.Remark = fmt.Sprintf("返还未领取完的红包金额%v元,增加会员余额%v", redInfo.RedEnvelopeAmount, redInfo.RedEnvelopeAmount)
				recordInfo.CreateTime = utility.GetNowTimestamp()
			}
			// 写入现金记录
			count, err := memberCashRecordBo.Inster(sess, recordInfo)
			if err != nil || count <= 0 {
				sess.Rollback()
				if err != nil {
					golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "写入现金记录失败:", err)
				}
				sess.Close()
				return nil, &validate.Err{Code: code.WRITING_CASH_RECORD_FAILED}
			}

			sess.Commit()
			var memberWinNum, adminWinNum, adminNum int
			sendData := make([]structs.RedPacketLogInfoResp, 0)
			for _, v := range logList {
				vData := make(map[string]string)
				json.Unmarshal([]byte(v.Extra), &vData)
				var memberNum, thunderNum int
				if _, ok := vData["adminNum"]; ok {
					adminNum, _ = strconv.Atoi(vData["adminNum"])
				}
				if _, ok := vData["memberNum"]; ok {
					memberNum, _ = strconv.Atoi(vData["memberNum"])
				}
				if _, ok := vData["thunderNum"]; ok {
					thunderNum, _ = strconv.Atoi(vData["thunderNum"])
				}
				if v.RedId == int(redInfo.Id) && v.Account != redInfo.Account {
					if v.GameTime == model.NIUNIU_RED_ENVELOPE {
						info := structs.RedPacketLogInfoResp{
							Account:   v.Account,
							RealMoney: v.RealMoney,
							AdminNum:  adminNum,
							MemberNum: memberNum,
						}
						sendData = append(sendData, info)
					} else if v.GameTime == model.MINESWEEPER_RED_PACKET {
						if v.Status == model.RED_RESULT_LOSE {
							info := structs.RedPacketLogInfoResp{
								Account:    v.Account,
								RealMoney:  v.RealMoney,
								MemberMine: int(common.DecimalMul(v.ReceiveMoney, 100)) % 10,
								ThunderNum: thunderNum,
							}
							sendData = append(sendData, info)
						}
					}
					if v.Status == model.RED_RESULT_WIN {
						memberWinNum = memberWinNum + 1
					} else if v.Status == model.RED_RESULT_LOSE {
						adminWinNum = adminWinNum + 1

					}
				} else if v.RedId == int(redInfo.Id) && v.Account == redInfo.Account {
					adminNum = memberNum
				}
				// 将红包领取记录状态改为无效
				_, err = RedPacketLogBo.UpdateRedLogStatus(sess, logList[0].Id, model.RED_STATUS_INVALID)
				if err != nil {
					sess.Rollback()
					sess.Close()
					if err != nil {
						golog.Error("RedPlay", "invalidBetSlip", "err:", err)
					}
					return nil, err
				}
			}
			// 解锁
			redisClient.Do("Del", fmt.Sprintf("settleBag_%v", redID))
			sess.Close()
			return nil, err
		} else {
			if redInfo.RedType != model.ORDINARY_RED_ENVELOPE {
				golog.Error("RedPlay", "invalidBetSlip", "没有查到注单数据", nil)
				sess.Close()
				return nil, err
			}
		}

	}
	var robotOrders, memberOrders int // 机器人领包个数，会员领包个数
	var robotCash, memberCash float64 // 机器人盈利，会员盈利
	// 开启事务
	sess.Begin()
	// 根据游戏类型和玩法判断计
	cashList := make([]*structs.MemberCashRecord, 0)
	switch redInfo.RedType {
	case model.NIUNIU_RED_ENVELOPE:
		cashList, robotOrders, robotCash, memberOrders, memberCash, err = NiuNiuSettlement(sess, lineInfo, redInfo, logList, userInfo)
		if err != nil || len(cashList) == 0 {
			sess.Rollback()
			sess.Close()
			return nil, err
		}
	case model.MINESWEEPER_RED_PACKET:
		cashList, robotOrders, robotCash, memberOrders, memberCash, err = MineSettlement(sess, lineInfo, redInfo, logList, room, userInfo)
		if err != nil || len(cashList) == 0 {
			sess.Rollback()
			sess.Close()
			return nil, err
		}
	case model.ORDINARY_RED_ENVELOPE:
		cashList, err = OrdinarySettlement(sess, lineInfo, redInfo, logList)
		if err != nil || len(cashList) == 0 {
			sess.Rollback()
			sess.Close()
			return nil, err
		}
	}

	// 插入现金流水
	_, err = new(bo.MemberCashRecord).Inster(sess, cashList...)
	if err != nil {
		sess.Rollback()
		sess.Close()
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		return nil, &validate.Err{Code: code.WRITING_CASH_RECORD_FAILED}
	}

	//插入注单采集
	rps, err := RedPacketLogBo.GetRedOrderRecordByRedIdWithStatus12(sess, lineId, agencyId, redID)
	if err != nil {
		sess.Rollback()
		sess.Close()
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		return nil, &validate.Err{Code: code.WRITING_CASH_RECORD_FAILED}
	}
	redPacketCollectBo := new(bo.RedPacketCollect)
	redPCs := make([]*structs.RedPacketCollect, 0)
	tm := utility.GetNowTimestamp()
	for _, rp := range rps {
		if rp.Account == rp.RedSender {
			rp.Remark = fmt.Sprintf("AI领包%d个，盈利%v元，会员领包%d个，盈利%v元", robotOrders, robotCash, memberOrders, memberCash)
		} else {
			if userInfo.IsRobot != model.USER_IS_ROBOT_YES {
				rp.Remark = "领取会员发包"
			} else {
				rp.Remark = "领取AI发包"
			}
		}
		rpStr, err := json.Marshal(rp)
		if err != nil {
			sess.Rollback()
			sess.Close()
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
			return nil, &echo.Err{Code: code.WRITING_CASH_RECORD_FAILED}
		}
		redPc := new(structs.RedPacketCollect)
		redPc.LineId = lineId
		redPc.AgencyId = agencyId
		redPc.CreateTime = tm
		redPc.CollectStatus = 1
		redPc.SettlementInfo = string(rpStr)
		redPCs = append(redPCs, redPc)
		redPCStr, err := json.Marshal(redPc)
		if err != nil {
			sess.Rollback()
			sess.Close()
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
			return nil, &echo.Err{Code: code.WRITING_CASH_RECORD_FAILED}
		}
		err = InsertRedPacketCollect(lineId, string(redPCStr))
		if err != nil {
			sess.Rollback()
			sess.Close()
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
			return nil, &echo.Err{Code: code.WRITING_CASH_RECORD_FAILED}
		}
	}
	_, err = redPacketCollectBo.InsertRedPacketCollects(sess, redPCs...)
	if err != nil {
		sess.Rollback()
		sess.Close()
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		return nil, &echo.Err{Code: code.WRITING_CASH_RECORD_FAILED}
	}
	RedPacketBo := new(bo.RedPacket)
	RedPacketBo.UpdateRedStatus(sess, redID, model.RED_STATUS_OVER)
	// 提交事务
	sess.Commit()

	// 推送结算信息
	var memberWinNum, adminWinNum, adminNum int
	sendData := make([]structs.RedPacketLogInfoResp, 0)
	for _, v := range logList {
		vData := make(map[string]string)
		json.Unmarshal([]byte(v.Extra), &vData)
		var memberNum, thunderNum int
		if _, ok := vData["adminNum"]; ok {
			adminNum, _ = strconv.Atoi(vData["adminNum"])
		}
		if _, ok := vData["memberNum"]; ok {
			memberNum, _ = strconv.Atoi(vData["memberNum"])
		}
		if _, ok := vData["thunderNum"]; ok {
			thunderNum, _ = strconv.Atoi(vData["thunderNum"])
		}
		if v.RedId == int(redInfo.Id) && v.Account != redInfo.Account {
			if v.GameType == model.NIUNIU_RED_ENVELOPE {
				info := structs.RedPacketLogInfoResp{
					Account:   v.Account,
					RealMoney: v.RealMoney,
					AdminNum:  adminNum,
					MemberNum: memberNum,
				}
				sendData = append(sendData, info)
			} else if v.GameType == model.MINESWEEPER_RED_PACKET {
				if v.Status == model.RED_RESULT_LOSE {
					info := structs.RedPacketLogInfoResp{
						Account:    v.Account,
						RealMoney:  v.RealMoney,
						MemberMine: int(common.DecimalMul(v.ReceiveMoney, 100)) % 10,
						ThunderNum: thunderNum,
					}
					sendData = append(sendData, info)
				}
			}
			if v.Status == model.RED_RESULT_WIN {
				memberWinNum = memberWinNum + 1
			} else if v.Status == model.RED_RESULT_LOSE {
				adminWinNum = adminWinNum + 1

			}
		} else if v.RedId == int(redInfo.Id) && v.Account == redInfo.Account {
			adminNum = memberNum
		}
	}
	err = OrderSendMessage(sess, room.LineId, room.AgencyId, redInfo.UserId, roomId, redInfo.Account,
		redInfo.RedEnvelopeAmount,
		redInfo.RedEnvelopesNum, int(redInfo.Id), redInfo.RedType, redInfo.RedPlay, room.GameTime, redInfo.CreateTime, redInfo.Mine, memberWinNum, adminWinNum, adminNum, model.RED_STATUS_OVER, room.Odds, sendData)
	// 解锁
	redisClient.Do("Del", fmt.Sprintf("settleBag_%v", redID))
	sess.Close()
	return nil, err
}

// 注单消息发送
func OrderSendMessage(sess *xorm.Session, lineId, agencyId string, userId, roomId int, account string,
	redAmount float64,
	redNum, redId, gameType, gamePlay, gameTime, createTime, mine, memberWinNum, adminWinNum, adminNum, redStatus int, odds float64, sendData []structs.RedPacketLogInfoResp) error {
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
	msgData["redSender"] = account
	msgData["redMoney"] = redAmount
	msgData["redNum"] = redNum
	msgData["redId"] = redId
	msgData["gameType"] = gameType
	msgData["gameTypeName"] = model.RED_ENVELOPE_TYPE[gameType]
	msgData["gamePlay"] = gamePlay
	msgData["gamePlayName"] = model.RED_ENVELOPE_TYPE_PLAY[gameType][gamePlay]
	msgData["gameTime"] = gameTime
	msgData["createTime"] = createTime
	msgData["mine"] = mine
	msgData["memberWinNum"] = memberWinNum
	msgData["adminWinNum"] = adminWinNum
	msgData["adminNum"] = adminNum
	msgData["sendData"] = sendData
	msgData["odds"] = odds
	msgData["redStatus"] = redStatus
	b, _ := json.Marshal(msgData)
	data := &RoomReq{
		Operation:   4,
		RoomId:      roomId,
		LineId:      lineId,
		AgencyId:    agencyId,
		Msg:         string(b),
		MsgType:     2,
		SendId:      userId,
		Key:         "",
		ReceiveType: 2,
	}
	msgHis := &structs.MessageHistory{
		LineId:     lineId,
		AgencyId:   agencyId,
		MsgType:    2,
		MsgContent: string(b),
		SenderId:   1,
		SenderName: account,
		Status:     1,
		SendTime:   utility.GetNowTimestamp(),
		RoomId:     roomId,
	}
	// 转发im发送信息
	if gameType != model.ORDINARY_RED_ENVELOPE {
		err := new(bo.MessageHistory).SaveMessageHistory(sess, msgHis)
		if err != nil {
			return err
		}
		_, err = SendRoomMessageFunc("/push/room", data)
		return err
	}
	return nil
}

// 注单无效
func invalidBetSlip(sess *xorm.Session, redId int) error {
	// 查询红包信息
	RedPacketBo := new(bo.RedPacket)

	// 获取红包领取记录（下注记录）
	//RedPacketLogBo := new(bo.RedPacketLog)
	n, err := RedPacketBo.UpdateRedStatus(sess, redId, model.RED_STATUS_INVALID)
	if err != nil || n <= 0 {
		sess.Rollback()
		if err != nil {
			golog.Error("RedPlay", "invalidBetSlip", "err:", err)
		}
		return err
	}

	return nil
}

// 牛牛红包算法
func NiuniuAlgorithmt(red *structs.RedPacket, redLog []structs.OrderRecord, userInfo *structs.User) ([]structs.OrderRecord, error) {
	// 牛牛经典算法
	// 经典算法 牛1-6是1倍，牛7-9是2倍，牛牛是3倍
	// 平倍算法 倍数固定为1倍
	userNiu := make(map[string]int, 0)
	var adminNum int
	var adminReceive float64
	for _, v := range redLog {
		userNiu[v.Account] = model.NIUNIU_PLAY_ODDS[red.RedPlay][NiuNiuCalculation(v.ReceiveMoney)]
		if v.Account == red.Account {
			adminNum = NiuNiuCalculation(v.ReceiveMoney)
			adminReceive = v.ReceiveMoney
		}
	}
	if _, ok := userNiu[red.Account]; !ok {
		return nil, &validate.Err{Code: code.DEALER_RESULTS_GET_FAILED}
	}
	var adminMoney float64
	var memberCash float64 // 本局游戏会员盈利
	for k, v := range redLog {
		if v.Account != red.Account {
			if ((NiuNiuCalculation(v.ReceiveMoney) > adminNum || NiuNiuCalculation(v.ReceiveMoney) == 0) && adminNum != 0) ||
				(NiuNiuCalculation(v.ReceiveMoney) == adminNum && v.ReceiveMoney > adminReceive) {
				// 赢
				redLog[k].Status = model.RED_RESULT_WIN
				redLog[k].Money = 0 + common.DecimalMul(red.RedEnvelopeAmount, float64(userNiu[v.Account]))
				if v.IsRobot == model.USER_IS_ROBOT_NO {
					redLog[k].RoyaltyMoney = common.DecimalMul(redLog[k].Money, common.DecimalDiv(redLog[k].Royalty, 100))
				}
				redLog[k].RealMoney = common.DecimalSub(redLog[k].Money, redLog[k].RoyaltyMoney)
				if v.IsRobot != model.USER_IS_ROBOT_YES {
					// 不是机器人,记录有效投注
					redLog[k].ValidBet = redLog[k].Money
				} else {
					// 是机器人，判断发包人是不是会员记录机器人盈利
					if userInfo.IsRobot != model.USER_IS_ROBOT_YES {
						// 发包人是会员，记录机器人盈利金额
						redLog[k].RobotWin = redLog[k].Money
					}
				}
			} else {
				// 输
				if v.IsFreeDeath == model.ROOM_FREE_FROM_DEATH_OFF {
					redLog[k].Status = model.RED_RESULT_INVALID
					redLog[k].Money = 0
					redLog[k].RealMoney = 0
				} else {
					redLog[k].Status = model.RED_RESULT_LOSE
					redLog[k].Money = 0 - common.DecimalMul(red.RedEnvelopeAmount, float64(userNiu[red.Account]))
					redLog[k].RealMoney = redLog[k].Money
					if v.IsRobot != model.USER_IS_ROBOT_YES {
						// 不是机器人,记录有效投注
						redLog[k].ValidBet = common.DecimalSub(0, v.Money)
					} else {
						// 是机器人，判断发包人是不是会员记录机器人盈利
						if userInfo.IsRobot != model.USER_IS_ROBOT_YES {
							// 发包人是会员，记录机器人盈利金额
							redLog[k].RobotWin = redLog[k].Money
						}
					}
				}
			}
			adminMoney = common.DecimalSub(adminMoney, redLog[k].Money)
		}
		vData := make(map[string]string)
		json.Unmarshal([]byte(v.Extra), &vData)
		vData["adminNum"] = strconv.Itoa(adminNum)
		vData["memberNum"] = strconv.Itoa(NiuNiuCalculation(v.ReceiveMoney))
		b, _ := json.Marshal(vData)
		redLog[k].Extra = string(b)
		if v.IsRobot != model.USER_IS_ROBOT_YES {
			// 不是机器人记录会员盈利
			memberCash = common.DecimalSum(memberCash, redLog[k].Money)
		}
	}
	for k, v := range redLog {
		if v.Account == red.Account {
			redLog[k].Money = adminMoney
			if redLog[k].Money != 0 {
				if redLog[k].Money > 0 {
					if v.IsRobot != model.USER_IS_ROBOT_YES {
						// 不是机器人,记录有效投注
						redLog[k].ValidBet = redLog[k].Money
						redLog[k].RoyaltyMoney = common.DecimalMul(redLog[k].Money, common.DecimalDiv(redLog[k].Royalty, 100))
					}
					redLog[k].RealMoney = common.DecimalSub(redLog[k].Money, redLog[k].RoyaltyMoney)
					redLog[k].Status = model.RED_RESULT_WIN
				} else if redLog[k].Money < 0 {
					redLog[k].RealMoney = adminMoney
					redLog[k].Status = model.RED_RESULT_LOSE
					if v.IsRobot != model.USER_IS_ROBOT_YES {
						// 不是机器人,记录有效投注
						redLog[k].ValidBet = common.DecimalSub(0, redLog[k].Money)
					}
				}
				if v.IsRobot == model.USER_IS_ROBOT_YES {
					// 发包人是机器人，记录机器人盈利
					redLog[k].RobotWin = common.DecimalSub(0, memberCash)
				}
			} else {
				redLog[k].RealMoney = 0
				redLog[k].Status = model.RED_STATUS_INVALID
			}
		}
	}
	// 超倍算法 牛几就是几倍  牛牛10倍
	return redLog, nil
}

// 计算尾部3位数字的合的个位数
func NiuNiuCalculation(money float64) int {
	num := int(common.DecimalMul(money, 100))
	b2 := num / 100
	b1 := (num - b2*100) / 10
	b0 := num % 10
	sum := b2 + b1 + b0
	n := sum % 10
	return n
}

// 订单号生成
func OderNo(gameType int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var order string
	if gameType == model.NIUNIU_RED_ENVELOPE {
		order = order + "nnhb"
	} else if gameType == model.MINESWEEPER_RED_PACKET {
		order = order + "slhb"
	}

	order += time.Now().Format("20060102150405") + strconv.Itoa(100000+r.Intn(899999))
	return order
}

// 牛牛结算
func NiuNiuSettlement(sess *xorm.Session, lineInfo *structs.SystemLine, redInfo *structs.RedPacket, logList []structs.OrderRecord, userInfo *structs.User) ([]*structs.MemberCashRecord, int, float64, int, float64, error) {
	redResList, err := NiuniuAlgorithmt(redInfo, logList, userInfo)
	if err != nil {
		sess.Rollback()
		return nil, 0, 0, 0, 0, &validate.Err{Code: code.UPDATE_FAILED}
	}
	var adminMoney, adminRoyaltyMoneym, adminRealMoney, capital float64 // 庄家输赢 庄家抽水 庄家实际输赢 总抽水
	var odds int                                                        // 闲家扣押金时的赔率
	for _, v := range model.NIUNIU_PLAY_ODDS[redInfo.RedPlay] {
		if v > odds {
			odds = v
		}
	}
	var robotOrders, memberOrders int // 机器人领包个数，会员领包个数
	var robotCash, memberCash float64 // 机器人盈利，会员盈利
	for k, v := range redResList {
		if v.Account == redInfo.Account {
			adminMoney = v.Money
			adminRoyaltyMoneym = v.RoyaltyMoney
			adminRealMoney = v.RealMoney
			capital = redInfo.Capital
		} else {
			capital = common.DecimalMul(redInfo.RedEnvelopeAmount, float64(odds))
			if v.IsRobot != model.USER_IS_ROBOT_YES {
				memberOrders += 1
				memberCash = common.DecimalSum(memberCash, v.Money)
			} else {
				robotOrders += 1
				robotCash = common.DecimalSum(robotCash, v.Money)
			}
		}
		if userInfo.IsRobot != model.USER_IS_ROBOT_YES {
			// 会员发包
			if v.IsRobot == model.USER_IS_ROBOT_YES {
				//机器人抢包
				v.RobotWin = v.Money
			}
		}
		remark := fmt.Sprintf("注单号:%v,牛牛红包结算并返还保证金, 会员%v盈利金额%v, 保证金金额%v", v.OrderNo, v.Account, v.RealMoney, capital)
		if v.IsRobot != model.USER_IS_ROBOT_YES { // 不为机器人时进行金额变更操作
			err = new(UserServer).ChangeAmount(sess, redInfo.LineId, redInfo.AgencyId, v.Account, lineInfo.ApiUrl, lineInfo.TransType, v.UserId, common.DecimalSum(redResList[k].RealMoney, capital), remark, common.DecimalSub(0, capital))
			if err != nil {
				sess.Rollback()
				return nil, 0, 0, 0, 0, err
			}
		}
		// 修改红包领取记录状态
		_, err = redPacketLogBo.UpdateRedLogListStatus(sess, v)
		if err != nil {
			sess.Rollback()
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
			return nil, 0, 0, 0, 0, &validate.Err{Code: code.UPDATE_FAILED}
		}
	}
	// 修改红包状态
	redInfo.Status = model.RED_STATUS_OVER
	redInfo.Money = adminMoney
	redInfo.RoyaltyMoney = adminRoyaltyMoneym
	redInfo.RealMoney = adminRealMoney
	num, err := redPacketBo.UpdateRed(sess, redInfo, "status,money,royalty_money,real_money")
	if err != nil || num == 0 {
		sess.Rollback()
		if err != nil {
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		}
		return nil, 0, 0, 0, 0, &validate.Err{Code: code.UPDATE_FAILED}
	}
	cashList := make([]*structs.MemberCashRecord, 0)
	for _, v := range redResList {
		info := &structs.MemberCashRecord{
			LineId:     redInfo.LineId,
			AgencyId:   redInfo.AgencyId,
			GameType:   redInfo.RedType,
			GameName:   model.RED_ENVELOPE_TYPE[redInfo.RedType],
			CreateTime: utility.GetNowTimestamp(),
			Money:      v.RealMoney,
			FlowType:   model.MEMBER_RECORD_WIN,
			Remark:     fmt.Sprintf("牛牛红包结束结算"),
			UserId:     v.UserId,
			OrderNo:    v.OrderNo,
			Account:    v.Account,
		}
		cashList = append(cashList, info)
	}
	return cashList, robotOrders, robotCash, memberOrders, memberCash, nil
}

// 扫雷结算
func MineSettlement(sess *xorm.Session, lineInfo *structs.SystemLine, redInfo *structs.RedPacket, logList []structs.OrderRecord, roomInfo *structs.Room, user *structs.User) ([]*structs.MemberCashRecord, int, float64, int, float64, error) {
	// 获取已经领取了多少钱 庄家盈利
	var claimedMoney, winMoney float64
	var memberTotal float64           // 领包人总盈利
	var adminOrderNo string           // 发包人注单号
	var robotOrders, memberOrders int // 机器人领包个数，会员领包个数
	var robotCash, memberCash float64 // 机器人盈利，会员盈利
	adminLogInfo := structs.OrderRecord{}
	for _, v := range logList {
		// 处理领包人逻辑
		if v.UserId != redInfo.UserId && v.Account != redInfo.Account {
			if v.IsRobot != model.USER_IS_ROBOT_YES {
				memberOrders += 1
				// 累加会员抽水前盈利
				memberCash = common.DecimalSum(memberCash, v.Money)
			} else {
				robotOrders += 1
				// 累加机器人抽水前盈利
				robotCash = common.DecimalSum(robotCash, v.Money)
			}
			claimedMoney = common.DecimalSum(claimedMoney, v.ReceiveMoney) // 获取已经领取的总金额
			memberTotal = common.DecimalSum(memberTotal, v.Money)          // 累加领包人盈利
		}
	}
	for _, v := range logList {
		// 处理发包人逻辑
		if v.UserId == redInfo.UserId && v.Account == redInfo.Account {
			// 庄家注单
			adminLogInfo = v
			adminOrderNo = v.OrderNo
			adminLogInfo.ReceiveMoney = common.DecimalSub(redInfo.RedEnvelopeAmount, claimedMoney) // 发包人领取金额
			adminLogInfo.Status = model.RED_STATUS_OVER
			// 庄家真实盈利取反
			winMoney = common.DecimalSub(0, memberTotal)
			// 计算抽水金额
			if winMoney > 0 { // 庄家赢
				if v.IsRobot != model.USER_IS_ROBOT_YES {
					// 庄家不是机器人，计算抽水
					adminLogInfo.RoyaltyMoney = common.DecimalMul(winMoney, common.DecimalDiv(roomInfo.Royalty, 100))
				}
				adminLogInfo.RealMoney = common.DecimalSub(winMoney, adminLogInfo.RoyaltyMoney)
				adminLogInfo.Money = winMoney
			} else { // 庄家输
				adminLogInfo.Money = winMoney
				adminLogInfo.RealMoney = winMoney
			}
			vData := make(map[string]string)
			json.Unmarshal([]byte(v.Extra), &vData)
			vData["receive_money"] = strconv.FormatFloat(adminLogInfo.ReceiveMoney, 'f', 2, 64)
			vData["real_money"] = strconv.FormatFloat(adminLogInfo.RealMoney, 'f', 2, 64)
			b, _ := json.Marshal(vData)
			adminLogInfo.Extra = string(b)
			if v.IsRobot != model.USER_IS_ROBOT_YES {
				// 发包人不是机器人，记录有效投注
				adminLogInfo.ValidBet = common.DecimalSub(adminLogInfo.RedMoney, adminLogInfo.ReceiveMoney)
			} else {
				// 发包人是机器人，记录机器人盈利
				adminLogInfo.RobotWin = common.DecimalSub(0, memberCash)
			}
		}
	}

	// 修改注单
	_, err := redPacketLogBo.UpdateAdminRedLogInfo(sess, &adminLogInfo, int(redInfo.Id))
	if err != nil {
		sess.Rollback()
		golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		return nil, 0, 0, 0, 0, &validate.Err{Code: code.UPDATE_FAILED}
	}

	cashList := make([]*structs.MemberCashRecord, 0)
	recordInfo := new(structs.MemberCashRecord)
	recordInfo.LineId = redInfo.LineId
	recordInfo.AgencyId = redInfo.AgencyId
	recordInfo.GameType = redInfo.RedType
	recordInfo.GameName = model.RED_ENVELOPE_TYPE[redInfo.RedType]
	recordInfo.CreateTime = utility.GetNowTimestamp()
	recordInfo.Money = common.DecimalSum(winMoney, redInfo.RedEnvelopeAmount)
	recordInfo.FlowType = model.MEMBER_RECORD_WIN
	recordInfo.OrderNo = adminOrderNo
	if winMoney > 0 {
		recordInfo.Remark = fmt.Sprintf("注单号:%v,扫雷红包结算,返还发包领取金额%v元,盈利%v元,系统抽水%v元", adminOrderNo, redInfo.RedEnvelopeAmount, adminLogInfo.Money, adminLogInfo.RoyaltyMoney)
	} else {
		recordInfo.Remark = fmt.Sprintf("注单号:%v,扫雷红包结算,返还未领取金额%v元", adminOrderNo, common.DecimalSub(redInfo.RedEnvelopeAmount, claimedMoney))
	}
	recordInfo.UserId = redInfo.UserId
	recordInfo.Account = redInfo.Account
	cashList = append(cashList, recordInfo)
	// 修改红包状态
	num, err := redPacketBo.UpdateRed(sess, &structs.RedPacket{
		Id:          redInfo.Id,
		LineId:      redInfo.LineId,
		AgencyId:    redInfo.AgencyId,
		Status:      model.RED_STATUS_OVER,
		ReturnMoney: common.DecimalSub(redInfo.RedEnvelopeAmount, claimedMoney),
		RealMoney:   common.DecimalSub(winMoney, adminLogInfo.RoyaltyMoney),
	}, "status,return_money,real_money")
	if err != nil || num == 0 {
		sess.Rollback()
		if err != nil {
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		}
		return nil, 0, 0, 0, 0, &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 更新庄家金额
	// 查询会员是否为机器人
	has, userInfo := userBo.GetOneByAccount(sess, redInfo.LineId, redInfo.AgencyId, redInfo.Account)
	if !has {
		sess.Rollback()
		return nil, 0, 0, 0, 0, &validate.Err{Code: code.MEMBER_INFORMATION_QUERY_FAILED}
	}
	if userInfo.IsRobot != model.USER_IS_ROBOT_YES { // 不是机器人的时候才会进行金额变动
		var remark string
		if winMoney > 0 {
			remark = fmt.Sprintf("注单号:%v,扫雷红包结算,返还发包领取金额%v元,盈利%v元,系统抽水%v元", adminOrderNo, redInfo.RedEnvelopeAmount, adminLogInfo.Money, adminLogInfo.RoyaltyMoney)
		} else {
			remark = fmt.Sprintf("注单号:%v,扫雷红包结算,返还未领取金额%v元", adminOrderNo, common.DecimalSub(redInfo.RedEnvelopeAmount, claimedMoney))
		}
		if adminLogInfo.RealMoney > 0 {
			err = new(UserServer).ChangeAmount(sess, redInfo.LineId, redInfo.AgencyId, redInfo.Account, lineInfo.ApiUrl, lineInfo.TransType, redInfo.UserId, common.DecimalSum(redInfo.RedEnvelopeAmount, adminLogInfo.RealMoney), remark, 0)
			if err != nil {
				sess.Rollback()
				return nil, 0, 0, 0, 0, err
			}
		} else {
			err = new(UserServer).ChangeAmount(sess, redInfo.LineId, redInfo.AgencyId, redInfo.Account, lineInfo.ApiUrl, lineInfo.TransType, redInfo.UserId, common.DecimalSub(redInfo.RedEnvelopeAmount, claimedMoney), remark, 0)
			if err != nil {
				sess.Rollback()
				return nil, 0, 0, 0, 0, err
			}
		}
	}
	return cashList, robotOrders, robotCash, memberOrders, memberCash, nil
}

// 普通红包结算
func OrdinarySettlement(sess *xorm.Session, lineInfo *structs.SystemLine, redInfo *structs.RedPacket, logList []structs.OrderRecord) ([]*structs.MemberCashRecord, error) {
	cashList := make([]*structs.MemberCashRecord, 0)
	// 获取已经领取了多少钱 并返回领取金额（写入返还红包金额的现金记录）
	var returnMoney, claimedMoney float64
	for _, v := range logList {
		claimedMoney = common.DecimalSum(claimedMoney, v.ReceiveMoney) // 获取已经领取的总金额
	}
	returnMoney = common.DecimalSub(redInfo.RedEnvelopeAmount, claimedMoney) // 需要返还的金额
	// 实际盈利金额扣除本金
	redRealMoney := common.DecimalSub(redInfo.RealMoney, common.DecimalSub(redInfo.RedEnvelopeAmount, returnMoney))

	recordInfo := new(structs.MemberCashRecord)
	recordInfo.LineId = redInfo.LineId
	recordInfo.AgencyId = redInfo.AgencyId
	recordInfo.GameType = redInfo.RedType
	recordInfo.GameName = model.RED_ENVELOPE_TYPE[redInfo.RedType]
	recordInfo.CreateTime = utility.GetNowTimestamp()
	recordInfo.Money = returnMoney
	recordInfo.FlowType = model.MEMBER_RECORD_ORDINARY_RETURBN
	recordInfo.Remark = fmt.Sprintf("普通红包结束返还,返还金额%v,增加会员余额%v", returnMoney, returnMoney)
	recordInfo.UserId = redInfo.UserId
	recordInfo.Account = redInfo.Account
	cashList = append(cashList, recordInfo)
	// 修改红包状态
	num, err := redPacketBo.UpdateRed(sess, &structs.RedPacket{Id: redInfo.Id, LineId: redInfo.LineId, AgencyId: redInfo.AgencyId, Status: model.RED_STATUS_OVER, ReturnMoney: returnMoney, RealMoney: redRealMoney}, "status,return_money,real_money")
	if err != nil || num == 0 {
		sess.Rollback()
		if err != nil {
			golog.Error("RedPlay", "RedEnvelopeAmountCalculation", "err:", err)
		}
		return nil, &validate.Err{Code: code.UPDATE_FAILED}
	}
	return cashList, err
}
