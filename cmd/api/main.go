package main

import (
    "fmt"
    "net/http"
    "os"
    "strconv"

    "github.com/redis/go-redis/v9"
    "github.com/go-chi/chi/v5"
    "github.com/dash-xd/logma/internal/app"
)

func main() {
    fmt.Printf("Starting server...\n")
    redisClient := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URI"),
        Password: os.Getenv("REDISCLI_AUTH"),
        DB:       0,
    })
    defer redisClient.Close()
    entityRegistrar := entity.NewRedisEntityRegistrar(redisClient)
    router := app.NewRouter(redisClient, entityRegistrar)
    port := getPortFromArgs()
    addr := fmt.Sprintf(":%d", port)
    fmt.Printf("Go Server is listening on http://localhost%s\n", addr)
    if err := http.ListenAndServe(addr, router); err != nil {
        panic(err)
    }
}

func getPortFromArgs() int {
    defaultPort := 8080
    if len(os.Args) > 1 {
        port, err := strconv.Atoi(os.Args[1])
        if err != nil {
            fmt.Println("Invalid port provided. Using default port:", defaultPort)
            return defaultPort
        }
        return port
    }
    return defaultPort
}
