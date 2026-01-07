CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "public"."role" AS ENUM ('user', 'admin', 'event organizer');

CREATE TABLE "public"."users"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"role" role NOT NULL,
    "is_verified" boolean DEFAULT false,
	"created_at" timestamp DEFAULT NOW(),
	"updated_at" timestamp DEFAULT NOW(),
	CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "public"."sessions"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"user_id" uuid NOT NULL,
	"refresh_token" varchar(255) NOT NULL,
	"ip_address" varchar(255) NOT NULL,
	"user_agent" varchar(255) NOT NULL,
	"created_at" timestamp DEFAULT NOW(),
	"updated_at" timestamp DEFAULT NOW(),
	CONSTRAINT "session_pkey" PRIMARY KEY ("id"),
	CONSTRAINT "fk_sessions_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);