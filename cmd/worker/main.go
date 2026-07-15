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

	defer worker.InspectorWorker.Close()
	defer worker.ClientWoker.Close()

	emailTaskPublisher := tasks.NewEmailTaskPublisher(worker.ClientWoker, worker.InspectorWorker)

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
	eventExpiryService := service.NewEventExpiryService(eventRepo, updatedEventRepo, emailTaskPublisher, worker.RDB)
	eventService := service.NewEventService(eventRepo, userRepo, nil, nil, nil, nil, worker.RDB)

	// Init tasks handler
	userTaskHandler := tasks.NewUserTaskHandler(userService)
	otpTaskHandler := tasks.NewOTPTaskHandler(authService)
	eventExpiryHandler := tasks.NewEventExpiryHandler(eventExpiryService, eventService)
	emailTaskHandler := tasks.NewEmailTaskHandler()

	// Setup handlers — task eksekusi (bukan email)
	worker.Mux.HandleFunc(model.TypeVerifiedUser, userTaskHandler.HandlerUnverifiedTask)
	worker.Mux.HandleFunc(model.TypeVerifiedOTP, otpTaskHandler.HandlerUnusedOTP)
	worker.Mux.HandleFunc(model.TypeEventExpired, eventExpiryHandler.HandleEventExpired)
	worker.Mux.HandleFunc(model.TypeUpdatedEventExpired, eventExpiryHandler.HandleUpdatedEventExpired)
	worker.Mux.HandleFunc(model.TypeDraftRevert, eventExpiryHandler.HandleDraftRevert)
	worker.Mux.HandleFunc(model.TypeUpdatedDraftRevert, eventExpiryHandler.HandleUpdatedDraftRevert)

	// Email handlers — auth/profile
	worker.Mux.HandleFunc(model.TypeEMailVerify, emailTaskHandler.HandlerUserVerification)
	worker.Mux.HandleFunc(model.TypeEmailForgetPassword, emailTaskHandler.HandlerUserForgetPassword)
	worker.Mux.HandleFunc(model.TypeEmailAdminVerification, emailTaskHandler.HandlerAdminVerifyProfile)
	worker.Mux.HandleFunc(model.TypeEmailOrganizerProfile, emailTaskHandler.HandlerOrganizerProfileNotif)
	worker.Mux.HandleFunc(model.TypeEmailOrganizerProfileVerified, emailTaskHandler.HandlerOrganizerProfileVerified)
	worker.Mux.HandleFunc(model.TypeEmailOrganizerProfileRejected, emailTaskHandler.HandlerOrganizerProfileRejected)

	// Email handlers — event
	worker.Mux.HandleFunc(model.TypeEventAdminNotification, emailTaskHandler.HandlerEventAdminNotification)
	worker.Mux.HandleFunc(model.TypeEventEONotification, emailTaskHandler.HandlerEventOrganizerNotification)
	worker.Mux.HandleFunc(model.TypeEventEOVerification, emailTaskHandler.HandlerEventOrganizerVerification)
	worker.Mux.HandleFunc(model.TypeEventUpdateNotification, emailTaskHandler.HandlerUpdateEventOrganizerNotif)
	worker.Mux.HandleFunc(model.TypeEventEORevertNotification, emailTaskHandler.HandlerEventOrganizerRevertNotif)
	worker.Mux.HandleFunc(model.TypeEventEOUpdateRevertNotification, emailTaskHandler.HandlerEventOrganizerUpdateRevertNotif)

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
