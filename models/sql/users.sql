CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    email         TEXT UNIQUE  NOT NULL
);

