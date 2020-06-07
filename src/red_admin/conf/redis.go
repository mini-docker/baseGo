package conf

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/FZambia/sentinel"
	"github.com/garyburd/redigo/redis"
)

var (
	// 私有redis连接
	redisClient   *redis.Pool
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
func GetRedis() *redis.Pool {
	redisClientRW.RLock()
	defer redisClientRW.RUnlock()
	return redisClient
}

// 获取新的redis客户端连接（直连）
func NewRedis(addr, password string, dbNum int) (*redis.Pool, error) {
	var errs error
	pool := &redis.Pool{
		MaxIdle:     dbNum, // 最大等待连接中的数量,设 0 为没有限制
		MaxActive:   1000,  // 最大连接数据库连接数,设 0 为没有限制
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				errs = err
				return nil, err
			}
			fmt.Println(password, "password")
			// if _, err := c.Do("AUTH", password); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			errs = err
			return err
		},
	}
	return pool, errs
}

// 获取新的redis客户端连接（哨兵）
func NewRedisBySentinel(masterName, addr, password string, dbNum int) (*redis.Pool, error) {
	// options := &redis.FailoverOptions{
	// 	MasterName:    masterName,
	// 	SentinelAddrs: strings.Split(addr, ","),
	// 	Password:      password,
	// 	DB:            dbNum,
	// }
	// cliNewBySentinel := redis.NewFailoverClient(options)
	// _, err := cliNewBySentinel.Ping().Result()
	// if err != nil {
	// 	return nil, err
	// }
	// return cliNewBySentinel, nil

	// redisAddr := "192.168.1.11:26378,192.168.1.22:26378"
	// redisAddrs := strings.Split(redisAddr, ",")
	// masterName := "master1" // 根据redis集群具体配置设置
	fmt.Println(dbNum, "dbNum")
	var errs error
	sntnl := &sentinel.Sentinel{
		Addrs:      strings.Split(addr, ","),
		MasterName: masterName,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	cliNewBySentinel := &redis.Pool{
		MaxIdle:     dbNum,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				errs = err
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				errs = err
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: CheckRedisRole,
	}

	return cliNewBySentinel, errs

}
func CheckRedisRole(c redis.Conn, t time.Time) error {
	if !sentinel.TestRole(c, "master") {
		return fmt.Errorf("Role check failed")
	} else {
		return nil
	}
}
