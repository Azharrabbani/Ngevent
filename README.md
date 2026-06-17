# 🎟️ Ngevent

![Go](https://img.shields.io/badge/Go-1.25-blue?style=flat&logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7?style=flat)
![GORM](https://img.shields.io/badge/GORM-v1.31-59666C?style=flat)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-blue?style=flat&logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-red?style=flat&logo=redis)

![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?style=flat&logo=typescript)
![Vite](https://img.shields.io/badge/Vite-7-646CFF?style=flat&logo=vite)
![TailwindCSS](https://img.shields.io/badge/TailwindCSS-4.2-06B6D4?style=flat&logo=tailwindcss)
![TanStack Query](https://img.shields.io/badge/TanStack%20Query-v5-FF4154?style=flat)
![Leaflet](https://img.shields.io/badge/Leaflet-Maps-199900?style=flat&logo=leaflet)

> ⚠️ **This project is still ongoing / under active development.**

A full-stack event management platform built with **Go (Fiber)** on the backend and **React (Vite + TypeScript)** on the frontend.

---

## 🧱 Tech Stack

| Layer      | Technology                                    |
|------------|-----------------------------------------------|
| Backend    | Go ≥ 1.21, Fiber v2, GORM, PostgreSQL         |
| Frontend   | React, Vite, TypeScript                       |
| Auth       | JWT (Access + Refresh Token), OTP             |
| Queue      | Redis + Asynq (background worker)             |
| Email      | SMTP via gomail                               |
| DevOps     | Docker, Docker Compose, Air (hot reload)      |

---

## 🚀 Installation & Setup Guide

### 🧩 Prerequisites

Make sure you have installed:

- **Go** ≥ 1.21
- **Node.js** + **npm**
- **Docker** & **Docker Compose**
- **Make**
- **Redis** (atau jalankan via Docker)

---

### ⚙️ 1. Clone Repository

```bash
git clone https://github.com/Azharrabbani/ngevent.git
cd ngevent
```

---

### 📄 2. Prepare the Configuration File

Create a `.env` file in the project root directory based on `.env.example`:

```env
PORT=8080
APP_ENV=local

DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=ngevent
DB_USERNAME=postgres
DB_PASSWORD=your_password

# JWT
SECRET_KEY=your_secret_key
REFRESH_SECRET_KEY=your_refresh_secret_key

# SMTP
SMTP_HOST=your_smtp_host
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

---

### 🐘 3. Run with Docker Compose

```bash
docker compose up --build -d
```

This will run:
- **PostgreSQL** on port `5432`
- **Redis** on port `6379`
- **Backend API** on port `8080`
- **Worker** (background job processor)
- **Frontend** on port `5173`

---

### 🧰 4. Run Locally (Without Docker)

#### Build binary

```bash
make build
```

#### Run the backend and frontend at the same time

```bash
make run
```

> This command runs the Go backend in the background and the frontend dev server (`npm run dev`).

---

### 🔄 5. Live Reload (Hot Reload)

```bash
make watch
```

> Using [`air`](https://github.com/air-verse/air). If it isn't already installed, this command will install it automatically.

---

### 🧪 6. Running the test

```bash
# Unit tests
make test

# Integration tests
make itest
```

---

### 🧼 7. Clean the Build

```bash
make clean
```

---

### 🐳 8. Manage Docker Container

```bash
# run the container
make docker-run

# stop the container
make docker-down
```

---

## 📚 API Documentation - Routing Map

Base URL: `http://localhost:8080/api/v1`

---

## 🔐 Authentication Routes

| Method | Endpoint                        | Description                     | Request Body                                                               |
|--------|---------------------------------|---------------------------------|----------------------------------------------------------------------------|
| POST   | `/api/v1/user/register`         | Register user (attendee/organizer) | `{"name": "string", "email": "string", "password": "string", "role": "attendee\|organizer"}` |
| GET    | `/api/v1/phone-codes`           | Fetch international phone codes | -                                                                          |
| POST   | `/api/v1/login`                 | Login user                      | `{"email": "string", "password": "string"}`                                |
| POST   | `/api/v1/refresh`               | Refresh access token            | -                                                                          |
| PUT    | `/api/v1/verify-email`          | Verify user's email             | -                                                                          |
| POST   | `/api/v1/forgot-password`       | Request password reset link     | `{"email": "string"}`                                                      |
| PUT    | `/api/v1/reset-password/:id`    | Reset password using token      | `{"new_password": "string", "confirm_password": "string"}`                 |
| POST   | `/api/v1/resend-otp`            | Resend OTP                      | -                                                                          |
| POST   | `/api/v1/logout`                | Logout (Protected by JWT)       | -                                                                          |

---

## 👤 User Management

*Route prefix: `/api/v1/user/*`*

| Method | Endpoint                          | Description                    | Middleware            | Request Body                                                |
|--------|-----------------------------------|--------------------------------|-----------------------|-------------------------------------------------------------|
| POST   | `/api/v1/user/register`           | Register new user              | -                     | `{"name": "string", "email": "string", "password": "string", "role": "attendee\|organizer"}` |
| GET    | `/api/v1/user/me`                 | Get current user               | JWT                   | -                                                           |
| POST   | `/api/v1/user/admin/register`     | Create admin user              | JWT + Admin only      | `{"name": "string", "email": "string", "password": "string"}` |
| GET    | `/api/v1/user/`                   | List all users                 | JWT + Admin only      | -                                                           |
| PUT    | `/api/v1/user/role`               | Select / switch user role      | JWT                   | `{"role": "string"}`                                        |
| GET    | `/api/v1/user/id`                 | Find user by ID                | JWT                   | -                                                           |

---

## 👤 Attendee Profile Management

*Route prefix: `/api/v1/attendee/*`*

| Method | Endpoint                          | Description                    | Middleware            | Request Body / Params                                                    |
|--------|-----------------------------------|--------------------------------|-----------------------|--------------------------------------------------------------------------|
| GET    | `/api/v1/attendee/photo/{filename}` | Get attendee profile photo   | -                     | -                                                                        |
| POST   | `/api/v1/attendee/`               | Create attendee profile        | JWT + Attendee role   | `{"username": "string", "bio": "string", ...}`                           |
| GET    | `/api/v1/attendee/`               | Get all attendee profiles      | JWT + Admin only      | -                                                                        |
| GET    | `/api/v1/attendee/:id`            | Get attendee profile by ID     | JWT                   | -                                                                        |
| GET    | `/api/v1/attendee/check-profile`  | Check if user has a profile    | JWT                   | -                                                                        |
| PUT    | `/api/v1/attendee/photo`          | Update attendee profile photo  | JWT + Attendee role   | FormData: `photo(file)`                                                  |
| PUT    | `/api/v1/attendee/`               | Update attendee profile        | JWT + Attendee role   | `{"username": "string", "bio": "string", ...}`                           |

---

## 🏢 Organizer Profile Management

*Route prefix: `/api/v1/organizer/*`*

| Method | Endpoint                              | Description                           | Middleware                    | Request Body / Notes                                    |
|--------|---------------------------------------|---------------------------------------|-------------------------------|---------------------------------------------------------|
| GET    | `/api/v1/organizer/photo/{filename}`  | Get organizer profile photo           | -                             | -                                                       |
| GET    | `/api/v1/organizer/npwp/{filename}`   | Get organizer NPWP document           | -                             | -                                                       |
| GET    | `/api/v1/organizer/nib/{filename}`    | Get organizer NIB document            | -                             | -                                                       |
| GET    | `/api/v1/organizer/`                  | List all organizers (public)          | -                             | -                                                       |
| GET    | `/api/v1/organizer/public/:slug`      | Get organizer profile by slug (public)| -                             | -                                                       |
| POST   | `/api/v1/organizer/`                  | Create organizer profile              | JWT + Organizer role          | `{"name": "string", "npwp": file, "nib": file, ...}`   |
| GET    | `/api/v1/organizer/profiles`          | Get all organizer profiles (admin)    | JWT + Admin only              | -                                                       |
| GET    | `/api/v1/organizer/me`                | Get own organizer profile             | JWT                           | -                                                       |
| PUT    | `/api/v1/organizer/photo`             | Update organizer profile photo        | JWT + Organizer/Admin         | FormData: `photo(file)`                                 |
| PUT    | `/api/v1/organizer/`                  | Update organizer profile              | JWT + Organizer/Admin         | `{"name": "string", ...}`                               |
| PUT    | `/api/v1/organizer/approve/:id`       | Approve organizer profile             | JWT + Admin only              | -                                                       |
| PUT    | `/api/v1/organizer/reject/:id`        | Reject organizer profile              | JWT + Admin only              | -                                                       |
| DELETE | `/api/v1/organizer/close-account`     | Close organizer account               | JWT + Organizer only          | -                                                       |
| GET    | `/api/v1/organizer/:id`               | Get organizer profile by ID           | JWT                           | -                                                       |

---

## 🔄 Organizer Profile Update Staging

*Route prefix: `/api/v1/staging-organizer/*`* — Admin Only (Protected by JWT)

| Method | Endpoint                                        | Description                             |
|--------|-------------------------------------------------|-----------------------------------------|
| GET    | `/api/v1/staging-organizer/npwp/{filename}`     | Get staged NPWP document                |
| GET    | `/api/v1/staging-organizer/nib/{filename}`      | Get staged NIB document                 |
| PUT    | `/api/v1/staging-organizer/:id`                 | Validate / approve a profile update     |
| GET    | `/api/v1/staging-organizer/update/:id`          | Find staged update by profile ID        |
| GET    | `/api/v1/staging-organizer/updates`             | List all staged profile updates         |
| GET    | `/api/v1/staging-organizer/:id`                 | Find staged update by ID                |

---

## 📂 Category Management

*Route prefix: `/api/v1/category/*`*

| Method | Endpoint                       | Description                       | Middleware       | Request Body                                     |
|--------|--------------------------------|-----------------------------------|------------------|--------------------------------------------------|
| GET    | `/api/v1/category/`            | List all categories (public)      | -                | -                                                |
| GET    | `/api/v1/category/filter`      | Filter categories by name (public)| -                | -                                                |
| POST   | `/api/v1/category/`            | Create category                   | JWT + Admin only | `{"name": "string"}`                             |
| GET    | `/api/v1/category/list`        | List categories with pagination   | JWT + Admin only | -                                                |
| PUT    | `/api/v1/category/:id`         | Update category                   | JWT + Admin only | `{"name": "string"}`                             |
| DELETE | `/api/v1/category/:id`         | Delete category                   | JWT + Admin only | -                                                |

---

## 🎪 Event Management

*Route prefix: `/api/v1/event/*`*

| Method | Endpoint                              | Description                              | Middleware                         | Request Body / Notes                                |
|--------|---------------------------------------|------------------------------------------|------------------------------------|-----------------------------------------------------|
| GET    | `/api/v1/event/banner/{filename}`     | Get event banner image (public)          | -                                  | -                                                   |
| GET    | `/api/v1/event/active`                | List all active events (public)          | -                                  | -                                                   |
| GET    | `/api/v1/event/nearest`               | Find nearest events (public)             | -                                  | -                                                   |
| GET    | `/api/v1/event/route/:id`             | Get event route/directions (rate-limited)| -                                  | -                                                   |
| GET    | `/api/v1/event/view/:slug`            | Get event detail by slug (public)        | -                                  | -                                                   |
| GET    | `/api/v1/event/public/:id`            | Get events by organizer ID (public)      | -                                  | -                                                   |
| POST   | `/api/v1/event/`                      | Create event                             | JWT + Organizer only               | `{"title": "string", "banner": file, "categories": ["uuid"], ...}` |
| GET    | `/api/v1/event/`                      | List all events (admin view)             | JWT + Admin only                   | -                                                   |
| GET    | `/api/v1/event/organizer-events`      | List events by organizer                 | JWT + Organizer/Admin              | -                                                   |
| PUT    | `/api/v1/event/review/:id`            | Review / approve event                   | JWT + Admin only                   | `{"status": "approved\|rejected", "note": "string"}` |
| PUT    | `/api/v1/event/cancel/:id`            | Cancel event                             | JWT + Organizer only               | -                                                   |
| GET    | `/api/v1/event/:id`                   | Get event by ID                          | JWT + Organizer/Admin              | -                                                   |
| PUT    | `/api/v1/event/:id`                   | Update event                             | JWT + Organizer only               | `{"title": "string", ...}`                          |
| DELETE | `/api/v1/event/:id`                   | Delete event                             | JWT + Organizer only               | -                                                   |

---

## 🔄 Updated Event Management

*Route prefix: `/api/v1/updated-event/*`* — Admin Only (Protected by JWT)

| Method | Endpoint                                          | Description                                 |
|--------|---------------------------------------------------|---------------------------------------------|
| GET    | `/api/v1/updated-event/banner/{filename}`         | Get updated event banner (static)           |
| GET    | `/api/v1/updated-event/`                          | List all pending event updates              |
| PUT    | `/api/v1/updated-event/review/:id`                | Review / validate an event update           |
| GET    | `/api/v1/updated-event/update-list/:event_id`     | List all updates for a specific event       |
| GET    | `/api/v1/updated-event/:event_id`                 | Get latest update for a specific event      |
| PUT    | `/api/v1/updated-event/:id`                       | Cancel an event update (Organizer only)     |

---

## 📍 Location

| Method | Endpoint                  | Description                            | Query Parameters     |
|--------|---------------------------|----------------------------------------|----------------------|
| GET    | `/api/v1/location/`       | Search location (for event creation)   | `q=string`           |

---

## 🔧 Project Structure

```
ngevent/
├── cmd/
│   └── app/              # Application entrypoint (main.go)
│   └── worker/           # worker entrypoint (main.go)
├── internal/
│   ├── config/           # App configuration (env loading)
        ├── database.go   # DB connection
        ├── redis.go      # redis connection
        ├── asynq.go      # Background job queue (Asynq) configuration
│   ├── dto/              # Data Transfer Objects
│   ├── handler/          # HTTP request handlers
│   ├── migrations/       # Database migration files
│   ├── model/            # GORM models / domain entities
│   ├── repository/       # Data access layer
│   ├── server/           # Fiber server setup & route registration
│   ├── service/          # Business logic layer
│   ├── tasks/            # Asynq background task definitions
│   └── utils/            # Helpers (middleware, JWT, etc.)
├── frontend/             # React + Vite + TypeScript frontend
├── storage/              # Uploaded files (profiles, NPWP, NIB, banners)
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── .env.example
```

---

## 📝 Notes

- API uses **JWT** for authentication. Include the token as a cookie or `Authorization: Bearer <token>` header.
- Rate limiting is applied globally (100 req/min) and on sensitive auth routes (10 req/min).
- Static files (photos, documents, banners) are served directly from the API under their respective routes.
- Background jobs (e.g., sending emails) are processed by a separate **worker** service using Redis + Asynq.

---

> ⚠️ **This project is still under active development. Features and APIs are subject to change.**
