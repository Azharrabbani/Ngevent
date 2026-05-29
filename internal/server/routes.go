package server

import (
	"ngevent/internal/handler"
	"ngevent/internal/model"
	"ngevent/internal/utils/middleware"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"

	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var v1 fiber.Router

func (s *FiberServer) RegisterFiberRoutes() {
	// Panic recovery
	s.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Unique ID per request for tracing logs
	s.App.Use(requestid.New())

	// Logger
	s.App.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} ${method} ${path} | ${ip} | ${latency} | reqid=${locals:requestid}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// HTTP headers
	s.App.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "SAMEORIGIN",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
	}))

	// Override CORP for static media
	s.App.Use(func(c *fiber.Ctx) error {
		staticPaths := []string{
			"/api/v1/organizer/photo",
			"/api/v1/organizer/npwp",
			"/api/v1/organizer/nib",
			"/api/v1/attendee/photo",
			"/api/v1/event/banner",
			"/api/v1/updated-event/banner",
			"/api/v1/staging-organizer/npwp",
			"/api/v1/staging-organizer/nib",
		}

		for _, path := range staticPaths {
			if strings.HasPrefix(c.Path(), path) {
				c.Set("Cross-Origin-Resource-Policy", "cross-origin")
				return c.Next()
			}
		}
		return c.Next()
	})

	// CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://127.0.0.1:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type,Origin,X-Requested-With",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// grzip response compression
	s.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Global rate limiter
	s.App.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-Forwaded-For", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    429,
				"status":  "error",
				"message": "too many request",
			})
		},
	}))

	// Idempotency
	s.App.Use(idempotency.New())

	v1 = s.App.Group("api/v1")

}

func (s *FiberServer) RegisterAuthRoutes(h *handler.AuthHandler) {
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-Forwaded-For", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    429,
				"status":  "error",
				"message": "too many request",
			})
		},
	})

	v1.Put("/verify-email", authLimiter, h.VerififyEmail)

	v1.Post("/login", authLimiter, h.Login)

	v1.Post("/refresh", authLimiter, h.Refresh)

	v1.Post("/forgot-password", authLimiter, h.ForgotPassword)

	v1.Put("/reset-password/:id", authLimiter, h.ResetPassword)

	v1.Get("/phone-codes", authLimiter, h.ListPhoneCodes)

	logout := v1.Group("logout")
	logout.Use(middleware.AuthMiddleware())
	{
		logout.Post("/", h.Logout)
	}
}

func (s *FiberServer) RegisterUserRoutes(h *handler.UserHandler) {
	user := v1.Group("/user")

	user.Use(middleware.OptionalAuthMiddleware())

	user.Post("/register", h.Register)

	user.Use(middleware.AuthMiddleware())
	{
		user.Get("/me", h.FindCurrentUser)
		user.Get("/", middleware.AuthorizeRoles(string(model.Admin)), h.ListUsers)
		user.Put("/role", h.SelectRole)
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
		profile.Post("/", middleware.AuthorizeRoles(string(model.Attendee)), h.CreateProfile)
		profile.Get(("/"), middleware.AuthorizeRoles(string(model.Admin)), h.GetAllProfiles)
		profile.Get("/:id", h.GetProfileByID)
		profile.Get("/", h.GetProfileByUserID)
		profile.Get("/check-profile", h.HasProfile)
		profile.Put("/photo", middleware.AuthorizeRoles(string(model.Attendee)), h.UpdateProfilePhoto)
		profile.Put("/", middleware.AuthorizeRoles(string(model.Attendee)), h.UpdateProfile)
	}
}

func (s *FiberServer) RegisterOrganizerProfileRoutes(h *handler.OrganizerProfileHandler) {
	profile := v1.Group("/organizer")
	profile.Static("/photo", "./storage/profiles")
	profile.Static("/npwp", "./storage/npwp")
	profile.Static("/nib", "./storage/nib")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.Post("/", middleware.AuthorizeRoles(string(model.Organizer)), h.CreateProfile)
		profile.Get("/profiles", h.GetAllProfile)
		profile.Get("/", h.GetProfileByUserID)
		profile.Get("/:id", h.GetProfileByID)
		profile.Get("/filter", h.FilterProfile)
		profile.Put("/photo", middleware.AuthorizeRoles(string(model.Organizer), string(model.Admin)), h.UpdatePhotoProfile)
		profile.Put("/", middleware.AuthorizeRoles(string(model.Organizer), string(model.Admin)), h.UpdateProfile)
		profile.Put("/approve/:id", middleware.AuthorizeRoles(string(model.Admin)), h.ApprovedProfile)
		profile.Put("/reject/:id", middleware.AuthorizeRoles(string(model.Admin)), h.RejectProfile)
	}
}

