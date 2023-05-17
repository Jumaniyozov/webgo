-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    email         TEXT UNIQUE  NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
