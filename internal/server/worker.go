package server

import (
	"ngevent/internal/config"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type Worker struct {
	Srv *asynq.Server
	DB  *gorm.DB
	Mux *asynq.ServeMux
}

func NewWorker() *Worker {
	return &Worker{
		Srv: config.GetAsynqServer(),
		DB:  config.ConnectDB(),
		Mux: config.GetAsynqMux(),
	}
}
