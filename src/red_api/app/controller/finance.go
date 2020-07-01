package controller

import (
	"baseGo/src/model"
	"baseGo/src/red_api/app/controller/common"
	"baseGo/src/red_api/app/server"
	"baseGo/src/red_api/app/services"
)

type FinanceController struct {
	TransferredReq struct {
		// 账号
		Account string `json:"account"`
		// 转入金额
		Amount float64 `json:"Amount"`
	}
}

func (ac FinanceController) TransferredIn(ctx server.Context) error {
	req := &ac.TransferredReq
	err := ctx.Validate(req)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	lineId := ctx.Get(model.LineId).(string)
	agencyId := ctx.Get(model.AgencyId).(string)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	orderNo, err := new(services.FinanceService).TransferredIn(lineId, agencyId, req.Account, req.Amount)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, orderNo)
}

func (ac FinanceController) TransferredOut(ctx server.Context) error {
	req := &ac.TransferredReq
	err := ctx.Validate(req)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	lineId := ctx.Get(model.LineId).(string)
	agencyId := ctx.Get(model.AgencyId).(string)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	balance, err := new(services.FinanceService).TransferredOut(lineId, agencyId, req.Account, req.Amount)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, balance)
}
