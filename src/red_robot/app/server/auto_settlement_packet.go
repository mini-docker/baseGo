package server

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/model/structs"
	"baseGo/src/red_robot/conf"
	"fmt"

	"github.com/jasonlvhit/gocron"
)

func AutoSettlementPacket() {
	s := gocron.NewScheduler()
	s.Every(20).Seconds().Do(settlementPacket)
	<-s.Start()
}

func settlementPacket() {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取所有未结算扫雷红包
	slRedInfos, err := new(bo.RedPacket).GetAllNeedSettlementPacket(sess, model.MINESWEEPER_RED_PACKET)
	if err != nil {
		golog.Error("AutoSettlementPacket", "settlementPacket", "获取未结算扫雷红包失败：", err)
		return
	}
	// 获取所有未结算牛牛红包
	nnRedInfos, err := new(bo.RedPacket).GetAllNeedSettlementPacket(sess, model.NIUNIU_RED_ENVELOPE)
	if err != nil {
		golog.Error("AutoSettlementPacket", "settlementPacket", "获取未结算牛牛红包失败：", err)
		return
	}
	// 获取所有未结算普通红包
	ptRedInfos, err := new(bo.RedPacket).GetAllNeedSettlementPacket(sess, model.ORDINARY_RED_ENVELOPE)
	if err != nil {
		golog.Error("AutoSettlementPacket", "settlementPacket", "获取未结算普通红包失败：", err)
		return
	}
	redInfos := make([]structs.OrderRecord, 0)
	for _, v := range slRedInfos {
		if v.RedStartTime+v.GameTime*60 < utility.GetNowTimestamp() {
			redInfos = append(redInfos, v)
		}
	}

	for _, v := range nnRedInfos {
		if v.RedStartTime+v.GameTime*60 < utility.GetNowTimestamp() {
			redInfos = append(redInfos, v)
		}
	}

	for _, v := range ptRedInfos {
		if v.RedStartTime+v.GameTime*60*60 < utility.GetNowTimestamp() {
			redInfos = append(redInfos, v)
		}
	}

	var redIds []int
	var i int
	if len(redInfos) > 0 {
		for _, v := range redInfos {
			// 过滤已结算过的注单
			if len(redIds) > 0 {
				for _, redId := range redIds {
					if redId == v.RedId {
						continue
					}
				}
			}
			// 循环结算红包
			_, err = new(RedPlay).RedEnvelopeAmountCalculation(v.LineId, v.AgencyId, v.RedId, v.RoomId)
			if err != nil {
				golog.Error("AutoSettlementPacketUser", "settlementPacket", "结算红包失败：", err, v.Id, v.Account)
				continue
			}
			// 删除红包缓存
			logKey := fmt.Sprintf("%v_%v_redLog", v.RedId, v.RoomId)
			conf.GetRedis().Get().Do("Del", logKey)
			redIds = append(redIds, v.RedId)
			i++
		}
	}
}
