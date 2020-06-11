package model

// 标准session信息
type AdminSession struct {
	SessionId    string     `json:"sessionId"`    //sessionId,登陆后返回整个信息给前端,所以把sessionId放到这里面
	User         *AdminUser `json:"user"`         // 用户信息
	TimeOut      int        `json:"time_out"`     // 超时时间
	IsKeepOnline bool       `json:"isKeepOnline"` // 保持登录
	IsAdmin      int        `json:"isAdmin"`      // 是否是超级管理员
}

type AdminUser struct {
	Id            int    `json:"id"`
	Account       string `json:"account"`       // 登陆账号
	RoleId        int    `json:"roleId"`        // 角色id
	RoleName      string `json:"roleName"`      // 角色名称
	LastIp        string `json:"lastIp"`        // 上次登陆ip
	LastLoginTime int    `json:"lastLoginTime"` // 上次登陆时间
}

func (c *AdminSession) Info() interface{} {
	return c
}

func (c *AdminSession) Uid() int {
	return c.User.Id
}

func (c *AdminSession) Account() string {
	return c.User.Account
}
