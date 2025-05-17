-- Create "meetings" table
CREATE TABLE "public"."meetings" (
  "id" bigserial NOT NULL,
  "host_user_id" text NOT NULL,
  "title" text NULL,
  "scheduled_start" timestamptz NULL,
  "scheduled_end" timestamptz NULL,
  "is_recurring" boolean NULL DEFAULT false,
  "meeting_code" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "duration_minutes" text NOT NULL,
  "status" text NOT NULL DEFAULT 'in-progress',
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_meetings_meeting_code" UNIQUE ("meeting_code")
);
-- Create "participants" table
CREATE TABLE "public"."participants" (
  "id" bigserial NOT NULL,
  "meeting_code" character varying NOT NULL,
  "user_id" character varying NOT NULL,
  "role" text NULL DEFAULT 'participant',
  "joined_at" timestamptz NULL,
  "left_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_participants_meeting_code" to table: "participants"
CREATE INDEX "idx_participants_meeting_code" ON "public"."participants" ("meeting_code");
-- Create index "idx_participants_user_id" to table: "participants"
CREATE INDEX "idx_participants_user_id" ON "public"."participants" ("user_id");
