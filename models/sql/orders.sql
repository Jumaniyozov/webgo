CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    user_Id     INT NOT NULL,
    amount      INT,
    description TEXT
);