package echo

type (
	//错误返回
	Err struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	//单条数据返回
	Item struct {
		Data interface{} `json:"data"`
	}
	// 集合分页返回
	Pagination struct {
		Data  interface{} `json:"data"`
		Meta  Meta        `json:"meta"`
		Links Links       `json:"links"`
	}
	// 页码数据
	Meta struct {
		Count           int `json:"count"`             //数据总数
		PageCount       int `json:"page_count"`        //页码总数
		CurrentPage     int `json:"current_page"`      //当前页码
		PageSize        int `json:"page_size"`         //每页数据数量
		CurrentPageSize int `json:"current_page_size"` //当前页数据数量
	}
	// 页码链接
	Links struct {
		First   string `json:"first"`   //首页链接
		Last    string `json:"last"`    //末页链接
		Current string `json:"current"` //当前链接
		Prev    string `json:"prev"`    //上一页链接
		Next    string `json:"next"`    //下一页链接
	}
)

func (m *Err) Error() string {
	return m.Msg
}

const (
	ZH = "zh"
	EN = "en"

	DEFALT_LANG = ZH // 默认中文

	DEFAULT_CODE           = 1000 // 参数错误
	SYSTEM_ERROR       int = 1001 // 系统错误
	PARAMS_BIND_ERROR      = 1002 // 参数绑定错误
	PARAMS_CHICK_ERROR     = 1003 // 参数校验错误
)

type Code struct {
	Zh string //中文
	En string //英文
}

type ErrCode = map[int]Code

var BASE_ERROR_CODE = ErrCode{
	DEFAULT_CODE:       {"参数错误", "params error"},
	SYSTEM_ERROR:       {"系统错误", "system error"},
	PARAMS_BIND_ERROR:  {"参数绑定错误", "Parameter binding error"},
	PARAMS_CHICK_ERROR: {"参数校验错误", "Parameter check error"},
}

func getMsg4Lang(lang string, codeObject Code) string {
	if lang == EN {
		return codeObject.En
	}
	return codeObject.Zh
}

func getLang(c Context) string {
	lang, _ := c.Get(TRANSLATE_LANGUAGE_HEADER_KEY).(string)
	//if !ok || (lang != ZH && lang != EN) {
	//	lang = DEFALT_LANG //没有获取到语种就默认中文
	//}
	if lang == "" {
		lang = DEFALT_LANG
	}
	return lang
}

// errAppendMsg err信息附加msg
func errAppendMsg(c Context, err *Err) {
	lang := getLang(c) //获取语言

	//先看有没有在框架中存在的错误码,再去项目中查找
	// 框架中找错误码
	codeObject, ok := BASE_ERROR_CODE[err.Code]
	if ok {
		err.Msg = getMsg4Lang(lang, codeObject)
	}
	if err.Msg == "" {
		// 项目中找错误码
		if len(c.GetGlobalErrorCodes()) > 0 {
			codeObject, ok := c.GetGlobalErrorCodes()[err.Code]
			if ok {
				err.Msg = getMsg4Lang(lang, codeObject)
			}
		}
	}
	// 始终没找到
	if err.Msg == "" {
		err.Code = DEFAULT_CODE
		err.Msg = getMsg4Lang(lang, BASE_ERROR_CODE[DEFAULT_CODE])
	}
}
