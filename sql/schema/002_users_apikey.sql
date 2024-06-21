-- +goose Up
ALTER TABLE users
ADD api_key VARCHAR(64) DEFAULT encode(sha256(random()::text::bytea), 'hex') UNIQUE NOT NULL;

-- +goose Down
ALTER TABLE users
DROP api_key;