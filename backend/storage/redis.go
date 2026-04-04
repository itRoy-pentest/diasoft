package storage

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ConnectRedis() *redis.Client {
	addr := os.Getenv("REDIS_HOST")
	if addr == "" {
		addr = "localhost:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // Если ставил пароль на VPS, впиши сюда
		DB:       0,
	})
}

func CreatePublicToken(rdb *redis.Client, token string, idHash string, ttl time.Duration) error {
	return rdb.Set(ctx, "token:"+token, idHash, ttl).Err()
}

func GetHashByToken(rdb *redis.Client, token string) (string, error) {
	return rdb.Get(ctx, "token:"+token).Result()
}
