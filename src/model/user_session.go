package model

// 标准session信息
type UserSession struct {
	SessionId    string `json:"sessionId"`    //sessionId,登陆后返回整个信息给前端,所以把sessionId放到这里面
	User         *User  `json:"user"`         // 用户信息
	TimeOut      int    `json:"timeOut"`      // 超时时间
	IsKeepOnline bool   `json:"isKeepOnline"` // 保持登录
}

type User struct {
	Id            int    `json:"id"`
	Account       string `json:"account"`       // 登陆账号
	LastIp        string `json:"lastIp"`        // 上次登陆ip
	LastLoginTime int    `json:"lastLoginTime"` // 上次登陆时间
	LineId        string `json:"lineId"`        // 线路id
	AgencyId      string `json:"agencyId"`      // 超管id
}

func (c *UserSession) Info() interface{} {
	return c
}

func (c *UserSession) Uid() int {
	return c.User.Id
}

func (c *UserSession) Account() string {
	return c.User.Account
}

func (c *UserSession) LineId() string {
	return c.User.LineId
}
func (c *UserSession) AgencyId() string {
	return c.User.AgencyId
}
