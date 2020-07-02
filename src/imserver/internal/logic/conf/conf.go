package conf

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/fecho/xorm"
	"baseGo/src/fecho/xorm/core"
	"baseGo/src/imserver/internal/logic/model"
	xtime "baseGo/src/imserver/pkg/time"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

var (
	confPath string
	host     string
	weight   int64
	Conf     *Config
)

var Localtest bool

// Init init config.
func Init(cfg1 *Config) error {
	Conf = cfg1
	//日志初始化
	err := InitLog(Conf.Log)
	if err != nil {
		golog.Error("im-logic", "Init", "config log error %v", err, "log", Conf.Log)
		return err
	}
	cfg1.Kafka.Brokers = strings.Split(cfg1.Kafka.BrokersBySplit, ",")
	cfg1.Env.Host, _ = os.Hostname()

	if false {
		_, err = toml.DecodeFile(fmt.Sprintf("/Users/js129/go/src/chat-server/js108/imserver/cmd/logic/%v", confPath), &Conf)
	}
	// 初始化子库
	err = InitMysql(Conf.Mysql)
	if err != nil {
		golog.Error("im-logic", "Init", "config SiteMysql error:%v", err)
		return err
	}

	err = model.InitRedis(&Conf.Redis)
	if err != nil {
		golog.Error("im-logic", "Init", "config redis %v %v", err)
		return err
	}
	return nil
}

// cdn配置
type CDNConfig struct {
	Host string `yaml:"host"` // 域名
}

// Config config.
type Config struct {
	Env        *Env              `yaml:"env"`
	RPCClient  *RPCClient        `yaml:"rpcClient"`
	RPCServer  *RPCServer        `yaml:"rpcServer"`
	HTTPServer *HTTPServer       `yaml:"httpServer"`
	Kafka      *Kafka            `yaml:"kafka"`
	CDN        CDNConfig         `yaml:"cdn"`
	Node       *Node             `yaml:"node"`
	Backoff    *Backoff          `yaml:"backoff"`
	Mysql      MysqlConfig       `yaml:"mysql"`
	Redis      model.RedisConfig `yaml:"redis"`
	Log        LogConfig         `yaml:"log"`
	Registry   *RegistryConfig   `yaml:"registry"`
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
	Host   string `yaml:"host"`
	Weight int64  `yaml:"weight"`
}

// Node node config.
type Node struct {
	DefaultDomain string         `yaml:"defaultDomain"`
	HostDomain    string         `yaml:"hostDomain"`
	TCPPort       int            `yaml:"tcpPort"`
	WSPort        int            `yaml:"wsPort"`
	WSSPort       int            `yaml:"wssPort"`
	HeartbeatMax  int            `yaml:"heartbeatMax"`
	Heartbeat     xtime.Duration `yaml:"heartbeat"`
	RegionWeight  float64        `yaml:"regionWeight"`
}

// Backoff backoff.
type Backoff struct {
	MaxDelay  int32   `yaml:"maxDelay"`
	BaseDelay int32   `yaml:"baseDelay"`
	Factor    float32 `yaml:"factor"`
	Jitter    float32 `yaml:"jitter"`
}

// Kafka .
type Kafka struct {
	Topic          string `yaml:"topic"`
	Brokers        []string
	BrokersBySplit string `yaml:"brokersBySplit"`
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
	IdleTimeout       xtime.Duration `yaml:"idleTimeout"`
	MaxLifeTime       xtime.Duration `yaml:"maxLifeTime"`
	ForceCloseWait    xtime.Duration `yaml:"forceCloseWait"`
	KeepAliveInterval xtime.Duration `yaml:"keepAliveInterval"`
	KeepAliveTimeout  xtime.Duration `yaml:"keepAliveTimeout"`
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string         `yaml:"network"`
	Addr         string         `yaml:"addr"`
	ReadTimeout  xtime.Duration `yaml:"readTimeout"`
	WriteTimeout xtime.Duration `yaml:"writeTimeout"`
}
type ServsConfig struct {
	Tcp  string `yaml:"tcp"`
	Ws   string `yaml:"ws"`
	Grpc string `yaml:"grpc"`
}

//Mysql配置
type MysqlConfig struct {
	Host        string `yaml:"mysqlHost"`
	Username    string `yaml:"mysqlUsername" `
	Password    string `yaml:"mysqlPassword,omitempty" `
	Timeout     string `yaml:"mysqlTimeout"`
	TablePrefix string `yaml:"tablePrefix"`
	ShowSql     int    `yaml:"showSql,omitempty"`
	DbName      string `yaml:"mysqlDbName"`
}

type RegistryConfig struct {
	Addr string `yaml:"addr"` // etcd 地址. 通过 "," 分割.
	TTL  int    `yaml:"ttl"`  // ttl 节点过期时间.
}

// InitMysql初始化mysql
func InitMysql(c MysqlConfig) error {
	engine, err := initMysql(&c)
	if err != nil {
		return err
	}
	primaryEngine = engine
	//fmt.Println(engines)
	return nil
}

func initMysql(c *MysqlConfig) (*xorm.Engine, error) {
	var err error
	//data := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&timeout=%s&interpolateParams=true",
	data := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&timeout=%s&interpolateParams=true&allowNativePasswords=true",
		c.Username,
		c.Password,
		c.Host,
		c.DbName,
		c.Timeout,
	)
	var engineTemp *xorm.Engine
	engineTemp, err = xorm.NewEngine("mysql", data)
	if err != nil {
		golog.Error("im-logic", "initMysql", "err :", err)
		return nil, err
	}
	engineTemp.ShowSQL(c.ShowSql == 1)
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, c.TablePrefix)
	engineTemp.SetTableMapper(tbMapper)
	engineTemp.SetLogger(new(golog.LoggingXorm))
	err = engineTemp.Ping()
	if err != nil {
		golog.Error("im-logic", "initMysql", "err :", err)
		return nil, err
	}
	//golog.Info( "mysql", "InitMysql", "mysql connect success", "host", c.Host)
	return engineTemp, nil
}

var (
	primaryEngine   *xorm.Engine
	primaryEngineRW sync.RWMutex

	KEY = []byte("9a8v2d5o") // 加密秘钥
)

// GetXormSession 获取xorm session
func GetXormSession() *xorm.Session {
	primaryEngineRW.RLock()
	defer primaryEngineRW.RUnlock()
	return primaryEngine.NewSession()
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
