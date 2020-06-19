package services

import (
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/conf"
)

type SystemGameService struct{}

var (
	GameBo = new(bo.SystemGameBo)
)

// 查询游戏列表
func (SystemGameService) QueryGameList(page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部游戏信息
	count, games, err := GameBo.QuerySystemGameList(sess, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = games
	pageResp.Count = count
	return pageResp, nil
}

// 添加游戏
func (SystemGameService) AddGame(gameName string, gameType, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	game := new(structs.Game)
	// 添加游戏
	game.GameName = gameName
	game.GameType = gameType
	game.Status = status
	_, err := GameBo.AddGame(sess, game)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 根据id查询游戏信息
func (SystemGameService) QueryGameOne(id int) (*structs.Game, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	game, has, _ := GameBo.QueryGameById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	return game, nil
}

// 修改游戏
func (SystemGameService) EditGame(id int, gameName string, gameType, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断游戏是否存在
	game, has, _ := GameBo.QueryGameById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	game.GameName = gameName
	game.GameType = gameType
	game.Status = status
	err := GameBo.EditGame(sess, game)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 修改游戏状态
func (SystemGameService) EditGameStatus(id int, status int) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断游戏是否存在
	game, has, _ := GameBo.QueryGameById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	game.Status = status
	err := GameBo.EditGame(sess, game)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}
