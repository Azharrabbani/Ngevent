package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"ngevent/internal/database"
)

type FiberServer struct {
	*fiber.App

	DB *gorm.DB
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ngevent",
			AppName:      "ngevent",
		}),

		DB: database.ConnectDB(),
	}

	return server
}
