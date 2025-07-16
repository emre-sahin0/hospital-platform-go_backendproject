package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-11740.c240.us-east-1-3.ec2.redns.redis-cloud.com:11740",
		Password: "TdpPFktE2kFeUFlLXR4lhLsSzAUflKJm",
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis bağlantısı başarısız: %v", err))
	}

	fmt.Println("✅ Redis bağlantısı başarılı")
}
func SetResetCode(phone string, code string) error {
	err := RedisClient.Set(Ctx, phone, code, 5*time.Minute).Err() // kod 5 dakika geçerli
	return err
}
