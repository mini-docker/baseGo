package services

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/conf"
	"fmt"

	"strings"
)

type NewSiteService struct {
}

func (NewSiteService) NewSite(lineId string, sites string) error {
	// 批量初始化站点数据
	siteList := strings.Split(sites, ",")
	if len(siteList) > 0 {
		// 连接数据库
		sess := conf.GetXormSession()
		defer sess.Close()
		// 获取当前最大房间号
		roomNo, err := RoomBo.GetMaxRoomNo(sess)
		if err != nil {
			golog.Info("NewSiteService", "NewSite", "获取当前最大房间号失败")
			return nil
		}
		for _, v := range siteList {
			// 添加站点信息
			site := new(structs.RedPacketSite)
			site.LineId = lineId
			site.AgencyId = v
			site.SiteName = v
			site.Status = 1
			site.DeleteTime = 0

			_, err := RedPacketSiteBo.AddSite(sess, site)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加站点失败:", err, "站点id:", v)
				continue
			}

			// 添加代理
			account := fmt.Sprintf("%v001", v)
			password := "123456"
			agency := new(structs.Agency)
			agency.Account = account
			agency.Password = utility.NewPasswordEncrypt(account, password)
			agency.LineId = lineId
			agency.Status = 1
			agency.AgencyId = v
			agency.IsAdmin = 2
			agency.IsOnline = 2
			agency.WhiteIpAddress = "127.0.0.1,127.0.0.2"
			agency.CreateTime = utility.GetNowTimestamp()
			_, err = SystemAgencyBo.AddAgency(sess, agency)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加站点代理失败:", err, "站点id:", v)
				continue
			}

			// 生成群主
			for i := 0; i < 5; i++ {
				groupOwner := new(structs.User)
				groupOwner.AgencyId = v
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
				groupOwner.Account = common.RandSeq(sess, lineId, v, m)
				groupOwner.Password = utility.NewPasswordEncrypt(groupOwner.Account, common.RandPassword(8))
				// 创建群主
				UserBo.SaveUser(sess, groupOwner)
			}
			// 获取群主列表
			groupOwnerList, _ := UserBo.GetGroupListByAgenncyId(sess, lineId, v)
			if len(groupOwnerList) != 5 {
				golog.Info("NewSiteService", "NewSite", "群主数量不足:", nil, "站点id:", v)
				continue
			}
			// 初始化群信息
			// 组装牛牛经典玩法群数据
			nnjdroom := new(structs.Room)
			nnjdroom.LineId = lineId
			nnjdroom.AgencyId = v
			nnjdroom.RoomName = fmt.Sprintf("牛牛经典群0-%d分钟", 1)
			nnjdroom.GameType = model.NIUNIU_RED_ENVELOPE
			nnjdroom.GamePlay = model.NIUNIU_CLASSIC_PLAY
			nnjdroom.MaxMoney = 100
			nnjdroom.MinMoney = 10
			nnjdroom.RedNum = 10
			nnjdroom.RedMinNum = 2
			nnjdroom.GameTime = 1
			nnjdroom.Royalty = 5
			nnjdroom.Odds = 1.6
			nnjdroom.Status = 1
			nnjdroom.RoomType = 1
			nnjdroom.FreeFromDeath = 1
			nnjdroom.RobotSendPacketTime = 0
			nnjdroom.RobotSendPacket = 2
			nnjdroom.RobotGrabPacket = 2
			nnjdroom.RobotSendPacketTime = 1
			nnjdroom.RoomNo = roomNo + 1
			nnjdroom.CreateTime = utility.GetNowTimestamp()
			nnjdroom.RobotId = groupOwnerList[0].Id
			err = RoomBo.SaveRoom(sess, nnjdroom)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加牛牛经典房间失败:", err, "站点id:", v)
				continue
			}
			roomNo++

			// 组装牛牛平倍玩法群数据
			nnpbroom := new(structs.Room)
			nnpbroom.LineId = lineId
			nnpbroom.AgencyId = v
			nnpbroom.RoomName = fmt.Sprintf("牛牛平倍群0-%d分钟", 1)
			nnpbroom.GameType = model.NIUNIU_RED_ENVELOPE
			nnpbroom.GamePlay = model.NIUNIU_FLAT_PLAY
			nnpbroom.MaxMoney = 100
			nnpbroom.MinMoney = 10
			nnpbroom.RedNum = 10
			nnpbroom.RedMinNum = 2
			nnpbroom.GameTime = 1
			nnpbroom.Royalty = 5
			nnpbroom.Odds = 1.6
			nnpbroom.Status = 1
			nnpbroom.RoomType = 1
			nnpbroom.FreeFromDeath = 1
			nnpbroom.RobotSendPacketTime = 0
			nnpbroom.RobotSendPacket = 2
			nnpbroom.RobotGrabPacket = 2
			nnpbroom.RobotSendPacketTime = 1
			nnpbroom.RoomNo = roomNo + 1
			nnpbroom.CreateTime = utility.GetNowTimestamp()
			nnpbroom.RobotId = groupOwnerList[1].Id
			err = RoomBo.SaveRoom(sess, nnpbroom)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加牛牛平倍房间失败:", err, "站点id:", v)
				continue
			}

			roomNo++

			// 组装牛牛超倍玩法群数据
			nncbroom := new(structs.Room)
			nncbroom.LineId = lineId
			nncbroom.AgencyId = v
			nncbroom.RoomName = fmt.Sprintf("牛牛超倍群0-%d分钟", 1)
			nncbroom.GameType = model.NIUNIU_RED_ENVELOPE
			nncbroom.GamePlay = model.NIUNIU_SUPER_PLAY
			nncbroom.MaxMoney = 100
			nncbroom.MinMoney = 10
			nncbroom.RedNum = 10
			nncbroom.RedMinNum = 2
			nncbroom.GameTime = 1
			nncbroom.Royalty = 5
			nncbroom.Odds = 1.6
			nncbroom.Status = 1
			nncbroom.RoomType = 1
			nncbroom.FreeFromDeath = 1
			nncbroom.RobotSendPacketTime = 0
			nncbroom.RobotSendPacket = 2
			nncbroom.RobotGrabPacket = 2
			nncbroom.RobotSendPacketTime = 1
			nncbroom.RoomNo = roomNo + 1
			nncbroom.CreateTime = utility.GetNowTimestamp()
			nncbroom.RobotId = groupOwnerList[2].Id
			err = RoomBo.SaveRoom(sess, nncbroom)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加牛牛超倍房间失败:", err, "站点id:", v)
				continue
			}

			roomNo++

			// 组装扫雷固定赔率玩法群数据
			slroom := new(structs.Room)
			slroom.LineId = lineId
			slroom.AgencyId = v
			slroom.RoomName = fmt.Sprintf("扫雷红包群0-%d分钟", 1)
			slroom.GameType = model.MINESWEEPER_RED_PACKET
			slroom.GamePlay = model.MINESWEEPER_FIXED_ODDS
			slroom.MaxMoney = 100
			slroom.MinMoney = 10
			slroom.RedNum = 10
			slroom.RedMinNum = 2
			slroom.GameTime = 1
			slroom.Royalty = 5
			slroom.Odds = 1.6
			slroom.Status = 1
			slroom.RoomType = 1
			slroom.FreeFromDeath = 1
			slroom.RobotSendPacketTime = 0
			slroom.RobotSendPacket = 2
			slroom.RobotGrabPacket = 2
			slroom.RobotSendPacketTime = 1
			slroom.RoomNo = roomNo + 1
			slroom.CreateTime = utility.GetNowTimestamp()
			slroom.RobotId = groupOwnerList[3].Id
			err = RoomBo.SaveRoom(sess, slroom)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加扫雷固定赔率房间失败:", err, "站点id:", v)
				continue
			}

			roomNo++

			// 组装扫雷不固定赔率玩法群数据
			slsroom := new(structs.Room)
			slsroom.LineId = lineId
			slsroom.AgencyId = v
			slsroom.RoomName = fmt.Sprintf("扫雷红包群0-%d分钟", 1)
			slsroom.GameType = model.MINESWEEPER_RED_PACKET
			slsroom.GamePlay = model.MINESWEEPER_UNFIXED_ODDS
			slsroom.MaxMoney = 100
			slsroom.MinMoney = 10
			slsroom.RedNum = 10
			slsroom.RedMinNum = 2
			slsroom.GameTime = 1
			slsroom.Royalty = 5
			slsroom.Odds = 1.5
			slsroom.Status = 1
			slsroom.RoomType = 1
			slsroom.FreeFromDeath = 1
			slsroom.RobotSendPacketTime = 0
			slsroom.RobotSendPacket = 2
			slsroom.RobotGrabPacket = 2
			slsroom.RobotSendPacketTime = 1
			slsroom.RoomNo = roomNo + 1
			slsroom.CreateTime = utility.GetNowTimestamp()
			slsroom.RobotId = groupOwnerList[4].Id
			err = RoomBo.SaveRoom(sess, slsroom)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加扫雷不固定赔率房间失败:", err, "站点id:", v)
				continue
			}

			// 初始化banner图
			activePicture := new(structs.ActivePicture)
			activePicture.LineId = lineId
			activePicture.AgencyId = v
			activePicture.ActiveName = "first"
			activePicture.StartTime = 1579798800
			activePicture.EndTime = 1611421200
			activePicture.Status = 1
			activePicture.Picture = "/redimgs/activePictures/20200118/2020/1/18/1579327938834987557.png"
			// 保存主表
			err = ActivePictureBo.AddActive(sess, activePicture)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加活动图片失败:", err, "站点id:", v)
				continue
			}

			// 初始化banner图
			activePicture1 := new(structs.ActivePicture)
			activePicture1.LineId = lineId
			activePicture1.AgencyId = v
			activePicture1.ActiveName = "second"
			activePicture1.StartTime = 1579798800
			activePicture1.EndTime = 1611421200
			activePicture1.Status = 1
			activePicture1.Picture = "/redimgs/activePictures/20200118/2020/1/18/1579327962940142269.png"
			// 保存主表
			err = ActivePictureBo.AddActive(sess, activePicture1)
			if err != nil {
				golog.Info("NewSiteService", "NewSite", "添加活动图片1失败:", err, "站点id:", v)
				continue
			}
		}
	}
	return nil
}

