-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    "id" SERIAL  PRIMARY KEY,
    "telegram_id" BIGINT UNIQUE NOT NULL,
    "updated_at" TIMESTAMP DEFAULT current_timestamp,
    "created_at" TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS characters (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER REFERENCES users("id") ON DELETE CASCADE UNIQUE,
    "name" VARCHAR(255) NOT NULL,
    "coins" DOUBLE PRECISION DEFAULT 0,
    "max_energy" SMALLINT DEFAULT 5,
    "experience" DOUBLE PRECISION DEFAULT 0,
    "updated_at" TIMESTAMP DEFAULT current_timestamp,
    "created_at" TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS inventory_items (
    "id" SERIAL PRIMARY KEY,
    "character_id" INTEGER REFERENCES characters("id") ON DELETE CASCADE NOT NULL,
    "item_id" VARCHAR(255) NOT NULL,
    "quantity" SMALLINT NOT NULL,
    "updated_at" TIMESTAMP DEFAULT current_timestamp,
    "created_at" TIMESTAMP DEFAULT current_timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS inventory_items;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
