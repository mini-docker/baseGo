package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type SystemGameBo struct{}

// 返回所有游戏列表
func (*SystemGameBo) QuerySystemGameList(sess *xorm.Session, page, pageSize int) (int64, []*structs.Game, error) {
	rows := make([]*structs.Game, 0)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&rows)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}

// 添加游戏
func (*SystemGameBo) AddGame(sess *xorm.Session, game *structs.Game) (int64, error) {
	return sess.Insert(game)
}

// 修改游戏信息
func (*SystemGameBo) EditGame(sess *xorm.Session, game *structs.Game) error {
	_, err := sess.Table(new(structs.Game).TableName()).
		ID(game.Id).
		Cols("game_name", "game_type", "status").
		Update(game)
	return err
}

// 修改游戏状态
func (*SystemGameBo) EditGameStatus(sess *xorm.Session, game *structs.Game) error {
	_, err := sess.Table(new(structs.Game).TableName()).
		ID(game.Id).
		Cols("status").
		Update(game)
	return err
}

// 根据id查询单个游戏
func (*SystemGameBo) QueryGameById(sess *xorm.Session, id int) (*structs.Game, bool, error) {
	game := new(structs.Game)
	has, err := sess.Where("id = ?", id).Get(game)
	return game, has, err
}