func (NewSiteService) NewRobot(lineId string) error {
	// 连接数据库
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取没有机器人的站点
	//sites, err := RedPacketSiteBo.GetNoRobotsSite(sess, lineId)
	// 获取全部站点
	sites, err := RedPacketSiteBo.SiteCode(sess, lineId)
	if err != nil {
		return err
	}
	for _, v := range sites {
		robots := make([]*structs.User, 0)
		for i := 0; i < 100; i++ {
			robot := new(structs.User)
			m := randInt(6, 8)
			if m < 6 {
				m = 6
			}
			robot.Account = common.RandSeq(sess, lineId, v.AgencyId, m)
			robot.AgencyId = v.AgencyId
			robot.LineId = lineId
			robot.Status = 1
			robot.CreateTime = utility.GetNowTimestamp()
			robot.IsGroupOwner = 2
			robot.IsRobot = 1
			robot.IsOnline = 2
			robot.Password = utility.NewPasswordEncrypt(robot.Account, common.RandPassword(8))
			robots = append(robots, robot)
		}
		UserBo.SaveRobots(sess, robots)
	}
	return nil
}

func (NewSiteService) InitStatistical(dayTime int, nextTime int) error {
	// 连接数据库
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取全部线路信息
	lines, err := new(bo.SystemLineBo).QuerySystemLineIdList(sess)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取线路信息失败", err)
		return err
	}
	// 获取当天日期
	dateStr := utility.Format10(utility.GetUnixTime(int64(dayTime)))
	for _, v := range lines {
		// 获取线路下面站点信息
		sites, err := new(bo.RedPacketSite).SiteCode(sess, v.LineId)
		if err != nil {
			golog.Error("middleware", "orderStatistical", "获取站点信息失败", err, "线路id：", v.LineId)
			continue
		}
		// 统计站点当日盈利分析
		for _, site := range sites {
			nnStatistical := middleware.Statistical(sess, v.LineId, site.AgencyId, model.NIUNIU_RED_ENVELOPE, dayTime, nextTime)
			nnStatistical.StatisticalDate = dateStr
			slStatistical := middleware.Statistical(sess, v.LineId, site.AgencyId, model.MINESWEEPER_RED_PACKET, dayTime, nextTime)
			slStatistical.StatisticalDate = dateStr
			// 判断当天是否已经添加统计数据
			has := OrderStatisticalBo.CheckStatistical(sess, v.LineId, site.AgencyId, dateStr)
			if has {
				// 数据存在，更新
				OrderStatisticalBo.UpdateStatistic(sess, nnStatistical)
				OrderStatisticalBo.UpdateStatistic(sess, slStatistical)
			} else {
				// 不存在，添加
				OrderStatisticalBo.InsertStatistic(sess, nnStatistical)
				OrderStatisticalBo.InsertStatistic(sess, slStatistical)
			}
		}
	}
	return nil
}
