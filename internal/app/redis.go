package app

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(
		&redis.Options{
			Addr:     os.Getenv("REDIS_URI"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		},
	)
}
