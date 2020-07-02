package conf

import (
	xtime "baseGo/src/imserver/pkg/time"
	"fecho/golog"
	"fecho/utility"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

var (
	confPath string
	Conf     *Config
)
var LocalTest bool

// Init init config.
func Init(cfg *Config) error {
	Conf = cfg
	//日志初始化
	err := InitLog(Conf.Log)
	if err != nil {
		golog.Error("initd", "InitLog", "config log error %v", err, "log", Conf.Log)
		return err
	}
	cfg.Kafka.Brokers = strings.Split(cfg.Kafka.BrokersBySplit, ",")
	cfg.Env.Host, _ = os.Hostname()

	if false {
		_, err = toml.DecodeFile(confPath, &Conf)
	}

	return nil
}

// Config is job config.
type Config struct {
	Env      *Env            `yaml:"env"`
	Kafka    *Kafka          `yaml:"kafka"`
	Registry *RegistryConfig `yaml:"registry"`
	Comet    *Comet          `yaml:"comet"`
	Room     *Room           `yaml:"room"`
	Log      LogConfig       `yaml:"log"`
}

//log配置
type LogConfig struct {
	Level      string `yaml:"level" `
	MaxLogSize int    `yaml:"maxLogSize" `
	Path       string `yaml:"path" `
	RecordPath string `yaml:"recordPath"` // 新加--统计日志文件路径
}

// Room is room config.
type Room struct {
	Batch  int            `yaml:"batch"`
	Signal xtime.Duration `yaml:"signal"`
	Idle   xtime.Duration `yaml:"idle"`
}

// Comet is comet config.
type Comet struct {
	RoutineChan int `yaml:"routineChan"`
	RoutineSize int `yaml:"routineSize"`
}

// Kafka is kafka config.
type Kafka struct {
	Topic          string `yaml:"topic"`
	Group          string `yaml:"group"`
	Brokers        []string
	BrokersBySplit string `yaml:"brokersBySplit"`
}

// Env is env config.
type Env struct {
	Host string `yaml:"host"`
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
