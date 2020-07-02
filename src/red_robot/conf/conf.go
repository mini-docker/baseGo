package conf

import (
	"baseGo/src/fecho/utility"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	Version  string
	mConfig  *Config
	ConfigRW sync.RWMutex //配置因为会有变动并且会有外部访问,所以加锁

	siteIdsRWMutex  sync.RWMutex
	siteIds         = make([]string, 0)
	DwMemberCanPlay = []string{"ll", "pk", "eg", "egtc", "egcs", "ig", "pkspt"}

	IndexPageCacheMapBy4G sync.Map // 4g首页缓存 [4gIndexPageUrl] -> string
	Conf                  *Config
)

// 配置
type Config struct {
	App       AppConfig         `yaml:"app"`
	Log       LogConfig         `yaml:"log"`
	Mysql     MysqlConfig       `yaml:"mysql"`
	Redis     RedisConfig       `yaml:"redis"`
	Session   SessionConfig     `yaml:"session"`
	Registry  RegistryConfig    `yaml:"registry"`
	Storage   StorageServConfig `yaml:"storage"`
	Listening ListeningConfig   `yaml:"listening"`
	CDN       CDNConfig         `yaml:"cdn"`
}

type ListeningConfig struct {
	Add     string `yaml:"add"`
	Md5key  string `yaml:"md5key"`
	Deskey  string `yaml:"deskey"`
	SendAdd string `yaml:"sendAdd"`
}

// session信息配置参数
type SessionConfig struct {
	ExpiredTime  int `yaml:"sessionExpiredTime" `  //token过期时间,外面配置文件过期时间以分钟为单位
	GoroutineNum int `yaml:"sessionGoroutineNum" ` //可以开的协程数量
	Runtime      int `yaml:"sessionRuntime" `      //定时任务多长时间执行一次，以秒为单位
}

// 新文件服务配置
type StorageServConfig struct {
	Host      string `yaml:"host"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
}

//Mysql配置
type MysqlConfig struct {
	Host        string `yaml:"mysqlHost"`
	Username    string `yaml:"mysqlUsername" `
	Password    string `yaml:"mysqlPassword,omitempty" `
	DbName      string `yaml:"mysqlDbName"`
	Timeout     string `yaml:"mysqlTimeout"`
	TablePrefix string `yaml:"tablePrefix"`
	ShowSql     int    `yaml:"showSql,omitempty"`
}

// redis配置
type RedisConfig struct {
	Mode       int    `yaml:"mode,omitempty"`       // redis模式：0-直连 1-哨兵
	MasterName string `yaml:"masterName,omitempty"` // 哨兵模式的节点名称，直连模式则忽略该项
	Addrs      string `yaml:"addrs,omitempty"`      // 直连或者哨兵地址（哨兵地址使用逗号分割多地址）
	Password   string `yaml:"password,omitempty"`   // 密码
	DB         int    `yaml:"db,omitempty"`         // 数据库
}

//log配置
type LogConfig struct {
	Level      string `yaml:"level" `
	MaxLogSize int    `yaml:"maxLogSize" `
	Path       string `yaml:"path" `
	RecordPath string `yaml:"recordPath"` // 新加--统计日志文件路径
}

//本应用配置
type AppConfig struct {
	Addr           string `yaml:"addr"`
	ApiPort        int    `yaml:"apiPort"`
	AppEnvironment string `yaml:"appEnvironment"` // 系统运行环境
	ZkUrl          string `yaml:"zkUrl"`          // zookeeper 地址
	UniqueId       int    `yaml:"uniqueId"`       //唯一id 生成uuid？
	BcacheExpire   int    `yaml:"bcacheExpire"`   // 缓存过期时间. 2的倍数.
}

type RegistryConfig struct {
	Addr string `yaml:"addr"` // etcd 地址. 通过 "," 分割.
	TTL  int    `yaml:"ttl"`  // ttl 节点过期时间.
}

// cdn配置
type CDNConfig struct {
	Host string `yaml:"host"` // 域名
}

func GetConfig() *Config {
	return mConfig
}

// GetConfig 保障线程安全,加锁获取配置信息
func GetAppConfig() AppConfig {
	ConfigRW.RLock()
	defer ConfigRW.RUnlock()
	return mConfig.App
}

func GetRegistryConfig() RegistryConfig {
	ConfigRW.RLock()
	defer ConfigRW.RUnlock()
	return mConfig.Registry
}

func GetStorageServConfig() StorageServConfig {
	ConfigRW.RLock()
	defer ConfigRW.RUnlock()
	return mConfig.Storage
}

func GetCDNConfig() CDNConfig {
	ConfigRW.RLock()
	defer ConfigRW.RUnlock()
	return mConfig.CDN
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
