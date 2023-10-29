CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Users
(
    id             UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username       VARCHAR(20)              NOT NULL,
    login          VARCHAR(20)       UNIQUE NOT NULL,
    password_hash  VARCHAR(256)             NOT NULL,
	planned_budget numeric(10, 2),
    avatar_url     UUID
);

CREATE TABLE Accounts (
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID REFERENCES Users(id),
    balance numeric(10, 2),
    mean_payment TEXT
);

CREATE TABLE Category (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES Users(id),
    name VARCHAR(15) UNIQUE NOT NULL
);


CREATE TABLE Transaction (
	id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id      UUID REFERENCES Users(id),
	category_id  UUID REFERENCES Category(id),
	account_id   UUID REFERENCES Accounts(id),
	total        numeric(10, 2),
	is_income    bool,
	date         DATE DEFAULT now(),
	payer        VARCHAR(20),
	description  VARCHAR(100)
);

--=============================================================================

ALTER TABLE Users
ALTER COLUMN planned_budget SET DEFAULT 0.0;

--=============================================================================

INSERT INTO "users"(login, username, password_hash, salt, planned_budget)
VALUES ('kossmatof','komarov', 'hash', 'fdsf', 10000);

INSERT INTO "users"(login, username, password_hash, salt, planned_budget)
VALUES ('test','test1', 'hash', 'fdsf', 10000);

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 533, 'Кошелек');

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 599, 'Наличка');

INSERT INTO "category"(user_id, name)
VALUES ((SELECT id FROM Users limit 1), 'ЖКХ');

INSERT INTO "category"(user_id, name)
VALUES ((SELECT id FROM Users limit 1), 'Стипендия');

INSERT INTO "transaction" (user_id, category_id, account_id, total, is_income, date, payer, description)
VALUES ((SELECT id FROM Users limit 1), (SELECT id FROM Category limit 1), (SELECT id FROM Accounts limit 1), 12400.50, false, '2023-10-01', 'МОСЖКХ', 'Оплата недвижки');

INSERT INTO "transaction"(user_id, category_id, account_id, total, is_income, date, payer, description)
VALUES ((SELECT id FROM Users limit 1), (SELECT id FROM Category limit 1), (SELECT id FROM Accounts limit 1), 12450.50, false, '2023-10-02',  'МОСЖКХ', 'Оплата недвижки');

INSERT INTO "transaction"(user_id, category_id, account_id, total, is_income, date, payer, description)
VALUES ((SELECT id FROM Users limit 1), (SELECT id FROM Category limit 1), (SELECT id FROM Accounts limit 1), 12450.50, false, '2023-10-02',  'МОСЖКХ', 'Оплата недвижки');

INSERT INTO "transaction"(user_id, category_id, account_id, total, is_income, date, payer, description)
VALUES ((SELECT id FROM Users limit 1), (SELECT id FROM Category limit 1), (SELECT id FROM Accounts limit 1), 40000, true, '2023-10-02', 'VK', 'Зарплата');

INSERT INTO "transaction"(user_id, category_id, account_id, total, is_income, date, payer, description)
VALUES ((SELECT id FROM Users limit 1), (SELECT id FROM Category limit 1), (SELECT id FROM Accounts limit 1), 40000, true, '2023-09-02', 'VK', 'Зарплата');
