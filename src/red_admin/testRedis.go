package main

import (
	"github.com/go-redis/redis"
)

var redisdb *redis.Client

// 初始化连接
func initClient() (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// defer redisdb.Close()
	_, err = redisdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	initClient()
}
