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
	otpRepo := repository.NewOtpRepository(worker.DB)

	// Init service
	authService := service.NewAuthService(userRepo, sessionRepo, otpRepo, nil, nil)
	userService := service.NewUserService(userRepo, otpRepo, nil, nil)

	// Init tasks handler
	userTaskHandler := tasks.NewUserTaskHandler(userService)
	otpTaskHandler := tasks.NewOTPTaskHandler(authService)

	// Setup handlers
	worker.Mux.HandleFunc(model.TypeVerifiedUser, userTaskHandler.HandlerUnverifiedTask)
	worker.Mux.HandleFunc(model.TypeVerifiedOTP, otpTaskHandler.HandlerUnusedOTP)

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
