package app

import (
    "context"
    "os"
  
    "github.com/redis/go-redis/v9"
)

type Callback func(*PublishRequest) error

type Listener struct {
    ChannelName string
    RedisClient *redis.Client
    Callbacks   []Callback
}

func NewListener(redisClient *redis.Client, channelName string) *Listener {
    return &Listener{
        ChannelName: channelName,
        RedisClient: redisClient,
        Callbacks:   make([]Callback, 0),
    }
}

func (l *Listener) AddCallback(cb Callback) {
    l.Callbacks = append(l.Callbacks, cb)
}

func (l *Listener) Start(ctx context.Context) {
    pubsub := l.RedisClient.Subscribe(ctx, l.ChannelName)
    defer pubsub.Close()

    pubsubChannel := pubsub.Channel()
    for msg := range pubsubChannel {
        var message PublishRequest
        for _, cb := range l.Callbacks {
            if err := cb(&message); err != nil {
            }
        }
    }
}
