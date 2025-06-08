package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rds *redis.Client
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

func MustNewRedis(conf Config) *Redis {
	rds := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Ping测试连接是否成功
	if err := rds.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	result := Redis{
		rds: rds,
	}

	return &result
}
