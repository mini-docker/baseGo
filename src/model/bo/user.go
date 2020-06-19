package bo

import (
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm"
	"baseGo/src/model"
	"baseGo/src/model/structs"
	"fmt"
	"strconv"
)

type User struct{}

func (*User) SaveUser(sess *xorm.Session, user *structs.User) error {
	_, err := sess.Insert(user)
	return err
}

func (*User) SaveRobots(sess *xorm.Session, users []*structs.User) error {
	_, err := sess.Insert(users)
	return err
}

func (*User) ExistUser(sess *xorm.Session, account string) (bool, *structs.User) {
	sess.Where("account = ? and delete_time = 0", account)
	user := new(structs.User)
	has, _ := sess.Get(user)
	return has, user
}

func (*User) GetOne(sess *xorm.Session, id int) (bool, *structs.User) {
	sess.Where("id = ?  and delete_time = 0", id)
	user := new(structs.User)
	has, _ := sess.Table(new(structs.User).TableName()).Get(user)
	return has, user
}

func (*User) GetById(sess *xorm.Session, lineId, agencyId string, id int) (bool, *structs.User) {
	sess.Where("line_id = ? ", lineId)
	if agencyId != "" {
		sess.Where("agency_id = ? ", agencyId)
	}
	sess.Where("id = ?  and delete_time = 0", id)
	user := new(structs.User)
	has, _ := sess.Table(new(structs.User).TableName()).Get(user)
	return has, user
}

func (*User) GetOneByAccount(sess *xorm.Session, lineId string, agencyId string, account string) (bool, *structs.User) {
	sess.Where("account = ?  and delete_time  = 0", account)
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ? ", agencyId)
	user := new(structs.User)
	has, _ := sess.Get(user)
	fmt.Println(sess.LastSQL())
	return has, user
}

func (*User) ModifyUser(sess *xorm.Session, user *structs.User) error {
	sess.Where("id = ?  and delete_time = 0", user.Id)
	_, err := sess.Cols("nike_name", "photo", "sex", "birthday", "sign", "address", "open_notice").Update(user)
	return err
}

func (*User) UpdatePassword(sess *xorm.Session, user *structs.User) error {
	sess.Where("id = ? and delete_time = 0", user.Id)
	_, err := sess.Cols("password").Update(user)
	return err
}
func (*User) UpdateCapitalIncr(sess *xorm.Session, userId int, capital float64) error {
	sess.Where("id = ? and delete_time = 0", userId)
	sess.Incr("capital", capital)
	_, err := sess.Cols("capital").Update(&structs.User{})
	return err
}

func (*User) UpdateCapitalDecr(sess *xorm.Session, lineId string, userId int, capital float64) error {
	sess.Where("id = ? and delete_time = 0", userId)
	sess.Decr("capital", capital)
	_, err := sess.Cols("capital").Update(&structs.User{})
	return err
}

// 批量更新会员
func (*User) UpdateCapitalDecrList(sess *xorm.Session, userIds []int, capital float64) error {
	sess.In("id", userIds)
	sess.Where("delete_time = 0")
	sess.Decr("capital", capital)
	_, err := sess.Cols("capital").Update(&structs.User{})
	return err
}

func (*User) FindByIds(sess *xorm.Session, lineId, agencyId string, ids []int) ([]*structs.User, error) {
	sess.In("id ", ids)
	sess.Where("line_id = ? ", lineId)
	sess.Where("agency_id = ? ", agencyId)
	sess.Where("delete_time = 0")
	users := make([]*structs.User, 0)
	err := sess.Table(new(structs.User).TableName()).Find(&users)
	return users, err
}

