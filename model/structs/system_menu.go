package structs

type SystemMenu struct {
	Id         int    `xorm:"'id' PK autoincr" json:"id"`    // 主键
	ParentId   int    `xorm:"parent_id" json:"parentId""`    // 父id
	Name       string `xorm:"name" json:"name""`             // 菜单名称
	Route      string `xorm:"route" json:"route"`            // 菜单路由
	Icon       string `xorm:"icon" json:"icon"`              // 菜单图标
	Level      int    `xorm:"level" json:"level"`            // 菜单级别 1 一级 2 二级 3 三级
	Status     int    `xorm:"status" json:"status"`          // 状态 1 启用 2 停用
	IsShow     int    `xorm:"is_show" json:"isShow"`         // 是否可见 1可见 2不可见
	Sort       int    `xorm:"sort" json:"sort"`              // 排列序号
	CreateTime int    `xorm:"create_time" json:"createTime"` // 创建时间
	DeleteTime int    `xorm:"delete_time" json:"deleteTime"` // 软删除时间
	UpdateTime int    `xorm:"update_time" json:"updateTime"` // 更新时间
}

func (*SystemMenu) TableName() string {
	return TABLE_SYSTEM_MENU
}

type SystemMenuResp struct {
	Id       int               `json:"id"`        // 主键
	ParentId int               `json:"parentId""` // 父id
	Name     string            `json:"name""`     // 菜单名称
	Route    string            `json:"route"`     // 菜单路由
	Icon     string            `json:"icon"`      // 菜单图标
	Level    int               `json:"level"`     // 菜单级别 1 一级 2 二级 3 三级
	Status   int               `json:"status"`    // 状态 1 启用 2 停用
	IsShow   int               `json:"isShow"`    // 是否可见 1可见 2不可见
	Sort     int               `json:"sort"`      // 排列序号
	Children []*SystemMenuResp `json:"children"`  // 子级
}

type SystemMenuCode struct {
	Id   int    `xorm:"'id' PK autoincr" json:"id"` // 主键
	Name string `xorm:"name" json:"name""`          // 菜单名称
}

func (*SystemMenuCode) TableName() string {
	return TABLE_SYSTEM_MENU
}

type SystemMenuPermission struct {
	Id int `xorm:"'id' PK autoincr" json:"id"` // 主键
}
