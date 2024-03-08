package app

import (
    "context"
    "strings"
  
    "github.com/redis/go-redis/v9"
)

type EntityRegistrar interface {
    Register(args ...string) (string, error)
}

type RedisEntityRegistrar struct {
    Client *redis.Client
}

func NewRedisEntityRegistrar(client *redis.Client) *RedisEntityRegistrar {
    return &RedisEntityRegistrar{Client: client}
}

func (r *RedisEntityRegistrar) Register(args ...string) (string, error) {
    result, err := r.Client.Do(context.Background(), append("FCALL", args...)).Result()
    if err != nil {
        return "", err
    }
    
    key := strings.TrimSpace(result.(string))
    return key, nil
}
