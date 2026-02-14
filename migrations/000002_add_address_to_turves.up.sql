-- Add address column to turves (safe for existing rows: use DEFAULT so no NULLs)
ALTER TABLE turves ADD COLUMN IF NOT EXISTS address VARCHAR(255) NOT NULL DEFAULT '';
