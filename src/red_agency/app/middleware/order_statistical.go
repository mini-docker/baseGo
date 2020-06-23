package middleware

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/conf"
	"fmt"

	"baseGo/src/github.com/jasonlvhit/gocron"
)

var (
	systemLineBo       = new(bo.SystemLineBo)
	siteBo             = new(bo.RedPacketSite)
	OrderStatisticalBo = new(bo.OrderStatistical)
)

// 注单盈利分析
func InitOrderStatistical() {
	s := gocron.NewScheduler()
	s.Every(10).Minutes().Do(orderStatistical)
	s.Every(1).Days().At("1:00").Do(lastDayOrderStatistical)
	<-s.Start()
}

// 注单盈利分析
func orderStatistical() {
	// 连接数据库
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取全部线路信息
	lines, err := systemLineBo.QuerySystemLineIdList(sess)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取线路信息失败", err)
		return
	}
	// 获取当天0点时间戳
	nightTime := utility.GetNightTimestamp(0)
	// 获取当天日期
	dateStr := utility.Format10(utility.GetNowTime())
	for _, v := range lines {
		// 获取线路下面站点信息
		sites, err := siteBo.SiteCode(sess, v.LineId)
		if err != nil {
			golog.Error("middleware", "orderStatistical", "获取站点信息失败", err, "线路id：", v.LineId)
			continue
		}
		// 统计站点当日盈利分析
		for _, site := range sites {
			nnStatistical := Statistical(sess, v.LineId, site.AgencyId, model.NIUNIU_RED_ENVELOPE, nightTime, 0)
			nnStatistical.StatisticalDate = dateStr
			slStatistical := Statistical(sess, v.LineId, site.AgencyId, model.MINESWEEPER_RED_PACKET, nightTime, 0)
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
}

// 前一天注单盈利分析
func lastDayOrderStatistical() {
	fmt.Println("统计前一天盈利分析数据")
	// 连接数据库
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取全部线路信息
	lines, err := systemLineBo.QuerySystemLineIdList(sess)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取线路信息失败", err)
		return
	}
	// 获取当天0点时间戳
	nightTime := utility.GetNightTimestamp(0)
	// 获取前一天0点时间戳
	lastNightTime := utility.GetNightTimestamp(-1)
	// 获取前一天日期
	dateStr := utility.Format10(utility.GetUnixTime(int64(lastNightTime)))
	for _, v := range lines {
		// 获取线路下面站点信息
		sites, err := siteBo.SiteCode(sess, v.LineId)
		if err != nil {
			golog.Error("middleware", "orderStatistical", "获取站点信息失败", err, "线路id：", v.LineId)
			continue
		}
		// 统计站点当日盈利分析
		for _, site := range sites {
			nnStatistical := Statistical(sess, v.LineId, site.AgencyId, model.NIUNIU_RED_ENVELOPE, lastNightTime, nightTime)
			nnStatistical.StatisticalDate = dateStr
			slStatistical := Statistical(sess, v.LineId, site.AgencyId, model.MINESWEEPER_RED_PACKET, lastNightTime, nightTime)
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
}

func Statistical(sess *xorm.Session, lineId, agencyId string, gameType, lastNightTime, nightTime int) *structs.RedOrderStatistical {
	// 统计有效打码金额
	validBet, err := OrderStatisticalBo.GetValidBet(sess, lineId, agencyId, gameType, lastNightTime, nightTime)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取站点有效投注信息失败", err, "线路id：", lineId, "站点id：", agencyId)
	}

	// 统计总局数
	redNum, err := OrderStatisticalBo.GetRedNum(sess, lineId, agencyId, gameType, lastNightTime, nightTime)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取站点总局数信息失败", err, "线路id：", lineId, "站点id：", agencyId)
	}

	// 统计总注单
	orderNum, err := OrderStatisticalBo.GetOrderNum(sess, lineId, agencyId, gameType, lastNightTime, nightTime)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取站点总注单信息失败", err, "线路id：", lineId, "站点id：", agencyId)
	}

	// 统计抽水金额
	royaltyMoney, err := OrderStatisticalBo.GetRoyaltyMoney(sess, lineId, agencyId, gameType, lastNightTime, nightTime)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取站点总抽水信息失败", err, "线路id：", lineId, "站点id：", agencyId)
	}

	// 统计免死号盈利金额
	freeDeathWin, err := OrderStatisticalBo.GetFreeDeathWin(sess, lineId, agencyId, gameType, lastNightTime, nightTime)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取站点免死号盈利信息失败", err, "线路id：", lineId, "站点id：", agencyId)
	}

	// 统计机器人盈利
	robotWin, err := OrderStatisticalBo.GetRobotWin(sess, lineId, agencyId, gameType, lastNightTime, nightTime)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取站点机器人盈利信息失败", err, "线路id：", lineId, "站点id：", agencyId)
	}

	// 统计总盈利
	totalWin, err := OrderStatisticalBo.GetTotalWin(sess, lineId, agencyId, gameType, lastNightTime, nightTime)
	if err != nil {
		golog.Error("middleware", "orderStatistical", "获取站点总盈利信息失败", err, "线路id：", lineId, "站点id：", agencyId)
	}

	// 初始化统计结构体
	orderStatistical := new(structs.RedOrderStatistical)
	orderStatistical.LineId = lineId
	orderStatistical.AgencyId = agencyId
	orderStatistical.GameType = model.NIUNIU_RED_ENVELOPE
	orderStatistical.ValidBet = validBet
	orderStatistical.RedNum = redNum
	orderStatistical.OrderNum = orderNum
	orderStatistical.RoyaltyMoney = royaltyMoney
	orderStatistical.FreeDeathWin = freeDeathWin
	orderStatistical.RobotWin = robotWin
	orderStatistical.TotalWin = common.DecimalSub(0, totalWin)
	orderStatistical.GameType = gameType
	return orderStatistical
}
