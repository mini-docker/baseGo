package services

import (
	"fecho/common"
	"fecho/golog"
	"fecho/utility"
	"fecho/xorm"
	"model/bo"
	"model/code"
	"model/structs"
	"red_admin/app/middleware/validate"
	"red_admin/conf"
)

type OrderStatistical struct{}

var (
	OrderStatisticalBo = new(bo.OrderStatistical)
)

// 获取盈利分析数据
func (*OrderStatistical) QueryOrderStatistical(lineId, agencyId string, startTime, endTime, gameType int64) (*structs.OrderStatisticalResp, error) {
	// 连接数据库
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取总计数据
	totalData, err := getTotalData(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, err
	}
	// 获取趋势分析
	orderSeries, err := getOrderSeries(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, err
	}

	// 获取盈利列表
	orderStatistical, err := getOrderStatistical(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, err
	}
	resp := new(structs.OrderStatisticalResp)
	resp.TotalData = totalData
	resp.OrderStatistical = orderStatistical
	resp.OrderSeries = orderSeries
	return resp, nil
}

// 获取总计数据
func getTotalData(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) (*structs.TotalData, error) {
	// 获取有效游戏人数
	validMember, err := OrderStatisticalBo.GetValidMember(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.GET_VALIDMEMBER_NUM_ERROR}
	}
	// 获取统计数据
	statistical, err := OrderStatisticalBo.QueryStatistical(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.GET_STATISTICAL_DATA_ERROR}
	}
	// 初始化结构体
	totalData := new(structs.TotalData)
	totalData.ValidMember = validMember
	for _, v := range statistical {
		totalData.ValidBet = common.DecimalSum(totalData.ValidBet, v.ValidBet)
		totalData.RedNum += v.RedNum
		totalData.OrderNum += v.OrderNum
		totalData.RoyaltyMoney = common.DecimalSum(totalData.RoyaltyMoney, v.RoyaltyMoney)
		totalData.FreeDeathWin = common.DecimalSum(totalData.FreeDeathWin, v.FreeDeathWin)
		totalData.RobotWin = common.DecimalSum(totalData.RobotWin, v.RobotWin)
		totalData.TotalWin = common.DecimalSum(totalData.TotalWin, v.TotalWin)
	}
	return totalData, nil
}

// 获取趋势分析数据
func getOrderSeries(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) (*structs.OrderSeries, error) {
	// 获取有效游戏人数
	validMember, err := OrderStatisticalBo.GetSeriesValidMember(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.GET_VALIDMEMBER_NUM_ERROR}
	}
	// 获取统计数据
	statistical, err := OrderStatisticalBo.QuerySeriesStatistical(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.GET_STATISTICAL_DATA_ERROR}
	}
	golog.Info("statistical","getOrderSeries","统计数据：",len(statistical),&statistical)
	// 拆分日期
	everyDays := utility.GetEveryDay(int(startTime), int(endTime))
	// 初始化数据
	orderSeries := new(structs.OrderSeries)
	orderSeries.Data = everyDays
	orderSeries.NameData = []string{"有效人数","有效打码","总局数","总注单","总抽水","免死号盈利","机器人盈利","总盈利"}
	series := make([]*structs.Series, 0)
	validMemberSeries := &structs.Series{
		Name:  "有效人数",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	validBetSeries := &structs.Series{
		Name:  "有效打码",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	redNumSeries := &structs.Series{
		Name:  "总局数",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	orderNumSeries := &structs.Series{
		Name:  "总注单",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	royaltySeries := &structs.Series{
		Name:  "总抽水",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	freeDeathSeries := &structs.Series{
		Name:  "免死号盈利",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	robotSeries := &structs.Series{
		Name:  "机器人盈利",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	totalWinSeries := &structs.Series{
		Name:  "总盈利",
		Type:  "line",
		Data:  make([]float64, len(everyDays)),
	}
	for i, v := range everyDays {
		for _, m := range validMember {
			if v == m.GameTime {
				validMemberSeries.Data[i] = float64(m.Total)
			}
		}
		for _, s := range statistical {
			if v == s.StatisticalDate {
				validBetSeries.Data[i] = common.DecimalSum(validBetSeries.Data[i],s.ValidBet)
				redNumSeries.Data[i] = common.DecimalSum(redNumSeries.Data[i],float64(s.RedNum))
				orderNumSeries.Data[i] = common.DecimalSum(orderNumSeries.Data[i],float64(s.OrderNum))
				royaltySeries.Data[i] = common.DecimalSum(royaltySeries.Data[i],s.RoyaltyMoney)
				freeDeathSeries.Data[i] = common.DecimalSum(freeDeathSeries.Data[i],s.FreeDeathWin)
				robotSeries.Data[i] = common.DecimalSum(robotSeries.Data[i],s.RobotWin)
				totalWinSeries.Data[i] = common.DecimalSum(totalWinSeries.Data[i],s.TotalWin)
			}
		}
	}
	series = append(series, validMemberSeries)
	series = append(series, validBetSeries)
	series = append(series, redNumSeries)
	series = append(series, orderNumSeries)
	series = append(series, royaltySeries)
	series = append(series, freeDeathSeries)
	series = append(series, robotSeries)
	series = append(series, totalWinSeries)
	orderSeries.Series = series
	return orderSeries, nil
}

// 获取盈利列表
func getOrderStatistical(sess *xorm.Session, lineId, agencyId string, startTime, endTime, gameType int64) ([]*structs.RedOrderStatistical, error) {
	// 获取全部数据
	statistical, err := OrderStatisticalBo.QuerySiteStatistic(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.GET_STATISTICAL_DATA_ERROR}
	}
	// 获取所有站点有效游戏人数
	validMenbers,err := OrderStatisticalBo.QuerySiteValidMember(sess, lineId, agencyId, startTime, endTime, gameType)
	if err != nil {
		return nil, &validate.Err{Code: code.GET_VALIDMEMBER_NUM_ERROR}
	}
	for _,v := range statistical {
		for _,s := range validMenbers {
			if v.AgencyId == s.AgencyId {
				v.ValidMember = s.Total
			}
		}
	}
	return statistical, nil
}
