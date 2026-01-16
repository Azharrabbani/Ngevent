package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	redisOnce sync.Once
	rdb       *redis.Client
)

func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}

		// Load redis config
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		if host == "" || port == "" {
			log.Fatal("❌ Missing Redis configuration: REDIS_HOST or REDIS_PORT not set")
		}

		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: "",
			DB:       0,
		})

		ctx := context.Background()

		// Test connection
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatal("Redis connection failed:", err)
		}
	})

	log.Println("✅ Redis connected")

	return rdb
}
