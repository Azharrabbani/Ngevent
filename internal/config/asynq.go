package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

var (
	asynqClientOnce    sync.Once
	asynqInspectorOnce sync.Once
	asynqServerOnce    sync.Once
	asynqMuxOnce       sync.Once

	client    *asynq.Client
	inspector *asynq.Inspector
	server    *asynq.Server
	mux       *asynq.ServeMux
)

func opt() *asynq.RedisClientOpt {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load redis config
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	if host == "" || port == "" {
		log.Fatal("❌ Missing Redis configuration: REDIS_HOST  or REDIS_PORT")
	}

	return &asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", host, port),
	}
}

func GetAsynqClient() *asynq.Client {
	asynqClientOnce.Do(func() {
		client = asynq.NewClient(opt())

		// Test connection
		if err := client.Ping(); err != nil {
			log.Fatal("❌ Asynq Client connection failed:", err)
		}
	})

	return client
}

func GetAsynqInspector() *asynq.Inspector {
	asynqInspectorOnce.Do(func() {
		inspector = asynq.NewInspector(opt())
	})

	return inspector
}

func GetAsynqServer() *asynq.Server {
	cfg := asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Printf("task %s failed: %v", task.Type(), err)
		}),
		ShutdownTimeout: 30 * time.Second,
	}

	asynqServerOnce.Do(func() {
		server = asynq.NewServer(opt(), cfg)

		// Test connection
		if err := server.Ping(); err != nil {
			log.Fatal("❌ Asynq Server connection failed:", err)
		}
	})

	return server
}

func GetAsynqMux() *asynq.ServeMux {
	asynqMuxOnce.Do(func() {
		mux = asynq.NewServeMux()
	})

	return mux
}
