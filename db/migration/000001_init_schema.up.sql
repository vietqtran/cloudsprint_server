CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "user_id" uuid NOT NULL,
  "status" int NOT NULL DEFAULT 1,
  "reset_password_token" varchar NULL,
  "reset_password_token_expires_at" timestamptz NULL,
  "login_failed_attempts" int NOT NULL DEFAULT 0,
  "first_failed_login_at" timestamptz NULL,
  "verify_email_token" varchar NULL,
  "email_verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL CHECK (LENGTH("first_name") > 1) CHECK (LENGTH("first_name") < 256),
  "last_name" varchar NOT NULL CHECK (LENGTH("last_name") > 1) CHECK (LENGTH("last_name") < 256),
  "status" int NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "account_id" uuid NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "status" int NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "password_resets" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "token" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "used" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "email_otps" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "otp" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "used" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS "users_email_idx" ON "users" ("email");
CREATE INDEX IF NOT EXISTS "accounts_email_idx" ON "accounts" ("email");
CREATE INDEX IF NOT EXISTS "accounts_user_id_idx" ON "accounts" ("user_id");
CREATE INDEX IF NOT EXISTS "sessions_account_id_idx" ON "sessions" ("account_id");
CREATE INDEX IF NOT EXISTS "password_resets_token_idx" ON "password_resets" ("token");
CREATE INDEX IF NOT EXISTS "password_resets_email_idx" ON "password_resets" ("email"); 
CREATE INDEX IF NOT EXISTS "email_otps_email_idx" ON "email_otps" ("email");
CREATE INDEX IF NOT EXISTS "email_otps_otp_idx" ON "email_otps" ("otp"); 