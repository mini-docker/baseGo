package main

import "github.com/go-redis/redis"

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
