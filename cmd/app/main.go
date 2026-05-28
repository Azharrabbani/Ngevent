package main

import (
	"context"
	"fmt"
	"log"
	"ngevent/internal/handler"
	"ngevent/internal/repository"
	"ngevent/internal/server"
	"ngevent/internal/service"
	"ngevent/internal/tasks"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	fiberServer.Close()

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	server := server.New()

	// ensure to close the worker inspector & client connection
	defer server.InspectorWorker.Close()
	defer server.ClientWoker.Close()

	// Create a new validator instance
	validate := validator.New()

	// init Task Publisher Worker
	unverifiedUserTaskPublisher := tasks.NewUserTaskPublisher(server.ClientWoker, server.InspectorWorker)
	unusedOTPTaskPublisher := tasks.NewOTPTaskPublisher(server.ClientWoker, server.InspectorWorker)
	emailTaskPublisher := tasks.NewEmailTaskPublisher(server.ClientWoker, server.InspectorWorker)
	eventExpiryPublisher := tasks.NewEventExpiryPublisher(server.ClientWoker, server.InspectorWorker)

	// Init repository
	userRepo := repository.NewUsersRepository(server.DB)
	sessionRepo := repository.NewSessionRepository(server.DB)
	otpRepo := repository.NewOtpRepository(server.DB)
	attendeeProfileRepo := repository.NewAttendeeProfileRepository(server.DB)
	organizerProfileRepo := repository.NewOrganizerRepository(server.DB)
	organizerUpdateRepo := repository.NewOrganizerProfileUpdateRepository(server.DB)
	categoryRepo := repository.NewCategoriesRepository(server.DB)
	eventRepo := repository.NewEventsRepository(server.DB)
	updateEventRepo := repository.NewUpdatedEventsRepository(server.DB)

	// Init service
	authService := service.NewAuthService(userRepo, sessionRepo, otpRepo, unverifiedUserTaskPublisher, unusedOTPTaskPublisher, emailTaskPublisher)
	userService := service.NewUserService(userRepo, attendeeProfileRepo, organizerProfileRepo, otpRepo, unverifiedUserTaskPublisher, unusedOTPTaskPublisher, emailTaskPublisher, server.RDB)
	otpService := service.NewOTPService(userRepo, otpRepo, unusedOTPTaskPublisher, unverifiedUserTaskPublisher, emailTaskPublisher)
	attendeeProfileService := service.NewAttendeeProfileService(attendeeProfileRepo, server.RDB)
	organizerProfileService := service.NewOrganizerProfileService(organizerProfileRepo, userRepo, organizerUpdateRepo, emailTaskPublisher, server.RDB)
	organizerUpdateService := service.NewOrganizerUpdateService(userRepo, organizerProfileRepo, organizerUpdateRepo, emailTaskPublisher, server.RDB)
	categoryService := service.NewCategoryService(categoryRepo, server.RDB)
	eventService := service.NewEventService(eventRepo, updateEventRepo, userRepo, organizerProfileRepo, categoryRepo, eventExpiryPublisher, emailTaskPublisher, server.RDB)
	updateEventService := service.NewUpdatedEventService(updateEventRepo, eventRepo, organizerProfileRepo, eventExpiryPublisher, server.RDB, emailTaskPublisher)
	locationService := service.NewLocationService(server.RDB)

	// Init handler
	authHandler := handler.NewAuthHandler(authService, validate)
	userHandler := handler.NewUserHandler(userService, validate)
	otpHandler := handler.NewOTPHandler(otpService, validate)
	attendeeProfileHandler := handler.NewAttendeeProfileService(attendeeProfileService, validate)
	organizerProfileHandler := handler.NewOrganizerProfileHandler(organizerProfileService, validate)
	organizerUpdateHandler := handler.NewOrganizerUpdateHandler(organizerUpdateService, validate)
	categoryHandler := handler.NewCategoriesHandler(categoryService, validate)
	eventHandler := handler.NewEventHandler(eventService, validate)
	updateEventHandler := handler.NewUpdatedEventHandler(updateEventService, validate)
	locationHandler := handler.NewLocationHandler(locationService)

	// Register routes
	server.RegisterFiberRoutes()
	server.RegisterAuthRoutes(authHandler)
	server.RegisterUserRoutes(userHandler)
	server.RegisterOTPRoutes(otpHandler)
	server.RegisterAttendeeProfileRoutes(attendeeProfileHandler)
	server.RegisterOrganizerProfileRoutes(organizerProfileHandler)
	server.RegisterOrganizerUpdateRoutes(organizerUpdateHandler)
	server.RegisterCategoriesRoutes(categoryHandler)
	server.RegisterEventRoutes(eventHandler)
	server.RegisterUpdatedEventRoutes(updateEventHandler)
	server.RegisterLocationRoutes(locationHandler)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
