package bo

import (
	"baseGo/src/fecho/xorm"
	"baseGo/src/model/structs"
)

type SystemLineMealBo struct{}

// 返回所有套餐列表
func (*SystemLineMealBo) QuerySystemLineMealList(sess *xorm.Session, page, pageSize int) (int64, []*structs.SystemLineMeal, error) {
	rows := make([]*structs.SystemLineMeal, 0)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&rows)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}

// 添加套餐
func (*SystemLineMealBo) AddLineMeal(sess *xorm.Session, LineMeal *structs.SystemLineMeal) (int64, error) {
	return sess.Insert(LineMeal)
}

// 修改套餐信息
func (*SystemLineMealBo) EditLineMeal(sess *xorm.Session, LineMeal *structs.SystemLineMeal) error {
	_, err := sess.Table(new(structs.SystemLineMeal).TableName()).
		ID(LineMeal.Id).
		Cols("meal_name", "nn_royalty", "sl_royalty", "edit_time").
		Update(LineMeal)
	return err
}

// 根据id查询单个套餐
func (*SystemLineMealBo) QueryLineMealById(sess *xorm.Session, id int) (*structs.SystemLineMeal, bool, error) {
	LineMeal := new(structs.SystemLineMeal)
	has, err := sess.Where("id = ?", id).Get(LineMeal)
	return LineMeal, has, err
}

func (*SystemLineMealBo) QuerySystemLineMealCode(sess *xorm.Session) ([]*structs.SystemLineMealCode, error) {
	rows := make([]*structs.SystemLineMealCode, 0)
	err := sess.Table(new(structs.SystemLineMeal).TableName()).Find(&rows)
	return rows, err
}
