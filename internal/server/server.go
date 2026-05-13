package server

import (
	"time"

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
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
			BodyLimit:    10 * 1024 * 1024,
			ServerHeader: "ngevent",
			AppName:      "ngevent",

			ErrorHandler: func(c *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				msg := "internal server error"

				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
					msg = e.Message
				}

				return c.Status(code).JSON(fiber.Map{
					"code":    code,
					"status":  "error",
					"message": msg,
				})
			},
		}),

		DB:              config.ConnectDB(),
		ClientWoker:     config.GetAsynqClient(),
		InspectorWorker: config.GetAsynqInspector(),
		RDB:             config.GetRedisClient(),
	}

	return server
}
