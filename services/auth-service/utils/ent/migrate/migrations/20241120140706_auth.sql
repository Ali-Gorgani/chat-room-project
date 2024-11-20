-- Create "auths" table
CREATE TABLE "auths" ("id" character varying NOT NULL, "user_id" bigint NOT NULL, "refresh_token" character varying NOT NULL, "is_revoked" boolean NOT NULL DEFAULT false, "created_at" timestamptz NOT NULL, "expires_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
