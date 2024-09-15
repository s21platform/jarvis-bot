-- +goose Up

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    service VARCHAR,
    task_type VARCHAR,
    task_title VARCHAR,
    assignee VARCHAR
);

-- +goose Down
DROP TABLE IF EXISTS tasks;