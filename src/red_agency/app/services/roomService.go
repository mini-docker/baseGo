package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

type RoomService struct{}

// 查询房间列表
func (RoomService) QueryRoomList(lineId string, agencyId string, startTime, endTime, gameType, status int, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部房间信息
	count, rooms, err := RoomBo.QueryRoomList(sess, lineId, agencyId, startTime, endTime, gameType, status, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = rooms
	pageResp.Count = count
	return pageResp, nil
}

func randInt(min, max int64) int {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		randInt(min, max)
	}
	return int(i.Int64())
}

// 添加房间
//func (RoomService) AddRoom(lineId string, agencyId string, roomName string, gameType int, gamePlay int, maxMoney string,
//	minMoney string, odds string, redNum int, royalty string, gameTime int, redMinNum int) error {
func (RoomService) AddRoom(lineId string, agencyId string, roomName string, gameType int, gamePlay int, maxMoney string,
	minMoney string, odds string, redNum int, royalty string, gameTime int, redMinNum int, roomType int, freeFromDeath int, creatorId int, creator, ip string, isAdmin int) error {
	maMoney, _ := strconv.ParseFloat(maxMoney, 64)
	miMoney, _ := strconv.ParseFloat(minMoney, 64)
	modds, _ := strconv.ParseFloat(odds, 64)
	mRoyalty, _ := strconv.ParseFloat(royalty, 64)
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取当前最大房间号
	roomNo, err := RoomBo.GetMaxRoomNo(sess)
	if err != nil {
		return &validate.Err{Code: code.QUERY_FAILED}
	}

	room := new(structs.Room)
	// 添加房间
	room.LineId = lineId
	room.AgencyId = agencyId
	if roomName == "" {
		if gameType == 1 {
			room.RoomName = "牛牛"
			// 1经典牛牛 2平倍牛牛 3超倍牛牛
			if gamePlay == 1 {
				room.RoomName += fmt.Sprintf("经典群0-%d分钟", gameTime)
			} else if gamePlay == 2 {
				room.RoomName += fmt.Sprintf("平倍群0-%d分钟", gameTime)
			} else {
				room.RoomName += fmt.Sprintf("超倍群0-%d分钟", gameTime)
			}
		} else {
			room.RoomName = fmt.Sprintf("扫雷红包群0-%d分钟", gameTime)
		}
	} else {
		room.RoomName = roomName
	}

	// 生成群主
	groupOwner := new(structs.User)
	groupOwner.AgencyId = agencyId
	groupOwner.LineId = lineId
	groupOwner.Status = 1
	groupOwner.CreateTime = utility.GetNowTimestamp()
	groupOwner.IsGroupOwner = 1
	groupOwner.IsRobot = 1
	groupOwner.IsOnline = 2
	m := randInt(5, 8)
	if m < 5 {
		m = 6
	}
	groupOwner.Account = common.RandSeq(sess, lineId, agencyId, m)
	groupOwner.Password = utility.NewPasswordEncrypt(groupOwner.Account, common.RandPassword(8))
	// 创建群主
	UserBo.SaveUser(sess, groupOwner)
	// 组装群数据
	room.GameType = gameType
	room.GamePlay = gamePlay
	room.MaxMoney = maMoney
	room.MinMoney = miMoney
	room.RedNum = redNum
	if redMinNum < 2 {
		redMinNum = 2
	}
	room.RedMinNum = redMinNum
	room.GameTime = gameTime
	room.Royalty = mRoyalty
	room.Odds = modds
	room.Status = 1
	room.RoomType = roomType
	room.FreeFromDeath = freeFromDeath
	room.RobotId = groupOwner.Id
	room.RobotSendPacketTime = 0
	room.RobotSendPacket = 2
	room.RobotGrabPacket = 1
	room.RobotSendPacketTime = 1
	room.RoomNo = roomNo + 1
	room.ControlKill = 1 // 默认启用控杀
	room.CreateTime = utility.GetNowTimestamp()

	err = RoomBo.SaveRoom(sess, room)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = lineId
	log.AgencyId = agencyId
	log.CreatorIp = ip
	if isAdmin == 1 {
		log.Remark = fmt.Sprintf("代理%v添加了群%v,群号%v:", creator, room.RoomName, room.RoomNo)
	} else {
		log.Remark = fmt.Sprintf("代理%v添加了群%v,群号%v:", creator, room.RoomName, room.RoomNo)
	}
	RedLogBo.AddLog(sess, log)
	return nil
}

