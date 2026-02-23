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

	v1.Post("/refresh", h.Refresh)

	v1.Post("/forgot-password", h.ForgotPassword)

	v1.Put("/reset-password/:id", h.ResetPassword)

	v1.Get("/phone-codes", h.ListPhoneCodes)

	logout := v1.Group("logout")
	logout.Use(middleware.AuthMiddleware())
	{
		logout.Post("/", h.Logout)
	}
}

func (s *FiberServer) RegisterUserRoutes(h *handler.UserHandler) {
	v1.Post("/register", h.Register)
}

func (s *FiberServer) RegisterOTPRoutes(h *handler.OTPHandler) {
	v1.Post("/resend-otp", h.ResendOTP)
}

func (s *FiberServer) RegisterAttendeeProfileRoutes(h *handler.AttendeeProfileHandler) {
	profile := v1.Group("/attendee")
	profile.Static("/photo", "./storage/profiles")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.Post("/", h.CreateProfile)
		profile.Get("/:id", h.GetProfileByID)
		profile.Get("/", h.GetProfileByUserID)
		profile.Put("/photo", h.UpdateProfilePhoto)
		profile.Put("/", h.UpdateProfile)
	}
}

func (s *FiberServer) RegisterOrganizerProfileRoutes(h *handler.OrganizerProfileHandler) {
	profile := v1.Group("/organizer")
	profile.Static("/photo", "./storage/profiles")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.Post("/", h.CreateProfile)
		profile.Get("/:id", h.GetProfileByID)
		profile.Get("/profiles", h.GetAllProfile)
		profile.Get("/profile", h.GetProfileByUserID)
		profile.Get("/filter", h.FilterProfile)
		profile.Put("/photo", h.UpdatePhotoProfile)
		profile.Put("/", h.UpdateProfile)
		profile.Put("/approve/:id", h.ApprovedProfile)
		profile.Put("/reject/:id", h.RejectProfile)
	}
}

func (s *FiberServer) RegisterOrganizerUpdateRoutes(h *handler.OrganizerUpdateHandler) {
	update := v1.Group("/staging-organizer")
	update.Use(middleware.AuthMiddleware())
	{
		update.Put("/:id", h.ValidateUpdate)
		update.Get("/:id", h.FindUpdateByID)
		update.Get("/", h.FindUpdateByProfileID)
	}
}
