CREATE TABLE IF NOT EXISTS Users
(
    id            SERIAL       PRIMARY KEY,
    username      VARCHAR(20)  UNIQUE      NOT NULL,
    password_hash VARCHAR(256)             NOT NULL,
    first_name    VARCHAR(20)              NOT NULL,
    last_name     VARCHAR(20)              NOT NULL,
    avatar_url    TEXT
);

CREATE TABLE IF NOT EXISTS Accounts (
    id UUID PRIMARY KEY,
    user_id INT,
    balance NUMERIC,
    mean_payment TEXT
);
