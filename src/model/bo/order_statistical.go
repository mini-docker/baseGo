package bo

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
	"fmt"
)

type OrderStatistical struct{}

// 获取站点有效游戏人数
func (*OrderStatistical) GetValidMember(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) (int64, error) {
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if startTime != 0 {
		sess.Where("receive_time >= ?", startTime)
	}
	if endTime != 0 {
		sess.Where("receive_time < ?", endTime)
	}
	if gameType != 0 {
		sess.Where("game_type = ?", gameType)
	}
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	return sess.Table(new(structs.OrderRecord).TableName()).Distinct("account").Count()
}

// 获取分析有效游戏人数
func (*OrderStatistical) GetSeriesValidMember(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) ([]*structs.ValidMemberCount, error) {
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if startTime != 0 {
		sess.Where("receive_time >= ?", startTime)
	}
	if endTime != 0 {
		sess.Where("receive_time < ?", endTime)
	}
	if gameType != 0 {
		sess.Where("game_type = ?", gameType)
	}
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	validMember := make([]*structs.ValidMemberCount, 0)
	err := sess.Table(new(structs.OrderRecord).TableName()).Select("count(DISTINCT(account)) as total,FROM_UNIXTIME(receive_time,'%Y-%m-%d') as gameTime").GroupBy("FROM_UNIXTIME(receive_time,'%Y-%m-%d')").Find(&validMember)
	return validMember, err
}

// 获取站点有效游戏人数
func (*OrderStatistical) QuerySiteValidMember(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) ([]*structs.SiteValidMemberCount, error) {
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if startTime != 0 {
		sess.Where("receive_time >= ?", startTime)
	}
	if endTime != 0 {
		sess.Where("receive_time < ?", endTime)
	}
	if gameType != 0 {
		sess.Where("game_type = ?", gameType)
	}
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	validMember := make([]*structs.SiteValidMemberCount, 0)
	err := sess.Table(new(structs.OrderRecord).TableName()).Select("count(DISTINCT(account)) as total,agency_id").GroupBy("agency_id").Find(&validMember)
	return validMember, err
}

// 获取站点有效投注金额
func (*OrderStatistical) GetValidBet(sess *xorm.Session, lineId, agencyId string, gameType int, lastNightTime, nightTime int) (float64, error) {
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	if lastNightTime != 0 {
		sess.Where("receive_time >= ?", lastNightTime)
	}
	if nightTime != 0 {
		sess.Where("receive_time < ?", nightTime)
	}
	sess.Where("game_type = ?", gameType)
	return sess.Table(new(structs.OrderRecord).TableName()).Sum(new(structs.OrderRecord), "valid_bet")
}

// 获取站点总局数
func (*OrderStatistical) GetRedNum(sess *xorm.Session, lineId, agencyId string, gameType int, lastNightTime, nightTime int) (int64, error) {
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	if lastNightTime != 0 {
		sess.Where("receive_time >= ?", lastNightTime)
	}
	if nightTime != 0 {
		sess.Where("receive_time < ?", nightTime)
	}
	sess.Where("game_type = ?", gameType)
	count, err := sess.Table(new(structs.OrderRecord).TableName()).Select("count(DISTINCT(red_id))").Count()
	return count, err
}

// 获取站点总注单
func (*OrderStatistical) GetOrderNum(sess *xorm.Session, lineId, agencyId string, gameType int, lastNightTime, nightTime int) (int64, error) {
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	if lastNightTime != 0 {
		sess.Where("receive_time >= ?", lastNightTime)
	}
	if nightTime != 0 {
		sess.Where("receive_time < ?", nightTime)
	}
	sess.Where("game_type = ?", gameType)
	return sess.Table(new(structs.OrderRecord).TableName()).Count()
}

// 获取站点总抽水金额
func (*OrderStatistical) GetRoyaltyMoney(sess *xorm.Session, lineId, agencyId string, gameType int, lastNightTime, nightTime int) (float64, error) {
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	if lastNightTime != 0 {
		sess.Where("receive_time >= ?", lastNightTime)
	}
	if nightTime != 0 {
		sess.Where("receive_time < ?", nightTime)
	}
	sess.Where("game_type = ?", gameType)
	return sess.Table(new(structs.OrderRecord).TableName()).Sum(new(structs.OrderRecord), "royalty_money")
}

// 获取站点免死号盈利
func (*OrderStatistical) GetFreeDeathWin(sess *xorm.Session, lineId, agencyId string, gameType int, lastNightTime, nightTime int) (float64, error) {
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_free_death = ?", model.USER_IS_ROBOT_YES)
	if lastNightTime != 0 {
		sess.Where("receive_time >= ?", lastNightTime)
	}
	if nightTime != 0 {
		sess.Where("receive_time < ?", nightTime)
	}
	sess.Where("game_type = ?", gameType)
	return sess.Table(new(structs.OrderRecord).TableName()).Sum(new(structs.OrderRecord), "real_money")
}

