package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"com.lc.go.codepush/server/config"
	"github.com/redis/go-redis/v9"
)

var redisDB *redis.Client
var ctx context.Context = context.Background()

func GetRedis() (redisD *redis.Client, err error) {
	if redisDB != nil {
		redisD = redisDB
		return
	}
	redisConfig := config.GetConfig().Redis
	redisDB = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + strconv.Itoa(int(redisConfig.Port)),
		Username: redisConfig.UserName,
		Password: redisConfig.Password,
		DB:       int(redisConfig.DBIndex),
	})
	redisD = redisDB
	return
}

func SetRedisObj(key string, obj any, duration time.Duration) {
	redis, _ := GetRedis()
	jData, _ := json.Marshal(obj)
	status := redis.Set(ctx, key, string(jData), duration)
	if err := status.Err(); err != nil {
		fmt.Println(err.Error())
	}
}

func DelRedisObj(key string) {
	redis, _ := GetRedis()

	iter := redis.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		err := redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic("Redis error:" + err.Error())
		}
	}
	if err := iter.Err(); err != nil {
		panic("Redis error:" + err.Error())
	}
}

func GetRedisObj[T any](key string) *T {

	redis, _ := GetRedis()
	status := redis.Get(ctx, key)
	if err := status.Err(); err != nil {
		log.Println(err.Error())
		return nil
	}
	var obj T
	err := json.Unmarshal([]byte(status.Val()), &obj)
	if err != nil {
		panic("Redis error:" + err.Error())
	}
	return &obj
}
