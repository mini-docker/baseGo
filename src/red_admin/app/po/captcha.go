package po

type CaptchaBean struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code,omitempty"`
}

type PuzzleCaptchaBean struct {
	Y      int    `json:"y,omitempty"`
	Id     string `json:"id,omitempty"`
	Imgx   int    `json:"imgx,omitempty"`
	Imgy   int    `json:"imgy,omitempty"`
	Small  string `json:"small,omitempty"`
	SImgx  int    `json:"sImgx,omitempty"`
	SImgy  int    `json:"sImgy,omitempty"`
	Normal string `json:"normal,omitempty"`
}

// 拼图验证返回
type PuzzleCaptchaResp struct {
	Y      int    `json:"y"`      // 裁剪的小图相对左上角的y轴坐标
	Array  string `json:"array"`  // 验证码图片混淆规律
	Imgx   int    `json:"imgx"`   // 验证码图片宽度
	Imgy   int    `json:"imgy"`   // 验证码图片高度
	Small  string `json:"small"`  // 裁剪的小图片
	SImgx  int    `json:"sImgx"`  // 验证码小图片宽度
	SImgy  int    `json:"sImgy"`  // 验证码小图片高度
	Normal string `json:"normal"` // 验证码混淆后的图片
}

// 拼图验证返回
type CaptchaReq struct {
	X        int    `json:"x"`                                                   // 图片位置
	CodeType int    `json:"codeType" valid:"Must;Min(1);Max(2);ErrorCode(5311)"` // 验证码类型 1 注册 2 登陆
	Account  string `json:"account" valid:"Min(5);Max(12);ErrorCode(8201)"`      // 账号错误
	//Account  string `json:"Account" valid:"Must;Min(5);Max(12);ErrorCode(8201)"` // 账号错误
}
