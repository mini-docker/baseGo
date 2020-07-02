package server

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
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

// 机器人抢红包
// 添加机器人抢包时间轮，
// 红包发送后将红包id加入该时间轮，
// 时间为3/4红包结束时间，
// 例如：红包游戏时间为2分钟120秒，那机器人介入抢红包时间为发送红包后90秒（也可以放在同一个时间轮根据类型判断是结算还是机器人介入抢包）

// 2.机器人抢包逻辑设计：

// 90%的概率抢红包;
// 根据剩余红包个数随机抢n个红包，n>0;
// 80%的概率抢赢钱的红包（如果剩余红包全为赢或输，是否必须抢？） todo

type RobotGrabRed struct{}

// 抢红包操作
func (RobotGrabRed) GrabRed(lineId, agencyId string, roomId, redId int) error {
	golog.Info("RobotGrabRed", "GrabRed", "进入机器人抢包逻辑", redId)
	redisClient := conf.GetRedis().Get()
	// 检查结算红包是否加锁
	res, err := redisClient.Do("Exists", fmt.Sprintf("grabBag_%v", redId))
	fResult, _ := strconv.Atoi(string(res.([]byte)))
	if fResult == 1 {
		golog.Info("RobotGrabRed", "GrabRed", "红包被锁", redId)
		return nil
	}
	// 加锁
	redisClient.Do("Set", fmt.Sprintf("grabBag_%v", redId), 1, time.Second*10)
	// 获取线路信息 判断线路是否存在
	_, err = GetStytemLineInfo(lineId)
	if err != nil {
		golog.Error("RobotGrabRed", "GrabRed", "error:", err)
		return &validate.Err{Code: code.LINE_QUERY_FAILED}
	}

	// 根据房间ID和红包ID 获取红包类型和红包玩法 和会员下注金额等数据
	sess := conf.GetXormSession()
	defer sess.Close()
	has, room := roomBo.GetOne(sess, roomId)
	if !has {
		golog.Error("RobotGrabRed", "GrabRed", "error:游戏群查询失败", nil, lineId, agencyId, roomId)
		return &validate.Err{Code: code.GAME_GROUP_QUERY_FAILED}
	}
	if room.Status != 1 {
		golog.Error("RobotGrabRed", "GrabRed", "error:游戏群停用", nil, room.Id)
		return &validate.Err{Code: code.ROOM_IS_FROBIDDEN}
	}
	if room.RobotGrabPacket == model.ROBOT_GRAB_PACKET_NO {
		// 机器人抢红包未开启 结束此轮抢红包
		golog.Info("RobotGrabRed", "GrabRed", "机器人抢包未开启", redId)
		return nil
	}

	// 查询红包信息
	has, redInfo := redPacketBo.ByRoomIdGetRedInfo(sess, redId)
	if !has && redInfo.DeleteTime != 0 {
		golog.Error("RobotGrabRed", "GrabRed", "error:红包信息查询失败", nil)
		return &validate.Err{Code: code.RED_PACKET_QUERY_FAILED}
	}
	if redInfo.Status != model.RED_STATUS_NORMAL { // 红包状态不为进行中 结束此轮抢红包
		golog.Info("RobotGrabRed", "GrabRed", "红包状态不为进行中", redId)
		return nil
	}

	if redInfo.EndTime <= utility.GetNowTimestamp() {
		// 红包已结束 结束此轮抢红包
		golog.Info("RobotGrabRed", "GrabRed", "红包已结束", redId)
		return nil
	}

	// 查询该代理下是否存在机器人
	userList, err := userBo.GetRobotListByAgenncyId(sess, lineId, agencyId)
	if err != nil {
		golog.Error("RobotGrabRed", "GrabRed", "error:红包信息查询失败", nil)
		return &validate.Err{Code: code.RED_PACKET_QUERY_FAILED}
	}
	if len(userList) <= 0 {
		// 没有机器人 结束此轮抢红包
		golog.Info("RobotGrabRed", "GrabRed", "没有机器人", agencyId)
		return nil
	}

	// 去掉发包人
	robots := make([]structs.User, 0)
	for _, r := range userList {
		if r.Account == redInfo.Account {
			continue
		}
		robots = append(robots, r)
	}

	var (
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	// 随机是否抢红包
	if r.Intn(10) == 9 { // 随机0-10 实际随机 0-9 的数字
		// 90%的概率抢红包 所以只要数字等于9 就不抢红包 结束此轮抢红包
		golog.Info("RobotGrabRed", "GrabRed", "这么倒霉90%的概率，你都能不抢", nil)
		return nil
	}

	logKey := fmt.Sprintf("%d_%d_redLog", redInfo.Id, roomId)
	fs, err := conf.GetRedis().Get().Do("HGetAll", logKey)
	keys := fs.(map[string]string)
	if err != nil {
		// 红包数据不存在
		golog.Error("RobotGrabRed", "GrabRed", "红包缓存不存在", err)
		return &validate.Err{Code: code.RED_ENVELOPE_CLAIM_FAILED}
	}
	if len(keys) == 0 {
		// 修改红包状态为已领完
		redPacketBo.UpdateRedStatus(sess, int(redInfo.Id), model.RED_STATUS_FINISHED)
		// 此轮没有剩余红包 结束此轮抢红包
		golog.Info("RobotGrabRed", "GrabRed", "红包已抢完", redId)
		return nil
	}

	// 根据红包个数随机出抢红包的数量
	// 随机范围
	var grabRed int // 抢红包个数
	{
		var rRangeNum int
		userCount := len(robots)
		redLogCount := len(keys)
		if userCount > redLogCount {
			rRangeNum = redLogCount
		} else {
			rRangeNum = userCount
		}
		grabRed = r.Intn(rRangeNum)
	}
	// 至少抢一个包
	if grabRed == 0 {
		grabRed = 1
	}
	if grabRed > len(keys) {
		grabRed = len(keys)
	}

	// 判断群是否开启控杀
	contorlKill := false
	if room.ControlKill == 1 {
		// 开启控杀
		if r.Intn(10) < 9 {
			contorlKill = true
		}
	}

	// 抢出来的红包
	grabRedLog := make([]*structs.OrderRecord, 0)
	grabReds := make([]*structs.OrderRecord, 0)
	var i = 0
	if contorlKill {
		// 控杀成功
		switch redInfo.RedType {
		case model.NIUNIU_RED_ENVELOPE:
			fmt.Println(fmt.Sprintf("牛牛控杀,抢包%v个", grabRed))
			// 查询发包人注单
			senderOrder, err := redPacketBo.GetSenderOrder(sess, redId, redInfo.UserId)
			if err != nil {
				golog.Error("RobotGrabRed", "GrabRed", "获取牛牛发包人注单失败", err)
				return &validate.Err{Code: code.QUERY_FAILED}
			}
			// 计算发包人牛数
			senderNiu := NiuNiuCalculation(senderOrder.ReceiveMoney)
			// 挑选牛数最大的红包
			for j := 0; j < grabRed; j++ {
				var maxNum int
				var key string
				var money float64
				maxNum = 1
				for k, v := range keys {
					redLog := new(structs.OrderRecord)
					json.Unmarshal([]byte(v), &redLog)
					newNum := NiuNiuCalculation(redLog.ReceiveMoney)
					// 当前包牛数大于历史最大牛数
					if newNum == 0 && maxNum != 0 {
						key = k
						money = redLog.ReceiveMoney
						maxNum = newNum
					} else if newNum > maxNum {
						// 更新缓存数据
						maxNum = newNum
						key = k
						money = redLog.ReceiveMoney
					} else if newNum == maxNum {
						// 当前包牛数等于历史最大牛数，且当前包金额大于历史最大金额
						if redLog.ReceiveMoney > money {
							// 更新缓存数据
							maxNum = newNum
							key = k
							money = redLog.ReceiveMoney
						}
					}
				}
				fmt.Println(fmt.Sprintf("最大牛数:%v,金额:%v,key:%v", maxNum, money, key))
				if maxNum < senderNiu || (maxNum == senderNiu && money < senderOrder.ReceiveMoney) {
					fmt.Println(fmt.Sprintf("最大牛数:%v,金额:%v,小于发包人牛数:%v,金额:%v,退出循环", maxNum, money, senderNiu, senderOrder.ReceiveMoney))
					break
				}
				redLogData := new(structs.OrderRecord)

				json.Unmarshal([]byte(keys[key]), &redLogData)
				grabReds = append(grabReds, redLogData)
				delete(keys, key)                   // 删除已取出的红包
				redisClient.Do("HDel", logKey, key) // 删除已取出的包缓存
			}
		case model.MINESWEEPER_RED_PACKET:
			fmt.Println(fmt.Sprintf("扫雷控杀,抢包%v个", grabRed))
			// 挑选不中雷的包
			for k, v := range keys {
				if i >= grabRed {
					break
				}
				redLogData := new(structs.OrderRecord)
				json.Unmarshal([]byte(v), &redLogData)
				if (int(common.DecimalMul(redLogData.ReceiveMoney, 100)) % 10) != redInfo.Mine {
					fmt.Println(fmt.Sprintf("抢包金额:%v,key:%v,雷值:%v", redLogData.ReceiveMoney, k, redInfo.Mine))
					grabReds = append(grabReds, redLogData)
					delete(keys, k) // 删除已取出的红包
					// 删除当前已取出的红包
					conf.GetRedis().Get().Do("HDel", logKey, k)
					i++
				}
			}
		}
		if len(grabReds) == 0 {
			golog.Info("RobotGrabRed", "GrabRed", "红包不存在赢钱的包，控杀情况下机器人不抢包", redId)
			return &validate.Err{Code: code.RED_ENVELOPE_CLAIM_FAILED}
		}
	} else {
		// 不控杀，随机出包
		for k, v := range keys {
			if i >= grabRed {
				break
			}
			redLogData := new(structs.OrderRecord)
			json.Unmarshal([]byte(v), &redLogData)
			grabReds = append(grabReds, redLogData)
			// 删除当前已取出的红包
			conf.GetRedis().Get().Do("HDel", logKey, k)
			i++
		}
	}

	for _, v := range grabReds {
		redLog := v
		// 随机一个机器人
		userIndex := r.Intn(len(robots) - 1) // 随机会员下标
		user := robots[userIndex]
		redLog.LineId = lineId
		redLog.AgencyId = agencyId
		redLog.RedSender = redInfo.Account
		redLog.GameType = redInfo.RedType
		redLog.GamePlay = redInfo.RedPlay
		redLog.RoomId = int(room.Id)
		redLog.RoomName = room.RoomName
		redLog.OrderNo = OderNo(redLog.GameType)
		redLog.RedId = int(redInfo.Id)
		redLog.RedMoney = redInfo.RedEnvelopeAmount
		redLog.RedNum = redInfo.RedEnvelopesNum
		redLog.Royalty = room.Royalty
		redLog.GameTime = room.GameTime
		redLog.ReceiveTime = utility.GetNowTimestamp()
		redLog.RedStartTime = redInfo.CreateTime
		redLog.Status = 0
		redLog.Account = user.Account
		redLog.UserId = user.Id
		redLog.IsRobot = user.IsRobot
		redLog.IsFreeDeath = user.IsGroupOwner
		vData := make(map[string]string)
		switch redInfo.RedType {
		case model.NIUNIU_RED_ENVELOPE:
			vData["memberNum"] = strconv.Itoa(NiuNiuCalculation(redLog.ReceiveMoney))
		case model.MINESWEEPER_RED_PACKET:
			if room.GamePlay == model.MINESWEEPER_UNFIXED_ODDS {
				// 不固定赔率，根据红包个数获取赔率
				room.Odds = model.MINESWEEPER_UNFIXED_PLAY[redInfo.RedEnvelopesNum]
			}
			vData["thunderNum"] = strconv.Itoa(redInfo.Mine)
			vData["odds"] = strconv.FormatFloat(room.Odds, 'f', 2, 64)
		}
		vDataByte, _ := json.Marshal(vData)
		redLog.Extra = string(vDataByte)
		if redLog.GameType == model.MINESWEEPER_RED_PACKET {
			if (int(common.DecimalMul(redLog.ReceiveMoney, 100)) % 10) == redInfo.Mine {
				// 中雷
				redLog.Status = model.RED_STATUS_OVER
				// 中雷扣款 红包本金*赔率-红包领取金额
				redLog.Money = common.DecimalSub(0, common.DecimalSub(common.DecimalMul(redInfo.RedEnvelopeAmount, float64(room.Odds)), redLog.ReceiveMoney))
				redLog.RealMoney = redLog.Money
				if redInfo.IsRobot == model.USER_IS_ROBOT_NO {
					// 发包人是会员,记录机器人盈利
					redLog.RobotWin = redLog.Money
				}
			} else {
				// 未中雷
				redLog.Status = model.RED_STATUS_OVER
				redLog.Money = common.DecimalSum(0, redLog.ReceiveMoney)
				redLog.RealMoney = redLog.Money
				if redInfo.IsRobot == model.USER_IS_ROBOT_NO {
					// 发包人是会员,记录机器人盈利
					redLog.RobotWin = redLog.Money
				}
			}
		}
		grabRedLog = append(grabRedLog, redLog)                      // 将抢到的红包放入抢到的红包列表中
		robots = append(robots[:userIndex], robots[userIndex+1:]...) // 将选中的机器人从机器人列表中删除
	}

	//红包数据处理好之后插入红包注单记录
	sess.Begin()
	count, err := redPacketLogBo.InsertRedPacketLogList(sess, grabRedLog...)
	if count <= 0 || err != nil {
		sess.Rollback()
		golog.Info("RobotGrabRed", "GrabRed", "插入注单数据失败:", err, count, "len(grabRedLog)", len(grabRedLog))
		return &validate.Err{Code: code.RED_ENVELOPE_CLAIM_FAILED}
	}

	// 写入红包记录之后  写现金记录
	switch redInfo.RedType {
	case model.MINESWEEPER_RED_PACKET:
		// 扫雷红包需要写现金记录和会员金额操作
		// 增加/扣除会员余额
		userMoney := make(map[int]float64)
		data := make([]*structs.MemberCashRecord, 0)
		var roomRoyaltyMoney float64
		var recordMoney float64
		for _, v := range grabRedLog {
			userMoney[v.UserId] = v.RealMoney

			// 现金记录
			var royaltyMoney float64 // 抽水金额
			royaltyMoney = v.RoyaltyMoney
			// 写入现金记录
			// 中雷了需要写入2条现金记录 1条闲家输钱的 1条庄家赢钱的 并且变更红包的相关信息(输赢金额)
			recordInfo := new(structs.MemberCashRecord)
			recordInfo.LineId = lineId
			recordInfo.AgencyId = redInfo.AgencyId
			recordInfo.GameType = redInfo.RedType
			recordInfo.GameName = model.RED_ENVELOPE_TYPE[redInfo.RedType]
			recordInfo.CreateTime = utility.GetNowTimestamp()
			if v.RealMoney > 0 {
				recordInfo.Money = v.RealMoney
				recordInfo.FlowType = model.MEMBER_RECORD_WIN
				recordInfo.Remark = fmt.Sprintf("领取扫雷红包,红包金额%v,增加会员余额%v", v.ReceiveMoney, recordInfo.Money)
				recordInfo.UserId = v.UserId
				recordInfo.OrderNo = v.OrderNo
				recordInfo.Account = v.Account
				data = append(data, recordInfo)
			} else {
				// 闲家现金记录
				recordInfo.Money = v.RealMoney
				recordInfo.FlowType = model.MEMBER_RECORD_LOSE
				recordInfo.Remark = fmt.Sprintf("领取扫雷红包,中雷,红包金额%v,扣除会员余额%v", v.ReceiveMoney, v.RealMoney)
				recordInfo.UserId = v.UserId
				recordInfo.Account = v.Account
				recordInfo.OrderNo = v.OrderNo
				data = append(data, recordInfo)
				recordMoney = recordMoney + recordInfo.Money
			}
			roomRoyaltyMoney = common.DecimalSub(roomRoyaltyMoney, royaltyMoney)
		}

		// 写入现金记录
		num, err := memberCashRecordBo.Inster(sess, data...)
		if err != nil || num <= 0 {
			sess.Rollback()
			golog.Error("RedPacketService", "GrabRedEnvelope", "写入现金记录失败:", err)
			return &validate.Err{Code: code.WRITING_CASH_RECORD_FAILED}
		}
	}
	sess.Commit()
	if len(keys) == 0 {
		// 将结算时间提前到5秒后
		keyData := map[string]string{
			"key":      "redPacketSettle",
			"redId":    strconv.Itoa(int(redInfo.Id)),
			"roomId":   fmt.Sprint(room.Id),
			"gameTime": strconv.Itoa(room.GameTime),
			"lineId":   redInfo.LineId,
			"agencyId": redInfo.AgencyId,
		}
		b, _ := json.Marshal(keyData)
		AddTimeWheel(string(b), 5)
	}
	// 解锁
	redisClient.Do("Del", fmt.Sprintf("grabBag_%v", redId))
	golog.Info("GrabEnd", "GrabEnd", "自动抢包逻辑结束", redId)
	return nil
}
