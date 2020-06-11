package structs

type Agency struct {
	Id             int    `xorm:"'id' PK autoincr" json:"id"`             // 主键id
	LineId         string `xorm:"line_id" json:"lineId"`                  // 线路id
	AgencyId       string `xorm:"agency_id" json:"agencyId"`              // 超管id
	Account        string `xorm:"account" json:"account"`                 // 账号
	Password       string `xorm:"password" json:"password"`               // 密码
	IsOnline       int    `xorm:"is_online" json:"isOnline"`              // 在线状态 1在线 2离线
	IsAdmin        int    `xorm:"is_admin" json:"isAdmin"`                // 是否是超管
	Status         int    `xorm:"status" json:"status"`                   // 状态 1正常 2停用
	CreateTime     int    `xorm:"create_time" json:"createTime"`          // 创建时间
	DeleteTime     int    `xorm:"delete_time" json:"deleteTime"`          // 软删除时间
	EditTime       int    `xorm:"edit_time" json:"editTime"`              // 修改时间
	WhiteIpAddress string `xorm:"white_ip_address" json:"whiteIpAddress"` // ip白名单
}

func (Agency) TableName() string {
	return TABLE_AGENCY
}

type AgencyReq struct {
	Id             int    `json:"id"`             // 主键id
	LineId         string `json:"lineId"`         // 线路id
	AgencyId       string `json:"agencyId"`       // 超管id
	Account        string `json:"account"`        // 账号
	IsOnline       int    `json:"isOnline"`       // 在线状态 1在线 2离线
	IsAdmin        int    `json:"isAdmin"`        // 是否是超管
	Status         int    `json:"status"`         // 状态 1正常 2停用
	CreateTime     int    `json:"createTime"`     // 创建时间
	DeleteTime     int    `json:"deleteTime"`     // 软删除时间
	EditTime       int    `json:"editTime"`       // 修改时间
	WhiteIpAddress string `json:"whiteIpAddress"` // ip白名单
}
