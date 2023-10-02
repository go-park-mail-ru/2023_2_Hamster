CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Users
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username      VARCHAR(20)  UNIQUE      NOT NULL,
    password_hash VARCHAR(256)             NOT NULL,
	planned_budget numeric(10,2),
    avatar_url    TEXT DEFAULT '/static/img/img1.png'
);

CREATE TABLE Accounts (
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID REFERENCES Users(id),
    balance MONEY,
    mean_payment TEXT
);

--=============================================================================

INSERT INTO "users"(username, password_hash, first_name, last_name, planned_budget, avatar_url)
VALUES ('kosmatoff', 'hash','Дмитрий', 'Комаров', 1000000, 'image/img1.png');


INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 533, 'Кошелек');

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 599, 'Наличка');