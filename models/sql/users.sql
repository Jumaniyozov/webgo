CREATE TABLE IF NOT EXISTS users
(
    id    SERIAL PRIMARY KEY,
    name  TEXT,
    email TEXT UNIQUE NOT NULL
);

