ALTER TABLE url ADD COLUMN IF NOT EXISTS is_deleted boolean NOT NULL default false;
