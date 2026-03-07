CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "public"."role" AS ENUM ('user', 'admin', 'event organizer');
CREATE TYPE "public"."type_verification" AS ENUM ('verified_email', 'reset_password');
CREATE TYPE "public"."organizer_profile_status" AS ENUM ('pending', 'approved', 'rejected');

SET TIME ZONE 'UTC';

CREATE TABLE "public"."users"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"email" varchar(255) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"role" role NOT NULL,
    "is_verified" boolean DEFAULT false,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"deleted_at" TIMESTAMPTZ NULL,
	CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "public"."sessions"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"user_id" uuid NOT NULL,
	"jti" varchar(255) NOT NULL,
	"refresh_token" text NOT NULL,
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

CREATE TABLE "public"."organizer_profiles"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"user_id" uuid NOT NULL UNIQUE,
	"status" organizer_profile_status NOT NULL DEFAULT 'pending',
	"rejected_reason" TEXT,
	"reviewed_by" uuid,
	"reviewed_at" TIMESTAMPTZ,
	"name" VARCHAR(255) NOT NULL,
	"photo_profile" TEXT,
	"email" VARCHAR(255),
	"instagram" VARCHAR(255),
	"phone_number" VARCHAR(100) NOT NULL,
	"country" VARCHAR(120) NOT NULL,
	"address" TEXT,
	"description" TEXT,
	"npwp_number" VARCHAR(100) NOT NULL,
	"npwp_document" TEXT NOT NULL,
	"nib_number" VARCHAR(100) NOT NULL,
	"nib_document" TEXT NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT "organizer_profiles_pk" PRIMARY KEY("id"),
	CONSTRAINT "fk_users_organizer_profiles" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

CREATE TABLE "public"."organizer_profiles_updates"(
	"id" uuid DEFAULT uuid_generate_v4() NOT NULL,
	"profile_id" uuid NOT NULL,
	"status" organizer_profile_status NOT NULL DEFAULT 'pending',
	"name" VARCHAR(255) NOT NULL,
	"phone_number" VARCHAR(100) NOT NULL,
	"country" VARCHAR(120) NOT NULL,
	"npwp_number" VARCHAR(100) NOT NULL,
	"npwp_document" TEXT NOT NULL,
	"nib_number" VARCHAR(100) NOT NULL,
	"nib_document" TEXT NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT "organizer_profiles_updates_pk" PRIMARY KEY("id"),
	CONSTRAINT "fk_organizer_profiles_fk" FOREIGN KEY ("profile_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE
);

CREATE TABLE "public"."categories"(
	"id" SERIAL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"slug" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Create Unique Index
CREATE UNIQUE INDEX "unique_approved_npwp"
ON "public"."organizer_profiles"("npwp_number")
WHERE status = 'approved';

CREATE UNIQUE INDEX "unique_approved_nib"
ON "public"."organizer_profiles"("nib_number")
WHERE status = 'approved';

CREATE UNIQUE INDEX unique_pending_update
ON "public"."organizer_profile_updates" ("organizer_profile_id")
WHERE status = 'pending';


-- Crete index
CREATE INDEX "idx_users_role" ON "public"."users"("role");
CREATE INDEX "idx_users_deleted_at" ON "public"."users"("deleted_at");
CREATE INDEX "idx_sessions_user_id" ON "public"."sessions"("user_id");
CREATE INDEX "idx_otp_user_id" ON "public"."otp_verifications"("user_id");
CREATE INDEX "idx_eo_profile_id" ON "public"."organizer_profiles_updates"("profile_id");
CREATE INDEX "idx_categories_slug" ON "public"."categories"("slug");