// 获取站点机器人盈利
func (*OrderStatistical) GetRobotWin(sess *xorm.Session, lineId, agencyId string, gameType int, lastNightTime, nightTime int) (float64, error) {
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_YES)
	sess.Where("is_free_death != ?", model.USER_IS_ROBOT_YES)
	if lastNightTime != 0 {
		sess.Where("receive_time >= ?", lastNightTime)
	}
	if nightTime != 0 {
		sess.Where("receive_time < ?", nightTime)
	}
	sess.Where("game_type = ?", gameType)
	return sess.Table(new(structs.OrderRecord).TableName()).Sum(new(structs.OrderRecord), "robot_win")
}

// 获取站点总盈利
func (*OrderStatistical) GetTotalWin(sess *xorm.Session, lineId, agencyId string, gameType int, lastNightTime, nightTime int) (float64, error) {
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ?", agencyId)
	if lastNightTime != 0 {
		sess.Where("receive_time >= ?", lastNightTime)
	}
	if nightTime != 0 {
		sess.Where("receive_time < ?", nightTime)
	}
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_NO)
	sess.Where("game_type = ?", gameType)
	return sess.Table(new(structs.OrderRecord).TableName()).Sum(new(structs.OrderRecord), "real_money")
}

// 验证数据是否存在
func (*OrderStatistical) CheckStatistical(sess *xorm.Session, lineId, agencyId, dateStr string) bool {
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("statistical_date = ?", dateStr)
	total, _ := sess.Table(new(structs.RedOrderStatistical).TableName()).Count()
	if total > 0 {
		return true
	}
	return false
}

// 添加统计信息
func (*OrderStatistical) InsertStatistic(sess *xorm.Session, info *structs.RedOrderStatistical) (int64, error) {
	return sess.Insert(info)
}

// 修改统计数据
func (*OrderStatistical) UpdateStatistic(sess *xorm.Session, info *structs.RedOrderStatistical) (int64, error) {
	sess.Where("line_id = ?", info.LineId)
	sess.Where("agency_id = ?", info.AgencyId)
	sess.Where("statistical_date = ?", info.StatisticalDate)
	sess.Where("game_type = ?", info.GameType)
	return sess.Cols("valid_bet", "red_num", "order_num", "royalty_money", "free_death_win", "robot_win", "total_win").Update(info)
}

// 获取统计总计数据
func (*OrderStatistical) QueryStatistical(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) ([]*structs.RedOrderStatistical, error) {
	var startTimeStr, endTimeStr string
	if startTime != 0 {
		startTimeStr = utility.Format10(utility.GetUnixTime(startTime))
	}
	if endTime != 0 {
		endTimeStr = utility.Format10(utility.GetUnixTime(endTime))
	}
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if startTimeStr != "" {
		sess.Where("statistical_date >= ?", startTimeStr)
	}
	if endTimeStr != "" {
		sess.Where("statistical_date <= ?", endTimeStr)
	}
	if gameType != 0 {
		sess.Where("game_type = ?", gameType)
	}
	data := make([]*structs.RedOrderStatistical, 0)
	err := sess.Find(&data)
	return data, err
}

// 获取每个站点总计数据
func (*OrderStatistical) QuerySiteStatistic(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) ([]*structs.RedOrderStatistical, error) {
	var startTimeStr, endTimeStr string
	if startTime != 0 {
		startTimeStr = utility.Format10(utility.GetUnixTime(startTime))
	}
	if endTime != 0 {
		endTimeStr = utility.Format10(utility.GetUnixTime(endTime))
	}
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if startTimeStr != "" {
		sess.Where("statistical_date >= ?", startTimeStr)
	}
	if endTimeStr != "" {
		sess.Where("statistical_date <= ?", endTimeStr)
	}
	if gameType != 0 {
		sess.Where("game_type = ?", gameType)
	}
	data := make([]*structs.RedOrderStatistical, 0)
	err := sess.Select("line_id,agency_id,sum(valid_bet) as valid_bet,sum(red_num) as red_num,sum(order_num) as order_num,sum(royalty_money) as royalty_money,sum(free_death_win) as free_death_win,sum(robot_win) as robot_win,sum(total_win) as total_win").GroupBy("line_id,agency_id").Find(&data)
	return data, err
}

// 获取分析总计数据
func (*OrderStatistical) QuerySeriesStatistical(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) ([]*structs.RedOrderStatistical, error) {
	var startTimeStr, endTimeStr string
	if startTime != 0 {
		startTimeStr = utility.Format10(utility.GetUnixTime(startTime))
	}
	if endTime != 0 {
		endTimeStr = utility.Format10(utility.GetUnixTime(endTime))
	}
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if startTimeStr != "" {
		sess.Where("statistical_date >= ?", startTimeStr)
	}
	if endTimeStr != "" {
		sess.Where("statistical_date <= ?", endTimeStr)
	}
	if gameType != 0 {
		sess.Where("game_type = ?", gameType)
	}
	data := make([]*structs.RedOrderStatistical, 0)
	err := sess.Select("line_id,agency_id,statistical_date,sum(valid_bet) as valid_bet,sum(red_num) as red_num,sum(order_num) as order_num,sum(royalty_money) as royalty_money,sum(free_death_win) as free_death_win,sum(robot_win) as robot_win,sum(total_win) as total_win").GroupBy("line_id,agency_id,statistical_date").Find(&data)
	fmt.Println(sess.LastSQL())
	return data, err
}
