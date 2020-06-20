package model

// 标准session信息
type AgencySession struct {
	SessionId    string      `json:"sessionId"`    //sessionId,登陆后返回整个信息给前端,所以把sessionId放到这里面
	User         *AgencyUser `json:"user"`         // 用户信息
	TimeOut      int         `json:"time_out"`     // 超时时间
	IsKeepOnline bool        `json:"isKeepOnline"` // 保持登录
	IsAdmin      int         `json:"isAdmin"`      // 是否是超级管理员
}

type AgencyUser struct {
	Id       int    `json:"id"`
	Account  string `json:"account"`  // 登陆账号
	LineId   string `json:"lineId"`   // 线路id
	AgencyId string `json:"agencyId"` // 代理id
}

func (c *AgencySession) Info() interface{} {
	return c
}

func (c *AgencySession) Uid() int {
	return c.User.Id
}

func (c *AgencySession) Account() string {
	return c.User.Account
}

func (c *AgencySession) LineId() string {
	return c.User.LineId
}

func (c *AgencySession) AgencyId() string {
	return c.User.AgencyId
}
