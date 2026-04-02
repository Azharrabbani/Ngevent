package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InvalidateCache(client *redis.Client, patterns []string) {
	ctx := context.Background()

	for _, pattern := range patterns {
		iter := client.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			client.Del(ctx, iter.Val())
		}
	}

	// Use SCAN for pattern keys to avoid blocking
	log.Println("[CACHE] cache invalidated")
}
