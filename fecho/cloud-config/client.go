package cloud_config

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

//type ConfData interface {
//	ParseConfig(dd *Data, conf *interface{}) error
//}

type ConfigForm struct {
	EnvName string `json:"e"`  // 开发环境
	AppId   string `json:"a"`  // 应用id
	Version int    `json:"v" ` // 版本号
	Data    string `json:"d"`
}

//""code":70000,"msg":"获取配置失败""
type Data struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Data    string `json:"d,omitempty"`
	Version int64  `json:"v,omitempty"`
}

//type CloudConfig struct {
//	ServerUrl string //url
//	AppID     string
//	Env       string
//	Other     []string
//}

// configRemote, env, name string, other ...string
// 远程获取配置
//func (cc *CloudConfig) BootConfig(conf interface{}) error {
//	data := ConfigForm{}
//	data.EnvName = cc.Env
//	data.AppId = cc.AppID
//	data.Version = 0
//	if len(cc.Other) > 0 {
//		data.Data = strings.Join(cc.Other, "@@")
//	} else {
//		data.Data = ""
//	}
//
//	buf, err := httpPostForm(cc.ServerUrl, data)
//	if err != nil {
//		return err
//	}
//	//解密
//	dd := new(Data)
//	err = json.Unmarshal(buf, dd)
//	if err != nil {
//		return err
//	}
//	if dd.Code != 0 {
//		return fmt.Errorf("%s error %d %s", cc.AppID, dd.Code, dd.Msg)
//	}
//	return parseConfig(dd, conf)
//}

//func (cc *CloudConfig) UpdateConfig(ver int64, conf interface{}) (bool, error) {
//	data := ConfigForm{}
//	data.EnvName = cc.Env
//	data.AppId = cc.AppID
//	data.Version = 0
//	if len(cc.Other) > 0 {
//		data.Data = strings.Join(cc.Other, "@@")
//	} else {
//		data.Data = ""
//	}
//
//	buf, err := httpPostForm(cc.ServerUrl, data)
//	if err != nil {
//		return false, err
//	}
//	//解密
//	dd := new(Data)
//	err = json.Unmarshal(buf, dd)
//	if err != nil {
//		return false, err
//	}
//
//	if ver == dd.Version {
//		return false, nil
//	} else if ver < dd.Version {
//		return true, parseConfig(dd, conf)
//	}
//	return false, nil
//}

func httpPostForm(postUrl string, data ConfigForm) ([]byte, error) {
	client := &http.Client{}
	formData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", postUrl, bytes.NewReader(formData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return body, nil
}

//func parseConfig(dd *Data, conf interface{}) error {
//	b, err := base64.StdEncoding.DecodeString(dd.Data)
//	if err != nil {
//		return err
//	}
//	confData, err := DoZlibUnCompress(b)
//	if err != nil {
//		return err
//	}
//
//	if len(confData) > 0 {
//		err = json.Unmarshal(confData, conf)
//		if err != nil {
//			return err
//		}
//	} else {
//		return errors.New("DoZlibUnCompress config error")
//	}
//	return nil
//}

const DICT = `{"log":{"log":{"log":{"level":"debug","maxLogSize":"1000000000","path":"/tmp/log.log"},"mogon":{"host":"127.0.0.1","password":"name112"}},"redis":""}}`

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevelDict(&in, 9, []byte(DICT)) //Dict() //NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReaderDict(b, []byte(DICT)) //.NewReader(b)
	_, err := io.Copy(&out, r)
	return out.Bytes(), err
}

//找到对应的时间戳
const SEC = 6 * 60

func GetRealTimeLine(tl int64) int64 {
	var allDate int64 = 24 * 60 * 60
	startTimeDayStartTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" 00:00:00")
	startTimeDayStartTimeLine := startTimeDayStartTime.Unix()
	if tl == 0 {
		tl = time.Now().Unix()
	}
	timeNowLine := tl - startTimeDayStartTimeLine
	k := allDate / SEC
	var i int64 = 0
	//切割时间
	for i = 0; i < k; i++ {
		min := i * (allDate / k)
		max := (i + 1) * (allDate / k) //最大时间
		if timeNowLine >= min && timeNowLine <= max {
			return startTimeDayStartTimeLine + min
		}
	}
	return 0
}

//func (cc *CloudConfig) BootConfigCustomer(parseFunc func([]byte) (interface{}, error)) (interface{}, error) {
//	data := ConfigForm{}
//	data.EnvName = cc.Env
//	data.AppId = cc.AppID
//	data.Version = 0
//	if len(cc.Other) > 0 {
//		data.Data = strings.Join(cc.Other, "@@")
//	} else {
//		data.Data = ""
//	}
//	buf, err := httpPostForm(cc.ServerUrl, data)
//	if err != nil {
//		return nil, err
//	}
//	//解密
//	dd := new(Data)
//	err = json.Unmarshal(buf, dd)
//	if err != nil {
//		return nil, err
//	}
//	if dd.Code != 0 {
//		return nil, fmt.Errorf("%s error %d %s", cc.AppID, dd.Code, dd.Msg)
//	}
//
//	b, err := base64.StdEncoding.DecodeString(dd.Data)
//	if err != nil {
//		return nil, err
//	}
//	confData, err := DoZlibUnCompress(b)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(confData) > 0 {
//		return parseFunc(confData)
//	} else {
//		return nil, errors.New("DoZlibUnCompress config error")
//	}
//}
//
//func (cc *CloudConfig) UpdateConfigCustomer(ver int64, parseFunc func([]byte) (interface{}, error)) (interface{}, bool, error) {
//	data := ConfigForm{}
//	data.EnvName = cc.Env
//	data.AppId = cc.AppID
//	data.Version = 0
//	if len(cc.Other) > 0 {
//		data.Data = strings.Join(cc.Other, "@@")
//	} else {
//		data.Data = ""
//	}
//
//	buf, err := httpPostForm(cc.ServerUrl, data)
//	if err != nil {
//		return nil, false, err
//	}
//	//解密
//	dd := new(Data)
//	err = json.Unmarshal(buf, dd)
//	if err != nil {
//		return nil, false, err
//	}
//
//	if ver == dd.Version {
//		return nil, false, nil
//	} else if ver < dd.Version {
//		b, err := base64.StdEncoding.DecodeString(dd.Data)
//		if err != nil {
//			return nil, false, err
//		}
//		confData, err := DoZlibUnCompress(b)
//		if err != nil {
//			return nil, false, err
//		}
//
//		if len(confData) > 0 {
//			res, err := parseFunc(confData)
//			if err != nil {
//				return nil, false, err
//			}
//			return res, true, nil
//		} else {
//			return nil, false, errors.New("DoZlibUnCompress config error")
//		}
//	}
//
//	return nil, false, nil
//}