// 会员余额变更
func (*User) UpdateBalance(sess *xorm.Session, data map[int]float64) (int64, error) {
	ids := make([]int, 0)
	var balanceSql, idStr string
	for k, v := range data {
		balanceSql += " WHEN " + strconv.Itoa(k) + " THEN balance + " + strconv.FormatFloat(v, 'f', -1, 64)
		ids = append(ids, k)
	}
	for k, v := range ids {
		if k == 0 {
			idStr = strconv.Itoa(v)
		} else {
			idStr += "," + strconv.Itoa(v)
		}
	}
	sqlStr := "UPDATE red_" + structs.TABLE_USER + " SET "
	sqlStr += " balance = CASE id" + balanceSql + " END"
	sqlStr += " where id IN (" + idStr + ")"
	res, err := sess.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 会员余额变更
func (*User) UpdateUserBalance(sess *xorm.Session, lineId string, agencyId string, userId int, money float64) error {
	sess.ID(userId)
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("delete_time = 0")
	sess.Decr("balance", money)
	_, err := sess.Cols("balance").Update(&structs.User{})
	return err
}

func (*User) UpdateLoginIp(sess *xorm.Session, user *structs.User) error {
	sess.Where("id = ? and delete_time = 0", user.Id)
	_, err := sess.Cols("last_login_ip", "last_login_time", "is_online").Update(user)
	return err
}
func (*User) UpdateUserBalanceIncr(sess *xorm.Session, lineId, agencyId string, userId int, money, capital float64) error {
	sess.ID(userId)
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("delete_time = 0")
	sess.Incr("balance", money)
	sess.Incr("capital", capital)
	_, err := sess.Cols("balance").Update(&structs.User{})
	return err
}

// 查询会员列表
func (*User) QueryUserList(sess *xorm.Session, lineId string, agencyId string, status, isOnline int, account string, page, pageSize int) (int64, []*structs.UserListResp, error) {
	if status != 0 {
		sess.Where("status = ?", status)
	}
	if isOnline != 0 {
		sess.Where("is_online = ?", isOnline)
	}
	if account != "" {
		sess.Where("account like ?", account+"%")
	}
	if agencyId != "" {
		sess.Where("agency_id = ?", agencyId)
	}
	if lineId != "" {
		sess.Where("line_id = ? ", lineId)
	}
	sess.Where("delete_time = ? and is_robot = ?", model.UNDEL, model.UNDEL)
	users := make([]*structs.UserListResp, 0)
	count, err := sess.Table(new(structs.User).TableName()).Limit(pageSize, (page-1)*pageSize).OrderBy("id desc").FindAndCount(&users)
	return count, users, err
}

// 批量更新会员状态
func (*User) EditUsersStatus(sess *xorm.Session, ids []int, status int) error {
	_, err := sess.Table(new(structs.User).TableName()).In("id", ids).Cols("status").Update(&structs.User{
		Status:   status,
		EditTime: utility.GetNowTimestamp(),
	})
	return err
}

// 更新会员在线状态
func (*User) EditUserOnlineStatu(sess *xorm.Session, id int64, status int) error {
	_, err := sess.Table(new(structs.User).TableName()).ID(id).Cols("is_online").Update(&structs.User{
		IsOnline: status,
		EditTime: utility.GetNowTimestamp(),
	})
	return err
}

// 批量踢线
func (*User) KickUsers(sess *xorm.Session, ids []int) error {
	_, err := sess.Table(new(structs.User).TableName()).In("id", ids).Cols("is_online").Update(&structs.User{
		IsOnline: 2,
		EditTime: utility.GetNowTimestamp(),
	})
	return err
}

// 批量删除
func (*User) DelUserByIds(sess *xorm.Session, ids []string) error {
	sess.In("id ", ids)
	sess.Where("is_group_owner != 1")
	_, err := sess.Table(new(structs.User).TableName()).Update(&structs.User{
		IsOnline:   2,
		DeleteTime: utility.GetNowTimestamp(),
	})
	return err
}

// 查询agencyId下机器人数据
func (*User) GetRobotListByAgenncyId(sess *xorm.Session, lineId, agencyId string) ([]structs.User, error) {
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_YES)
	sess.Where("is_group_owner != ?", model.USER_IS_GROUP_OWNER_YES)
	sess.Where("delete_time = 0 and status = 1") // 未删除且未停用的机器人
	data := make([]structs.User, 0)
	err := sess.Table(new(structs.User).TableName()).Find(&data)
	return data, err
}

// 查询agencyId下群主数据
func (*User) GetGroupListByAgenncyId(sess *xorm.Session, lineId, agencyId string) ([]structs.User, error) {
	sess.Where("line_id = ?", lineId)
	sess.Where("agency_id = ?", agencyId)
	sess.Where("is_robot = ?", model.USER_IS_ROBOT_YES)
	sess.Where("is_group_owner = ?", model.USER_IS_GROUP_OWNER_YES)
	sess.Where("delete_time = 0 and status = 1") // 未删除且未停用的机器人
	data := make([]structs.User, 0)
	err := sess.Table(new(structs.User).TableName()).Find(&data)
	return data, err
}

// 查询所有需退还保证金会员
func (*User) GetAllBackCapital(sess *xorm.Session) ([]*structs.User, error) {
	users := make([]*structs.User, 0)
	err := sess.SQL("select * from red_user where id not in (select user_id from red_order_record where status = 0 and game_type = 1 ) and capital != 0").Find(&users)
	return users, err
}
