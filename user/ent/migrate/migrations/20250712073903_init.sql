-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "sync_code" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "is_archived" boolean NOT NULL DEFAULT false,
  "archived_at" timestamptz NULL,
  "delete_after" timestamptz NULL,
  PRIMARY KEY ("id")
);
