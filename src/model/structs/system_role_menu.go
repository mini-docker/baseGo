package structs

type SystemRoleMenu struct {
	Id     int `xorm:"'id' PK autoincr" json:"id"` // 主键
	RoleId int `xorm:"role_id" json:"roleId"`      // 角色id
	MenuId int `xorm:"menu_id" json:"menuId"`      // 菜单id
}

func (*SystemRoleMenu) TableName() string {
	return TABLE_SYSTEM_ROLE_MENU
}
