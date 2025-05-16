package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedisClient(host, port, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})
}

func SetBlacklistToken(ctx context.Context, rdb *redis.Client, token string, duration time.Duration) error {
	return rdb.Set(ctx, token, "blacklisted", duration).Err()
}

func IsTokenBlacklisted(ctx context.Context, rdb *redis.Client, token string) (bool, error) {
	res, err := rdb.Get(ctx, token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return res == "blacklisted", nil
}
