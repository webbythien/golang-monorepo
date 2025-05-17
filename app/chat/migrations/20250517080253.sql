-- Modify "meetings" tableBEGIN;

ALTER TABLE public.meetings
    ALTER COLUMN duration_minutes
    TYPE bigint
    USING NULLIF(trim(duration_minutes), '')::bigint;

COMMIT;
