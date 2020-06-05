package structs

type SystemAdmin struct {
	Id            int    `xorm:"'id' PK autoincr" json:"id"`           // 主键
	Account       string `xorm:"account" json:"account"`               // 账号
	Password      string `xorm:"password" json:"password"`             // 密码
	RoleId        int    `xorm:"role_id" json:"roleId"`                // 角色id
	IsOnline      int    `xorm:"is_online" json:"isOnline"`            // 在线状态 1 在线  2 离线
	LastIp        string `xorm:"last_ip" json:"lastIp"`                // 上次登陆ip
	LastLoginTime int    `xorm:"last_login_time" json:"lastLoginTime"` // 上次登陆时间
	CreateTime    int    `xorm:"create_time" json:"createTime"`        // 创建时间
	DeleteTime    int    `xorm:"delete_time" json:"deleteTime"`        // 是否删除或删除时间
}

func (*SystemAdmin) TableName() string {
	return TABLE_SYSTEM_ADMIN
}

type SystemAdminReq struct {
	Id            int    `xorm:"'id' PK autoincr" json:"id"`           // 主键
	Account       string `xorm:"account" json:"account"`               // 账号
	RoleId        int    `xorm:"role_id" json:"roleId"`                // 角色id
	RoleName      string `xorm:"-" json:"roleName"`                    // 角色名称
	IsOnline      int    `xorm:"is_online" json:"isOnline"`            // 在线状态 1 在线  2 离线
	LastIp        string `xorm:"last_ip" json:"lastIp"`                // 上次登陆ip
	LastLoginTime int    `xorm:"last_login_time" json:"lastLoginTime"` // 上次登陆时间
	CreateTime    int    `xorm:"create_time" json:"createTime"`        // 创建时间
	DeleteTime    int    `xorm:"delete_time" json:"deleteTime"`        // 是否删除或删除时间
}
