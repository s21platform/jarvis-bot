-- +goose Up
CREATE TABLE IF NOT EXISTS config (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS config;