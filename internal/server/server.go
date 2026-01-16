package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ngevent/internal/config"
)

type FiberServer struct {
	*fiber.App

	DB              *gorm.DB
	ClientWoker     *asynq.Client
	InspectorWorker *asynq.Inspector
	RDB             *redis.Client
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ngevent",
			AppName:      "ngevent",
		}),

		DB:              config.ConnectDB(),
		ClientWoker:     config.GetAsynqClient(),
		InspectorWorker: config.GetAsynqInspector(),
		RDB:             config.GetRedisClient(),
	}

	return server
}
