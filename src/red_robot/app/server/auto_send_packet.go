package server

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/structs"
	"baseGo/src/red_robot/conf"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
)

var (
	REDIS_KEY = "lastSendPacket"
)

func InitAutoSendPacket() {
	s := gocron.NewScheduler()
	s.Every(1).Minutes().Do(sendPacket)
	<-s.Start()
}

var (
	r = mrand.New(mrand.NewSource(time.Now().UnixNano()))
)

func sendPacket() {
	// 获取开启自动发包全部红包群
	redisClient := conf.GetRedis().Get()
	cmd, err := redisClient.Do("HKeys", REDIS_KEY)
	if err != nil {
		golog.Error("AutoSendPacket", "sendPacket", "redis err:", err)
		return
	}
	keys := cmd.([]interface{})
	// 遍历获取数据
	if len(keys) > 0 {
		sess := conf.GetXormSession()
		defer sess.Close()
		for _, v := range keys {
			key := string(v.([]byte))
			fResult, err := redisClient.Do("HGet", REDIS_KEY, key)
			roomCmd := string(fResult.([]byte))
			// redis获取失败，继续执行下一个群
			if err != nil {
				golog.Error("AutoSendPacket", "sendPacket", "redis err:", err)
				continue
			} else {
				// 解析群信息
				room := new(structs.Room)
				err := json.Unmarshal([]byte(roomCmd), &room)
				if err != nil {
					golog.Error("AutoSendPacket", "sendPacket", "json unmarshal err:", err)
					continue
				}
				// 判断群是否存在
				has, roomdb := roomBo.GetOne(sess, int(room.Id))
				if !has {
					golog.Error("AutoSendPacket", "sendPacket", "error:游戏群查询失败", nil, room.Id)
					redisClient.Do("HDel", REDIS_KEY, key)
					continue
				}
				if roomdb.Status != 1 {
					golog.Error("AutoSendPacket", "sendPacket", "error:游戏群停用", nil, room.Id)
					redisClient.Do("HDel", REDIS_KEY, key)
					continue
				}
				// 判断当前时间是否超过该群设置的自动发包时间
				if utility.GetNowTimestamp()-room.LastTime > roomdb.RobotSendPacketTime*60 {
					// 符合发包时间，自动发包
					// 获取机器人信息
					robots, err := new(bo.User).GetRobotListByAgenncyId(sess, roomdb.LineId, roomdb.AgencyId)
					if err != nil {
						golog.Error("AutoSendPacket", "sendPacket", "get robots err:", err)
						continue
					}
					if len(robots) <= 0 {
						golog.Error("AutoSendPacket", "sendPacket", "no robots ", nil)
						continue
					}
					num := r.Intn(len(robots))
					user := robots[num]
					//红包数据组装
					redInfo := new(structs.RedPacket)
					// 随机发包金额
					redAmount := randInt(int64(roomdb.MinMoney), int64(roomdb.MaxMoney+1))
					if redAmount < int(roomdb.MinMoney) {
						redAmount = int(roomdb.MinMoney)
					}
					// 随机发包个数
					if roomdb.RedMinNum < 2 {
						roomdb.RedMinNum = 2
					}
					redNum := randInt(int64(roomdb.RedMinNum), int64(roomdb.RedNum+1))
					if redNum < roomdb.RedMinNum {
						redNum = roomdb.RedMinNum
					}
					redInfo.RedEnvelopeAmount = float64(redAmount)
					redInfo.RedEnvelopesNum = redNum
					redInfo.AgencyId = roomdb.AgencyId
					redInfo.LineId = roomdb.LineId
					redInfo.UserId = user.Id
					redInfo.Account = user.Account
					redInfo.CreateTime = utility.GetNowTimestamp()
					redInfo.RedType = roomdb.GameType
					redInfo.RedPlay = roomdb.GamePlay
					redInfo.RoomId = int(roomdb.Id)
					redInfo.RoomName = roomdb.RoomName
					redInfo.IsRobot = model.USER_IS_ROBOT_YES
					redInfo.Status = model.RED_STATUS_NORMAL
					// 判断游戏类型
					switch roomdb.GameType {
					case model.NIUNIU_RED_ENVELOPE:
						// 牛牛

					case model.MINESWEEPER_RED_PACKET:
						// 扫雷
						mine := r.Intn(10)
						redInfo.Mine = mine
						if roomdb.GamePlay == model.MINESWEEPER_UNFIXED_ODDS {
							// 不固定赔率，根据红包个数获取赔率
							roomdb.Odds = model.MINESWEEPER_UNFIXED_PLAY[redInfo.RedEnvelopesNum]
						}
					}
					redInfo.Capital = 0
					redInfo.EndTime = redInfo.CreateTime + (roomdb.GameTime * 60)
					// 插入红包数据
					count, err := new(bo.RedPacket).InsertRedPacket(sess, redInfo)
					if err != nil || count == 0 {
						sess.Rollback()
						if err != nil {
							golog.Error("AutoSendPacket", "CreateRedPacket", "err:", err)
						}
						continue
					}
					// 生成红包
					redLogList, err := redPacketAllocation(redNum, float64(redAmount), 0.01)
					if err != nil || len(redLogList) == 0 {
						sess.Rollback()
						golog.Error("AutoSendPacket", "CreateRedPacket", "new packet err", err)
						continue
					}

					// 判断群是否开启控杀
					var controlKill = false
					if roomdb.ControlKill == 1 {
						// 开启控杀
						if r.Intn(10) < 7 {
							controlKill = true
						}
					}

					switch redInfo.RedType {
					// 牛牛红包
					case model.NIUNIU_RED_ENVELOPE:
						// 更新会员已抵押的保证金
						user.Capital = 0
						err = new(bo.User).UpdateCapitalIncr(sess, user.Id, 0)
						if err != nil {
							sess.Rollback()
							golog.Error("AutoSendPacket", "CreateRedPacket", "err:", err)
							continue
						}
						logInfo := new(structs.OrderRecord)
						// 牛牛红包时发包人需要提前抢一个红包
						if controlKill {
							fmt.Println("牛牛发包控杀")
							var num, redIndex int
							var money float64
							num = 1
							// 控杀成功，领取牛数最大的包
							for k, v := range redLogList {
								newNum := NiuNiuCalculation(v.ReceiveMoney)
								if newNum == 0 && num != 0 {
									redIndex = k
									money = v.ReceiveMoney
									num = newNum
								} else if newNum > num && newNum != 0 {
									redIndex = k
									money = v.ReceiveMoney
									num = newNum
								} else if newNum == num {
									if v.ReceiveMoney > money {
										redIndex = k
										money = v.ReceiveMoney
										num = newNum
									}
								}
							}
							logInfo = redLogList[redIndex]
							if redIndex == len(redLogList)-1 {
								redLogList = redLogList[:len(redLogList)-1]
							} else {
								redLogList = append(redLogList[:redIndex], redLogList[redIndex+1:]...)
							}
							fmt.Println(fmt.Sprintf("发包人领取最大牛数红包,牛数为：%v,金额为：%v", num, money))
						} else {
							// 控杀失败失败，随机领取一个包
							var redIndex int
							if len(redLogList) > 1 {
								redIndex = r.Intn(len(redLogList))
							} else {
								redIndex = 0
							}
							logInfo = redLogList[redIndex]
							if redIndex == len(redLogList)-1 {
								redLogList = redLogList[:len(redLogList)-1]
							} else {
								redLogList = append(redLogList[:redIndex], redLogList[redIndex+1:]...)
							}
						}

						logInfo.LineId = roomdb.LineId
						logInfo.AgencyId = roomdb.AgencyId
						logInfo.Account = user.Account
						logInfo.UserId = user.Id
						logInfo.RedSender = user.Account
						logInfo.GameType = redInfo.RedType
						logInfo.GamePlay = redInfo.RedPlay
						logInfo.RoomId = int(roomdb.Id)
						logInfo.RoomName = roomdb.RoomName
						logInfo.OrderNo = OderNo(logInfo.GameType)
						logInfo.RedId = int(redInfo.Id)
						logInfo.RedMoney = redInfo.RedEnvelopeAmount
						logInfo.RedNum = redInfo.RedEnvelopesNum
						logInfo.Royalty = roomdb.Royalty
						logInfo.GameTime = roomdb.GameTime
						logInfo.ReceiveTime = utility.GetNowTimestamp()
						logInfo.RedStartTime = redInfo.CreateTime
						logInfo.Status = 0
						logInfo.IsRobot = user.IsRobot
						logInfo.IsFreeDeath = user.IsGroupOwner
						vData := make(map[string]string)
						niuNum := strconv.Itoa(NiuNiuCalculation(logInfo.ReceiveMoney))
						vData["memberNum"] = niuNum
						vData["adminNum"] = niuNum
						vDataByte, _ := json.Marshal(vData)
						logInfo.Extra = string(vDataByte)
						//	写入红包记录
						count, err := new(bo.RedPacketLog).InsertRedPacketLog(sess, logInfo)
						if count <= 0 || err != nil {
							sess.Rollback()
							golog.Error("AutoSendPacket", "CreateRedPacket", "err:", err)
							continue
						}
						// 扫雷红包
					case model.MINESWEEPER_RED_PACKET:
						// 控杀成功判断已生成的红包是否存在雷包
						if controlKill {
							fmt.Println("扫雷发包控杀")
							var i = 0
							for _, v := range redLogList {
								if (int(common.DecimalMul(v.ReceiveMoney, 100)) % 10) == redInfo.Mine {
									i++
								}
							}
							if i == 0 {
								// 没有雷包,随机雷包个数
								n := r.Intn(redNum)
								if n == 0 {
									n = 1
								}
								// 遍历生成的红包修改金额
								for _, v := range redLogList {
									if (int(common.DecimalMul(v.ReceiveMoney, 100)) % 10) != redInfo.Mine {
										if i >= n {
											break
										}
										// 修改当前包金额为中雷
										addMoney, _ := strconv.ParseFloat(strconv.Itoa(redInfo.Mine-(int(common.DecimalMul(v.ReceiveMoney, 100))%10)), 64)
										v.ReceiveMoney = common.DecimalSum(v.ReceiveMoney, common.DecimalDiv(addMoney, 100))
										// 将添加的金额从下一个包内扣出来
										for _, x := range redLogList {
											if x.ReceiveMoney > common.DecimalDiv(addMoney, 100) {
												x.ReceiveMoney = common.DecimalSub(redLogList[0].ReceiveMoney, common.DecimalDiv(addMoney, 100))
												break
											}
										}
										i++
									}
								}
							}
							fmt.Println(fmt.Sprintf("红包个数为：%v个，有雷的红包个数为：%v个", redNum, i))
						}

						// 扫雷红包需要生成一条庄家的注单
						logInfo := new(structs.OrderRecord)
						logInfo.LineId = roomdb.LineId
						logInfo.AgencyId = roomdb.AgencyId
						logInfo.Account = user.Account
						logInfo.UserId = user.Id
						logInfo.RedSender = user.Account
						logInfo.GameType = redInfo.RedType
						logInfo.GamePlay = redInfo.RedPlay
						logInfo.RoomId = int(roomdb.Id)
						logInfo.RoomName = roomdb.RoomName
						logInfo.OrderNo = OderNo(logInfo.GameType)
						logInfo.RedId = int(redInfo.Id)
						logInfo.RedMoney = redInfo.RedEnvelopeAmount
						logInfo.RedNum = redInfo.RedEnvelopesNum
						logInfo.Royalty = roomdb.Royalty
						logInfo.GameTime = roomdb.GameTime
						logInfo.ReceiveTime = utility.GetNowTimestamp()
						logInfo.ReceiveMoney = 0
						logInfo.RedStartTime = redInfo.CreateTime
						logInfo.Status = 0
						logInfo.IsRobot = user.IsRobot
						logInfo.IsFreeDeath = user.IsGroupOwner
						vData := make(map[string]string)
						vData["thunderNum"] = strconv.Itoa(redInfo.Mine)
						vData["odds"] = strconv.FormatFloat(roomdb.Odds, 'f', 2, 64)
						vDataByte, _ := json.Marshal(vData)
						logInfo.Extra = string(vDataByte)
						//	写入红包记录
						count, err := new(bo.RedPacketLog).InsertRedPacketLog(sess, logInfo)
						if count <= 0 || err != nil {
							sess.Rollback()
							golog.Error("AutoSendPacket", "CreateRedPacket", "err:", err)
							continue
						}
						// 直接扣会员余额
						err = new(bo.User).UpdateUserBalance(sess, roomdb.LineId, roomdb.AgencyId, user.Id, float64(redAmount))
						if err != nil {
							sess.Rollback()
							golog.Error("AutoSendPacket", "CreateRedPacket", "err:", err)
							continue
						}
						// 写现金记录
						recordInfo := new(structs.MemberCashRecord)
						recordInfo.LineId = roomdb.LineId
						recordInfo.AgencyId = roomdb.AgencyId
						recordInfo.GameType = redInfo.RedType
						recordInfo.GameName = model.RED_ENVELOPE_TYPE[roomdb.GameType]
						recordInfo.FlowType = 1
						recordInfo.Money = float64(redAmount)
						recordInfo.Remark = fmt.Sprintf("发扫雷红包,红包金额%v,扣除会员余额%v", recordInfo.Money, recordInfo.Money)
						recordInfo.CreateTime = utility.GetNowTimestamp()
						recordInfo.UserId = redInfo.UserId
						recordInfo.Account = redInfo.Account
						count1, err := new(bo.MemberCashRecord).Inster(sess, recordInfo)
						if err != nil || count1 <= 0 {
							sess.Rollback()
							if err != nil {
								golog.Error("AutoSendPacket", "CreateRedPacket", "err:", err)
							}
							continue
						}
					}
					// 将剩余的红包存入redis中
					logKey := fmt.Sprintf("%d_%d_redLog", redInfo.Id, roomdb.Id)
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
						return
					}

					// 红包发送成功后更新红包最近发包金额
					if roomdb.RobotSendPacket == model.ROBOT_SEND_PACKET_OFF {
						robotKey := "lastSendPacket"
						roomdb.LastTime = utility.GetNowTimestamp()
						// 更新最近发红包时间
						str, _ := json.Marshal(roomdb)
						conf.GetRedis().Get().Do("HSet", robotKey, fmt.Sprint(roomdb.Id), string(str))
					}

					// 写入时间轮
					keyData := map[string]string{
						"key":      "redPacketSettle",
						"redId":    strconv.Itoa(int(redInfo.Id)),
						"roomId":   fmt.Sprint(roomdb.Id),
						"gameTime": strconv.Itoa(roomdb.GameTime),
						"lineId":   roomdb.LineId,
						"agencyId": roomdb.AgencyId,
					}
					b, _ := json.Marshal(keyData)
					AddTimeWheel(string(b), redInfo.EndTime-utility.GetNowTimestamp())
					if roomdb.RobotGrabPacket == model.ROBOT_GRAB_PACKET_OFF { // 机器人抢红包开启
						GrabKeyData := map[string]string{
							"key":      "robotGrabRedPacket",
							"redId":    strconv.Itoa(int(redInfo.Id)),
							"roomId":   fmt.Sprint(roomdb.Id),
							"lineId":   roomdb.LineId,
							"agencyId": roomdb.AgencyId,
						}
						b, _ = json.Marshal(GrabKeyData)
						AddTimeWheel(string(b), (redInfo.EndTime - utility.GetNowTimestamp() - 20))
					}

					// 给会员推送红包信息
					{
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
						msgData["redSender"] = user.Account
						msgData["redMoney"] = redAmount
						msgData["redNum"] = redNum
						msgData["redId"] = int(redInfo.Id)
						msgData["gameType"] = redInfo.RedType
						msgData["gameTypeName"] = model.RED_ENVELOPE_TYPE[redInfo.RedType]
						msgData["gamePlay"] = redInfo.RedPlay
						msgData["gamePlayName"] = model.RED_ENVELOPE_TYPE_PLAY[redInfo.RedType][redInfo.RedPlay]
						msgData["gameTime"] = roomdb.GameTime
						msgData["createTime"] = redInfo.CreateTime
						msgData["mine"] = redInfo.Mine
						msgData["odds"] = roomdb.Odds
						msgData["redStatus"] = 1
						b, _ := json.Marshal(msgData)
						data := &RoomReq{
							Operation:   4,
							RoomId:      int(roomdb.Id),
							LineId:      roomdb.LineId,
							AgencyId:    roomdb.AgencyId,
							Msg:         string(b),
							MsgType:     1,
							SendId:      user.Id,
							ReceiveType: 2,
						}
						msgHis := &structs.MessageHistory{
							LineId:     roomdb.LineId,
							AgencyId:   roomdb.AgencyId,
							MsgType:    1,
							MsgContent: string(b),
							SenderId:   user.Id,
							SenderName: user.Account,
							Status:     1,
							SendTime:   utility.GetNowTimestamp(),
							RoomId:     int(roomdb.Id),
						}
						err = new(bo.MessageHistory).SaveMessageHistory(sess, msgHis)
						if err != nil {
							continue
						}

						// 转发im发送信息
						SendRoomMessageFunc("/push/room", data)
					}
				} else {
					continue
				}
			}
		}
	}
}

func randInt(min, max int64) int {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		randInt(min, max)
	}
	return int(i.Int64())
}
