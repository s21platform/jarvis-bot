-- +goose Up
CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    thread_id VARCHAR NOT NULL,
    channel VARCHAR NOT NULL,
    message_id VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(thread_id)
);


-- +goose Down
DROP TABLE IF EXISTS news; 