package po

type HttpResult struct {
	//    [{
	//        "success":"",	 --->是否成功
	//        "code":"",	 --->状态码
	//        "message":"",    --->消息框
	//        "version":"", --->版本信息
	//        "data":[{"key":"value"}],	--->返回数据
	//        "is_native" true // 是否跳转h5
	//    }]
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//data的正确格式: 3&4.0.1&black&xxhdpi&zh&aaa_a&53496937f5bc02807d3934ac9fcf5ac9
type HttpHeaderData struct {
	Terminal   int // 3 android端 4 ios端
	Version    string
	Theme      string
	Resolution string
	Locale     string
	Sign       string
}

//{"success":false,"message":"账号或密码错误！","username":null,"password":null,"isOpenCaptcha":true,"propMessages":{}}
//{"success":false,"message":null,"username":null,"password":null,"isOpenCaptcha":true,"propMessages":{"captcha":"验证码不正确!"}}
//{"success":true,"message":null,"username":null,"password":null,"isOpenCaptcha":false,"propMessages":{}}

type LoginBean struct {
	Success       bool             `json:"success"`
	Message       string           `json:"message"`
	Username      string           `json:"username"`
	Password      string           `json:"password"`
	IsOpenCaptcha bool             `json:"isOpenCaptcha"`
	PropMessages  PropMessagesBean `json:"propMessages"`
}

type PropMessagesBean struct {
	Captcha  string `json:"captcha,omitempty"`
	Location string `json:"location,omitempty"`
	GbToken  string `json:"gbToken,omitempty"`
}

type HttpRes struct {
	Data string `json:"data"`
}
