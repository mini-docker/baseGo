package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
	"fmt"
)

type SystemLineRoyalty struct {
}

func (SystemLineRoyalty) QueryLineRoyaltyList(sess *xorm.Session, startTime, endTime int) ([]*structs.LineRoyaltyListResp, error) {
	if startTime != 0 {
		sess.Where("red_start_time >= ?", startTime)
	}
	if endTime != 0 {
		sess.Where("red_start_time <= ?", endTime)
	}
	sess.Where("status < ?", model.RED_RESULT_INVALID)
	winRoyaltyList := make([]*structs.LineRoyaltyListResp, 0)
	err := sess.
		Select("line_id as lineId,sum(case when game_type = 1 and royalty_money > 0 and is_robot = 0 then royalty_money else 0.00 end) as nnRoyalty," +
			"sum(case when game_type = 2 and royalty_money > 0 and is_robot = 0 then royalty_money else 0.00 end) as slRoyalty").
		GroupBy("lineId").Find(&winRoyaltyList)
	fmt.Println(sess.LastSQL())
	return winRoyaltyList, err
}

func (SystemLineRoyalty) QueryLineAgencyRoyaltyList(sess *xorm.Session, startTime, endTime int, lineId string) ([]*structs.AgencyRoyaltyListResp, error) {
	if startTime != 0 {
		sess.Where("red_start_time >= ?", startTime)
	}
	if endTime != 0 {
		sess.Where("red_start_time <= ?", endTime)
	}
	if lineId != "" {
		sess.Where("line_id = ?", lineId)
	}
	sess.Where("status < ?", model.RED_RESULT_INVALID)
	winRoyaltyList := make([]*structs.AgencyRoyaltyListResp, 0)
	err := sess.
		Select("line_id as lineId,agency_id as agencyId,sum(case when game_type = 1 and royalty_money > 0 and is_robot = 0 then royalty_money else 0.00 end) as nnRoyalty," +
			"sum(case when game_type = 2 and royalty_money > 0 and is_robot = 0 then royalty_money else 0.00 end) as slRoyalty").
		GroupBy("lineId,agencyId").Find(&winRoyaltyList)
	return winRoyaltyList, err
}
