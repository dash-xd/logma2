package app

import (
    "context"
    "strings"
  
    "github.com/redis/go-redis/v9"
)

type EntityRegistrar interface {
    Register(ctx context.Context, args ...string) (string, error)
}

type RedisEntityRegistrar struct {
    Client *redis.Client
}

func NewRedisEntityRegistrar(client *redis.Client) *RedisEntityRegistrar {
    return &RedisEntityRegistrar{Client: client}
}

func (r *RedisEntityRegistrar) Register(ctx context.Context, args ...string) (string, error) {
    redisCmd := []string{"FCALL", "entity_registration_function"} // Adjust with your actual function name
    redisCmd = append(redisCmd, args...)
    result, err := r.Client.Do(ctx, redisCmd...).Result()
    if err != nil {
        return "", err
    }
    key := strings.TrimSpace(result.(string))
    return key, nil
}
