-- modify "expenses" table
ALTER TABLE "expenses" ALTER COLUMN "create_time" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "update_time" SET DEFAULT CURRENT_TIMESTAMP;
-- modify "personal_settings" table
ALTER TABLE "personal_settings" ALTER COLUMN "create_time" SET DEFAULT CURRENT_TIMESTAMP, ALTER COLUMN "update_time" SET DEFAULT CURRENT_TIMESTAMP;
