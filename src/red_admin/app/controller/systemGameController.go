package controller

import (
	"baseGo/src/red_admin/app/controller/common"
	"baseGo/src/red_admin/app/server"
	"baseGo/src/red_admin/app/services"
)

// 游戏控制器
type SystemGameController struct {
	AddGameReq struct {
		GameName string `json:"gameName" valid:"Must;ErrorCode(3033)` // 游戏名称
		GameType int    `json:"gameType" valid:"Must;ErrorCode(3032)` // 游戏类型  1 红包
		Status   int    `json:"status" valid:"Must;ErrorCode(3028)`   // 状态 1正常 2停用
	}

	QueryGameReq struct {
		PageIndex int `json:"pageIndex"` // 页码
		PageSize  int `json:"pageSize"`  // 每页条数
	}

	QueryGameOneReq struct {
		Id int `json:"id"` // 游戏id
	}

	EditGameReq struct {
		Id       int    `json:"id" valid:"Must;ErrorCode(3031)`       // 游戏id
		GameName string `json:"gameName" valid:"Must;ErrorCode(3033)` // 游戏名称
		GameType int    `json:"gameType" valid:"Must;ErrorCode(3032)` // 游戏类型  1 红包
		Status   int    `json:"status" valid:"Must;ErrorCode(3028)`   // 状态 1正常 2停用
	}

	EditGameStatusReq struct {
		Id     int `json:"id"`     // 游戏id
		Status int `json:"status"` // 状态
	}
}

var (
	SystemGameService = new(services.SystemGameService)
)

// 查询全部游戏
func (m SystemGameController) QueryGameList(ctx server.Context) error {
	req := &m.QueryGameReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := SystemGameService.QueryGameList(req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加游戏
func (m SystemGameController) AddGame(ctx server.Context) error {
	req := &m.AddGameReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemGameService.AddGame(req.GameName, req.GameType, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询游戏信息
func (m SystemGameController) QueryGameOne(ctx server.Context) error {
	req := &m.QueryGameOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	game, err := SystemGameService.QueryGameOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, game)
}

// 修改游戏
func (m SystemGameController) EditGame(ctx server.Context) error {
	req := &m.EditGameReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemGameService.EditGame(req.Id, req.GameName, req.GameType, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改代理状态
func (m SystemGameController) EditGameStatus(ctx server.Context) error {
	req := &m.EditGameStatusReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemGameService.EditGameStatus(req.Id, req.Status)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
