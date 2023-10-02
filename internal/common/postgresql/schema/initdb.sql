CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Users
(
    id             UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username       VARCHAR(20)  UNIQUE      NOT NULL,
    password_hash  VARCHAR(256)             NOT NULL,
    salt           VARCHAR(16)              NOT NULL,
	planned_budget numeric(10, 2),
    avatar_url     TEXT DEFAULT '/static/img/img1.png'
);

CREATE TABLE Accounts (
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID REFERENCES Users(id),
    balance numeric(10, 2),
    mean_payment TEXT
);

--=============================================================================

ALTER TABLE Users
ALTER COLUMN planned_budget SET DEFAULT 0.0;

--=============================================================================

INSERT INTO "users"(username, password_hash, planned_budget, avatar_url)
VALUES ('kosmatoff', 'hash', 10000, 'image/img1.png');


INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 533, 'Кошелек');

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 599, 'Наличка');