-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    command TEXT NOT NULL,
    payload TEXT,
    status TEXT DEFAULT 'pending',
    scheduled_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_tasks_scheduled_at ON tasks (scheduled_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tasks_scheduled_at;
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