// 根据id查询房间信息
func (RoomService) QueryRoomOne(id int) (*structs.Room, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	has, room := RoomBo.GetOne(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}

	// 获取群主信息
	has, roomOwner := UserBo.GetOne(sess, room.RobotId)
	if has && roomOwner.IsGroupOwner == 1 {
		room.RobotAccount = roomOwner.Account
	}

	return room, nil
}

// 修改房间
//func (RoomService) EditRoom(lineId string, agencyId string, id int, roomName string, gameType int, gamePlay int, maxMoney string,
//	minMoney string, odds string, redNum int, royalty string, gameTime int, roomSort int, redMinNum int) error {
func (RoomService) EditRoom(id int, roomName string, gameType int, gamePlay int, maxMoney string,
	minMoney string, odds string, redNum int, royalty string, gameTime int, roomSort int, redMinNum int, roomType int,
	freeFromDeath int, robotSendPacket int, robotSendPacketTime int, robotGrabPacket int, controlKill int, creatorId int, creator, ip string, isAdmin int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	if robotSendPacket == 1 && robotSendPacketTime == 0 {
		return &validate.Err{Code: code.ROBOT_SEND_PACKET_TIME_CAN_NOT_BE_EMPTY}
	}
	// 判断房间是否存在
	has, room := RoomBo.GetOne(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	oldRoomName := room.RoomName
	var remark string
	var changeType bool
	if gameType != 0 && gameType != room.GameType {
		room.GameType = gameType
		switch gameType {
		case 1:
			remark += "将群游戏类型改为牛牛;"
		case 2:
			remark += "将群游戏类型改为扫雷;"
		}
		changeType = true
	}
	if gamePlay != 0 && room.GamePlay != gamePlay {
		room.GamePlay = gamePlay
	}
	if roomName == "" {
		if changeType {
			switch gameType {
			case 1:
				switch gamePlay {
				case 1:
					remark += "将群游戏玩法改为经典玩法;"
				case 2:
					remark += "将群游戏玩法改为平倍玩法;"
				case 3:
					remark += "将群游戏玩法改为超倍玩法;"
				}
			case 2:
				switch gamePlay {
				case 1:
					remark += "将群游戏玩法改为固定赔率;"
				case 2:
					remark += "将群游戏玩法改为不固定赔率;"
				}
			}
		}
		if gameType == 1 {
			roomName := "牛牛"
			// 1经典牛牛 2平倍牛牛 3超倍牛牛
			if gamePlay == 1 {
				roomName += fmt.Sprintf("经典群%d-%d分钟", roomSort, gameTime)
				if roomName != room.RoomName {
					room.RoomName = roomName
					remark += fmt.Sprintf("将群名称改为%v;", roomName)
				}
			} else if gamePlay == 2 {
				roomName += fmt.Sprintf("平倍群%d-%d分钟", roomSort, gameTime)
				if roomName != room.RoomName {
					room.RoomName = roomName
					remark += fmt.Sprintf("将群名称改为%v;", roomName)
				}
			} else {
				roomName += fmt.Sprintf("超倍群%d-%d分钟", roomSort, gameTime)
				if roomName != room.RoomName {
					room.RoomName = roomName
					remark += fmt.Sprintf("将群名称改为%v;", roomName)
				}
			}
		} else {
			roomName := fmt.Sprintf("扫雷红包群%d-%d分钟", roomSort, gameTime)
			if roomName != room.RoomName {
				room.RoomName = roomName
				remark += fmt.Sprintf("将群名称改为了%v;", roomName)
			}
		}
	} else {
		room.RoomName = roomName
		remark += fmt.Sprintf("将群名称改为了%v;", roomName)
	}

	if maxMoney != "" {
		maMoney, _ := strconv.ParseFloat(maxMoney, 64)
		if maMoney != room.MaxMoney {
			room.MaxMoney = maMoney
			remark += fmt.Sprintf("将群最大发包金额改为%v;", maMoney)
		}
	}
	if minMoney != "" {
		miMoney, _ := strconv.ParseFloat(minMoney, 64)
		if miMoney != room.MinMoney {
			room.MinMoney = miMoney
			remark += fmt.Sprintf("将群最小发包金额改为%v;", miMoney)
		}
	}
	if redNum != 0 && redNum != room.RedNum {
		room.RedNum = redNum
		remark += fmt.Sprintf("将群最大发包个数改为%v;", redNum)
	}
	if redMinNum != 0 && redMinNum != room.RedMinNum {
		room.RedMinNum = redMinNum
		remark += fmt.Sprintf("将群最小发包个数改为%v;", redMinNum)
	}
	if gameTime != 0 && room.GameTime != gameTime {
		room.GameTime = gameTime
		remark += fmt.Sprintf("将群游戏时间改为%v分钟;", gameTime)
	}
	if royalty != "" {
		mRoyalty, _ := strconv.ParseFloat(royalty, 64)
		if mRoyalty != room.Royalty {
			room.Royalty = mRoyalty
			remark += fmt.Sprintf("将群抽水比例改为%v;", mRoyalty)
		}
	}
	if odds != "" {
		modds, _ := strconv.ParseFloat(odds, 64)
		if modds != room.Odds {
			room.Odds = modds
			remark += fmt.Sprintf("将群赔率改为%v;", modds)
		}
	}
	if roomSort != room.RoomSort {
		room.RoomSort = roomSort
		remark += fmt.Sprintf("将群排序改为%v;", roomSort)
	}
	if roomType != 0 && roomType != room.RoomType {
		room.RoomType = roomType
		switch roomType {
		case 1:
			remark += fmt.Sprintf("将群排序改为%v;", roomSort)
		case 2:
			remark += fmt.Sprintf("将群排序改为%v;", roomSort)
		}
	}
	if freeFromDeath != 0 && freeFromDeath != room.FreeFromDeath {
		room.FreeFromDeath = freeFromDeath
		switch freeFromDeath {
		case 1:
			remark += "将群免死号改为启用;"
		case 2:
			remark += "将群免死号改为停用;"
		}
	}
	var robotSendHas bool
	if robotSendPacket != 0 && robotSendPacket != room.RobotSendPacket {
		room.RobotSendPacket = robotSendPacket
		switch robotSendPacket {
		case 1:
			remark += "将群自动发包改为启用;"
		case 2:
			remark += "将群自动发包改为停用;"
		}

		robotSendHas = true
	}
	if robotSendPacketTime != 0 && robotSendPacketTime != room.RobotSendPacketTime {
		room.RobotSendPacketTime = robotSendPacketTime
		remark += fmt.Sprintf("将群自动发包时间改为%v分钟;", robotSendPacketTime)
	}
	if robotGrabPacket != 0 && room.RobotGrabPacket != robotGrabPacket {
		room.RobotGrabPacket = robotGrabPacket
		switch robotGrabPacket {
		case 1:
			remark += "将群自动抢包改为启用;"
		case 2:
			remark += "将群自动抢包改为停用;"
		}
	}
	if controlKill != 0 && room.ControlKill != controlKill {
		room.ControlKill = controlKill
		switch controlKill {
		case 1:
			remark += "将群控杀改为启用;"
		case 2:
			remark += "将群控杀改为停用;"
		}
	}

	err := RoomBo.ModifyRoom(sess, room)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 发包时间缓存
	if robotSendHas {
		redisKey := "lastSendPacket"
		if robotSendPacket == model.ROBOT_SEND_PACKET_OFF { // 开启自动发红包
			room.LastTime = utility.GetNowTimestamp()
			str, _ := json.Marshal(room)
			conf.GetRedis().Get().Do("HSet", redisKey, fmt.Sprint(room.Id)) // 设置最后一个红包时间redis redis生存时间无限
			conf.GetRedis().Get().Do("expire", string(str))
		} else { // 关闭自动发红包
			conf.GetRedis().Get().Do("HDel", redisKey, fmt.Sprint(room.Id))
		}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = room.LineId
	log.AgencyId = room.AgencyId
	log.CreatorIp = ip
	if isAdmin == 1 {
		log.Remark = fmt.Sprintf("超管%v修改了群%v,群号%v信息:", creator, oldRoomName, room.RoomNo) + remark
	} else {
		log.Remark = fmt.Sprintf("代理%v修改了群%v,群号%v信息:", creator, oldRoomName, room.RoomNo) + remark
	}

	RedLogBo.AddLog(sess, log)
	return nil
}

// 修改房间状态
func (RoomService) EditRoomStatus(id int, status int, creatorId int, creator, ip string, isAdmin int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断房间是否存在
	has, room := RoomBo.GetOne(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	var remark string
	if room.Status != status {
		room.Status = status
		switch status {
		case 1:
			remark += fmt.Sprintf("将群状态修改为启用️;")
		case 2:
			remark += fmt.Sprintf("将群状态修改为停用;")
		}
	}
	err := RoomBo.ModifyRoomStatus(sess, room)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = room.LineId
	log.AgencyId = room.AgencyId
	log.CreatorIp = ip
	if isAdmin == 1 {
		log.Remark = fmt.Sprintf("超管%v修改了群%v,群号%v信息:", creator, room.RoomName, room.RoomNo) + remark
	} else {
		log.Remark = fmt.Sprintf("代理%v修改了群%v,群号%v信息:", creator, room.RoomName, room.RoomNo) + remark
	}
	RedLogBo.AddLog(sess, log)
	return nil
}

// 删除房间信息
func (RoomService) DelRoom(id int, creatorId int, creator, ip string, isAdmin int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断房间是否存在
	has, room := RoomBo.GetOne(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	// 判断是否存在未结算的注单
	count, err := RedPacketLogBo.GetOrdersByRoomId(sess, id, room.LineId, room.AgencyId)
	if err == nil && count > 0 {
		return &validate.Err{Code: code.GROUP_HAS_ORDER_CAN_NOT_BE_DEL}
	}
	room.DeleteTime = utility.GetNowTimestamp()
	err = RoomBo.DelRoom(sess, room)
	if err != nil {
		return &validate.Err{Code: code.DELETE_FAILED}
	}
	// 写入操作日志
	log := new(structs.RedLog)
	log.Creator = creator
	log.CreatorId = creatorId
	log.CreateTime = utility.GetNowTimestamp()
	log.LogType = model.LOG_TYPE_OTHER
	log.LineId = room.LineId
	log.AgencyId = room.AgencyId
	log.CreatorIp = ip
	if isAdmin == 1 {
		log.Remark = fmt.Sprintf("超管%v删除了群%v,群号%v", creator, room.RoomName, room.RoomNo)
	} else {
		log.Remark = fmt.Sprintf("代理%v删除了群%v,群号%v", creator, room.RoomName, room.RoomNo)
	}
	RedLogBo.AddLog(sess, log)
	return nil
}

func (RoomService) RoomCode(lineId, agencyId string, gameType int) ([]*structs.RoomCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	roomCode, err := RoomBo.RoomCode(sess, lineId, agencyId, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return roomCode, err
}
