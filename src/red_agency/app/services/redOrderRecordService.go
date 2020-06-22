package services

import (
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
	"encoding/json"
	"strconv"
)

type RedOrderService struct {
}

// 查询注单列表
func (RedOrderService) QueryRedRecordList(lineId string, agencyId string, startTime, endTime, gameType, status int,
	orderNo, account, redSender string, page, pageSize int, redId, roomId, isRobot int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取注单列表
	count, orders, err := RedPacketLogBo.QueryRedRecordList(sess, lineId, agencyId, startTime, endTime, gameType, status, orderNo, account, redSender, page, pageSize, redId, roomId, isRobot)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	if len(orders) > 0 {
		for k, v := range orders {
			vData := make(map[string]string)
			json.Unmarshal([]byte(v.Extra), &vData)
			if _, ok := vData["adminNum"]; ok {
				orders[k].AdminNum, _ = strconv.Atoi(vData["adminNum"])
			}
			if _, ok := vData["memberNum"]; ok {
				orders[k].MemberNum, _ = strconv.Atoi(vData["memberNum"])
			}
			if _, ok := vData["thunderNum"]; ok {
				orders[k].ThunderNum, _ = strconv.Atoi(vData["thunderNum"])
			}
			if _, ok := vData["odds"]; ok {
				orders[k].Odds, _ = strconv.ParseFloat(vData["odds"], 64)
			}
			if _, ok := vData["memberMine"]; ok {
				orders[k].MemberMine, _ = strconv.Atoi(vData["memberMine"])
			}
		}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = orders
	pageResp.Count = count
	return pageResp, nil
}

// 获取红包领取记录
func (RedOrderService) GetRedInfo(lineId string, agencyId string, redId int) ([]*structs.RedOrderResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取注单列表
	orders, err := RedPacketLogBo.GetRedInfo(sess, lineId, agencyId, redId)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	for k, v := range orders {
		if v.GameType == 1 {
			// 牛牛
			if v.RedSender == v.Account {
				v.InfoType = 4 // 庄家
			}
			if v.IsFreeDeath == 1 {
				v.InfoType = 3 // 免死号
			}
		} else {
			// 扫雷
			if v.IsFreeDeath == 1 {
				v.InfoType = 3 // 免死号
			}
			if v.IsRobot == 1 && v.IsFreeDeath != 1 {
				v.InfoType = 2 // 机器人
			}
			if v.RealMoney < 0 {
				v.InfoType = 1 // 中雷
			}
			if v.RedSender == v.Account {
				orders = append(orders[:k], orders[k+1:]...) // 将发包人注单删除
			}
		}
	}
	return orders, nil
}
