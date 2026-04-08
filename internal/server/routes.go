package server

import (
	"ngevent/internal/handler"
	"ngevent/internal/model"
	"ngevent/internal/utils/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var v1 fiber.Router

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://127.0.0.1:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,Origin,X-Requested-With",
		AllowCredentials: true, // credentials require explicit origins
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
	user := v1.Group("/user")
	user.Post("/register", h.Register)

	user.Use(middleware.AuthMiddleware())
	{
		user.Get("/", middleware.AuthorizeRoles("admin"), h.ListUsers)
		user.Get("/id", h.FindUserByID)
	}
}

func (s *FiberServer) RegisterOTPRoutes(h *handler.OTPHandler) {
	v1.Post("/resend-otp", h.ResendOTP)
}

func (s *FiberServer) RegisterAttendeeProfileRoutes(h *handler.AttendeeProfileHandler) {
	profile := v1.Group("/attendee")
	profile.Static("/photo", "./storage/profiles")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.Post("/", middleware.AuthorizeRoles("user"), h.CreateProfile)
		profile.Get("/:id", h.GetProfileByID)
		profile.Get("/", h.GetProfileByUserID)
		profile.Put("/photo", middleware.AuthorizeRoles("user"), h.UpdateProfilePhoto)
		profile.Put("/", middleware.AuthorizeRoles("user"), h.UpdateProfile)
	}
}

func (s *FiberServer) RegisterOrganizerProfileRoutes(h *handler.OrganizerProfileHandler) {
	profile := v1.Group("/organizer")
	profile.Static("/photo", "./storage/profiles")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.Post("/", middleware.AuthorizeRoles("event organizer"), h.CreateProfile)
		profile.Get("/:id", h.GetProfileByID)
		profile.Get("/profiles", h.GetAllProfile)
		profile.Get("/profile", h.GetProfileByUserID)
		profile.Get("/filter", h.FilterProfile)
		profile.Put("/photo", middleware.AuthorizeRoles("event organizer", "admin"), h.UpdatePhotoProfile)
		profile.Put("/", middleware.AuthorizeRoles("event organizer", "admin"), h.UpdateProfile)
		profile.Put("/approve/:id", middleware.AuthorizeRoles("admin"), h.ApprovedProfile)
		profile.Put("/reject/:id", middleware.AuthorizeRoles("admin"), h.RejectProfile)
	}
}

func (s *FiberServer) RegisterOrganizerUpdateRoutes(h *handler.OrganizerUpdateHandler) {
	update := v1.Group("/staging-organizer")
	update.Use(middleware.AuthMiddleware())
	{
		update.Put("/:id", middleware.AuthorizeRoles("admin"), h.ValidateUpdate)
		update.Get("/:id", h.FindUpdateByID)
		update.Get("/", h.FindUpdateByProfileID)
	}
}

func (s *FiberServer) RegisterCategoriesRoutes(h *handler.CategoriesHandler) {
	category := v1.Group("/category")
	category.Use(middleware.AuthMiddleware())
	{
		category.Post("/", middleware.AuthorizeRoles("admin"), h.CreateCategory)
		category.Get("/", h.ListCategories)
		category.Get("/filter", h.ListByCatName)
		category.Put("/:id", middleware.AuthorizeRoles("admin"), h.UpdateCategory)
		category.Delete("/:id", middleware.AuthorizeRoles("admin"), h.DeleteCategory)
	}
}

func (s *FiberServer) RegisterEventRoutes(h *handler.EventHandler) {
	event := v1.Group("/event")
	event.Static("/banner", "./storage/event/banner")
	event.Use(middleware.AuthMiddleware())
	{
		event.Post("/", middleware.AuthorizeRoles(string(model.Organizer)), h.CreateEvent)
		event.Get("/", middleware.AuthorizeRoles(string(model.Attendee), string(model.Admin)), h.GetEvents)
		event.Get("/nearest", middleware.AuthorizeRoles(string(model.Attendee)), h.FindNearestEvents)
		event.Get("/organizer-events", middleware.AuthorizeRoles(string(model.Organizer), string(model.Admin)), h.GetEventsByProfileID)
		event.Put("/review", middleware.AuthorizeRoles(string(model.Admin)), h.ReviewEvent)
		event.Put("/cancel", middleware.AuthorizeRoles(string(model.Organizer)), h.CancelEvent)
		event.Get("/:id", h.GetEventByID)
		event.Put("/:id", middleware.AuthorizeRoles(string(model.Organizer)), h.UpdateEvent)
	}
}

func (s *FiberServer) RegisterUpdatedEventRoutes(h *handler.UpdatedEventHandler) {
	updateEvent := v1.Group("/updated-event")
	updateEvent.Static("/updated-banner", "./storage/updated_event/banner")
	updateEvent.Use(middleware.AuthMiddleware())
	{
		updateEvent.Get("/", middleware.AuthorizeRoles(string(model.Admin)), h.ListAllUpdated)
		updateEvent.Get("/:id", middleware.AuthorizeRoles(string(model.Admin)), h.GetUpdatedByID)
		updateEvent.Get("/update-list/:event_id", middleware.AuthorizeRoles(string(model.Admin)), h.ListAllUpdatedByEventID)
		updateEvent.Put("/", middleware.AuthorizeRoles(string(model.Admin)), h.ReviewUpdate)
		updateEvent.Put("/:id", middleware.AuthorizeRoles(string(model.Organizer)), h.CancelUpdate)
	}

}
