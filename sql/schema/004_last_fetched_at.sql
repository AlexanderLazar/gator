-- +goose up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose down

ALTER TABLE IF EXISTS feeds DROP COLUMN IF EXISTS last_fetched_at;
