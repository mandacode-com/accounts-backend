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
-- Create "sent_emails" table
CREATE TABLE "public"."sent_emails" (
  "id" uuid NOT NULL,
  "email" character varying NOT NULL,
  "sent_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "sent_emails_users_sent_emails" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
