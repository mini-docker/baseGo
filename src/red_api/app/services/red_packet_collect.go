package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_api/app/middleware/validate"
	"baseGo/src/red_api/app/server"
	"baseGo/src/red_api/conf"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

const QUERY_TIME = 5 * 60

type RedPacketCollectService struct{}

func (ms RedPacketCollectService) Collect(lineId string) ([]*structs.RedPacketCollect, error) {

	rpcs := make([]*structs.RedPacketCollect, 0)
	lastTime, err := server.QueryRedPacketCollectLastTime(lineId)
	if err != nil {
		if err != redis.Nil {
			return nil, &validate.Err{Code: code.NO_DATA}
		}
		goto OUT
	}

	if lastT, _ := strconv.Atoi(lastTime); (utility.GetNowTimestamp() - lastT) < QUERY_TIME {
		return nil, &validate.Err{Code: code.TIME_TOO_SHORT}
	}
OUT:

	rpcsStrs, err := server.QueryRedPacketCollect(lineId)

	if err != nil {
		if err == redis.Nil {
			return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
		}
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}

	for _, rpcsStr := range rpcsStrs {
		rpc := new(structs.RedPacketCollect)
		err := json.Unmarshal([]byte(rpcsStr), rpc)
		if err != nil {
			return nil, &validate.Err{Code: code.QUERY_FAILED}
		}
		rpcs = append(rpcs, rpc)
	}
	server.SetRedPacketCollectLastTime(lineId, fmt.Sprint(utility.GetNowTimestamp()))
	return rpcs, nil
}

func (ms RedPacketCollectService) CollectByDate(start, end int, lineId string) ([]*structs.RedPacketCollect, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	sess.Begin()
	rpcs, err := RedPacketCollectBo.QueryPacketCollects(sess, lineId, start, end)
	if err != nil {
		sess.Rollback()
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	err = RedPacketCollectBo.UpdatePacketCollects(sess, lineId, start, end)

	if err != nil {
		sess.Rollback()
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	sess.Commit()

	return rpcs, nil
}
