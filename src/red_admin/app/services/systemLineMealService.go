package services

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/model/bo"
	"baseGo/src/model/code"
	"baseGo/src/model/structs"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/conf"
	"fmt"
)

type SystemLineMealService struct{}

var (
	SystemLineMealBo = new(bo.SystemLineMealBo)
)

// 查询套餐列表
func (SystemLineMealService) QuerySystemLineMealList(page, pageSize int) (*structs.PageListResp, error) {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取全部套餐信息
	count, lineMeals, err := SystemLineMealBo.QuerySystemLineMealList(sess, page, pageSize)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	pageResp := new(structs.PageListResp)
	pageResp.Data = lineMeals
	pageResp.Count = count
	return pageResp, nil
}

// 添加套餐
func (SystemLineMealService) AddSystemLineMeal(mealName string, slRoyalty, NnRoyalty float64) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	systemLineMeal := new(structs.SystemLineMeal)
	// 添加套餐
	systemLineMeal.MealName = mealName
	systemLineMeal.NnRoyalty = NnRoyalty
	systemLineMeal.SlRoyalty = slRoyalty
	systemLineMeal.CreateTime = utility.GetNowTimestamp()
	fmt.Println(systemLineMeal, "systemLineMeal")
	_, err := SystemLineMealBo.AddLineMeal(sess, systemLineMeal)
	if err != nil {
		return &validate.Err{Code: code.INSET_ERROR}
	}
	return nil
}

// 根据id查询套餐信息
func (SystemLineMealService) QueryLineMealOne(id int) (*structs.SystemLineMeal, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	LineMeal, has, _ := SystemLineMealBo.QueryLineMealById(sess, id)
	if !has {
		return nil, &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	return LineMeal, nil
}

// 修改套餐
func (SystemLineMealService) EditSystemLineMeal(id int, mealName string, slRoyalty, NnRoyalty float64) error {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 判断套餐是否存在
	systemLineMeal, has, _ := SystemLineMealBo.QueryLineMealById(sess, id)
	if !has {
		return &validate.Err{Code: code.DATA_NOT_EXIST}
	}
	systemLineMeal.MealName = mealName
	systemLineMeal.NnRoyalty = NnRoyalty
	systemLineMeal.SlRoyalty = slRoyalty
	systemLineMeal.EditTime = utility.GetNowTimestamp()
	err := SystemLineMealBo.EditLineMeal(sess, systemLineMeal)
	if err != nil {
		return &validate.Err{Code: code.UPDATE_FAILED}
	}
	return nil
}

// 查询lineMeal枚举
func (SystemLineMealService) QueryAllLineMealCode() ([]*structs.SystemLineMealCode, error) {
	sess := conf.GetXormSession()
	defer sess.Close()
	// 获取全部线路信息
	lines, err := SystemLineMealBo.QuerySystemLineMealCode(sess)
	if err != nil {
		return nil, &validate.Err{Code: code.QUERY_FAILED}
	}
	return lines, nil
}
