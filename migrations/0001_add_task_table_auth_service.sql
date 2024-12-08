-- +goose Up

CREATE TABLE IF NOT EXISTS credentials (
    id SERIAL PRIMARY KEY,
    service VARCHAR,
    data  JSONB
);

-- +goose Down
DROP TABLE IF EXISTS credentials;