-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    users
ADD
    CONSTRAINT unique_username UNIQUE ("username");

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    users DROP CONSTRAINT unique_username;

-- +goose StatementEnd