package server

import (
	"ngevent/internal/config"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Worker struct {
	Srv             *asynq.Server
	DB              *gorm.DB
	Mux             *asynq.ServeMux
	ClientWoker     *asynq.Client
	InspectorWorker *asynq.Inspector
	RDB             *redis.Client
}

func NewWorker() *Worker {
	return &Worker{
		Srv:             config.GetAsynqServer(),
		DB:              config.ConnectDB(),
		Mux:             config.GetAsynqMux(),
		ClientWoker:     config.GetAsynqClient(),
		InspectorWorker: config.GetAsynqInspector(),
		RDB:             config.GetRedisClient(),
	}
}
