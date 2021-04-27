package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var userInfoDB = initUserInfoRedis()

func initUserInfoRedis() *redis.Client {
	//log.Println("Redis already connected")
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	var ctx = context.Background()

	n, e := client.Keys(ctx, "*").Result()
	log.Println("------- KEYS: ",n, e)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println("fail to Ping redis-server:", err)
		return nil
	}
	return client
}
