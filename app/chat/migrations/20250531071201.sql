-- Create "meetings" table
CREATE TABLE "public"."meetings" (
  "id" bigserial NOT NULL,
  "meeting_id" text NOT NULL,
  "host_user_id" text NOT NULL,
  "title" text NULL,
  "scheduled_start" timestamptz NULL,
  "scheduled_end" timestamptz NULL,
  "is_recurring" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "duration_minutes" bigint NOT NULL,
  "status" text NOT NULL DEFAULT 'in-progress',
  PRIMARY KEY ("meeting_id")
);
-- Create "participants" table
CREATE TABLE "public"."participants" (
  "id" bigserial NOT NULL,
  "meeting_id" character varying NOT NULL,
  "user_id" character varying NOT NULL,
  "role" text NULL DEFAULT 'participant',
  "joined_at" timestamptz NULL,
  "left_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_meeting_user" to table: "participants"
CREATE UNIQUE INDEX "idx_meeting_user" ON "public"."participants" ("meeting_id", "user_id");
