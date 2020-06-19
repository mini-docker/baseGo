package structs

// 用户基本信息
type User struct {
	Id            int     `xorm:"'id' PK autoincr" json:"id"`           // 主键id
	LineId        string  `xorm:"line_id" json:"lineId"`                // 线路id
	AgencyId      string  `xorm:"agency_id" json:"agencyId"`            // 超管ID
	Account       string  `xorm:"account" json:"account"`               // 账号
	Password      string  `xorm:"password" json:"password"`             // 密码
	IsOnline      int     `xorm:"is_online" json:"isOnline"`            // 在线状态 1在线 2离线
	Balance       float64 `xorm:"balance" json:"balance"`               // 会员余额
	Ip            string  `xorm:"ip" json:"ip"`                         // 注册ip
	Status        int     `xorm:"status" json:"status"`                 // 状态 1正常 2停用
	CreateTime    int     `xorm:"create_time" json:"createTime"`        // 创建时间
	DeleteTime    int     `xorm:"delete_time" json:"deleteTime"`        // 软删除时间
	EditTime      int     `xorm:"edit_time" json:"editTime"`            // 修改时间
	Capital       float64 `xorm:"capital" json:"capital"`               // 红包押金
	LastLoginIp   string  `xorm:"last_login_ip" json:"lastLoginIp"`     //上次登陆ip
	LastLoginTime int     `xorm:"last_login_time" json:"lastLoginTime"` //上次登陆时间
	IsRobot       int     `xorm:"is_robot" json:"isRobot"`              // 是否是机器人
	IsGroupOwner  int     `xorm:"is_group_owner" json:"isGroupOwner"`   // 是否是群主
}

func (*User) TableName() string {
	return TABLE_USER
}

// 用户基本信息
type UserResp struct {
	Id               int     `json:"id"`               // 主键id
	LineId           string  `json:"lineId"`           // 线路id
	AgencyId         string  `json:"agencyId"`         // 超管ID
	Account          string  `json:"account"`          // 账号
	Balance          float64 `json:"balance"`          // 会员余额
	CreateTime       int     `json:"createTime"`       // 创建时间
	EditTime         int     `json:"editTime"`         // 修改时间
	Capital          float64 `json:"capital"`          // 红包押金
	AvailableBalance float64 `json:"availableBalance"` // 可用金额
	LastLoginIp      string  `json:"lastLoginIp"`      // 上次登陆ip
	LastLoginTime    int     `json:"lastLoginTime"`    // 上次登陆时间
}

// 用户列表返回
type UserListResp struct {
	Id         int     `xorm:"'id' PK autoincr" json:"id"` // 主键id
	LineId     string  `xorm:"line_id" json:"lineId"`     // 线路id
	AgencyId   string  `xorm:"agency_id" json:"agencyId"`                   // 超管ID
	Account    string  `xorm:"account" json:"account"`     // 账号
	Balance    float64 `json:"balance"`                    // 会员余额
	CreateTime int     `json:"createTime"`                 // 创建时间
	Ip         string  `xorm:"ip" json:"ip"`               // 注册ip
	Status     int     `xorm:"status" json:"status"`       // 状态 1正常 2停用
	IsOnline   int     `xorm:"is_online" json:"isOnline"`  // 在线状态 1在线 2离线
}

type RobotAccounts struct {
	Account string `json:"account"`
}
