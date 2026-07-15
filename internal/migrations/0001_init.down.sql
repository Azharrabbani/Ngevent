-- Drop indexes yang eksplisit dibuat terpisah (opsional, karena akan otomatis
-- terhapus saat tabelnya di-drop, tapi ditulis eksplisit untuk kejelasan/simetri)
DROP INDEX IF EXISTS "public"."idx_events_active_coordinates";
DROP INDEX IF EXISTS "public"."event_update_categories_event_update_id";
DROP INDEX IF EXISTS "public"."event_update_tickets_event_update_id";
DROP INDEX IF EXISTS "public"."event_update_status";
DROP INDEX IF EXISTS "public"."event_updates_start_date";
DROP INDEX IF EXISTS "public"."event_update_event_id";
DROP INDEX IF EXISTS "public"."idx_events_coordinates";
DROP INDEX IF EXISTS "public"."idx_category_id";
DROP INDEX IF EXISTS "public"."idx_event_categories_deleted_at";
DROP INDEX IF EXISTS "public"."idx_events_start_date";
DROP INDEX IF EXISTS "public"."idx_events_deleted_at";
DROP INDEX IF EXISTS "public"."idx_event_id";
DROP INDEX IF EXISTS "public"."idx_event_city";
DROP INDEX IF EXISTS "public"."idx_event_status";
DROP INDEX IF EXISTS "public"."idx_organizer_profiles_updates_slug";
DROP INDEX IF EXISTS "public"."idx_organizer_profiles_slug";
DROP INDEX IF EXISTS "public"."idx_event_slug";
DROP INDEX IF EXISTS "public"."idx_categories_slug";
DROP INDEX IF EXISTS "public"."idx_eo_profile_id";
DROP INDEX IF EXISTS "public"."idx_otp_user_id";
DROP INDEX IF EXISTS "public"."idx_sessions_user_id";
DROP INDEX IF EXISTS "public"."idx_organizer_profiles_updates_deleted_at";
DROP INDEX IF EXISTS "public"."idx_organizer_profiles_deleted_at";
DROP INDEX IF EXISTS "public"."idx_users_deleted_at";
DROP INDEX IF EXISTS "public"."idx_users_role";

DROP INDEX IF EXISTS "public"."unique_users_email";
DROP INDEX IF EXISTS "public"."unique_ticket_type";
DROP INDEX IF EXISTS "public"."unique_event_updates";
DROP INDEX IF EXISTS "public"."unique_category";
DROP INDEX IF EXISTS "public"."unique_category_event";
DROP INDEX IF EXISTS "public"."unique_pending_update";
DROP INDEX IF EXISTS "public"."unique_approved_nib";
DROP INDEX IF EXISTS "public"."unique_approved_npwp";

-- Drop tables (urutan: child dulu, baru parent, sesuai dependency FK)
DROP TABLE IF EXISTS "public"."event_update_categories";
DROP TABLE IF EXISTS "public"."event_update_tickets";
DROP TABLE IF EXISTS "public"."event_updates";
DROP TABLE IF EXISTS "public"."event_categories";
DROP TABLE IF EXISTS "public"."tickets";
DROP TABLE IF EXISTS "public"."events";
DROP TABLE IF EXISTS "public"."categories";
DROP TABLE IF EXISTS "public"."organizer_profiles_updates";
DROP TABLE IF EXISTS "public"."organizer_profiles";
DROP TABLE IF EXISTS "public"."attendee_profiles";
DROP TABLE IF EXISTS "public"."otp_verifications";
DROP TABLE IF EXISTS "public"."sessions";
DROP TABLE IF EXISTS "public"."users";

-- Drop custom enum types
DROP TYPE IF EXISTS "public"."ticket_type";
DROP TYPE IF EXISTS "public"."event_status";
DROP TYPE IF EXISTS "public"."organizer_profile_status";
DROP TYPE IF EXISTS "public"."type_verification";
DROP TYPE IF EXISTS "public"."role";

-- Drop extensions
DROP EXTENSION IF EXISTS "postgis";
DROP EXTENSION IF EXISTS "uuid-ossp";