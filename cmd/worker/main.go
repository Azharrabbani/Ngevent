package main

import (
	"log"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/server"
	"ngevent/internal/service"
	"ngevent/internal/tasks"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize worker
	worker := server.NewWorker()

	// Init repository
	userRepo := repository.NewUsersRepository(worker.DB)
	sessionRepo := repository.NewSessionRepository(worker.DB)
	attendeeProfileRepo := repository.NewAttendeeProfileRepository(worker.DB)
	organizerProfileRepo := repository.NewOrganizerRepository(worker.DB)
	otpRepo := repository.NewOtpRepository(worker.DB)
	eventRepo := repository.NewEventsRepository(worker.DB)
	updatedEventRepo := repository.NewUpdatedEventsRepository(worker.DB)

	// Init service
	authService := service.NewAuthService(userRepo, sessionRepo, otpRepo, nil, nil, nil)
	userService := service.NewUserService(userRepo, attendeeProfileRepo, organizerProfileRepo, otpRepo, nil, nil, nil, nil)
	eventExpiryService := service.NewEventExpiryService(eventRepo, updatedEventRepo)

	// Init tasks handler
	userTaskHandler := tasks.NewUserTaskHandler(userService)
	otpTaskHandler := tasks.NewOTPTaskHandler(authService)
	eventExpiryHandler := tasks.NewEventExpiryHandler(eventExpiryService)

	// Setup handlers
	worker.Mux.HandleFunc(model.TypeVerifiedUser, userTaskHandler.HandlerUnverifiedTask)
	worker.Mux.HandleFunc(model.TypeVerifiedOTP, otpTaskHandler.HandlerUnusedOTP)
	worker.Mux.HandleFunc(model.TypeEventExpired, eventExpiryHandler.HandleEventExpired)
	worker.Mux.HandleFunc(model.TypeUpdatedEventExpired, eventExpiryHandler.HandleUpdatedEventExpired)

	// Email handlers
	worker.Mux.HandleFunc(model.TypeEMailVerify, tasks.NewEmailTaskHandler().HandlerUserVerification)
	worker.Mux.HandleFunc(model.TypeEmailForgetPassword, tasks.NewEmailTaskHandler().HandlerUserForgetPassword)
	worker.Mux.HandleFunc(model.TypeEmailAdminVerification, tasks.NewEmailTaskHandler().HandlerAdminVerifyProfile)
	worker.Mux.HandleFunc(model.TypeEmailOrganizerProfile, tasks.NewEmailTaskHandler().HandlerOrganizerProfileNotif)
	worker.Mux.HandleFunc(model.TypeEmailOrganizerProfileVerified, tasks.NewEmailTaskHandler().HandlerOrganizerProfileVerified)
	worker.Mux.HandleFunc(model.TypeEmailOrganizerProfileRejected, tasks.NewEmailTaskHandler().HandlerOrganizerProfileRejected)

	// Event email handler
	worker.Mux.HandleFunc(model.TypeEventAdminNotification, tasks.NewEmailTaskHandler().HandlerEventAdminNotification)
	worker.Mux.HandleFunc(model.TypeEventEONotification, tasks.NewEmailTaskHandler().HandlerEventOrganizerNotification)
	worker.Mux.HandleFunc(model.TypeEventEOVerification, tasks.NewEmailTaskHandler().HandlerEventOrganizerVerification)
	worker.Mux.HandleFunc(model.TypeEventUpdateNotification, tasks.NewEmailTaskHandler().HandlerUpdateEventOrganizerNotif)

	// Start Worker Server
	if err := worker.Srv.Run(worker.Mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down worker server...")

	defer worker.Srv.Shutdown()

	log.Println("worker server stopped")

}