func (s *FiberServer) RegisterOrganizerUpdateRoutes(h *handler.OrganizerUpdateHandler) {
	update := v1.Group("/staging-organizer")
	update.Static("/npwp", "./storage/npwp/stage")
	update.Static("/nib", "./storage/nib/stage")
	update.Use(middleware.AuthMiddleware(), middleware.AuthorizeRoles(string(model.Admin)))
	{
		update.Put("/:id", h.ValidateUpdate)
		update.Get("/update/:id", h.FindByProfileID)
		update.Get("/updates", h.FindUpdatesByProfileID)
		update.Get("/:id", h.FindUpdateByID)
	}
}

func (s *FiberServer) RegisterCategoriesRoutes(h *handler.CategoriesHandler) {
	category := v1.Group("/category")
	category.Use(middleware.AuthMiddleware())
	{
		category.Post("/", middleware.AuthorizeRoles(string(model.Admin)), h.CreateCategory)
		category.Get("/", h.ListCategories)
		category.Get("/list", middleware.AuthorizeRoles(string(model.Admin)), h.ListWithPagination)
		category.Get("/filter", h.ListByCatName)
		category.Put("/:id", middleware.AuthorizeRoles(string(model.Admin)), h.UpdateCategory)
		category.Delete("/:id", middleware.AuthorizeRoles(string(model.Admin)), h.DeleteCategory)
	}
}

func (s *FiberServer) RegisterEventRoutes(h *handler.EventHandler) {
	routeLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-Forwaded-For", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    429,
				"status":  "error",
				"message": "route computation is rate limited, please wait a moment",
			})
		},
	})

	event := v1.Group("/event")
	event.Static("/banner", "./storage/event/banner")
	event.Use(middleware.AuthMiddleware())
	{
		event.Post("/", middleware.AuthorizeRoles(string(model.Organizer)), h.CreateEvent)
		event.Get("/", middleware.AuthorizeRoles(string(model.Attendee), string(model.Admin)), h.GetEvents)
		event.Get("/nearest", middleware.AuthorizeRoles(string(model.Attendee)), h.FindNearestEvents)
		event.Get("/organizer-events", middleware.AuthorizeRoles(string(model.Organizer), string(model.Admin)), h.GetEventsByProfileID)
		event.Get("/route/:id", routeLimiter, h.GetEventRoute)
		event.Put("/review/:id", middleware.AuthorizeRoles(string(model.Admin)), h.ReviewEvent)
		event.Put("/cancel/:id", middleware.AuthorizeRoles(string(model.Organizer)), h.CancelEvent)
		event.Get("/:id", h.GetEventByID)
		event.Put("/:id", middleware.AuthorizeRoles(string(model.Organizer)), h.UpdateEvent)
		event.Delete("/:id", middleware.AuthorizeRoles(string(model.Organizer)), h.DeleteEvent)
	}
}

func (s *FiberServer) RegisterUpdatedEventRoutes(h *handler.UpdatedEventHandler) {
	updateEvent := v1.Group("/updated-event")
	updateEvent.Static("/banner", "./storage/updated/banner")
	updateEvent.Use(middleware.AuthMiddleware())
	{
		updateEvent.Get("/", middleware.AuthorizeRoles(string(model.Admin)), h.ListAllUpdated)
		updateEvent.Put("/review/:id", middleware.AuthorizeRoles(string(model.Admin)), h.ReviewUpdate)
		updateEvent.Get("/update-list/:event_id", middleware.AuthorizeRoles(string(model.Admin)), h.ListAllUpdatedByEventID)
		updateEvent.Get("/:event_id", middleware.AuthorizeRoles(string(model.Admin)), h.GetUpdatedByEventID)
		updateEvent.Put("/:id", middleware.AuthorizeRoles(string(model.Organizer)), h.CancelUpdate)
	}

}

func (s *FiberServer) RegisterLocationRoutes(h *handler.LocationHandler) {
	location := v1.Group("/location")
	location.Get("/", h.SearchLocation)
}
