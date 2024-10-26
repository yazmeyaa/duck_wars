-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    users RENAME COLUMN "name" TO "username";

ALTER TABLE
    users
ADD
    COLUMN "password" varchar(128) NOT NULL;

ALTER TABLE
    users
ADD
    COLUMN "gmail" varchar(128) DEFAULT NULL;

ALTER TABLE
    users
ALTER COLUMN
    "telegramId"
SET
    DEFAULT NULL;

ALTER TABLE
    users
ADD
    CONSTRAINT unique_telegramId UNIQUE ("telegramId");

ALTER TABLE
    users
ADD
    CONSTRAINT unique_gmail UNIQUE ("gmail");

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    users RENAME COLUMN "username" TO "name";

ALTER TABLE
    users DROP COLUMN "password";

ALTER TABLE
    users DROP COLUMN "gmail";

ALTER TABLE
    users
ALTER COLUMN
    "telegramId" DROP DEFAULT;

ALTER TABLE
    users DROP CONSTRAINT unique_telegramId;

ALTER TABLE
    users DROP CONSTRAINT unique_gmail;

-- +goose StatementEnd