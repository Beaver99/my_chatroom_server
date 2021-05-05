package logic

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

// port 6378
func initGroupInfoRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "",
		DB:       1,
	})

	var ctx = context.Background()

	n, e := client.Keys(ctx, "*").Result()
	log.Println("------- KEYS: ", n, e)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println("fail to Ping redis-server:", err)
		return nil
	}
	return client
}

func GetGroupInfoDB() *redis.Client {
	return groupInfoDB
}

func GetMembers(ctx context.Context,groupID UUID)([]string, bool){
	members, err:=groupInfoDB.SMembers(ctx,string(groupID)).Result()
	log.Println(members)
	if err == redis.Nil{
		log.Println("this group does not exsists")
		return nil, false
	}else if err != nil{
		log.Println(err)
		return nil, false
	}
	return members, true
}

