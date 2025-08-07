-- +goose Up
CREATE TABLE IF NOT EXISTS staff_birthday
(
    id         SERIAL PRIMARY KEY,
    user_id    VARCHAR NOT NULL,
    nickname   VARCHAR NOT NULL,
    channel_id VARCHAR NOT NULL,
    day        INT CHECK (day BETWEEN 1 AND 31),
    month      INT CHECK (month BETWEEN 1 AND 12),
    year       INT,
    notes      TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id)
);


-- +goose Down
DROP TABLE IF EXISTS staff_birthday;