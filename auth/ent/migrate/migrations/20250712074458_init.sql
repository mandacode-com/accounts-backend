-- Create "auth_accounts" table
CREATE TABLE "public"."auth_accounts" (
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "provider" character varying NOT NULL,
  "provider_id" character varying NULL,
  "is_verified" boolean NOT NULL DEFAULT false,
  "email" character varying NOT NULL,
  "password_hash" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "authaccount_email_provider" to table: "auth_accounts"
CREATE UNIQUE INDEX "authaccount_email_provider" ON "public"."auth_accounts" ("email", "provider");
-- Create index "authaccount_provider_provider_id" to table: "auth_accounts"
CREATE UNIQUE INDEX "authaccount_provider_provider_id" ON "public"."auth_accounts" ("provider", "provider_id");
-- Create index "authaccount_user_id_provider" to table: "auth_accounts"
CREATE UNIQUE INDEX "authaccount_user_id_provider" ON "public"."auth_accounts" ("user_id", "provider");
