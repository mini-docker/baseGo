package conf

import (
	"strings"
	"sync"

	"github.com/go-redis/redis"
)

var (
	// 私有redis连接
	redisClient   *redis.Client
	redisClientRW sync.RWMutex
)

// 初始化私有redis(存放登录session以外的数据)
func InitRedis(cfg *RedisConfig) error {
	// 哨兵模式
	if cfg.Mode == 1 {
		cliNew, err := NewRedisBySentinel(cfg.MasterName, cfg.Addrs, cfg.Password, cfg.DB)
		if err != nil {
			return err
		}
		redisClientRW.Lock()
		redisClient = cliNew
		redisClientRW.Unlock()
		return nil
	}

	// 直连模式
	cliNew, err := NewRedis(cfg.Addrs, cfg.Password, cfg.DB)
	if err != nil {
		return err
	}
	redisClientRW.Lock()
	redisClient = cliNew
	redisClientRW.Unlock()
	return nil
}

// 获取私有redis客户端连接
func GetRedis() *redis.Client {
	redisClientRW.RLock()
	defer redisClientRW.RUnlock()
	return redisClient
}

// 获取新的redis客户端连接（直连）
func NewRedis(addr, password string, dbNum int) (*redis.Client, error) {
	cliNew := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNum,
		// Addr:     "127.0.0.1:6379",
		// Password: "",
		// DB:       0,
	})

	_, err := cliNew.Ping().Result()
	if err != nil {
		return nil, err
	}
	return cliNew, nil
}

// 获取新的redis客户端连接（哨兵）
func NewRedisBySentinel(masterName, addr, password string, dbNum int) (*redis.Client, error) {
	options := &redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: strings.Split(addr, ","),
		Password:      password,
		DB:            dbNum,
	}
	cliNewBySentinel := redis.NewFailoverClient(options)
	_, err := cliNewBySentinel.Ping().Result()
	if err != nil {
		return nil, err
	}
	return cliNewBySentinel, nil
}
