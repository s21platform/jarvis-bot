-- +goose Up
ALTER TABLE tasks
    ADD COLUMN task_description TEXT;

-- +goose Down
ALTER TABLE tasks
DROP COLUMN task_description;