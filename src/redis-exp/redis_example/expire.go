package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// ps -ef | grep redis
// redis-cli -h 127.0.0.1 -p 6379
// -a "mypass"
func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}
	defer c.Close()
	_, err = c.Do("expire", "abc", 10)
	if err != nil {
		fmt.Println(err)
		return
	}
}
