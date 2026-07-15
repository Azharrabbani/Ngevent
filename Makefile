-include .env
export

# DB URL construction for migrations
DB_URL := postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable

# Build the application
all: build build-worker test

build:
	@echo "Building app..."
	@go build -o main.exe cmd/app/main.go

build-worker:
	@echo "Building worker..."
	@go build -o worker.exe cmd/worker/main.go

# Run the application
run:
	@echo "Running app..."
	@go run cmd/app/main.go

run-worker:
	@echo "Running worker..."
	@go run cmd/worker/main.go

run-frontend:
	@echo "Running frontend..."
	@npm install --prefer-offline --no-fund --prefix ./frontend
	@npm run dev --prefix ./frontend

# Docker Compose shortcuts
docker-run:
	@docker compose up -d

docker-down:
	@docker compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Database Migrations
migrate-up:
	@echo "Running migrations up..."
	@migrate -path internal/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	@echo "Running migrations down..."
	@migrate -path internal/migrations -database "$(DB_URL)" -verbose down

migrate-create:
	@echo "Creating new migration (Usage: make migrate-create name=init)..."
	@migrate create -ext sql -dir internal/migrations -seq "$(name)"

# Code quality
fmt:
	@echo "Formatting..."
	@go fmt ./...

vet:
	@echo "Vetting..."
	@go vet ./...

lint: fmt vet

# Dependency management
tidy:
	@echo "Tidying go modules..."
	@go mod tidy

deps:
	@echo "Downloading dependencies..."
	@go mod download

# Database / Redis shell access
psql:
	@docker exec -it postgres-ngevent psql -U $(DB_USERNAME) -d $(DB_DATABASE)

redis-cli:
	@docker exec -it redis-ngevent redis-cli

# Asynq queue inspection (requires asynq CLI installed)
asynq-scheduled:
	@asynq task ls --queue=default --state=scheduled

asynq-archived:
	@asynq task ls --queue=default --state=archived

asynq-pending:
	@asynq task ls --queue=default --state=pending

# Clean the binaries
clean:
	@echo "Cleaning..."
	@rm -f main.exe worker.exe

# Live Reload (Air)
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

.PHONY: all build build-worker run run-worker run-frontend docker-run docker-down test itest \
	migrate-up migrate-down migrate-create fmt vet lint tidy deps psql redis-cli \
	asynq-scheduled asynq-archived asynq-pending clean watch