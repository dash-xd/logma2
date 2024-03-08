package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/go-chi/chi/v5"
	"github.com/dash-xd/logma/internal/entity"
	"github.com/dash-xd/logma/internal/listener"
)

var (
	RedisClient     *redis.Client
	EntityRegistrar entity.EntityRegistrar
)

func init() {
	// Initialize the Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDISCLI_AUTH"),
		DB:       0,
	})

	// Initialize the entity registrar with the Redis client
	EntityRegistrar = entity.NewRedisEntityRegistrar(RedisClient)
}

type Subscriber struct {
	RedisClient     *redis.Client
	EntityRegistrar entity.EntityRegistrar
}

func NewSubscriber(redisClient *redis.Client, entityRegistrar entity.EntityRegistrar) *Subscriber {
	return &Subscriber{
		RedisClient:     redisClient,
		EntityRegistrar: entityRegistrar,
	}
}

func (s *Subscriber) Subscribe(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		ChannelName string `json:"channelName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	listener := listener.NewListener(requestBody.ChannelName, s.RedisClient)

	// Add callback
	listener.AddCallback(func(message *listener.PublishRequest) error {
		// Implement callback logic -- todo
	})

	// Start listener
	go listener.Start(r.Context())

	// Register entity
	entityKey, err := s.EntityRegistrar.Register(r.Context(), requestBody.ChannelName)
	if err != nil {
		// Handle error
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"entityKey": entityKey})
}

func NewRouter() http.Handler {
	s := NewSubscriber(RedisClient, EntityRegistrar)
	r := chi.NewRouter()
	r.Post("/subscribe", s.Subscribe)
	return r
}
