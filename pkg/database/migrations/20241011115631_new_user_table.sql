-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    "id" SERIAL PRIMARY KEY,
    "telegramId" integer unique,
    "name" varchar(255) not null,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd