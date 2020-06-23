package controllers

import (
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services"
)

var NewSiteService = new(services.NewSiteService)

type NewSiteController struct {
	NewSiteReq struct {
		LineId string `json:"lineId"` // 线路id
		Sites  string `json:"sites"`  // 站点id","分割数组
	}

	StatisticalReq struct {
		DayTime  int `json:"dayTime"`
		NextTime int `json:"nextTime"`
	}
}

// 初始化站点
func (m NewSiteController) NewSite(ctx server.Context) error {
	req := &m.NewSiteReq
	if err := ctx.Validate(req); err != nil {
		return err
	}
	// 初始化数据
	err := NewSiteService.NewSite(req.LineId, req.Sites)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 初始化机器人
func (m NewSiteController) NewRobot(ctx server.Context) error {
	req := &m.NewSiteReq
	if err := ctx.Validate(req); err != nil {
		return err
	}
	// 初始化数据
	err := NewSiteService.NewRobot(req.LineId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 初始化统计数据
func (m NewSiteController) InitStatistical(ctx server.Context) error {
	req := &m.StatisticalReq
	if err := ctx.Validate(req); err != nil {
		return err
	}
	err := NewSiteService.InitStatistical(req.DayTime, req.NextTime)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
