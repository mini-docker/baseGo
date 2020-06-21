package services

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/conf"
)

type RedLogService struct{}

func (RedLogService) GetRedLogList(lineId, agencyId string, logType, startTime, endTime, page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	count, logList, err := RedLogBo.GetRedLogList(sess, lineId, agencyId, logType, startTime, endTime, page, pageSize)
	if err != nil {
		golog.Error("RedPacketService", "CreateRedPacket", "err:", err)
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = logList
	pageResp.Count = count
	return pageResp, nil
}
