package services

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/conf"
	"encoding/json"
	"fmt"
)

type SystemLineService struct{}

var (
	SystemLineBo = new(bo.SystemLineBo)
)

// 查询线路列表
func (SystemLineService) QuerySystemLineList(lineId, lineName string, status, transType int, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部线路信息
	count, lines, err := SystemLineBo.QuerySystemLineList(sess, lineId, lineName, status, transType, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = lines
	pageResp.Count = count
	return pageResp, nil
}

// 添加线路
func (SystemLineService) AddSystemLine(lineId, lineName string, limitCost float64, mealId int, domain string, transType int, apiUrl, md5key, rsaPubKey, rsaPriKey string, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	if transType == 1 && apiUrl == "" {
		return &validate.Err{Code: code.APIURL_CAN_NOT_BE_EMPTY}
	}

	// 判断线路id是否已存在
	_, has, _ := SystemLineBo.QueryLineBylineId(sess, lineId)
	if has {
		return &validate.Err{Code: code.LINE_ID_EXIST}
	}

	systemLine := new(structs.SystemLine)
	// 添加线路
	systemLine.LineId = lineId
	systemLine.LineName = lineName
	systemLine.LimitCost = limitCost
	systemLine.MealId = mealId
	systemLine.Domain = domain
	systemLine.TransType = transType
	systemLine.ApiUrl = apiUrl
	systemLine.Md5key = md5key
	systemLine.RsaPubKey = rsaPubKey
	systemLine.RsaPriKey = rsaPriKey
	systemLine.Status = status
	systemLine.CreateTime = utility.GetNowTimestamp()
	_, err := SystemLineBo.AddLine(sess, systemLine)
	if err != nil {
		fmt.Println(err, "err_systemLine_1")
		return &validate.Err{Code: code.INSET_ERROR}
	}
	lineByte, _ := json.Marshal(systemLine)
	redisClient := conf.GetRedis().Get()
	// 写入redis
	_, err = redisClient.Do("HSet", model.SYSTEM_LINE_REDIS_KEY, systemLine.LineId, string(lineByte))
	if err != nil {
		fmt.Println(err, "err_systemLine_2")
		golog.Error("systemLineService", "EditSystemLine", "写入redis失败", err)
	}
	// 设置超时时间6分钟
	//conf.GetRedis().Expire(model.SYSTEM_LINE_REDIS_KEY, 6*time.Minute)
	return nil
}

// 根据id查询线路信息
func (SystemLineService) QueryLineOne(id int) (*structs.SystemLine, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	Line, has, _ := SystemLineBo.QueryLineById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	return Line, nil
}

// 修改线路
func (SystemLineService) EditSystemLine(id int, lineName string, limitCost float64, mealId int, domain string, transType int, apiUrl, md5key, rsaPubKey, rsaPriKey string, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	if transType == 1 && apiUrl == "" {
		return &validate.Err{Code: code.APIURL_CAN_NOT_BE_EMPTY}
	}
	// 判断线路是否存在
	systemLine, has, _ := SystemLineBo.QueryLineById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	systemLine.LineName = lineName
	systemLine.LimitCost = limitCost
	systemLine.MealId = mealId
	systemLine.Domain = domain
	systemLine.TransType = transType
	systemLine.ApiUrl = apiUrl
	systemLine.Md5key = md5key
	systemLine.RsaPubKey = rsaPubKey
	systemLine.RsaPriKey = rsaPriKey
	systemLine.Status = status
	systemLine.EditTime = utility.GetNowTimestamp()
	err := SystemLineBo.EditLine(sess, systemLine)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入redis
	lineByte, _ := json.Marshal(systemLine)
	_, err = conf.GetRedis().Get().Do("HSet", model.SYSTEM_LINE_REDIS_KEY, systemLine.LineId, string(lineByte))
	if err != nil {
		golog.Error("systemLineService", "EditSystemLine", "写入redis失败", err)
	}
	return nil
}

// 修改线路状态
func (SystemLineService) EditSystemLineStatus(id, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断线路是否存在
	systemLine, has, _ := SystemLineBo.QueryLineById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	systemLine.Status = status
	systemLine.EditTime = utility.GetNowTimestamp()
	err := SystemLineBo.EditLineStatus(sess, systemLine)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	// 写入redis
	_, err = conf.GetRedis().Get().Do("HSet", model.SYSTEM_LINE_REDIS_KEY, systemLine.LineId, systemLine)
	if err != nil {
		golog.Error("systemLineService", "EditSystemLineStatus", "写入redis失败", err)
	}
	return nil
}

// 查询lineId枚举
func (SystemLineService) QueryAllLineId() ([]*structs.SystemLineCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取全部线路信息
	lines, err := SystemLineBo.QuerySystemLineIdList(sess)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return lines, nil
}
