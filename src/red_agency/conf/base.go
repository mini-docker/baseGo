package conf

type ReqMember struct {
	SiteId   string  `json:"siteId,omitempty"`
	AgentId  string  `json:"agentId,omitempty"`
	Username string  `json:"username,omitempty"`
	Balance  float64 `json:"balance,omitempty"`
	Currency string  `json:"currency,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
}
type RequestBody struct {
	Cmd       string     `json:"cmd"`
	RequestId string     `json:"requestId,omitempty"`
	SubVt     string     `json:"vt,omitempty"` //彩票类型
	Data      string     `json:"data,omitempty"`
	Member    *ReqMember `json:"member,omitempty"`
}

type Packet struct {
	Buff     []byte `json:"-"`
	Platform string `json:"p"`
	Data     string `json:"data"`
	Key      string `json:"key"`
}

type Member struct {
	SiteId        string  `json:"siteId,omitempty"`
	AgentId       string  `json:"agentId,omitempty"`
	Username      string  `json:"username,omitempty"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency,omitempty"`
	Amount        float64 `json:"amount,omitempty"`        //操作的额度
	BeforeBalance float64 `json:"beforeBalance,omitempty"` //操作之前额度
}

type ResponseBody struct {
	Code      int     `json:"code"`
	Msg       string  `json:"msg"`
	RequestId string  `json:"requestId,omitempty"`
	TicketId  string  `json:"ticketId,omitempty"`
	Member    *Member `json:"member,omitempty"`
}



