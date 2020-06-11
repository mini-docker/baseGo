package structs

type SystemRole struct {
	Id         int    `xorm:"'id' PK autoincr" json:"id"`    // 主键
	RoleName   string `xorm:"role_name" json:"roleName"`     // 角色名称
	IsDefault  int    `xorm:"is_default" json:"isDefault"`   // 是否是默认角色 不允许 删除和禁用 1 是 2 不是
	Status     int    `xorm:"status" json:"status"`          // 角色状态 1 为启用 2为禁用
	Remark     string `xorm:"remark" json:"remark"`          // 备注
	EditTime   int    `xorm:"edit_time" json:"editTime"`     // 修改时间
	CreateTime int    `xorm:"create_time" json:"createTime"` // 创建时间
	DeleteTime int    `xorm:"delete_time" json:"deleteTime"` // 是否删除或删除时间
}

func (*SystemRole) TableName() string {
	return TABLE_SYSTEM_ROLE
}

type SystemRoleCode struct {
	Id       int    `json:"id"`       // 主键
	RoleName string `json:"roleName"` // 角色名称
}
