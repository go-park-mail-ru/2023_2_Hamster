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

CREATE TABLE SUB_Category (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    category_id UUID REFERENCES Category(id),
    name VARCHAR(15) UNIQUE NOT NULL
);

CREATE TABLE Transaction (
	id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id      UUID REFERENCES Users(id),
    account_income UUID REFERENCES Accounts(id),
    account_outcome UUID REFERENCES Accounts(id),
	income       numeric(10, 2),
    outcome      numeric(10, 2),
	date         timestamp DEFAULT now(),
	payer        VARCHAR(20),
	description  VARCHAR(100)
);

CREATE TABLE TransactionCategory (
    transaction_id UUID REFERENCES Transaction(id),
    category_id UUID REFERENCES Category(id),
    PRIMARY KEY (transaction_id, category_id)
);
--=============================================================================

ALTER TABLE Users
ALTER COLUMN planned_budget SET DEFAULT 0.0;

--=============================================================================

INSERT INTO "users"(login, username, password_hash, planned_budget)
VALUES ('kossmatof','komarov', '$argon2id$v=19$m=65536,t=1,p=4$m8qhM3XLae+RCTGirBFEww$Znu5RBnxlam2xRoVtwBzbdSrN4/sRCm1IMOVX4N2uxw', 10000);

INSERT INTO "users"(login, username, password_hash, planned_budget)
VALUES ('test','test1', '$argon2id$v=19$m=65536,t=1,p=4$m8qhM3XLae+RCTGirBFEww$Znu5RBnxlam2xRoVtwBzbdSrN4/sRCm1IMOVX4N2uxw', 10000);

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 1000, 'Кошелек');

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 1000, 'Наличка');

INSERT INTO "category"(user_id, name)
VALUES ((SELECT id FROM Users limit 1), 'ЖКХ');

INSERT INTO "category"(user_id, name)
VALUES ((SELECT id FROM Users limit 1), 'Стипендия');
