package controller

import (
	"baseGo/src/model"
	"baseGo/src/red_api/app/controller/common"
	"baseGo/src/red_api/app/server"
	"baseGo/src/red_api/app/services"
)

type RedPacketCollectController struct {
	ColletcByDateReq struct {
		// 开始时间
		Start int `json:"start"`
		// 结束时间
		End int `json:"end"`
	}
}

func (ac RedPacketCollectController) Collect(ctx server.Context) error {

	lineId := ctx.Get(model.LineId).(string)
	rpcs, err := new(services.RedPacketCollectService).Collect(lineId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, rpcs)
}

func (ac RedPacketCollectController) CollectByDate(ctx server.Context) error {
	req := &ac.ColletcByDateReq
	err := ctx.Validate(req)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	lineId := ctx.Get(model.LineId).(string)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	rpcs, err := new(services.RedPacketCollectService).CollectByDate(req.Start, req.End, lineId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	return common.HttpResultJson(ctx, rpcs)
}
