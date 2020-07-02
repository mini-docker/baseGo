package conf

import (
	"fecho/golog"
	"fecho/utility"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"

	//"github.com/bilibili/discovery/naming"
	xtime "baseGo/src/imserver/pkg/time"
)

var (
	confPath string
	Conf     *Config
)
var Localtest bool

// Init init config.
func Init(cfg1 *Config) error {
	Conf = cfg1
	//日志初始化
	err := InitLog(Conf.Log)
	if err != nil {
		golog.Error("initd", "InitLog", "config log error %v", err, "log", Conf.Log)
		return err
	}
	//处理云配置不支持数组和bool类型的情况
	cfg1.Env.Addrs = strings.Split(cfg1.Env.AddrsBySplit, ",")
	if cfg1.Env.OfflineByBool == 0 {
		cfg1.Env.Offline = false
	} else {
		cfg1.Env.Offline = true
	}
	cfg1.TCP.Bind = strings.Split(cfg1.TCP.BindBySplit, ",")
	if cfg1.TCP.KeepAliveByBool == 0 {
		cfg1.TCP.KeepAlive = false
	} else {
		cfg1.TCP.KeepAlive = true
	}

	cfg1.Websocket.Bind = strings.Split(cfg1.Websocket.BindBySplit, ",")
	cfg1.Websocket.TLSBind = strings.Split(cfg1.Websocket.TLSBindBySplit, ",")
	if cfg1.Websocket.TLSOpenByBool == 0 {
		cfg1.Websocket.TLSOpen = false
	} else {
		cfg1.Websocket.TLSOpen = true
	}

	whiteListTmp := make([]int64, 0)
	for _, v := range strings.Split(cfg1.Whitelist.WhitelistBySplit, ",") {
		ival, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		whiteListTmp = append(whiteListTmp, int64(ival))
	}
	cfg1.Whitelist.Whitelist = whiteListTmp

	cfg1.Env.Host, _ = os.Hostname()
	if false {
		_, err = toml.DecodeFile(confPath, &Conf)
	}

	return nil
}

// Config is comet config.
type Config struct {
	Debug     bool            `yaml:"debug"`
	Env       *Env            `yaml:"env"`
	Registry  *RegistryConfig `yaml:"registry"`
	TCP       *TCP            `yaml:"tcp"`
	Websocket *Websocket      `yaml:"websocket"`
	Protocol  *Protocol       `yaml:"protocol"`
	Bucket    *Bucket         `yaml:"bucket"`
	RPCClient *RPCClient      `yaml:"rpcClient"`
	RPCServer *RPCServer      `yaml:"rpcServer"`
	Whitelist *Whitelist      `yaml:"whitelist"`
	Log       LogConfig       `yaml:"log"`
}

//log配置
type LogConfig struct {
	Level      string `yaml:"level" `
	MaxLogSize int    `yaml:"maxLogSize" `
	Path       string `yaml:"path" `
	RecordPath string `yaml:"recordPath"` // 新加--统计日志文件路径
}

// Env is env config.
type Env struct {
	DeployEnv     string `yaml:"deployEnv"`
	Host          string `yaml:""`
	Weight        int64  `yaml:"weight"`
	Offline       bool
	OfflineByBool int `yaml:"offlineByBool"`
	Addrs         []string
	AddrsBySplit  string `yaml:"addrsBySplit"`
}

// RPCClient is RPC client config.
type RPCClient struct {
	Dial    xtime.Duration `yaml:"dial"`
	Timeout xtime.Duration `yaml:"timeout"`
}

// RPCServer is RPC server config.
type RPCServer struct {
	Network           string         `yaml:"network"`
	Addr              string         `yaml:"addr"`
	Timeout           xtime.Duration `yaml:"timeout"`
	IdleTimeout       xtime.Duration
	MaxLifeTime       xtime.Duration
	ForceCloseWait    xtime.Duration
	KeepAliveInterval xtime.Duration
	KeepAliveTimeout  xtime.Duration
}

// TCP is tcp config.
type TCP struct {
	Bind            []string
	BindBySplit     string `yaml:"bindBySplit"`
	Sndbuf          int    `yaml:"sndbuf"`
	Rcvbuf          int    `yaml:"rcvbuf"`
	KeepAlive       bool
	KeepAliveByBool int    `yaml:"keepAliveByBool"`
	Reader          int    `yaml:"reader"`
	ReadBuf         int    `yaml:"readBuf"`
	ReadBufSize     int    `yaml:"readBufSize"`
	Writer          int    `yaml:"writer"`
	WriteBuf        int    `yaml:"writeBuf"`
	WriteBufSize    int    `yaml:"writeBufSize"`
	PubUrl          string `yaml:"pubUrl"`
}

// Websocket is websocket config.
type Websocket struct {
	Bind           []string
	BindBySplit    string `yaml:"bindBySplit"`
	TLSOpen        bool
	TLSOpenByBool  int `yaml:"TLSOpenByBool"`
	TLSBind        []string
	TLSBindBySplit string `yaml:"tlsBindBySplit"`
	CertFile       string `yaml:"certFile"`
	PrivateFile    string `yaml:"privateFile"`
	PubUrl         string `yaml:"pubUrl"`
}

// Protocol is protocol config.
type Protocol struct {
	Timer            int            `yaml:"timer"`
	TimerSize        int            `yaml:"timerSize"`
	SvrProto         int            `yaml:"svrProto"`
	CliProto         int            `yaml:"cliProto"`
	HandshakeTimeout xtime.Duration `yaml:"handshakeTimeout"`
}

// Bucket is bucket config.
type Bucket struct {
	Size          int    `yaml:"size"`
	Channel       int    `yaml:"channel"`
	Room          int    `yaml:"room"`
	RoutineAmount uint64 `yaml:"routineAmount"`
	RoutineSize   int    `yaml:"routineSize"`
}

// Whitelist is white list config.
type Whitelist struct {
	Whitelist        []int64
	WhitelistBySplit string `yaml:"whitelistBySplit"`
	WhiteLog         string `yaml:"whiteLog"`
}

type RegistryConfig struct {
	Addr     string `yaml:"addr"`     // etcd 地址. 通过 "," 分割.
	TTL      int    `yaml:"ttl"`      // ttl 节点过期时间.
	Interval int    `yaml:"interval"` // interval 节点间隔时间.
}

//从文件解析配置
func ParseConfigFile(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ParseConfigData(data)
}

//从数据解析配置
func ParseConfigData(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	path, _ := utility.ExecPath()

	if !strings.HasPrefix(cfg.Log.Path, "/") && os.Getenv("LOG_OUTPUT") != "" {
		cfg.Log.Path = filepath.Join(path, cfg.Log.Path)
	}
	return &cfg, nil
}
