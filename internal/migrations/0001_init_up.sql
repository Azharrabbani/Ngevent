CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "public"."role" AS ENUM ('user', 'admin', 'event organizer');
CREATE TYPE "public"."type_verification" AS ENUM ('verified_email', 'reset_password');

SET TIME ZONE 'UTC';

CREATE TABLE "public"."users"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"role" role NOT NULL,
    "is_verified" boolean DEFAULT false,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "public"."sessions"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"user_id" uuid NOT NULL,
	"refresh_token" varchar(255) NOT NULL,
	"ip_address" varchar(255) NOT NULL,
	"user_agent" varchar(255) NOT NULL,
	"expired_at" TIMESTAMPTZ NOT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT "session_pkey" PRIMARY KEY ("id"),
	CONSTRAINT "fk_sessions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE
);

CREATE TABLE "public"."otp_verifications"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"user_id" uuid NOT NULL,
	"otp"	VARCHAR(255) NOT NULL,
	"is_used" BOOLEAN DEFAULT false,
	"expired_at" TIMESTAMPTZ NOT NULL,
	"type_verification" type_verification NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT "otp_pkey" PRIMARY KEY("id"),
	CONSTRAINT "fk_otp_users" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

CREATE TABLE "public"."attendee_profiles"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"user_id" uuid NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"username" VARCHAR(255),
	"photo_profile" TEXT,
	"phone_number" VARCHAR(100) NOT NULL,
	"country" VARCHAR(120) NOT NULL,
	"address" TEXT,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT "attendee_profiles_pk" PRIMARY KEY("id"),
	CONSTRAINT "fk_users_attendee_profiles" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);


CREATE INDEX "idx_users_role" ON "public"."users"("role");
CREATE INDEX "idx_sessions_user_id" ON "public"."sessions"("user_id");
CREATE INDEX "idx_otp_user_id" ON "public"."otp_verifications"("user_id");
