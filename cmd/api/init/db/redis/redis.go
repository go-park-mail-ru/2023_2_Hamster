package redis

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	addr string
	port string
}

func initRedisConfigFromEnv() (RedisConfig, error) {
	var cfg = RedisConfig{}
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}

	host, existHost := os.LookupEnv("REDIS_HOST")
	port, existPort := os.LookupEnv("REDIS_PORT")

	if !existHost || !existPort {
		return cfg, errors.New("existHost or existPort is Empty")
	}

	cfg = RedisConfig{
		addr: host,
		port: port,
	}
	return cfg, nil
}

func InitRedisCli(ctx context.Context) (*redis.Client, error) {
	cfg, err := initRedisConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return redis.NewClient(&redis.Options{
		Addr:     cfg.addr + cfg.port,
		Password: "",
		DB:       0,
	}), nil
}
