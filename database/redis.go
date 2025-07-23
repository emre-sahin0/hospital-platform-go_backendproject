package database

import (
	"context"
	"fmt"
	"hospital-platform/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	// Environment variables'dan Redis bağlantı bilgilerini al
	redisHost := config.GetEnv("REDIS_HOST", "localhost")
	redisPort := config.GetEnv("REDIS_PORT", "6379")
	redisPassword := config.GetEnv("REDIS_PASSWORD", "")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis bağlantısı başarısız: %v", err))
	}

	fmt.Printf("✅ Redis bağlantısı başarılı: %s\n", addr)
}

func SetResetCode(phone string, code string) error {
	err := RedisClient.Set(Ctx, phone, code, 5*time.Minute).Err() // kod 5 dakika geçerli
	return err
}

func GetResetCode(phone string) (string, error) {
	code, err := RedisClient.Get(Ctx, phone).Result()
	return code, err
}

func DeleteResetCode(phone string) error {
	err := RedisClient.Del(Ctx, phone).Err()
	return err
}
