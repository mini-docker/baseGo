package validate

import (
	"baseGo/src/fecho/valid"
	"baseGo/src/model/code"
)

var (
	codes = make(ErrCode)
	ZH    = "ZH"
	EN    = "EN"
)

func init() {
	all := code.ErrCodes()
	for k, v := range all {
		codes[k] = v
	}
}

type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e Err) Error() string {
	return e.Msg
}

type ErrCode map[int]code.Code

type FrontValidate struct{}

func (fv *FrontValidate) Validate(i interface{}) error {
	v := &valid.Validation{}
	ok, c, err := v.Valid(i)
	if err != nil {
		return err
	}
	if !ok {
		return &Err{Code: int(c)}
	}
	return nil
}

func Find(code int, lang string) string {
	if val, ok := codes[code]; ok {
		if lang == ZH {
			return val.Zh
		}
		return val.En
	}
	return ""
}
