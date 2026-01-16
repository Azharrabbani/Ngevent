package server

import (
	"ngevent/internal/handler"
	"ngevent/internal/utils/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var v1 fiber.Router

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,Origin,X-Requested-With",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	v1 = s.App.Group("api/v1")

}

func (s *FiberServer) RegisterAuthRoutes(h *handler.AuthHandler) {
	v1.Put("/verify-email/:id", h.VerififyEmail)

	v1.Post("/login", h.Login)

	v1.Post("/forgot-password", h.ForgotPassword)

	v1.Put("/reset-password/:id", h.ResetPassword)

	logout := v1.Group("logout")
	logout.Use(middleware.AuthMiddleware())
	{
		logout.Post("/", h.Logout)
	}
}

func (s *FiberServer) RegisterUserRoutes(h *handler.UserHandler) {
	v1.Post("/register", h.Register)
}
