# Backend
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc g++ make libwebp-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download go mod download
RUN go mod download

# Install asynq CLI
RUN go install github.com/hibiken/asynq/tools/asynq@latest

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o app ./cmd/app
RUN CGO_ENABLED=1 GOOS=linux go build -o worker ./cmd/worker


# Final stage backend
FROM alpine:latest AS prod

RUN apk add --no-cache libwebp tzdata ghostscript

# Workdir App
WORKDIR /app
COPY --from=builder /app/app . 
COPY --from=builder /app/worker .
COPY --from=builder /go/bin/asynq /usr/local/bin/asynq
COPY .env .env

RUN mkdir -p /app/storage/profiles /app/storage/npwp /app/storage/nib /app/storage/npwp/stage /app/storage/nib/stage /app/storage/event/banner

EXPOSE 8080
CMD ["./app"]

# Frontend
FROM node:25-alpine AS frontend_builder
WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm install
COPY frontend/. .
RUN npm run build

# Final stage frontend
FROM node:25-slim AS frontend
RUN npm install -g serve
COPY --from=frontend_builder /frontend/dist /app/dist
EXPOSE 5173
CMD ["serve", "-s", "/app/dist", "-l", "5173"]
